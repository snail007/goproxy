这里以vps centos 64位为例子  
Linux 部分  
1.Putty工具（或其他工具）  
root登入  
2.下载批量命令文件install_auto.sh（64位的话直接执行这个命令即可）  
#curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash  
注意  
这里的install_auto.sh 源码可以下载修改proxy版本,保存后执行.  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image001.png?raw=true"/>  
3.修改/etc/proxy/proxy.toml配置文件   
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image002.png?raw=true"/>
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image003.png?raw=true"/>
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image004.png?raw=true"/>  
#/usr/bin/proxyd status  
如果未运行那么执行调试命令：/usr/bin/proxy   
如果一切正常,可以使用proxyd命令管理proxy,执行 proxyd 可以查看用法.
后台启动proxy: proxyd start
4.下载证书加密文件/etc/proxy/proxy.crt和/etc/proxy/proxy.key到windows  
Windows部分  
5.https://github.com/snail007/goproxy/releases 下载对应windows版本   
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image005.jpg?raw=true"/>  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image006.png?raw=true"/>  
我的是d：盘  
6.修改windows下的proxy.toml  vps服务ip和上面设置的端口哦  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image007.png?raw=true"/>  
然后运行proxy.exe即可.  
这时候浏览器代理服务器就是127.0.0.1:9501啦,完毕!  

要隐藏windows命令用工具下载RunHiddenConsole.exe 写个bat文件都放proxy目录下就行
Start.bat  

@echo off  
echo Starting  
RunHiddenConsole D:/proxy/proxy.exe  

Stop.bat  
@echo off  
echo Stopping  
taskkill /F /IM proxy.exe > nul  
exit  
