#!/bin/bash
set -e
if [ -e /tmp/proxy ]; then
    rm -rf /tmp/proxy
fi
mkdir /tmp/proxy
cd /tmp/proxy

LAST_VERSION=$(curl --silent "https://api.github.com/repos/snail007/goproxy/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")')

wget "https://github.com/snail007/goproxy/releases/download/${LAST_VERSION}/proxy-linux-amd64.tar.gz"

# #install proxy
tar zxvf proxy-linux-amd64.tar.gz
cp proxy /usr/bin/
chmod +x /usr/bin/proxy
if [ ! -e /etc/proxy ]; then
    mkdir /etc/proxy
    cp blocked /etc/proxy
    cp direct  /etc/proxy
fi

if [ ! -e /etc/proxy/proxy.crt ]; then
    cd /etc/proxy/
    proxy keygen -C proxy >/dev/null 2>&1 
fi
rm -rf /tmp/proxy
echo "install done"
proxy help
