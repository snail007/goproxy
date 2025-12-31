package stealth

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// TrafficObfuscator obfuscates traffic to appear normal
type TrafficObfuscator struct {
	enabled bool
}

// NewTrafficObfuscator creates a new traffic obfuscator
func NewTrafficObfuscator(enabled bool) *TrafficObfuscator {
	return &TrafficObfuscator{
		enabled: enabled,
	}
}

// AddPadding adds random padding to data to avoid size-based fingerprinting
func (to *TrafficObfuscator) AddPadding(data []byte, minSize, maxSize int) []byte {
	if !to.enabled {
		return data
	}

	// Random padding size
	paddingSize := RandomInt(minSize, maxSize)

	// Generate random padding
	padding := make([]byte, paddingSize)
	rand.Read(padding)

	// Append padding
	return append(data, padding...)
}

// RandomDelay adds random delay to mimic human behavior
func (to *TrafficObfuscator) RandomDelay(minMs, maxMs int) {
	if !to.enabled {
		return
	}

	delay := RandomInt(minMs, maxMs)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ShapeTraffic shapes traffic to match normal patterns
func (to *TrafficObfuscator) ShapeTraffic(data []byte) []byte {
	if !to.enabled {
		return data
	}

	// Add small random delays between packets
	to.RandomDelay(10, 50)

	// Add minimal padding
	return to.AddPadding(data, 0, 16)
}

// GenerateRandomBytes generates random bytes for obfuscation
func GenerateRandomBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

// RandomInt generates a random integer in range [min, max]
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	if max-min > math.MaxInt32 {
		// Use big.Int for large ranges
		diff := big.NewInt(int64(max - min + 1))
		n, _ := rand.Int(rand.Reader, diff)
		return int(n.Int64()) + min
	}
	// Use simple random for small ranges
	diff := max - min + 1
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	return int(n.Int64()) + min
}
