package stealth

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// HeaderRandomizer randomizes HTTP headers to avoid fingerprinting
type HeaderRandomizer struct {
	userAgents []string
	browsers   map[string][]string
}

// NewHeaderRandomizer creates a new header randomizer
func NewHeaderRandomizer() *HeaderRandomizer {
	hr := &HeaderRandomizer{
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
		},
		browsers: make(map[string][]string),
	}

	// Chrome headers
	hr.browsers["chrome"] = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language: en-US,en;q=0.9",
		"Accept-Encoding: gzip, deflate, br",
		"Connection: keep-alive",
		"Upgrade-Insecure-Requests: 1",
		"Sec-Fetch-Dest: document",
		"Sec-Fetch-Mode: navigate",
		"Sec-Fetch-Site: none",
		"Sec-Fetch-User: ?1",
		"Cache-Control: max-age=0",
	}

	// Safari headers
	hr.browsers["safari"] = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language: en-US,en;q=0.9",
		"Accept-Encoding: gzip, deflate, br",
		"Connection: keep-alive",
		"Upgrade-Insecure-Requests: 1",
	}

	// Edge headers
	hr.browsers["edge"] = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language: en-US,en;q=0.9",
		"Accept-Encoding: gzip, deflate, br",
		"Connection: keep-alive",
		"Upgrade-Insecure-Requests: 1",
		"Sec-Fetch-Dest: document",
		"Sec-Fetch-Mode: navigate",
		"Sec-Fetch-Site: none",
		"Sec-Fetch-User: ?1",
	}

	return hr
}

// GetRandomUserAgent returns a random user agent
func (hr *HeaderRandomizer) GetRandomUserAgent() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(hr.userAgents))))
	return hr.userAgents[int(n.Int64())]
}

// GetRandomHeaders returns randomized headers for a browser type
func (hr *HeaderRandomizer) GetRandomHeaders(browserType string) []string {
	headers, ok := hr.browsers[browserType]
	if !ok {
		headers = hr.browsers["chrome"] // Default to Chrome
	}

	// Shuffle headers to avoid fingerprinting
	shuffled := make([]string, len(headers))
	copy(shuffled, headers)
	
	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(n.Int64())
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

// RemoveProxyHeaders removes proxy-related headers from request
func (hr *HeaderRandomizer) RemoveProxyHeaders(headers string) string {
	proxyHeaders := []string{
		"Via:",
		"X-Forwarded-For:",
		"X-Forwarded-Proto:",
		"X-Real-IP:",
		"Proxy-Connection:",
		"X-Proxy-Authorization:",
		"Forwarded:",
	}

	lines := strings.Split(headers, "\r\n")
	var filtered []string

	for _, line := range lines {
		shouldRemove := false
		lineUpper := strings.ToUpper(line)
		for _, proxyHeader := range proxyHeaders {
			if strings.HasPrefix(lineUpper, strings.ToUpper(proxyHeader)) {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\r\n")
}

// InjectHeaders injects realistic browser headers into request
func (hr *HeaderRandomizer) InjectHeaders(requestHeaders string, browserType string) string {
	// Remove proxy headers first
	cleaned := hr.RemoveProxyHeaders(requestHeaders)

	// Get random headers
	randomHeaders := hr.GetRandomHeaders(browserType)

	// Parse existing headers
	lines := strings.Split(cleaned, "\r\n")
	headerMap := make(map[string]string)

	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(strings.ToLower(parts[0]))
				headerMap[key] = strings.TrimSpace(parts[1])
			}
		}
	}

	// Inject random headers (don't override existing important headers)
	for _, header := range randomHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			if key != "user-agent" && key != "host" && key != "content-length" {
				if _, exists := headerMap[key]; !exists {
					headerMap[key] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Rebuild headers
	var result []string
	result = append(result, lines[0]) // Request line

	// Add User-Agent if not present
	if _, exists := headerMap["user-agent"]; !exists {
		result = append(result, "User-Agent: "+hr.GetRandomUserAgent())
	}

	// Add other headers
	for key, value := range headerMap {
		if key != "user-agent" {
			result = append(result, strings.Title(key)+": "+value)
		}
	}

	return strings.Join(result, "\r\n")
}
