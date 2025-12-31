# Quick Start - Deploy to Your Server

## Server Details
- **IP**: 194.233.73.36
- **IPv6**: 2407:3640:2299:7804::1
- **User**: root
- **Password**: JVxgREA3hbe6

## One-Command Deployment

```bash
cd deploy
./deploy.sh
```

## What Gets Installed

✅ Go 1.21.5  
✅ Game Proxy with Stealth Features  
✅ Systemd Service (auto-start)  
✅ Firewall Configuration  
✅ IPv6 Support  
✅ Health Monitoring  

## After Deployment

### Test Your Proxy
```bash
# Test IPv4
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# Test IPv6
curl -x http://[2407:3640:2299:7804::1]:8080 http://httpbin.org/ip
```

### Check Status
```bash
ssh root@194.233.73.36 'systemctl status game-proxy'
```

### View Logs
```bash
ssh root@194.233.73.36 'tail -f /opt/game-proxy/logs/proxy.log'
```

## Use in Mac App

Update `ProxyManager.swift` with your server:
```swift
func startProxy(port: Int = 8080, parentProxy: String? = nil) {
    // Use your server as parent proxy
    let parentProxy = "http://194.233.73.36:8080"
    // ... rest of code
}
```

## Proxy Endpoints

- **HTTP/HTTPS**: `194.233.73.36:8080`
- **IPv6**: `[2407:3640:2299:7804::1]:8080`

## Security Note

⚠️ **Important**: Change the root password after deployment!

```bash
ssh root@194.233.73.36
passwd root
```

## Next Steps

1. ✅ Deploy proxy server
2. ✅ Test connection
3. ✅ Configure Mac app to use server
4. ✅ Change server password
5. ✅ Setup SSH keys (optional but recommended)
