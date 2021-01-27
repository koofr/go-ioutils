package ioutils

type FuncWriter func(p []byte) (int, error)

func (r FuncWriter) Write(p []byte) (n int, err error) {
	return r(p)
}
