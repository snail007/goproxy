#!/bin/bash

# Run this script ON THE SERVER (194.233.73.36)
# Copy and paste these commands, or upload this script and run it

set -e

echo "=========================================="
echo "Game Proxy Server Deployment"
echo "=========================================="

# Update system
echo "Step 1: Updating system..."
apt-get update -qq
apt-get install -y wget curl git ufw iptables net-tools

# Install Go
echo "Step 2: Installing Go..."
if ! command -v go &> /dev/null; then
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
fi
export PATH=$PATH:/usr/local/go/bin
go version

# Create directories
echo "Step 3: Creating directories..."
mkdir -p /opt/game-proxy/{bin,config,logs}

# Extract source files (assumes proxy-deploy.tar.gz is in /tmp)
echo "Step 4: Extracting source files..."
if [ -f /tmp/proxy-deploy.tar.gz ]; then
    cd /tmp
    tar -xzf proxy-deploy.tar.gz
    mv stealth game services main.go config.go go.mod /opt/game-proxy/ 2>/dev/null || true
    cd /opt/game-proxy
else
    echo "ERROR: proxy-deploy.tar.gz not found in /tmp"
    echo "Please upload it first: scp proxy-deploy.tar.gz root@194.233.73.36:/tmp/"
    exit 1
fi

# Build proxy
echo "Step 5: Building proxy..."
export PATH=$PATH:/usr/local/go/bin
cd /opt/game-proxy
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy
ls -lh bin/proxy

# Create systemd service
echo "Step 6: Creating systemd service..."
cat > /etc/systemd/system/game-proxy.service <<'EOF'
[Unit]
Description=Game Proxy Server with Stealth Features
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/game-proxy
ExecStart=/opt/game-proxy/bin/proxy game-http -t tcp -p :8080 --forever --log /opt/game-proxy/logs/proxy.log
Restart=always
RestartSec=5
StandardOutput=append:/opt/game-proxy/logs/proxy.log
StandardError=append:/opt/game-proxy/logs/proxy-error.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable game-proxy

# Configure firewall
echo "Step 7: Configuring firewall..."
ufw --force enable
ufw allow 8080/tcp
ufw allow 1080/tcp
ufw allow 5353/udp
ufw allow 22/tcp
ufw status numbered | head -10

# Configure IPv6
echo "Step 8: Configuring IPv6..."
IFACE=$(ip link show | grep -E "eth0|ens3" | head -1 | cut -d: -f2 | tr -d ' ')
if [ -n "$IFACE" ]; then
    ip -6 addr add 2407:3640:2299:7804::1/64 dev $IFACE 2>/dev/null || echo "IPv6 may already be configured"
    ip -6 addr show | grep "2407:3640" || echo "IPv6 configuration check"
fi

# Start service
echo "Step 9: Starting proxy service..."
systemctl start game-proxy
sleep 3
systemctl status game-proxy --no-pager | head -20

# Test
echo "Step 10: Testing proxy..."
sleep 2
curl -s --connect-timeout 5 http://127.0.0.1:8080 http://httpbin.org/ip || echo "Proxy test - check logs if failed"

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo "Proxy: http://194.233.73.36:8080"
echo "IPv6: http://[2407:3640:2299:7804::1]:8080"
echo ""
echo "Check status: systemctl status game-proxy"
echo "View logs: tail -f /opt/game-proxy/logs/proxy.log"
echo ""
