package ioutils

import "io"

type FuncReadCloseCloner struct {
	io.ReadCloser
	CloneFunc func() (io.ReadCloser, error)
}

func (r *FuncReadCloseCloner) Clone() (io.ReadCloser, error) {
	return r.CloneFunc()
}
