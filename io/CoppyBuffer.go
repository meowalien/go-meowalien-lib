package io

import (
	"errors"
	"io"
)

var errInvalidWrite = errors.New("invalid write result")

func CopyBufferWithCallback(dst io.Writer, src io.Reader, buf []byte, callback func(count int,written int64)) (written int64, err error) {
	//if wt, ok := src.(io.WriterTo); ok {
	//	return wt.WriteTo(dst)
	//}
	//fmt.Printf("res.Body8: %v",src)
	if buf == nil {
		err = errors.New("bf is nil")
		return
	}

	var count int
	for {
		nr, er := src.Read(buf)
		//fmt.Println("nr: ",nr)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}

		count ++
		if callback != nil{
			callback(count , written)
		}
	}
	return written, err
}
