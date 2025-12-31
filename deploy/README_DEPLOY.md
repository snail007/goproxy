# Server Deployment Guide

## Server Information
- **IPv4**: 194.233.73.36
- **IPv6**: 2407:3640:2299:7804::1/64
- **User**: root
- **Password**: JVxgREA3hbe6
- **Specs**: 4 cores, 8GB RAM

## Quick Deployment

### Option 1: Automated Deployment (Recommended)

```bash
cd deploy
chmod +x deploy.sh
./deploy.sh
```

This will:
1. Update system packages
2. Install Go and dependencies
3. Build and deploy proxy
4. Configure firewall
5. Set up systemd service
6. Configure IPv6
7. Start the proxy service

### Option 2: Manual Deployment

#### Step 1: Connect to Server
```bash
ssh root@194.233.73.36
# Password: JVxgREA3hbe6
```

#### Step 2: Update System
```bash
apt-get update && apt-get upgrade -y
apt-get install -y wget curl git ufw iptables
```

#### Step 3: Install Go
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
rm go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

#### Step 4: Deploy Proxy
```bash
mkdir -p /opt/game-proxy/{bin,config,logs}
# Copy proxy binary and source files to /opt/game-proxy/
chmod +x /opt/game-proxy/bin/proxy
```

#### Step 5: Configure Firewall
```bash
ufw allow 8080/tcp
ufw allow 1080/tcp
ufw allow 5353/udp
ufw allow 22/tcp
ufw --force enable
```

#### Step 6: Setup Systemd Service
```bash
cp deploy/systemd/game-proxy.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable game-proxy
systemctl start game-proxy
```

#### Step 7: Configure IPv6
```bash
ip -6 addr add 2407:3640:2299:7804::1/64 dev eth0
# or
ip -6 addr add 2407:3640:2299:7804::1/64 dev ens3
```

## Service Management

### Check Status
```bash
systemctl status game-proxy
```

### View Logs
```bash
# Real-time logs
tail -f /opt/game-proxy/logs/proxy.log

# Error logs
tail -f /opt/game-proxy/logs/proxy-error.log
```

### Restart Service
```bash
systemctl restart game-proxy
```

### Stop Service
```bash
systemctl stop game-proxy
```

## Proxy Endpoints

### HTTP/HTTPS Proxy
- **IPv4**: `194.233.73.36:8080`
- **IPv6**: `[2407:3640:2299:7804::1]:8080`

### SOCKS5 Proxy (if enabled)
- **IPv4**: `194.233.73.36:1080`
- **IPv6**: `[2407:3640:2299:7804::1]:1080`

### UDP Proxy (if enabled)
- **IPv4**: `194.233.73.36:5353`
- **IPv6**: `[2407:3640:2299:7804::1]:5353`

## Testing

### Test HTTP Proxy
```bash
# From your local machine
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# With IPv6
curl -x http://[2407:3640:2299:7804::1]:8080 http://httpbin.org/ip
```

### Test from Server
```bash
curl -x http://127.0.0.1:8080 http://httpbin.org/ip
```

### Check Proxy Headers
```bash
curl -x http://194.233.73.36:8080 http://httpbin.org/headers
# Should show randomized headers, no proxy signatures
```

## Configuration

### Edit Proxy Configuration
```bash
nano /opt/game-proxy/config/proxy-config.json
systemctl restart game-proxy
```

### Enable Authentication (Optional)
Edit `/etc/systemd/system/game-proxy.service`:
```ini
ExecStart=/opt/game-proxy/bin/proxy game-http \
    -t tcp \
    -p :8080 \
    -a "username:password" \
    --forever \
    --log /opt/game-proxy/logs/proxy.log
```

### Enable Parent Proxy (Optional)
```ini
ExecStart=/opt/game-proxy/bin/proxy game-http \
    -t tcp \
    -p :8080 \
    -T tls \
    -P "parent.proxy.com:443" \
    -C /path/to/cert.crt \
    -K /path/to/key.key \
    --forever \
    --log /opt/game-proxy/logs/proxy.log
```

## Security Hardening

### 1. Change Default Password
```bash
passwd root
```

### 2. Setup SSH Key Authentication
```bash
# On your local machine
ssh-copy-id root@194.233.73.36

# On server, disable password auth
nano /etc/ssh/sshd_config
# Set: PasswordAuthentication no
systemctl restart sshd
```

### 3. Enable Fail2Ban
```bash
apt-get install -y fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### 4. Setup Firewall Rules
```bash
# Allow only necessary ports
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp
ufw allow 8080/tcp
ufw allow 1080/tcp
ufw allow 5353/udp
ufw enable
```

## Monitoring

### Health Check Script
The deployment includes a health check script that runs every 5 minutes:
```bash
/opt/game-proxy/monitor.sh
```

### View Health Logs
```bash
tail -f /opt/game-proxy/logs/health.log
```

### System Resources
```bash
# CPU and Memory
htop

# Network
iftop

# Disk
df -h
```

## Troubleshooting

### Proxy Not Starting
```bash
# Check logs
journalctl -u game-proxy -n 50

# Check if port is in use
netstat -tulpn | grep 8080

# Check permissions
ls -la /opt/game-proxy/bin/proxy
```

### Connection Issues
```bash
# Test connectivity
ping 194.233.73.36

# Check firewall
ufw status

# Check iptables
iptables -L -n
```

### IPv6 Not Working
```bash
# Check IPv6 address
ip -6 addr show

# Test IPv6 connectivity
ping6 2407:3640:2299:7804::1

# Check IPv6 firewall
ip6tables -L -n
```

## Performance Tuning

### Increase File Descriptors
```bash
echo "* soft nofile 65535" >> /etc/security/limits.conf
echo "* hard nofile 65535" >> /etc/security/limits.conf
```

### Optimize Network
```bash
# Add to /etc/sysctl.conf
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
net.ipv4.tcp_congestion_control = bbr

sysctl -p
```

## Backup

### Backup Configuration
```bash
tar -czf game-proxy-backup-$(date +%Y%m%d).tar.gz /opt/game-proxy
```

### Restore Configuration
```bash
tar -xzf game-proxy-backup-YYYYMMDD.tar.gz -C /
systemctl restart game-proxy
```

## Support

For issues:
1. Check logs: `/opt/game-proxy/logs/`
2. Check systemd: `systemctl status game-proxy`
3. Review firewall: `ufw status`
4. Test connectivity: `curl -x http://127.0.0.1:8080 http://httpbin.org/ip`
