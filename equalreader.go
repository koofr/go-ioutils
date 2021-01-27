package ioutils

import (
	"bytes"
	"io"
)

func CompareReaders(r1 io.Reader, r2 io.Reader) (bool, error) {
	bufSize := 32 * 1024
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	eof1 := false
	eof2 := false

	for !eof1 && !eof2 {
		n1, err1 := r1.Read(buf1)
		switch err1 {
		case io.EOF:
			eof1 = true
		case nil:
		default:
			return false, err1
		}

		n2, err2 := r2.Read(buf2)
		switch err2 {
		case io.EOF:
			eof2 = true
		case nil:
		default:
			return false, err2
		}

		switch {
		case n1 < n2:
			n, errPart := readPartial(r1, buf1, n1, n2)
			switch errPart {
			case io.EOF:
				eof1 = true
			case nil:
			default:
				return false, errPart
			}
			n1 = n
		case n2 < n1:
			n, errPart := readPartial(r2, buf2, n2, n1)
			switch errPart {
			case io.EOF:
				eof2 = true
			case nil:
			default:
				return false, errPart
			}
			n2 = n
		}

		if n1 != n2 {
			return false, nil
		}

		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
	}

	if eof1 && !eof2 {
		_, err2 := r2.Read(buf2)
		switch err2 {
		case io.EOF:
			// e1 EOF'ed before e2
			return true, nil
		case nil:
			return false, nil
		default:
			return false, err2
		}
	}

	if !eof1 && eof2 {
		_, err1 := r1.Read(buf1)
		switch err1 {
		case io.EOF:
			// e2 EOF'ed before e1
			return true, nil
		case nil:
			return false, nil
		default:
			return false, err1
		}
	}

	return true, nil
}

func readPartial(r io.Reader, buf []byte, n1, n2 int) (int, error) {
	for n1 < n2 {
		n, err := r.Read(buf[n1:n2])
		n1 += n
		if err != nil {
			return n1, err
		}
	}
	return n1, nil
}
