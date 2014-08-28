package ioutils

import (
	"io"
)

type rechunkingReader struct {
	buffer    []byte
	eof       bool
	fetchMore func() ([]byte, error)
}

func NewRechunkingReader(fetchMore func() ([]byte, error)) io.Reader {
	return &rechunkingReader{[]byte{}, false, fetchMore}
}

func (r *rechunkingReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	if r.eof {
		err = io.EOF
		return
	}

	if len(r.buffer) == 0 && !r.eof {
		var newBuf []byte
		newBuf, err = r.fetchMore()
		if err == io.EOF {
			r.eof = true
		}
		r.buffer = newBuf
	}

	n = copy(p, r.buffer)
	if n == 0 && err == nil {
		err = io.EOF
	}
	r.buffer = r.buffer[n:]
	return
}

type splitReader struct {
	rechunkingReader
	ch     chan splitReaderChunk
	close  chan struct{}
	closed bool
}

type splitReaderChunk struct {
	buf []byte
	err error
}

func (r *splitReader) Close() error {
	r.close <- struct{}{}
	return nil
}

func SplitReader(r io.ReadCloser, bufSize uint64, n uint) []io.ReadCloser {
	readers := make([]*splitReader, n)
	readClosers := make([]io.ReadCloser, n)
	for i := uint(0); i < n; i++ {
		ch := make(chan splitReaderChunk, 0)
		rc := rechunkingReader{[]byte{}, false, nil}
		r := &splitReader{rc, ch, make(chan struct{}, 1), false}
		func(r *splitReader) {
			r.fetchMore = func() (bytes []byte, err error) {
				if r.closed {
					err = io.EOF
					return
				}

				chunk, ok := <-ch
				if !ok {
					err = io.EOF
					return
				}
				bytes = chunk.buf
				err = chunk.err
				return
			}
		}(r)
		readers[i] = r
		readClosers[i] = r
	}
	go pumpSplitReader(r, bufSize, readers)
	return readClosers
}

func pumpSplitReader(source io.ReadCloser, bufSize uint64, readers []*splitReader) {
	bufs := [][]byte{make([]byte, bufSize), make([]byte, bufSize)}
	for i := 0; ; i = (i + 1) % 2 {
		buf := bufs[i]
		total := 0

		n, err := source.Read(buf)

		for _, rd := range readers {
			if rd.closed {
				continue
			}
			total += 1

			select {
			case rd.ch <- splitReaderChunk{buf[0:n], err}:
			case <-rd.close:
				rd.closed = true
				close(rd.ch)
			}
		}

		if total == 0 || err == io.EOF {
			break
		}
	}
	source.Close()
	for _, rd := range readers {
		if !rd.closed {
			close(rd.ch)
		}
	}
}
