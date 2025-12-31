# Server Information & Configuration

## Server Details

### Basic Information
- **IPv4 Address**: 194.233.73.36
- **IPv6 Address**: 2407:3640:2299:7804::1/64
- **Username**: root
- **Password**: JVxgREA3hbe6
- **Specifications**: 4 CPU cores, 8GB RAM

### Network Configuration

#### IPv4
- **IP**: 194.233.73.36
- **Gateway**: (Auto-configured)
- **Netmask**: (Auto-configured)

#### IPv6
- **IP**: 2407:3640:2299:7804::1
- **Prefix**: /64
- **Full Address**: 2407:3640:2299:7804:0000:0000:0000:0001/64

## Proxy Configuration

### Ports
- **HTTP/HTTPS Proxy**: 8080
- **SOCKS5 Proxy**: 1080 (if enabled)
- **UDP Proxy**: 5353 (if enabled)

### Endpoints

#### IPv4
```
HTTP Proxy:  http://194.233.73.36:8080
HTTPS Proxy: http://194.233.73.36:8080
SOCKS5:      socks5://194.233.73.36:1080
```

#### IPv6
```
HTTP Proxy:  http://[2407:3640:2299:7804::1]:8080
HTTPS Proxy: http://[2407:3640:2299:7804::1]:8080
SOCKS5:      socks5://[2407:3640:2299:7804::1]:1080
```

## Connection Examples

### curl
```bash
# IPv4
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# IPv6
curl -x http://[2407:3640:2299:7804::1]:8080 http://httpbin.org/ip
```

### wget
```bash
# IPv4
wget -e use_proxy=yes -e http_proxy=194.233.73.36:8080 http://example.com

# IPv6
wget -e use_proxy=yes -e http_proxy=[2407:3640:2299:7804::1]:8080 http://example.com
```

### Browser Settings
```
HTTP Proxy:  194.233.73.36:8080
HTTPS Proxy: 194.233.73.36:8080
SOCKS5:      194.233.73.36:1080
```

### Mac App Configuration
Update `ProxyManager.swift`:
```swift
let parentProxy = "http://194.233.73.36:8080"
```

## SSH Access

### Connect
```bash
ssh root@194.233.73.36
# Password: JVxgREA3hbe6
```

### Setup SSH Keys (Recommended)
```bash
# On your local machine
ssh-copy-id root@194.233.73.36

# Then disable password auth on server
ssh root@194.233.73.36
nano /etc/ssh/sshd_config
# Set: PasswordAuthentication no
systemctl restart sshd
```

## Security Checklist

- [ ] Change root password after deployment
- [ ] Setup SSH key authentication
- [ ] Configure firewall (UFW)
- [ ] Enable fail2ban
- [ ] Setup log rotation
- [ ] Configure monitoring
- [ ] Regular security updates

## Monitoring

### Check Proxy Status
```bash
ssh root@194.233.73.36 'systemctl status game-proxy'
```

### View Logs
```bash
ssh root@194.233.73.36 'tail -f /opt/game-proxy/logs/proxy.log'
```

### Check Resources
```bash
ssh root@194.233.73.36 'htop'
ssh root@194.233.73.36 'df -h'
ssh root@194.233.73.36 'free -h'
```

## Backup

### Backup Configuration
```bash
ssh root@194.233.73.36 'tar -czf /root/game-proxy-backup.tar.gz /opt/game-proxy'
scp root@194.233.73.36:/root/game-proxy-backup.tar.gz ./
```

## Troubleshooting

### Can't Connect
1. Check firewall: `ssh root@194.233.73.36 'ufw status'`
2. Check service: `ssh root@194.233.73.36 'systemctl status game-proxy'`
3. Check logs: `ssh root@194.233.73.36 'tail -50 /opt/game-proxy/logs/proxy.log'`

### IPv6 Not Working
1. Check IPv6 address: `ssh root@194.233.73.36 'ip -6 addr show'`
2. Test connectivity: `ping6 2407:3640:2299:7804::1`
3. Check firewall: `ssh root@194.233.73.36 'ip6tables -L -n'`

## Quick Commands Reference

```bash
# Deploy
cd deploy && ./deploy.sh

# Test
./deploy/test-connection.sh

# Status
ssh root@194.233.73.36 'systemctl status game-proxy'

# Restart
ssh root@194.233.73.36 'systemctl restart game-proxy'

# Logs
ssh root@194.233.73.36 'tail -f /opt/game-proxy/logs/proxy.log'
```
