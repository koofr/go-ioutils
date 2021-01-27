package ioutils_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("Equalreader", func() {
	Describe("CompareReaders", func() {
		It("should return true if reads are equal", func() {
			r1 := strings.NewReader("test")
			r2 := bytes.NewReader([]byte("test"))
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		It("should return the first reader's error", func() {
			err1 := fmt.Errorf("error1")
			r1 := io.MultiReader(
				strings.NewReader("test"),
				NewErrorReader(err1),
			)
			r2 := bytes.NewReader([]byte("test"))
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err1)).To(BeTrue())
		})

		It("should return the second reader's error", func() {
			err2 := fmt.Errorf("error2")
			r1 := strings.NewReader("test")
			r2 := io.MultiReader(
				bytes.NewReader([]byte("test")),
				NewErrorReader(err2),
			)
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err2)).To(BeTrue())
		})

		It("should read different chunk sizes (1)", func() {
			r1 := io.MultiReader(
				strings.NewReader("te"),
				strings.NewReader("st12"),
				strings.NewReader("34"),
			)
			r2 := io.MultiReader(
				bytes.NewReader([]byte("test")),
				bytes.NewReader([]byte("1234")),
			)
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		It("should read different chunk sizes (2)", func() {
			r1 := io.MultiReader(
				strings.NewReader("test"),
				strings.NewReader("1234"),
			)
			r2 := io.MultiReader(
				bytes.NewReader([]byte("te")),
				bytes.NewReader([]byte("st12")),
				bytes.NewReader([]byte("34")),
			)
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		It("should read different chunk sizes (3)", func() {
			r1 := io.MultiReader(
				strings.NewReader("te"),
				strings.NewReader("st"),
			)
			r2 := io.MultiReader(
				bytes.NewReader([]byte("test1")),
			)
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})

		It("should read different chunk sizes (4)", func() {
			err1 := fmt.Errorf("error1")
			r1 := io.MultiReader(
				strings.NewReader("te"),
				strings.NewReader("st"),
				NewErrorReader(err1),
			)
			r2 := io.MultiReader(
				bytes.NewReader([]byte("test1")),
			)
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err1)).To(BeTrue())
		})

		It("should read different chunk sizes (5)", func() {
			err2 := fmt.Errorf("error2")
			r1 := io.MultiReader(
				strings.NewReader("test1"),
			)
			r2 := io.MultiReader(
				bytes.NewReader([]byte("te")),
				bytes.NewReader([]byte("st")),
				NewErrorReader(err2),
			)
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err2)).To(BeTrue())
		})

		It("should handle e1 EOF before e2", func() {
			r1 := NewReadWithEOFReader()
			b2 := bytes.Repeat([]byte{1}, 32*1024)
			r2 := bytes.NewReader(b2)
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		It("should handle e1 EOF before e2 and e2 has more", func() {
			r1 := NewReadWithEOFReader()
			b2 := bytes.Repeat([]byte{1}, 32*1024+1)
			r2 := bytes.NewReader(b2)
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})

		It("should handle e1 EOF before e2 and e2 errors out", func() {
			err2 := fmt.Errorf("error2")
			r1 := NewReadWithEOFReader()
			r2 := io.MultiReader(
				bytes.NewReader(bytes.Repeat([]byte{1}, 32*1024)),
				NewErrorReader(err2),
			)
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err2)).To(BeTrue())
		})

		It("should handle e2 EOF before e1", func() {
			b1 := bytes.Repeat([]byte{1}, 32*1024)
			r1 := bytes.NewReader(b1)
			r2 := NewReadWithEOFReader()
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})

		It("should handle e2 EOF before e1 and e1 has more", func() {
			b1 := bytes.Repeat([]byte{1}, 32*1024+1)
			r1 := bytes.NewReader(b1)
			r2 := NewReadWithEOFReader()
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})

		It("should handle e2 EOF before e1 and e1 errors out", func() {
			err1 := fmt.Errorf("error1")
			r1 := io.MultiReader(
				bytes.NewReader(bytes.Repeat([]byte{1}, 32*1024)),
				NewErrorReader(err1),
			)
			r2 := NewReadWithEOFReader()
			_, err := CompareReaders(r1, r2)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, err1)).To(BeTrue())
		})

		It("should return false if content is different", func() {
			r1 := strings.NewReader("test")
			r2 := bytes.NewReader([]byte("foob"))
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})

		It("should return false if length is different", func() {
			r1 := strings.NewReader("test1")
			r2 := bytes.NewReader([]byte("foob"))
			ok, err := CompareReaders(r1, r2)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeFalse())
		})
	})
})
