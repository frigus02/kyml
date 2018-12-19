package fs

import (
	"bytes"
	"io/ioutil"
	"os"
)

type fakeFilesystem struct {
	files     map[string][]byte
	fileModes map[string]os.FileMode
}

// NewFakeFilesystem creates a new filesystem using the OS.
func NewFakeFilesystem() Filesystem {
	return &fakeFilesystem{
		files:     make(map[string][]byte),
		fileModes: make(map[string]os.FileMode),
	}
}

// NewFakeFilesystemFromDisk creates a new filesystem prefilled with the
// specified files, which are read from disk.
func NewFakeFilesystemFromDisk(files ...string) (Filesystem, error) {
	fs := NewFakeFilesystem()
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		fs.WriteFile(file, data, 0644)
	}

	return fs, nil
}

func (fs *fakeFilesystem) Open(name string) (File, error) {
	if data, ok := fs.files[name]; ok {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		return &fakeFile{bytes.NewBuffer(dataCopy), false}, nil
	}

	return nil, os.ErrNotExist
}

func (fs *fakeFilesystem) ReadFile(filename string) ([]byte, error) {
	if data, ok := fs.files[filename]; ok {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		return dataCopy, nil
	}

	return nil, os.ErrNotExist
}

func (fs *fakeFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fs.files[filename] = data
	fs.fileModes[filename] = perm
	return nil
}
