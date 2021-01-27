package ioutils

import (
	"errors"
	"io"
)

var ErrMaxSizeExceeded = errors.New("max size exceeded")

type SizeLimitedReader struct {
	r    io.ReadCloser
	left int64
}

func NewSizeLimitedReader(r io.ReadCloser, limit int64) io.ReadCloser {
	return &SizeLimitedReader{
		r:    r,
		left: limit,
	}
}

func (r *SizeLimitedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)

	if int64(n) > r.left {
		return n, ErrMaxSizeExceeded
	}

	r.left -= int64(n)

	return n, err
}

func (r *SizeLimitedReader) Close() (err error) {
	return r.r.Close()
}
