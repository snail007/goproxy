#!/bin/bash
F="proxy-linux-amd64.tar.gz"
manual="https://snail007.host900.com/goproxy/manual/"
set -e
WORKDIR="/tmp/proxy"
rm -rf $WORKDIR
mkdir $WORKDIR
cp $F $WORKDIR
cd /tmp/proxy
echo -e ">>> installing ... \n"
tar zxvf $F >/dev/null 2>&1
set +e
killall -9 proxy >/dev/null 2>&1
set -e
cp -f proxy /usr/bin/
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
version=`proxy --version 2>&1`
echo  -e ">>> install done, thanks for using snail007/goproxy $version\n"
echo  -e ">>> install path /usr/bin/proxy\n"
echo  -e ">>> configuration path /etc/proxy\n"
echo  -e ">>> uninstall just exec : rm /usr/bin/proxy && rm -rf /etc/proxy\n"
echo  -e ">>> How to using? Please visit : $manual\n"
