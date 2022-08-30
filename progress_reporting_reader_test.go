package ioutils_test

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("ProgressReportingReader", func() {
	It("should report progress", func() {
		r0 := strings.NewReader("0123456789")
		var changes []int64

		r := NewProgressReportingReader(r0, 2, func(bytesRead int64) {
			changes = append(changes, bytesRead)
		})

		_, err := r.Read(make([]byte, 1))
		Expect(err).NotTo(HaveOccurred())
		Expect(changes).To(BeEmpty())

		_, err = r.Read(make([]byte, 2))
		Expect(err).NotTo(HaveOccurred())
		Expect(changes).To(Equal([]int64{3}))

		_, err = r.Read(make([]byte, 3))
		Expect(err).NotTo(HaveOccurred())
		Expect(changes).To(Equal([]int64{3, 3}))

		n, err := r.Read(make([]byte, 10))
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(4))
		Expect(changes).To(Equal([]int64{3, 3, 4}))

		n, err = r.Read(make([]byte, 10))
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
		Expect(changes).To(Equal([]int64{3, 3, 4, 0}))
	})

	It("should report progress for an empty reader", func() {
		r0 := FuncReader(func(p []byte) (int, error) {
			return 0, io.EOF
		})
		var changes []int64

		r := NewProgressReportingReader(r0, 2, func(bytesRead int64) {
			changes = append(changes, bytesRead)
		})

		n, err := r.Read(make([]byte, 2))
		Expect(err).To(Equal(io.EOF))
		Expect(n).To(Equal(0))
		Expect(changes).To(Equal([]int64{0}))
	})
})
