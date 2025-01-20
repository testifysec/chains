package archivista

import (
	"context"
	"testing"

	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	fakePipelineClient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tektoncd/chains/pkg/chains/objects"
)

// --- Supported Branch Tests ---

// TestPatchTektonObjectAnnotations_TaskRunV1Beta1 exercises the v1beta1.TaskRun branch.
func TestPatchTektonObjectAnnotations_TaskRunV1Beta1(t *testing.T) {
	ctx := context.Background()
	// Seed the fake client with a v1beta1.TaskRun.
	tr := &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "taskrun-v1beta1",
			Namespace: "default",
		},
	}
	client := fakePipelineClient.NewSimpleClientset(tr)
	// Wrap the TaskRun using the helper.
	obj := objects.NewTaskRunObjectV1Beta1(tr)
	annotations := map[string]string{"foo": "bar"}
	if err := patchTektonObjectAnnotations(ctx, obj, annotations, client); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(client.Actions()) == 0 {
		t.Fatalf("expected a patch action, got none")
	}
}

// TestPatchTektonObjectAnnotations_PipelineRunV1Beta1 exercises the v1beta1.PipelineRun branch.
func TestPatchTektonObjectAnnotations_PipelineRunV1Beta1(t *testing.T) {
	ctx := context.Background()
	// Seed the fake client with a v1beta1.PipelineRun.
	pr := &v1beta1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipelinerun-v1beta1",
			Namespace: "default",
		},
	}
	client := fakePipelineClient.NewSimpleClientset(pr)
	// Wrap the PipelineRun.
	obj := objects.NewPipelineRunObjectV1Beta1(pr)
	annotations := map[string]string{"foo": "bar"}
	if err := patchTektonObjectAnnotations(ctx, obj, annotations, client); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(client.Actions()) == 0 {
		t.Fatalf("expected a patch action, got none")
	}
}

// --- Unsupported Branch Test ---

// unsupportedTaskRun is a type that embeds a v1beta1.TaskRun but is distinct.
type unsupportedTaskRun struct {
	v1beta1.TaskRun
}

// unsupportedTektonObject wraps a real TektonObject (from a v1beta1.TaskRun) but
// overrides GetObject() so that it returns an unsupported type.
type unsupportedTektonObject struct {
	inner objects.TektonObject
}

// GetObject returns an unsupported type by wrapping the inner object's underlying TaskRun.
func (u *unsupportedTektonObject) GetObject() interface{} {
	// Retrieve the inner object (which should be a *v1beta1.TaskRun).
	tr, ok := u.inner.GetObject().(*v1beta1.TaskRun)
	if !ok {
		return u.inner.GetObject()
	}
	// Wrap it in an unsupportedTaskRun so the type switch in PatchTektonObjectAnnotations doesn't match.
	return &unsupportedTaskRun{*tr}
}
