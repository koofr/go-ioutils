package ioutils_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("DelayedEOFReader", func() {
	It("should handle empty reader", func() {
		r := NewDelayedEOFReader(bytes.NewReader(nil))

		n, err := r.Read([]byte{0})
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
	})

	It("should handle reader", func() {
		r := NewDelayedEOFReader(bytes.NewReader([]byte{1}))

		p := []byte{0}
		n, err := r.Read(p)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(1))
		Expect(p[0]).To(Equal(uint8(1)))

		n, err = r.Read([]byte{0})
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
	})

	It("should delay the EOF", func() {
		r := NewDelayedEOFReader(NewReadWithEOFReaderLen(1))

		p := []byte{0}
		n, err := r.Read(p)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(1))
		Expect(p).To(Equal([]byte{1}))

		n, err = r.Read([]byte{0})
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
	})
})
