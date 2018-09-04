package ioutils

import (
	"io"
	logger "log"

	lbuf "github.com/snail007/goproxy/core/lib/buf"
)

func IoBind(dst io.ReadWriteCloser, src io.ReadWriteCloser, fn func(err interface{}), log *logger.Logger) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("bind crashed %s", err)
			}
		}()
		e1 := make(chan interface{}, 1)
		e2 := make(chan interface{}, 1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("bind crashed %s", err)
				}
			}()
			//_, err := io.Copy(dst, src)
			err := ioCopy(dst, src)
			e1 <- err
		}()
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("bind crashed %s", err)
				}
			}()
			//_, err := io.Copy(src, dst)
			err := ioCopy(src, dst)
			e2 <- err
		}()
		var err interface{}
		select {
		case err = <-e1:
			//log.Printf("e1")
		case err = <-e2:
			//log.Printf("e2")
		}
		src.Close()
		dst.Close()
		if fn != nil {
			fn(err)
		}
	}()
}
func ioCopy(dst io.ReadWriter, src io.ReadWriter) (err error) {
	buf := lbuf.LeakyBuffer.Get()
	defer lbuf.LeakyBuffer.Put(buf)
	n := 0
	for {
		n, err = src.Read(buf)
		if n > 0 {
			if _, e := dst.Write(buf[0:n]); e != nil {
				return e
			}
		}
		if err != nil {
			return
		}
	}
}
