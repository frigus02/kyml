package fs

import (
	"os"
	"time"
)

type fakeFileInfo struct {
	name string
	size int64
	mode os.FileMode
}

func (fi *fakeFileInfo) Name() string {
	return fi.name
}

func (fi *fakeFileInfo) Size() int64 {
	return fi.size
}

func (fi *fakeFileInfo) Mode() os.FileMode {
	return fi.mode
}

func (fi *fakeFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *fakeFileInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

func (fi *fakeFileInfo) Sys() interface{} {
	return nil
}
