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

rm -rf sdk-linux-*.tar.gz
rm -rf README.md libproxy-sdk.so libproxy-sdk.h libproxy-sdk.a

#linux 32bit
CGO_ENABLED=1 GOARCH=386 GOOS=linux go build -buildmode=c-archive $OPTS -ldflags "-s -w $X" -o libproxy-sdk.a sdk.go
CGO_ENABLED=1 GOARCH=386 GOOS=linux go build -buildmode=c-shared $OPTS -ldflags "-s -w $X" -o libproxy-sdk.so sdk.go
cp ../README.md .
tar zcf sdk-linux-32bit-${VERSION}.tar.gz README.md libproxy-sdk.so libproxy-sdk.a libproxy-sdk.h
rm -rf README.md libproxy-sdk.so libproxy-sdk.h libproxy-sdk.a

#linux 64bit
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-archive $OPTS -ldflags "-s -w $X" -o libproxy-sdk.a sdk.go
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-shared $OPTS -ldflags "-s -w $X" -o libproxy-sdk.so sdk.go
cp ../README.md .
tar zcf sdk-linux-64bit-${VERSION}.tar.gz README.md libproxy-sdk.so libproxy-sdk.a libproxy-sdk.h
rm -rf README.md libproxy-sdk.so libproxy-sdk.h libproxy-sdk.a

echo "done."



