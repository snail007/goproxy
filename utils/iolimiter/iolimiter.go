package iolimiter

import (
	"context"
	"io"
	"net"
	"time"

	"golang.org/x/time/rate"
)

const burstLimit = 1000 * 1000 * 1000

type Reader struct {
	r       io.Reader
	limiter *rate.Limiter
	ctx     context.Context
}

type Writer struct {
	w       io.Writer
	limiter *rate.Limiter
	ctx     context.Context
}

type conn struct {
	net.Conn
	r            io.Reader
	w            io.Writer
	readLimiter  *rate.Limiter
	writeLimiter *rate.Limiter
	ctx          context.Context
}

//NewtRateLimitConn sets rate limit (bytes/sec) to the Conn read and write.
func NewtConn(c net.Conn, bytesPerSec float64) net.Conn {
	s := &conn{
		Conn: c,
		r:    c,
		w:    c,
		ctx:  context.Background(),
	}
	s.readLimiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.readLimiter.AllowN(time.Now(), burstLimit) // spend initial burst
	s.writeLimiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.writeLimiter.AllowN(time.Now(), burstLimit) // spend initial burst
	return s
}

//NewtRateLimitReaderConn sets rate limit (bytes/sec) to the Conn read.
func NewReaderConn(c net.Conn, bytesPerSec float64) net.Conn {
	s := &conn{
		Conn: c,
		r:    c,
		w:    c,
		ctx:  context.Background(),
	}
	s.readLimiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.readLimiter.AllowN(time.Now(), burstLimit) // spend initial burst
	return s
}

//NewtRateLimitWriterConn sets rate limit (bytes/sec) to the Conn write.
func NewWriterConn(c net.Conn, bytesPerSec float64) net.Conn {
	s := &conn{
		Conn: c,
		r:    c,
		w:    c,
		ctx:  context.Background(),
	}
	s.writeLimiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.writeLimiter.AllowN(time.Now(), burstLimit) // spend initial burst
	return s
}

// Read reads bytes into p.
func (s *conn) Read(p []byte) (int, error) {
	if s.readLimiter == nil {
		return s.r.Read(p)
	}
	n, err := s.r.Read(p)
	if err != nil {
		return n, err
	}
	if err := s.readLimiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, nil
}

// Write writes bytes from p.
func (s *conn) Write(p []byte) (int, error) {
	if s.writeLimiter == nil {
		return s.w.Write(p)
	}
	n, err := s.w.Write(p)
	if err != nil {
		return n, err
	}
	if err := s.writeLimiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, err
}
func (s *conn) Close() error {
	if s.Conn != nil {
		e := s.Conn.Close()
		s.Conn = nil
		s.r = nil
		s.w = nil
		s.readLimiter = nil
		s.writeLimiter = nil
		s.ctx = nil
		return e
	}
	return nil
}

// NewReader returns a reader that implements io.Reader with rate limiting.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   r,
		ctx: context.Background(),
	}
}

// NewReaderWithContext returns a reader that implements io.Reader with rate limiting.
func NewReaderWithContext(r io.Reader, ctx context.Context) *Reader {
	return &Reader{
		r:   r,
		ctx: ctx,
	}
}

// NewWriter returns a writer that implements io.Writer with rate limiting.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:   w,
		ctx: context.Background(),
	}
}

// NewWriterWithContext returns a writer that implements io.Writer with rate limiting.
func NewWriterWithContext(w io.Writer, ctx context.Context) *Writer {
	return &Writer{
		w:   w,
		ctx: ctx,
	}
}

// SetRateLimit sets rate limit (bytes/sec) to the reader.
func (s *Reader) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit) // spend initial burst
}

// Read reads bytes into p.
func (s *Reader) Read(p []byte) (int, error) {
	if s.limiter == nil {
		return s.r.Read(p)
	}
	n, err := s.r.Read(p)
	if err != nil {
		return n, err
	}
	if err := s.limiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, nil
}

// SetRateLimit sets rate limit (bytes/sec) to the writer.
func (s *Writer) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit) // spend initial burst
}

// Write writes bytes from p.
func (s *Writer) Write(p []byte) (int, error) {
	if s.limiter == nil {
		return s.w.Write(p)
	}
	n, err := s.w.Write(p)
	if err != nil {
		return n, err
	}
	if err := s.limiter.WaitN(s.ctx, n); err != nil {
		return n, err
	}
	return n, err
}
