package k8syaml

import "k8s.io/apimachinery/pkg/runtime/schema"

// GVKEquals returns true if both specified struct have equal fields.
func GVKEquals(a, b schema.GroupVersionKind) bool {
	return a.Group == b.Group && a.Version == b.Version && a.Kind == b.Kind
}
