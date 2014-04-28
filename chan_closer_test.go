package ioutils_test

import (
	"bytes"
	. "git.koofr.lan/go-ioutils.git"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChanCloser", func() {
	It("write to close chanel when ReadCloser is closed", func() {
		br := bytes.NewReader([]byte("123"))

		doneCh := make(chan bool)

		r := NewChanCloser(br, doneCh)

		isDone := false

		go func() {
			<-doneCh
			isDone = true
		}()

		Expect(isDone).To(BeFalse())

		data, err := ioutil.ReadAll(r)

		Expect(data).To(Equal([]byte("123")))
		Expect(err).NotTo(HaveOccurred())

		Expect(isDone).To(BeFalse())

		err = r.Close()
		Expect(err).NotTo(HaveOccurred())

		Expect(isDone).To(BeTrue())
	})
})
