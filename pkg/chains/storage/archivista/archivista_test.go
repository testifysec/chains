package archivista

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	archivistaClient "github.com/in-toto/archivista/pkg/http-client"
	"github.com/in-toto/go-witness/dsse"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/chains/pkg/chains/objects"
	"github.com/tektoncd/chains/pkg/config"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	fakePipelineClient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// --------------------------
// Helper: setupEnv
// --------------------------

// setupEnv creates a fresh ArchivistaStorage test environment using a given TaskRun name.
func setupEnv(taskRunName string, cfg config.Config, archClient *archivistaClient.ArchivistaClient) (*ArchivistaStorage, objects.TektonObject, *fakePipelineClient.Clientset) {
	tr := &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskRunName,
			Namespace: "default",
		},
	}
	obj := objects.NewTaskRunObjectV1Beta1(tr)
	fakeClient := fakePipelineClient.NewSimpleClientset(tr)
	aStorage, err := NewArchivistaStorage(cfg, fakeClient)
	if err != nil {
		panic("failed to initialize ArchivistaStorage: " + err.Error())
	}
	// Override the Archivista client with the provided one.
	aStorage.client = archClient
	return aStorage, obj, fakeClient
}

// --------------------------
// StorePayload Tests
// --------------------------

// TestStorePayload_TaskRun tests the basic success path of StorePayload without certificate data.
func TestStorePayload_TaskRun(t *testing.T) {
	ctx := context.Background()

	// Create a v1beta1.TaskRun with minimal metadata and a dummy result.
	tr := &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-taskrun",
			Namespace: "default",
		},
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				TaskRunResults: []v1beta1.TaskRunResult{
					{
						Name:  "IMAGE_URL",
						Value: *v1beta1.NewStructuredValues("mockImage"),
					},
				},
			},
		},
	}
	fakeClient := fakePipelineClient.NewSimpleClientset(tr)

	// Set up an httptest server to simulate Archivista.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/upload" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"gitoid": "fake-gitoid"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	httpClient := &http.Client{}
	cfg := config.Config{
		Storage: config.StorageConfigs{
			Archivista: config.ArchivistaStorageConfig{
				URL: ts.URL,
			},
		},
	}
	archClient, err := archivistaClient.CreateArchivistaClient(httpClient, cfg.Storage.Archivista.URL)
	if err != nil {
		t.Fatalf("failed to create Archivista client: %v", err)
	}

	aStorage, obj, fakeClient := setupEnv("test-taskrun", cfg, archClient)

	// Prepare a valid payload.
	type mockPayload struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	payload := mockPayload{
		A: "foo",
		B: 3,
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err, "should marshal payload")
	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)
	signature := "test-signature"
	opts := config.StorageOpts{
		ShortKey: "mockpayload",
		Cert:     "",
		Chain:    "",
	}

	// Call StorePayload.
	err = aStorage.StorePayload(ctx, obj, []byte(encodedPayload), signature, opts)
	assert.NoError(t, err, "StorePayload should succeed")

	// Retrieve the updated TaskRun.
	updated, err := fakeClient.TektonV1beta1().TaskRuns("default").Get(ctx, "test-taskrun", metav1.GetOptions{})
	assert.NoError(t, err, "should retrieve updated TaskRun")
	assert.Equal(t, "fake-gitoid", updated.Annotations["chains.tekton.dev/archivista-gitoid"])
	assert.Equal(t, ts.URL, updated.Annotations["chains.tekton.dev/archivista-url"])
}

// TestStorePayload_ErrorCases exercises error branches in StorePayload.
func TestStorePayload_ErrorCases(t *testing.T) {
	ctx := context.Background()

	// Setup a common httptest server and configuration.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"gitoid": "fake-gitoid"}`))
	}))
	defer ts.Close()

	httpClient := &http.Client{}
	cfg := config.Config{
		Storage: config.StorageConfigs{
			Archivista: config.ArchivistaStorageConfig{
				URL: ts.URL,
			},
		},
	}
	archClient, err := archivistaClient.CreateArchivistaClient(httpClient, cfg.Storage.Archivista.URL)
	if err != nil {
		t.Fatalf("failed to create Archivista client: %v", err)
	}

	setup := func(name string) (*ArchivistaStorage, objects.TektonObject, *fakePipelineClient.Clientset) {
		return setupEnv(name, cfg, archClient)
	}

	t.Run("missing signature", func(t *testing.T) {
		aStorage, obj, _ := setup("missing-sig")
		err := aStorage.StorePayload(ctx, obj, []byte("dummy"), "", config.StorageOpts{})
		if err == nil || err.Error() != "missing signature" {
			t.Errorf("expected missing signature error, got: %v", err)
		}
	})

	t.Run("invalid certificate PEM decode", func(t *testing.T) {
		aStorage, obj, fakeClient := setup("invalid-cert-decode")
		opts := config.StorageOpts{
			Cert: "invalid pem",
		}
		payload := base64.StdEncoding.EncodeToString([]byte("dummy"))
		err := aStorage.StorePayload(ctx, obj, []byte(payload), "sig", opts)
		if err != nil {
			t.Errorf("expected success even with invalid cert PEM, got error: %v", err)
		}
		updated, err := fakeClient.TektonV1beta1().TaskRuns("default").Get(ctx, "invalid-cert-decode", metav1.GetOptions{})
		if err != nil {
			t.Errorf("failed to get updated TaskRun: %v", err)
		}
		if updated.Annotations["chains.tekton.dev/archivista-gitoid"] != "fake-gitoid" {
			t.Errorf("unexpected gitoid: %s", updated.Annotations["chains.tekton.dev/archivista-gitoid"])
		}
	})

	t.Run("certificate PEM parse failure", func(t *testing.T) {
		aStorage, obj, fakeClient := setup("cert-parse-failure")
		opts := config.StorageOpts{
			Cert: "-----BEGIN CERTIFICATE-----\nnotbase64\n-----END CERTIFICATE-----",
		}
		payload := base64.StdEncoding.EncodeToString([]byte("dummy"))
		err := aStorage.StorePayload(ctx, obj, []byte(payload), "sig", opts)
		if err != nil {
			t.Errorf("expected success even if certificate fails parsing, got error: %v", err)
		}
		updated, err := fakeClient.TektonV1beta1().TaskRuns("default").Get(ctx, "cert-parse-failure", metav1.GetOptions{})
		if err != nil {
			t.Errorf("failed to get updated TaskRun: %v", err)
		}
		if updated.Annotations["chains.tekton.dev/archivista-gitoid"] != "fake-gitoid" {
			t.Errorf("unexpected gitoid: %s", updated.Annotations["chains.tekton.dev/archivista-gitoid"])
		}
	})

	t.Run("payload base64 decode error", func(t *testing.T) {
		aStorage, obj, fakeClient := setup("payload-decode-error")
		// Provide rawPayload that is not valid base64.
		err := aStorage.StorePayload(ctx, obj, []byte("not-base64"), "sig", config.StorageOpts{})
		if err != nil {
			t.Errorf("expected success even if base64 decode fails, got error: %v", err)
		}
		updated, err := fakeClient.TektonV1beta1().TaskRuns("default").Get(ctx, "payload-decode-error", metav1.GetOptions{})
		if err != nil {
			t.Errorf("failed to get updated TaskRun: %v", err)
		}
		if updated.Annotations["chains.tekton.dev/archivista-gitoid"] != "fake-gitoid" {
			t.Errorf("unexpected gitoid: %s", updated.Annotations["chains.tekton.dev/archivista-gitoid"])
		}
	})
}

// TestStorePayload_CertificateSuccess_WithRecordingServer tests the certificate branch.
// It uses an httptest.Server to record the outgoing DSSE envelope.
func TestStorePayload_CertificateSuccess_WithRecordingServer(t *testing.T) {
	ctx := context.Background()

	// Generate a valid certificate.
	validCertPEM, _, _ := createCertificate(t, "dummy", "dummy", 123, time.Now(), nil, nil)

	// Override keyID generation to return "test-key-id".
	origFunc := generatePublicKeyIDFunc
	generatePublicKeyIDFunc = func(pub interface{}, hash crypto.Hash) (string, error) {
		return "test-key-id", nil
	}
	defer func() { generatePublicKeyIDFunc = origFunc }()

	// Create an httptest.Server that records the DSSE envelope.
	var recordedBody []byte
	recServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		recordedBody = body
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"gitoid": "fake-gitoid"}`))
	}))
	defer recServer.Close()

	archClient, err := archivistaClient.CreateArchivistaClient(&http.Client{}, recServer.URL)
	if err != nil {
		t.Fatalf("failed to create Archivista client: %v", err)
	}

	// Create a minimal TaskRun and wrap it.
	tr := &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cert-success",
			Namespace: "default",
		},
	}
	tektonObj := objects.NewTaskRunObjectV1Beta1(tr)

	// Build configuration using the recording server's URL.
	cfg := config.Config{
		Storage: config.StorageConfigs{
			Archivista: config.ArchivistaStorageConfig{
				URL: recServer.URL,
			},
		},
	}

	aStorage, _, _ := setupEnv("cert-success", cfg, archClient)

	// Prepare a valid payload.
	payload := []byte("dummy payload")
	encodedPayload := base64.StdEncoding.EncodeToString(payload)
	signature := "dummy-signature"
	opts := config.StorageOpts{
		ShortKey: "dummy",
		Cert:     validCertPEM,
		Chain:    "",
	}

	// Call StorePayload.
	err = aStorage.StorePayload(ctx, tektonObj, []byte(encodedPayload), signature, opts)
	if err != nil {
		t.Fatalf("StorePayload failed: %v", err)
	}

	// Unmarshal the recorded DSSE envelope.
	var env dsse.Envelope
	if err := json.Unmarshal(recordedBody, &env); err != nil {
		t.Fatalf("failed to unmarshal recorded envelope: %v", err)
	}

	// Verify the signature's KeyID.
	if len(env.Signatures) == 0 {
		t.Fatal("expected at least one signature in envelope")
	}
	if env.Signatures[0].KeyID != "test-key-id" {
		t.Errorf("expected keyID 'test-key-id', got %q", env.Signatures[0].KeyID)
	}
}

func TestBuildEnvelope_FallbackOnInvalidChain(t *testing.T) {
	// Prepare inputs.
	rawPayload := []byte("dummy")
	signature := "dummy-sig"
	keyID := "dummy-key"
	// Provide a non-empty chain that cannot be parsed as valid certificates.
	invalidChain := "invalid chain"

	// Call buildEnvelope.
	env := buildEnvelope(rawPayload, signature, keyID, invalidChain)

	// Expect that no certificate data was included.
	if env.Signatures[0].Certificate != nil {
		t.Errorf("expected certificate to be nil, got %v", env.Signatures[0].Certificate)
	}
	if len(env.Signatures[0].Intermediates) != 0 {
		t.Errorf("expected intermediates to be empty, got %v", env.Signatures[0].Intermediates)
	}
}
