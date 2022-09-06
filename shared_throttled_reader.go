package ioutils

import (
	"context"
	"io"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const sharedBurstLimit = 1000 * 1000 * 1000

type SharedLimiter struct {
	limiter      *rate.Limiter
	limiterMutex sync.RWMutex
}

func NewSharedLimiter(bytesPerSec int64) *SharedLimiter {
	var limiter *rate.Limiter

	if bytesPerSec > 0 {
		limiter = rate.NewLimiter(rate.Limit(bytesPerSec), sharedBurstLimit)
		limiter.AllowN(time.Now(), sharedBurstLimit) // spend initial burst
	}

	return &SharedLimiter{
		limiter: limiter,
	}
}

func (l *SharedLimiter) WaitN(ctx context.Context, n int) (err error) {
	l.limiterMutex.RLock()
	limiter := l.limiter
	l.limiterMutex.RUnlock()

	if limiter == nil {
		return nil
	}

	return limiter.WaitN(ctx, n)
}

func (l *SharedLimiter) SetLimit(newBytesPerSecond int64) {
	l.limiterMutex.Lock()
	defer l.limiterMutex.Unlock()

	if newBytesPerSecond == 0 {
		l.limiter = nil
	} else {
		if l.limiter == nil {
			l.limiter = rate.NewLimiter(rate.Limit(newBytesPerSecond), sharedBurstLimit)
			l.limiter.AllowN(time.Now(), sharedBurstLimit) // spend initial burst
		} else {
			l.limiter.SetLimit(rate.Limit(newBytesPerSecond))
		}
	}
}

type SharedThrottledReaderLimiter interface {
	WaitN(ctx context.Context, n int) error
}

type SharedThrottledReader struct {
	ctx     context.Context
	r       io.ReadCloser
	limiter SharedThrottledReaderLimiter
}

func NewSharedThrottledReader(ctx context.Context, r io.ReadCloser, limiter SharedThrottledReaderLimiter) *SharedThrottledReader {
	return &SharedThrottledReader{
		ctx:     ctx,
		r:       r,
		limiter: limiter,
	}
}

func (r *SharedThrottledReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}

	if err := r.limiter.WaitN(r.ctx, n); err != nil {
		return n, err
	}

	return n, nil
}

func (r *SharedThrottledReader) Close() (err error) {
	return r.r.Close()
}
