#/bin/bash
VER="v4.7"

rm -rf *.tar.gz
rm -rf README.md proxy-sdk.dylib proxy-sdk.h

#mac  , macos required
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -buildmode=c-shared -ldflags "-s -w" -o proxy-sdk.dylib sdk.go
cp ../README.md .
tar zcf sdk-mac-${VER}.tar.gz README.md proxy-sdk.dylib proxy-sdk.h
rm -rf README.md proxy-sdk.dylib proxy-sdk.h

echo "done."
