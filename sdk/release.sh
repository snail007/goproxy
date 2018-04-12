#/bin/bash
VER="v4.7"
rm -rf android
rm -rf ios
mkdir android
mkdir ios

#arm
gomobile bind -v -target=android/arm -javapkg=snail007 -ldflags="-s -w"
mkdir arm
mv proxy.aar arm/snail007.goproxy.sdk.aar
mv proxy-sources.jar arm/snail007.goproxy.sdk-sources.jar
tar zcfv sdk-arm-${VER}.tar.gz arm
mv sdk-arm-${VER}.tar.gz android
rm -rf arm


#arm64
gomobile bind -v -target=android/arm64 -javapkg=snail007 -ldflags="-s -w"
mkdir arm64
mv proxy.aar arm64/snail007.goproxy.sdk.aar
mv proxy-sources.jar arm64/snail007.goproxy.sdk-sources.jar
tar zcfv sdk-arm64-${VER}.tar.gz arm64
mv sdk-arm64-${VER}.tar.gz android
rm -rf arm64


#386
gomobile bind -v -target=android/386 -javapkg=snail007 -ldflags="-s -w"
mkdir 386
mv proxy.aar 386/snail007.goproxy.sdk.aar
mv proxy-sources.jar 386/snail007.goproxy.sdk-sources.jar
tar zcfv sdk-386-${VER}.tar.gz 386
mv sdk-386-${VER}.tar.gz android
rm -rf 386

#amd64
gomobile bind -v -target=android/amd64 -javapkg=snail007 -ldflags="-s -w"
mkdir amd64
mv proxy.aar amd64/snail007.goproxy.sdk.aar
mv proxy-sources.jar amd64/snail007.goproxy.sdk-sources.jar
tar zcfv sdk-amd64-${VER}.tar.gz amd64
mv sdk-amd64-${VER}.tar.gz android
rm -rf amd64


#all-in-one
gomobile bind -v -target=android -javapkg=snail007 -ldflags="-s -w"
mkdir all
mv proxy.aar all/snail007.goproxy.sdk.aar
mv proxy-sources.jar all/snail007.goproxy.sdk-sources.jar
tar zcfv sdk-all-${VER}.tar.gz all
mv sdk-all-${VER}.tar.gz android
rm -rf all

echo "done."
