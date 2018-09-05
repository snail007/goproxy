# 透传用户IP手册

说明:

通过Linux的TPROXY功能,可以实现源站服务程序可以看见客户端真实IP,实现该功能需要linux操作系统和程序都要满足一定的条件.

环境要求:

源站必须是运行在Linux上面的服务程序,同时Linux需要满足下面条件:

1.Linux内核版本 >= 2.6.28

2.判断系统是否支持TPROXY,执行:

    grep TPROXY /boot/config-`uname -r`

    如果输出有下面的结果说明支持.

    CONFIG_NETFILTER_XT_TARGET_TPROXY=m

部署步骤:

1.在源站的linux系统里面每次开机启动都要用root权限执行tproxy环境设置脚本:tproxy_setup.sh

2.在源站的linux系统里面使用root权限执行代理proxy

参数 -tproxy 是开启代理的tproxy功能.

./proxy -tproxy

2.源站的程序监听的地址IP需要使用:127.0.1.1

比如源站以前监听的地址是: 0.0.0.0:8800 , 现在需要修改为:127.0.1.1:8800

3.转发规则里面源站地址必须是对应的,比如上面的:127.0.1.1:8800
