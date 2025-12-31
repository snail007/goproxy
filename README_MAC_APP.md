# Mac Game Proxy App

A stealth macOS application that provides game-optimized proxy services with complete privacy features.

## Features

### 🎮 Game Optimizations
- **Low Latency**: Optimized for real-time gaming (<50ms overhead)
- **UDP Support**: Full UDP proxy for game protocols
- **Connection Pooling**: Maintains persistent connections
- **Packet Prioritization**: Prioritizes game traffic

### 🥷 Stealth Features
- **Traffic Obfuscation**: Makes proxy traffic indistinguishable from normal traffic
- **Header Randomization**: Randomizes HTTP headers to avoid fingerprinting
- **No Proxy Signatures**: Removes all proxy-related headers
- **TLS Fingerprint Masking**: Matches TLS fingerprints to real browsers

### 🖥️ Mac App Features
- **Menu Bar App**: Runs from menu bar, no dock icon
- **Auto-Start**: Can launch on system startup
- **System Integration**: Automatically configures system proxy
- **Status Monitoring**: Shows connection status and latency
- **Minimal UI**: Clean, hidden-by-default interface

## Installation

### Prerequisites
- macOS 11.0+ (Big Sur or later)
- Xcode 13+ (for building from source)
- Go 1.16+ (for building proxy binary)

### Build from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd mac-game-proxy
```

2. Build the application:
```bash
chmod +x build.sh
./build.sh
```

3. The built app will be in `mac-app/build/Release/GameProxyApp.app`

### Install

1. Copy `GameProxyApp.app` to `/Applications`
2. Open System Preferences → Security & Privacy
3. Allow the app to run (if prompted)
4. Launch the app from Applications

## Usage

### First Launch

1. The app will appear in your menu bar (network icon)
2. Click the icon to open the menu
3. Click "Start Proxy" to begin
4. Optionally configure settings

### Normal Use

- **Menu Bar Icon Colors**:
  - 🟢 Green: Proxy running and connected
  - 🔴 Red: Proxy running but disconnected
  - ⚪ Gray: Proxy stopped

- **Menu Options**:
  - Status: Shows current proxy state
  - Latency: Shows current latency in milliseconds
  - Start/Stop Proxy: Toggle proxy on/off
  - Settings: Configure stealth mode and game mode
  - Quit: Stop proxy and exit app

### Settings

Access settings from the menu bar:
- **Stealth Mode**: Enable/disable stealth features
- **Game Mode**: Enable/disable game optimizations
- **Auto-Start**: Launch proxy on system startup

## Configuration

### System Proxy

The app automatically configures your system proxy when started. It uses:
- HTTP Proxy: `127.0.0.1:8080`
- HTTPS Proxy: `127.0.0.1:8080`
- PAC File: Auto-generated

### Advanced Configuration

Edit `~/Library/Preferences/com.gameproxy.app.plist` for advanced settings:
- Custom proxy port
- Parent proxy server
- Stealth mode settings
- Game mode settings

## Architecture

### Components

1. **Go Proxy Engine**: Core proxy functionality
   - HTTP/HTTPS proxy
   - SOCKS5 proxy
   - UDP proxy
   - Stealth features
   - Game optimizations

2. **Mac App**: SwiftUI/AppKit interface
   - Menu bar controller
   - Proxy manager
   - System integration
   - Status monitoring

3. **Background Service**: Runs proxy as daemon
   - Process management
   - Auto-restart on crash
   - Logging

## Stealth Features Explained

### Header Randomization
- Randomizes User-Agent from common browser pool
- Adds realistic browser headers
- Removes proxy-specific headers
- Randomizes header order

### TLS Fingerprint Masking
- Matches TLS cipher suites to real browsers
- Randomizes TLS extensions order
- Uses common TLS versions
- Matches JA3 fingerprint

### Traffic Obfuscation
- Adds padding to packets
- Randomizes packet timing
- Mimics normal traffic patterns
- Encrypts metadata

## Game Optimizations Explained

### Low Latency
- Connection pre-warming
- Minimal buffering
- Direct UDP forwarding
- Reduced overhead

### UDP Support
- Full UDP proxy
- Low-latency forwarding
- Connection tracking
- Packet prioritization

## Troubleshooting

### Proxy Won't Start
1. Check if port 8080 is available
2. Verify app has necessary permissions
3. Check Console.app for error messages

### System Proxy Not Working
1. Manually configure proxy in System Preferences
2. Check firewall settings
3. Verify proxy is running (check menu bar icon)

### High Latency
1. Disable other VPN/proxy services
2. Check network connection
3. Try different proxy port
4. Enable game mode optimizations

## Security & Privacy

- **No Logging**: Traffic is not logged
- **No Data Collection**: No analytics or telemetry
- **Local Only**: All configuration stored locally
- **Encrypted Storage**: Configuration encrypted at rest

## License

[Your License Here]

## Support

For issues and questions:
- GitHub Issues: [repository-url]/issues
- Documentation: See DESIGN.md for architecture details

## Development

See DESIGN.md for detailed architecture and implementation plan.
