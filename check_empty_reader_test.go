package ioutils_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("CheckEmptyReader", func() {
	It("should return true if the reader is empty", func() {
		r := NewCheckEmptyReader(ioutil.NopCloser(bytes.NewReader(nil)))
		isEmpty, err := r.IsEmpty()
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeTrue())
		bb, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).NotTo(HaveOccurred())
		Expect(bb).To(BeEmpty())
	})

	It("should return false if the reader is not empty", func() {
		r := NewCheckEmptyReader(ioutil.NopCloser(bytes.NewReader([]byte{42})))
		isEmpty, err := r.IsEmpty()
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeFalse())
		bb, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).NotTo(HaveOccurred())
		Expect(bb).To(Equal([]byte{42}))
	})

	It("should return and not cache the error", func() {
		readCalls := 0
		r := NewCheckEmptyReader(ioutil.NopCloser(FuncReader(func(b []byte) (int, error) {
			readCalls++
			return 0, fmt.Errorf("custom error")
		})))
		_, err := r.IsEmpty()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("custom error"))
		Expect(readCalls).To(Equal(1))
		_, err = ioutil.ReadAll(r)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("custom error"))
		Expect(readCalls).To(Equal(2))
	})

	It("should call the original readers close", func() {
		closeCalled := true
		originalReader := NewPassCloseReader(bytes.NewReader(nil), func() error {
			closeCalled = true
			return nil
		})
		r := NewCheckEmptyReader(originalReader)
		isEmpty, err := r.IsEmpty()
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeTrue())
		bb, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).NotTo(HaveOccurred())
		Expect(bb).To(BeEmpty())
		Expect(closeCalled).To(BeTrue())
	})

	It("should fail if IsEmpty is called more than once", func() {
		r := NewCheckEmptyReader(ioutil.NopCloser(bytes.NewReader(nil)))
		isEmpty, err := r.IsEmpty()
		Expect(err).NotTo(HaveOccurred())
		Expect(isEmpty).To(BeTrue())
		_, err = r.IsEmpty()
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, ErrDirtyReader)).To(BeTrue())
	})

	It("should read even if IsEmpty is not called", func() {
		r := NewCheckEmptyReader(ioutil.NopCloser(bytes.NewReader([]byte{42})))
		bb, err := ioutil.ReadAll(r)
		r.Close()
		Expect(err).NotTo(HaveOccurred())
		Expect(bb).To(Equal([]byte{42}))
	})
})
