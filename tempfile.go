package ioutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TempFile struct {
	*os.File
}

func NewTempFile(prefix string) (t *TempFile, err error) {
	tmp, err := ioutil.TempFile(os.TempDir(), prefix)

	if err != nil {
		err = fmt.Errorf("NewTempFile error: %s", err)
		return
	}

	t = &TempFile{tmp}

	return
}

func (t *TempFile) remove() (err error) {
	err = os.Remove(t.File.Name())

	if err != nil {
		err = fmt.Errorf("TempFile remove error: %s", err)
	}

	return
}

func (t *TempFile) Close() error {
	err := t.File.Close()

	if err != nil {
		defer t.remove()
		return fmt.Errorf("TempFile close error: %s", err)
	}

	return t.remove()
}
