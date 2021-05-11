#!/bin/bash
set -e
ver=$1
if [ -z "$ver" ]; then
  echo -e "example:\n./build.sh 10.0"
exit
fi
CLEAN="goproxy proxy"
rm -rf $CLEAN
mkdir goproxy

cd goproxy
wget https://mirrors.host900.com/snail007/goproxy/proxy-linux-amd64.tar.gz
tar zxf proxy-linux-amd64.tar.gz
mv proxy ../
cd ..

docker build --no-cache -t snail007/goproxy:v$ver .
docker tag snail007/goproxy:v$ver snail007/goproxy:latest
docker images
rm -rf $CLEAN