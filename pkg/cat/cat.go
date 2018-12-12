package cat

import (
	"io"
	"os"

	"github.com/frigus02/kyml/pkg/k8syaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Cat prints YAML documents in the specified files and prints them one after
// another in the specified writer. If a YAML document has the same apiVersion,
// kind, namespace and name as a previous one it replaces it in the output.
func Cat(out io.Writer, files []string) error {
	var documents []unstructured.Unstructured
	for _, filename := range files {
		file, err := os.Open(filename)
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

		for _, doc := range docsInFile {
			found := false
			for i, seenDoc := range documents {
				if doc.GetAPIVersion() == seenDoc.GetAPIVersion() &&
					doc.GetKind() == seenDoc.GetKind() &&
					doc.GetNamespace() == seenDoc.GetNamespace() &&
					doc.GetName() == doc.GetName() {
					documents[i] = doc
					found = true
					break
				}
			}

			if !found {
				documents = append(documents, doc)
			}
		}
	}

	return k8syaml.Encode(out, documents)
}
