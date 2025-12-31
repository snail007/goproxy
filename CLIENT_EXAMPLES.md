# Client Connection Examples

## Quick Start

### Test Connection First
```bash
./test-proxy.sh
```

## Connection Methods

### 1. Mac App (Easiest)

**Option A: Use Built-in Quick Connect**
1. Launch GameProxyApp
2. Click menu bar icon
3. Select "Quick Connect"
4. Click "Connect" button

**Option B: Configure Manually**
1. Launch GameProxyApp
2. Click menu bar icon → Settings
3. Enter Parent Proxy: `http://194.233.73.36:8080`
4. Click "Start Proxy"

### 2. macOS System Proxy (All Apps)

```bash
# Set proxy for Wi-Fi
sudo networksetup -setwebproxy "Wi-Fi" 194.233.73.36 8080
sudo networksetup -setsecurewebproxy "Wi-Fi" 194.233.73.36 8080

# Enable
sudo networksetup -setwebproxystate "Wi-Fi" on
sudo networksetup -setsecurewebproxystate "Wi-Fi" on

# Test
curl http://httpbin.org/ip

# Disable when done
sudo networksetup -setwebproxystate "Wi-Fi" off
sudo networksetup -setsecurewebproxystate "Wi-Fi" off
```

### 3. Terminal Only (curl, wget, git)

```bash
# Set environment variables
export http_proxy=http://194.233.73.36:8080
export https_proxy=http://194.233.73.36:8080
export HTTP_PROXY=http://194.233.73.36:8080
export HTTPS_PROXY=http://194.233.73.36:8080

# Now all terminal commands use proxy
curl http://httpbin.org/ip
wget http://example.com
git clone https://github.com/user/repo.git

# Disable
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
```

### 4. Browser Only

#### Chrome/Edge
1. Settings → Advanced → System → Open proxy settings
2. Configure as per macOS System Proxy above

#### Firefox
1. Preferences → General → Network Settings
2. Manual proxy: `194.233.73.36:8080`

### 5. Specific Applications

#### Git
```bash
git config --global http.proxy http://194.233.73.36:8080
git config --global https.proxy http://194.233.73.36:8080
```

#### npm
```bash
npm config set proxy http://194.233.73.36:8080
npm config set https-proxy http://194.233.73.36:8080
```

#### Python
```python
import os
os.environ['http_proxy'] = 'http://194.233.73.36:8080'
os.environ['https_proxy'] = 'http://194.233.73.36:8080'
```

## Connection Test

Run the test script:
```bash
./test-proxy.sh
```

Or test manually:
```bash
# Basic test
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# Should return your server's IP, not your local IP
```

## Troubleshooting

### Can't Connect
1. Check server is running: `ssh root@194.233.73.36 'systemctl status game-proxy'`
2. Check firewall: `ssh root@194.233.73.36 'ufw status'`
3. Test connectivity: `ping 194.233.73.36`

### Connection Slow
1. Check latency: `ping 194.233.73.36`
2. Check server resources: `ssh root@194.233.73.36 'htop'`
3. Try IPv6: `curl -x http://[2407:3640:2299:7804::1]:8080 http://httpbin.org/ip`

### Proxy Not Working
1. Verify proxy settings: `networksetup -getwebproxy "Wi-Fi"`
2. Check logs: `ssh root@194.233.73.36 'tail -f /opt/game-proxy/logs/proxy.log'`
3. Restart proxy: `ssh root@194.233.73.36 'systemctl restart game-proxy'`

## Quick Reference Card

```
┌─────────────────────────────────────────┐
│  Proxy Server Connection               │
├─────────────────────────────────────────┤
│  IPv4: 194.233.73.36:8080              │
│  IPv6: [2407:3640:2299:7804::1]:8080   │
│                                         │
│  Quick Test:                            │
│  curl -x http://194.233.73.36:8080 \   │
│       http://httpbin.org/ip            │
│                                         │
│  System Proxy:                          │
│  networksetup -setwebproxy \           │
│    "Wi-Fi" 194.233.73.36 8080         │
│                                         │
│  Environment:                            │
│  export http_proxy=\                    │
│    http://194.233.73.36:8080          │
└─────────────────────────────────────────┘
```
