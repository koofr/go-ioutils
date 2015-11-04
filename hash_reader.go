package ioutils

import (
	"encoding/hex"
	"github.com/koofr/go-cryptoutils/bettermd5"
	"io"
)

type HashReader struct {
	io.Reader
	digest *bettermd5.BetterDigest
}

func NewHashReader(reader io.Reader) *HashReader {
	return &HashReader{reader, bettermd5.New()}
}

func NewHashReaderFromState(reader io.Reader, state []byte) *HashReader {
	return &HashReader{reader, bettermd5.NewFromState(state)}
}

func (r *HashReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	_, _ = r.digest.Write(p[0:n])

	return
}

func (r *HashReader) GetState() []byte {
	return r.digest.GetState()
}

func (r *HashReader) Hash() (hash string) {
	sum := r.digest.Sum(nil)
	return hex.EncodeToString(sum)
}
