#!/bin/bash
# Deployment script executed on remote server via Terraform

set -e

echo "=========================================="
echo "Terraform: Deploying Game Proxy"
echo "=========================================="

# Step 1: Install dependencies
echo "Step 1: Installing dependencies..."
apt-get update -qq
apt-get install -y wget curl git ufw iptables net-tools

# Step 2: Install Go
echo "Step 2: Installing Go..."
if ! command -v go > /dev/null 2>&1; then
    echo "Downloading Go..."
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

# Step 6: Start service
echo "Step 6: Starting proxy service..."
systemctl stop game-proxy 2>/dev/null || true
systemctl start game-proxy
sleep 3

# Step 7: Check status
echo "Step 7: Checking status..."
systemctl status game-proxy --no-pager | head -15 || true

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo "Proxy: http://194.233.73.36:8080"
echo ""
