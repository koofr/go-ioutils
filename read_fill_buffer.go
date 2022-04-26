package ioutils

import "io"

// ReadFillBuffer is similar to io.ReadFull but it does not return
// ErrUnexpectedEOF
func ReadFillBuffer(r io.Reader, buf []byte) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	return n, err
}
