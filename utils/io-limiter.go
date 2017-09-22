package utils

import (
	"context"
	"io"
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
