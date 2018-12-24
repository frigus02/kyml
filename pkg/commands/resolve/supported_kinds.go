package resolve

import (
	"github.com/frigus02/kyml/pkg/k8syaml"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// We only want to resolve images mentioned in the `image` property of
// containers. These currently only appear in PodSpec, which is under the path
// spec.template.spec in the listed resource kinds.
//
// See: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#container-v1-core
var supportedKinds = []schema.GroupVersionKind{
	schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"},
	schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
	schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"},
	schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
	schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
	schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "ReplicationController"},
}

func isSupportedKind(gvk schema.GroupVersionKind) bool {
	for _, supportedGvk := range supportedKinds {
		if k8syaml.GVKEquals(gvk, supportedGvk) {
			return true
		}
	}

	return false
}
