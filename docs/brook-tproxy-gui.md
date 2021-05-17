## \$ brook tproxy

> **Dependencies: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy**

$ brook tproxy can create Transparent Proxy on your linux router with `TPROXY mod`, **port 9999, 1080, 5353 will be used**. It must work with $ brook server.

> Only support IPv4 server, but if your server support IPv6 you can enable later, please see the introduction below

## Install ipk

1. Download the [ipk](https://github.com/txthinking/brook/releases) file for your router
2. Upload and install: OpenWrt Web -> System -> Software -> Upload Package...
3. Refresh page, the Brook menu will appear at the top
4. OpenWrt Web -> Brook -> type and Connect
5. And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353 or other port brook created
6. And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

## Error log files

* `/root/.brook.web.err`
* `/root/.brook.tproxy.err`
