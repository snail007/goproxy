
# Proxy SDK 使用说明  

proxy使用gombile实现了一份go代码编译为android和ios平台下面可以直接调用的sdk类库,  
基于这些类库,APP开发者可以轻松的开发出各种形式的代理工具.  

### 下面分平台介绍SDK的用法  

#### Android SDK
  
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https://img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total.svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style=plastic)](https://github.com/snail007/goproxy/releases)  
  
[点击下载Android-SDK](https://github.com/snail007/goproxy-sdk-android/releases)  
在Android系统提供的sdk形式是一个后缀为.aar的类库文件,开发的时候只需要把arr类库文件引入android项目即可.  

### Android-SDK使用实例

#### 1.导入包
```java
import snail007.proxy.Porxy
```

#### 2.启动一个服务
```java
String args="http -p :8080";
String err=Proxy.start(args);
if (err.isEmpty()){
    //启动失败
    System.out.println("start fail,error:"+err);
}else{
    //启动成功
}
```
#### 3.判断一个服务是否在运行

```java
String args="http -p :8080";
boolean isRunning=Proxy.isRunning(args);//这里传递http也可以,最终使用的就是args里面的第一个参数http
if(isRunning){
    //正在运行
}else{
    //没有运行
}
```
#### 4.停止一个服务

```java
String args="http -p :8080";
Proxy.stop(args);//这里传递http也可以,最终使用的就是args里面的第一个参数http
//停止完毕

```


### IOS-SDK使用实例

#### todo


### 关于服务  
proxy的服务有11种,分别是:  

```shell
http  
socks  
sps  
tcp  
udp  
bridge  
server  
client  
tbridge  
tserver  
tclient  
```
每个服务只能启动一个,如果相同的服务启动多次,那么之前的服务会被停掉,后面启动的服务覆盖之前的服务.  
上面这些服务的具体使用方式和具体参数,可以参考[proxy手册](https://github.com/snail007/goproxy/blob/master/README_ZH.md)  