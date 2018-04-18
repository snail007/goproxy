#/bin/bash
VER="v4.7"
rm -rf sdk-ios-*.tar.gz
rm -rf ios
mkdir ios

#ios  XCode required
gomobile bind -v -target=ios -ldflags="-s -w"
mv Proxy.framework ios
tar zcfv sdk-ios-${VER}.tar.gz ios
rm -rf ios

echo "done."
