#/bin/bash
VER="v5.0"

rm -rf *.tar.gz
rm -rf README.md libproxy-sdk.dylib libproxy-sdk.h

#mac  , macos required
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -buildmode=c-shared -ldflags "-s -w" -o libproxy-sdk.dylib sdk.go
cp ../README.md .
tar zcf sdk-mac-${VER}.tar.gz README.md libproxy-sdk.dylib libproxy-sdk.h
rm -rf README.md libproxy-sdk.dylib libproxy-sdk.h

echo "done."
