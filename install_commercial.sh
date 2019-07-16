#!/bin/bash
F="proxy-linux-amd64_commercial.tar.gz"

set -e

# #install proxy
tar zxvf $F
cp -f proxy /usr/bin/
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
