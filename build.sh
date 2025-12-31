#!/bin/bash

# Build script for Mac Game Proxy App

set -e

echo "Building Mac Game Proxy App..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build Go proxy binary
echo -e "${BLUE}Building Go proxy binary...${NC}"
cd "$(dirname "$0")/.."
go build -o mac-app/GameProxyApp/proxy -ldflags="-s -w" main.go config.go

# Build Mac app (requires Xcode)
echo -e "${BLUE}Building Mac app...${NC}"
if command -v xcodebuild &> /dev/null; then
    cd mac-app
    xcodebuild -project GameProxyApp.xcodeproj -scheme GameProxyApp -configuration Release
    echo -e "${GREEN}Build complete!${NC}"
else
    echo -e "${BLUE}Xcode not found. Go binary built. Install Xcode to build Mac app.${NC}"
fi

echo -e "${GREEN}Done!${NC}"
