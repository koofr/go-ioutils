package ioutils

import (
	"io"
)

type EofReader struct {
	io.Reader
	Eof bool
}

func NewEofReader(reader io.Reader) *EofReader {
	return &EofReader{
		Reader: reader,
		Eof:    false,
	}
}

func (r *EofReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	if err == io.EOF {
		r.Eof = true
	}

	return
}
