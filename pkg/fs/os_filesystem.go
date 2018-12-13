package fs

import (
	"io/ioutil"
	"os"
)

type osFilesystem struct{}

// NewOSFilesystem creates a new filesystem using the OS.
func NewOSFilesystem() Filesystem {
	return osFilesystem{}
}

func (fs osFilesystem) Open(name string) (File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return osFile{file}, nil
}

func (fs osFilesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs osFilesystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (fs osFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
