#!/bin/bash
# Complete deployment script - run this on the server

set -e

echo "=========================================="
echo "Game Proxy Deployment"
echo "=========================================="

# Step 1: Install dependencies
echo "Step 1: Installing dependencies..."
apt-get update -qq
apt-get install -y wget curl git ufw iptables net-tools

# Step 2: Install Go
echo "Step 2: Installing Go..."
if ! command -v go > /dev/null 2>&1; then
    wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
fi
export PATH=$PATH:/usr/local/go/bin
go version

# Step 3: Extract files
echo "Step 3: Extracting source files..."
cd /tmp
if [ ! -f proxy-deploy.tar.gz ]; then
    echo "ERROR: proxy-deploy.tar.gz not found!"
    exit 1
fi
tar -xzf proxy-deploy.tar.gz
mkdir -p /opt/game-proxy/{bin,config,logs}
mv stealth game services main.go config.go go.mod /opt/game-proxy/ 2>/dev/null || true
cd /opt/game-proxy

# Step 4: Build
echo "Step 4: Building proxy..."
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy
ls -lh bin/proxy

# Step 5: Create systemd service
echo "Step 5: Creating systemd service..."
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

# Step 6: Configure firewall
echo "Step 6: Configuring firewall..."
ufw allow 8080/tcp
ufw allow 1080/tcp
ufw allow 5353/udp
ufw allow 22/tcp
ufw --force enable
ufw status numbered | head -10

# Step 7: Start service
echo "Step 7: Starting proxy service..."
systemctl start game-proxy
sleep 3

# Step 8: Check status
echo "Step 8: Checking status..."
systemctl status game-proxy --no-pager | head -15

# Step 9: Test
echo "Step 9: Testing proxy..."
sleep 2
curl -s --connect-timeout 5 http://127.0.0.1:8080 http://httpbin.org/ip || echo "Test failed - check logs"

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo "Proxy: http://194.233.73.36:8080"
echo "Check logs: tail -f /opt/game-proxy/logs/proxy.log"
echo ""
