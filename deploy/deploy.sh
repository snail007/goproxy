#!/bin/bash

# Deployment script for Game Proxy Server
# Server: 194.233.73.36
# IPv6: 2407:3640:2299:7804::1/64

set -e

SERVER_IP="194.233.73.36"
SERVER_IPV6="2407:3640:2299:7804::1"
SERVER_USER="root"
SERVER_PASS="JVxgREA3hbe6"
PROXY_PORT="8080"
SOCKS5_PORT="1080"
UDP_PORT="5353"

echo "=========================================="
echo "Game Proxy Server Deployment"
echo "=========================================="
echo "Server: $SERVER_IP"
echo "IPv6: $SERVER_IPV6"
echo "=========================================="

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if sshpass is installed
if ! command -v sshpass &> /dev/null; then
    echo -e "${YELLOW}Installing sshpass...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install hudochenkov/sshpass/sshpass
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get update && sudo apt-get install -y sshpass
    fi
fi

# Function to execute remote commands
remote_exec() {
    sshpass -p "$SERVER_PASS" ssh -o StrictHostKeyChecking=no "$SERVER_USER@$SERVER_IP" "$1"
}

# Function to copy files
remote_copy() {
    sshpass -p "$SERVER_PASS" scp -o StrictHostKeyChecking=no "$1" "$SERVER_USER@$SERVER_IP:$2"
}

echo -e "${BLUE}Step 1: Updating system...${NC}"
remote_exec "apt-get update && apt-get upgrade -y"

echo -e "${BLUE}Step 2: Installing dependencies...${NC}"
remote_exec "apt-get install -y wget curl git ufw iptables"

echo -e "${BLUE}Step 3: Installing Go...${NC}"
remote_exec "wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz && \
    rm go1.21.5.linux-amd64.tar.gz && \
    echo 'export PATH=\$PATH:/usr/local/go/bin' >> /root/.bashrc"

echo -e "${BLUE}Step 4: Creating proxy directory...${NC}"
remote_exec "mkdir -p /opt/game-proxy/{bin,config,logs}"

echo -e "${BLUE}Step 5: Copying source files to server...${NC}"
# Create temporary directory for source files
mkdir -p /tmp/game-proxy-deploy
cp -r ../stealth /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: stealth directory not found, will build on server"
cp -r ../game /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: game directory not found, will build on server"
cp -r ../services /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: services directory not found, will build on server"
cp ../main.go /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: main.go not found"
cp ../config.go /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: config.go not found"
cp ../go.mod /tmp/game-proxy-deploy/ 2>/dev/null || echo "Note: go.mod not found"

# Copy source files
remote_exec "mkdir -p /tmp/game-proxy-source"
if [ -d "/tmp/game-proxy-deploy/stealth" ]; then
    remote_copy "/tmp/game-proxy-deploy/stealth" "/tmp/game-proxy-source/"
fi
if [ -d "/tmp/game-proxy-deploy/game" ]; then
    remote_copy "/tmp/game-proxy-deploy/game" "/tmp/game-proxy-source/"
fi
if [ -d "/tmp/game-proxy-deploy/services" ]; then
    remote_copy "/tmp/game-proxy-deploy/services" "/tmp/game-proxy-source/"
fi
if [ -f "/tmp/game-proxy-deploy/main.go" ]; then
    remote_copy "/tmp/game-proxy-deploy/main.go" "/tmp/game-proxy-source/"
fi
if [ -f "/tmp/game-proxy-deploy/config.go" ]; then
    remote_copy "/tmp/game-proxy-deploy/config.go" "/tmp/game-proxy-source/"
fi
if [ -f "/tmp/game-proxy-deploy/go.mod" ]; then
    remote_copy "/tmp/game-proxy-deploy/go.mod" "/tmp/game-proxy-source/"
fi

echo -e "${BLUE}Step 6: Building proxy on server...${NC}"
remote_exec "cd /tmp/game-proxy-source && \
    export PATH=\$PATH:/usr/local/go/bin && \
    go mod tidy && \
    go build -o /opt/game-proxy/bin/proxy -ldflags='-s -w' . && \
    chmod +x /opt/game-proxy/bin/proxy && \
    rm -rf /tmp/game-proxy-source"

echo -e "${BLUE}Step 7: Setting up configuration...${NC}"
remote_exec "chmod +x /opt/game-proxy/bin/proxy"

echo -e "${BLUE}Step 8: Configuring firewall...${NC}"
remote_exec "ufw allow $PROXY_PORT/tcp"
remote_exec "ufw allow $SOCKS5_PORT/tcp"
remote_exec "ufw allow $UDP_PORT/udp"
remote_exec "ufw allow 22/tcp"
remote_exec "ufw --force enable"

echo -e "${BLUE}Step 9: Creating systemd service...${NC}"
cat > /tmp/game-proxy.service <<EOF
[Unit]
Description=Game Proxy Server with Stealth Features
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/game-proxy
ExecStart=/opt/game-proxy/bin/proxy game-http -t tcp -p :$PROXY_PORT --forever --log /opt/game-proxy/logs/proxy.log
Restart=always
RestartSec=5
StandardOutput=append:/opt/game-proxy/logs/proxy.log
StandardError=append:/opt/game-proxy/logs/proxy-error.log

[Install]
WantedBy=multi-user.target
EOF

remote_copy "/tmp/game-proxy.service" "/tmp/game-proxy.service"
remote_exec "mv /tmp/game-proxy.service /etc/systemd/system/game-proxy.service"
remote_exec "systemctl daemon-reload"
remote_exec "systemctl enable game-proxy"
remote_exec "systemctl start game-proxy"

echo -e "${BLUE}Step 10: Configuring IPv6...${NC}"
remote_exec "ip -6 addr add $SERVER_IPV6/64 dev eth0 || ip -6 addr add $SERVER_IPV6/64 dev ens3 || true"

echo -e "${BLUE}Step 11: Checking service status...${NC}"
sleep 3
remote_exec "systemctl status game-proxy --no-pager"

echo ""
echo -e "${GREEN}=========================================="
echo "Deployment Complete!"
echo "==========================================${NC}"
echo -e "${GREEN}Proxy Server:${NC}"
echo "  HTTP/HTTPS: $SERVER_IP:$PROXY_PORT"
echo "  IPv6: [$SERVER_IPV6]:$PROXY_PORT"
echo ""
echo -e "${GREEN}Management:${NC}"
echo "  Status: ssh $SERVER_USER@$SERVER_IP 'systemctl status game-proxy'"
echo "  Logs: ssh $SERVER_USER@$SERVER_IP 'tail -f /opt/game-proxy/logs/proxy.log'"
echo "  Restart: ssh $SERVER_USER@$SERVER_IP 'systemctl restart game-proxy'"
echo ""
echo -e "${GREEN}Test Connection:${NC}"
echo "  curl -x http://$SERVER_IP:$PROXY_PORT http://httpbin.org/ip"
echo ""
