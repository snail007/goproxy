package encrypt

import (
	"crypto/cipher"
	"fmt"
	"io"
	"net"

	lbuf "github.com/snail007/goproxy/core/lib/buf"
)

var (
	lBuf = lbuf.NewLeakyBuf(2048, 2048)
)

type Conn struct {
	net.Conn
	*Cipher
	w io.Writer
	r io.Reader
}

func NewConn(c net.Conn, method, password string) (conn net.Conn, err error) {
	cipher0, err := NewCipher(method, password)
	if err != nil {
		return
	}
	conn = &Conn{
		Conn:   c,
		Cipher: cipher0,
		r:      &cipher.StreamReader{S: cipher0.ReadStream, R: c},
		w:      &cipher.StreamWriter{S: cipher0.WriteStream, W: c},
	}
	return
}
func (s *Conn) Read(b []byte) (n int, err error) {
	if s.r == nil {
		return 0, fmt.Errorf("use of closed network connection")
	}
	return s.r.Read(b)
}
func (s *Conn) Write(b []byte) (n int, err error) {
	if s.w == nil {
		return 0, fmt.Errorf("use of closed network connection")
	}
	return s.w.Write(b)
}
func (s *Conn) Close() (err error) {
	if s.Cipher != nil {
		err = s.Conn.Close()
		s.Cipher = nil
		s.r = nil
		s.w = nil
	}
	return
}
