---
title: "Cloudflare使用简明教程"
date: 2019-06-28T17:37:30+08:00
draft: false
description: "Cloudflare使用简明教程"
tags: [ "CDN" ]
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
#### 前言

GOPROXY底层传输协议支持ws，可以利用cdn加速，使用cdn需要域名解析知识以及cdn的一些知识，那么很多小伙伴对cdn不是很熟悉，为了更好的使用GOPROXY，就找了几个高质量的怎么使用cloudflare免费套餐的文章。  

#### 教程

1.不需要修改域名DNS，可以直接使用二级域名加速：https://www.geekzu.cn/archives/free-cname-cloudflare-with-ssl.html  

2.需要修改域名DNS，使用一级域名加速：https://ask.dobunkan.com/q-26886.html  