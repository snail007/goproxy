<img src="https://github.com/snail007/goproxy/blob/master/docs/images/logo.jpg?raw=true" width="200"/>  
Proxy是golang实现的高性能http,https,websocket,tcp,udp,socks5,ss代理服务器,支持正向代理、反向代理、透明代理、内网穿透、TCP/UDP端口映射、SSH中转、TLS加密传输、协议转换、防污染DNS代理。

[点击下载](https://github.com/snail007/goproxy/releases) 官方QQ交流群:189618940  

---  
  
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https://img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total.svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style=plastic)](https://github.com/snail007/goproxy/releases)  
  
**[English Manual](/README.md)**  

**[全平台图形界面版本](/gui/README.md)**  

**[全平台SDK](/sdk/README.md)**

**[GoProxy特殊授权](/AUTHORIZATION.md)**

### 如何贡献代码(Pull Request)?  

欢迎加入一起发展壮大proxy.首先需要clone本项目到自己的帐号下面,   
然后在dev分支上面修改代码,最后发Pull Request到goproxy项目的dev分支即可,  
为了高效贡献代码,pr的时候需要说明做了什么变更,原因是什么.  

### Features  
- 链式代理,程序本身可以作为一级代理,如果设置了上级代理那么可以作为二级代理,乃至N级代理.  
- 通讯加密,如果程序不是一级代理,而且上级代理也是本程序,那么可以加密和上级代理之间的通讯,采用底层tls高强度加密,安全无特征.  
- 智能HTTP,SOCKS5代理,会自动判断访问的网站是否屏蔽,如果被屏蔽那么就会使用上级代理(前提是配置了上级代理)访问网站;如果访问的网站没有被屏蔽,为了加速访问,代理会直接访问网站,不使用上级代理.  
- 域名黑白名单，更加自由的控制网站的访问方式。  
- 跨平台性,无论你是widows,linux,还是mac,甚至是树莓派,都可以很好的运行proxy.  
- 多协议支持,支持HTTP(S),TCP,UDP,Websocket,SOCKS5代理. 
- TCP/UDP端口转发. 
- 支持内网穿透,协议支持TCP和UDP.  
- SSH中转,HTTP(S),SOCKS5代理支持SSH中转,上级Linux服务器不需要任何服务端,本地一个proxy即可开心上网.  
- [KCP](https://github.com/xtaci/kcp-go)协议支持,HTTP(S),SOCKS5代理支持KCP协议传输数据,降低延迟,提升浏览体验.  
- 集成外部API，HTTP(S),SOCKS5代理认证功能可以与外部HTTP API集成，可以方便的通过外部系统控制代理用户．  
- 反向代理,支持直接把域名解析到proxy监听的ip,然后proxy就会帮你代理访问需要访问的HTTP(S)网站.
- 透明HTTP(S)代理,配合iptables,在网关直接把出去的80,443方向的流量转发到proxy,就能实现无感知的智能路由器代理.  
- 协议转换，可以把已经存在的HTTP(S)或SOCKS5或SS代理转换为一个端口同时支持HTTP(S)和SOCKS5和SS代理，转换后的SOCKS5和SS代理如果上级是SOCKS5代理,那么支持UDP功能,同时支持强大的级联认证功能。
- 自定义底层加密传输，http(s)\sps\socks代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据,除此之外还支持在tls和kcp之后进行自定义加密,也就是说自定义加密和tls|kcp是可以联合使用的,内部采用AES256加密,使用的时候只需要自己定义一个密码即可。
- 底层压缩高效传输，http(s)\sps\socks代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据,在加密之后还可以对数据进行压缩,也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的。
- 安全的DNS代理，可以通过本地的proxy提供的DNS代理服务器与上级代理加密通讯实现安全防污染的DNS查询。
- 负载均衡,高可用,HTTP(S)\SOCKS5\SPS代理支持上级负载均衡和高可用,多个上级重复-P参数即可.
- 指定出口IP,HTTP(S)\SOCKS5\SPS代理支持客户端用入口IP连接过来的,就用入口IP作为出口IP访问目标网站的功能。如果入口IP是内网IP，出口IP不会使用入口IP
- 支持限速,HTTP(S)\SOCKS5\SPS代理支持限速.
- SOCKS5代理支持级联认证.
- 证书参数使用base64数据,默认情况下-C,-K参数是crt证书和key文件的路径,如果是base64://开头,那么就认为后面的数据是base64编码的,会解码后使用.


### Why need these?  
- 当由于某某原因,我们不能访问我们在其它地方的服务,我们可以通过多个相连的proxy节点建立起一个安全的隧道访问我们的服务.  
- 微信接口本地开发,方便调试.  
- 远程访问内网机器.  
- 和小伙伴一起玩局域网游戏.  
- 以前只能在局域网玩的,现在可以在任何地方玩.  
- 替代圣剑内网通，显IP内网通，花生壳之类的工具.
- ...  

 
本页是v5.4手册,其他版本手册请点击[这里](docs/old-release.md)查看. 
 

### 怎么找到组织?  
[点击加入交流组织gitter](https://gitter.im/go-proxy/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)  
[点击加入交流组织TG](https://t.me/joinchat/GYHXghCDSBmkKZrvu4wIdQ)  
### 安装 
1. [快速安装](#自动安装)
1. [手动安装](#手动安装)
1. [Docker安装](#docker安装)

### 首次使用必看
- [环境](#首次使用必看-1)
- [使用配置文件](#使用配置文件)
- [调试输出](#调试输出)
- [使用日志文件](#使用日志文件)
- [后台运行](#后台运行)
- [守护运行](#守护运行)
- [生成通讯证书文件](#生成加密通讯需要的证书文件)
- [安全建议](#安全建议)

### 手册目录
- [负载均衡和高可用](#负载均衡和高可用)
- [1. HTTP代理](#1http代理)
    - [1.1 普通HTTP代理](#11普通一级http代理)
    - [1.2 普通二级HTTP代理](#12普通二级http代理)
    - [1.3 HTTP二级代理(加密)](#13http二级代理加密)
    - [1.4 HTTP三级代理(加密)](#14http三级代理加密)
    - [1.5 Basic认证](#15basic认证)
    - [1.6 强制走上级HTTP代理](#16http代理流量强制走上级http代理)
    - [1.7 通过SSH中转](#17https通过ssh中转)
        - [1.7.1 用户名和密码的方式](#171-ssh用户名和密码的方式)
        - [1.7.2 用户名和密钥的方式](#172-ssh用户名和密钥的方式)
    - [1.8 KCP协议传输](#18kcp协议传输)
    - [1.9 HTTP(S)反向代理](#19-https反向代理)
    - [1.10 HTTP(S)透明代理](#110-https透明代理)
    - [1.11 自定义DNS](#111-自定义dns)
    - [1.12 自定义加密](#112-自定义加密)
    - [1.13 压缩传输](#113-压缩传输)
    - [1.14 负载均衡](#114-负载均衡)
    - [1.15 限速](#115-限速)
    - [1.16 指定出口IP](#116-指定出口ip)
    - [1.17 证书参数使用base64数据](#117-证书参数使用base64数据)
    - [1.18 查看帮助](#118-查看帮助)
- [2. TCP代理(端口映射)](#2tcp代理)
    - [2.1 普通一级TCP代理](#21普通一级tcp代理)
    - [2.2 普通二级TCP代理](#22普通二级tcp代理)
    - [2.3 普通三级TCP代理](#23普通三级tcp代理)
    - [2.4 加密二级TCP代理](#24加密二级tcp代理)
    - [2.5 加密三级TCP代理](#25加密三级tcp代理)
    - [2.6 通过代理连接上级](#26通过代理连接上级)
    - [2.7 查看帮助](#27查看帮助)
- [3. UDP代理(端口映射)](#3udp代理)
    - [3.1 普通一级UDP代理](#31普通一级udp代理)
    - [3.2 普通二级UDP代理](#32普通二级udp代理)
    - [3.3 普通三级UDP代理](#33普通三级udp代理)
    - [3.4 加密二级UDP代理](#34加密二级udp代理)
    - [3.5 加密三级UDP代理](#35加密三级udp代理)
    - [3.6 查看帮助](#36查看帮助)
- [4. 内网穿透](#4内网穿透)
    - [4.1 原理说明](#41原理说明)
    - [4.2 TCP普通用法](#42tcp普通用法)
    - [4.3 微信接口本地开发](#43微信接口本地开发)
    - [4.4 UDP普通用法](#44udp普通用法)
    - [4.5 高级用法一](#45高级用法一)
    - [4.6 高级用法一](#46高级用法二)
    - [4.7 server的-r参数](#47server的-r参数)
    - [4.8 server和client通过代理连接bridge](#48server和client通过代理连接bridge)
    - [4.9 查看帮助](#49查看帮助)
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
    - [5.9 自定义DNS](#59自定义dns)
    - [5.10 自定义加密](#510-自定义加密)
    - [5.11 压缩传输](#511-压缩传输)
    - [5.12 负载均衡](#512-负载均衡)
    - [5.13 限速](#513-限速)
    - [5.14 指定出口IP](#514-指定出口ip)
    - [5.15 级联认证](#515-级联认证)
    - [5.16 证书参数使用base64数据](#516-证书参数使用base64数据)
    - [5.17 查看帮助](#517查看帮助)
- [6. 代理协议转换](#6代理协议转换)
    - [6.1 功能介绍](#61-功能介绍)
    - [6.2 HTTP(S)转HTTP(S)+SOCKS5+SS](#62-https转httpssocks5ss)
    - [6.3 SOCKS5转HTTP(S)+SOCKS5+SS](#63-socks5转httpssocks5ss)
    - [6.4 SS转HTTP(S)+SOCKS5+SS](#64-ss转httpssocks5ss)
    - [6.5 链式连接](#65-链式连接)
    - [6.6 监听多个端口](#66-监听多个端口)
    - [6.7 认证功能](#67-认证功能)
    - [6.8 自定义加密](#68-自定义加密)
    - [6.9 压缩传输](#69-压缩传输)
    - [6.10 禁用协议](#610-禁用协议)
    - [6.11 限速](#611-限速)
    - [6.12 指定出口IP](#612-指定出口ip)
    - [6.13 证书参数使用base64数据](#613-证书参数使用base64数据)
    - [6.14 查看帮助](#614-查看帮助)
- [7. KCP配置](#7kcp配置)
    - [7.1 配置介绍](#71-配置介绍)
    - [7.2 详细配置](#72-详细配置)
- [8. DNS防污染服务器](#8dns防污染服务器)
    - [8.1 介绍](#81-介绍)
    - [8.2 使用示例](#82-使用示例)

### Fast Start  
提示:所有操作需要root权限.  
#### 自动安装
#### **0.如果你的VPS是linux64位的系统,那么只需要执行下面一句,就可以完成自动安装和配置.**  

```shell  
curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash  
```  

安装完成,配置目录是/etc/proxy,更详细的使用方法请参考上面的手册目录,进一步了解你想要使用的功能.  
如果安装失败或者你的vps不是linux64位系统,请按照下面的半自动步骤安装:  
  
#### 手动安装  

#### **1.下载proxy**  
下载地址:https://github.com/snail007/goproxy/releases/latest   
下面以v6.2为例,如果有最新版,请使用最新版链接.   

```shell  
cd /root/proxy/  
wget https://github.com/snail007/goproxy/releases/download/v6.2/proxy-linux-amd64.tar.gz  
```  

#### **2.下载自动安装脚本**  

```shell  
cd /root/proxy/  
wget https://raw.githubusercontent.com/snail007/goproxy/master/install.sh  
chmod +x install.sh  
./install.sh  
```  

#### Docker安装 
项目根目录的Dockerfile文件用来构建,使用golang 1.10.3,构建基于goproxy的master分支最新版本,  
全部大小17.3MB,默认情况下使用master分支,不过可以通过修改配置文件Dockerfile  
或者使用参数GOPROXY_VERSION指定构建的goproxy版本.  

```
ARG GOPROXY_VERSION=v5.3
```

步骤:  
1. 克隆仓库,然后cd进入仓库文件夹,执行:
```
sudo docker build .
```
2. 镜像打标签:
```
sudo docker tag <上一步的结果ID> snail007/goproxy:latest
```
3. 运行 
参数OPTS的值就是传递给proxy的所有参数
比如下面的例子启动了一个http服务:

```
sudo docker run -d --restart=always --name goproxy -e OPTS="http -p :33080" -p 33080:33080 snail007/goproxy:latest
```
4. 查看日志:
```
sudo docker logs -f goproxy
```

## **首次使用必看**  
  
### **环境**  
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
### **调试输出**   
默认情况下,日志输出的信息不包含文件行数,某些情况下为了排除程序问题,快速定位问题,  
可以使用--debug参数,输出代码行数和毫秒时间.  

### **使用日志文件**   
默认情况下,日志是直接在控制台显示出来的,如果要保存到文件,可以使用--log参数,  
比如: --log proxy.log,日志就会输出到proxy.log方便排除问题.   


### **生成加密通讯需要的证书文件**  
http,tcp,udp代理过程会和上级通讯,为了安全我们采用加密通讯,当然可以选择不加密通信通讯,本教程所有和上级通讯都采用加密,需要证书文件.  

1.通过下面的命令生成自签名的证书和key文件.  
`./proxy keygen -C proxy`  
会在当前程序目录下面生成证书文件proxy.crt和key文件proxy.key。   

2.通过下面的命令生,使用自签名证书proxy.crt和key文件proxy.key签发新证书:goproxy.crt和goproxy.key.   
`./proxy keygen -s -C proxy -c goproxy`  
会在当前程序目录下面生成证书文件goproxy.crt和key文件goproxy.key。   

3.默认情况下证书的里面的域名是随机的,可以使用`-n test.com`参数指定.  

4.更多用法:`proxy keygen --help`。    

### **后台运行**
默认执行proxy之后,如果要保持proxy运行,不能关闭命令行.  
如果想在后台运行proxy,命令行可以关闭,只需要在命令最后加上--daemon参数即可.  
比如:  
`./proxy http -t tcp -p "0.0.0.0:38080" --daemon`   

### **守护运行**  
守护运行参数--forever,比如: `proxy http --forever` ,  
proxy会fork子进程,然后监控子进程,如果子进程异常退出,5秒后重启子进程.  
该参数配合后台运行参数--daemon和日志参数--log,可以保障proxy一直在后台执行不会因为意外退出,  
而且可以通过日志文件看到proxy的输出日志内容.  
比如: `proxy http -p ":9090" --forever --log proxy.log --daemon`  

### **安全建议**
当VPS在nat设备后面,vps上网卡IP都是内网IP,这个时候可以通过-g参数添加vps的外网ip防止死循环.  
假设你的vps外网ip是23.23.23.23,下面命令通过-g参数设置23.23.23.23  
`./proxy http -g "23.23.23.23"`  

### **负载均衡和高可用**

HTTP(S)\SOCKS5\SPS代理支持上级负载均衡和高可用,多个上级重复-P参数即可.

负载均衡策略支持5种,可以通过`--lb-method`参数指定:

roundrobin 轮流使用

leastconn  使用最小连接数的

leasttime  使用连接时间最小的

hash     使用根据客户端地址计算出一个固定上级

weight    根据每个上级的权重和连接数情况,选择出一个上级

提示:

负载均衡检查时间间隔可以通过`--lb-retrytime`设置,单位毫秒

负载均衡连接超时时间可以通过`--lb-timeout`设置,单位毫秒

如果负载均衡策略是权重(weight),-P格式为:2.2.2.2:3880@1,1就是权重,大于0的整数.

如果负载均衡策略是hash,默认是根据客户端地址选择上级,可以通过开关`--lb-hashtarget`使用访问的目标地址选择上级.

### **1.HTTP代理**  
#### **1.1.普通一级HTTP代理**  
![1.1](/docs/images/http-1.png)  
`./proxy http -t tcp -p "0.0.0.0:38080"`  
  
#### **1.2.普通二级HTTP代理**  
![1.2](/docs/images/http-2.png)  
使用本地端口8090,假设上级HTTP代理是`22.22.22.22:8080`  
`./proxy http -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" `  
我们还可以指定网站域名的黑白名单文件,一行一个域名,匹配规则是最右匹配,比如:baidu.com,匹配的是*.*.baidu.com,黑名单的域名直接走上级代理,白名单的域名不走上级代理.  
`./proxy http -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt`  
  
#### **1.3.HTTP二级代理(加密)**  
> 注意: 后面二级代理使用的`proxy.crt`和`proxy.key`应与一级代理一致  

![1.3](/docs/images/http-tls-2.png)  
一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key`  
  
二级HTTP代理(本地Linux)  
`./proxy http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080.  
  
二级HTTP代理(本地windows)  
`./proxy.exe http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为http模式，地址为：127.0.0.1，端口为：8080,程序即可通过加密通道通过vps上网。  
  
#### **1.4.HTTP三级代理(加密)**  
![1.3](/docs/images/http-tls-3.png)  
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
  
另外,http(s)代理还集成了外部HTTP API认证,我们可以通过--auth-url参数指定一个http url接口地址,  
然后有用户连接的时候,proxy会GET方式请求这url,带上下面四个参数,如果返回HTTP状态码204,代表认证成功  
其它情况认为认证失败.  
比如:  
`./proxy http -t tcp -p ":33080" --auth-url "http://test.com/auth.php"`  
用户连接的时候,proxy会GET方式请求这url("http://test.com/auth.php"),  
带上user,pass,ip,target四个参数:  
http://test.com/auth.php?user={USER}&pass={PASS}&ip={IP}&target={TARGET}  
user:用户名  
pass:密码  
ip:用户的IP,比如:192.168.1.200  
target:用户访问的URL,比如:http://demo.com:80/1.html或https://www.baidu.com:80  

如果没有-a或-F或--auth-url参数,就是关闭Basic认证.   

#### **1.6.HTTP代理流量强制走上级HTTP代理**  
默认情况下,proxy会智能判断一个网站域名是否无法访问,如果无法访问才走上级HTTP代理.通过--always可以使全部HTTP代理流量强制走上级HTTP代理.  
`./proxy http --always -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
  
#### **1.7.HTTP(S)通过SSH中转**  
![1.7](/docs/images/http-ssh-1.png)  
说明:ssh中转的原理是利用了ssh的转发功能,就是你连接上ssh之后,可以通过ssh代理访问目标地址.  
假设有:vps  
- IP是2.2.2.2, ssh端口是22, ssh用户名是:user, ssh用户密码是:demo  
- 用户user的ssh私钥名称是user.key    

##### ***1.7.1 ssh用户名和密码的方式***   
本地HTTP(S)代理28080端口,执行:  
`./proxy http -T ssh -P "2.2.2.2:22" -u user -D demo -t tcp -p ":28080"`  
##### ***1.7.2 ssh用户名和密钥的方式***   
本地HTTP(S)代理28080端口,执行:  
`./proxy http -T ssh -P "2.2.2.2:22" -u user -S user.key -t tcp -p ":28080"`  

#### **1.8.KCP协议传输**  
![1.8](/docs/images/http-kcp.png)  
KCP协议需要--kcp-key参数设置一个密码用于加密解密数据  

一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy http -t kcp -p ":38080" --kcp-key mypassword`  
  
二级HTTP代理(本地Linux)  
`./proxy http -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" --kcp-key mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080,数据通过kcp协议传输,注意kcp走的是udp协议协议,所以防火墙需放开38080的udp协议.   

#### **1.9 HTTP(S)反向代理**   
![1.9](/docs/images/fxdl.png)  
proxy不仅支持在其他软件里面通过设置代理的方式,为其他软件提供代理服务,而且支持直接把请求的网站域名解析到proxy监听的ip上,然后proxy监听80和443端口,那么proxy就会自动为你代理访问需要访问的HTTP(S)网站.  

使用方式:  
在"最后一级proxy代理"的机器上,因为proxy要伪装成所有网站,网站默认的端口HTTP是80,HTTPS是443,让proxy监听80和443端口即可.参数-p多个地址用逗号分割.  
`./proxy http -t tcp -p :80,:443`    

这个命令就在机器上启动了一个proxy代理,同时监听80和443端口,既可以当作普通的代理使用,也可以直接把需要代理的域名解析到这个机器的IP上. 

如果有上级代理那么参照上面教程设置上级即可,使用方式完全一样.  
`./proxy http -t tcp -p :80,:443 -T tls -P "2.2.2.2:33080" -C proxy.crt -K proxy.key`   

注意:  
proxy所在的服务器的DNS解析结果不能受到自定义的解析影响,不然就死循环了,proxy代理最好指定`--dns 8.8.8.8`参数.  
  
#### **1.10 HTTP(S)透明代理** 
该模式需要具有一定的网络基础,相关概念不懂的请自行搜索解决.  
假设proxy现在在路由器上运行,启动命令如下:  
`./proxy http -t tcp -p :33080 -T tls -P "2.2.2.2:33090" -C proxy.crt -K proxy.key`   

然后添加iptables规则,下面是参考规则:  
```shell
#上级proxy服务端服务器IP地址:
proxy_server_ip=2.2.2.2

#路由器运行proxy监听的端口:
proxy_local_port=33080

#下面的就不用修改了
#create a new chain named PROXY
iptables -t nat -N PROXY

# Ignore your PROXY server's addresses
# It's very IMPORTANT, just be careful.

iptables -t nat -A PROXY -d $proxy_server_ip -j RETURN

# Ignore LANs IP address
iptables -t nat -A PROXY -d 0.0.0.0/8 -j RETURN
iptables -t nat -A PROXY -d 10.0.0.0/8 -j RETURN
iptables -t nat -A PROXY -d 127.0.0.0/8 -j RETURN
iptables -t nat -A PROXY -d 169.254.0.0/16 -j RETURN
iptables -t nat -A PROXY -d 172.16.0.0/12 -j RETURN
iptables -t nat -A PROXY -d 192.168.0.0/16 -j RETURN
iptables -t nat -A PROXY -d 224.0.0.0/4 -j RETURN
iptables -t nat -A PROXY -d 240.0.0.0/4 -j RETURN

# Anything to port 80 443 should be redirected to PROXY's local port
iptables -t nat -A PROXY -p tcp --dport 80 -j REDIRECT --to-ports $proxy_local_port
iptables -t nat -A PROXY -p tcp --dport 443 -j REDIRECT --to-ports $proxy_local_port

# Apply the rules to nat client
iptables -t nat -A PREROUTING -p tcp -j PROXY
# Apply the rules to localhost
iptables -t nat -A OUTPUT -p tcp -j PROXY
```
- 清空整个链 iptables -F 链名比如iptables -t nat -F PROXY
- 删除指定的用户自定义链 iptables -X 链名 比如 iptables -t nat -X PROXY
- 从所选链中删除规则 iptables -D 链名 规则详情 比如 iptables -t nat -D PROXY -d 223.223.192.0/255.255.240.0 -j RETURN

#### **1.11 自定义DNS** 
--dns-address和--dns-ttl参数,用于自己指定proxy访问域名的时候使用的dns（--dns-address）  
以及解析结果缓存时间（--dns-ttl）秒数，避免系统dns对proxy的干扰，另外缓存功能还能减少dns解析时间提高访问速度.    
比如：  
`./proxy http -p ":33080" --dns-address "8.8.8.8:53" --dns-ttl 300`  

#### **1.12 自定义加密**  
proxy的http(s)代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据,除此之外还支持在tls和kcp之后进行自定义  
加密,也就是说自定义加密和tls|kcp是可以联合使用的,内部采用AES256加密,使用的时候只需要自己定义一个密码即可,  
加密分为两个部分,一部分是本地(-z)是否加密解密,一部分是与上级(-Z)传输是否加密解密.    
自定义加密要求两端都是proxy才可以,下面分别用二级,三级为例:  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy http -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy http -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888` 
本地三级执行:  
`proxy http -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  

#### **1.13 压缩传输**  
proxy的http(s)代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据,在自定义加密之前还可以对数据进行压缩,  
也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的,压缩分为两个部分,一部分是本地(-m)是否压缩传输,  
一部分是与上级(-M)传输是否压缩.    
压缩要求两端都是proxy才可以,压缩也在一定程度上保护了(加密)数据,下面分别用二级,三级为例:  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -m -p :7777`  
本地二级执行:  
`proxy http -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy http -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888` 
本地三级执行:  
`proxy http -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  

### **1.14 负载均衡**  

HTTP(S)代理支持上级负载均衡,多个上级重复-P参数即可.

`proxy http --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080`

#### **1.14.1 设置重试间隔和超时时间**  

`proxy http --lb-method=leastconn --lb-retrytime 300 --lb-timeout 300 -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -t tcp -p :33080`

#### **1.14.2 设置权重**  

`proxy http --lb-method=weight -T tcp -P 1.1.1.1:33080@1 -P 2.1.1.1:33080@2 -P 3.1.1.1:33080@1 -t tcp -p :33080`

#### **1.14.3 使用目标地址选择上级**  

`proxy http --lb-hashtarget --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -t tcp -p :33080`

### **1.15 限速**  

限速100K,通过`-l`参数即可指定,比如:100K 1.5M . 0意味着无限制.

`proxy http -t tcp -p 2.2.2.2:33080 -l 100K`

### **1.16 指定出口IP**  

`--bind-listen`参数，就可以开启客户端用入口IP连接过来的,就用入口IP作为出口IP访问目标网站的功能。如果入口IP是内网IP，出口IP不会使用入口IP。

`proxy http -t tcp -p 2.2.2.2:33080 --bind-listen`

### **1.17 证书参数使用base64数据**  

默认情况下-C,-K参数是crt证书和key文件的路径,

如果是base64://开头,那么就认为后面的数据是base64编码的,会解码后使用.

#### **1.18 查看帮助**  
`./proxy help http`  
  
### **2.TCP代理**  
  
#### **2.1.普通一级TCP代理**  
![2.1](/docs/images/tcp-1.png)  
本地执行:  
`./proxy tcp -p ":33080" -T tcp -P "192.168.22.33:22"`  
那么访问本地33080端口就是访问192.168.22.33的22端口.  
  
#### **2.2.普通二级TCP代理**  
![2.2](/docs/images/tcp-2.png)  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -p ":33080" -T tcp -P "127.0.0.1:8080"`  
本地执行:  
`./proxy tcp -p ":23080" -T tcp -P "22.22.22.33:33080"`  
那么访问本地23080端口就是访问22.22.22.33的8080端口.  
  
#### **2.3.普通三级TCP代理**  
![2.3](/docs/images/tcp-3.png)  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -p ":38080" -T tcp -P "66.66.66.66:8080"`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`./proxy tcp -p ":8080" -T tcp -P "33.33.33.33:28080"`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口.  
  
#### **2.4.加密二级TCP代理**  
![2.4](/docs/images/tcp-tls-2.png)  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -t tls -p ":33080" -T tcp -P "127.0.0.1:8080" -C proxy.crt -K proxy.key`  
本地执行:  
`./proxy tcp -p ":23080" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地23080端口就是通过加密TCP隧道访问22.22.22.33的8080端口.  
  
#### **2.5.加密三级TCP代理**  
![2.5](/docs/images/tcp-tls-3.png)  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -t tls -p ":38080" -T tcp -P "66.66.66.66:8080" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`./proxy tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口.  
  
#### **2.6.通过代理连接上级**  
有时候proxy所在的网络不能直接访问外网,需要通过一个https或者socks5代理才能上网,那么这个时候  
-J参数就可以帮助你让proxy的tcp端口映射的时候通过https或者socks5代理去连接上级-P,将外部端口映射到本地.    
-J参数格式如下:  

https代理写法:  
代理需要认证,用户名:username 密码:password  
https://username:password@host:port  
代理不需要认证  
https://host:port  

socks5代理写法:
代理需要认证,用户名:username 密码:password
socks5://username:password@host:port
代理不需要认证
socks5://host:port

host:代理的IP或者域名
port:代理的端口

#### **2.7.查看帮助**  
`./proxy help tcp`  
  
### **3.UDP代理**  
  
#### **3.1.普通一级UDP代理**  
![3.1](/docs/images/udp-1.png)  
本地执行:  
`./proxy udp -p ":5353" -T udp -P "8.8.8.8:53"`  
那么访问本地UDP:5353端口就是访问8.8.8.8的UDP:53端口.  
  
#### **3.2.普通二级UDP代理**  
![3.2](/docs/images/udp-2.png)  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -p ":33080" -T udp -P "8.8.8.8:53"`  
本地执行:  
`./proxy udp -p ":5353" -T tcp -P "22.22.22.33:33080"`  
那么访问本地UDP:5353端口就是通过TCP隧道,通过VPS访问8.8.8.8的UDP:53端口.  
  
#### **3.3.普通三级UDP代理**  
![3.3](/docs/images/udp-3.png)  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -p ":38080" -T udp -P "8.8.8.8:53"`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`./proxy udp -p ":5353" -T tcp -P "33.33.33.33:28080"`  
那么访问本地5353端口就是通过TCP隧道,通过VPS访问8.8.8.8的53端口.  
  
#### **3.4.加密二级UDP代理**  
![3.4](/docs/images/udp-tls-2.png)  
VPS(IP:22.22.22.33)执行:  
`./proxy tcp -t tls -p ":33080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
本地执行:  
`./proxy udp -p ":5353" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地UDP:5353端口就是通过加密TCP隧道,通过VPS访问8.8.8.8的UDP:53端口.  
  
#### **3.5.加密三级UDP代理**  
![3.5](/docs/images/udp-tls-3.png)  
一级TCP代理VPS_01,IP:22.22.22.22  
`./proxy tcp -t tls -p ":38080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02,IP:33.33.33.33  
`./proxy tcp -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`./proxy udp -p ":5353" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地5353端口就是通过加密TCP隧道,通过VPS_01访问8.8.8.8的53端口.  
  
#### **3.6.查看帮助**  
`./proxy help udp`  
  
### **4.内网穿透**  
#### **4.1、原理说明**  
内网穿透,分为两个版本，“多链接版本”和“多路复用版本”，一般像web服务这种不是长时间连接的服务建议用“多链接版本”，如果是要保持长时间连接建议使用“多路复用版本”。
1. 多链接版本，对应的子命令是tserver，tclient，tbridge。  
1. 多路复用版本，对应的子命令是server，client，bridge。  
1. 多链接版本和多路复用版本的参数和使用方式完全一样。  
1. **多路复用版本的server，client可以开启压缩传输，参数是--c。**   
1. **server，client要么都开启压缩，要么都不开启，不能只开一个。**    

下面的教程以“多路复用版本”为例子，说明使用方法。    
内网穿透由三部分组成:client端,server端,bridge端；client和server主动连接bridge端进行桥接.    
当用户访问server端,流程是:   
1. 首先server端主动和bridge端建立连接；  
1. 然后bridge端通知client端连接bridge端和目标端口;  
1. 然后client端绑定“client端到bridge端”和“client端到目标端口”的连接；  
1. 然后bridge端把“client过来的连接”与“server端过来的连接”绑定；  
1. 整个通道建立完成；  
  
#### **4.2、TCP普通用法**  
背景:  
- 公司机器A提供了web服务80端口  
- 有VPS一个,公网IP:22.22.22.22  

需求:  
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口  
  
步骤:  
1. 在vps上执行  
    `./proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy server -r ":28080@:80" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  
  
1. 在公司机器A上面执行  
    `./proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

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
    `./proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy server -r ":80@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 在自己笔记本上面执行  
    `./proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

1. 完成  
  
#### **4.4、UDP普通用法**  
背景:  
- 公司机器A提供了DNS解析服务,UDP:53端口  
- 有VPS一个,公网IP:22.22.22.22  
  
需求:  
在家里能够通过设置本地dns为22.22.22.22,使用公司机器A进行域名解析服务.  
  
步骤:  
1. 在vps上执行  
    `./proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy server --udp -r ":53@:53" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  

1. 在公司机器A上面执行  
    `./proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  

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
    `./proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
  
1. 在公司机器A上面执行  
    `./proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  
  
1. 在家里电脑上执行  
    `./proxy server -r ":28080@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`  
  
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
    `./proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
    `./proxy server -r ":28080@:80" -r ":29090@:21" --k test -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`  

1. 在公司机器A上面执行  
    `./proxy client --k test -P "22.22.22.22:33080" -C proxy.crt -K proxy.key` 

1. 完成  
  
#### **4.7.server的-r参数**  
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

#### **4.8.server和client通过代理连接bridge**  
有时候server或者client所在的网络不能直接访问外网,需要通过一个https或者socks5代理才能上网,那么这个时候  
-J参数就可以帮助你让server或者client通过https或者socks5代理去连接bridge.  
-J参数格式如下:  

https代理写法:  
代理需要认证,用户名:username 密码:password  
https://username:password@host:port  
代理不需要认证  
https://host:port  

socks5代理写法:
代理需要认证,用户名:username 密码:password
socks5://username:password@host:port
代理不需要认证
socks5://host:port

host:代理的IP或者域名
port:代理的端口

#### **4.9.查看帮助**  
`./proxy help bridge`  
`./proxy help server`  
`./proxy help client`  
  
### **5.SOCKS5代理**  
提示:

SOCKS5代理,支持CONNECT,UDP协议,不支持BIND,支持用户名密码认证.  

***如果你的VPS是阿里云，腾讯云这种VPS，就是ifconfig看不见你的公网IP，只能看见内网IP，***

***那么需要加上`-g VPS公网IP`参数，SOCKS5代理的UDP功能才能正常工作。***

#### **5.1.普通SOCKS5代理**  
`./proxy socks -t tcp -p "0.0.0.0:38080"`  
  
#### **5.2.普通二级SOCKS5代理**  
![5.2](/docs/images/socks-2.png)  
使用本地端口8090,假设上级SOCKS5代理是`22.22.22.22:8080`  
`./proxy socks -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" `  
我们还可以指定网站域名的黑白名单文件,一行一个域名,匹配规则是最右匹配,比如:baidu.com,匹配的是*.*.baidu.com,黑名单的域名域名直接走上级代理,白名单的域名不走上级代理;如果域名即在黑名单又在白名单中,那么黑名单起作用.  
`./proxy socks -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt`  
  
#### **5.3.SOCKS二级代理(加密)**  
![5.3](/docs/images/socks-tls-2.png)  
一级SOCKS代理(VPS,IP:22.22.22.22)  
`./proxy socks -t tls -p ":38080" -C proxy.crt -K proxy.key`  
  
二级SOCKS代理(本地Linux)  
`./proxy socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080.  
  
二级SOCKS代理(本地windows)  
`./proxy.exe socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为socks5模式，地址为：127.0.0.1，端口为：8080,程序即可通过加密通道通过vps上网。  
  
#### **5.4.SOCKS三级代理(加密)**  
![5.4](/docs/images/socks-tls-3.png)  
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
![5.6](/docs/images/socks-ssh.png)  
说明:ssh中转的原理是利用了ssh的转发功能,就是你连接上ssh之后,可以通过ssh代理访问目标地址.  
假设有:vps  
- IP是2.2.2.2, ssh端口是22, ssh用户名是:user, ssh用户密码是:demo
- 用户user的ssh私钥名称是user.key   

##### ***5.6.1 ssh用户名和密码的方式***  
本地SOCKS5代理28080端口,执行:  
`./proxy socks -T ssh -P "2.2.2.2:22" -u user -D demo -t tcp -p ":28080"`  
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

另外,socks5代理还集成了外部HTTP API认证,我们可以通过--auth-url参数指定一个http url接口地址,  
然后有用户连接的时候,proxy会GET方式请求这url,带上下面三个参数,如果返回HTTP状态码204,代表认证成功  
其它情况认为认证失败.  
比如:  
`./proxy socks -t tcp -p ":33080" --auth-url "http://test.com/auth.php"`  
用户连接的时候,proxy会GET方式请求这url("http://test.com/auth.php"),  
带上user,pass,ip,三个参数:  
http://test.com/auth.php?user={USER}&pass={PASS}&ip={IP}  
user:用户名  
pass:密码  
ip:用户的IP,比如:192.168.1.200  

如果没有-a或-F或--auth-url参数,就是关闭认证.    

#### **5.8.KCP协议传输**  
KCP协议需要--kcp-key参数设置一个密码用于加密解密数据  

一级HTTP代理(VPS,IP:22.22.22.22)  
`./proxy socks -t kcp -p ":38080" --kcp-key mypassword`  
  
二级HTTP代理(本地Linux)  
`./proxy socks -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" --kcp-key mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080,数据通过kcp协议传输.  

#### **5.9.自定义DNS** 
--dns-address和--dns-ttl参数,用于自己指定proxy访问域名的时候使用的dns（--dns-address）  
以及解析结果缓存时间（--dns-ttl）秒数，避免系统dns对proxy的干扰，另外缓存功能还能减少dns解析时间提高访问速度.    
比如：  
`./proxy socks -p ":33080" --dns-address "8.8.8.8:53" --dns-ttl 300`  

#### **5.10 自定义加密**  
proxy的socks代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据,除此之外还支持在tls和kcp之后进行自定义加密,也就是说自定义加密和tls|kcp是可以联合使用的,内部采用AES256加密,使用的时候只需要自己定义一个密码即可,  
加密分为两个部分,一部分是本地(-z)是否加密解密,一部分是与上级(-Z)传输是否加密解密.    

自定义加密要求两端都是proxy才可以.  

下面分别用二级,三级为例:  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy socks -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy socks -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888` 
本地三级执行:  
`proxy socks -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  

#### **5.11 压缩传输**  
proxy的socks代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据,在自定义加密之前还可以  
对数据进行压缩,也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的,压缩分为两个部分,   
一部分是本地(-m)是否压缩传输,一部分是与上级(-M)传输是否压缩.    

压缩要求两端都是proxy才可以,压缩也在一定程度上保护了(加密)数据.  

下面分别用二级,三级为例:  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -m -p :7777`  
本地二级执行:  
`proxy socks -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy socks -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888` 
本地三级执行:  
`proxy socks -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  


#### **5.12 负载均衡**  

SOCKS代理支持上级负载均衡,多个上级重复-P参数即可.

`proxy socks --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080  -p :33080 -t tcp`

#### **5.12.1 设置重试间隔和超时时间**  

`proxy socks --lb-method=leastconn --lb-retrytime 300 --lb-timeout 300 -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -p :33080 -t tcp`

#### **5.12.2 设置权重**  

`proxy socks --lb-method=weight -T tcp -P 1.1.1.1:33080@1 -P 2.1.1.1:33080@2 -P 3.1.1.1:33080@1 -p :33080 -t tcp`

#### **5.12.3 使用目标地址选择上级**  

`proxy socks --lb-hashtarget --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -p :33080 -t tcp`

#### **5.13 限速**  

限速100K,通过`-l`参数即可指定,比如:100K 1.5M . 0意味着无限制.

`proxy socks -t tcp -p 2.2.2.2:33080 -l 100K`

#### **5.14 指定出口IP**  

`--bind-listen`参数，就可以开启客户端用入口IP连接过来的,就用入口IP作为出口IP访问目标网站的功能。如果入口IP是内网IP，出口IP不会使用入口IP。

`proxy socks -t tcp -p 2.2.2.2:33080 --bind-listen`

#### **5.15 级联认证**  

SOCKS5支持级联认证,-A可以设置上级认证信息.

上级:

`proxy socks -t tcp -p 2.2.2.2:33080 -a user:pass`

本地:

`proxy socks -T tcp -P 2.2.2.2:33080 -A user:pass -t tcp -p :33080`

#### **5.16 证书参数使用base64数据**  

默认情况下-C,-K参数是crt证书和key文件的路径,

如果是base64://开头,那么就认为后面的数据是base64编码的,会解码后使用.


#### **5.17.查看帮助**  
`./proxy help socks`  

### **6.代理协议转换** 

#### **6.1 功能介绍** 
代理协议转换使用的是sps子命令，sps本身不提供代理功能，只是接受代理请求"转换并转发"给已经存在的http(s)代理或者socks5代理或者ss代理；sps可以把已经存在的http(s)代理或者socks5代理或ss代理转换为一个端口同时支持http(s)和socks5和ss的代理，而且http(s)代理支持正向代理和反向代理(SNI)，转换后的SOCKS5代理，当上级是SOCKS5或者SS时仍然支持UDP功能；另外对于已经存在的http(s)代理或者socks5代理，支持tls、tcp、kcp三种模式，支持链式连接，也就是可以多个sps结点层级连接构建加密通道。

#### **6.2 HTTP(S)转HTTP(S)+SOCKS5+SS** 
假设已经存在一个普通的http(s)代理：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080,ss加密方式:aes-192-cfb,ss密码:pass。  
命令如下：  
`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`

假设已经存在一个tls的http(s)代理：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080，tls需要证书文件,ss加密方式:aes-192-cfb,ss密码:pass。  
命令如下：  
`./proxy sps -S http -T tls -P 127.0.0.1:8080 -t tcp -p :18080 -C proxy.crt -K proxy.key -h aes-192-cfb -j pass`   

假设已经存在一个kcp的http(s)代理（密码是：demo123）：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080,ss加密方式:aes-192-cfb,ss密码:pass。  
命令如下：  
`./proxy sps -S http -T kcp -P 127.0.0.1:8080 -t tcp -p :18080 --kcp-key demo123 -h aes-192-cfb -j pass`  

#### **6.3 SOCKS5转HTTP(S)+SOCKS5+SS** 
假设已经存在一个普通的socks5代理：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080,ss加密方式:aes-192-cfb,ss密码:pass。    
命令如下：  
`./proxy sps -S socks -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`

假设已经存在一个tls的socks5代理：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080，tls需要证书文件,ss加密方式:aes-192-cfb,ss密码:pass。   
命令如下：  
`./proxy sps -S socks -T tls -P 127.0.0.1:8080 -t tcp -p :18080 -C proxy.crt -K proxy.key -h aes-192-cfb -j pass`   

假设已经存在一个kcp的socks5代理（密码是：demo123）：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080,ss加密方式:aes-192-cfb,ss密码:pass。  
命令如下：  
`./proxy sps -S socks -T kcp -P 127.0.0.1:8080 -t tcp -p :18080 --kcp-key demo123 -h aes-192-cfb -j pass`  

#### **6.4 SS转HTTP(S)+SOCKS5+SS** 
SPS上级和本地支持ss协议,上级可以是SPS或者标准的ss服务.  
SPS本地默认提供HTTP(S)\SOCKS5\SPS三种代理,当上级是SOCKS5时转换后的SOCKS5和SS支持UDP功能.  
假设已经存在一个普通的SS或者SPS代理(开启了ss,加密方式:aes-256-cfb,密码:demo)：127.0.0.1:8080,现在我们把它转为同时支持http(s)和socks5和ss的普通代理,转换后的本地端口为18080,转换后的ss加密方式:aes-192-cfb,ss密码:pass。  
命令如下：  
`./proxy sps -S ss -H aes-256-cfb -J pass -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`.  

#### **6.5 链式连接**   
![6.4](/docs/images/sps-tls.png)  
上面提过多个sps结点可以层级连接构建加密通道，假设有如下vps和家里的pc电脑。  
vps01：2.2.2.2  
vps02：3.3.3.3  
现在我们想利用pc和vps01和vps02构建一个加密通道，本例子用tls加密也可以用kcp，在pc上访问本地18080端口就是访问vps01的本地8080端口。  
首先在vps01(2.2.2.2)上我们运行一个只有本地可以访问的http(s)代理,执行：  
`./proxy http -t tcp -p 127.0.0.1:8080`  

然后在vps01(2.2.2.2)上运行一个sps结点，执行：  
`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tls -p :8081 -C proxy.crt -K proxy.key`  

然后在vps02(3.3.3.3)上运行一个sps结点，执行：  
`./proxy sps -S http -T tls -P 2.2.2.2:8081 -t tls -p :8082 -C proxy.crt -K proxy.key`  

然后在pc上运行一个sps结点，执行：  
`./proxy sps -S http -T tls -P 3.3.3.3:8082 -t tcp -p :18080 -C proxy.crt -K proxy.key`  

完成。  

#### **6.6 监听多个端口**   
一般情况下监听一个端口就可以，不过如果作为反向代理需要同时监听80和443两个端口，那么-p参数是支持的，  
格式是：`-p 0.0.0.0:80,0.0.0.0:443`，多个绑定用逗号分隔即可。  

#### **6.7 认证功能**   
sps支持http(s)\socks5代理认证,可以级联认证,有四个重要的信息:  
1:用户发送认证信息`user-auth`。   
2:设置的本地认证信息`local-auth`。  
3:设置的连接上级使用的认证信息`parent-auth`。  
4:最终发送给上级的认证信息`auth-info-to-parent`。  
他们的情况关系如下:  

| user-auth | local-auth | parent-auth | auth-info-to-paren 
| ------ | ------ | ------ | ------  
| 有/没有  | 有     |     有   |   来自parent-auth  
| 有/没有  | 没有    |    有    |   来自parent-auth  
| 有/没有  | 有     |     没有  |   无  
| 没有   | 没有    |   没有    |   无  
| 有    | 没有    |   没有    |   来自user-auth  

对于sps代理我们可以进行用户名密码认证,认证的用户名和密码可以在命令行指定    
`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -a "user1:pass1" -a "user2:pass2"`  
多个用户,重复-a参数即可.  
也可以放在文件中,格式是一行一个"用户名:密码",然后用-F指定.  
`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -F auth-file.txt`  

如果上级有认证,下级可以通过-A参数设置认证信息,比如:  
上级:`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -a "user1:pass1" -a "user2:pass2"`  
下级:`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -A "user1:pass1" -t tcp -p ":33080" `  

另外,sps代理,本地认证集成了外部HTTP API认证,我们可以通过--auth-url参数指定一个http url接口地址,    
然后有用户连接的时候,proxy会GET方式请求这url,带上下面四个参数,如果返回HTTP状态码204,代表认证成功  
其它情况认为认证失败.  
比如:  
`./proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" --auth-url "http://test.com/auth.php"`  
用户连接的时候,proxy会GET方式请求这url("http://test.com/auth.php"),  
带上user,pass,ip,target四个参数:  
http://test.com/auth.php?user={USER}&pass={PASS}&ip={IP}&target={TARGET}  
user:用户名   
pass:密码   
ip:用户的IP,比如:192.168.1.200   
target:如果客户端是http(s)代理请求,这里代表的是请求的完整url,其它情况为空.  

如果没有-a或-F或--auth-url参数,就是关闭本地认证.  
如果没有-A参数,连接上级不使用认证.  


#### **6.8 自定义加密**  
proxy的sps代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据,除此之外还支持在tls和kcp之后进行  
自定义加密,也就是说自定义加密和tls|kcp是可以联合使用的,内部采用AES256加密,使用的时候只需要自己定义  
一个密码即可,加密分为两个部分,一部分是本地(-z)是否加密解密,一部分是与上级(-Z)传输是否加密解密.      

自定义加密要求两端都是proxy才可以.  

下面分别用二级,三级为例:  

假设已经存在一个http(s)代理:`6.6.6.6:6666`  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy sps -S http -T tcp -P 6.6.6.6:6666 -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy sps -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy sps -S http -T tcp -P 6.6.6.6:6666 -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy sps -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888` 
本地三级执行:  
`proxy sps -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站.  

#### **6.9 压缩传输**  
proxy的sps代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据,在自定义加密之前还可以  
对数据进行压缩,也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的,压缩分为两个部分,   
一部分是本地(-m)是否压缩传输,一部分是与上级(-M)传输是否压缩.    

压缩要求两端都是proxy才可以,压缩也在一定程度上保护了(加密)数据.  

下面分别用二级,三级为例:  

**二级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy sps -t tcp -m -p :7777`  
本地二级执行:  
`proxy sps -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  


**三级实例**  
一级vps(ip:2.2.2.2)上执行:  
`proxy sps -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy sps -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888` 
本地三级执行:  
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站.  

#### **6.10 禁用协议**  
SPS默认情况下一个端口支持http(s)和socks5两种代理协议,我们可以通过参数禁用某个协议  
比如:  
1.禁用HTTP(S)代理功能只保留SOCKS5代理功能,参数:`--disable-http`.   
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080 --disable-http`   

1.禁用SOCKS5代理功能只保留HTTP(S)代理功能,参数:`--disable-socks`.   
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080 --disable-http`   

#### **6.11 限速**  

假设存在SOCKS5上级:

`proxy socks -p 2.2.2.2:33080 -z password -t tcp`

sps下级,限速100K

`proxy sps -S socks -P 2.2.2.2:33080 -T tcp -Z password -l 100K -t tcp -p :33080`

通过`-l`参数即可指定,比如:100K 1.5M . 0意味着无限制.

#### **6.12 指定出口IP**  

`--bind-listen`参数，就可以开启客户端用入口IP连接过来的,就用入口IP作为出口IP访问目标网站的功能。如果入口IP是内网IP，出口IP不会使用入口IP。

`proxy sps -S socks -P 2.2.2.2:33080 -T tcp -Z password -l 100K -t tcp --bind-listen -p :33080`

#### **6.13 证书参数使用base64数据**  

默认情况下-C,-K参数是crt证书和key文件的路径,

如果是base64://开头,那么就认为后面的数据是base64编码的,会解码后使用.

#### **6.14 查看帮助** 
`./proxy help sps`  

### **7.KCP配置**   

#### **7.1 配置介绍**   
proxy的很多功能都支持kcp协议，凡是使用了kcp协议的功能都支持这里介绍的配置参数。  
所以这里统一对KCP配置参数进行介绍。  

#### **7.2 详细配置**   
所有的KCP配置参数共有17个，你可以都不用设置，他们都有默认值，如果为了或者最好的效果，  
就需要自己根据自己根据网络情况对参数进行配置。由于kcp配置很复杂需要一定的网络基础知识，  
如果想获得kcp参数更详细的配置和解说，请自行搜索。每个参数的命令行名称以及默认值和简单的功能说明如下：  
```
--kcp-key="secrect"        pre-shared secret between client and server
--kcp-method="aes"         encrypt/decrypt method, can be: aes, aes-128, aes-192, salsa20, blowfish, 
                           twofish, cast5, 3des, tea, xtea, xor, sm4, none
--kcp-mode="fast"       profiles: fast3, fast2, fast, normal, manual
--kcp-mtu=1350             set maximum transmission unit for UDP packets
--kcp-sndwnd=1024          set send window size(num of packets)
--kcp-rcvwnd=1024          set receive window size(num of packets)
--kcp-ds=10                set reed-solomon erasure coding - datashard
--kcp-ps=3                 set reed-solomon erasure coding - parityshard
--kcp-dscp=0               set DSCP(6bit)
--kcp-nocomp               disable compression
--kcp-acknodelay           be carefull! flush ack immediately when a packet is received
--kcp-nodelay=0            be carefull!
--kcp-interval=50          be carefull!
--kcp-resend=0             be carefull!
--kcp-nc=0                 be carefull! no congestion
--kcp-sockbuf=4194304      be carefull!
--kcp-keepalive=10         be carefull!
```
提示：  
参数：--kcp-mode中的四种fast3, fast2, fast, normal模式，  
相当于设置了下面四个参数：  
normal：`--nodelay=0 --interval=40 --resend=2 --nc=1`  
fast ：`--nodelay=0 --interval=30 --resend=2 --nc=1`  
fast2：`--nodelay=1 --interval=20 --resend=2 --nc=1`  
fast3：`--nodelay=1 --interval=10 --resend=2 --nc=1`  

### **8.DNS防污染服务器** 

#### **8.1 介绍** 
众所周知DNS是UDP端口53提供的服务，但是随着网络的发展一些知名DNS服务器也支持TCP方式dns查询，比如谷歌的8.8.8.8，proxy的DNS防污染服务器原理就是在本地启动一个proxy的DNS代理服务器，它用TCP的方式通过上级代理进行dns查询。如果它和上级代理通讯采用加密的方式，那么就可以进行安全无污染的DNS解析。

#### **8.2 使用示例** 

***8.2.1 普通HTTP(S)上级代理***   
假设有一个上级代理：2.2.2.2:33080  
本地执行：  
`proxy dns -S http -T tcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了DNS解析功能。  

***8.2.2 普通SOCKS5上级代理***   
假设有一个上级代理：2.2.2.2:33080  
本地执行：  
`proxy dns -S socks -T tcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了DNS解析功能。 

***8.2.3 TLS加密的HTTP(S)上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy http -t tls -C proxy.crt -K proxy.key -p :33080`
本地执行：  
`proxy dns -S http -T tls -P 2.2.2.2:33080  -C proxy.crt -K proxy.key -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。 

***8.2.4 TLS加密的SOCKS5上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy socks -t tls -C proxy.crt -K proxy.key -p :33080`
本地执行：  
`proxy dns -S socks -T tls -P 2.2.2.2:33080  -C proxy.crt -K proxy.key -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。  

***8.2.5 KCP加密的HTTP(S)上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy http -t kcp -p :33080`
本地执行：  
`proxy dns -S http -T kcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。 

***8.2.6 KCP加密的SOCKS5上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy socks -t kcp -p :33080`
本地执行：  
`proxy dns -S socks -T kcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。 

***8.2.7 自定义加密的HTTP(S)上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy http -t tcp -p :33080 -z password`
本地执行：  
`proxy dns -S http -T tcp -Z password -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。 

***8.2.8 自定义加密的SOCKS5上级代理***   
假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：
`proxy socks -t kcp -p :33080 -z password`
本地执行：  
`proxy dns -S socks -T tcp -Z password -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

### TODO  
- http,socks代理多个上级负载均衡?
- http(s)代理增加pac支持?
- 欢迎加群反馈...

### 如何使用源码?   
建议go1.10.1.       
`go get github.com/snail007/goproxy`   
cd进入你的go src目录  
cd进入`github.com/snail007/goproxy`即可.    
编译直接:`go build -o proxy`        
运行: `go run *.go`       
utils是工具包,service是具体的每个服务类. 

### License  
Proxy is licensed under GPLv3 license.  
### Contact  
QQ交流群:189618940  
  
  
### Donation  
如果proxy帮助你解决了很多问题,你可以通过下面的捐赠更好的支持proxy.  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/alipay.jpg?raw=true" width="200"/>  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/wxpay.jpg?raw=true" width="200"/>  
