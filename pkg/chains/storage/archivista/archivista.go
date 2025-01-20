package archivista

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	archivistaClient "github.com/in-toto/archivista/pkg/http-client"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/tektoncd/chains/pkg/chains/objects"
	"github.com/tektoncd/chains/pkg/config"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1" // if needed
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/logging"
)

const (
	StorageBackendArchivista = "archivista"
)

// generatePublicKeyIDFunc is a package-level variable wrapping the public key ID generation.
// It allows tests to simulate errors.
var generatePublicKeyIDFunc = cryptoutil.GeneratePublicKeyID

// buildEnvelope constructs a DSSE envelope from the raw payload, signature, keyID, and certificate chain.
// If a valid chain is provided, it parses it into a leaf and intermediates; otherwise, certificate data is omitted.
func buildEnvelope(rawPayload []byte, signature, keyID string, chain string) dsse.Envelope {
	var leaf []byte
	var inters [][]byte

	chain = strings.TrimSpace(chain)
	if chain != "" {
		var err error
		leaf, inters, err = parseAndOrderCertificateChain(chain)
		if err != nil {
			// Log error if needed and fall back to no certificate data.
			leaf = nil
			inters = [][]byte{}
		}
	}
	return dsse.Envelope{
		Payload:     rawPayload,
		PayloadType: "application/vnd.in-toto+json",
		Signatures: []dsse.Signature{
			{
				KeyID:         keyID,
				Signature:     []byte(signature),
				Certificate:   leaf,
				Intermediates: inters,
			},
		},
	}
}

// Backend is the interface that all storage backends must implement.
type Backend interface {
	StorePayload(ctx context.Context, obj objects.TektonObject, rawPayload []byte, signature string, opts config.StorageOpts) error
	RetrievePayloads(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string]string, error)
	RetrieveSignatures(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string][]string, error)
	Type() string
}

// ArchivistaStorage implements the Backend interface for Archivista.
type ArchivistaStorage struct {
	client       *archivistaClient.ArchivistaClient
	url          string
	cfg          config.ArchivistaStorageConfig
	tektonClient tektonclient.Interface // Injected Tekton client for patching objects
}

// NewArchivistaStorage initializes a new ArchivistaStorage backend.
func NewArchivistaStorage(cfg config.Config, tektonClient tektonclient.Interface) (*ArchivistaStorage, error) {
	archCfg := cfg.Storage.Archivista
	if strings.TrimSpace(archCfg.URL) == "" {
		return nil, fmt.Errorf("missing archivista URL in storage configuration")
	}

	client, err := archivistaClient.CreateArchivistaClient(&http.Client{}, archCfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Archivista client: %w", err)
	}

	return &ArchivistaStorage{
		client:       client,
		url:          archCfg.URL,
		cfg:          archCfg,
		tektonClient: tektonClient,
	}, nil
}

// patchTektonObjectAnnotations patches the Tekton object's annotations with the given key/value pairs
// in one single patch call.
func PatchTektonObjectAnnotations(ctx context.Context, obj objects.TektonObject, annotations map[string]string, tektonClient tektonclient.Interface) error {
	patchData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": annotations,
		},
	}
	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("failed to marshal patch data: %w", err)
	}

	switch o := obj.GetObject().(type) {
	case *tektonv1.TaskRun:
		_, err = tektonClient.TektonV1().TaskRuns(o.Namespace).Patch(ctx, o.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	case *tektonv1.PipelineRun:
		_, err = tektonClient.TektonV1().PipelineRuns(o.Namespace).Patch(ctx, o.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	case *v1beta1.TaskRun:
		_, err = tektonClient.TektonV1beta1().TaskRuns(o.Namespace).Patch(ctx, o.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	case *v1beta1.PipelineRun:
		_, err = tektonClient.TektonV1beta1().PipelineRuns(o.Namespace).Patch(ctx, o.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
		return err
	default:
		return fmt.Errorf("unsupported Tekton object type for patching")
	}
}

// StorePayload builds a DSSE envelope from the raw payload and signature,
// logs the envelope, uploads it via the Archivista client API, and patches the
// Tekton object with the returned gitoid and Archivista URL.
func (a *ArchivistaStorage) StorePayload(ctx context.Context, obj objects.TektonObject, rawPayload []byte, signature string, opts config.StorageOpts) error {
	logger := logging.FromContext(ctx)

	// Validate signature.
	if strings.TrimSpace(signature) == "" {
		return fmt.Errorf("missing signature")
	}

	var keyID string
	certPEM := strings.TrimSpace(opts.Cert)
	if certPEM != "" {
		block, _ := pem.Decode([]byte(certPEM))
		if block != nil {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				// Generate keyID from the public key.
				keyID, err = generatePublicKeyIDFunc(cert.PublicKey, crypto.SHA256)
				if err != nil {
					logger.Errorw("Failed to generate KeyID", "error", err)
					keyID = ""
				}
			} else {
				logger.Errorw("Failed to parse certificate", "error", err)
			}
		} else {
			logger.Error("Failed to decode certificate PEM")
		}
	} // if no certificate provided, keyID remains blank

	// Optionally decode the payload for logging.
	decodedPayload, err := base64.StdEncoding.DecodeString(string(rawPayload))
	if err != nil {
		logger.Errorw("Failed to base64 decode payload", "keyID", keyID, "error", err)
		logger.Infof("Raw payload (not base64 decoded): %s", string(rawPayload))
	} else {
		logger.Infof("Decoded payload: %s", string(decodedPayload))
	}

	env := buildEnvelope(rawPayload, signature, keyID, opts.Chain)

	// Upload the envelope using the Archivista client's Store method.
	uploadResp, err := a.client.Store(ctx, env)
	if err != nil {
		logger.Errorw("Failed to upload DSSE envelope to Archivista", "error", err)
		return fmt.Errorf("failed to upload envelope to Archivista: %w", err)
	}
	logger.Infof("Successfully uploaded DSSE envelope to Archivista, response: %+v", uploadResp)

	// Update the in-memory Tekton object with Archivista annotations.
	annotations := map[string]string{
		"chains.tekton.dev/archivista-gitoid": uploadResp.Gitoid,
		"chains.tekton.dev/archivista-url":    a.url,
	}
	obj.SetAnnotations(annotations)

	// Patch the live Tekton object in one call.
	if err := PatchTektonObjectAnnotations(ctx, obj, annotations, a.tektonClient); err != nil {
		logger.Errorw("Failed to patch Tekton object with Archivista annotations", "error", err)
		return fmt.Errorf("failed to patch Tekton object: %w", err)
	}

	return nil
}

// RetrievePayload is not implemented for Archivista.
func (a *ArchivistaStorage) RetrievePayload(ctx context.Context, key string) ([]byte, []byte, error) {
	return nil, nil, fmt.Errorf("RetrievePayload not implemented for Archivista")
}

// RetrievePayloads is not implemented for Archivista.
func (a *ArchivistaStorage) RetrievePayloads(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string]string, error) {
	return nil, fmt.Errorf("RetrievePayloads not implemented for Archivista")
}

// RetrieveSignatures is not implemented for Archivista.
func (a *ArchivistaStorage) RetrieveSignatures(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string][]string, error) {
	return nil, fmt.Errorf("RetrieveSignatures not implemented for Archivista")
}

// Type returns the storage backend type.
func (a *ArchivistaStorage) Type() string {
	return StorageBackendArchivista
}
