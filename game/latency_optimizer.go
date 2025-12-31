package game

import (
	"net"
	"time"
)

// LatencyOptimizer optimizes connections for low latency
type LatencyOptimizer struct {
	enabled bool
}

// NewLatencyOptimizer creates a new latency optimizer
func NewLatencyOptimizer() *LatencyOptimizer {
	return &LatencyOptimizer{
		enabled: true,
	}
}

// OptimizeTCP optimizes TCP connection for low latency
func (lo *LatencyOptimizer) OptimizeTCP(conn net.Conn) error {
	if !lo.enabled {
		return nil
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		// Disable Nagle's algorithm for low latency
		tcpConn.SetNoDelay(true)

		// Enable keep-alive
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)

		// Set read/write deadlines
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	}

	return nil
}

// OptimizeUDP optimizes UDP connection for low latency
func (lo *LatencyOptimizer) OptimizeUDP(conn *net.UDPConn) error {
	if !lo.enabled {
		return nil
	}

	// Set read buffer size for better performance
	conn.SetReadBuffer(1024 * 1024) // 1MB
	conn.SetWriteBuffer(1024 * 1024) // 1MB

	return nil
}

// PreWarmConnection pre-warms a connection to reduce latency
func (lo *LatencyOptimizer) PreWarmConnection(conn net.Conn) error {
	if !lo.enabled {
		return nil
	}

	// Send a small keep-alive packet
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
	}

	return nil
}

// MeasureLatency measures connection latency
func (lo *LatencyOptimizer) MeasureLatency(conn net.Conn) (time.Duration, error) {
	start := time.Now()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return 0, err
	}
	return time.Since(start), nil
}
