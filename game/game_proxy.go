package game

import (
	"log"
	"net"
	"sync"
	"time"
)

// GameProxy provides game-optimized proxy functionality
type GameProxy struct {
	udpConnections map[string]*UDPConnection
	tcpConnections map[string]*TCPConnection
	mu             sync.RWMutex
	enabled        bool
}

// UDPConnection represents a UDP game connection
type UDPConnection struct {
	conn        *net.UDPConn
	lastActive  time.Time
	remoteAddr  *net.UDPAddr
	localAddr   *net.UDPAddr
	packetCount int64
}

// TCPConnection represents a TCP game connection
type TCPConnection struct {
	conn       net.Conn
	lastActive time.Time
	remoteAddr net.Addr
	localAddr  net.Addr
	bytesTx    int64
	bytesRx    int64
}

// NewGameProxy creates a new game proxy
func NewGameProxy() *GameProxy {
	return &GameProxy{
		udpConnections: make(map[string]*UDPConnection),
		tcpConnections: make(map[string]*TCPConnection),
		enabled:        true,
	}
}

// Enable enables game proxy optimizations
func (gp *GameProxy) Enable() {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	gp.enabled = true
}

// Disable disables game proxy optimizations
func (gp *GameProxy) Disable() {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	gp.enabled = false
}

// HandleUDP handles UDP game traffic with low latency
func (gp *GameProxy) HandleUDP(localConn *net.UDPConn, remoteAddr *net.UDPAddr, data []byte) error {
	if !gp.enabled {
		return nil
	}

	gp.mu.Lock()
	connKey := remoteAddr.String()
	conn, exists := gp.udpConnections[connKey]
	if !exists {
		conn = &UDPConnection{
			conn:        localConn,
			remoteAddr:  remoteAddr,
			lastActive:  time.Now(),
			packetCount: 0,
		}
		gp.udpConnections[connKey] = conn
	}
	conn.lastActive = time.Now()
	conn.packetCount++
	gp.mu.Unlock()

	// Forward packet immediately with minimal processing
	_, err := localConn.WriteToUDP(data, remoteAddr)
	if err != nil {
		log.Printf("UDP forward error: %v", err)
		return err
	}

	return nil
}

// HandleTCP handles TCP game traffic with optimizations
func (gp *GameProxy) HandleTCP(conn net.Conn) error {
	if !gp.enabled {
		return nil
	}

	gp.mu.Lock()
	connKey := conn.RemoteAddr().String()
	tcpConn := &TCPConnection{
		conn:       conn,
		lastActive: time.Now(),
		remoteAddr: conn.RemoteAddr(),
		localAddr:  conn.LocalAddr(),
		bytesTx:    0,
		bytesRx:    0,
	}
	gp.tcpConnections[connKey] = tcpConn
	gp.mu.Unlock()

	// Set TCP options for low latency
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true) // Disable Nagle's algorithm
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	return nil
}

// CleanupStaleConnections removes stale connections
func (gp *GameProxy) CleanupStaleConnections(maxIdle time.Duration) {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	now := time.Now()

	// Cleanup UDP connections
	for key, conn := range gp.udpConnections {
		if now.Sub(conn.lastActive) > maxIdle {
			delete(gp.udpConnections, key)
		}
	}

	// Cleanup TCP connections
	for key, conn := range gp.tcpConnections {
		if now.Sub(conn.lastActive) > maxIdle {
			if conn.conn != nil {
				conn.conn.Close()
			}
			delete(gp.tcpConnections, key)
		}
	}
}

// GetStats returns proxy statistics
func (gp *GameProxy) GetStats() map[string]interface{} {
	gp.mu.RLock()
	defer gp.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["udp_connections"] = len(gp.udpConnections)
	stats["tcp_connections"] = len(gp.tcpConnections)
	stats["enabled"] = gp.enabled

	var totalPackets int64
	for _, conn := range gp.udpConnections {
		totalPackets += conn.packetCount
	}
	stats["total_udp_packets"] = totalPackets

	return stats
}

// StartCleanup starts the cleanup goroutine
func (gp *GameProxy) StartCleanup(interval time.Duration, maxIdle time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			gp.CleanupStaleConnections(maxIdle)
		}
	}()
}
