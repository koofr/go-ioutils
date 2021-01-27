package ioutils

import (
	"bytes"
	"errors"
	"io"
)

var ErrDirtyReader = errors.New("dirty reader")

type CheckEmptyReader struct {
	r      io.Reader
	closer io.Closer

	isDirty bool
}

func NewCheckEmptyReader(r io.ReadCloser) *CheckEmptyReader {
	return &CheckEmptyReader{
		r:      r,
		closer: r,

		isDirty: false,
	}
}

func (r *CheckEmptyReader) IsEmpty() (isEmpty bool, err error) {
	if r.isDirty {
		return false, ErrDirtyReader
	}
	r.isDirty = true

	buf := []byte{0}

	n, err := r.r.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	if n == 0 {
		return true, nil
	}

	r.r = io.MultiReader(bytes.NewBuffer(buf), r.r)

	return false, nil
}

func (r *CheckEmptyReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *CheckEmptyReader) Close() (err error) {
	return r.closer.Close()
}
