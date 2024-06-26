package ioutils

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

type HashReader struct {
	io.Reader
	digest hash.Hash
}

func NewHashReader(reader io.Reader) *HashReader {
	return &HashReader{reader, md5.New()}
}

func (r *HashReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	_, _ = r.digest.Write(p[0:n])

	return n, err
}

func (r *HashReader) Hash() (hash string) {
	sum := r.digest.Sum(nil)
	return hex.EncodeToString(sum)
}
