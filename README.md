<img src="https://github.com/snail007/goproxy/blob/master/docs/images/logo.jpg?raw=true" width="200"/>      
Proxy是golang实现的高性能http,https,websocket,tcp代理服务器.   
<hr/>

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https://img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total.svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style=plastic)](https://github.com/snail007/goproxy/releases)

**1.开始代理,使用本地端口38080.**  
`./proxy -p 38080`  
**2.作为二级代理,使用本地端口8090,假设上级http代理是`22.22.22.22:8080`**  
`./proxy -P "22.22.22.22:8080" -p 8090`  
**3.生成加密通讯需要的证书和密钥**  
在linux上并安装了openssl命令，可以直接通过下面的命令生成证书和key文件.   
`./proxy keygen`   
会在当前目录下面生成一个证书文件proxy.crt和key文件proxy.key。  
**4.一级和二级代理加密通讯.**  
一级代理(VPS,IP:22.22.22.22)   
`./proxy -x -p 38080 -c proxy.crt -k proxy.key`  
二级代理(本地)  
`./proxy -X -P "22.22.22.22:38080" -p 8080 -c proxy.crt -k proxy.key`  
那么访问本地的8080端口就是访问VPS上面的代理端口38080.  
**5.TCP隧道**  
VPS,IP:22.22.22.22  
`./proxy --tcp -P "127.0.0.1:22" -x -p 22022 -c proxy.crt -k proxy.key`  
本地  
`./proxy --tcp -X -P "22.22.22.22:22022" -p 12222 -c proxy.crt -k proxy.key`  
那么服务本地的12222端口就是通过加密隧道安全的访问vps的本地22端口.  

[图文教程](docs/faststart.md)

### Features
- 程序本身可以作为一级代理,如果设置了上级代理那么可以作为二级代理,乃至N级代理.
- 如果程序不是一级代理,而且上级代理也是本程序,那么可以加密和上级代理之间的通讯,采用底层tls高强度加密,安全无特征.
- 代理时会自动判断访问的网站是否屏蔽,如果被屏蔽那么就会使用上级代理(前提是配置了上级代理)访问网站;如果访问的网站没有被屏蔽,为了加速访问,代理会直接访问网站,不使用上级代理.
- 可以设置域名黑白名单，更加自由的控制网站的访问方式。
- 跨平台性,无论你是widows,linux,还是mac,甚至是树莓派,都可以很好的运行proxy.  

### Why need these?
当由于安全因素或者限制,我们不能顺畅的访问我们在其它地方的服务,我们可以通过多个相连的proxy节点建立起一个安全的隧道,顺畅的访问我们的服务.

### Fast Start
提示:所有操作需要root权限.  
**0.如果你的VPS是linux64位的系统,那么只需要执行下面一句,就可以完成自动安装和配置.**   
```shell
curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash
```
安装完成,配置目录是/etc/proxy,更详细的使用方法参考下面的进一步了解.  
如果你的vps不是linux64位系统,请按照下面的半自动步骤安装:  

**1.登录你的VPS,下载守护进程monexec,选择合适你的版本,vps一般选择"linux_amd64.tar.gz"的即可.**      
下载地址:https://github.com/reddec/monexec/releases   
比如下载到/root/proxy/  
执行:  
```shell
mkdir /root/proxy/  
cd /root/proxy/  
wget https://github.com/reddec/monexec/releases/download/v0.1.1/monexec_0.1.1_linux_amd64.tar.gz   
```
**2.下载proxy**  
下载地址:https://github.com/snail007/goproxy/releases   
```shell
cd /root/proxy/  
wget https://github.com/snail007/goproxy/releases/download/v2.0/proxy-linux-amd64.tar.gz    
```
**3.下载自动安装脚本** 
```shell
cd /root/proxy/   
wget https://raw.githubusercontent.com/snail007/goproxy/master/install.sh
chmod +x install.sh   
./install.sh   
```
### More...
**1、作为普通一级代理。**   
默认监听0.0.0.0:33080端口，可以使用-p修改端口，-i修改绑定ip。  
默认情况  
`./proxy`  
指定ip和端口  
`./proxy  -i 192.168.1.100 -p 60080`  

**2、作为普通二级代理。**  
可以通过-P指定上级代理，格式是IP:端口  
`./proxy -P "192.168.1.100:60080" -p 33080`   

**3、作为加密一级代理。**  
加密模式的一级代理需要和加密的二级代理配合。  
加密模式需要证书和key文件，在linux上并安装了openssl命令，可以直接通过下面的命令生成证书和key文件。  
`./proxy keygen`  
会在当前目录下面生成一个证书文件proxy.crt和key文件proxy.key。  
比如在你的vps上运行加密一级代理，使用参数-x即可，默认会使用程序相同目录下面的证书文件proxy.crt和key文件proxy.key。  
`./proxy -x`   
或者使用-c和-k指定证书和key文件,ip和端口。   
`./proxy -x -c "proxy.crt" -k "proxy.key" -p 58080`   

**4、作为加密二级代理。**  
加密模式的二级代理需要和加密的一级代理配合。加密模式的二级代理和加密模式的一级代理要使用相同的证书和key文件。  
默认会使用程序相同目录下面的证书文件proxy.crt和key文件proxy.key。    
比如在你的windows电脑上允许二级加密代理,需要-P指定上级代理，同时设置-X代表是加密的上级代理。   
假设一级代理vps外网IP是：115.34.9.63。    
`./proxy.exe -X -P "115.34.9.63:58080" -c "proxy.crt" -k "proxy.key"  -p 18080`     
然后设置你的windos系统中，需要通过代理上网的程序的代理为http模式，地址为：127.0.0.1，端口为：18080，    
然后程序即可通过加密通道通过vps上网。 
### License 
Proxy is licensed under GPLv3 license.
### Contact 
QQ交流群:189618940



