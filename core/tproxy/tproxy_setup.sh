#!/bin/bash
SOURCE_BIND_IP="127.0.1.1"

echo 0 > /proc/sys/net/ipv4/conf/lo/rp_filter
echo 2 > /proc/sys/net/ipv4/conf/default/rp_filter
echo 2 > /proc/sys/net/ipv4/conf/all/rp_filter
echo 1 > /proc/sys/net/ipv4/conf/all/send_redirects
echo 1 > /proc/sys/net/ipv4/conf/all/forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# 本地的话,貌似这段不需要
# iptables -t mangle -N DIVERT >/dev/null 2>&1
# iptables -t mangle -F DIVERT
# iptables -t mangle -D PREROUTING -p tcp -m socket -j DIVERT >/dev/null 2>&1
# iptables -t mangle -A PREROUTING -p tcp -m socket -j DIVERT
# iptables -t mangle -A DIVERT -j MARK --set-mark 1
# iptables -t mangle -A DIVERT -j ACCEPT

ip rule del fwmark 1 lookup 100
ip rule add fwmark 1 lookup 100
ip route del local 0.0.0.0/0 dev lo table 100
ip route add local 0.0.0.0/0 dev lo table 100

ip rule del from ${SOURCE_BIND_IP} table 101
ip rule add from ${SOURCE_BIND_IP} table 101
ip route del default via 127.0.0.1 dev lo table 101
ip route add default via 127.0.0.1 dev lo table 101

ip route flush cache
ip ro flush cache