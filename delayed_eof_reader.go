package ioutils

import "io"

// DelayedEOFReader fixes Read behaviour when n > 0 and err == io.EOF.
// Read usually returns (1, nil) and (0, io.EOF) but it can also just end
// with (1, io.EOF). DelayedEOFReader fixes the second behavior to behave like
// the first one.
type DelayedEOFReader struct {
	r     io.Reader
	isEOF bool
}

func NewDelayedEOFReader(r io.Reader) *DelayedEOFReader {
	return &DelayedEOFReader{
		r:     r,
		isEOF: false,
	}
}

func (r *DelayedEOFReader) Read(p []byte) (int, error) {
	if r.isEOF {
		return 0, io.EOF
	}
	n, err := r.r.Read(p)
	if err == io.EOF {
		r.isEOF = true
		if n > 0 {
			return n, nil
		}
	}
	return n, err
}
