#!/bin/bash

# SSH deployment script using expect
# Server: 194.233.73.36

SERVER_IP="194.233.73.36"
SERVER_USER="root"
SERVER_PASS="JVxgREA3hbe6"
SERVER_IPV6="2407:3640:2299:7804::1"

# Check if expect is installed
if ! command -v expect &> /dev/null; then
    echo "Installing expect..."
    sudo apt-get update && sudo apt-get install -y expect
fi

# Create expect script for SSH commands
expect_ssh() {
    expect <<EOF
set timeout 30
spawn ssh -o StrictHostKeyChecking=no ${SERVER_USER}@${SERVER_IP}
expect {
    "password:" {
        send "${SERVER_PASS}\r"
        exp_continue
    }
    "# " {
        send "$1\r"
        expect "# "
        send "exit\r"
    }
    timeout {
        puts "Connection timeout"
        exit 1
    }
}
expect eof
EOF
}

# Create expect script for SCP
expect_scp() {
    expect <<EOF
set timeout 60
spawn scp -o StrictHostKeyChecking=no $1 ${SERVER_USER}@${SERVER_IP}:$2
expect {
    "password:" {
        send "${SERVER_PASS}\r"
        exp_continue
    }
    "100%" {
        puts "File uploaded successfully"
    }
    timeout {
        puts "Upload timeout"
        exit 1
    }
}
expect eof
EOF
}

echo "=========================================="
echo "Deploying to Server: $SERVER_IP"
echo "=========================================="

# Test connection
echo "Testing connection..."
expect_ssh "echo 'Connected successfully'; uname -a"

# Update system
echo "Updating system..."
expect_ssh "apt-get update -qq && apt-get install -y wget curl git ufw iptables net-tools"

# Install Go
echo "Installing Go..."
expect_ssh "export PATH=\$PATH:/usr/local/go/bin 2>/dev/null; if ! command -v go &> /dev/null; then wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz && tar -C /usr/local -xzf /tmp/go.tar.gz && rm /tmp/go.tar.gz && echo 'export PATH=\$PATH:/usr/local/go/bin' >> /root/.bashrc; fi; /usr/local/go/bin/go version 2>/dev/null || go version"

# Create directories
echo "Creating directories..."
expect_ssh "mkdir -p /opt/game-proxy/{bin,config,logs}"

# Upload files
echo "Uploading source files..."
cd /workspace
tar -czf /tmp/proxy-deploy.tar.gz stealth/ game/ services/game_http.go main.go config.go go.mod 2>/dev/null
expect_scp "/tmp/proxy-deploy.tar.gz" "/tmp/"

# Extract and setup
echo "Setting up files..."
expect_ssh "cd /tmp && tar -xzf proxy-deploy.tar.gz && mv stealth game services main.go config.go go.mod /opt/game-proxy/ 2>&1"

# Build proxy
echo "Building proxy..."
expect_ssh "cd /opt/game-proxy && export PATH=\$PATH:/usr/local/go/bin && go mod tidy && go build -o bin/proxy -ldflags='-s -w' . && chmod +x bin/proxy && ls -lh bin/proxy"

# Setup systemd service
echo "Configuring systemd service..."
expect_scp "deploy/systemd/game-proxy.service" "/tmp/game-proxy.service"
expect_ssh "mv /tmp/game-proxy.service /etc/systemd/system/ && systemctl daemon-reload && systemctl enable game-proxy"

# Configure firewall
echo "Configuring firewall..."
expect_ssh "ufw --force enable; ufw allow 8080/tcp; ufw allow 1080/tcp; ufw allow 5353/udp; ufw allow 22/tcp; ufw status numbered | head -10"

# Configure IPv6
echo "Configuring IPv6..."
expect_ssh "ip -6 addr show | grep '2407:3640' || (IFACE=\$(ip link show | grep -E 'eth0|ens3' | head -1 | cut -d: -f2 | tr -d ' '); ip -6 addr add ${SERVER_IPV6}/64 dev \$IFACE 2>&1); ip -6 addr show | grep '2407:3640'"

# Start service
echo "Starting proxy service..."
expect_ssh "systemctl start game-proxy && sleep 3 && systemctl status game-proxy --no-pager | head -20"

# Test connection
echo "Testing proxy..."
sleep 2
curl -s --connect-timeout 5 -x http://${SERVER_IP}:8080 http://httpbin.org/ip || echo "Proxy test failed - check logs"

echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo "Proxy: http://${SERVER_IP}:8080"
echo "IPv6: http://[${SERVER_IPV6}]:8080"
echo ""
echo "Check status: ssh ${SERVER_USER}@${SERVER_IP} 'systemctl status game-proxy'"
echo "View logs: ssh ${SERVER_USER}@${SERVER_IP} 'tail -f /opt/game-proxy/logs/proxy.log'"
echo ""
