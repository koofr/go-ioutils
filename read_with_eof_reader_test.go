package ioutils_test

import (
	"io"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("ReadWithEOFReader", func() {
	It("should return n > 0 and io.EOF", func() {
		r := NewReadWithEOFReader()

		p := []byte{0}
		n, err := r.Read(p)
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(1))
	})

	Describe("ioutil.ReadAll", func() {
		It("should read ReadWithEOFReader", func() {
			r := NewReadWithEOFReader()
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(b)).To(BeNumerically(">", 0))
		})

		It("should read ReadWithEOFReader with limited length", func() {
			r := NewReadWithEOFReaderLen(2)
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal([]byte{1, 1}))
		})
	})
})
