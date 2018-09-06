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

#sudo rm /usr/local/go
#sudo ln -s /usr/local/go1.10.1 /usr/local/go
rm -rf sdk-windows-*.tar.gz
rm -rf README.md proxy-sdk.h proxy-sdk.dll


#apt-get install gcc-multilib 
#apt-get install gcc-mingw-w64

#windows 64bit
CC=x86_64-w64-mingw32-gcc GOARCH=amd64 CGO_ENABLED=1 GOOS=windows go build $OPTS -buildmode=c-shared -ldflags "-s -w $X" -o proxy-sdk.dll sdk.go
cp ../README.md .
tar zcf sdk-windows-64bit-${VERSION}.tar.gz README.md proxy-sdk.dll proxy-sdk.h ieshims.dll
rm -rf README.md proxy-sdk.h proxy-sdk.dll

#windows 32bit
CC=i686-w64-mingw32-gcc-win32 GOARCH=386 CGO_ENABLED=1 GOOS=windows go build $OPTS -buildmode=c-shared -ldflags "-s -w $X" -o proxy-sdk.dll sdk.go
cp ../README.md .
tar zcf sdk-windows-32bit-${VERSION}.tar.gz README.md proxy-sdk.dll proxy-sdk.h ieshims.dll
rm -rf README.md proxy-sdk.h proxy-sdk.dll

#sudo rm /usr/local/go
#sudo ln -s /usr/local/go1.8.5 /usr/local/go

echo "done."
