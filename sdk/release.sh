#/bin/bash
VER="v4.7"
rm -rf proxy-sdk-release-*
#arm
gomobile bind -v -target=android/arm 
mkdir proxy-sdk-arm
mv sdk.aar proxy-sdk-arm/proxy-sdk-arm.aar
mv sdk-sources.jar proxy-sdk-arm/proxy-sdk-arm-sources.jar
tar zcfv proxy-sdk-arm-${VER}.tar.gz proxy-sdk-arm
rm -rf proxy-sdk-arm
#arm64
gomobile bind -v -target=android/arm64
mkdir proxy-sdk-arm64
mv sdk.aar proxy-sdk-arm64/proxy-sdk-arm64.aar
mv sdk-sources.jar proxy-sdk-arm64/proxy-sdk-arm64-sources.jar
tar zcfv proxy-sdk-arm64-${VER}.tar.gz proxy-sdk-arm64
rm -rf proxy-sdk-arm64
#386
gomobile bind -v -target=android/386
mkdir proxy-sdk-386
mv sdk.aar proxy-sdk-386/proxy-sdk-386.aar
mv sdk-sources.jar proxy-sdk-386/proxy-sdk-386-sources.jar
tar zcfv proxy-sdk-386-${VER}.tar.gz proxy-sdk-386
rm -rf proxy-sdk-386
#amd64
gomobile bind -v -target=android/amd64
mkdir proxy-sdk-amd64
mv sdk.aar proxy-sdk-amd64/proxy-sdk-amd64.aar
mv sdk-sources.jar proxy-sdk-amd64/proxy-sdk-amd64-sources.jar
tar zcfv proxy-sdk-amd64-${VER}.tar.gz proxy-sdk-amd64
rm -rf proxy-sdk-amd64
#all-in-one
gomobile bind -v -target=android
mkdir proxy-sdk-all
mv sdk.aar proxy-sdk-all/proxy-sdk-all.aar
mv sdk-sources.jar proxy-sdk-all/proxy-sdk-all-sources.jar
tar zcfv proxy-sdk-all-${VER}.tar.gz proxy-sdk-all
rm -rf proxy-sdk-all
mkdir proxy-sdk-release-${VER}
mv *.tar.gz proxy-sdk-release-${VER}
echo "done."
