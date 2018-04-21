#/bin/bash
VER="v4.7"

rm -rf sdk-linux-*.tar.gz
rm -rf README.md libproxy-sdk.so libproxy-sdk.h libproxy-sdk.a

#linux
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-archive -ldflags "-s -w" -o libproxy-sdk.a sdk.go
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -ldflags "-s -w" -o libproxy-sdk.so sdk.go
cp ../README.md .
tar zcf sdk-linux-${VER}.tar.gz README.md libproxy-sdk.so libproxy-sdk.a libproxy-sdk.h
rm -rf README.md libproxy-sdk.so libproxy-sdk.h libproxy-sdk.a

echo "done."
