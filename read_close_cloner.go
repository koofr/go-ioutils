package ioutils

import "io"

type ReadCloseCloner interface {
	io.ReadCloser
	Clone() (io.ReadCloser, error)
}
