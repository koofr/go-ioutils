package ioutils

import (
	"context"
	"io"
)

type CtxReader struct {
	r   io.Reader
	ctx context.Context
}

func NewCtxReader(ctx context.Context, r io.Reader) *CtxReader {
	return &CtxReader{
		r:   r,
		ctx: ctx,
	}
}

func (r *CtxReader) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
	}

	return r.r.Read(p)
}
