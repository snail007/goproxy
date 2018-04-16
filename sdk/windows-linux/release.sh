#/bin/bash
VER="v4.7"

rm -rf *.tar.gz
rm -rf proxy-sdk.so proxy-sdk.h proxy-sdk.a proxy-sdk.dll

#windows 
#apt-get install gcc-multilib 
#apt-get install gcc-mingw-w64
#32bit CC=i686-w64-mingw32-gcc-win32 GOARCH=386
#64bit CC=x86_64-w64-mingw32-gcc GOARCH=amd64
CC=i686-w64-mingw32-gcc-win32 GOARCH=386 CGO_ENABLED=1 GOOS=windows go build -buildmode=c-shared -ldflags "-s -w" -o proxy-sdk.dll sdk.go
cp ../README.md .
tar zcf sdk-windows-${VER}.tar.gz README.md proxy-sdk.dll proxy-sdk.h ieshims.dll
rm -rf proxy-sdk.h proxy-sdk.dll


#linux
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-archive -ldflags "-s -w" -o proxy-sdk.a sdk.go
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -ldflags "-s -w" -o proxy-sdk.so sdk.go
cp ../README.md .
tar zcf sdk-linux-${VER}.tar.gz README.md proxy-sdk.so proxy-sdk.a proxy-sdk.h

rm -rf README.md proxy-sdk.so proxy-sdk.h proxy-sdk.a

echo "done."
