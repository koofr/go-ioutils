package ioutils_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("HashReader", func() {
	It("should calculate hash", func() {
		br := bytes.NewReader([]byte("123"))

		r := NewHashReader(br)

		Expect(r.Hash()).To(Equal("d41d8cd98f00b204e9800998ecf8427e"))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Hash()).To(Equal("c4ca4238a0b923820dcc509a6f75849b"))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Hash()).To(Equal("c20ad4d76fe97759aa27a0c99bff6710"))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Hash()).To(Equal("202cb962ac59075b964b07152d234b70"))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(io.EOF))
		Expect(r.Hash()).To(Equal("202cb962ac59075b964b07152d234b70"))
	})

	It("should resume hash calculation", func() {
		br := bytes.NewReader([]byte("123"))

		r := NewHashReader(br)

		Expect(r.Hash()).To(Equal("d41d8cd98f00b204e9800998ecf8427e"))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Hash()).To(Equal("c4ca4238a0b923820dcc509a6f75849b"))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Hash()).To(Equal("c20ad4d76fe97759aa27a0c99bff6710"))

		r1 := NewHashReaderFromState(br, r.GetState())

		n, err = r1.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(r1.Hash()).To(Equal("202cb962ac59075b964b07152d234b70"))

		n, err = r1.Read([]byte{0})
		Expect(n).To(Equal(0))
		Expect(err).To(Equal(io.EOF))
		Expect(r1.Hash()).To(Equal("202cb962ac59075b964b07152d234b70"))
	})
})
