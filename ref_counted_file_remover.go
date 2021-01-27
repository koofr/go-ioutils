package ioutils

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type RefCountedFileRemover struct {
	path string

	refCount         int
	refCountMutex    sync.Mutex
	totalClonedCount int64

	OnRemoveError func(path string, err error)
}

func NewRefCountedFileRemover(path string) *RefCountedFileRemover {
	return &RefCountedFileRemover{
		path: path,

		refCount: 0,
	}
}

func (r *RefCountedFileRemover) TotalClonedCount() int64 {
	return atomic.LoadInt64(&r.totalClonedCount)
}

func (r *RefCountedFileRemover) Open() (*RefCountedFileRemoveReader, error) {
	file, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open dest path: %w", err)
	}
	r.refCountMutex.Lock()
	r.refCount++
	r.refCountMutex.Unlock()
	atomic.AddInt64(&r.totalClonedCount, 1)

	reader := &RefCountedFileRemoveReader{
		File:    file,
		remover: r,

		isDereferenced: false,
	}

	return reader, nil
}

func (r *RefCountedFileRemover) derefAndRemoveIfUnused() {
	r.refCountMutex.Lock()
	r.refCount--
	isUnused := r.refCount == 0
	r.refCountMutex.Unlock()

	if isUnused {
		err := os.Remove(r.path)
		if err != nil {
			if r.OnRemoveError != nil {
				r.OnRemoveError(r.path, err)
			}
		}
	}
}

type RefCountedFileRemoveReader struct {
	*os.File

	remover        *RefCountedFileRemover
	isDereferenced bool
	onAfterClose   func()
}

func (r *RefCountedFileRemoveReader) Clone() (io.ReadCloser, error) {
	return r.remover.Open()
}

func (r *RefCountedFileRemoveReader) deref() {
	if r.isDereferenced {
		return
	}
	r.isDereferenced = true

	r.remover.derefAndRemoveIfUnused()
}

func (r *RefCountedFileRemoveReader) afterClose() {
	if r.onAfterClose != nil {
		r.onAfterClose()
	}

	r.deref()
}

func (r *RefCountedFileRemoveReader) Close() error {
	defer r.afterClose()

	return r.File.Close()
}

func (r *RefCountedFileRemoveReader) SetOnAfterClose(onAfterClose func()) {
	r.onAfterClose = onAfterClose
}

var _ ReadCloseCloner = (*RefCountedFileRemoveReader)(nil)
