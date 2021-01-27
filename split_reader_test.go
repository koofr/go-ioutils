package ioutils_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

func checkOneByte(r io.Reader, expected byte) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	Expect(n).To(Equal(1))
	Expect(buf[0]).To(Equal(expected))
	Expect(err).To(BeNil())
}

func checkEof(r io.Reader) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	Expect(n).To(Equal(0))
	Expect(err).To(Equal(io.EOF))
}

type simpleCloser struct {
	r      io.Reader
	closed bool
}

func (r *simpleCloser) Read(p []byte) (n int, err error) {
	if r.closed {
		err = io.EOF
		return
	}
	return r.r.Read(p)
}

func (r *simpleCloser) Close() error {
	r.closed = true
	return nil
}

var _ = Describe("RechunkingReader", func() {
	It("returns eof on empty input", func() {
		r := NewRechunkingReader(func() ([]byte, error) {
			return []byte{}, nil
		})
		checkEof(r)
	})

	It("rechunks big buffers", func() {
		done := false
		r := NewRechunkingReader(func() ([]byte, error) {
			if done {
				return []byte{}, nil
			} else {
				done = true
				return []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, nil
			}
		})

		buf := make([]byte, 4)

		n, err := r.Read(buf)
		Expect(n).To(Equal(4))
		Expect(buf[0]).To(Equal(byte(0)))
		Expect(err).To(BeNil())

		n, err = r.Read(buf)
		Expect(n).To(Equal(4))
		Expect(buf[0]).To(Equal(byte(4)))
		Expect(err).To(BeNil())

		n, err = r.Read(buf)
		Expect(n).To(Equal(2))
		Expect(buf[0]).To(Equal(byte(8)))
		Expect(err).To(BeNil())

		checkEof(r)
	})

	It("does partial reads for small buffers", func() {
		output := []byte{1, 7, 5}
		r := NewRechunkingReader(func() ([]byte, error) {
			if len(output) > 0 {
				n := output[0]
				output = output[1:]
				return []byte{byte(n)}, nil
			} else {
				return []byte{}, io.EOF
			}
		})

		buf := make([]byte, 4)

		n, err := r.Read(buf)
		Expect(n).To(Equal(1))
		Expect(buf[0]).To(Equal(byte(1)))
		Expect(err).To(BeNil())

		n, err = r.Read(buf)
		Expect(n).To(Equal(1))
		Expect(buf[0]).To(Equal(byte(7)))
		Expect(err).To(BeNil())

		n, err = r.Read(buf)
		Expect(n).To(Equal(1))
		Expect(buf[0]).To(Equal(byte(5)))

		checkEof(r)
	})

	It("propagnates errors", func() {
		myErr := errors.New("my custom error")
		r := NewRechunkingReader(func() ([]byte, error) {
			return nil, myErr
		})
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(myErr))
	})

	It("nops for empty buffer", func() {
		r := NewRechunkingReader(func() ([]byte, error) {
			return nil, nil
		})
		buf := make([]byte, 0)
		n, err := r.Read(buf)
		Expect(n).To(Equal(0))
		Expect(err).To(BeNil())

		n, err = r.Read(nil)
		Expect(n).To(Equal(0))
		Expect(err).To(BeNil())
	})

	It("short circuts after EOF", func() {
		r := NewRechunkingReader(func() ([]byte, error) {
			return nil, io.EOF
		})
		checkEof(r)
		checkEof(r)
	})
})

var _ = Describe("SplitReader", func() {
	It("interleaves reads on two readers", func() {
		source := ioutil.NopCloser(bytes.NewReader([]byte{1, 2}))
		rds := SplitReader(source, 1, 2)
		Expect(len(rds)).To(Equal(2))
		r1 := rds[0]
		defer r1.Close()
		r2 := rds[1]
		defer r2.Close()

		checkOneByte(r1, 1)
		checkOneByte(r2, 1)
		checkOneByte(r1, 2)
		checkOneByte(r2, 2)
		checkEof(r1)
		checkEof(r2)
	})

	It("continues with one reader closed", func() {
		source := &simpleCloser{bytes.NewReader([]byte{1, 2, 3, 4, 5}), false}
		rds := SplitReader(source, 2, 2)
		Expect(len(rds)).To(Equal(2))
		r1 := rds[0]
		r2 := rds[1]

		checkOneByte(r1, 1)
		checkOneByte(r2, 1)
		checkOneByte(r1, 2)
		checkOneByte(r2, 2)
		r2.Close()
		checkOneByte(r1, 3)
		checkOneByte(r1, 4)
		checkOneByte(r1, 5)
		checkEof(r1)
		r1.Close()

		time.Sleep(10 * time.Millisecond)
		Expect(source.closed).To(BeTrue())
	})

	It("closes source if all readers close immediately", func() {
		source := &simpleCloser{bytes.NewReader([]byte{1, 2, 3, 4, 5}), false}
		rds := SplitReader(source, 2, 2)
		Expect(len(rds)).To(Equal(2))
		r1 := rds[0]
		r2 := rds[1]
		r1.Close()
		r2.Close()

		time.Sleep(10 * time.Millisecond)
		Expect(source.closed).To(BeTrue())
	})

	It("handles concurrency correctly", func() {
		L := 16 * 1024 * 1024
		data := make([]byte, L)
		for i := 0; i < L; i++ {
			data[i] = byte(i)
		}
		source := &simpleCloser{bytes.NewReader(data), false}

		n := uint(3)
		var wg sync.WaitGroup
		wg.Add(int(n))

		rds := SplitReader(source, 1024, n)
		Expect(len(rds)).To(Equal(int(n)))
		for _, r := range rds {
			go func(r io.ReadCloser) {
				buf := make([]byte, 1024)
				for i := 0; i < L; {
					n, err := r.Read(buf)
					if err != nil && err != io.EOF {
						Fail("Read error " + err.Error())
					}
					for j := 0; j < n; j++ {
						if buf[j] != data[i] {
							Fail(fmt.Sprintf("Wrong data read %d %d %d %d", i, j, data[i], buf[j]))
						}
						i++
					}
				}
				r.Close()
				wg.Done()
			}(r)
		}

		wg.Wait()
		time.Sleep(10 * time.Millisecond)
		Expect(source.closed).To(BeTrue())
	})
})
