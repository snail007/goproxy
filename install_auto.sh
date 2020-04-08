#!/bin/bash
F="proxy-linux-amd64.tar.gz"
set -e
if [ -e /tmp/proxy ]; then
    rm -rf /tmp/proxy
fi
mkdir /tmp/proxy
cd /tmp/proxy
echo -e "\n>>> downloading ... $F\n"
set +e
CN=$(wget -O - -t 3 --dns-timeout 1 --connect-timeout 2 --read-timeout 2 myip.ipip.net)
if [ "$CN" != "" ];then
CN=$(echo $CN| grep "中国" |grep -v grep)
fi
set -e
manual=""
if [ -z "$CN" ];then
manual="https://snail007.github.io/goproxy/manual/"
LAST_VERSION=$(curl --silent "https://api.github.com/repos/snail007/goproxy/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")')
wget  -t 1 "https://github.com/snail007/goproxy/releases/download/${LAST_VERSION}/$F"
else
manual="https://snail007.github.io/goproxy/manual/zh/"
wget  -t 1 "http://mirrors.host900.com:9090/snail007/goproxy/$F"
fi
echo -e ">>> installing ... \n"
# #install proxy
tar zxvf $F >/dev/null
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
echo  -e ">>> uninstall just exec : rm /usr/bin/proxy && rm /etc/proxy\n"
echo  -e ">>> How to using? Please visit : $manual\n"
