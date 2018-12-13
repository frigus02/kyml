package diff

import "github.com/pmezard/go-difflib/difflib"

// Diff returns a diff between the two specified strings A and B.
func Diff(nameA string, a string, nameB string, b string) (string, error) {
	linesA := difflib.SplitLines(a)
	linesB := difflib.SplitLines(b)

	return difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        linesA,
		FromFile: nameA,
		B:        linesB,
		ToFile:   nameB,
		Context:  0,
	})
}
