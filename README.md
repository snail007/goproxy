## GOPROXY Introduction
<img src="https://github.com/snail007/goproxy/blob/master/doc/images/logo.jpg?raw=true" width="200"/>
Proxy is a high-performance http, https, websocket, tcp, udp, socks5, ss proxy server implemented by golang, supporting forward proxy, reverse proxy, transparent proxy, intranet penetration, TCP/UDP port mapping, SSH relay, TLS Encrypted transmission, protocol conversion, anti-pollution DNS proxy.

[Click to download] (https://github.com/snail007/goproxy/releases) Official QQ exchange group: 42805407

[Free version VS commercial version] (https://snail007.github.io/goproxy/free_vs_commercial/)

---

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/snail007/goproxy/) [![license](https:/ /img.shields.io/github/license/snail007/goproxy.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/snail007/goproxy/total. Svg?style=plastic)](https://github.com/snail007/goproxy/releases) [![download](https://img.shields.io/github/release/snail007/goproxy.svg?style= Plastic)](https://github.com/snail007/goproxy/releases)
  
- [English Manual](/README.md)
- [GORPOXY Help Manual] (https://snail007.github.io/goproxy/manual/)
- [GORPOXY Practical Tutorial] (https://snail007.github.io/goproxy)
- [Desktop Edition] (/gui/README_ZH.md)
- [SDK](https://github.com/snail007/goproxy-sdk)

### What can it do?
- Chained agent, the program itself can be used as an agent, and if it is set up, it can be used as a secondary agent or even an N-level agent.
- Communication encryption, if the program is not a level one agent, and the upper level agent is also the program, then the communication between the upper level agent and the upper level agent can be encrypted, and the underlying tls high-intensity encryption is used, and the security is featureless.
- Smart HTTP, SOCKS5 proxy, will automatically determine whether the visited website is blocked. If it is blocked, it will use the superior agent (provided that the superior agent is configured) to access the website; if the visited website is not blocked, in order to speed up the access, the agent will Direct access to the website without using a superior agent.
- Domain name black and white list, more free to control the way the website is accessed.
- Cross-platform, whether you are widows, linux, mac, or even raspberry pie, you can run the proxy very well.
- Multi-protocol support, support for HTTP(S), TCP, UDP, Websocket, SOCKS5 proxy.
- TCP/UDP port forwarding.
- Support intranet penetration, protocol supports TCP and UDP.
- SSH relay, HTTP (S), SOCKS5 proxy supports SSH relay, the upper Linux server does not need any server, a local proxy can be happy online.
- [KCP] (https://github.com/xtaci/kcp-go) protocol support, HTTP(S), SOCKS5 proxy supports KCP protocol to transmit data, reduce latency and improve browsing experience.
- Integrated external API, HTTP(S), SOCKS5 proxy authentication function can be integrated with external HTTP API, which can easily control proxy users through external systems.
- Reverse proxy, which supports direct parsing of the domain name to the IP of the proxy listener, and then the proxy will help you access the HTTP(S) website that needs to be accessed.
- Transparent HTTP (S) proxy, in conjunction with iptables, forwards the outgoing 80, 443 traffic directly to the proxy at the gateway, enabling non-aware intelligent router proxy.
- Protocol conversion, which can convert existing HTTP(S) or SOCKS5 or SS proxy into one port and support HTTP(S) and SOCKS5 and SS proxy at the same time. Converted SOCKS5 and SS proxy. If the superior is SOCKS5 proxy, then UDP is supported. Features while supporting powerful cascading authentication.
- Custom underlying encrypted transmission, http(s)\sps\socks proxy can encrypt tcp data via tls standard encryption and kcp protocol on top of tcp, in addition to support custom encryption after tls and kcp, that is Said custom encryption and tls|kcp can be used in combination, the internal AES256 encryption, you only need to define a password when you use it.
- Underlying compression efficient transmission, http(s)\sps\socks proxy can encrypt tcp data through custom encryption and tls standard encryption and kcp protocol on tcp, and can also compress data after encryption, that is, compression function And custom encryption and tls|kcp can be used in combination.
- Secure DNS proxy, which can secure and prevent pollution DNS queries through encrypted proxy communication between the DNS proxy server provided by the local proxy and the superior proxy.
- Load balancing, high availability, HTTP(S)\SOCKS5\SPS agent supports superior load balancing and high availability, and multiple superior repeat-P parameters can be used.
- Specify the egress IP. The HTTP(S)\SOCKS5\SPS proxy supports the client to connect with the ingress IP, and uses the ingress IP as the egress IP to access the target website. If the ingress IP is an intranet IP, the egress IP does not use the ingress IP.
- Support speed limit, HTTP(S)\SOCKS5\SPS proxy supports speed limit.
- SOCKS5 agent supports cascading certification.
- The certificate parameter uses base64 data. By default, the -C, -K parameter is the path of the crt certificate and the key file. If it is the beginning of base64://, then the latter data is considered to be base64 encoded and will be used after decoding.


### Why do you need it?

- When for some reason we are unable to access our services elsewhere, we can establish a secure tunnel to access our services through multiple connected proxy nodes.
- WeChat interface is developed locally for easy debugging.
- Remote access to intranet machines.
- Play LAN games with your friends.
- I used to play only on the LAN, and now I can play anywhere.
- Replace the sword inside Netnet, show IP internal Netcom, peanut shell and other tools.
- ..

 
The manual on this page applies to the latest version of goproxy. Other versions may not be applicable. Please use the command according to your own instructions.
 

### Joining the organization
[Click to join the gitter] (https://gitter.im/go-proxy/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)

[Click to join the TG] (https://t.me/snail007_goproxy)

## Download and install 

### Quick installation

0. If your VPS is a Linux 64-bit system, you only need to execute the following sentence to complete the automatic installation and configuration.

Tip: All operations require root privileges.

The free version performs this:

```shell
Curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto.sh | bash
```

The commercial version performs this:

```shell
Curl -L https://raw.githubusercontent.com/snail007/goproxy/master/install_auto_commercial.sh | bash
```

The installation is complete, the configuration directory is /etc/proxy. For more detailed usage, please refer to the manual directory above to learn more about the features you want to use.
If the installation fails or your vps is not a linux64-bit system, follow the semi-automatic steps below to install:
  
### Manual installation

1. Download the proxy

Download address: https://github.com/snail007/goproxy/releases/latest

Let's take v7.9 as an example. If you have the latest version, please use the latest version of the link. Note that the version number in the download link below is the latest version number.

The free version performs this:

```shell
Cd /root/proxy/
Wget https://github.com/snail007/goproxy/releases/download/v7.9/proxy-linux-amd64.tar.gz
```

The commercial version performs this:

```shell
Cd /root/proxy/
Wget https://github.com/snail007/goproxy/releases/download/v7.9/proxy-linux-amd64_commercial.tar.gz
```

2. Download the automatic installation script

The free version performs this:

```shell
Cd /root/proxy/
Wget https://raw.githubusercontent.com/snail007/goproxy/master/install.sh
Chmod +x install.sh
./install.sh
```

The commercial version performs this:

```shell
Cd /root/proxy/
Wget https://raw.githubusercontent.com/snail007/goproxy/master/install_commercial.sh
Chmod +x install_commercial.sh
./install_commercial.sh
```

## TODO
- http,socks proxy multiple superior load balancing?
- http(s) proxy to increase pac support?
- Welcome to add group feedback..

## License
Proxy is licensed under GPLv3 license.

## Contact
Official QQ exchange group: 42805407

## Donation
If the proxy helps you solve a lot of problems, you can better support the proxy through the donation below.
<img src="https://github.com/snail007/goproxy/blob/master/doc/images/alipay.jpg?raw=true" width="200"/>
<img src="https://github.com/snail007/goproxy/blob/master/doc/images/wxpay.jpg?raw=true" width="200"/>

### Source code declaration

The author of this project found that a large number of developers based on the project for secondary development or using a large number of core code of the project without complying with the GPLv3 agreement, which seriously violates the original intention of using the GPLv3 open source agreement in this project. In view of this situation, the project adopts the source. The code delays the release strategy, to a certain extent, to curb these behaviors that do not respect open source and do not respect the labor results of others.
This project will continue to update the iterations and continue to release the full platform binary program, providing you with powerful and convenient agent tools.
If you have customized, business needs, please send an email to `arraykeys@gmail.com`