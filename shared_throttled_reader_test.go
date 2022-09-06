package ioutils_test

import (
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/koofr/go-ioutils"
)

var _ = Describe("SharedLimiter", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("should limit", func() {
		l := NewSharedLimiter(1024 * 1024)

		start := time.Now()
		Expect(l.WaitN(ctx, 1)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 50*time.Millisecond))

		start = time.Now()
		Expect(l.WaitN(ctx, 100)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 200*time.Millisecond))

		start = time.Now()
		Expect(l.WaitN(ctx, 500*1024)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 600*time.Millisecond))

		l.SetLimit(10)

		start = time.Now()
		Expect(l.WaitN(ctx, 20)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 2000*time.Millisecond))
	})

	It("should disable limit", func() {
		l := NewSharedLimiter(0)

		start := time.Now()
		Expect(l.WaitN(ctx, 1)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 20*time.Millisecond))

		l.SetLimit(10)

		start = time.Now()
		Expect(l.WaitN(ctx, 1)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically(">", 20*time.Millisecond))
		Expect(time.Since(start)).To(BeNumerically("<", 200*time.Millisecond))

		l.SetLimit(0)

		start = time.Now()
		Expect(l.WaitN(ctx, 1)).To(Succeed())
		Expect(time.Since(start)).To(BeNumerically("<", 20*time.Millisecond))
	})

	Describe("SharedThrottledReader", func() {
		It("should limit the reader", func() {
			l := NewSharedLimiter(10 * 1024)

			r := NewSharedThrottledReader(ctx, io.NopCloser(bytes.NewReader(make([]byte, 1024*1024))), l)

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer GinkgoRecover()
				n, err := r.Read(make([]byte, 15*1024))
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(15 * 1024))
				wg.Done()
			}()

			go func() {
				defer GinkgoRecover()
				n, err := r.Read(make([]byte, 15*1024))
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(15 * 1024))
				wg.Done()
			}()

			start := time.Now()
			wg.Wait()
			Expect(time.Since(start)).To(BeNumerically(">", 500*time.Millisecond))
			Expect(time.Since(start)).To(BeNumerically("<", 4000*time.Millisecond))
		})
	})
})
