package ioutils

import "io"

type FileSpan struct {
	Start, End int64
}

func ApplyFileSpan(r io.ReadSeeker, span *FileSpan) (io.Reader, error) {
	if span == nil {
		return r, nil
	}

	if _, err := r.Seek(span.Start, io.SeekStart); err != nil {
		return nil, err
	}

	length := span.End - span.Start + 1

	return io.LimitReader(r, length), nil
}
