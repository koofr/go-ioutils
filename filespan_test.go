package ioutils_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("FileSpan", func() {
	Describe("ApplyFileSpan", func() {
		It("should return the original reader if the span is nil", func() {
			r, err := ApplyFileSpan(bytes.NewReader([]byte("12345")), nil)
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal([]byte("12345")))
		})

		It("should skip bytes", func() {
			r, err := ApplyFileSpan(bytes.NewReader([]byte("12345")), &FileSpan{
				Start: 1,
				End:   4,
			})
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal([]byte("2345")))
		})

		It("should limit the reader", func() {
			r, err := ApplyFileSpan(bytes.NewReader([]byte("12345")), &FileSpan{
				Start: 0,
				End:   2,
			})
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal([]byte("123")))
		})

		It("should skip bytes and limit the reader", func() {
			r, err := ApplyFileSpan(bytes.NewReader([]byte("12345")), &FileSpan{
				Start: 1,
				End:   2,
			})
			Expect(err).NotTo(HaveOccurred())
			b, err := ioutil.ReadAll(r)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(Equal([]byte("23")))
		})
	})
})
