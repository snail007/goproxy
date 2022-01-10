## 稳定优质IP代理商推荐
<a href="https://brightdata.grsm.io/9gnzdzgpvtmo">
<img width="auto" height="100" src="https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/doc/images/brightdata.png"/>
</a>

## GOPROXY简介

<div align="center">
<img src="https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/doc/images/logo.jpg" width="500" height="auto"/>

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https://img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total.svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style=plastic)](https://github.com/snail007/goproxy/releases)

---

GoProxy是一款轻量级、功能强大、高性能的http代理、https代理、socks5代理、内网穿透代理服务器、ss代理、游戏盾、游戏代理，支持API代理认证。websocke代理、tcp代理、udp代理、socket代理、高防服务器。支持正向代理、反向代理、透明代理、TCP内网穿透、UDP内网穿透、HTTP内网穿透、HTTPS内网穿透、https代理负载均衡、http代理负载均衡、socks5代理负载均衡、socket代理负载均衡、ss代理负载均衡、TCP/UDP端口映射、SSH中转、TLS加密传输、协议转换、防污染DNS代理，限速，限连接数。官方QQ交流群: 608062193。

</div>

---

### [点击我观看视频教程](https://space.bilibili.com/472844633)
- [下载地址](https://github.com/snail007/goproxy/releases)
- [参考手册](https://snail007.host900.com/goproxy/manual/zh/)
- [桌面版](https://github.com/snail007/proxy_admin_free/blob/master/README_ZH.md)
- [安卓全局代理版](https://github.com/snail007/goproxy-ss-plugin-android) 
- [安卓全能代理版](https://github.com/snail007/goproxy-android) 
- [安卓内网穿透客户端](https://github.com/snail007/lanass) 
- [SDK](https://github.com/snail007/goproxy-sdk)
- [GORPOXY实战教程](https://snail007.host900.com/goproxy/)
- [免费版VS商业版(安装、激活)](https://snail007.host900.com/goproxy/page/free_vs_commercial/)

## 国内下载

请在github的下载链接前面加上: `https://mirrors.host900.com/` 。

比如`v10.4`的github下载链接是：

`https://github.com/snail007/goproxy/releases/download/v10.4/proxy-linux-amd64.tar.gz`

那么国内下载地址就是：

`https://mirrors.host900.com/https://github.com/snail007/goproxy/releases/download/v10.4/proxy-linux-amd64.tar.gz`

此地址也适用于wget，curl直接命令行下载。

## ProxyAdmin介绍预览（这不是goproxy，是控制面板友情链接；安装使用goproxy请往下看，谢谢！）
`ProxyAdmin` 是强大的代理服务工具 snail007/goproxy 的控制面板，运行了它，一秒让你的服务器变为强大的代理服务器，友好的交互界面，小白也能轻松上手，让你用起来得心应手，心情舒畅。

![](https://mirrors.host900.com/https://github.com/snail007/proxy_admin_free/blob/master/res/images/socks5_cn.gif)

### goproxy能干什么？
- 链式代理，程序本身可以作为一级代理，如果设置了上级代理那么可以作为二级代理，乃至N级代理。  
- 通讯加密，如果程序不是一级代理，而且上级代理也是本程序，那么可以加密和上级代理之间的通讯，采用底层tls高强度加密，安全无特征。  
- 智能HTTP代理，HTTPS代理，SOCKS5代理，会自动判断访问的网站是否屏蔽，如果被屏蔽那么就会使用上级代理(前提是配置了上级代理)访问网站;如果访问的网站没有被屏蔽，为了加速访问，代理会直接访问网站，不使用上级代理。  
- 域名黑白名单，更加自由的控制网站的访问方式。  
- 跨平台性，无论你是windows，linux，还是mac，甚至是树莓派，都可以很好的运行proxy。  
- 多协议支持，支持HTTP(S)，TCP，UDP，Websocket，SOCKS5代理。  
- TCP/UDP端口转发。 
- 游戏盾，游戏代理，高防服务器。 
- 内网穿透，P2P传输，协议支持TCP和UDP，针对HTTP的优化穿透。  
- SSH中转，HTTP(S)，SOCKS5代理支持SSH中转，上级Linux服务器不需要任何服务端，本地一个proxy即可开心上网。  
- [KCP](https://github.com/xtaci/kcp-go)协议支持，HTTP(S)，SOCKS5代理支持KCP协议传输数据，降低延迟，提升浏览体验。  
- 动态选择上级代理，通过外部API，HTTP(S)，SOCKS5，SPS代理可以实现基于用户或者IP的限速，连接数限制，动态获取上级。
- 灵活的上级分配，HTTP(S)，SOCKS5，SPS代理可以通过配置文件实现基于用户或者IP的限速，连接数限制，指定上级。
- 反向代理，支持直接把域名解析到proxy监听的ip，然后proxy就会帮你代理访问需要访问的HTTP(S)网站。  
- 透明HTTP(S)代理，配合iptables，在网关直接把出去的80，443方向的流量转发到proxy，就能实现无感知的智能路由器代理。  
- 协议转换，可以把已经存在的HTTP(S)或SOCKS5或SS代理转换为一个端口同时支持HTTP(S)和SOCKS5和SS代理，转换后的SOCKS5和SS代理如果上级是SOCKS5代理，那么支持UDP功能，同时支持强大的级联认证功能。
- 自定义底层加密传输，http(s)\sps\socks代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据，除此之外还支持在tls和kcp之后进行自定义加密，也就是说自定义加密和tls|kcp是可以联合使用的，内部采用AES256加密，使用的时候只需要自己定义一个密码即可。
- 底层压缩高效传输，http(s)\sps\socks代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据，在加密之后还可以对数据进行压缩，也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的。
- 安全的DNS代理，可以通过本地的proxy提供的DNS代理服务器与上级代理加密通讯实现安全防污染的DNS查询。
- 负载均衡，高可用，HTTP(S)\SOCKS5\SPS代理支持上级负载均衡和高可用，多个上级重复-P参数即可。  
- 指定出口IP，HTTP(S)\SOCKS5\SPS代理支持客户端用入口IP连接过来的，就用入口IP作为出口IP访问目标网站的功能。如果入口IP是内网IP，出口IP不会使用入口IP
- 支持限速，HTTP(S)\SOCKS5\SPS\TCP代理支持限速。  
- 支持限连接数，HTTP(S)\SOCKS5\SPS\TCP代理支持限连接数。  
- SOCKS5代理支持级联认证。  
- 证书参数使用base64数据，默认情况下-C，-K参数是crt证书和key文件的路径，如果是base64://开头，那么就认为后面的数据是base64编码的，会解码后使用。  
- 支持客户端IP黑白名单，更加安全的控制客户端对代理服务的访问，如果黑白名单同时设置，那么只有白名单生效。socks/http(s)/sps/tcp/udp/dns/内网穿透bridge/内网穿透tbridge，都支持客户端IP黑白名单。 
- 端口范围批量监听，HTTP(S)\SOCKS5\SPS\TCP代理支持指定端口范围监听，避免启动过多进程，提高性能。

### 为什么需要它？

- 当由于某某原因，我们不能访问我们在其它地方的服务，我们可以通过多个相连的proxy节点建立起一个安全的隧道访问我们的服务。  
- 微信接口本地开发，方便调试。  
- 远程访问内网机器。  
- 和小伙伴一起玩局域网游戏。  
- 以前只能在局域网玩的，现在可以在任何地方玩。  
- 替代圣剑内网通，显IP内网通，花生壳之类的工具。 
- 有大量IP资源，想变现，对外提供IP代理服务。 
- 有大量拨号VPS，想对外提供IP代理服务。
- 公司安全要求，审计员工对互联网的访问。
- 想要一个高性能稳定的，认证功能齐全的代理服务。
- 想一个固定入口，实现动态IP出口。
- ..。  

 
本页手册适用于最新版goproxy，其他版本可能有的地方不再适用，请自己根据命令帮助使用。  
 

### 加入组织

[点击加入 Telegram 交流群](https://t.me/snail007_goproxy)  

## 下载安装 goproxy

### 快速安装 goproxy

如果你的VPS是linux64位的系统，那么只需要执行下面一句，就可以完成自动安装和配置.

提示:所有操作需要root权限。 

免费版执行这个：  

```shell  
curl -L https://mirrors.host900.com/https://github.com/snail007/goproxy/blob/master/install_auto.sh | bash  
```  

商业版执行这个：  

```shell  
curl -L https://mirrors.host900.com/https://github.com/snail007/goproxy/blob/master/install_auto_commercial.sh | bash  
```  

安装完成，配置目录是/etc/proxy，更详细的使用方法请参考上面的手册目录，进一步了解你想要使用的功能。  
如果安装失败或者你的vps不是linux64位系统，请按照下面的半自动步骤安装:  
  
### 手动安装 goproxy

1.下载goproxy

根据你的平台和CPU类型选择,下载地址: https://github.com/snail007/goproxy/releases ,  

这里以 `proxy-linux-amd64.tar.gz` `v10.4` 为例,具体使用的时候,请根据你的平台和CPU类型选择具体文件名称.   

免费版执行这个：  

```shell  
cd /root/proxy/  
wget https://mirrors.host900.com/https://github.com/snail007/goproxy/releases/download/v10.4/proxy-linux-amd64.tar.gz  
```  

商业版执行这个：  

```shell  
cd /root/proxy/  
wget https://mirrors.host900.com/https://github.com/snail007/goproxy/releases/download/v10.4/proxy-linux-amd64_commercial.tar.gz  
```  

2.下载自动安装脚本

免费版执行这个：  

```shell  
cd /root/proxy/  
wget https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/install.sh  
chmod +x install.sh  
./install.sh  
```  

商业版执行这个：  

```shell  
cd /root/proxy/  
wget https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/install_commercial.sh  
chmod +x install_commercial.sh  
./install_commercial.sh  
```  

## 升级更新

### Linux
用`root`打开一个终端，如果proxy不再系统PATH里面，需要cd进入proxy目录执行`./proxy`。
下面假设proxy在/usr/bin/proxy，执行用的是`proxy`。

```shell
proxy update
```

已经安装了最新的版本，默认不会更新，如果想强制更新加上 -f 参数即可。

```shell
proxy update -f
```

### Windows
用`管理员`权限打开命令提示符窗口，如果proxy不再系统PATH里面，需要cd进入proxy目录执行。
这里假设proxy在c:\gp\proxy,根据你的情况调整命令。

```bat
c:\
cd gp
proxy update
```

已经安装了最新的版本，默认不会更新，如果想强制更新加上 -f 参数即可。

```shell
c:\
cd gp
proxy update -f
```

## TODO  
- http，socks代理多个上级负载均衡?
- http(s)代理增加pac支持?
- 欢迎加群反馈..。  

## License  
Proxy is licensed under GPLv3 license。  

## Contact  
官方QQ交流群: 608062193  

## Donation  
如果proxy帮助你解决了很多问题，你可以通过下面的捐赠更好的支持proxy。  
<img src="https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/doc/images/alipay.jpg" width="200"  height="auto"/>  
<img src="https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/doc/images/wxpay.jpg" width="200"  height="auto"/>  

### 源代码申明

本项目作者发现大量的开发者基于本项目进行二次开发或使用大量本项目核心代码而不遵循GPLv3协议，这严重违背了本项目使用GPLv3开源协议的初衷，鉴于这种情况，本项目采取源代码延迟发布策略，在一定程度上遏制这些不尊重开源，不尊重他人劳动成果的行为。  
本项目会持续更新迭代，持续发布全平台的二进制程序，给大家提供强大便捷的代理工具。  
如果你有定制，商业需求请发邮件至`arraykeys@gmail.com`

## goproxy使用手册


### 稳定优质IP代理商推荐

<a href="https://brightdata.grsm.io/9gnzdzgpvtmo">
<img width="auto" height="100" src="https://mirrors.host900.com/https://raw.githubusercontent.com/snail007/goproxy/master/doc/images/brightdata.png"/>
</a>

## 如何安装

### 1. Linux安装

[点击查看Linux安装教程](https://github.com/snail007/goproxy/blob/master/README_ZH.md#%E4%B8%8B%E8%BD%BD%E5%AE%89%E8%A3%85-goproxy)

### 2. 苹果Mac系统安装

[点击查看苹果Mac系统安装教程](https://github.com/snail007/proxy_admin_free/blob/master/README_ZH.md#%E8%A7%86%E9%A2%91%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

### 3. Windows安装

为了方便操作,推荐Windows用户使用proxy-admin面板,[点击查看Windows安装教程](https://github.com/snail007/proxy_admin_free/blob/master/README_ZH.md#%E8%A7%86%E9%A2%91%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

当然你也可以使用命令行goproxy[点击查看手动安装](https://github.com/snail007/goproxy/blob/master/README_ZH.md#%E6%89%8B%E5%8A%A8%E5%AE%89%E8%A3%85-goproxy)

### 4. 其它平台安装

[点击查看其它安装教程](https://github.com/snail007/goproxy/blob/master/README_ZH.md#%E6%89%8B%E5%8A%A8%E5%AE%89%E8%A3%85-goproxy)

## 免费版、商业版说明

本手册描述功能,`proxyadmin商业版`和`goproxy商业版`全部包含;认证等高级功能参数免费版和VIP控制面板不包含;

如果您用`goproxy免费版`执行某些命令或者`proxyadmin 免费版`、`proxyadmin VIP版`服务启动失败的时候遇到,类似如下xxx参数不存在的提示,说明这个参数是商业版的功能,您需要下载并购买商业版授权才能使用.  
` err : unknown short flag '-a'`

以下软件都是相互独立使用，没有依赖关系，需要购买的也是分别单独购买使用。

软件名称 | 免费版功能 | 商业版功能 | 购买使用 | 特点 | 传送门
:----------- | :---: | :---: | :---: | :---: | :---:
proxyadmin 免费版 |   √ |   x |   x| Web界面操作，支持Linux，Windows，macOS，免费使用，服务数量有限制，适合个人，小白和白嫖党 | [下载安装](https://github.com/snail007/proxy_admin_free)
proxyadmin VIP版 |   √ |   x |   √| Web界面操作，支持更多平台，功能和免费版一样，服务数量无限制，适合个体户或者小集体 | [下载安装](https://github.com/snail007/proxy-admin-vip)
proxyadmin 商业版 |   √ |   √ |   √|Web界面操作，支持更多平台，无限制，适合集体或公司对外提供服务 | [下载安装](https://github.com/snail007/proxy-admin-commercial)
proxyadmin 集群版|   √ |   √ |   √|Web界面操作，支持Linux，Windows，macOS，无限制，配合`goproxy 商业版`,可以实现以组为单位，管理海量机器上的proxy服务，适合有大量机器的集体或公司对外提供服务 | [下载安装](https://github.com/snail007/proxy-admin-cluster)
goproxy 免费版 |   √ |   x |   x|命令行操作，全平台支持，免费使用，稳定且灵活，适合一切熟悉命令行的大佬或集体或公司白嫖，自用或者大规模部署对外服务 | [下载安装](https://github.com/snail007/goproxy)
goproxy 商业版|   √ |   √ |   √|命令行操作，全平台支持，稳定且灵活，适合集体或公司大规模部署对外提供服务 | [下载安装](https://github.com/snail007/goproxy)

关于免费版功能和商业版功能对比请看这里,[`免费版`和`商业版`功能对比](https://snail007.host900.com/goproxy/page/free_vs_commercial/).

[商业版激活绑定教程](https://snail007.host900.com/goproxy/page/free_vs_commercial/)

提示：

**免费和付费软件均没有额外技术支持，不按着手册操作的一切问题请自行解决。付费软件免费享有`手册功能`使用指导和`首次安装配置`指导。**

## FAQ

[别点我](https://snail007.host900.com/goproxy/page/faq/goproxy%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98%E8%A7%A3%E7%AD%94/)

## 首次使用必看,谢谢!!！

### 1. 环境

该手册教程，默认系统是linux，程序是proxy；所有操作需要root权限；

如果你的是windows，请使用windows版本的proxy.exe即可。

### 2. 使用配置文件

接下来的教程都是通过命令行参数介绍使用方法，也可以通过读取配置文件获取参数。

具体格式是通过@符号指定配置文件，例如:proxy @configfile.txt

configfile.txt里面的格式是，第一行是子命令名称，第二行开始一行一个参数，

格式:`参数 参数值`，没有参数值的直接写参数，比如:--nolog

比如configfile.txt内容如下:

```shell  
http  
-t tcp  
-p :33080  
--forever  
```  

### 3. 调试输出

默认情况下，日志输出的信息不包含文件行数，某些情况下为了排除程序问题，快速定位问题，

可以使用--debug参数，输出代码行数和毫秒时间。

### 4. 使用日志文件

默认情况下，日志是直接在控制台显示出来的，如果要保存到文件，可以使用--log参数，

比如: --log proxy.log，日志就会输出到proxy.log方便排除问题。

### 5. 生成加密通讯需要的证书文件

http(s)代理、tcp代理、udp代理、socks5代理、内网穿透等功能和上级通讯的时候，为了安全我们采用TLS加密通讯，当然可以选择不加密通信通讯，本教程所有和上级通讯都采用加密，需要证书文件。

***所有端必须使用相同的proxy.crt和proxy.key***

1.通过下面的命令生成自签名的证书和key文件。  
`proxy keygen -C proxy`  
会在当前程序目录下面生成证书文件proxy.crt和key文件proxy.key。

2.通过下面的命令生，使用自签名证书proxy.crt和key文件proxy.key签发新证书:goproxy.crt和goproxy.key。  
`proxy keygen -s -C proxy -c goproxy`  
会在当前程序目录下面生成证书文件goproxy.crt和key文件goproxy.key。

3.默认情况下证书的里面的域名是随机的，可以使用`-n test.com`参数指定。

4.更多用法:`proxy keygen --help`。

### 6. 后台运行

默认执行proxy之后，如果要保持proxy运行，不能关闭命令行。

如果想在后台运行proxy，命令行可以关闭，只需要在命令最后加上--daemon参数即可。

比如:

`proxy http -t tcp -p "0.0.0.0:38080" --daemon`

### 7. 守护运行

守护运行参数--forever，比如: `proxy http --forever` ，

proxy会fork子进程，然后监控子进程，如果子进程异常退出，5秒后重启子进程。

该参数配合后台运行参数--daemon和日志参数--log，可以保障proxy一直在后台执行不会因为意外退出，

而且可以通过日志文件看到proxy的输出日志内容。

比如: `proxy http -p ":9090" --forever --log proxy.log --daemon`

### 8. 安全建议

当VPS在nat设备后面，vps上网卡IP都是内网IP，这个时候可以通过-g参数添加vps的外网ip防止死循环。

假设你的vps外网ip是23.23.23.23，下面命令通过-g参数设置23.23.23.23

`proxy http -g "23.23.23.23"`

### 9. 负载均衡和高可用

HTTP(S)\SOCKS5\SPS\TCP代理支持上级负载均衡和高可用，多个上级重复-P参数即可。

负载均衡策略支持5种，可以通过`--lb-method`参数指定:

roundrobin 轮流使用

leastconn 使用最小连接数的

leasttime 使用连接时间最小的

hash 使用根据客户端地址计算出一个固定上级

weight 根据每个上级的权重和连接数情况，选择出一个上级

提示:

1.负载均衡检查时间间隔可以通过`--lb-retrytime`设置，单位毫秒

2.负载均衡连接超时时间可以通过`--lb-timeout`设置，单位毫秒

3.如果负载均衡策略是权重(weight)，-P格式为:2.2.2.2:3880?w=1，1就是权重，大于0的整数。

4.如果负载均衡策略是hash，默认是根据客户端地址选择上级，可以通过开关`--lb-hashtarget`使用访问的目标地址选择上级。

5.TCP代理没有参数`--lb-hashtarget`.

6.默认是负载均衡+高可用模式,如果使用了参数`--lb-onlyha`就只使用高可用模式,依据负载均衡策略选择一个节点,之后就一直使用这个节点,直到这个节点不再存活,那么会依据负载均衡策略再选择一个节点使用,以此循环.

7.如果检查节点全部不再存活,那么每次连接都会随机选取一个节点使用.

### 10. 代理跳板跳转

http(s)代理，SPS代理，内网穿透，tcp代理都支持通过中间第三方代理连接上级，

参数是:--jumper，所有格式如下:

```text  
http://username:password@host:port  
http://host:port  
https://username:password@host:port  
https://host:port  
socks5://username:password@host:port  
socks5://host:port  
socks5s://username:password@host:port  
socks5s://host:port  
ss://method:password@host:port  
```  

http，socks5代表的是普通的http和socks5代理。

https，socks5s代表的是通过tls保护的http和socks5代理，

也就是http代理 over TLS ， socks over TLS。

### 11. 域名黑白名单

socks/http(s)/sps代理都支持域名黑白名单。

用--stop参数指定一个域名黑名单列表文件，那么当用户连接文件里面这些域名的时候连接就会被断开。

用--only参数指定一个域名白名单列表文件，那么当用户连接文件里面这些域名之外的域名的时候连接就会被断开。

如果同时设置了--stop和--only，那么只有--only会起作用。

黑白域名名单文件内容格式如下:

```text  
**.baidu.com  
*.taobao.com  
a.com  
192.168.1.1  
192.168.*.*  
?.qq.com  
```  

说明:

1.一行一个域名，域名写法支持通配符`*`和`?`，`*`代表任意个字符，`?`代表一个任意字符，

2.`**.baidu.com` 匹配无论是多少级所有后缀是`.baidu.com`的域名。

3.`*.taobao.com` 匹配后缀是`.taobao.com`的三级域名。

4.还可以直接是IP地址。

5.`#`开头的为注释。

### 12. 客户端IP黑白名单

socks/http(s)/sps/tcp/udp/dns/内网穿透bridge/内网穿透tbridge，都支持客户端IP黑白名单。

用--ip-deny参数指定一个客户端IP黑名单列表文件，那么当用户的IP在这个文件里面的时候连接就会被断开。

用--ip-allow参数指定一个客户端IP白名单列表文件，那么当用户的IP不在这个文件里面的时候连接就会被断开。

如果同时设置了--ip-deny和--ip-allow，那么只有--ip-allow会起作用。

客户端IP黑白名单文件内容格式如下:

```text  
192.168.1.1  
192.168.*.*  
192.168.1?.*  
```  

说明:

1.一行一个域名，域名写法支持通配符`*`和`?`，`*`代表任意个字符，`?`代表一个任意字符。

2.`#`开头的为注释。

### 13. 协议加载文件

proxy的各种代理功能里面很多地方都有参数设置一个文件，比如：--blocked 指定一个直接走上级的域名列表文件，参数值是文件的路径，

如果参数支持协议加载文件，那么文件路径不仅可以是文件路径，还可以是：

a.“base64://”开头的base64编码的上面说明的文件内容，比如：base64://ajfpoajsdfa=

b.”str://“开头的英文逗号分割的多个，比如：str://xxx,yyy

proxy的blocked，direct，stop，only，hosts，resolve.rules，rewriter.rules，ip.allow，ip.deny 文件支持协议加载。

### 14. 客户端并发连接数

socks5\sps\http代理,控制客户端并发连接数参数是:`--max-conns-rate`,控制每秒客户端的最大连接数,默认20, 0为不限制.

### 15 监听多个端口

`tcp/http/socks/sps`支持同时监听多个端口以及范围端口。 一般情况下监听一个端口就可以，不过如果需要同时监听多个两个端口，或者范围端口，那么-p参数是支持的，  
格式是：`-p 0.0.0.0:80,0.0.0.0:443,0.0.0.0:8000-9000,:5000-6000`，多个绑定用逗号分隔即可。

## 1.HTTP代理

### 1.1.普通一级HTTP代理

![1.1](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-1.png)

`proxy http -t tcp -p "0.0.0.0:38080"`

-p参数支持的写法：

```text
  -p ":8081"  监听8081
  -p ":8081,:8082"  监听8081和8082
  -p ":8081,:8082,:9000-9999" 监听8081和8082以及9000,9001至9999，共1002个端口
```

### 1.2.普通二级HTTP代理

![1.2](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-2.png)

使用本地端口8090，假设上级HTTP代理是`22.22.22.22:8080`

`proxy http -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080" `

我们还可以指定网站域名的黑白名单文件，一行一个域名，匹配规则是最右匹配，比如:baidu.com，匹配的是*.*.baidu.com，黑名单的域名直接走上级代理，白名单的域名不走上级代理。

`proxy http -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt`

### 1.3.HTTP二级代理(加密)

> 注意: 后面二级代理使用的`proxy.crt`和`proxy.key`应与一级代理一致

![1.3](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-tls-2.png)  
一级HTTP代理(VPS，IP:22.22.22.22)  
`proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key`

二级HTTP代理(本地Linux)  
`proxy http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080。

二级HTTP代理(本地windows)  
`proxy.exe http -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为http模式，地址为：127.0.0.1，端口为：8080，程序即可通过加密通道通过vps上网。

### 1.4.HTTP三级代理(加密)

![1.3](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-tls-3.png)  
一级HTTP代理VPS_01，IP:22.22.22.22  
`proxy http -t tls -p ":38080" -C proxy.crt -K proxy.key`  
二级HTTP代理VPS_02，IP:33.33.33.33  
`proxy http -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级HTTP代理(本地)  
`proxy http -t tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地的8080端口就是访问一级HTTP代理上面的代理端口38080。

### 1.5.Basic认证，API认证

请参考`9.API认证` 和 `10.本地认证`

### 1.6.HTTP代理流量强制走上级HTTP代理

默认情况下，proxy会智能判断一个网站域名是否无法访问，如果无法访问才走上级HTTP代理.通过--always可以使全部HTTP代理流量强制走上级HTTP代理。  
`proxy http --always -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`

### 1.7.HTTP(S)通过SSH中转

![1.7](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-ssh-1.png)  
说明:ssh中转的原理是利用了ssh的转发功能，就是你连接上ssh之后，可以通过ssh代理访问目标地址。  
假设有:vps

- IP是2.2.2.2， ssh端口是22， ssh用户名是:user， ssh用户密码是:demo
- 用户user的ssh私钥名称是user.key

#### *1.7.1 ssh用户名和密码的方式*

本地HTTP(S)代理28080端口，执行:  
`proxy http -T ssh -P "2.2.2.2:22" -u user -A demo -t tcp -p ":28080"`

#### *1.7.2 ssh用户名和密钥的方式*

本地HTTP(S)代理28080端口，执行:  
`proxy http -T ssh -P "2.2.2.2:22" -u user -S user.key -t tcp -p ":28080"`

### 1.8.KCP协议传输

![1.8](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/http-kcp.png)  
KCP协议需要--kcp-key参数设置一个密码用于加密解密数据

一级HTTP代理(VPS，IP:22.22.22.22)  
`proxy http -t kcp -p ":38080" --kcp-key mypassword`

二级HTTP代理(本地Linux)  
`proxy http -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" --kcp-key mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080，数据通过kcp协议传输，注意kcp走的是udp协议协议，所以防火墙需放开38080的udp协议。

### 1.9 HTTP(S)反向代理

![1.9](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/fxdl.png)  
proxy不仅支持在其他软件里面通过设置代理的方式，为其他软件提供代理服务，而且支持直接把请求的网站域名解析到proxy监听的ip上，然后proxy监听80和443端口，那么proxy就会自动为你代理访问需要访问的HTTP(S)网站。

使用方式:  
在"最后一级proxy代理"的机器上，因为proxy要伪装成所有网站，网站默认的端口HTTP是80，HTTPS是443，让proxy监听80和443端口即可.参数-p多个地址用逗号分割。  
`proxy http -t tcp -p :80,:443`

这个命令就在机器上启动了一个proxy代理，同时监听80和443端口，既可以当作普通的代理使用，也可以直接把需要代理的域名解析到这个机器的IP上。

如果有上级代理那么参照上面教程设置上级即可，使用方式完全一样。  
`proxy http -t tcp -p :80,:443 -T tls -P "2.2.2.2:33080" -C proxy.crt -K proxy.key`

注意:  
proxy所在的服务器的DNS解析结果不能受到自定义的解析影响，不然就死循环了，proxy代理最好指定`--dns 8.8.8.8`参数。

### 1.10 HTTP(S)透明代理

该模式需要具有一定的网络基础，相关概念不懂的请自行搜索解决。  
假设proxy现在在路由器上运行，启动命令如下:  
`proxy http -t tcp -p :33080 -T tls -P "2.2.2.2:33090" -C proxy.crt -K proxy.key`

然后添加iptables规则，下面是参考规则:

```shell  
#上级proxy服务端服务器IP地址:  
proxy_server_ip=2.2.2.2  

#路由器运行proxy监听的端口:  
proxy_local_port=33080  

#下面的就不用修改了  
#create a new chain named PROXY  
iptables -t nat -N PROXY  

# Ignore your PROXY server's addresses  
# It's very IMPORTANT， just be careful。  

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

### 1.11 自定义DNS

--dns-address和--dns-ttl参数，用于自己指定proxy访问域名的时候使用的dns（--dns-address）  
以及解析结果缓存时间（--dns-ttl）秒数，避免系统dns对proxy的干扰，另外缓存功能还能减少dns解析时间提高访问速度。  
比如：  
`proxy http -p ":33080" --dns-address "8.8.8.8:53" --dns-ttl 300`

### 1.12 自定义加密

proxy的http(s)代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据，除此之外还支持在tls和kcp之后进行自定义  
加密，也就是说自定义加密和tls|kcp是可以联合使用的，内部采用AES256加密，使用的时候只需要自己定义一个密码即可，  
加密分为两个部分，一部分是本地(-z)是否加密解密，一部分是与上级(-Z)传输是否加密解密。  
自定义加密要求两端都是proxy才可以，下面分别用二级，三级为例:

二级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy http -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

三级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy http -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888`  
本地三级执行:  
`proxy http -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

### 1.13 压缩传输

proxy的http(s)代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据，在自定义加密之前还可以对数据进行压缩，  
也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的，压缩分为两个部分，一部分是本地(-m)是否压缩传输，  
一部分是与上级(-M)传输是否压缩。  
压缩要求两端都是proxy才可以，压缩也在一定程度上保护了(加密)数据，下面分别用二级，三级为例:

二级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -m -p :7777`  
本地二级执行:  
`proxy http -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

三级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy http -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy http -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888`  
本地三级执行:  
`proxy http -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

### 1.14 负载均衡

HTTP(S)代理支持上级负载均衡，多个上级重复-P参数即可。

`proxy http --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080`

### 1.14.1 设置重试间隔和超时时间

`proxy http --lb-method=leastconn --lb-retrytime 300 --lb-timeout 300 -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -t tcp -p :33080`

### 1.14.2 设置权重

`proxy http --lb-method=weight -T tcp -P 1.1.1.1:33080?w=1 -P 2.1.1.1:33080?w=2 -P 3.1.1.1:33080?w=1 -t tcp -p :33080`

### 1.14.3 使用目标地址选择上级

`proxy http --lb-hashtarget --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -t tcp -p :33080`

### 1.15 限速

限速100K，通过`-l`参数即可指定，比如:100K 2000K 1M . 0意味着无限制。

`proxy http -t tcp -p 2.2.2.2:33080 -l 100K`

### 1.16 指定出口IP

`--bind-listen`参数，就可以开启客户端用`入口IP`连接过来的，就用`入口IP`作为`出口IP`访问目标网站的功能。如果绑定了不正确的IP会导致代理不能工作。

普通多IP用法：

`proxy http -t tcp -p :33080 --bind-listen`

上面的写法适合本机有大量IP，不想一个一个单独写，比如本机有IP`3.3.3.3`,当客户端使用`3.3.3.3`连接代理，那么`3.3.3.3`就会作为客户端的出口IP， 这里的`3.3.3.3`就是`入口IP`和`出口IP`
，你会发现`入口IP`和`出口IP`是一样的，`--bind-listen`这种模式人为不能干预`入口IP`和`出口IP`。

单独用法：

本机有IP：`5.5.5.5`，开一个代理端口`6666`，凡是通过`5.5.5.5:6666`端口进来的，都使用`5.5.5.5`作为出口IP。

命令如下：

`proxy http -t tcp -p 5.5.5.5:6666 --bind-listen`

#### 自由指定出口IP

上面的`--bind-listen`参数虽然可以指定出口IP，但是`入口IP`和`出口IP`不能人为干预。如果想入口IP和出口IP不一样， 可以使用`--bind-ip`参数，格式：`IP:端口`，比如： `1.1.1.1:8080`
，`[2000:0:0:0:0:0:0:1]:8080`。 多个绑定需求，可以重复`--bind-ip`参数即可。

比如本机有IP`5.5.5.5`,`6.6.6.6`，监听了`8888`和`7777`两个端口，命令如下：

`proxy http -t tcp -p :8888,:7777 --bind-ip 5.5.5.5:7777 --bind-ip 6.6.6.6:8888`

那么客户端访问`7777`端口，出口IP就是`5.5.5.5`，访问`8888`端口，出口IP就是`6.6.6.6`，如果同时设置了`--bind-ip`和`--bind-listen`,`--bind-ip`的优先级更高。

另外，`--bind-ip`参数的`IP`部分，支持指定`网卡名称`，`通配符`，还能指定多个，详细说明如下：

- 指定网卡名称，比如：`--bind-ip eth0:7777`，那么客户端访问`7777`端口，出口IP就是eth0网卡的IP。
- 网卡名称支持通配符，比如：`--bind-ip eth0.*:7777`，那么客户端访问`7777`端口，出口IP就是`eth0.`开头网卡的IP中随机选择的一个。
- IP支持通配符，比如：`--bind-ip 192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器所有IP中，匹配`192.168.?.*`的IP中随机选择的一个。
- 还可以是网卡名称和IP的多个组合，多个使用半角逗号分割，比如：`--bind-ip pppoe??,192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器网卡名称匹配`pppoe??`
  和机器所有IP中匹配`192.168.?.*`的IP中随机选择的一个。
- 通配符`*`代表0至任意多个字符，`?`代表1个字符。
- 如果网卡IP发生变化，也会实时生效。
- 可以通过`--bind-refresh`参数，指定刷新本地网卡信息的间隔，默认`5`，单位秒。

### 1.17 证书参数使用base64数据

默认情况下-C，-K参数是crt证书和key文件的路径，

如果是base64://开头，那么就认为后面的数据是base64编码的，会解码后使用。

### 1.18 智能模式

智能模式设置，可以是intelligent|direct|parent三者之一。  
默认是:parent。  
每个值的含义如下:  
`--intelligent=direct`，不在blocked里面的目标都直连。  
`--intelligent=parent`，不在direct里面的目标都走上级。  
`--intelligent=intelligent`，blocked和direct里面都没有的目标，智能判断是否使用上级访问目标。

### 1.19 查看帮助

`proxy help http`

## 2.TCP代理

### 2.1 普通一级TCP代理

![2.1](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/tcp-1.png)  
本地执行:  
`proxy tcp -p ":33080" -T tcp -P "192.168.22.33:22"`  
那么访问本地33080端口就是访问192.168.22.33的22端口。

`-p`参数支持的写法：

```text
  -p ":8081"  监听8081
  -p ":8081,:8082"  监听8081和8082
  -p ":8081,:8082,:9000-9999" 监听8081和8082以及9000,9001至9999，共1002个端口
```

如果本地监听端口数量大于1，那么将会连接与本地端口一致的对应上级端口，忽略`-P`里面的端口。

如果需要所有端口进来的连接，都连接到上级指定端口，可以加上参数`--lock-port`。

比如：

`proxy tcp -p ":33080-33085" -T tcp -P "192.168.22.33:0"`

那么`33080`端口进来的连接，将会连接192.168.22.33的`33080`端口，其它端口以此类推，本地和上级端口一致，此时参数`-P`里面的端口用`0`。

如果想无论是`33080`，`33081`等端口进来的连接都连接到192.168.22.33的`22`端口，可以加上参数`--lock-port`

`proxy tcp -p ":33080-33085" -T tcp -P "192.168.22.33:22" --lock-port`

### 2.2 普通二级TCP代理

![2.2](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/tcp-2.png)  
VPS(IP:22.22.22.33)执行:  
`proxy tcp -p ":33080" -T tcp -P "127.0.0.1:8080"`  
本地执行:  
`proxy tcp -p ":23080" -T tcp -P "22.22.22.33:33080"`  
那么访问本地23080端口就是访问22.22.22.33的8080端口。

### 2.3 普通三级TCP代理

![2.3](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/tcp-3.png)  
一级TCP代理VPS_01，IP:22.22.22.22  
`proxy tcp -p ":38080" -T tcp -P "66.66.66.66:8080"`  
二级TCP代理VPS_02，IP:33.33.33.33  
`proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`proxy tcp -p ":8080" -T tcp -P "33.33.33.33:28080"`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口。

### 2.4 加密二级TCP代理

![2.4](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/tcp-tls-2.png)  
VPS(IP:22.22.22.33)执行:  
`proxy tcp -t tls -p ":33080" -T tcp -P "127.0.0.1:8080" -C proxy.crt -K proxy.key`  
本地执行:  
`proxy tcp -p ":23080" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地23080端口就是通过加密TCP隧道访问22.22.22.33的8080端口。

### 2.5 加密三级TCP代理

![2.5](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/tcp-tls-3.png)  
一级TCP代理VPS_01，IP:22.22.22.22  
`proxy tcp -t tls -p ":38080" -T tcp -P "66.66.66.66:8080" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02，IP:33.33.33.33  
`proxy tcp -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`proxy tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地8080端口就是通过加密TCP隧道访问66.66.66.66的8080端口。

### 2.6 通过代理连接上级

有时候proxy所在的网络不能直接访问外网，需要通过一个https或者socks5代理才能上网，那么这个时候  
-J参数就可以帮助你让proxy的tcp端口映射的时候通过https或者socks5代理去连接上级-P，将外部端口映射到本地。  
-J参数格式如下:

https代理写法:  
代理需要认证，用户名:username 密码:password  
https://username:password@host:port  
代理不需要认证  
https://host:port

socks5代理写法:  
代理需要认证，用户名:username 密码:password  
socks5://username:password@host:port  
代理不需要认证  
socks5://host:port

host:代理的IP或者域名  
port:代理的端口

### 2.7 指定出口IP

当TCP代理当上级类型（参数：-T）是tcp当时候，支持指定`出口IP`。`--bind-listen`参数，就可以开启客户端用`入口IP`连接过来的，就用`入口IP`作为`出口IP`
访问目标网站的功能。如果绑定了不正确的IP会导致代理不能工作。

普通多IP用法：

`proxy tcp -t tcp -p :33080 --bind-listen -T tcp -P 2.2.2.2:3322`

上面的写法适合本机有大量IP，不想一个一个单独写，比如本机有IP`3.3.3.3`,当客户端使用`3.3.3.3`连接代理，那么`3.3.3.3`就会作为客户端的出口IP， 这里的`3.3.3.3`就是`入口IP`和`出口IP`
，你会发现`入口IP`和`出口IP`是一样的，`--bind-listen`这种模式人为不能干预`入口IP`和`出口IP`。

单独用法：

本机有IP：`5.5.5.5`，开一个代理端口`6666`，凡是通过`5.5.5.5:6666`端口进来的，都使用`5.5.5.5`作为出口IP。

命令如下：

`proxy tcp -t tcp -p 5.5.5.5:6666 --bind-listen -T tcp -P 2.2.2.2:3322`

#### 自由指定出口IP

上面的`--bind-listen`参数虽然可以指定出口IP，但是`入口IP`和`出口IP`不能人为干预。如果想入口IP和出口IP不一样， 可以使用`--bind-ip`参数，格式：`IP:端口`，比如： `1.1.1.1:8080`
，`[2000:0:0:0:0:0:0:1]:8080`。 多个绑定需求，可以重复`--bind-ip`参数即可。

比如本机有IP`5.5.5.5`,`6.6.6.6`，监听了`8888`和`7777`两个端口，命令如下：

`proxy tcp -t tcp -p :8888,:7777 --bind-ip 5.5.5.5:7777 --bind-ip 6.6.6.6:8888 -T tcp -P 2.2.2.2:3322`

那么客户端访问`7777`端口，出口IP就是`5.5.5.5`，访问`8888`端口，出口IP就是`6.6.6.6`，如果同时设置了`--bind-ip`和`--bind-listen`,`--bind-ip`的优先级更高。

另外，`--bind-ip`参数的`IP`部分，支持指定`网卡名称`，`通配符`，还能指定多个，详细说明如下：

- 指定网卡名称，比如：`--bind-ip eth0:7777`，那么客户端访问`7777`端口，出口IP就是eth0网卡的IP。
- 网卡名称支持通配符，比如：`--bind-ip eth0.*:7777`，那么客户端访问`7777`端口，出口IP就是`eth0.`开头网卡的IP中随机选择的一个。
- IP支持通配符，比如：`--bind-ip 192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器所有IP中，匹配`192.168.?.*`的IP中随机选择的一个。
- 还可以是网卡名称和IP的多个组合，多个使用半角逗号分割，比如：`--bind-ip pppoe??,192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器网卡名称匹配`pppoe??`
  和机器所有IP中匹配`192.168.?.*`的IP中随机选择的一个。
- 通配符`*`代表0至任意多个字符，`?`代表1个字符。
- 如果网卡IP发生变化，也会实时生效。
- 可以通过`--bind-refresh`参数，指定刷新本地网卡信息的间隔，默认`5`，单位秒。

### 2.8 限速，限制连接数

参数`--max-conns`可以限制每个端口的最大连接数。   
比如限制每个端口最多1000个连接数：   
`proxy tcp -p ":33080" -T tcp -P "192.168.22.33:22" --max-conns 1000`    
参数`--rate-limit`可以限制每个tcp连接的速率。   
比如限制每个tcp连接速率为100k/s：  
`proxy tcp -p ":33080" -T tcp -P "192.168.22.33:22" --rate-limit 100k`

### 2.9 压缩传输

`--c`控制本地和客户端之间是否压缩传输，默认false；`--C`控制本地和上级之间是否压缩传输，默认false。

示例：

VPS(IP:22.22.22.33)执行:  
`proxy tcp -t tcp --c -p ":33080" -T tcp -P "127.0.0.1:8080"`  
本地执行:  
`proxy tcp -t tcp -p ":23080" -T tcp -P "22.22.22.33:33080" --C`

### 2.10 查看帮助

`proxy help tcp`

## 3.UDP代理

### 3.1.普通一级UDP代理

![3.1](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/udp-1.png)  
本地执行:  
`proxy udp -p ":5353" -T udp -P "8.8.8.8:53"`  
那么访问本地UDP:5353端口就是访问8.8.8.8的UDP:53端口。

`-p`参数支持的写法：

```text
  -p ":8081"  监听8081
  -p ":8081,:8082"  监听8081和8082
  -p ":8081,:8082,:9000-9999" 监听8081和8082以及9000,9001至9999，共1002个端口
```

如果本地监听端口数量大于1，那么将会连接与本地端口一致的对应上级端口，忽略`-P`里面的端口。

如果需要所有端口进来的连接，都连接到上级指定端口，可以加上参数`--lock-port`。

比如：

`proxy udp -p ":33080-33085" -T udp -P "192.168.22.33:0"`

那么`33080`端口进来的连接，将会连接192.168.22.33的`33080`端口，其它端口以此类推，本地和上级端口一致，此时参数`-P`里面的端口用`0`。

如果想无论是`33080`，`33081`等端口进来的连接都连接到192.168.22.33的`2222`端口，可以加上参数`--lock-port`

`proxy udp -p ":33080-33085" -T udp -P "192.168.22.33:2222" --lock-port`

### 3.2.普通二级UDP代理

![3.2](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/udp-2.png)  
VPS(IP:22.22.22.33)执行:  
`proxy tcp -p ":33080" -T udp -P "8.8.8.8:53"`  
本地执行:  
`proxy udp -p ":5353" -T tcp -P "22.22.22.33:33080"`  
那么访问本地UDP:5353端口就是通过TCP隧道，通过VPS访问8.8.8.8的UDP:53端口。

### 3.3.普通三级UDP代理

![3.3](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/udp-3.png)  
一级TCP代理VPS_01，IP:22.22.22.22  
`proxy tcp -p ":38080" -T udp -P "8.8.8.8:53"`  
二级TCP代理VPS_02，IP:33.33.33.33  
`proxy tcp -p ":28080" -T tcp -P "22.22.22.22:38080"`  
三级TCP代理(本地)  
`proxy udp -p ":5353" -T tcp -P "33.33.33.33:28080"`  
那么访问本地5353端口就是通过TCP隧道，通过VPS访问8.8.8.8的53端口。

### 3.4.加密二级UDP代理

![3.4](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/udp-tls-2.png)  
VPS(IP:22.22.22.33)执行:  
`proxy tcp -t tls -p ":33080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
本地执行:  
`proxy udp -p ":5353" -T tls -P "22.22.22.33:33080" -C proxy.crt -K proxy.key`  
那么访问本地UDP:5353端口就是通过加密TCP隧道，通过VPS访问8.8.8.8的UDP:53端口。

### 3.5.加密三级UDP代理

![3.5](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/udp-tls-3.png)  
一级TCP代理VPS_01，IP:22.22.22.22  
`proxy tcp -t tls -p ":38080" -T udp -P "8.8.8.8:53" -C proxy.crt -K proxy.key`  
二级TCP代理VPS_02，IP:33.33.33.33  
`proxy tcp -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key`  
三级TCP代理(本地)  
`proxy udp -p ":5353" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key`  
那么访问本地5353端口就是通过加密TCP隧道，通过VPS_01访问8.8.8.8的53端口。

### 3.6 指定出口IP

当UDP代理当上级类型（参数：-T）是udp当时候，支持指定`出口IP`。使用`--bind-listen`参数，就可以开启客户端用`入口IP`连接过来的，就用`入口IP`作为`出口IP`
访问目标的功能。如果绑定了不正确的IP会导致代理不能工作。

`proxy udp -p ":33080" -T udp -P "192.168.22.33:2222" -B`

### 3.7 查看帮助

`proxy help udp`

## 4.内网穿透

### 4.1、原理说明

内网穿透，分为两个版本，“多链接版本”和“多路复用版本”，一般像web服务这种不是长时间连接的服务建议用“多链接版本”，如果是要保持长时间连接建议使用“多路复用版本”。

1. 多链接版本，对应的子命令是tserver，tclient，tbridge。
1. 多路复用版本，对应的子命令是server，client，bridge。
1. 多链接版本和多路复用版本的参数和使用方式完全一样。
1. 多路复用版本的server，client可以开启压缩传输，参数是--c。
1. server，client要么都开启压缩，要么都不开启，不能只开一个。

下面的教程以“多路复用版本”为例子，说明使用方法。  
内网穿透由三部分组成:client端，server端，bridge端；client和server主动连接bridge端进行桥接。

### 4.2、TCP普通用法

背景:

- 公司机器A提供了web服务80端口
- 有VPS一个，公网IP:22.22.22.22

需求:  
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口

步骤:

1. 在vps上执行  
   `proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
   `proxy server -r ":28080@:80" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`

1. 在公司机器A上面执行  
   `proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 完成

### 4.3、微信接口本地开发

背景:

- 自己的笔记本提供了nginx服务80端口
- 有VPS一个，公网IP:22.22.22.22

需求:  
在微信的开发帐号的网页回调接口配置里面填写地址:http://22.22.22.22/calback.php  
然后就可以访问到笔记本的80端口下面的calback.php，如果需要绑定域名，可以用自己的域名  
比如:wx-dev.xxx.com解析到22.22.22.22，然后在自己笔记本的nginx里  
配置域名wx-dev.xxx.com到具体的目录即可。

步骤:

1. 在vps上执行，确保vps的80端口没被其它程序占用。  
   `proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
   `proxy server -r ":80@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 在自己笔记本上面执行  
   `proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 完成

### 4.4、UDP普通用法

背景:

- 公司机器A提供了DNS解析服务，UDP:53端口
- 有VPS一个，公网IP:22.22.22.22

需求:  
在家里能够通过设置本地dns为22.22.22.22，使用公司机器A进行域名解析服务。

步骤:

1. 在vps上执行  
   `proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
   `proxy server --udp -r ":53@:53" -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`

1. 在公司机器A上面执行  
   `proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 完成

### 4.5、高级用法一

背景:

- 公司机器A提供了web服务80端口
- 有VPS一个，公网IP:22.22.22.22

需求:  
为了安全，不想在VPS上能够访问到公司机器A，在家里能够通过访问本机的28080端口，  
通过加密隧道访问到公司机器A的80端口。

步骤:

1. 在vps上执行  
   `proxy bridge -p ":33080" -C proxy.crt -K proxy.key`

1. 在公司机器A上面执行  
   `proxy client -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 在家里电脑上执行  
   `proxy server -r ":28080@:80" -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 完成

### 4.6、高级用法二

提示:  
如果同时有多个client连接到同一个bridge，需要指定不同的key，可以通过--k参数设定，--k可以是任意唯一字符串，  
只要在同一个bridge上唯一即可。  
server连接到bridge的时候，如果同时有多个client连接到同一个bridge，需要使用--k参数选择client。  
暴露多个端口重复-r参数即可.-r格式是:"本地IP:本地端口@clientHOST:client端口"。

背景:

- 公司机器A提供了web服务80端口，ftp服务21端口
- 有VPS一个，公网IP:22.22.22.22

需求:  
在家里能够通过访问VPS的28080端口访问到公司机器A的80端口  
在家里能够通过访问VPS的29090端口访问到公司机器A的21端口

步骤:

1. 在vps上执行  
   `proxy bridge -p ":33080" -C proxy.crt -K proxy.key`  
   `proxy server -r ":28080@:80" -r ":29090@:21" --k test -P "127.0.0.1:33080" -C proxy.crt -K proxy.key`

1. 在公司机器A上面执行  
   `proxy client --k test -P "22.22.22.22:33080" -C proxy.crt -K proxy.key`

1. 完成

### 4.7.server的-r参数

-r完整格式是:`PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT`

4.7.1.协议PROTOCOL:tcp、udp、ptcp、pudp。  
比如: `-r "udp://:10053@:53" -r "tcp://:10800@:1080" -r ":8080@:80"`  
如果指定了--udp参数，PROTOCOL默认为udp，那么:`-r ":8080@:80"`默认为udp;  
如果没有指定--udp参数，PROTOCOL默认为tcp，那么:`-r ":8080@:80"`默认为tcp;

4.7.2.CLIENT_KEY:默认是default。  
比如: -r "udp://:10053@[test1]:53" -r "tcp://:10800@[test2]:1080" -r ":8080@:80"  
如果指定了--k参数，比如--k test，那么:`-r ":8080@:80"`CLIENT_KEY默认为test;  
如果没有指定--k参数，那么:`-r ":8080@:80"`CLIENT_KEY默认为default;

4.7.3.LOCAL_IP为空默认是:`0.0.0.0`，CLIENT_LOCAL_HOST为空默认是:`127.0.0.1`;

### 4.8.server和client通过代理连接bridge

有时候server或者client所在的网络不能直接访问外网，需要通过一个https或者socks5代理才能上网，那么这个时候  
-J参数就可以帮助你让server或者client通过https或者socks5代理去连接bridge。  
-J参数格式如下:

https代理写法:  
代理需要认证，用户名:username 密码:password  
https://username:password@host:port  
代理不需要认证  
https://host:port

socks5代理写法:  
代理需要认证，用户名:username 密码:password  
socks5://username:password@host:port  
代理不需要认证  
socks5://host:port

host:代理的IP或者域名  
port:代理的端口

### 4.9.内网穿透HTTP服务

通常HTTP请求客户端会使用server的ip和端口去设置HOST字段，但是与期望的后端实际HOST不一样，这样就造成了tcp是通的， 但后端依赖HOST字段定位虚拟主机就不能工作.现在用`--http-host`
参数强制设置http头部的HOST字段值为后端实际的域名和端口即可轻松解决,另外使用了`--http-host`参数之后,会在每个HTTP请求的头部增加两个头部 `X-Forwarded-For` 和 `X-Real-IP`
值是客户端IP, 这样后端http服务就能方便的获取用户端真实IP地址。

`server`的--http-host参数格式如下:

`--http-host www.test.com:80@2200`，如果server监听多个端口，只需要重复`--http-host`参数设置每个端口的HOST即可。

实例:

比如client本地nginx，127.0.0.1:80提供了web服务，其中绑定了一个域名`local.com`。

那么server的启动参数可以如下:

`proxy server -P :30000 -r :2500@127.0.0.1:80 --http-host local.com@2500`

解释:

`-r :2500@127.0.0.1:80` 和 `--http-host local.com:80@2500` 里面的2500端口是server本地监听的端口

当使用http协议请求server的ip:2500端口的时候，http的头部HOST字段就会被设置为`local.com`。

另外只有设置了`--http-host`参数

### 4.10 关于流量统计

如果单独启动一个server对接上级是proxy-admin控制面板，需要在上级控制面板里面新建一个映射，获得个映射规则的ID，

然后启动server的时候加上参数 --server-id=映射规则的ID 才能统计到流量。

### 4.11 关于p2p

内网穿透支持在server和client网络情况满足的情况下，server和client之间通过p2p直接连接，开启方法是：

在启动bridge，server，client到时候都加上`--p2p`参数即可。server的-r参数可以针对端口是否启用p2p（ptcp和pudp）。

如果server和client之间p2p打洞失败，那么会自动切换使用bridge中转传输数据。

### 4.12 客户端key白名单

内网穿透bridge可以设置客户端key白名单，参数是--client-keys，格式可以是：

a.文件名,文件内容一行一个客户端key只能包含数字字母下划线，也就是客户端启动参数--k的值，只有客户端key在此白名单的客户端才能连接。# 开头的行，为注释。

b.“base64://”开头的base64编码的上面a说明的文件内容，比如：base64://ajfpoajsdfa=

c.”str://“开头的英文逗号分割的多个key，比如：str://default,company,school

默认是空，允许所有key。

### 4.13 网络NAT类型判断

nat类型判断,方便查看网络是否支持p2p，可以执行：`proxy tools -a nattype`

### 4.14 查看帮助

`proxy help bridge`  
`proxy help server`  
`proxy help client`

## 5.SOCKS5代理

提示:

SOCKS5代理，支持CONNECT，UDP协议，不支持BIND，支持用户名密码认证。

***如果你的VPS是阿里云，腾讯云这种VPS，就是ifconfig看不见你的公网IP，只能看见内网IP，***

***那么需要加上`-g VPS公网IP`参数，SOCKS5代理的UDP功能才能正常工作。***

***socks5的udp功能默认关闭，可以通过--udp开启，默认是握手随机端口，可以通过固定一个端口提高性能， 通过参数--udp-port 0设置，0代表随机选择一个空闲端口，也可以手动指定一个具体端口。***

### 5.1 普通SOCKS5代理

`proxy socks -t tcp -p "0.0.0.0:38080" --udp-port 0 --udp`

-p参数支持的写法：

```text
  -p ":8081"  监听8081
  -p ":8081,:8082"  监听8081和8082
  -p ":8081,:8082,:9000-9999" 监听8081和8082以及9000,9001至9999，共1002个端口
```

### 5.2.普通二级SOCKS5代理

![5.2](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/socks-2.png)  
使用本地端口8090，假设上级SOCKS5代理是`22.22.22.22:8080`  
`proxy socks -t tcp -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  --udp-port 0 --udp`  
我们还可以指定网站域名的黑白名单文件，一行一个域名，匹配规则是最右匹配，比如:baidu.com，匹配的是*.*
.baidu.com，黑名单的域名域名直接走上级代理，白名单的域名不走上级代理;如果域名即在黑名单又在白名单中，那么黑名单起作用。  
`proxy socks -p "0.0.0.0:8090" -T tcp -P "22.22.22.22:8080"  -b blocked.txt -d direct.txt --udp-port 0 --udp`

### 5.3.SOCKS二级代理(加密)

![5.3](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/socks-tls-2.png)  
一级SOCKS代理(VPS，IP:22.22.22.22)  
`proxy socks -t tls -p ":38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`

二级SOCKS代理(本地Linux)  
`proxy socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080。

二级SOCKS代理(本地windows)  
`proxy.exe socks -t tcp -p ":8080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`  
然后设置你的windos系统中，需要通过代理上网的程序的代理为socks5模式，地址为：127.0.0.1，端口为：8080，程序即可通过加密通道通过vps上网。

### 5.4.SOCKS三级代理(加密)

![5.4](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/socks-tls-3.png)  
一级SOCKS代理VPS_01，IP:22.22.22.22  
`proxy socks -t tls -p ":38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`  
二级SOCKS代理VPS_02，IP:33.33.33.33  
`proxy socks -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`  
三级SOCKS代理(本地)  
`proxy socks -t tcp -p ":8080" -T tls -P "33.33.33.33:28080" -C proxy.crt -K proxy.key --udp-port 0 --udp`  
那么访问本地的8080端口就是访问一级SOCKS代理上面的代理端口38080。

### 5.5.SOCKS代理流量强制走上级SOCKS代理

默认情况下，proxy会智能判断一个网站域名是否无法访问，如果无法访问才走上级SOCKS代理.通过--always可以使全部SOCKS代理流量强制走上级SOCKS代理。  
`proxy socks --always -t tls -p ":28080" -T tls -P "22.22.22.22:38080" -C proxy.crt -K proxy.key --udp-port 0 --udp`

### 5.6.SOCKS通过SSH中转

![5.6](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/socks-ssh.png)  
说明:ssh中转的原理是利用了ssh的转发功能，就是你连接上ssh之后，可以通过ssh代理访问目标地址。  
假设有:vps

- IP是2.2.2.2， ssh端口是22， ssh用户名是:user， ssh用户密码是:demo
- 用户user的ssh私钥名称是user.key

#### *5.6.1 ssh用户名和密码的方式*

本地SOCKS5代理28080端口，执行:  
`proxy socks -T ssh -P "2.2.2.2:22" -u user -D demo -t tcp -p ":28080"  --udp-port 0 --udp`

#### *5.6.2 ssh用户名和密钥的方式*

本地SOCKS5代理28080端口，执行:  
`proxy socks -T ssh -P "2.2.2.2:22" -u user -S user.key -t tcp -p ":28080"  --udp-port 0 --udp`

那么访问本地的28080端口就是通过VPS访问目标地址。

### 5.7.认证

对于socks5代理协议我们可以进行用户名密码认证，认证的用户名和密码可以在命令行指定  
`proxy socks -t tcp -p ":33080" -a "user1:pass1" -a "user2:pass2"`  
多个用户，重复-a参数即可。  
也可以放在文件中，格式是一行一个"用户名:密码"，然后用-F指定。  
`proxy socks -t tcp -p ":33080" -F auth-file.txt`

### 5.8.KCP协议传输

KCP协议需要--kcp-key参数设置一个密码用于加密解密数据

一级HTTP代理(VPS，IP:22.22.22.22)  
`proxy socks -t kcp -p ":38080" --kcp-key mypassword -g 22.22.22.22`

二级HTTP代理(本地Linux)  
`proxy socks -t tcp -p ":8080" -T kcp -P "22.22.22.22:38080" --kcp-key mypassword`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080，数据通过kcp协议传输。

提示:

当本地使用kcp协议时,需要用-g指定vps公网IP,socks5的UDP功能才能正常使用.这时-g是返回给客户端的UDP地址中的IP地址.

### 5.9.自定义DNS

--dns-address和--dns-ttl参数，用于自己指定proxy访问域名的时候使用的dns（--dns-address）  
以及解析结果缓存时间（--dns-ttl）秒数，避免系统dns对proxy的干扰，另外缓存功能还能减少dns解析时间提高访问速度。  
比如：  
`proxy socks -p ":33080" --dns-address "8.8.8.8:53" --dns-ttl 300`

### 5.10 自定义加密

proxy的socks代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据，除此之外还支持在tls和kcp之后进行自定义加密，也就是说自定义加密和tls|kcp是可以联合使用的，内部采用AES256加密，使用的时候只需要自己定义一个密码即可，  
加密分为两个部分，一部分是本地(-z)是否加密解密，一部分是与上级(-Z)传输是否加密解密。

自定义加密要求两端都是proxy才可以。

下面分别用二级，三级为例:

二级实例  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy socks -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

三级实例  
一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy socks -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888`  
本地三级执行:  
`proxy socks -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

### 5.11 压缩传输

proxy的socks代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据，在自定义加密之前还可以  
对数据进行压缩，也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的，压缩分为两个部分，  
一部分是本地(-m)是否压缩传输，一部分是与上级(-M)传输是否压缩。

压缩要求两端都是proxy才可以，压缩也在一定程度上保护了(加密)数据。

下面分别用二级，三级为例:

二级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -m -p :7777`  
本地二级执行:  
`proxy socks -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

三级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy socks -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy socks -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888`  
本地三级执行:  
`proxy socks -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

### 5.12 负载均衡

SOCKS代理支持上级负载均衡，多个上级重复-P参数即可。

`proxy socks --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -p :33080 -t tcp`

### 5.12.1 设置重试间隔和超时时间

`proxy socks --lb-method=leastconn --lb-retrytime 300 --lb-timeout 300 -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -p :33080 -t tcp`

### 5.12.2 设置权重

`proxy socks --lb-method=weight -T tcp -P 1.1.1.1:33080?w=1 -P 2.1.1.1:33080?w=2 -P 3.1.1.1:33080?w=1 -p :33080 -t tcp`

### 5.12.3 使用目标地址选择上级

`proxy socks --lb-hashtarget --lb-method=hash -T tcp -P 1.1.1.1:33080 -P 2.1.1.1:33080 -P 3.1.1.1:33080 -p :33080 -t tcp`

### 5.13 限速

限速100K，通过`-l`参数即可指定，比如:100K 2000K 1M . 0意味着无限制。

`proxy socks -t tcp -p 2.2.2.2:33080 -l 100K`

### 5.14 指定出口IP

`--bind-listen`参数，就可以开启客户端用`入口IP`连接过来的，就用`入口IP`作为`出口IP`访问目标网站的功能。如果绑定了不正确的IP会导致代理不能工作。

普通多IP用法：

`proxy socks -t tcp -p :33080 --bind-listen`

上面的写法适合本机有大量IP，不想一个一个单独写，比如本机有IP`3.3.3.3`,当客户端使用`3.3.3.3`连接代理，那么`3.3.3.3`就会作为客户端的出口IP， 这里的`3.3.3.3`就是`入口IP`和`出口IP`
，你会发现`入口IP`和`出口IP`是一样的，`--bind-listen`这种模式人为不能干预`入口IP`和`出口IP`。

单独用法：

本机有IP：`5.5.5.5`，开一个代理端口`6666`，凡是通过`5.5.5.5:6666`端口进来的，都使用`5.5.5.5`作为出口IP。

命令如下：

`proxy socks -t tcp -p 5.5.5.5:6666 --bind-listen`

#### 自由指定出口IP

上面的`--bind-listen`参数虽然可以指定出口IP，但是`入口IP`和`出口IP`不能人为干预。如果想入口IP和出口IP不一样， 可以使用`--bind-ip`参数，格式：`IP:端口`，比如： `1.1.1.1:8080`
，`[2000:0:0:0:0:0:0:1]:8080`。 多个绑定需求，可以重复`--bind-ip`参数即可。

比如本机有IP`5.5.5.5`,`6.6.6.6`，监听了`8888`和`7777`两个端口，命令如下：

`proxy socks -t tcp -p :8888,:7777 --bind-ip 5.5.5.5:7777 --bind-ip 6.6.6.6:8888`

那么客户端访问`7777`端口，出口IP就是`5.5.5.5`，访问`8888`端口，出口IP就是`6.6.6.6`，如果同时设置了`--bind-ip`和`--bind-listen`,`--bind-ip`的优先级更高。

另外，`--bind-ip`参数的`IP`部分，支持指定`网卡名称`，`通配符`，还能指定多个，详细说明如下：

- 指定网卡名称，比如：`--bind-ip eth0:7777`，那么客户端访问`7777`端口，出口IP就是eth0网卡的IP。
- 网卡名称支持通配符，比如：`--bind-ip eth0.*:7777`，那么客户端访问`7777`端口，出口IP就是`eth0.`开头网卡的IP中随机选择的一个。
- IP支持通配符，比如：`--bind-ip 192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器所有IP中，匹配`192.168.?.*`的IP中随机选择的一个。
- 还可以是网卡名称和IP的多个组合，多个使用半角逗号分割，比如：`--bind-ip pppoe??,192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器网卡名称匹配`pppoe??`
  和机器所有IP中匹配`192.168.?.*`的IP中随机选择的一个。
- 通配符`*`代表0至任意多个字符，`?`代表1个字符。
- 如果网卡IP发生变化，也会实时生效。
- 可以通过`--bind-refresh`参数，指定刷新本地网卡信息的间隔，默认`5`，单位秒。

### 5.15 级联认证

SOCKS5支持级联认证，-A可以设置上级认证信息。

上级:

`proxy socks -t tcp -p 2.2.2.2:33080 -a user:pass`

本地:

`proxy socks -T tcp -P 2.2.2.2:33080 -A user:pass -t tcp -p :33080`

请更多认证细节，请参考`9.API认证` 和 `10.本地认证`

### 5.16 证书参数使用base64数据

默认情况下-C，-K参数是crt证书和key文件的路径，

如果是base64://开头，那么就认为后面的数据是base64编码的，会解码后使用。

### 5.17 智能模式

智能模式设置，可以是intelligent|direct|parent三者之一。  
默认是:parent。  
每个值的含义如下:  
`--intelligent=direct`，不在blocked里面的目标都直连。  
`--intelligent=parent`，不在direct里面的目标都走上级。  
`--intelligent=intelligent`，blocked和direct里面都没有的目标，智能判断是否使用上级访问目标。

### 5.18 固定UDP功能端口

默认情况下socks5的UDP功能的端口号，proxy是按着`rfc1982草案`要求，使用协议握手过程中随机指定，不需要提前指定。

但是某些情况下，需要固定UDP功能端口，可以通过参数`--udp-port 端口号`用来固定UDP功能的端口号，比如：

`proxy socks -t tcp -p "0.0.0.0:38080" --udp-port 38080`

### 5.19 查看帮助

`proxy help socks`

## 6.SPS协议转换

### 6.1 功能介绍

代理协议转换使用的是sps子命令，sps可以把已经存在的http(s)代理或者socks5代理或ss代理转换为一个端口同时支持http(s)和socks5和ss的代理，而且http(s)代理支持正向代理和反向代理(SNI)
，当上级是SOCKS5时,转换后的SOCKS5或者SS代理仍然支持UDP功能；另外对于已经存在的http(s)代理或者socks5代理，支持tls、tcp、kcp三种模式，支持链式连接，也就是可以多个sps结点层级连接构建加密通道。

`ss`功能支持的加密方法为:aes-128-cfb ， aes-128-ctr ， aes-128-gcm ， aes-192-cfb ， aes-192-ctr ， aes-192-gcm ， aes-256-cfb ，
aes-256-ctr ， aes-256-gcm ， bf-cfb ， cast5-cfb ， chacha20 ， chacha20-ietf ， chacha20-ietf-poly1305 ， des-cfb ， rc4-md5 ，
rc4-md5-6 ， salsa20 ， xchacha20

本地监听端口-p参数支持的写法：

```text
  -p ":8081"  监听8081
  -p ":8081,:8082"  监听8081和8082
  -p ":8081,:8082,:9000-9999" 监听8081和8082以及9000,9001至9999，共1002个端口
```

提示:

当本地使用kcp协议时,需要用-g指定vps公网IP,socks5的UDP功能才能正常使用.这时-g是返回给客户端的UDP地址中的IP地址.

ss的udp功能默认关闭，可以通过--ssudp开启。socks5的udp功能默认关闭，可以通过--udp开启，默认是握手随机端口，可以通过固定一个端口提高性能， 通过参数--udp-port
0设置，0代表随机选择一个空闲端口，也可以手动指定一个具体端口。

### 6.2 HTTP(S)转HTTP(S)+SOCKS5+SS

转换后的SOCKS5或者SS代理不具备UDP功能。    
假设已经存在一个普通的http(s)代理：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，ss加密方式:aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`

假设已经存在一个tls的http(s)代理：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，tls需要证书文件，ss加密方式:
aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S http -T tls -P 127.0.0.1:8080 -t tcp -p :18080 -C proxy.crt -K proxy.key -h aes-192-cfb -j pass`

假设已经存在一个kcp的http(s)代理（密码是：demo123）：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，ss加密方式:
aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S http -T kcp -P 127.0.0.1:8080 -t tcp -p :18080 --kcp-key demo123 -h aes-192-cfb -j pass`

### 6.3 SOCKS5转HTTP(S)+SOCKS5+SS

如果上级是SOCKS5而且支持UDP功能，那么转换后的SOCKS5或者SS代理仍然支持UDP功能。  
假设已经存在一个普通的socks5代理：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，ss加密方式:aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S socks -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`

假设已经存在一个tls的socks5代理：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，tls需要证书文件，ss加密方式:aes-192-cfb，ss密码:
pass。  
命令如下：  
`proxy sps -S socks -T tls -P 127.0.0.1:8080 -t tcp -p :18080 -C proxy.crt -K proxy.key -h aes-192-cfb -j pass`

假设已经存在一个kcp的socks5代理（密码是：demo123）：127.0.0.1:8080，现在我们把它转为同时支持http(s)和socks5和ss的普通代理，转换后的本地端口为18080，ss加密方式:
aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S socks -T kcp -P 127.0.0.1:8080 -t tcp -p :18080 --kcp-key demo123 -h aes-192-cfb -j pass`

### 6.4 SS转HTTP(S)+SOCKS5+SS

SPS上级和本地支持ss协议，上级可以是SPS或者标准的ss服务，转换后的SOCKS5或者SS代理不具备UDP功能。  
假设已经存在一个普通的SS或者SPS代理(开启了ss，加密方式:aes-256-cfb，密码:demo)：127.0.0.1:8080，现在我们把它转为同时支持http(s)
和socks5和ss的普通代理，转换后的本地端口为18080，转换后的ss加密方式:aes-192-cfb，ss密码:pass。  
命令如下：  
`proxy sps -S ss -H aes-256-cfb -J pass -T tcp -P 127.0.0.1:8080 -t tcp -p :18080 -h aes-192-cfb -j pass`。

### 6.5 链式连接

![6.4](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/sps-tls.png)  
上面提过多个sps结点可以层级连接构建加密通道，假设有如下vps和家里的pc电脑。  
vps01：2.2.2.2  
vps02：3.3.3.3  
现在我们想利用pc和vps01和vps02构建一个加密通道，本例子用tls加密也可以用kcp，在pc上访问本地18080端口就是访问vps01的本地8080端口。  
首先在vps01(2.2.2.2)上我们运行一个只有本地可以访问的http(s)代理，执行：  
`proxy http -t tcp -p 127.0.0.1:8080`

然后在vps01(2.2.2.2)上运行一个sps结点，执行：  
`proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tls -p :8081 -C proxy.crt -K proxy.key`

然后在vps02(3.3.3.3)上运行一个sps结点，执行：  
`proxy sps -S http -T tls -P 2.2.2.2:8081 -t tls -p :8082 -C proxy.crt -K proxy.key`

然后在pc上运行一个sps结点，执行：  
`proxy sps -S http -T tls -P 3.3.3.3:8082 -t tcp -p :18080 -C proxy.crt -K proxy.key`

完成。

### 6.6 认证功能

sps支持http(s)\socks5代理认证，可以级联认证，有四个重要的信息:  
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

对于sps代理我们可以进行用户名密码认证，认证的用户名和密码可以在命令行指定  
`proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -a "user1:pass1:0:0:" -a "user2:pass2:0:0:"`  
多个用户，重复-a参数即可。  
也可以放在文件中，格式是一行一个`用户名:密码:连接数:速率:上级`，然后用-F指定。  
`proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -F auth-file.txt`

如果上级有认证，下级可以通过-A参数设置认证信息，比如:  
上级:`proxy sps -S http -T tcp -P 127.0.0.1:8080 -t tcp -p ":33080" -a "user1:pass1:0:0:" -a "user2:pass2:0:0:"`  
下级:`proxy sps -S http -T tcp -P 127.0.0.1:8080 -A "user1:pass1" -t tcp -p ":33080" `

请更多认证细节，请参考`9.API认证` 和 `10.本地认证`

### 6.7 多个上级

如果存在多个上级，可以通过多个-P指定。

比如：

`proxy sps -P http://127.0.0.1:3100  -P socks5://127.0.0.1:3200`

`-P`完整格式如下：

`protocol://a:b@2.2.2.2:33080#1`

下面对每部分进行解释：

`protocol://` 是协议类型，可能的类型以及含有如下:

```text  
http  等价于 -S http -T tcp  
https  等价于 -S http -T tls --parent-tls-single ， 也就是http(s)代理 over TLS  
https2 等价于 -S http -T tls  
socks5  等价于 -S socks -T tcp  
socks5s 等价于 -S socks -T tls --parent-tls-single ， 也就是socks over TLS  
socks5s2 等价于 -S socks -T tls  
ss     等价于 -S ss -T tcp  
httpws    等价于  -S http -T ws  
httpwss   等价于  -S http -T wss  
socks5ws  等价于  -S socks -T ws  
socks5wss 等价于  -S socks -T wss  
```  

`a:b`是代理认证的用户名密码，如果是ss，`a`是加密方法，`b`是密码，没有用户名密码可以留空比如：`http://2.2.2.2:33080`，如果用户名和密码保护特殊符号可以使用urlencode进行编码。

`2.2.2.2:33080`是上级地址，格式是：`IP（或域名）:端口`，如果底层是ws/wss协议还可以带上路径，比如：`2.2.2.2:33080/ws`;  
还能通过附加查询参数`m`和`k`设置`ws\wss`的`加密方法`和`密码`，比如：`2.2.2.2:33080/ws?m=aes-192-cfb&k=password`

`#1`多个上级的负载均衡是权重策略时候，设置的权重，很少用到。

### 6.8 自定义加密

proxy的sps代理在tcp之上可以通过tls标准加密以及kcp协议加密tcp数据，除此之外还支持在tls和kcp之后进行  
自定义加密，也就是说自定义加密和tls|kcp是可以联合使用的，内部采用AES256加密，使用的时候只需要自己定义  
一个密码即可，加密分为两个部分，一部分是本地(-z)是否加密解密，一部分是与上级(-Z)传输是否加密解密。

自定义加密要求两端都是proxy才可以。

下面分别用二级，三级为例:

假设已经存在一个http(s)代理:`6.6.6.6:6666`

二级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy sps -S http -T tcp -P 6.6.6.6:6666 -t tcp -z demo_password -p :7777`  
本地二级执行:  
`proxy sps -T tcp -P 2.2.2.2:777 -Z demo_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

三级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy sps -S http -T tcp -P 6.6.6.6:6666 -t tcp -z demo_password -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy sps -T tcp -P 2.2.2.2:7777 -Z demo_password -t tcp -z other_password -p :8888`  
本地三级执行:  
`proxy sps -T tcp -P 3.3.3.3:8888 -Z other_password -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级加密传输访问目标网站。

### 6.9 压缩传输

proxy的sps代理在tcp之上可以通过自定义加密和tls标准加密以及kcp协议加密tcp数据，在自定义加密之前还可以  
对数据进行压缩，也就是说压缩功能和自定义加密和tls|kcp是可以联合使用的，压缩分为两个部分，  
一部分是本地(-m)是否压缩传输，一部分是与上级(-M)传输是否压缩。

压缩要求两端都是proxy才可以，压缩也在一定程度上保护了(加密)数据。

下面分别用二级，三级为例:

二级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy sps -t tcp -m -p :7777`  
本地二级执行:  
`proxy sps -T tcp -P 2.2.2.2:777 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

三级实例

一级vps(ip:2.2.2.2)上执行:  
`proxy sps -t tcp -m -p :7777`  
二级vps(ip:3.3.3.3)上执行:  
`proxy sps -T tcp -P 2.2.2.2:7777 -M -t tcp -m -p :8888`  
本地三级执行:  
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080`  
这样通过本地代理8080访问网站的时候就是通过与上级压缩传输访问目标网站。

### 6.10 禁用协议

SPS默认情况下一个端口支持http(s)和socks5两种代理协议，我们可以通过参数禁用某个协议  
比如:  
1.禁用HTTP(S)代理功能只保留SOCKS5代理功能，参数:`--disable-http`。  
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080 --disable-http`

1.禁用SOCKS5代理功能只保留HTTP(S)代理功能，参数:`--disable-socks`。  
`proxy sps -T tcp -P 3.3.3.3:8888 -M -t tcp -p :8080 --disable-socks`

### 6.11 限速

假设存在SOCKS5上级:

`proxy socks -p 2.2.2.2:33080 -z password -t tcp`

sps下级，限速100K

`proxy sps -S socks -P 2.2.2.2:33080 -T tcp -Z password -l 100K -t tcp -p :33080`

通过`-l`参数即可指定，比如:100K 2000K 1M . 0意味着无限制。

### 6.12 指定出口IP

`--bind-listen`参数，就可以开启客户端用`入口IP`连接过来的，就用`入口IP`作为`出口IP`访问目标网站的功能。如果绑定了不正确的IP会导致代理不能工作。

普通多IP用法：

`proxy sps -t tcp -p :33080 --bind-listen`

上面的写法适合本机有大量IP，不想一个一个单独写，比如本机有IP`3.3.3.3`,当客户端使用`3.3.3.3`连接代理，那么`3.3.3.3`就会作为客户端的出口IP， 这里的`3.3.3.3`就是`入口IP`和`出口IP`
，你会发现`入口IP`和`出口IP`是一样的，`--bind-listen`这种模式人为不能干预`入口IP`和`出口IP`。

单独用法：

本机有IP：`5.5.5.5`，开一个代理端口`6666`，凡是通过`5.5.5.5:6666`端口进来的，都使用`5.5.5.5`作为出口IP。

命令如下：

`proxy sps -t tcp -p 5.5.5.5:6666 --bind-listen`

#### 自由指定出口IP

上面的`--bind-listen`参数虽然可以指定出口IP，但是`入口IP`和`出口IP`不能人为干预。如果想入口IP和出口IP不一样， 可以使用`--bind-ip`参数，格式：`IP:端口`，比如： `1.1.1.1:8080`
，`[2000:0:0:0:0:0:0:1]:8080`。 多个绑定需求，可以重复`--bind-ip`参数即可。

比如本机有IP`5.5.5.5`,`6.6.6.6`，监听了`8888`和`7777`两个端口，命令如下：

`proxy sps -t tcp -p :8888,:7777 --bind-ip 5.5.5.5:7777 --bind-ip 6.6.6.6:8888`

那么客户端访问`7777`端口，出口IP就是`5.5.5.5`，访问`8888`端口，出口IP就是`6.6.6.6`，如果同时设置了`--bind-ip`和`--bind-listen`,`--bind-ip`的优先级更高。

另外，`--bind-ip`参数的`IP`部分，支持指定`网卡名称`，`通配符`，还能指定多个，详细说明如下：

- 指定网卡名称，比如：`--bind-ip eth0:7777`，那么客户端访问`7777`端口，出口IP就是eth0网卡的IP。
- 网卡名称支持通配符，比如：`--bind-ip eth0.*:7777`，那么客户端访问`7777`端口，出口IP就是`eth0.`开头网卡的IP中随机选择的一个。
- IP支持通配符，比如：`--bind-ip 192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器所有IP中，匹配`192.168.?.*`的IP中随机选择的一个。
- 还可以是网卡名称和IP的多个组合，多个使用半角逗号分割，比如：`--bind-ip pppoe??,192.168.?.*:7777`，那么客户端访问`7777`端口，出口IP就是机器网卡名称匹配`pppoe??`
  和机器所有IP中匹配`192.168.?.*`的IP中随机选择的一个。
- 通配符`*`代表0至任意多个字符，`?`代表1个字符。
- 如果网卡IP发生变化，也会实时生效。
- 可以通过`--bind-refresh`参数，指定刷新本地网卡信息的间隔，默认`5`，单位秒。

### 6.13 证书参数使用base64数据

默认情况下-C，-K参数是crt证书和key文件的路径，

如果是base64://开头，那么就认为后面的数据是base64编码的，会解码后使用。

### 6.14 独立服务

sps一个端口可完成完整功能的代理`http(s)\socks\ss`功能。  
下面这个命令，就是一键开启http(s)\ss\socks服务,同时开启socks5的udp和ss的udp功能。  
`proxy sps -p :33080 --ssudp --udp --udp-port 0`

### 6.15 目标重定向

sps功能提供的http(s)\socks5\ss代理功能，客户端通过sps代理去连接指定的“目标”，这个“目标”一般是网站也可能是任意的tcp地址，  
网站”目标“一般是foo.com:80,foo.com:443，sps支持使用--rewrite参数指定一个“目标”重定向规则文件，对目标进行重定向，客户端是无感知的，  
比如你对“目标”：demo.com:80重定向到192.168.0.12:80，那么客户端访问网站demo.com，其实访问的是192.168.0.12提供的网站服务。  
“目标”重定向规则文件示例：

```text  
# example  
www.a.com:80     10.0.0.2:8080  
**.b.com:80      10.0.0.2:80  
192.168.0.11:80  10.0.0.2:8080  
```  

### 6.16 固定UDP功能端口

默认情况下sps的socks5的UDP功能的端口号是按着`rfc1982草案`要求，使用协议握手过程中随机指定，不需要提前指定。

但是某些情况下，需要固定UDP功能端口，可以通过参数`--udp-port 端口号`用来固定UDP功能的端口号，比如：

`proxy sps -t tcp -p "0.0.0.0:38080" --udp-port 38081`

需要注意的是，sps的ss功能也有UDP功能，而且ss的UDP端口和tcp端口是一样的，所以要避免socks5的UDP端口和ss的UDP端口冲突，

要指定和tcp端口不一样的端口。

### 6.17 iptables 透明代理

sps模式支持Linux系统的iptables转发支持，也就是通常所说的iptables透明代理，如果在网关设备上进行iptables透明代理，那么对通过网关联网的设备就能实现无感知的代理。

启动命令实例：

`proxy sps --redir -p :8888 -P httpws://1.1.1.1:33080`

这里假设存在一个http的上级代理1.1.1.1:33080，使用ws传输数据。

然后添加iptables规则，下面是参考规则:

```shell  
#上级proxy服务端服务器IP地址:  
proxy_server_ip=1.1.1.1

#路由器运行proxy监听的端口:  
proxy_local_port=33080  

#下面的就不用修改了  
#create a new chain named PROXY  
iptables -t nat -N PROXY  

# Ignore your PROXY server's addresses  
# It's very IMPORTANT， just be careful。  

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
iptables -t nat -A PROXY -p tcp  -j REDIRECT --to-ports $proxy_local_port
# Apply the rules to nat client  
iptables -t nat -A PREROUTING -p tcp -j PROXY  
# Apply the rules to localhost  
iptables -t nat -A OUTPUT -p tcp -j PROXY  
```  

- 清空整个链 iptables -F 链名比如iptables -t nat -F PROXY
- 删除指定的用户自定义链 iptables -X 链名 比如 iptables -t nat -X PROXY
- 从所选链中删除规则 iptables -D 链名 规则详情 比如 iptables -t nat -D PROXY -d 223.223.192.0/255.255.240.0 -j RETURN

### 6.18 查看帮助

`proxy help sps`

## 7.KCP配置

### 7.1 配置介绍

proxy的很多功能都支持kcp协议，凡是使用了kcp协议的功能都支持这里介绍的配置参数。  
所以这里统一对KCP配置参数进行介绍。

### 7.2 详细配置

所有的KCP配置参数共有17个，你可以都不用设置，他们都有默认值，如果为了或者最好的效果，  
就需要自己根据自己根据网络情况对参数进行配置。由于kcp配置很复杂需要一定的网络基础知识，  
如果想获得kcp参数更详细的配置和解说，请自行搜索。每个参数的命令行名称以及默认值和简单的功能说明如下：

```  
--kcp-key="secrect"        pre-shared secret between client and server  
--kcp-method="aes"         encrypt/decrypt method， can be: aes， aes-128， aes-192， salsa20， blowfish，  
twofish， cast5， 3des， tea， xtea， xor， sm4， none  
--kcp-mode="fast"       profiles: fast3， fast2， fast， normal， manual  
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
参数：--kcp-mode中的四种fast3， fast2， fast， normal模式，  
相当于设置了下面四个参数：  
normal：`--nodelay=0 --interval=40 --resend=2 --nc=1`  
fast ：`--nodelay=0 --interval=30 --resend=2 --nc=1`  
fast2：`--nodelay=1 --interval=20 --resend=2 --nc=1`  
fast3：`--nodelay=1 --interval=10 --resend=2 --nc=1`

## 8.安全DNS

### 8.1 介绍

众所周知DNS是UDP端口53提供的服务，但是随着网络的发展一些知名DNS服务器也支持TCP方式dns查询，比如谷歌的8.8.8.8，proxy的DNS防污染服务器原理就是在本地启动一个proxy的DNS代理服务器，它用TCP的方式通过上级代理进行dns查询。如果它和上级代理通讯采用加密的方式，那么就可以进行安全无污染的DNS解析，还支持独立服务，并发解析，支持增强的hosts文件功能支持灵活的并发解析转发。

dns解析顺序：  
1.使用参数--hosts解析。  
2.要解析的域名在1中没有找到，就使用参数--forward规则解析。  
3.要解析的域名在1和2中都没有找到，就使用默认--default解析，--default默认行为参数值有三种：proxy，direct，system。  
三种参数值解释如下：  
proxy：通过上级去连接-q参数指定的dns服务器去解析域名。  
direct：通过本地网络去连接-q参数指定的dns服务器去解析域名。  
system：通过系统dns去解析域名。

提示：  
--hosts 参数指定的host文件格式和系统hosts文件一致，而且域名支持通配符，可以参考hosts文件。  
--forward 参数指定的解析转发规则文件，格式可以参考resolve.rules文件，域名支持通配符，支持为每个域名指定多个dns服务器并发解析，谁最快解析成功就用谁的解析结果。  
-q 参数可以指定多个远程dns服务器，执行并发解析谁最快解析成功就用谁的解析结果，默认是：1.1.1.1，8.8.8.8，9.9.9.9，多个用逗号分割，  
比如还可以带端口：1.1.1.1，8.8.8.8#53，9.9.9.9

如果独立服务，不需要上级：  
可以执行：  
`proxy dns --default system -p :5353`  
or  
`proxy dns --default direct -p :5353`

### 8.2 使用示例

#### 8.2.1 普通HTTP(S)上级代理

假设有一个上级代理：2.2.2.2:33080  
本地执行：  
`proxy dns -S http -T tcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了DNS解析功能。

#### 8.2.2 普通SOCKS5上级代理

假设有一个上级代理：2.2.2.2:33080  
本地执行：  
`proxy dns -S socks -T tcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了DNS解析功能。

#### 8.2.3 TLS加密的HTTP(S)上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy http -t tls -C proxy.crt -K proxy.key -p :33080`  
本地执行：  
`proxy dns -S http -T tls -P 2.2.2.2:33080 -C proxy.crt -K proxy.key -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

#### 8.2.4 TLS加密的SOCKS5上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy socks -t tls -C proxy.crt -K proxy.key -p :33080`  
本地执行：  
`proxy dns -S socks -T tls -P 2.2.2.2:33080 -C proxy.crt -K proxy.key -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

#### 8.2.5 KCP加密的HTTP(S)上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy http -t kcp -p :33080`  
本地执行：  
`proxy dns -S http -T kcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

#### 8.2.6 KCP加密的SOCKS5上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy socks -t kcp -p :33080`  
本地执行：  
`proxy dns -S socks -T kcp -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

#### 8.2.7 自定义加密的HTTP(S)上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy http -t tcp -p :33080 -z password`  
本地执行：  
`proxy dns -S http -T tcp -Z password -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

#### 8.2.8 自定义加密的SOCKS5上级代理

假设有一个上级代理：2.2.2.2:33080  
上级代理执行的命令是：  
`proxy socks -t kcp -p :33080 -z password`  
本地执行：  
`proxy dns -S socks -T tcp -Z password -P 2.2.2.2:33080 -p :53`  
那么本地的UDP端口53就提供了安全防污染DNS解析功能。

## 9.API认证，限速，控制连接数，限流量

proxy的http(s)/socks5/sps代理功能，支持通过API控制用户对代理对访问。

### 通过API可以干什么？

- 用户维度，控制单个连接速率，控制最大连接数。
- IP维度，控制单个连接速率，控制最大连接数。
- 动态上级，可以根据用户或者客户端IP，动态的从API获取其上级，支持http(s)/socks5/ss上级。
- 认证每一个连接，无论是否要求客户端认证。
- 缓存认证结果，时间可以设置，减轻API压力。

#### 具体使用

proxy的http(s)/socks5/sps代理API功能，通过`--auth-url`和`--auth-nouser`和`--auth-cache`三个参数控制。  
参数`--auth-url`是HTTP API接口地址，客户端连接的时候，proxy会GET方式请求这url，带上下面参数，如果返回HTTP状态码204，代表认证成功，其它情况认为认证失败。

一个完整的请求API的示例：  
`http://test.com/auth.php?user=a&pass=b&client_addr=127.0.0.1:49892&local_addr=127.0.0.1:8100&target=http%3A%2F%2Fwww.baidu.com&service=http&sps=0`

#### 参数说明

`user和pass` 当代理开启了认证功能，这里就是客户端提供的用户名和密码。  
`client_addr` 客户端访问代理时的用的地址，格式IP:端口。  
`local_addr` 客户端访问的代理地址，格式IP:端口。  
`service` 代理类型，分为：http、socks。  
`sps` 代理是否是sps提供的，1:是，0:否。  
`target`  客户端要访问的目标，如果是http(s)代理，target是访问的具体url；如果是socks5代理，target是空。

#### 示例

假设--auth-url http://127.0.0.1:333/auth.php 指向了一个php接口地址.  
auth.php内容如下：

```php  
<?php  
#设置所有的用户名密码,前面是用户名,后面是对应密码  
$alluser=[  
    "user1"=>"pass1",  
    "user2"=>"pass2",  
    "user3"=>"pass3",  
    "user4"=>"pass4",  
];  
$proxy_ip=$_GET['local_addr'];   
$user_ip=$_GET['client_addr'];  
$service=$_GET['service'];   
$is_sps=$_GET['sps']=='1';  
$user=$_GET['user'];  
$pass=$_GET['pass'];  
$target=$_GET['target'];  
//业务逻辑判断  
//....  
$ok=false;  
foreach ($alluser as $dbuser => $dbpass) {  
    if ($user==$dbuser&&$pass==$dbpass){  
        $ok=true;  
        break;  
    }  
}  
//设置认证结果  
if($ok){  
    header("userconns:1000");  
    header("ipconns:2000");  
    header("userrate:3000");  
    header("iprate:8000");  
    header("upstream:http://127.0.0.1:3500?parent-type=tcp");  
    header("outgoing:1.1.1.1");  
    header("HTTP/1.1 204 No Content");  
}
```  

#### 解释

userconns：用户的最大连接数，不限制为0或者不设置这个头部。  
ipconns：用户IP的最大连接数，不限制为0或者不设置这个头部。  
userrate：用户的单个TCP连接速率限制，单位：字节/秒，不限制为0或者不设置这个头部。  
iprate：用户IP的单个TCP连接速率限制，单位：字节/秒，不限制为0或者不设置这个头部。  
upstream：使用的上级，没有为空，或者不设置这个头部。
outgoing: 使用的出口IP，这个设置，只有在upstream为空的的时候才有效，这里设置的IP必须是proxy所在机器具有的IP。

#### 提示

1.默认情况下，设置了`--auth-url`是强制要求客户端提供用户名和密码的；如果不需要强制要求客户端提供用户名密码，并认证，可以加上`--auth-nouser`，每次访问仍然会访问认证地址`--auth-url`
进行认证，当客户端没有发生认证信息当时候，php接口里面接收的$user认证用户名和$pass认证密码都为空。  
2.连接数限制优先级：用户认证文件连接数限制-》文件ip.limit连接数限制-》API用户连接数限制-》API的IP连接数限制-》命令行全局连接数限制。  
3.速率限制优先级：用户认证文件速率限制-》文件ip.limit速率限制-》API用户速率限制-》API的IP速率限制-》命令行全局速率限制。  
3.上级获取优先级：用户认证文件的upstream-》文件ip.limit的upstream-》API的upstream-》命令行指定的上级。  
4.`--auth-cache`认证缓存，对认证结果缓存一定时间，提升性能，降低认证接口压力，--auth-cache 单位秒，默认0, 设置0是关闭缓存。

#### upstream详细说明

1.当参数`sps`是0的时候。  
service是http时，upstream只支持http(s)代理，不支持认证，如需要认证可以用sps代替，格式：  
`http://127.0.0.1:3100?argk=argv`  
service是socks时，upstream只支持socks5代理，格式：  
`socks5://127.0.0.1:3100?argk=argv`

解释：`http://`，`socks5://` 是固定的，`127.0.0.1:3100`是上级的地址

2.当`sps`是1的时候。  
upstream支持socks5、http(s)代理，支持认证，格式：`protocol://a:b@2.2.2.2:33080?argk=argv`，具体介绍请参考SPS章节的，**多个上级**，`-P`参数的说明。  
3.参数，`?`后面`argk=argv`是参数：参数名称=参数值，多个参数用`&`连接。  
支持的所有参数如下,和命令行同名参数意义一致。

1. parent-type : 上级底层传输类型，支持 tcp,tls,ws,wss。
2. parent-ws-method : 上级底层ws传输类型的加密方法，支持的值和命令行支持的值范围一样。
3. parent-ws-password : 上级底层ws传输类型的加密密码，数字字母组成的密码。
4. parent-tls-single : 上级底层tls传输类型是否是单向tls，可以是：true | false。
5. timeout : 建立tcp连接的超时时间，数字，单位毫秒。
6. ca : 上级底层tls传输类型的ca证书文件经过base64编码后的字符串。
7. cert : 上级底层tls传输类型的证书文件经过base64编码后的字符串。
8. key : 上级底层tls传输类型的证书密钥文件经过base64编码后的字符串。
9. luminati:上级是否是luminati代理，可以是：true | false。

### 流量上报/流量统计/流量限制

proxy的http(s)/socks5/sps/tcp/udp代理功能支持流量上报,可以通过参数`--traffic-url`设置一个http接口地址。  
proxy会把连接使用的流量上报到这个地址,具体情况是,proxy发送一个HTTP到GET请求到`--traffic-url`设置的HTTP URL地址.  
上报模式有两种，可以通过`--traffic-mode`参数指定，可以是`normal`普通模式上报，也可以是`fast`快速模式上报。

1. `normal`普通模式上报  
   当连接释放时,proxy会把此次连接使用的流量上报到这个`--traffic-url`地址.
2. `fast`快速模式上报  
   对已经建立的每个连接，proxy会`定时`把这个连接产生的流量上报到这个这个`--traffic-url`地址.    
   `定时`默认是5秒，可以通过参数`--traffic-interval`修改`定时`为合适的秒数。   
   流量上报功能结合上面的API认证功能可以实现实时控制用户的流量使用,流量统计,流量限制;流量上报到接口,接口把流量数据写入数据库,然后认证API查询数据库判断用户或者IP流量使用情况,用来确定用户是否可以认证成功.

下面是一个完整的URL请求实例:

`http://127.0.0.1:8080/auth.php?act=traffic&bytes=7627&client_addr=127.0.0.1%3A63637
&id=http&out_local_addr=127.0.0.1%3A63640&out_remote_addr=127.0.0.1%3A63639 &server_addr=127.0.0.1%3A33080&target_addr=www.baidu.com%3A443 &upstream=http%3A%2F%2F127.0.0.1%3A3100&username=a`

请求参数说明:
id: 服务id标志。  
server_addr: 客户端请求的代理地址,格式: `IP:端口`。  
client_addr: 客户端地址,格式: `IP:端口`。  
target_addr: 目标地址,格式: `IP:端口`,tcp/udp代理时,这个是空。  
username: 代理认证用户名,tcp/udp代理时,这个是空。  
bytes: 此次使用的流量字节数。  
out_local_addr: 代理对外建立的TCP连接的本地地址，格式: `IP:端口`。  
out_remote_addr: 代理对外建立的TCP连接的远程地址，格式: `IP:端口`。  
upstream: 使用的上级，格式是标准URL格式，如果没有使用上级，这里是空。

#### 提示

访问的`--traffic-url`地址接口必须返回HTTP状态码`204`，流量上报才认为上报成功，返回其它状态码将认为上报流量失败，会输出相关日志。

#### 流量上报原理

![traffic](https://cdn.jsdelivr.net/gh/snail007/goproxy@master/doc/images/traffic.png)

### 主动断开用户连接

proxy的http(s)/socks5/sps代理功能支持`控制接口`,可以通过参数`--control-url`指定的http接口地址, 那么proxy就会`定期`
的把当前连接到proxy的全部用户名或客户端IP发送到此URL，具体情况是,proxy发送一个HTTP的POST请求到`--control-url`设置的HTTP URL地址.`定期`默认是30秒，可以通过`--control-sleep`
参数修改此值。  
当用户过期，或者用户流量已经用完，通过认证API只能控制用户不能新建连接，但是已经和proxy建立当连接没法立刻断开， 那么通过控制接口可以解决这个问题，控制接口会在最慢`定期`时间内通过控制接口返回的内容，结束已经无效的用户建立的连接。

#### 控制接口请求说明

proxy会向控制接口URL发送一个HTTP POST请求，表单数据中有两个字段：user和ip。

user：当前连接到proxy的用户名，多个使用英文逗号分割，比如：user1,user2

ip：当前连接到proxy的客户端ip地址，多个使用英文逗号分割，比如：1.1.1.1,2.2.2.2

#### 控制接口返回数据说明

控制接口返回的数据是无效的用户和IP，格式是一个json对象数据，有两个字段user和ip。

比如：{"user":"a,b","ip":""}

user：当前连接到proxy的用户名，多个使用英文逗号分割，没有留空，比如：user1,user2

ip：当前连接到proxy的客户端ip地址，多个使用英文逗号分割，没有留空，比如：1.1.1.1,2.2.2.2

返回的用户和ip已经建立的连接会被proxy断开。

#### 示例

假设--control-url `http://127.0.0.1:33088/user/control.php`  指向了一个php接口地址.  
control.php内容如下：

```php  
<?php  
#接收proxy post过来的数据
$userArr=explode(",",$_POST['user']);   
$ipArr=$_GET['ip'];  

//无效用户列表
$badUsers=[]; 

foreach ($userArr as $user) {  
    //逻辑判断用户$user是否无效,如果无效就放入$badUsers
    $badUsers[]=$user;
}  
$data=["user"=>implode(","$badUsers),"ip"=>""];

echo json_encode($data);
```  

### 使用代理

http(s)/socks5/sps认证API(`--auth-url`)、控制API(`--control-url`)、流量上报API(`--traffic-url`)三个地址的访问，支持通过代理访问，参数是:`--auth-proxy`
,格式和`--jumper`格式一致，详细详细请参考`--jumper`。

## 10.本地认证，限速，控制连接数

proxy的http(s)/socks5/sps代理功能，支持通过配置文件控制用户对代理对访问，支持开启http(s)代理``Proxy Basic 代理认证`,socks5代理认证。

### 开始使用

proxy的http(s)/socks5/sps代理功能可以通过  
`--auth-file`，`--max-conns`，`--ip-limit`，`--rate-limit`，`-a`这五个参数控制。

#### 参数详细解释

##### `--auth-file`

认证的用户名和密码文件，该参数指定一个文件，一行一条规则,格式是:"用户名:密码:连接数:速率:上级"。  
`连接数`是用户最大连接数，`速率`是用户每个tcp连接最大速度，单位是：字节/秒，`上级`是用户使用的上级。  
不仅可以通过`--auth-file`设置认证用户，也可以直接用`-a`参数设置，多个用户就重复多个`-a`参数即可。  
比如：`proxy http -a a:b:0:0: -a c:d:0:0:`

实例解释：  
比如：`user:pass:100:10240:http://192.168.1.1:3100`  
`user`是认证用户名  
`pass`是认证用户密码（不能包含冒号:）  
`100`是这个用户的最大连接数，不限制写0  
`10240`是这个用户单个tcp连接的速率限制，单位是：字节/秒，不限制写0  
`http://192.168.1.1:3100`是这个用户使用的上级，没有就留空

##### `--max-conns`

限制代理服务的全局最大连接数，一个数字，0为不限制，默认0。

##### `--ip-limit`

控制客户端IP的连接数和连接速率，该参数指定一个文件，一行一条规则，#开头的是注视。  
示例文件ip.limit，规则格式如下：  
`127.0.0.1:100:10240:http://192.168.1.1:3100`  
规则解释：  
`127.0.0.1`是要限制的IP  
`100`是这个IP的最大连接数，不限制写0  
`10240`是IP单个tcp连接的速率限制，单位是：字节/s，不限制写0  
`http://192.168.1.1:3100`是这个IP使用的上级，没有就留空

##### `--rate-limit`

限制服务的每一个tcp连接的速度，比如:100K 2000K 1M . 0意味着无限制，默认0。

## 11.集群管理

proxy支持集群管理功能，此时proxy安装在每个机器节点上作为agent，配合控制面板[`proxyadmin集群版`](https://github.com/snail007/proxy-admin-cluster)，就能够实现统一管理海量机器上的proxy服务。

proxy如果要作为agent运行，假设控制面板的集群端口地址是：`1.1.1.1:55333`.

命令示例如下：

`proxy agent -k xxx -c 1.1.1.1:55333 -i test`

命令解释：

agent ： 是功能参数，表示运行agent模式。

-k ：和`proxyadmin集群版`通讯加密解密的密钥，这个密钥是在`proxyadmin集群版`的配置文件里面设置的。

-c ：`proxyadmin集群版`的集群端口地址，格式：IP:端口。

-i ： agent唯一标识，保证每个agent的都不一样，此处指定的`唯一标识`,在控制面板添加节点的时候用到，IP一栏填写这个`唯一标识`。  
如果-i不指定，默认是空，控制面板添加节点的时候IP一栏填写：agent的出口IP。

-u： 代理参数，默认是空。可以指定一个代理，agent会通过这个代理与cluster通讯。  
格式和`--jumper`格式一致，详细情况请参考手册中`--jumper`部分。

提示：

1.当控制面板里面配置了client服务，那么所有节点都使用相同都key那么就导致，只有一个client可以工作，为了解决这个问题，  
client服务参数里面可以通过使用占位符：`{AGENT_ID}`引用agent的id当作client的key，这样就可以保证每个客户端都有唯一的key。

比如client服务参数：

`client -T tcp -P 1.1.1.1:30000 --k {AGENT_ID}`

2.如果使用了代理参数`-u`,务必设置唯一标示参数`-i`,不然cluster无法区分不同的agent，agent将不能正常工作。

## 12.http、https 网站反向代理

proxy可以对 http 和 https 网站进行反向代理。

支持特性如下：

- http和https相互转换。
- 多个后端。
- 多个后端负载均衡。
- 多个后端高可用。
- 目录映射转发。
- 路径保护。
- 绑定的域名，可以另外指定别名。

使用示例，配置文件`rhttp.toml`。

```shell
proxy rhttp -c rhttp.toml
```

具体使用可以参考配置文件 [rhttp.toml](https://github.com/snail007/goproxy/blob/master/rhttp.toml) ，里面有完整配置说明。

## 13.典型用法-拨号VPS

为了方便说明，假设背景情况如下：

1、vps有一个主网卡，名称是eth0，ip是1.1.1.1，然后vps可以pppoe拨号，拨号建立的网卡名称前缀都是`pppoe_`。

实现的效果：

- 提供动态IP代理服务，客户端访问`1.1.1.1`的代理端口`7777`，可以使用拨号的ip随机选择一个作为出口，`7777`端口支持代理认证。
- 7777端口同时支持http/socks5代理。

操作步骤：

1. vps执行`proxy sps -p :7777 --bind-ip pppoe_*:7777 -a user1:password1 -a user2:password2 `。
2. 命令中`-a`是设置`代理认证用户`，多个用户，可以重复`-a`参数，格式是：`用户名:密码`。
3. 更多认证方式可以参考手册`API认证`，`认证`部分。

## 14.典型用法-多IP的VPS

为了方便说明，假设背景情况如下：

1、vps有一个主网卡，名称是`eth0`，ip是1.1.1.1，然后`eth0`配置了255个子网卡并设置了IP，比如：`eth0:1`,`eth0:255`，子网卡名称前缀都是：`eth0:`。

实现的效果：

- 提供动态IP代理服务，客户端访问`1.1.1.1`的代理端口`7777`，可以使用配置的255个子网卡的ip中随机选择一个作为出口，`7777`端口支持代理认证。
- 7777端口同时支持http/socks5代理。

操作步骤：

1. vps执行`proxy sps -p :7777 --bind-ip eth0:*:7777 -a user1:password1 -a user2:password2 `。
2. 命令中`-a`是设置`代理认证用户`，多个用户，可以重复`-a`参数，格式是：`用户名:密码`。
3. 更多认证方式可以参考手册`API认证`，`认证`部分。

## 15.典型用法-拨号的VPS集群

为了方便说明，假设背景情况如下：

1、有一批vps，它们每个都有一个主网卡配置了固定的IP：x.x.x.x，然后vps可以pppoe拨号，拨号建立的网卡名称前缀都是`pppoe_`。 

2、有一个VPS作为代理入口，它的ip是2.2.2.2。

实现的效果：

- `2.2.2.2`提供动态IP代理服务，客户端访问`2.2.2.2`的代理端口`8888`，可以随机选择拨号vps集群中一个，然后使用拨号vps的拨号ip随机选择一个作为出口ip，`8888`端口支持代理认证。
- `8888`端口同时支持http/socks5代理。

操作步骤：

拨号VPS：

1. 拨号的vps执行`proxy sps -p :7777 --bind-ip pppoe_*:7777 `。

入口VPS:

1. `2.2.2.2`
   的vps执行`proxy sps -p :8888 -a user1:password1 -a user2:password2 -P http://x.x.x.1:7777  -P http://x.x.x.2:7777`。
2. 命令中`-a`是设置`代理认证用户`，多个用户，可以重复`-a`参数，格式是：`用户名:密码`。
3. `x.x.x.1`,`x.x.x.2`是拨号vps的固定IP，有多个，重复`-P`参数即可

## 16.典型用法-对接Luminati

Luminati提供了高质量稳定的住宅IP可以做很多事情，但是价格也不菲，如果买了一个线路，可以通过proxy把它转成有子账号的代理，二次分发。

为了方便说明，假设背景情况如下：

1、购买了一个http代理线路，http代理账号信息格式是：`xxx.com:8888:username:password`，这里面包含了代理服务器地址：`xxx.com:8888`，认证用户账号和密码：`username:password`。  
2、有一个自己的VPS，IP是1.1.1.1。

实现的效果：

1. 访问1.1.1.1的7777端口，使用Luminati的线路出去。
2. 7777端口具有认证功能。

操作步骤：
1. 1.1.1.1的VPS执行：`proxy sps -p :7777 --disable-socks --disable-ss -a user1:pass1 -a user2:pass2 -P http://username:password@xxx.com:8888` 
2. 命令中`-a`是设置`代理认证用户`，多个用户，可以重复`-a`参数，格式是：`用户名:密码`。

提醒，此种做法，有可能带来Luminati限制你的账号风险，VPS的IP也有可能被可能Luminati屏蔽导致转换不能工作，风险自负。

另外上面只能把luminati的http代理转成自定义认证的二级http代理，那么如如果要把luminati的http代理转成http(s)+socks5代理可以使用参数`--luminati`。

上面步骤 `1.` 执行的命令改为：`proxy sps -p :7777 --luminati --disable-ss -a user1:pass1 -a user2:pass2 -P http://username:password@xxx.com:8888`

## 17.典型用法-认证转为无认证

通过proxy可以把有认证功能的代理转成没有认证的代理。

为了方便说明，假设背景情况如下：

1. 有一个http认证代理，地址是：2.2.2.2:8888，认证账号：abc，密码是：123。
2. 有一个自己的VPS，IP是1.1.1.1。   

实现的效果：

1. 访问1.1.1.1的7777端口，无需认证，会转给2.2.2.2:8888代理。

操作步骤：
1. 1.1.1.1的VPS执行：`proxy sps -p :7777 -P http://abc:123@2.2.2.2:8888` 

## 18.典型用法-镜像网站

有时候访问GitHub很慢，导致学习困难，可以用proxy镜像GitHub加速访问。

为了方便说明，假设背景情况如下：

1. 有一个自己的VPS，IP是1.1.1.1，可以访问GitHub，而且速度还可以。

实现的效果：

1. 访问1.1.1.1的7777端口，就是访问GitHub。

操作步骤：
1. 准备配置文件github.toml，写入如下内容：
```toml
[[host]]
bind="http://1.1.1.1:7777/"
target="https://github.com/"
upstream="github.com:443"
timeout=5000
```
1. 在1.1.1.1的VPS执行：`proxy rhttp -c github.toml`  
2. 访问 `http://1.1.1.1:7777/snail007/` ，如果一切正常，就把命令加上后台运行参数即可。
3. 命令改成：`proxy rhttp -c github.toml --daemon --log /tmp/github.log`  
4. 此功能需要`proxy`版本大于等于`v11.2`。

如果是想使用域名比如a.com，那么：
1. 把a.com解析到IP 1.1.1.1。
2. 把`bind="http://1.1.1.1:7777/"` 改成 `bind="http://a.com:7777/"`
3. 访问 `http://a.com:7777/snail007/` 即可。
4. 如果想使用默认80端口，配置bind为 `bind="http://a.com/"` 。
5. 访问 `http://a.com/snail007/` 即可。