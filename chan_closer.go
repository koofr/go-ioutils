package ioutils

import (
	"io"
	"sync"
)

type ChanCloser struct {
	io.Reader
	ch          chan bool
	closed      bool
	closedMutex sync.RWMutex
}

func (c *ChanCloser) Close() error {
	c.closedMutex.RLock()
	closed := c.closed
	c.closedMutex.RUnlock()

	if !closed {
		c.closedMutex.Lock()
		c.closed = true
		c.closedMutex.Unlock()

		c.ch <- true
	}

	return nil
}

func NewChanCloser(r io.Reader, ch chan bool) *ChanCloser {
	return &ChanCloser{
		Reader: r,
		ch:     ch,
		closed: false,
	}
}
