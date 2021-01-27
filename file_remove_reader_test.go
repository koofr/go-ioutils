package ioutils_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("FileRemoveReader", func() {
	It("should remove file after closing", func() {
		file, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())

		_, err = file.Write([]byte("123"))
		Expect(err).NotTo(HaveOccurred())
		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		_, err = os.Stat(file.Name())
		Expect(err).NotTo(HaveOccurred())

		newfile, err := os.Open(file.Name())
		Expect(err).NotTo(HaveOccurred())

		reader := NewFileRemoveReader(newfile)
		content, err := ioutil.ReadAll(reader)
		Expect(err).NotTo(HaveOccurred())
		Expect(content).To(Equal([]byte("123")))

		err = reader.Close()
		Expect(err).NotTo(HaveOccurred())

		_, err = os.Stat(file.Name())
		Expect(os.IsNotExist(err)).To(BeTrue())
	})
})
