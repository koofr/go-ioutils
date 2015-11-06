package ioutils

import (
	"io"
	"sync"
)

type LazyReader struct {
	builder func() (io.Reader, error)
	r       io.Reader
	err     error
	once    sync.Once
}

func (lr *LazyReader) Read(p []byte) (n int, err error) {
	lr.once.Do(func() {
		r, err := lr.builder()

		if err != nil {
			lr.err = err
			return
		}

		lr.r = r
	})

	if lr.err != nil {
		return 0, lr.err
	}

	return lr.r.Read(p)
}

func (lr *LazyReader) Close() error {
	if lr.r != nil {
		if c, ok := lr.r.(io.Closer); ok {
			return c.Close()
		}

		return nil
	}

	return nil
}

func NewLazyReader(builder func() (io.Reader, error)) io.Reader {
	lr := &LazyReader{}
	lr.builder = builder
	return lr
}
