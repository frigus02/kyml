package fs

import (
	"io"
)

// File is an interface abstraction for os.File.
type File interface {
	io.ReadCloser
}
