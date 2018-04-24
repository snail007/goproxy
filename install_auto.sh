#!/bin/bash
set -e
if [ -e /tmp/proxy ]; then
    rm -rf /tmp/proxy
fi
mkdir /tmp/proxy
cd /tmp/proxy
wget https://github.com/onetwotrip/goproxy/releases/download/v4.6/proxy-linux-amd64.tar.gz

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
    proxy keygen >/dev/null 2>&1 
fi
rm -rf /tmp/proxy
echo "install done"
proxy help
