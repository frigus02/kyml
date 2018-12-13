package fs

import (
	"bytes"
	"os"
)

type fakeFile struct {
	data     *bytes.Buffer
	isClosed bool
}

func (f *fakeFile) Read(p []byte) (n int, err error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	return f.data.Read(p)
}

func (f *fakeFile) Close() error {
	if f.isClosed {
		return os.ErrClosed
	}

	f.isClosed = true
	return nil
}
