#!/bin/bash
VERSION=$(cat VERSION)
VER="${VERSION}_$(date '+%Y%m%d%H%M%S')"
X="-X main.APP_VERSION=$VER"
RELEASE="release-${VERSION}"
TRIMPATH1="/Users/snail/go/src/github.com/snail007"
TRIMPATH=$(dirname ~/go/src/github.com/snail007)/snail007
if [ -d "$TRIMPATH1" ];then
	TRIMPATH=$TRIMPATH1
fi
OPTS="-gcflags=-trimpath=$TRIMPATH -asmflags=-trimpath=$TRIMPATH"

rm -rf .cert
mkdir .cert
go build $OPTS -ldflags "$X" -o proxy 
cd .cert
../proxy keygen -C proxy
cd ..
rm -rf ${RELEASE}
mkdir ${RELEASE}
#linux
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-386.tar.gz" proxy direct blocked
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm-v6.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=6 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm64-v6.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm-v7.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=7 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm64-v7.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm-v5.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=5 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm64-v5.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm64-v8.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-arm-v8.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips64le.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mipsle.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips-softfloat.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 GOMIPS=softfloat go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips64-softfloat.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le GOMIPS=softfloat go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mips64le-softfloat.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-mipsle-softfloat.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-ppc64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-ppc64le.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-linux-s390x.tar.gz" proxy direct blocked 
#android
CGO_ENABLED=0 GOOS=android GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-android-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=android GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-android-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=android GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-android-arm.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-android-arm64.tar.gz" proxy direct blocked 
#darwin
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-darwin-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-darwin-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=darwin GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-darwin-arm.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-darwin-arm64.tar.gz" proxy direct blocked 
#dragonfly
CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-dragonfly-amd64.tar.gz" proxy direct blocked 
#freebsd
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-freebsd-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-freebsd-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-freebsd-arm.tar.gz" proxy direct blocked 
#nacl
CGO_ENABLED=0 GOOS=nacl GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-nacl-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=nacl GOARCH=amd64p32 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-nacl-amd64p32.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=nacl GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-nacl-arm.tar.gz" proxy direct blocked 
#netbsd
CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-netbsd-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-netbsd-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=netbsd GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-netbsd-arm.tar.gz" proxy direct blocked 
#openbsd
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-openbsd-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-openbsd-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-openbsd-arm.tar.gz" proxy direct blocked 
#plan9
CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-plan9-386.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-plan9-amd64.tar.gz" proxy direct blocked 
CGO_ENABLED=0 GOOS=plan9 GOARCH=arm go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-plan9-arm.tar.gz" proxy direct blocked 
#solaris
CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build $OPTS -ldflags "$X" -o proxy && tar zcfv "${RELEASE}/proxy-solaris-amd64.tar.gz" proxy direct blocked 
#windows
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build $OPTS -ldflags="-H=windowsgui $X" -o proxy-noconsole.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build $OPTS  -ldflags "$X" -o proxy.exe && tar zcfv "${RELEASE}/proxy-windows-386.tar.gz" proxy.exe proxy-noconsole.exe direct blocked .cert/proxy.crt .cert/proxy.key
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $OPTS -ldflags="-H=windowsgui $X" -o proxy-noconsole.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $OPTS  -ldflags "$X" -o proxy.exe && tar zcfv "${RELEASE}/proxy-windows-amd64.tar.gz" proxy.exe proxy-noconsole.exe direct blocked .cert/proxy.crt .cert/proxy.key

rm -rf proxy proxy.exe proxy-noconsole.exe .cert

