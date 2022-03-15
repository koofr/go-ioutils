package ioutils

// from https://stackoverflow.com/a/52153000/141214

import (
	"bytes"
	"io"
)

func CountLines(r io.Reader) (int64, error) {
	var count int64
	const lineBreak = '\n'

	buf := make([]byte, 64*1024)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var pos int

		for {
			i := bytes.IndexByte(buf[pos:], lineBreak)
			if i == -1 || bufferSize == pos {
				break
			}
			pos += i + 1
			count++
		}

		if err == io.EOF {
			break
		}
	}

	return count, nil
}
