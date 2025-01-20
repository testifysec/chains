package archivista

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	archivistaClient "github.com/in-toto/archivista/pkg/http-client"
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

type Backend struct {
	client       *archivistaClient.ArchivistaClient
	url          string
	cfg          config.ArchivistaStorageConfig
	tektonClient tektonclient.Interface // Injected Tekton client for patching objects
}

// NewStorageBackend returns a new Archivista StorageBackend that can store Payloaders that are signed
// and wrapped in a DSSE envelope
func NewStorageBackend(cfg config.Config, tektonClient tektonclient.Interface) (*Backend, error) {
	archCfg := cfg.Storage.Archivista
	if strings.TrimSpace(archCfg.URL) == "" {
		return nil, fmt.Errorf("missing archivista URL in storage configuration")
	}

	client, err := archivistaClient.CreateArchivistaClient(&http.Client{}, archCfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Archivista client: %w", err)
	}

	return &Backend{
		client:       client,
		url:          archCfg.URL,
		cfg:          archCfg,
		tektonClient: tektonClient,
	}, nil
}

func patchTektonObjectAnnotations(ctx context.Context, obj objects.TektonObject, annotations map[string]string, tektonClient tektonclient.Interface) error {
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

// StorePayload attempts to parse `signature` as a DSSE envelope, and if successful
// sends it to an Archivista server for storage. Annotations pointing to the stored
// payload will be added to the Tekton object
func (a *Backend) StorePayload(ctx context.Context, obj objects.TektonObject, rawPayload []byte, signature string, opts config.StorageOpts) error {
	logger := logging.FromContext(ctx)
	var env dsse.Envelope
	if err := json.Unmarshal([]byte(signature), &env); err != nil {
		logger.Errorf("Failed to parse DSSE envelope: %w", err)
		return errors.Join(errors.New("Failed to parse DSSE envelope"), err)
	}

	uploadResp, err := a.client.Store(ctx, env)
	if err != nil {
		logger.Errorw("Failed to upload DSSE envelope to Archivista", "error", err)
		return err
	}
	logger.Infof("Successfully uploaded DSSE envelope to Archivista, response: %+v", uploadResp)

	// Update the Tekton object with Archivista annotations.
	annotations := map[string]string{
		"chains.tekton.dev/archivista-gitoid": uploadResp.Gitoid,
		"chains.tekton.dev/archivista-url":    a.url,
	}
	obj.SetAnnotations(annotations)

	if err := patchTektonObjectAnnotations(ctx, obj, annotations, a.tektonClient); err != nil {
		logger.Errorw("Failed to patch Tekton object with Archivista annotations", "error", err)
		return fmt.Errorf("failed to patch Tekton object: %w", err)
	}

	return nil
}

// RetrievePayload is not implemented for Archivista.
func (a *Backend) RetrievePayload(ctx context.Context, key string) ([]byte, []byte, error) {
	return nil, nil, fmt.Errorf("RetrievePayload not implemented for Archivista")
}

// RetrievePayloads is not implemented for Archivista.
func (a *Backend) RetrievePayloads(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string]string, error) {
	return nil, fmt.Errorf("RetrievePayloads not implemented for Archivista")
}

// RetrieveSignatures is not implemented for Archivista.
func (a *Backend) RetrieveSignatures(ctx context.Context, obj objects.TektonObject, opts config.StorageOpts) (map[string][]string, error) {
	return nil, fmt.Errorf("RetrieveSignatures not implemented for Archivista")
}

func (a *Backend) Type() string {
	return StorageBackendArchivista
}
