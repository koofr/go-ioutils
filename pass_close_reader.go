package ioutils

import "io"

type PassCloseReader struct {
	r     io.Reader
	close func() error
}

func (r *PassCloseReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *PassCloseReader) Close() error {
	return r.close()
}
