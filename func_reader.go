package ioutils

type FuncReader func(p []byte) (int, error)

func (r FuncReader) Read(p []byte) (n int, err error) {
	return r(p)
}
