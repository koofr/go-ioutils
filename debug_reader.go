package ioutils

import (
	"fmt"
	"io"
	"time"
)

type DebugReader struct {
	r               io.Reader
	lastReadStarted time.Time
	lastReadEnded   time.Time
	lastReadErr     error
	lastReadN       int
	totalN          int64
}

func NewDebugReader(r io.Reader) *DebugReader {
	return &DebugReader{
		r:               r,
		lastReadStarted: time.Time{},
		lastReadEnded:   time.Time{},
		lastReadErr:     nil,
		lastReadN:       0,
		totalN:          0,
	}
}

func (r *DebugReader) Read(p []byte) (n int, err error) {
	r.lastReadStarted = time.Now()
	n, err = r.r.Read(p)
	r.lastReadEnded = time.Now()
	r.lastReadErr = err
	r.lastReadN = n
	r.totalN += int64(n)
	return n, err
}

func (r *DebugReader) String() string {
	return fmt.Sprintf(
		"lastReadStarted=%s lastReadEnded=%s lastReadDuration=%d ms lastReadErr=%v lastReadN=%d totalN=%d",
		r.lastReadStarted.UTC().Format(time.RFC3339Nano),
		r.lastReadEnded.UTC().Format(time.RFC3339Nano),
		r.lastReadEnded.Sub(r.lastReadStarted)/time.Millisecond,
		r.lastReadErr,
		r.lastReadN,
		r.totalN,
	)
}
