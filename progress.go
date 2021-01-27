package ioutils

import (
	"io"

	"github.com/koofr/pb"
)

type ReaderAtSeeker interface {
	io.ReaderAt
	io.ReadSeeker
}

type ProgressReader struct {
	ReaderAtSeeker
	bar *pb.ProgressBar
}

func NewProgressReader(r ReaderAtSeeker) (pr *ProgressReader, err error) {
	len, err := r.Seek(0, 2)

	if err != nil {
		return
	}

	if _, err = r.Seek(0, 0); err != nil {
		return
	}

	bar := pb.New(0)
	bar.Total = len

	bar.Units = pb.U_BYTES
	bar.Start()

	pr = &ProgressReader{
		ReaderAtSeeker: r,
		bar:            bar,
	}

	return
}

func (pr *ProgressReader) Read(p []byte) (len int, err error) {
	defer pr.bar.Read(p)
	return pr.ReaderAtSeeker.Read(p)
}

func (pr *ProgressReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off == 0 {
		pr.bar.Set(0)
	}
	defer pr.bar.Read(p)
	return pr.ReaderAtSeeker.ReadAt(p, off)
}
