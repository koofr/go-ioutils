package ioutils_test

import (
	"bytes"
	"errors"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("LazyReader", func() {
	It("should only build reader once", func() {
		c := 0

		r := NewLazyReader(func() (io.Reader, error) {
			c += 1
			return bytes.NewReader([]byte("123")), nil
		})

		Expect(c).To(Equal(0))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(c).To(Equal(1))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(c).To(Equal(1))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(c).To(Equal(1))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(io.EOF))
		Expect(c).To(Equal(1))
	})

	It("should only build error once", func() {
		e := errors.New("builder error")
		c := 0

		r := NewLazyReader(func() (io.Reader, error) {
			c += 1
			return nil, e
		})

		Expect(c).To(Equal(0))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(e))
		Expect(c).To(Equal(1))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(e))
		Expect(c).To(Equal(1))
	})
})
