package transport

import (
	"net"
	"time"

	"github.com/golang/snappy"
)

func NewCompStream(conn net.Conn) *CompStream {
	c := new(CompStream)
	c.conn = conn
	c.w = snappy.NewBufferedWriter(conn)
	c.r = snappy.NewReader(conn)
	return c
}
func NewCompConn(conn net.Conn) net.Conn {
	c := CompStream{}
	c.conn = conn
	c.w = snappy.NewBufferedWriter(conn)
	c.r = snappy.NewReader(conn)
	return &c
}

type CompStream struct {
	net.Conn
	conn net.Conn
	w    *snappy.Writer
	r    *snappy.Reader
}

func (c *CompStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *CompStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

func (c *CompStream) Close() (err error) {
	err = c.conn.Close()
	c.r = nil
	c.w = nil
	return
}
func (c *CompStream) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *CompStream) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *CompStream) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *CompStream) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *CompStream) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
