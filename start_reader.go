package ioutils

import (
	"io"
)

type StartReader struct {
	io.Reader
	OnStart func()
	started bool
}

func NewStartReader(reader io.Reader, onStart func()) *StartReader {
	return &StartReader{
		Reader:  reader,
		OnStart: onStart,
		started: false,
	}
}

func (r *StartReader) Read(p []byte) (n int, err error) {
	if !r.started {
		r.started = true
		r.OnStart()
	}

	return r.Reader.Read(p)
}
