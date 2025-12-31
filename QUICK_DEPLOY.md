# Quick Deployment Guide

## Server: 194.233.73.36

### Fastest Method: Copy-Paste Commands

**1. Upload source files:**
```bash
# From your local machine
cd /workspace
scp proxy-deploy.tar.gz root@194.233.73.36:/tmp/
scp deploy/deploy-on-server.sh root@194.233.73.36:/tmp/
```

**2. SSH to server:**
```bash
ssh root@194.233.73.36
# Password: JVxgREA3hbe6
```

**3. Run deployment script:**
```bash
chmod +x /tmp/deploy-on-server.sh
/tmp/deploy-on-server.sh
```

### Or Manual Step-by-Step:

**After SSH to server, run:**

```bash
# 1. Update & install
apt-get update && apt-get install -y wget curl git ufw iptables net-tools

# 2. Install Go
wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz
tar -C /usr/local -xzf /tmp/go.tar.gz && rm /tmp/go.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# 3. Extract source (after uploading proxy-deploy.tar.gz)
cd /tmp && tar -xzf proxy-deploy.tar.gz
mkdir -p /opt/game-proxy/{bin,config,logs}
mv stealth game services main.go config.go go.mod /opt/game-proxy/

# 4. Build
cd /opt/game-proxy
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy

# 5. Create service
cat > /etc/systemd/system/game-proxy.service <<'EOF'
[Unit]
Description=Game Proxy Server
After=network.target
[Service]
Type=simple
WorkingDirectory=/opt/game-proxy
ExecStart=/opt/game-proxy/bin/proxy game-http -t tcp -p :8080 --forever --log /opt/game-proxy/logs/proxy.log
Restart=always
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable game-proxy

# 6. Firewall
ufw allow 8080/tcp && ufw allow 22/tcp && ufw --force enable

# 7. Start
systemctl start game-proxy

# 8. Test
curl -x http://127.0.0.1:8080 http://httpbin.org/ip
```

## Files Ready

✅ `proxy-deploy.tar.gz` - Source code (7.8KB)  
✅ `deploy/deploy-on-server.sh` - Deployment script  
✅ `deploy/MANUAL_DEPLOY.md` - Detailed guide  

## Test After Deployment

From your local machine:
```bash
curl -x http://194.233.73.36:8080 http://httpbin.org/ip
```

Should return: `{"origin":"194.233.73.36"}`
