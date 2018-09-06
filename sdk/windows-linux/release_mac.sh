#/bin/bash
VERSION=$(cat ../../VERSION)
VER="${VERSION}_$(date '+%Y%m%d%H%M%S')"
X="-X github.com/snail007/goproxy/sdk/android-ios.SDK_VERSION=$VER -X main.APP_VERSION=$VER"
TRIMPATH1="/Users/snail/go/src/github.com/snail007"
TRIMPATH=$(dirname ~/go/src/github.com/snail007)/snail007
if [ -d "$TRIMPATH1" ];then
	TRIMPATH=$TRIMPATH1
fi
OPTS="-gcflags=-trimpath=$TRIMPATH -asmflags=-trimpath=$TRIMPATH"

rm -rf *.tar.gz
rm -rf README.md libproxy-sdk.dylib libproxy-sdk.h

#mac  , macos required
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -buildmode=c-shared $OPTS -ldflags "-s -w $X" -o libproxy-sdk.dylib sdk.go
cp ../README.md .
tar zcf sdk-mac-${VERSION}.tar.gz README.md libproxy-sdk.dylib libproxy-sdk.h
rm -rf README.md libproxy-sdk.dylib libproxy-sdk.h

echo "done."
