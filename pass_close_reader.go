package ioutils

import "io"

type PassCloseReader struct {
	r     io.Reader
	close func() error
}

func NewPassCloseReader(reader io.Reader, close func() error) *PassCloseReader {
	return &PassCloseReader{
		r:     reader,
		close: close,
	}
}

func (r *PassCloseReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *PassCloseReader) Close() error {
	return r.close()
}
