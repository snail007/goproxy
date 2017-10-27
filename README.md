<img src="https://github.com/snail007/goproxy/blob/master/docs/images/logo.jpg?raw=true" width="200"/>  
Proxy是golang实现的高性能http,https,websocket,tcp,udp,socks5代理服务器,支持正向代理、内网穿透、SSH中转。  
  
---  
  
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https://img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total.svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style=plastic)](https://github.com/snail007/goproxy/releases)  
  
### Features  
- 链式代理,程序本身可以作为一级代理,如果设置了上级代理那么可以作为二级代理,乃至N级代理.  
- 通讯加密,如果程序不是一级代理,而且上级代理也是本程序,那么可以加密和上级代理之间的通讯,采用底层tls高强度加密,安全无特征.  
- 智能HTTP,SOCKS5代理,会自动判断访问的网站是否屏蔽,如果被屏蔽那么就会使用上级代理(前提是配置了上级代理)访问网站;如果访问的网站没有被屏蔽,为了加速访问,代理会直接访问网站,不使用上级代理.  
- 域名黑白名单，更加自由的控制网站的访问方式。  
- 跨平台性,无论你是widows,linux,还是mac,甚至是树莓派,都可以很好的运行proxy.  
- 多协议支持,支持HTTP(S),TCP,UDP,Websocket,SOCKS5代理.  
- 支持内网穿透,协议支持TCP和UDP.  
- SSH中转,HTTP(S),SOCKS5代理支持SSH中转,上级Linux服务器不需要任何服务端,本地一个proxy即可开心上网.  
- 支持[KCP](https://github.com/xtaci/kcp-go)协议,HTTP(S),SOCKS5代理支持KCP协议传输数据,降低延迟,提升浏览体验.  
  
### Why need these?  
- 当由于某某原因,我们不能访问我们在其它地方的服务,我们可以通过多个相连的proxy节点建立起一个安全的隧道访问我们的服务.  
- 微信接口本地开发,方便调试.  
- 远程访问内网机器.  
- 和小伙伴一起玩局域网游戏.  
- 以前只能在局域网玩的,现在可以在任何地方玩.  
- 替代圣剑内网通，显IP内网通，花生壳之类的工具.
- ...  

 
本页是v3.4手册,其他版本手册请点击下面链接查看.  
- [v3.3手册](https://github.com/snail007/goproxy/tree/v3.3)
- [v3.2手册](https://github.com/snail007/goproxy/tree/v3.2)
- [v3.1手册](https://github.com/snail007/goproxy/tree/v3.1)
- [v3.0手册](https://github.com/snail007/goproxy/tree/v3.0)
- [v2.x手册](https://github.com/snail007/goproxy/tree/v2.2)  

### 安装 
1. [快速安装](#自动安装)
1. [手动安装](#手动安装)

### 首次使用必看
- [环境](#使用教程)
- [使用配置文件](#使用配置文件)
- [生成通讯证书文件](#生成加密通讯需要的证书文件)
- [安全建议](#安全建议)

### 手册目录
- [1. HTTP代理](#1http代理)
    - [1.1 普通HTTP代理](#11普通http代理)
    - [1.2 普通二级HTTP代理](#12普通二级http代理)
    - [1.3 HTTP二级代理(加密)](#13http二级代理加密)
    - [1.4 HTTP三级代理(加密)](#14http三级代理加密)
    - [1.5 Basic认证](#15basic认证)
    - [1.6 强制走上级HTTP代理](#16http代理流量强制走上级http代理)
    - [1.7 通过SSH中转](#17https通过ssh中转)
        - [1.7.1 用户名和密码的方式](#171-ssh用户名和密码的方式)
        - [1.7.2 用户名和密钥的方式](#172-ssh用户名和密钥的方式)
    - [1.8 KCP协议传输](#18kcp协议传输)
    - [1.9 查看帮助](#19查看帮助)
- [2. TCP代理](#2tcp代理)
    - [2.1 普通一级TCP代理](#21普通一级tcp代理)
    - [2.2 普通二级TCP代理](#22普通二级tcp代理)
    - [2.3 普通三级TCP代理](#23普通三级tcp代理)
    - [2.4 加密二级TCP代理](#24加密二级tcp代理)
    - [2.5 加密三级TCP代理](#25加密三级tcp代理)
    - [2.6 查看帮助](#26查看帮助)
- [3. UDP代理](#3udp代理)
    - [3.1 普通一级TCP代理](#31普通一级udp代理)
    - [3.2 普通二级TCP代理](#32普通二级udp代理)
    - [3.3 普通三级TCP代理](#33普通三级udp代理)
    - [3.4 加密二级TCP代理](#34加密二级udp代理)
    - [3.5 加密三级TCP代理](#35加密三级udp代理)
    - [3.6 查看帮助](#36查看帮助)
- [4. 内网穿透](#4内网穿透)
    - [4.1 原理说明](#41原理说明)
    - [4.2 TCP普通用法](#42tcp普通用法)
    - [4.3 微信接口本地开发](#43微信接口本地开发)
    - [4.4 UDP普通用法](#44udp普通用法)
    - [4.5 高级用法一](#45高级用法一)
    - [4.6 高级用法一](#46高级用法二)
    - [4.7 tserver的-r参数](#47tserver的-r参数)
    - [4.8 查看帮助](#48查看帮助)
- [5. SOCKS5代理](#5socks5代理)
    - [5.1 普通SOCKS5代理](#51普通socks5代理)
    - [5.2 普通二级SOCKS5代理](#52普通二级socks5代理)
    - [5.3 SOCKS二级代理(加密)](#53socks二级代理加密)
    - [5.4 SOCKS三级代理(加密)](#54socks三级代理加密)
    - [5.5 流量强制走上级SOCKS代理](#55socks代理流量强制走上级socks代理)
    - [5.6 通过SSH中转](#56socks通过ssh中转)
        - [5.6.1 用户名和密码的方式](#561-ssh用户名和密码的方式)
        - [5.6.2 用户名和密钥的方式](#562-ssh用户名和密钥的方式)
    - [5.7 认证](#57认证)
    - [5.8 KCP协议传输](#58kcp协议传输)
    - [5.9 查看帮助](#59查看帮助)

### Fast Start  
提示:所有操作需要root权限.  
#### 自动安装
#### **0.如果你的VPS是linux64位的系统,那么只需要执行下面一句,就可以完成自动安装和配置.**  
```shell  
curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash  
```  
安装完成,配置目录是/etc/proxy,更详细的使用方法参考下面的进一步了解.  
如果安装失败或者你的vps不是linux64位系统,请按照下面的半自动步骤安装:  
  
#### 手动安装
#### **1.登录你的VPS,下载守护进程monexec,选择合适你的版本,vps一般选择"linux_amd64.tar.gz"的即可.**  
下载地址:https://github.com/reddec/monexec/releases  
比如下载到/root/proxy/  
执行:  
```shell  
mkdir /root/proxy/  
cd /root/proxy/  
wget https://github.com/reddec/monexec/releases/download/v0.1.1/monexec_0.1.1_linux_amd64.tar.gz  
```  
#### **2.下载proxy**  
下载地址:https://github.com/snail007/goproxy/releases  
```shell  
cd /root/proxy/  
wget https://github.com/snail007/goproxy/releases/download/v3.1fix/proxy-linux-amd64.tar.gz  
```  
#### **3.下载自动安装脚本**  
```shell  
cd /root/proxy/  
wget https://raw.githubusercontent.com/snail007/goproxy/master/install.sh  
chmod +x install.sh  
./install.sh  
```  
  
## 使用教程  
  
#### **提示**  
接下来的教程,默认系统是linux,程序是proxy；所有操作需要root权限；  
如果你的是windows,请使用windows版本的proxy.exe即可.  
  
### **使用配置文件**  
接下来的教程都是通过命令行参数介绍使用方法,也可以通过读取配置文件获取参数.  
具体格式是通过@符号指定配置文件,例如:./proxy @configfile.txt  
configfile.txt里面的格式是,第一行是子命令名称,第二行开始一行一个:参数的长格式=参数值,前后不能有空格和双引号.  
参数的长格式都是--开头的,短格式参数都是-开头,如果你不知道某个短格式参数对应长格式参数,查看帮助命令即可.  
比如configfile.txt内容如下:
```shell
http
--local-type=tcp
--local=:33080
```
### 0.生成加密通讯需要的证书文件  
  
http,tcp,udp代理过程会和上级通讯,为了安全我们采用加密通讯,当然可以选择不加密通信通讯,本教程所有和上级通讯都采用加密,需要证书文件.  
在linux上并安装了openssl命令，可以直接通过下面的命令生成证书和key文件.  
`./proxy keygen`  
默认会在当前程序目录下面生成证书文件proxy.crt和key文件proxy.key。  
  
### 安全建议
当VPS在nat设备后面,vps上网卡IP都是内网IP,这个时候可以通过-g参数添加vps的外网ip防止死循环.  
假设你的vps外网ip是23.23.23.23,下面命令通过-g参数设置23.23.23.23  
`./proxy http -g "23.23.23.23"`  

### 1.HTTP代理  
#### **1.1.普通HTTP代理**  
`./proxy http -t tcp -p "0.0.0.0:38080"`  
  
#### **1.2.普通二级HTTP代理**  
使用本地端口8090,假设上级HTTP代理是`22.22.22.22:8080`  
`./proxy http -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" `  
默认关闭了连接池,如果要加快访问速度,-L可以开启连接池,10就是连接池大小,0为关闭,  
开启连接池在网络不好的情况下,稳定不是很好.   
`./proxy http -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" -L 10`  
我们还可以指定网站域名的黑白名单文件,一行一个域名,怕匹配规则是最右批评匹配,比如:baidu.com,匹配的是*.*.baidu.com,黑名单的域名域名直接走上级代理,白名单的域名不走上级代理.  
`./proxy http -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt`  
  
#### **1.3.HTTP二级代理(加密)**  
一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key`  
  
二级HTTP代理(本地Linux)  
`./proxy http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080.  
  
二级HTTP代理(本地windows)  
`./proxy.exe http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为http模式，地址为：127.0.0.1，端口为：8080,程序即可通过加密通道通过vps上网。  
  
#### **1.4.HTTP三级代理(加密)**  
一级HTTP代理VPS_01,IP:22.22.22.22  
`./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key`  
二级HTTP代理VPS_02,IP:33.33.33.33  
`./proxy http -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级HTTP代理(本地)  
`./proxy http -t tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问一级HTTP代理上面的代理端口38080.  
  
#### **1.5.Basic认证**  
对于代理HTTP协议我们可以basic进行Basic认证,认证的用户名和密码可以在命令行指定  
`./proxy http -t tcp -p ":33080" -a "user1:pass1" -a "user2:pass2"`  
多个用户,重复-a参数即可.  
也可以放在文件中,格式是一行一个"用户名:密码",然后用-F指定.  
`./proxy http -t tcp -p ":33080" -F auth-file.txt`  
如果没有-a或-F参数,就是关闭Basic认证.  
  
#### **1.6.HTTP代理流量强制走上级HTTP代理**  
默认情况下,proxy会智能判断一个网站域名是否无法访问,如果无法访问才走上级HTTP代理.通过--always可以使全部HTTP代理流量强制走上级HTTP代理.  
`./proxy http --always -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
  
#### **1.7.HTTP(S)通过SSH中转**  
说明:ssh中转的原理是利用了ssh的转发功能,就是你连接上ssh之后,可以通过ssh代理访问目标地址.  
假设有:vps  
- IP是2.2.2.2, ssh端口是22, ssh用户名是:user, ssh用户密码是:demo  
- 用户user的ssh私钥名称是user.key    

##### ***1.7.1 ssh用户名和密码的方式***   
本地HTTP(S)代理28080端口,执行:  
`./proxy http -T ssh -P "2.2.2.2:22" -u user -A demo -t tcp -p ":28080"`  
##### ***1.7.2 ssh用户名和密钥的方式***   
本地HTTP(S)代理28080端口,执行:  
`./proxy http -T ssh -P "2.2.2.2:22" -u user -S user.key -t tcp -p ":28080"`  

#### **1.8.KCP协议传输**  
KCP协议需要-B参数设置一个密码用于加密解密数据  

一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy http -t kcp -p ":38080" -B mypassword  
  
二级HTTP代理(本地Linux)  
`./proxy http -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" -B mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080,数据通过kcp协议传输.  

#### **1.9.查看帮助**  
`./proxy help http`  
  
### 2.TCP代理  
  
#### **2.1.普通一级TCP代理**  
本地执行:  
`./proxy tcp -p ":33080" -T tcp -P "192.168.22.33:22" -L 0`  
那么访问本地33080端口就是访问192.168.22.33的22端口.  
  
#### **2.2.普通二级TCP代理**  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -p ":33080" -T tcp -P "127.0.0.1:8080" -L 0`  
本地执行:  
`./proxy tcp -p ":23080" -T tcp -P "22.22.22.33:33080"`  
那么访问本地23080端口就是访问22.22.22.33的8080端口.  
  
#### **2.3.普通三级TCP代理**  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -p ":38080" -T tcp -P "66.66.66.66:8080" -L 0`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`./proxy tcp -p ":8080" -T tcp -P "33.33.33.33:28080"`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口.  
  
#### **2.4.加密二级TCP代理**  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp --tls -p ":33080" -T tcp -P "127.0.0.1:8080" -L 0 -C proxy.crt -K proxy.key`  
本地执行:  
`./proxy tcp -p ":23080" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地23080端口就是通过加密TCP隧道访问22.22.22.33的8080端口.  
  
#### **2.5.加密三级TCP代理**  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp --tls -p ":38080" -T tcp -P "66.66.66.66:8080" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp --tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`./proxy tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口.  
  
#### **2.6.查看帮助**  
`./proxy help tcp`  
  
### 3.UDP代理  
  
#### **3.1.普通一级UDP代理**  
本地执行:  
`./proxy udp -p ":5353" -T udp -P "8.8.8.8:53"`  
那么访问本地UDP:5353端口就是访问8.8.8.8的UDP:53端口.  
  
#### **3.2.普通二级UDP代理**  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -p ":33080" -T udp -P "8.8.8.8:53"`  
本地执行:  
`./proxy udp -p ":5353" -T tcp -P "22.22.22.33:33080"`  
那么访问本地UDP:5353端口就是通过TCP隧道,通过VPS访问8.8.8.8的UDP:53端口.  
  
#### **3.3.普通三级UDP代理**  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -p ":38080" -T udp -P "8.8.8.8:53"`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`./proxy udp -p ":5353" -T tcp -P "33.33.33.33:28080"`  
那么访问本地5353端口就是通过TCP隧道,通过VPS访问8.8.8.8的53端口.  
  
#### **3.4.加密二级UDP代理**  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp --tls -p ":33080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
本地执行:  
`./proxy udp -p ":5353" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地UDP:5353端口就是通过加密TCP隧道,通过VPS访问8.8.8.8的UDP:53端口.  
  
#### **3.5.加密三级UDP代理**  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp --tls -p ":38080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp --tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`./proxy udp -p ":5353" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地5353端口就是通过加密TCP隧道,通过VPS_01访问8.8.8.8的53端口.  
  
#### **3.6.查看帮助**  
`./proxy help udp`  
  
### 4.内网穿透  
#### **4.1、原理说明**  
内网穿透,由三部分组成:client端,server端,bridge端；client和server主动连接bridge端进行桥接.  
当用户访问server端,流程是:  
1. server主动和bridge端建立连接；  
1. 然后bridge端通知client端连接bridge端,并连接内网目标端口;  
1. 然后绑定client端到bridge端和client端到内网端口的连接；  
1. 然后bridge端把client过来的连接与server端过来的连接绑定；  
1. 整个通道建立完成；  
  
#### **4.2、TCP普通用法**  
背景:  
- 公司机器A提供了web服务80端口  
- 有VPS一个,公网IP:22.22.22.22  

需求:  
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口  
  
步骤:  
1. 在vps上执行  
    `./proxy tbridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy tserver -r ":28080@:80" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  
  
1. 在公司机器A上面执行  
    `./proxy tclient -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 完成  
  
#### **4.3、微信接口本地开发**  
背景:  
- 自己的笔记本提供了nginx服务80端口  
- 有VPS一个,公网IP:22.22.22.22  

需求:  
在微信的开发帐号的网页回调接口配置里面填写地址:http://22.22.22.22/calback.php  
然后就可以访问到笔记本的80端口下面的calback.php,如果需要绑定域名,可以用自己的域名  
比如:wx-dev.xxx.com解析到22.22.22.22,然后在自己笔记本的nginx里  
配置域名wx-dev.xxx.com到具体的目录即可.  

  
步骤:  
1. 在vps上执行,确保vps的80端口没被其它程序占用.  
    `./proxy tbridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy tserver -r ":80@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 在自己笔记本上面执行  
    `./proxy tclient -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 完成  
  
#### **4.4、UDP普通用法**  
背景:  
- 公司机器A提供了DNS解析服务,UDP:53端口  
- 有VPS一个,公网IP:22.22.22.22  
  
需求:  
在家里能够通过设置本地dns为22.22.22.22,使用公司机器A进行域名解析服务.  
  
步骤:  
1. 在vps上执行  
    `./proxy tbridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy tserver --udp -r ":53@:53" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  

1. 在公司机器A上面执行  
    `./proxy tclient --udp -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 完成  
  
#### **4.5、高级用法一**  
背景:  
- 公司机器A提供了web服务80端口  
- 有VPS一个,公网IP:22.22.22.22  
  
需求:  
为了安全,不想在VPS上能够访问到公司机器A,在家里能够通过访问本机的28080端口,  
通过加密隧道访问到公司机器A的80端口.  
  
步骤:  
1. 在vps上执行  
    `./proxy tbridge -p ":33080" -C proxy.crt -K proxy.key`  
  
1. 在公司机器A上面执行  
    `./proxy tclient -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  
  
1. 在家里电脑上执行  
    `./proxy tserver -r ":28080@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  
  
1. 完成  
  
#### **4.6、高级用法二**  
提示:  
如果同时有多个client连接到同一个bridge,需要指定不同的key,可以通过--k参数设定,--k可以是任意唯一字符串,  
只要在同一个bridge上唯一即可.  
server连接到bridge的时候,如果同时有多个client连接到同一个bridge,需要使用--k参数选择client.   
暴露多个端口重复-r参数即可.-r格式是:"本地IP:本地端口@clientHOST:client端口".   
  
背景:  
- 公司机器A提供了web服务80端口,ftp服务21端口  
- 有VPS一个,公网IP:22.22.22.22  
  
需求:  
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口  
在家里能够通过访问VPS的29090端口访问到公司机器A的21端口  
  
步骤:  
1. 在vps上执行  
    `./proxy tbridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy tserver -r ":28080@:80" -r ":29090@:21" --k test -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  

1. 在公司机器A上面执行  
    `./proxy tclient --k test -P "22.22.22.22:33080" -C proxy.crt -K proxy.key` 

1. 完成  
  
#### **4.7.tserver的-r参数**  
  -r完整格式是:`PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT`  
  
  4.7.1.协议PROTOCOL:tcp或者udp.  
  比如: `-r "udp://:10053@:53" -r "tcp://:10800@:1080" -r ":8080@:80"`  
  如果指定了--udp参数,PROTOCOL默认为udp,那么:`-r ":8080@:80"`默认为udp;  
  如果没有指定--udp参数,PROTOCOL默认为tcp,那么:`-r ":8080@:80"`默认为tcp;  
  
  4.7.2.CLIENT_KEY:默认是default.  
  比如: -r "udp://:10053@[test1]:53" -r "tcp://:10800@[test2]:1080" -r ":8080@:80"  
  如果指定了--k参数,比如--k test,那么:`-r ":8080@:80"`CLIENT_KEY默认为test;  
  如果没有指定--k参数,那么:`-r ":8080@:80"`CLIENT_KEY默认为default;  
  
  4.7.3.LOCAL_IP为空默认是:`0.0.0.0`,CLIENT_LOCAL_HOST为空默认是:`127.0.0.1`; 

#### **4.8.查看帮助**  
`./proxy help tbridge`  
`./proxy help tserver`  
`./proxy help tserver`  
  
### 5.SOCKS5代理  
提示:SOCKS5代理,只支持TCP协议,不支持UDP协议,不支持用户名密码认证.  
#### **5.1.普通SOCKS5代理**  
`./proxy socks -t tcp -p "0.0.0.0:38080"`  
  
#### **5.2.普通二级SOCKS5代理**  
使用本地端口8090,假设上级SOCKS5代理是`22.22.22.22:8080`  
`./proxy socks -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" `  
我们还可以指定网站域名的黑白名单文件,一行一个域名,怕匹配规则是最右批评匹配,比如:baidu.com,匹配的是*.*.baidu.com,黑名单的域名域名直接走上级代理,白名单的域名不走上级代理.  
`./proxy socks -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt`  
  
#### **5.3.SOCKS二级代理(加密)**  
一级SOCKS代理(VPS,IP:22.22.22.22)  
`./proxy socks -t tls -p ":38080" -C proxy.crt -K proxy.key`  
  
二级SOCKS代理(本地Linux)  
`./proxy socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080.  
  
二级SOCKS代理(本地windows)  
`./proxy.exe socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为socks5模式，地址为：127.0.0.1，端口为：8080,程序即可通过加密通道通过vps上网。  
  
#### **5.4.SOCKS三级代理(加密)**  
一级SOCKS代理VPS_01,IP:22.22.22.22  
`./proxy socks -t tls -p ":38080" -C proxy.crt -K proxy.key`  
二级SOCKS代理VPS_02,IP:33.33.33.33  
`./proxy socks -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级SOCKS代理(本地)  
`./proxy socks -t tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问一级SOCKS代理上面的代理端口38080.  
  
#### **5.5.SOCKS代理流量强制走上级SOCKS代理**  
默认情况下,proxy会智能判断一个网站域名是否无法访问,如果无法访问才走上级SOCKS代理.通过--always可以使全部SOCKS代理流量强制走上级SOCKS代理.  
`./proxy socks --always -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
  
#### **5.6.SOCKS通过SSH中转**  
说明:ssh中转的原理是利用了ssh的转发功能,就是你连接上ssh之后,可以通过ssh代理访问目标地址.  
假设有:vps  
- IP是2.2.2.2, ssh端口是22, ssh用户名是:user, ssh用户密码是:demo
- 用户user的ssh私钥名称是user.key   

##### ***5.6.1 ssh用户名和密码的方式***  
本地SOCKS5代理28080端口,执行:  
`./proxy socks -T ssh -P "2.2.2.2:22" -u user -A demo -t tcp -p ":28080"`  
##### ***5.6.2 ssh用户名和密钥的方式***  
本地SOCKS5代理28080端口,执行:  
`./proxy socks -T ssh -P "2.2.2.2:22" -u user -S user.key -t tcp -p ":28080"`  

那么访问本地的28080端口就是通过VPS访问目标地址.  

#### **5.7.认证**  
对于socks5代理协议我们可以进行用户名密码认证,认证的用户名和密码可以在命令行指定  
`./proxy socks -t tcp -p ":33080" -a "user1:pass1" -a "user2:pass2"`  
多个用户,重复-a参数即可.  
也可以放在文件中,格式是一行一个"用户名:密码",然后用-F指定.  
`./proxy socks -t tcp -p ":33080" -F auth-file.txt`  
如果没有-a或-F参数,就是关闭认证.  

#### **5.8.KCP协议传输**  
KCP协议需要-B参数设置一个密码用于加密解密数据  

一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy socks -t kcp -p ":38080" -B mypassword  
  
二级HTTP代理(本地Linux)  
`./proxy socks -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" -B mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080,数据通过kcp协议传输.  

#### **5.9.查看帮助**  
`./proxy help socks`  

### TODO  
- http,socks代理多个上级负载均衡?
- 内网穿透server<->bridge心跳机制?
- 欢迎加群反馈...

### 如何使用源码?   
cd进入你的go src目录,然后git clone https://github.com/snail007/goproxy.git ./proxy 即可.   
编译直接:go build     
运行: go run *.go    
utils是工具包,service是具体的每个服务类.   

### License  
Proxy is licensed under GPLv3 license.  
### Contact  
QQ交流群:189618940  
  
  
  
