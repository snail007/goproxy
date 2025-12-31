# Client Connection Guide

## Proxy Server Information
- **IPv4**: 194.233.73.36:8080
- **IPv6**: [2407:3640:2299:7804::1]:8080
- **Protocol**: HTTP/HTTPS Proxy

## Quick Connection Examples

### 1. Mac App (Game Proxy App)

#### Update ProxyManager.swift
```swift
func startProxy(port: Int = 8080, parentProxy: String? = nil) {
    // Connect to your server
    let parentProxy = "http://194.233.73.36:8080"
    
    // Or use IPv6
    // let parentProxy = "http://[2407:3640:2299:7804::1]:8080"
    
    startProxy(port: port, parentProxy: parentProxy)
}
```

#### Manual Configuration in App
1. Launch GameProxyApp
2. Click menu bar icon
3. Select "Settings"
4. Enter parent proxy: `http://194.233.73.36:8080`
5. Click "Start Proxy"

### 2. Command Line (curl)

#### Basic Usage
```bash
# HTTP request through proxy
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# HTTPS request through proxy
curl -x http://194.233.73.36:8080 https://httpbin.org/ip

# With IPv6
curl -x http://[2407:3640:2299:7804::1]:8080 http://httpbin.org/ip
```

#### With Authentication (if enabled)
```bash
curl -x http://username:password@194.233.73.36:8080 http://httpbin.org/ip
```

#### Environment Variables
```bash
export http_proxy=http://194.233.73.36:8080
export https_proxy=http://194.233.73.36:8080
export HTTP_PROXY=http://194.233.73.36:8080
export HTTPS_PROXY=http://194.233.73.36:8080

# Now curl will use proxy automatically
curl http://httpbin.org/ip

# Disable proxy
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY
```

### 3. macOS System Proxy Settings

#### Via System Preferences
1. Open **System Preferences** → **Network**
2. Select your network interface (Wi-Fi/Ethernet)
3. Click **Advanced** → **Proxies** tab
4. Check **Web Proxy (HTTP)**
   - Server: `194.233.73.36`
   - Port: `8080`
5. Check **Secure Web Proxy (HTTPS)**
   - Server: `194.233.73.36`
   - Port: `8080`
6. Click **OK** → **Apply**

#### Via Command Line (networksetup)
```bash
# Set HTTP proxy
sudo networksetup -setwebproxy "Wi-Fi" 194.233.73.36 8080

# Set HTTPS proxy
sudo networksetup -setsecurewebproxy "Wi-Fi" 194.233.73.36 8080

# Enable proxy
sudo networksetup -setwebproxystate "Wi-Fi" on
sudo networksetup -setsecurewebproxystate "Wi-Fi" on

# Disable proxy
sudo networksetup -setwebproxystate "Wi-Fi" off
sudo networksetup -setsecurewebproxystate "Wi-Fi" off

# Check status
networksetup -getwebproxy "Wi-Fi"
networksetup -getsecurewebproxy "Wi-Fi"
```

#### PAC File (Automatic Configuration)
Create `proxy.pac`:
```javascript
function FindProxyForURL(url, host) {
    return "PROXY 194.233.73.36:8080; DIRECT";
}
```

Apply PAC file:
```bash
sudo networksetup -setautoproxyurl "Wi-Fi" file:///path/to/proxy.pac
```

### 4. Browser Configuration

#### Chrome/Edge
1. Open Settings → Advanced → System
2. Click "Open proxy settings"
3. Configure as per macOS System Proxy Settings above

#### Firefox
1. Open Preferences → General → Network Settings
2. Select "Manual proxy configuration"
3. HTTP Proxy: `194.233.73.36` Port: `8080`
4. HTTPS Proxy: `194.233.73.36` Port: `8080`
5. Check "Also use this proxy for HTTPS"
6. Click OK

#### Safari
Uses system proxy settings (configure via System Preferences)

### 5. Game Clients

#### Steam
1. Open Steam → Settings → Network
2. Check "Use proxy server"
3. Proxy server: `194.233.73.36:8080`
4. Click OK and restart Steam

#### Epic Games Launcher
1. Open Settings → Network
2. Enable "Use proxy server"
3. Proxy server: `194.233.73.36:8080`

#### Minecraft
Edit `launcher_profiles.json`:
```json
{
  "profiles": {
    "profile_name": {
      "javaArgs": "-Dhttp.proxyHost=194.233.73.36 -Dhttp.proxyPort=8080 -Dhttps.proxyHost=194.233.73.36 -Dhttps.proxyPort=8080"
    }
  }
}
```

### 6. Terminal/Shell (All Applications)

#### Set Proxy for All Terminal Commands
Add to `~/.zshrc` or `~/.bashrc`:
```bash
# Proxy settings
export http_proxy=http://194.233.73.36:8080
export https_proxy=http://194.233.73.36:8080
export HTTP_PROXY=http://194.233.73.36:8080
export HTTPS_PROXY=http://194.233.73.36:8080
export no_proxy=localhost,127.0.0.1,::1
export NO_PROXY=localhost,127.0.0.1,::1

# Reload shell
source ~/.zshrc  # or source ~/.bashrc
```

#### Git Proxy
```bash
# Set proxy for Git
git config --global http.proxy http://194.233.73.36:8080
git config --global https.proxy http://194.233.73.36:8080

# Unset
git config --global --unset http.proxy
git config --global --unset https.proxy
```

#### npm/yarn Proxy
```bash
npm config set proxy http://194.233.73.36:8080
npm config set https-proxy http://194.233.73.36:8080

yarn config set proxy http://194.233.73.36:8080
yarn config set https-proxy http://194.233.73.36:8080
```

### 7. Python Applications

#### Using requests library
```python
import requests

proxies = {
    'http': 'http://194.233.73.36:8080',
    'https': 'http://194.233.73.36:8080'
}

response = requests.get('http://httpbin.org/ip', proxies=proxies)
print(response.json())
```

#### Environment Variables
```python
import os
os.environ['http_proxy'] = 'http://194.233.73.36:8080'
os.environ['https_proxy'] = 'http://194.233.73.36:8080'
```

### 8. Node.js Applications

```javascript
const axios = require('axios');

const proxyConfig = {
    host: '194.233.73.36',
    port: 8080,
    protocol: 'http'
};

axios.get('http://httpbin.org/ip', {
    proxy: proxyConfig
}).then(response => {
    console.log(response.data);
});
```

### 9. Go Applications

```go
package main

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "fmt"
)

func main() {
    proxyURL, _ := url.Parse("http://194.233.73.36:8080")
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }
    
    resp, _ := client.Get("http://httpbin.org/ip")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### 10. Mobile Devices

#### iOS (iPhone/iPad)
1. Settings → Wi-Fi
2. Tap (i) next to your network
3. Scroll to "HTTP Proxy"
4. Select "Manual"
5. Server: `194.233.73.36`
6. Port: `8080`
7. Save

#### Android
1. Settings → Wi-Fi
2. Long press your network → Modify
3. Advanced Options
4. Proxy: Manual
5. Hostname: `194.233.73.36`
6. Port: `8080`
7. Save

### 11. SOCKS5 Proxy (if enabled)

If SOCKS5 is enabled on port 1080:

```bash
# curl with SOCKS5
curl --socks5 194.233.73.36:1080 http://httpbin.org/ip

# Environment variable
export ALL_PROXY=socks5://194.233.73.36:1080
```

### 12. Test Connection

#### Quick Test Script
```bash
#!/bin/bash
PROXY="194.233.73.36:8080"

echo "Testing proxy connection..."
echo "=========================="

# Test HTTP
echo -n "HTTP: "
curl -s --connect-timeout 5 -x http://$PROXY http://httpbin.org/ip && echo " ✅" || echo " ❌"

# Test HTTPS
echo -n "HTTPS: "
curl -s --connect-timeout 5 -x http://$PROXY https://httpbin.org/ip && echo " ✅" || echo " ❌"

# Test latency
echo -n "Latency: "
time curl -s -x http://$PROXY http://httpbin.org/ip > /dev/null
```

## Troubleshooting

### Connection Refused
```bash
# Check if proxy is running
curl -x http://194.233.73.36:8080 http://httpbin.org/ip

# Check firewall
ssh root@194.233.73.36 'ufw status'

# Check service
ssh root@194.233.73.36 'systemctl status game-proxy'
```

### Timeout
```bash
# Test connectivity
ping 194.233.73.36

# Test port
nc -zv 194.233.73.36 8080

# Check server logs
ssh root@194.233.73.36 'tail -f /opt/game-proxy/logs/proxy.log'
```

### IPv6 Not Working
```bash
# Test IPv6 connectivity
ping6 2407:3640:2299:7804::1

# Use IPv4 instead
curl -x http://194.233.73.36:8080 http://httpbin.org/ip
```

## Security Notes

- ✅ Proxy uses stealth mode (no headers leaked)
- ✅ Traffic is encrypted (HTTPS)
- ⚠️ HTTP traffic is not encrypted (use HTTPS when possible)
- ⚠️ Change server password after setup

## Quick Reference

| Platform | Configuration |
|----------|--------------|
| **Mac App** | Settings → Parent Proxy: `http://194.233.73.36:8080` |
| **System Proxy** | `194.233.73.36:8080` |
| **curl** | `curl -x http://194.233.73.36:8080 URL` |
| **Environment** | `export http_proxy=http://194.233.73.36:8080` |
| **Browser** | Manual proxy: `194.233.73.36:8080` |
| **Games** | Proxy settings: `194.233.73.36:8080` |
