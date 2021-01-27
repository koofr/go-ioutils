package ioutils_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("DebugReader", func() {
	It("should display debug info", func() {
		c := 0
		r := FuncReader(func(p []byte) (int, error) {
			c++
			time.Sleep(100 * time.Millisecond)
			if c == 1 {
				return len(p), nil
			}
			return 3, fmt.Errorf("custom error")
		})

		dr := NewDebugReader(r)

		n, err := dr.Read([]byte{0, 0, 0, 0, 0})
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(5))

		n, err = dr.Read([]byte{0, 0, 0, 0, 0})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("custom error"))
		Expect(n).To(Equal(3))

		display := fmt.Sprintf("%v", dr)
		Expect(display).To(MatchRegexp("^lastReadStarted=[^ ]+ lastReadEnded=[^ ]+ lastReadDuration=\\d+ ms lastReadErr=custom error lastReadN=3 totalN=8$"))
	})
})
