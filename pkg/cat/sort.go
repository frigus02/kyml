package cat

import (
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var gvkOrder = []schema.GroupVersionKind{
	// Most resources require a namespace. A namespace has no requirements.
	{Group: "", Version: "v1", Kind: "Namespace"},

	// Custom resources require the definition.
	{Group: "apiextensions.k8s.io", Version: "v1beta1", Kind: "CustomResourceDefinition"},

	// StorageClasses can be configured as default, so that PVCs can use them
	// without an explicit reference.
	{Group: "storage.k8s.io", Version: "v1", Kind: "StorageClass"},

	// Creation of a service account fails if a secret referenced in
	// imagePullSecrets does not exist.
	{Group: "", Version: "v1", Kind: "ConfigMap"},
	{Group: "", Version: "v1", Kind: "Secret"},

	// Creation of pods fail if the service account referenced in
	// serviceAccountName does not exist. Role bindings require the referenced
	// service account and role.
	{Group: "", Version: "v1", Kind: "ServiceAccount"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"},

	// It’s best to specify the service first, since that will ensure the
	// scheduler can spread the pods associated with the service as they are
	// created by the controller(s), such as Deployment.
	// https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/
	{Group: "", Version: "v1", Kind: "Service"},

	// Several resources require pods, e.g. HorizontalPodAutoscaler. These
	// resources will create pods.
	{Group: "apps", Version: "v1", Kind: "DaemonSet"},
	{Group: "apps", Version: "v1", Kind: "Deployment"},
	{Group: "apps", Version: "v1", Kind: "ReplicaSet"},
	{Group: "apps", Version: "v1", Kind: "StatefulSet"},
	{Group: "batch", Version: "v1", Kind: "Job"},
	{Group: "batch", Version: "v1beta1", Kind: "CronJob"},
	{Group: "", Version: "v1", Kind: "ReplicationController"},
}
var gvkOrderMap = func() map[string]int {
	m := map[string]int{}
	for i, n := range gvkOrder {
		m[n.String()] = i
	}
	return m
}()

func sortDocs(docs []*unstructured.Unstructured) {
	sort.SliceStable(docs, func(i, j int) bool {
		indexI, foundI := gvkOrderMap[docs[i].GroupVersionKind().String()]
		indexJ, foundJ := gvkOrderMap[docs[j].GroupVersionKind().String()]
		if foundI && foundJ {
			return indexI < indexJ
		}
		if foundI && !foundJ {
			return true
		}
		return false
	})
}
