#!/bin/bash

# Test script for proxy server connection

SERVER_IP="194.233.73.36"
SERVER_IPV6="2407:3640:2299:7804::1"
PROXY_PORT="8080"

echo "Testing Game Proxy Server Connection..."
echo "========================================"

# Test IPv4 HTTP Proxy
echo ""
echo "1. Testing IPv4 HTTP Proxy..."
if curl -s --connect-timeout 10 -x http://$SERVER_IP:$PROXY_PORT http://httpbin.org/ip > /dev/null; then
    echo "✅ IPv4 HTTP Proxy: Working"
    curl -s -x http://$SERVER_IP:$PROXY_PORT http://httpbin.org/ip
else
    echo "❌ IPv4 HTTP Proxy: Failed"
fi

# Test IPv6 HTTP Proxy
echo ""
echo "2. Testing IPv6 HTTP Proxy..."
if curl -s --connect-timeout 10 -x http://[$SERVER_IPV6]:$PROXY_PORT http://httpbin.org/ip > /dev/null; then
    echo "✅ IPv6 HTTP Proxy: Working"
    curl -s -x http://[$SERVER_IPV6]:$PROXY_PORT http://httpbin.org/ip
else
    echo "❌ IPv6 HTTP Proxy: Failed (may not be configured)"
fi

# Test Headers (Stealth Check)
echo ""
echo "3. Testing Stealth Features (Header Check)..."
HEADERS=$(curl -s -x http://$SERVER_IP:$PROXY_PORT http://httpbin.org/headers)
if echo "$HEADERS" | grep -q "Via\|X-Forwarded-For\|Proxy-Connection"; then
    echo "⚠️  Warning: Proxy headers detected!"
else
    echo "✅ Stealth Mode: Headers look clean"
fi

# Test Latency
echo ""
echo "4. Testing Latency..."
LATENCY=$(curl -o /dev/null -s -w '%{time_total}' -x http://$SERVER_IP:$PROXY_PORT http://httpbin.org/ip)
LATENCY_MS=$(echo "$LATENCY * 1000" | bc | cut -d. -f1)
echo "Latency: ${LATENCY_MS}ms"

if [ "$LATENCY_MS" -lt 500 ]; then
    echo "✅ Latency: Good (<500ms)"
else
    echo "⚠️  Latency: High (>500ms)"
fi

# Test HTTPS
echo ""
echo "5. Testing HTTPS Proxy..."
if curl -s --connect-timeout 10 -x http://$SERVER_IP:$PROXY_PORT https://httpbin.org/ip > /dev/null; then
    echo "✅ HTTPS Proxy: Working"
else
    echo "❌ HTTPS Proxy: Failed"
fi

echo ""
echo "========================================"
echo "Test Complete!"
echo ""
echo "Proxy Endpoints:"
echo "  IPv4: http://$SERVER_IP:$PROXY_PORT"
echo "  IPv6: http://[$SERVER_IPV6]:$PROXY_PORT"
echo ""
