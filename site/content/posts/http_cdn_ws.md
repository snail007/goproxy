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
1. 在VPS上下载对应平台的执行文件，这里以Linux为例。						
2. 在电脑上下载对应平台的执行文件，这里以Windows为例。						

功能配置简介 | VPS上的命令 | Cloudflare上的设置 | windows计算机上的命令 | 计算机上的代理设置 | 相关参数介绍 
:--- | :--- | :--- | :--- | :--- | :--- | :---			
proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http代理  | 	./proxy http -t ws -p "0.0.0.0:8080" | 	Crypto -> SSL ->设置为Off	 | proxy.exe http -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080"	 | "127.0.0.1:8090 http"	 |  
proxy本地http代理 <--ws\wss--> CDN  <--ws\wss--> VPS的proxy提供的http代理  |	./proxy http -t wss -p "0.0.0.0:8443"	 | Crypto -> SSL ->设置为Full | 	proxy.exe http -t tcp -p "0.0.0.0:8090" -T wss -P "your.domain.com:8443"	  | "127.0.0.1:8090 http"	 |  
proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http+Basic认证 代理 （成功） | 	./proxy http -t ws -p "0.0.0.0:8080" -a "user:pass" 	 | Crypto -> SSL ->设置为Off | 	proxy.exe http -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080"	 | "127.0.0.1:8090 http"	 |   
proxy本地http代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的http+加密方式+密码 代理 （成功） | 	./proxy http -t wss -p "0.0.0.0:8443" --parent-ws-method="chacha20-ietf-poly1305" --parent-ws-password="snail007/goproxy" | 	Crypto -> SSL ->设置为Full | 	proxy.exe http -t tcp -p "0.0.0.0:8090" -T wss -P "your.domain.com:8443" --local-ws-method="chacha20-ietf-poly1305" --local-ws-password="snail007/goproxy"	 | "127.0.0.1:8090 http"	 | "--local-ws-method:加密方式--parent-ws-password:设置密码"	
proxy本地socks5代理 <--ws\wss--> CDN  <--ws\wss-->  VPS的proxy提供的socks5代理 | 	./proxy socks -t ws -p "0.0.0.0:8080" | 	Crypto -> SSL ->设置为Off | 	proxy.exe socks --always -t tcp -p "0.0.0.0:8090" -T ws -P "your.domain.com:8080" | 	"127.0.0.1:8090 socks5 [Remote DNS]"

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
