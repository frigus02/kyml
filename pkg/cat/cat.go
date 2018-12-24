package cat

import (
	"io"

	"github.com/frigus02/kyml/pkg/fs"
	"github.com/frigus02/kyml/pkg/k8syaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Cat reads YAML documents from the specified files and prints them one after
// another in the specified writer. If a YAML document has the same apiVersion,
// kind, namespace and name as a previous one it replaces it in the output.
func Cat(out io.Writer, files []string, fs fs.Filesystem) error {
	var documents []*unstructured.Unstructured
	for _, filename := range files {
		file, err := fs.Open(filename)
		if err != nil {
			return err
		}

		docsInFile, err := k8syaml.Decode(file)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}

		documents = addOrReplaceExistingDocs(documents, docsInFile)
	}

	return k8syaml.Encode(out, documents)
}

// Stream reads YAML documents from the specified reader and prints them one
// after another in the specified writer. If a YAML document has the same
// apiVersion, kind, namespace and name as a previous one it replaces it in the
// output.
func Stream(out io.Writer, stream io.Reader) error {
	documents, err := StreamDecodeOnly(stream)
	if err != nil {
		return err
	}

	return k8syaml.Encode(out, documents)
}

// StreamDecodeOnly works like Stream, but returns a slice of unstructured
// objects instead of writing them to an output.
func StreamDecodeOnly(stream io.Reader) ([]*unstructured.Unstructured, error) {
	docsInStream, err := k8syaml.Decode(stream)
	if err != nil {
		return nil, err
	}

	var documents []*unstructured.Unstructured
	documents = addOrReplaceExistingDocs(documents, docsInStream)

	return documents, nil
}

func addOrReplaceExistingDocs(existingDocs, newDocs []*unstructured.Unstructured) []*unstructured.Unstructured {
	for _, doc := range newDocs {
		docGVK := doc.GroupVersionKind()
		found := false
		for i, seenDoc := range existingDocs {
			if k8syaml.GVKEquals(docGVK, seenDoc.GroupVersionKind()) &&
				doc.GetNamespace() == seenDoc.GetNamespace() &&
				doc.GetName() == doc.GetName() {
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
