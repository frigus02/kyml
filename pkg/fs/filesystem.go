package fs

import "os"

// Filesystem is an interface abstraction for os.Open.
type Filesystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}
