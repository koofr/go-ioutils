package ioutils_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("ExactSizeReader", func() {
	It("should succeed if the size matches", func() {
		b := []byte("12345")
		r := ioutil.NopCloser(bytes.NewReader(b))
		r = NewExactSizeReader(r, 5)
		bb, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).NotTo(HaveOccurred())
		Expect(bb).To(Equal(b))
	})

	It("should fail if a reader is too small", func() {
		b := []byte("1234")
		r := ioutil.NopCloser(bytes.NewReader(b))
		r = NewExactSizeReader(r, 5)
		_, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ErrExactSizeMismatch)).To(BeTrue())
	})

	It("should fail if a reader is too big", func() {
		b := []byte("123456")
		r := ioutil.NopCloser(bytes.NewReader(b))
		r = NewExactSizeReader(r, 5)
		_, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ErrExactSizeMismatch)).To(BeTrue())
	})

	It("should fail if a reader is too big (check EOF)", func() {
		b := []byte("123456")
		r := ioutil.NopCloser(bytes.NewReader(b))
		r = NewExactSizeReader(r, 5)
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ErrExactSizeMismatch)).To(BeTrue())
		Expect(n).To(Equal(0))
	})

	It("should fail if a reader is too big (check EOF larger buffer)", func() {
		r := ioutil.NopCloser(FuncReader(func(p []byte) (n int, err error) {
			if len(p) == 5 {
				return 5, nil
			}
			return 1, io.EOF
		}))
		r = NewExactSizeReader(r, 5)
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ErrExactSizeMismatch)).To(BeTrue())
		Expect(n).To(Equal(0))
	})
})
