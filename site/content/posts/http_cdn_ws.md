---
title: "代理之套用CDN"
date: 2019-06-14T16:25:06+08:00
draft: false
description: "http代理前置CDN，保护你的后端代理"
tags: [ "cdn","ws" ]
categories: [ "默认分类" ]
isCJKLanguage: true
weight: 0

# keywords: [ "keyword" ]
# linkTitle: ""

# 使用这两个参数将会重置permalink，默认使用文件名
#url: 
#slug: ""

# 设置文章的过期时间，如果是已过期的文章则不会发布，除非使用 --buildExpired 参数
# expiryDate: 2020-01-01

# 设置文章的发布时间，如果是未来的时间则不会发布，除非使用 --buildFuture 参数
# publishDate: 2020-01-01

# 别名将通过重定向实现
# aliases:
#   - 别名1
#   - /posts/my-original-url/
#   - /2010/01/01/another-url.html

# type 与 layout 参数将会改变 Hugo 寻找该文章模板的顺序
# type: review
# layout: reviewarticle
---

### goproxy代理之套用CDN

#### 准备	
1. 在VPS上[下载](https://github.com/snail007/goproxy/releases)对应平台的执行文件，这里以Linux为例。						
2. 在电脑上[下载](https://github.com/snail007/goproxy/releases)对应平台的执行文件，这里以Windows为例。						

#### 1. proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http代理  

1.1. VPS上的命令 ./proxy http -t ws -p "0.0.0.0:8080" 

1.2 Cloudflare上的设置 Crypto -> SSL ->设置为Off	 

1.3 windows计算机上的命令 proxy.exe http -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080"	 

1.4 计算机上的代理设置 "127.0.0.1:8090 http"	 
  

#### 2. proxy本地http代理 <--ws\wss--> CDN  <--ws\wss--> VPS的proxy提供的http代理  

2.1 VPS上的命令 ./proxy http -t wss -p "0.0.0.0:8443"	 

2.2 Cloudflare上的设置 Crypto -> SSL ->设置为Full 

2.3 windows计算机上的命令 proxy.exe http -t tcp -p "0.0.0.0:8090" -T wss -P "your.domain.com:8443"	  

2.4 计算机上的代理设置 "127.0.0.1:8090 http"	 
  

#### 3. proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http+Basic认证 代理 （成功） 

3.1 VPS上的命令 ./proxy http -t ws -p "0.0.0.0:8080" -a "user:pass" 	 

3.2 Cloudflare上的设置 Crypto -> SSL ->设置为Off 

3.3 windows计算机上的命令 proxy.exe http -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080"	 

3.4 计算机上的代理设置 "127.0.0.1:8090 http"	 
   

#### 4. proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http+加密方式+密码 代理 （成功） 

4.1 VPS上的命令 ./proxy http -t wss -p "0.0.0.0:8443" --parent-ws-method="chacha20-ietf-poly1305" --parent-ws-password="snail007/goproxy" 

4.2 Cloudflare上的设置 Crypto -> SSL ->设置为Full 

4.3 windows计算机上的命令 proxy.exe http -t tcp -p "0.0.0.0:8090" -T wss -P "your.domain.com:8443" --local-ws-method="chacha20-ietf-poly1305" --local-ws-password="snail007/goproxy"	 

4.4 计算机上的代理设置 "127.0.0.1:8090 http"	 

4.5 相关参数介绍 "--local-ws-method:加密方式--parent-ws-password:设置密码"	

#### 5. proxy本地socks5代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的socks5代理 

5.1 VPS上的命令 ./proxy socks -t ws -p "0.0.0.0:8080" 

5.2 Cloudflare上的设置 Crypto -> SSL ->设置为Off 

5.3 windows计算机上的命令 proxy.exe socks --always -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080" 

5.4 计算机上的代理设置 "127.0.0.1:8090 socks5 [Remote DNS]"
 


#### 备注
1、Cloudflare支持的回源端口  
```text
HTTP 端口为：
            80	 
            8080	 
            8880	 
            2052	 
            2082	 
            2086	 
            2095	
HTTPs 端口为：
            443
            2053
            2083
            2087
            2096
            8443	
```
