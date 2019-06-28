---
title: "内网穿透也能用CDN？？？！"
date: 2019-06-28T17:02:57+08:00
draft: false
description: "内网穿透通过CDN服务http服务"
tags: [ "内网穿透","CND" ]
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

#### 内网穿透之套CDN
**好处就不用说了吧，两个字稳！**   

假如在windows上配置了可以访问的连接 http://127.0.0.1:801/abc/index.php , 如何能使用VPS+CDN做内网穿透呢？  
配置如下：   
1. 在VPS上执行（配置无误后，再加上--daemon）  
./proxy bridge -t ws -p ":8080" -C proxy.crt -K proxy.key  
./proxy server -T ws -r ":8880@:801" -P "127.0.0.1:8080" -C proxy.crt -K proxy.key  
2. 在windows上执行  
proxy.exe client -T ws -P "your.domain.com:8080" -C proxy.crt -K proxy.key  
3. 这时在浏览器里访问：http://your.domain.com:8880/abc/index.php ，就会看到电脑中的页面了。  