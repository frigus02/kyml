package fs

import "os"

type osFile struct {
	file *os.File
}

func (f osFile) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f osFile) Close() error {
	return f.file.Close()
}
