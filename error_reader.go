package ioutils

type ErrorReader struct {
	Error error
}

func NewErrorReader(err error) *ErrorReader {
	return &ErrorReader{err}
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.Error
}
