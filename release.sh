#!/bin/bash
VER="3.4"
RELEASE="release-${VER}"
rm -rf .cert
mkdir .cert
go build 
cd .cert
../proxy keygen
cd ..
rm -rf ${RELEASE}
mkdir ${RELEASE}
set CGO_ENABLED=0
#linux
GOOS=linux GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-linux-386.tar.gz" proxy direct blocked
GOOS=linux GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-linux-amd64.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=arm GOARM=7 go build && tar zcfv "${RELEASE}/proxy-linux-arm.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=arm64 GOARM=7 go build && tar zcfv "${RELEASE}/proxy-linux-arm64.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=mips go build && tar zcfv "${RELEASE}/proxy-linux-mips.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=mips64 go build && tar zcfv "${RELEASE}/proxy-linux-mips64.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=mips64le go build && tar zcfv "${RELEASE}/proxy-linux-mips64le.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=mipsle go build && tar zcfv "${RELEASE}/proxy-linux-mipsle.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=ppc64 go build && tar zcfv "${RELEASE}/proxy-linux-ppc64.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=ppc64le go build && tar zcfv "${RELEASE}/proxy-linux-ppc64le.tar.gz" proxy direct blocked 
GOOS=linux GOARCH=s390x go build && tar zcfv "${RELEASE}/proxy-linux-s390x.tar.gz" proxy direct blocked 
#android
GOOS=android GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-android-386.tar.gz" proxy direct blocked 
GOOS=android GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-android-amd64.tar.gz" proxy direct blocked 
GOOS=android GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-android-arm.tar.gz" proxy direct blocked 
GOOS=android GOARCH=arm64 go build && tar zcfv "${RELEASE}/proxy-android-arm64.tar.gz" proxy direct blocked 
#darwin
GOOS=darwin GOARCH=386 go build go build && tar zcfv "${RELEASE}/proxy-darwin-386.tar.gz" proxy direct blocked 
GOOS=darwin GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-darwin-amd64.tar.gz" proxy direct blocked 
GOOS=darwin GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-darwin-arm.tar.gz" proxy direct blocked 
GOOS=darwin GOARCH=arm64 go build && tar zcfv "${RELEASE}/proxy-darwin-arm64.tar.gz" proxy direct blocked 
#dragonfly
GOOS=dragonfly GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-dragonfly-amd64.tar.gz" proxy direct blocked 
#freebsd
GOOS=freebsd GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-freebsd-386.tar.gz" proxy direct blocked 
GOOS=freebsd GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-freebsd-amd64.tar.gz" proxy direct blocked 
GOOS=freebsd GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-freebsd-arm.tar.gz" proxy direct blocked 
#nacl
GOOS=nacl GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-nacl-386.tar.gz" proxy direct blocked 
GOOS=nacl GOARCH=amd64p32 go build && tar zcfv "${RELEASE}/proxy-nacl-amd64p32.tar.gz" proxy direct blocked 
GOOS=nacl GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-nacl-arm.tar.gz" proxy direct blocked 
#netbsd
GOOS=netbsd GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-netbsd-386.tar.gz" proxy direct blocked 
GOOS=netbsd GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-netbsd-amd64.tar.gz" proxy direct blocked 
GOOS=netbsd GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-netbsd-arm.tar.gz" proxy direct blocked 
#openbsd
GOOS=openbsd GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-openbsd-386.tar.gz" proxy direct blocked 
GOOS=openbsd GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-openbsd-amd64.tar.gz" proxy direct blocked 
GOOS=openbsd GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-openbsd-arm.tar.gz" proxy direct blocked 
#plan9
GOOS=plan9 GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-plan9-386.tar.gz" proxy direct blocked 
GOOS=plan9 GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-plan9-amd64.tar.gz" proxy direct blocked 
GOOS=plan9 GOARCH=arm go build && tar zcfv "${RELEASE}/proxy-plan9-arm.tar.gz" proxy direct blocked 
#solaris
GOOS=solaris GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-solaris-amd64.tar.gz" proxy direct blocked 
#windows
GOOS=windows GOARCH=386 go build && tar zcfv "${RELEASE}/proxy-windows-386.tar.gz" proxy.exe  direct blocked .cert/proxy.crt .cert/proxy.key
GOOS=windows GOARCH=amd64 go build && tar zcfv "${RELEASE}/proxy-windows-amd64.tar.gz" proxy.exe  direct blocked .cert/proxy.crt .cert/proxy.key

rm -rf proxy proxy.exe .cert
