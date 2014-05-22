package ioutils

import (
	"os"
)

type FileRemoveReader struct {
	*os.File
}

func NewFileRemoveReader(file *os.File) *FileRemoveReader {
	return &FileRemoveReader{file}
}

func (r *FileRemoveReader) Close() error {
	defer os.Remove(r.File.Name())
	return r.File.Close()
}
