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
var supportedKinds = []struct {
	GroupVersionKind schema.GroupVersionKind
	PathToPodSpec    []string
}{
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "ReplicaSet"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "batch", Version: "v1beta1", Kind: "CronJob"},
		PathToPodSpec:    []string{"spec", "jobTemplate", "spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
	{
		GroupVersionKind: schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "ReplicationController"},
		PathToPodSpec:    []string{"spec", "template", "spec"},
	},
}

func getPathToPodSpec(gvk schema.GroupVersionKind) []string {
	for _, kind := range supportedKinds {
		if k8syaml.GVKEquals(gvk, kind.GroupVersionKind) {
			return kind.PathToPodSpec
		}
	}

	return nil
}
