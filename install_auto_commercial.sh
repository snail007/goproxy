#!/bin/bash
F="proxy-linux-amd64_commercial.tar.gz"
V="v8.9"
set -e
if [ -e /tmp/proxy ]; then
    rm -rf /tmp/proxy
fi
mkdir /tmp/proxy
cd /tmp/proxy
echo -e "downloading ... $F-$V\n"
CN=$(wget -O - myip.ipip.net | grep "中国" |grep -v grep)
if [ -z "$CN" ];then
LAST_VERSION=$(curl --silent "https://api.github.com/repos/snail007/goproxy/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")')
wget "https://github.com/snail007/goproxy/releases/download/${LAST_VERSION}/$F"
else
wget "http://mirrors.host900.com/snail007/goproxy/$V/$F"
fi
echo -e "installing ... \n"
# #install proxy
tar zxvf $F
killall -9 proxy
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
echo -e "\n#######################\n"
proxy --version
echo  -e ">>> install done, thanks for using snail007/goproxy\n"
echo  -e ">>> install path /usr/bin/proxy\n"
echo  -e ">>> configuration path /etc/proxy\n\n"
echo  -e ">>> uninstall just exec : rm /usr/bin/proxy && rm /etc/proxy\n\n"
echo  -e ">>> How to using? Please visit : https://snail007.github.io/goproxy/manual/zh/\n"
echo  -e "#######################\n"
