# Deployment Ready! 🚀

## Quick Deploy to Your Server

Your server is ready for deployment. Here are the options:

### Option 1: Manual Deployment (Recommended)

**Step 1: Upload files**
```bash
# From your local machine
cd /workspace
tar -czf proxy-deploy.tar.gz stealth/ game/ services/game_http.go main.go config.go go.mod
scp proxy-deploy.tar.gz root@194.233.73.36:/tmp/
```

**Step 2: SSH and deploy**
```bash
ssh root@194.233.73.36
# Password: JVxgREA3hbe6

# Then run these commands:
cd /tmp && tar -xzf proxy-deploy.tar.gz
mkdir -p /opt/game-proxy/{bin,config,logs}
mv stealth game services main.go config.go go.mod /opt/game-proxy/
cd /opt/game-proxy

# Install Go if needed
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
rm go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Build
go mod tidy
go build -o bin/proxy -ldflags="-s -w" .
chmod +x bin/proxy

# Setup service
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
systemctl start game-proxy

# Configure firewall
ufw allow 8080/tcp
ufw allow 22/tcp
ufw --force enable

# Test
curl -x http://127.0.0.1:8080 http://httpbin.org/ip
```

### Option 2: Use Deployment Script

The file `proxy-deploy.tar.gz` is ready in `/workspace/`

Copy it to server and extract:
```bash
scp proxy-deploy.tar.gz root@194.233.73.36:/tmp/
ssh root@194.233.73.36
cd /tmp && tar -xzf proxy-deploy.tar.gz
# Follow manual steps above
```

## Files Ready for Deployment

✅ Source code packaged: `proxy-deploy.tar.gz`  
✅ Deployment scripts: `deploy/` directory  
✅ Documentation: `deploy/MANUAL_DEPLOY.md`  

## Server Details

- **IPv4**: 194.233.73.36:8080
- **IPv6**: [2407:3640:2299:7804::1]:8080
- **User**: root
- **Password**: JVxgREA3hbe6

## After Deployment

Test your proxy:
```bash
curl -x http://194.233.73.36:8080 http://httpbin.org/ip
```

Should return your server's IP address!

## Need Help?

See `deploy/MANUAL_DEPLOY.md` for detailed step-by-step instructions.
