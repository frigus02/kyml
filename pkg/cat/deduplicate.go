package cat

import (
	"github.com/frigus02/kyml/pkg/k8syaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func addOrReplaceExistingDocs(existingDocs, newDocs []*unstructured.Unstructured) []*unstructured.Unstructured {
	for _, doc := range newDocs {
		docGVK := doc.GroupVersionKind()
		found := false
		for i, seenDoc := range existingDocs {
			if k8syaml.GVKEquals(docGVK, seenDoc.GroupVersionKind()) &&
				doc.GetNamespace() == seenDoc.GetNamespace() &&
				doc.GetName() == seenDoc.GetName() {
				existingDocs[i] = doc
				found = true
				break
			}
		}

		if !found {
			existingDocs = append(existingDocs, doc)
		}
	}

	return existingDocs
}
