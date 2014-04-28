package ioutils

import (
	"io"
)

type SizeReader struct {
	io.Reader
	size int64
}

func NewSizeReader(reader io.Reader) *SizeReader {
	return &SizeReader{reader, 0}
}

func (r *SizeReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.size += int64(n)

	return
}

func (r *SizeReader) Size() (size int64) {
	return r.size
}
