package ioutils_test

import (
	"bytes"
	"context"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("CtxReader", func() {
	It("should stop reading when the context is canceled", func() {
		br := bytes.NewReader([]byte{0, 1, 2, 3, 4})

		ctx, cancel := context.WithCancel(context.Background())

		ctxReader := NewCtxReader(ctx, br)

		p := []byte{0, 0, 0}
		n, err := ctxReader.Read(p)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(3))
		Expect(p).To(Equal([]byte{0, 1, 2}))

		cancel()

		n, err = ctxReader.Read(p)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("context canceled"))
		Expect(n).To(Equal(0))
	})

	It("should stop reading when the context is canceled after EOF", func() {
		br := bytes.NewReader([]byte{0, 1, 2, 3, 4})

		ctx, cancel := context.WithCancel(context.Background())

		ctxReader := NewCtxReader(ctx, br)

		p := []byte{0, 0, 0, 0, 0, 0}
		n, err := ctxReader.Read(p)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(5))
		Expect(p).To(Equal([]byte{0, 1, 2, 3, 4, 0}))

		n, err = ctxReader.Read(p)
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))

		cancel()

		n, err = ctxReader.Read(p)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("context canceled"))
		Expect(n).To(Equal(0))
	})
})
