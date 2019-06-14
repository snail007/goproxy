---
title: "常见问题解答"
date: 2019-06-14T16:29:13+08:00
draft: false
description: "goproxy常见问答"
categories: [ "默认分类" ]
isCJKLanguage: true
weight: 10000

# keywords: [ "keyword" ]
# linkTitle: ""

# 使用这两个参数将会重置permalink，默认使用文件名
#url: "faq/" 
slug: "faq/goproxy常见问题解答"

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

### 问：http代理支持https网站吗？
答：goproxy的http代理同时支持http和https网站。

### 问：socks5代理支持UDP吗？
答：goproxy的socks5代理是完全按着[rfc1982](https://tools.ietf.org/html/rfc1928)草案实现的，支持UDP功能，草案规定UDP端口是在代理握手请求过程中服务端动态指定的，不需要提前指定。如果你使用的socks5客户端是标准客户端，是不需要你手动指定UDP端口的。  

### 问：什么是http代理？https代理？
答：无论是`http代理`还是`https代理`，都支持同时代理访问`http`和`https`网站。  
`http代理`和`https代理`里的http和https和你访问的网站是否是http（https）是无关的。  
`http代理`：客户端和代理服务器之间是`tcp`传输数据。  
`https代理`：客户端和代理服务器之间是`tls`加密传输数据。  
http代理是使用最广泛的代理，大部分客户端都不支持https代理。  

### 问：goproxy的http代理是http代理还是https代理？
答：goproxy支持http代理，https代理，https双向认证代理。  
至于提供什么类型的代理，是参数`-t`决定的：  
http代理 ：`-t tcp`    
https代理：`-t tls --local-tls-single`，服务端需要设置tls证书。  
https双向认证代理：`-t tls`，服务端和客户端都需要设置一样的tls证书。    
