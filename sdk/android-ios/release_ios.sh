#/bin/bash
VERSION=$(cat ../../VERSION)
VER="${VERSION}_$(date '+%Y%m%d%H%M%S')"
X="-X github.com/snail007/goproxy/sdk/android-ios.SDK_VERSION=$VER -X main.APP_VERSION=$VER"

rm -rf sdk-ios-*.tar.gz
rm -rf ios
mkdir ios

#ios  XCode required
gomobile bind -v -target=ios -ldflags="-s -w $X"
mv Proxy.framework ios
cp ../README.md ios
tar zcfv sdk-ios-${VERSION}.tar.gz ios
rm -rf ios

echo "done."
