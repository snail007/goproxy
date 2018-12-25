# Pass-through user IP manual

## Description:

By Linux TPROXY function,you can achieve the source Station service program can see the client's real IP, to achieve this feature requires linux operating systems and programs must meet certain conditions.

## Environmental requirements:

The source station must be a service program running on Linux, and Linux needs to meet the following conditions:

1. Linux kernel version >= 2.6.28

2. Determine whether the system supports TPROXY, execute:

    grep TPROXY /boot/config-`uname -r`

    If the output has the following result description is supported.

    CONFIG_NETFILTER_XT_TARGET_TPROXY=m

## Deployment steps:

1. The tproxy environment setup script should be executed with root privileges every time the boot from the source Linux system: tproxy_setup.sh

2. Execute proxy proxy with root access on the source Linux system

## Parameter-tproxy is the tproxy function that turns on the proxy.

./proxy -tproxy

2. The IP address of the source station to listen to the program requires the use of: `127.0.1.1`

For example, the address of the source station before listening is: `0.0.0.0:8800`, now need to be modified to: `127.0.1.1:8800`

3. Forwarding rules inside the source address must be the corresponding, such as the above: `127.0.1.1:8800`
