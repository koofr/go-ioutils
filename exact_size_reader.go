package ioutils

import (
	"errors"
	"io"
)

var ErrExactSizeMismatch = errors.New("exact size mismatch")

type ExactSizeReader struct {
	r    io.ReadCloser
	left int64
}

func NewExactSizeReader(r io.ReadCloser, expectedSize int64) io.ReadCloser {
	return &ExactSizeReader{
		r:    r,
		left: expectedSize,
	}
}

func (r *ExactSizeReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)

	r.left -= int64(n)

	if r.left < 0 || (err == io.EOF && r.left > 0) {
		return 0, ErrExactSizeMismatch
	}
	if r.left == 0 && err != io.EOF {
		eofN, _ := r.r.Read([]byte{0})
		if eofN > 0 {
			return 0, ErrExactSizeMismatch
		}
	}

	return n, err
}

func (r *ExactSizeReader) Close() (err error) {
	return r.r.Close()
}
