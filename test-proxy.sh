#!/bin/bash

# Quick proxy connection test script
# Tests connection to proxy server: 194.233.73.36:8080

PROXY_SERVER="194.233.73.36"
PROXY_PORT="8080"
PROXY_URL="http://${PROXY_SERVER}:${PROXY_PORT}"

echo "=========================================="
echo "Proxy Connection Test"
echo "=========================================="
echo "Server: $PROXY_SERVER:$PROXY_PORT"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test 1: Basic connectivity
echo -n "1. Testing connectivity: "
if ping -c 1 -W 2 $PROXY_SERVER > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Server reachable${NC}"
else
    echo -e "${RED}❌ Server unreachable${NC}"
    exit 1
fi

# Test 2: Port check
echo -n "2. Testing port $PROXY_PORT: "
if nc -zv -w 2 $PROXY_SERVER $PROXY_PORT > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Port open${NC}"
else
    echo -e "${RED}❌ Port closed or filtered${NC}"
fi

# Test 3: HTTP Proxy
echo -n "3. Testing HTTP proxy: "
RESPONSE=$(curl -s --connect-timeout 5 -x $PROXY_URL http://httpbin.org/ip 2>&1)
if echo "$RESPONSE" | grep -q "origin"; then
    echo -e "${GREEN}✅ Working${NC}"
    echo "   Response: $RESPONSE"
else
    echo -e "${RED}❌ Failed${NC}"
    echo "   Error: $RESPONSE"
fi

# Test 4: HTTPS Proxy
echo -n "4. Testing HTTPS proxy: "
RESPONSE=$(curl -s --connect-timeout 5 -x $PROXY_URL https://httpbin.org/ip 2>&1)
if echo "$RESPONSE" | grep -q "origin"; then
    echo -e "${GREEN}✅ Working${NC}"
else
    echo -e "${YELLOW}⚠️  May need configuration${NC}"
fi

# Test 5: Latency
echo -n "5. Testing latency: "
START=$(date +%s%N)
curl -s -x $PROXY_URL http://httpbin.org/ip > /dev/null 2>&1
END=$(date +%s%N)
LATENCY=$(( (END - START) / 1000000 ))
if [ $LATENCY -lt 500 ]; then
    echo -e "${GREEN}✅ ${LATENCY}ms (Good)${NC}"
elif [ $LATENCY -lt 1000 ]; then
    echo -e "${YELLOW}⚠️  ${LATENCY}ms (Acceptable)${NC}"
else
    echo -e "${RED}❌ ${LATENCY}ms (High)${NC}"
fi

# Test 6: Stealth check (headers)
echo -n "6. Testing stealth mode: "
HEADERS=$(curl -s -x $PROXY_URL http://httpbin.org/headers)
if echo "$HEADERS" | grep -q "Via\|X-Forwarded-For\|Proxy-Connection"; then
    echo -e "${YELLOW}⚠️  Proxy headers detected${NC}"
else
    echo -e "${GREEN}✅ No proxy headers (Stealth active)${NC}"
fi

# Test 7: IPv6 (if available)
echo -n "7. Testing IPv6: "
IPV6_URL="http://[2407:3640:2299:7804::1]:${PROXY_PORT}"
RESPONSE=$(curl -s --connect-timeout 5 -x $IPV6_URL http://httpbin.org/ip 2>&1)
if echo "$RESPONSE" | grep -q "origin"; then
    echo -e "${GREEN}✅ Working${NC}"
else
    echo -e "${YELLOW}⚠️  IPv6 not configured or unavailable${NC}"
fi

echo ""
echo "=========================================="
echo "Connection Examples:"
echo "=========================================="
echo ""
echo "curl:"
echo "  curl -x $PROXY_URL http://httpbin.org/ip"
echo ""
echo "Environment variables:"
echo "  export http_proxy=$PROXY_URL"
echo "  export https_proxy=$PROXY_URL"
echo "  curl http://httpbin.org/ip"
echo ""
echo "Browser:"
echo "  HTTP Proxy: $PROXY_SERVER:$PROXY_PORT"
echo "  HTTPS Proxy: $PROXY_SERVER:$PROXY_PORT"
echo ""
echo "Mac System Proxy:"
echo "  networksetup -setwebproxy \"Wi-Fi\" $PROXY_SERVER $PROXY_PORT"
echo "  networksetup -setsecurewebproxy \"Wi-Fi\" $PROXY_SERVER $PROXY_PORT"
echo ""
