这里以vps centos 64位为例子  
Linux 部分  
1.Putty工具（或其他工具）  
root登入  
2.下载批量命令文件install_auto.sh（64位的话直接执行这个命令即刻）  
#curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash  
注意  
这里的install_auto.sh 源码可以下载修改对面版本proxy保存后执行批量命令  
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image001.png?raw=true"/>  
3.修改/etc/proxy/proxy.toml配置文件   
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image002.png?raw=true"/>
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image003.png?raw=true"/>
<img src="https://github.com/snail007/goproxy/blob/master/docs/images/image004.png?raw=true"/>  
#/usr/bin/proxyd status  
如果未运行那么执行：/usr/bin/proxy   
