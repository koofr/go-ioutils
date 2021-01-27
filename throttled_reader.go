package ioutils

import (
	"context"
	"io"
	"time"

	"golang.org/x/time/rate"
)

const burstLimit = 1000 * 1000 * 1000

type ThrottledReader struct {
	r       io.ReadCloser
	limiter *rate.Limiter
	ctx     context.Context
}

func NewThrottledReader(ctx context.Context, r io.ReadCloser, bytesPerSec int64) *ThrottledReader {
	limiter := rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	limiter.AllowN(time.Now(), burstLimit) // spend initial burst

	return &ThrottledReader{
		r:       r,
		limiter: limiter,
		ctx:     ctx,
	}
}

func (r *ThrottledReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}

	if err := r.limiter.WaitN(r.ctx, n); err != nil {
		return n, err
	}

	return n, nil
}

func (r *ThrottledReader) Close() (err error) {
	return r.r.Close()
}
