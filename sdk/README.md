
# Proxy SDK 使用说明  

支持以下平台:  
- Android,`.arr`库
- IOS,`.framework`库
- Windows,`.dll`库
- Linux,`.so`库

proxy使用gombile实现了一份go代码编译为android和ios平台下面可以直接调用的sdk类库, 
另外还为linux和windows提供sdk支持，基于这些类库,APP开发者可以轻松的开发出各种形式的代理工具。    

# 下面分平台介绍SDK的用法  

## Android SDK
  
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy-sdk-android/) [![license](https://img.shields.io/github/license/snail007/goproxy-sdk-android.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy-sdk-android/total.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-android/releases) [![download](https://img.shields.io/github/release/snail007/goproxy-sdk-android.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-android/releases)  
  
[点击下载Android-SDK](https://github.com/snail007/goproxy-sdk-android/releases)  
在Android系统提供的sdk形式是一个后缀为.aar的类库文件,开发的时候只需要把arr类库文件引入android项目即可.  

### Android-SDK使用实例

#### 1.导入包

```java
import snail007.proxy.Porxy
```

#### 2.启动一个服务

```java
String serviceID="http01";//这里serviceID是自定义的唯一标识字符串,保证每个启动的服务不一样即可
String serviceArgs="http -p :8080";
String err=Proxy.start(serviceID,serviceArgs);
if (!err.isEmpty()){
    //启动失败
    System.out.println("start fail,error:"+err);
}else{
    //启动成功
}
```
#### 3.判断一个服务是否在运行

```java
String serviceID="http01";
boolean isRunning=Proxy.isRunning(serviceID);
if(isRunning){
    //正在运行
}else{
    //没有运行
}
```

由于tclient和client服务的特性,目前这个方法对于服务tclient和client永远返回false.  

#### 4.停止一个服务

```java
String serviceID="http01";
Proxy.stop(serviceID);
//停止完毕

```

## IOS SDK
  
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy-sdk-ios/) [![license](https://img.shields.io/github/license/snail007/goproxy-sdk-ios.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy-sdk-ios/total.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-ios/releases) [![download](https://img.shields.io/github/release/snail007/goproxy-sdk-ios.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-ios/releases)  
  
[点击下载IOS-SDK](https://github.com/snail007/goproxy-sdk-ios/releases)  
在IOS系统提供的sdk形式是一个后缀为.framework的类库文件夹,开发的时候只需要把类库文件引入项目,然后调用方法即可.  

### IOS-SDK使用实例

#### 导入包

```objc
#import <Proxy/Proxy.objc.h>
```

#### 2.启动一个服务

```objc
-(IBAction)doStart:(id)sender
{
	//这里serviceID是自定义的唯一标识字符串,保证每个启动的服务不一样
	NSString *serviceID = @"http01";
    NSString *serviceArgs = @"http -p :8080";
    NSString *error = ProxyStart(serviceID,serviceArgs);
    
    if (error != nil && error.length > 0)
    {
        NSLog(@"start error %@",error);
    }else{
        NSLog(@"启动成功");
    }
}
```

#### 3.判断一个服务是否在运行

```objc
-(IBAction)hasRunning:(id)sender;
{
    NSString *serviceID = @"http01";
    if (ProxyIsRunning(serviceID))
    {
        NSLog(@"正在运行");
    }else{
        NSLog(@"没有运行");
    }
}
```

由于tclient和client服务的特性,目前这个方法对于服务tclient和client永远返回false.  

#### 4.停止一个服务

```objc
-(IBAction)doStop:(id)sender
{
    NSString *serviceID = @"http01";
    ProxyStop(serviceID);
    //停止完毕
}
```


## Windows SDK
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy-sdk-windows/) [![license](https://img.shields.io/github/license/snail007/goproxy-sdk-windows.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy-sdk-windows/total.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-windows/releases) [![download](https://img.shields.io/github/release/snail007/goproxy-sdk-windows.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-windows/releases)  
  
[点击下载Windows-SDK](https://github.com/snail007/goproxy-sdk-windows/releases)  
在Windows系统提供的sdk形式是一个后缀为.dll的类库文件,开发的时候只需要把dll类库文件加载,然后调用方法即可.  

### Windows-SDK使用实例  
C++示例，不需要包含头文件，只需要加载proxy-sdk.dll即可，ieshims.dll需要和proxy-sdk.dll在一起。  
作者：[yjbdsky](https://github.com/yjbdsky)     

```cpp
#include <stdio.h>
#include<stdlib.h>
#include <string.h>
#include<pthread.h>
#include<Windows.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef char *(*GOSTART)(char *s);
typedef char *(*GOSTOP)(char *s);
typedef int(*GOISRUN)(char *s);
HMODULE GODLL = LoadLibrary("proxy-sdk.dll");

char * Start(char * p0,char * p1)
{
	if (GODLL != NULL)
	{
		GOSTART gostart = *(GOSTART)(GetProcAddress(GODLL, "Start"));
		if (gostart != NULL){
			printf("%s:%s\n",p0, p1);
			char *ret = gostart(p0,p1);
			return ret;
		}
	}
	return "Cannot Find dll";
}
char * Stop(char * p)
{
	if (GODLL != NULL)
	{
		GOSTOP gostop = *(GOSTOP)(GetProcAddress(GODLL, "Stop"));
		if (gostop != NULL){
			printf("%s\n", p);
			char *ret = gostop(p);
			return ret;
		}
	}
	return "Cannot Find dll";
}

int IsRunning(char * p)
{
	
	if (GODLL != NULL)
	{
		GOISRUN isrun = *(GOISRUN)(GetProcAddress(GODLL, "IsRunning"));
		if (isrun != NULL){
			int ret = isrun(p);
			return ret;
		}
		FreeLibrary(GODLL);
	}
	return 0;
}

int main()
{
	//这里p0是自定义的唯一标识字符串,保证每个启动的服务不一样
	char *p0 = "http01";
	char *p1 = "http -t tcp -p :38080";
	printf("This is demo application.\n");
	//启动服务,返回空字符串说明启动成功;返回非空字符串说明启动失败,返回的字符串是错误原因
	printf("start result %s\n", Start(p0,p1));
	//停止服务,没有返回值
	Stop(p0);
	//服务是否在运行,返回0是没有运行,返回1正在运行
	printf("is running result %d\n", IsRunning(p0));
	return 0;
}


#ifdef __cplusplus
}
#endif
```

C++示例2，请移步：[GoProxyForC](https://github.com/SuperPowerLF2/GoProxyForC)   

## Linux SDK
[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy-sdk-linux/) [![license](https://img.shields.io/github/license/snail007/goproxy-sdk-linux.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy-sdk-linux/total.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-linux/releases) [![download](https://img.shields.io/github/release/snail007/goproxy-sdk-linux.svg?style=plastic)](https://github.com/snail007/goproxy-sdk-linux/releases)  
  
[点击下载Linux-SDK](https://github.com/snail007/goproxy-sdk-linux/releases)  
在Linux系统提供的sdk形式是一个后缀为.so的类库文件,开发的时候只需要把so类库加载,调用方法即可.  

### Linux-SDK使用实例
Linux下面使用的sdk是so文件即proxy-sdk.so,下面写一个简单的C程序示例,调用so库里面的方法.  

`vi test-proxy.c`  

```c
#include <stdio.h>
#include "proxy-sdk.h"

int main() {
     printf("This is demo application.\n");
	 //这里p0是自定义的唯一标识字符串,保证每个启动的服务不一样
	 char *p0 = "http01";
     char *p1 = "http -t tcp -p :38080";
     //启动服务,返回空字符串说明启动成功;返回非空字符串说明启动失败,返回的字符串是错误原因
     printf("start result %s\n",Start(p0,p1));
     //停止服务,没有返回值
     Stop(p0);
     //服务是否在运行,返回0是没有运行,返回1正在运行
     printf("is running result %d\n",IsRunning(p0));
     return 0;
}
```

#### 编译test-proxy.c ####  
`export LD_LIBRARY_PATH=./ && gcc -o test-proxy test.c proxy-sdk.so`  

#### 执行 ####  
`./test-proxy`  


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
sdk里面的服务不支持手册里面的：--daemon和--forever参数.  


