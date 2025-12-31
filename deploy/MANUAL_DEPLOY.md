# Manual Deployment Instructions

Since automated deployment tools may not be available, here are manual steps to deploy to your server.

## Server Information
- **IP**: 194.233.73.36
- **IPv6**: 2407:3640:2299:7804::1
- **User**: root
- **Password**: JVxgREA3hbe6

## Step-by-Step Deployment

### 1. Connect to Server
```bash
ssh root@194.233.73.36
# Password: JVxgREA3hbe6
```

### 2. Update System
```bash
apt-get update
apt-get install -y wget curl git ufw iptables net-tools
```

### 3. Install Go
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
rm go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
go version
```

### 4. Create Directories
```bash
mkdir -p /opt/game-proxy/{bin,config,logs}
```

### 5. Upload Source Files

**From your local machine:**
```bash
cd /workspace
tar -czf proxy-deploy.tar.gz stealth/ game/ services/game_http.go main.go config.go go.mod
scp proxy-deploy.tar.gz root@194.233.73.36:/tmp/
```

**On server:**
```bash
cd /tmp
tar -xzf proxy-deploy.tar.gz
mv stealth game services main.go config.go go.mod /opt/game-proxy/
cd /opt/game-proxy
```

### 6. Build Proxy
```bash
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy
ls -lh bin/proxy
```

### 7. Create Systemd Service
```bash
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

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable game-proxy
```

### 8. Configure Firewall
```bash
ufw allow 8080/tcp
ufw allow 1080/tcp
ufw allow 5353/udp
ufw allow 22/tcp
ufw --force enable
ufw status
```

### 9. Configure IPv6
```bash
# Find network interface
ip link show

# Add IPv6 address (replace eth0 with your interface name)
ip -6 addr add 2407:3640:2299:7804::1/64 dev eth0
# or
ip -6 addr add 2407:3640:2299:7804::1/64 dev ens3

# Verify
ip -6 addr show | grep 2407:3640
```

### 10. Start Service
```bash
systemctl start game-proxy
systemctl status game-proxy
```

### 11. Test Connection
```bash
# From server
curl -x http://127.0.0.1:8080 http://httpbin.org/ip

# From your local machine
curl -x http://194.233.73.36:8080 http://httpbin.org/ip
```

## Troubleshooting

### Check Service Status
```bash
systemctl status game-proxy
journalctl -u game-proxy -n 50
```

### Check Logs
```bash
tail -f /opt/game-proxy/logs/proxy.log
```

### Check Port
```bash
ss -tulpn | grep 8080
netstat -tulpn | grep 8080
```

### Restart Service
```bash
systemctl restart game-proxy
```

## Quick Commands Reference

```bash
# Status
systemctl status game-proxy

# Logs
tail -f /opt/game-proxy/logs/proxy.log

# Restart
systemctl restart game-proxy

# Stop
systemctl stop game-proxy

# Start
systemctl start game-proxy
```

## Verify Deployment

1. **Service Running**: `systemctl is-active game-proxy` should return "active"
2. **Port Listening**: `ss -tulpn | grep 8080` should show port 8080
3. **Proxy Working**: `curl -x http://127.0.0.1:8080 http://httpbin.org/ip` should return IP

## Next Steps

After deployment:
1. Test connection from your Mac app
2. Configure client connections
3. Monitor logs for any issues
4. Change root password for security
