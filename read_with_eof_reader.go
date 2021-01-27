package ioutils

import "io"

// ReadWithEOFReader is used for testing Read behaviour when n > 0 and err ==
// io.EOF. Read usually returns (1, nil) and (0, io.EOF) but it can also just
// end with (1, io.EOF).
type ReadWithEOFReader struct {
	len int
}

func NewReadWithEOFReaderLen(l int) *ReadWithEOFReader {
	return &ReadWithEOFReader{
		len: l,
	}
}

func NewReadWithEOFReader() *ReadWithEOFReader {
	return NewReadWithEOFReaderLen(0)
}

func (r *ReadWithEOFReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 1
	}
	if r.len > 0 {
		return r.len, io.EOF
	}
	return len(p), io.EOF
}
