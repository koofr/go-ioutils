package ioutils

import (
	"io"
)

type ChanCloser struct {
	io.Reader
	ch chan bool
}

func (c *ChanCloser) Close() error {
	c.ch <- true
	return nil
}

func NewChanCloser(r io.Reader, ch chan bool) *ChanCloser {
	return &ChanCloser{r, ch}
}
