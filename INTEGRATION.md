# Integration Guide

## Overview
This document explains how to integrate the stealth game proxy features into the existing proxy codebase.

## Integration Steps

### 1. Update Go Modules

Add the new packages to your `go.mod`:

```bash
go mod tidy
```

### 2. Update Main Entry Point

Modify `main.go` to support game mode:

```go
// Add game mode flag
gameMode := app.Flag("game", "enable game mode optimizations").Default("false").Bool()
stealthMode := app.Flag("stealth", "enable stealth mode").Default("true").Bool()
```

### 3. Register Game HTTP Service

In `config.go`, add registration for game HTTP service:

```go
services.Regist("game-http", services.NewGameHTTP(), httpArgs)
```

### 4. Build Mac App

1. Build Go binary:
```bash
go build -o mac-app/GameProxyApp/proxy main.go config.go
```

2. Create Xcode project (if not exists):
```bash
cd mac-app
# Create Xcode project structure
```

3. Build Mac app:
```bash
./build.sh
```

## Usage Examples

### Command Line

```bash
# Start with game mode and stealth
./proxy game-http -t tcp -p :8080 --game --stealth

# Start with parent proxy
./proxy game-http -t tcp -p :8080 -T tls -P parent.proxy.com:443 --game --stealth
```

### Mac App

1. Launch `GameProxyApp.app`
2. Click menu bar icon
3. Select "Start Proxy"
4. Configure settings as needed

## Testing

### Test Stealth Features

```bash
# Test header randomization
curl -x http://127.0.0.1:8080 http://httpbin.org/headers

# Verify no proxy headers in response
```

### Test Game Mode

```bash
# Test UDP support
nc -u -p 5353 127.0.0.1 8080

# Test low latency
ping -c 10 127.0.0.1
```

## Configuration Files

### Stealth Configuration

Create `stealth.conf`:
```ini
[stealth]
enabled=true
header_randomization=true
tls_fingerprint_masking=true
traffic_obfuscation=true
```

### Game Configuration

Create `game.conf`:
```ini
[game]
enabled=true
low_latency=true
udp_support=true
connection_pooling=true
```

## Troubleshooting

### Build Issues

1. **Missing packages**: Run `go mod tidy`
2. **CGO errors**: Ensure Xcode command line tools installed
3. **Swift errors**: Check Xcode version compatibility

### Runtime Issues

1. **Port already in use**: Change port in settings
2. **Permission denied**: Grant network permissions
3. **Proxy not working**: Check system proxy settings

## Next Steps

1. Complete system proxy integration (todo #5)
2. Add advanced stealth features (todo #6)
3. Performance optimization
4. User testing and feedback
