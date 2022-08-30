package ioutils

import "io"

type ProgressReportingReader struct {
	r           io.Reader
	minChange   int64
	onChange    func(bytesRead int64)
	currentRead int64
}

func NewProgressReportingReader(r io.Reader, minChange int64, onChange func(bytesRead int64)) *ProgressReportingReader {
	return &ProgressReportingReader{
		r:           r,
		minChange:   minChange,
		onChange:    onChange,
		currentRead: 0,
	}
}

func (r *ProgressReportingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)

	r.currentRead += int64(n)

	if r.currentRead > r.minChange || err == io.EOF {
		r.onChange(r.currentRead)
		r.currentRead = 0
	}

	return n, err
}
