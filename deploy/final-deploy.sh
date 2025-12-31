#!/bin/bash
set -e

echo "=========================================="
echo "Final Game Proxy Deployment"
echo "=========================================="

cd /tmp
tar -xzf proxy-deploy.tar.gz
mkdir -p /opt/game-proxy/{bin,config,logs}
cp -r stealth game services utils main.go config.go go.mod /opt/game-proxy/

cd /opt/game-proxy
export PATH=$PATH:/usr/local/go/bin

echo "Building proxy..."
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy
ls -lh bin/proxy

echo "Creating systemd service..."
cat > /etc/systemd/system/game-proxy.service <<'EOF'
[Unit]
Description=Game Proxy Server
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/game-proxy
ExecStart=/opt/game-proxy/bin/proxy game-http -t tcp -p :8080 --forever --log /opt/game-proxy/logs/proxy.log
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable game-proxy
systemctl start game-proxy
sleep 3

echo "Checking status..."
systemctl status game-proxy --no-pager | head -15

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
