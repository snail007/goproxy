#/bin/bash
VER="v4.7"
rm -rf android
rm -rf ios
mkdir android
mkdir ios

#all-in-one
gomobile bind -v -target=android -javapkg=snail007 -ldflags="-s -w"
mv proxy.aar android/snail007.goproxy.sdk.aar
mv proxy-sources.jar android/snail007.goproxy.sdk-sources.jar
cp README.md android
tar zcfv sdk-android-${VER}.tar.gz android
rm -rf android
echo "done."
