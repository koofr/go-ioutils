package ioutils_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("StartReader", func() {
	It("should call OnStart when Read is first called", func() {
		br := bytes.NewReader([]byte("123"))

		counter := 0

		onStart := func() {
			counter += 1
		}

		r := NewStartReader(br, onStart)

		Expect(counter).To(Equal(0))

		n, err := r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(counter).To(Equal(1))

		n, err = r.Read([]byte{0})
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(counter).To(Equal(1))
	})
})
