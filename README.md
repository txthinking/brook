# Brook
A cross-platform programmable network tool. 一个跨平台可编程网络工具

# Sponsor
**❤️  [Shiliew - China Optimized VPN](https://www.txthinking.com/shiliew.html)**

Table of Contents
=================

* [Brook](#brook)
* [Sponsor](#sponsor)
* [Getting Started 快速上手](#getting-started-快速上手)
   * [Server](#server)
   * [GUI Client](#gui-client)
   * [CLI Client](#cli-client)
* [Install CLI 安装命令行](#install-cli-安装命令行)
   * [nami](#nami)
   * [brook](#brook-1)
   * [joker](#joker)
   * [jinbe](#jinbe)
   * [tun2brook](#tun2brook)
   * [via pacman](#via-pacman)
   * [via brew](#via-brew)
   * [via docker](#via-docker)
* [Daemon 守护进程](#daemon-守护进程)
* [Auto Start at Boot 开机自启](#auto-start-at-boot-开机自启)
* [One Click Script 一键脚本](#one-click-script-一键脚本)
* [CLI Documentation 命令行文档](#cli-documentation-命令行文档)
* [NAME](#name)
* [SYNOPSIS](#synopsis)
* [GLOBAL OPTIONS](#global-options)
* [COMMANDS](#commands)
   * [server](#server-1)
   * [client](#client)
   * [wsserver](#wsserver)
   * [wsclient](#wsclient)
   * [wssserver](#wssserver)
   * [wssclient](#wssclient)
   * [quicserver](#quicserver)
   * [quicclient](#quicclient)
   * [relayoverbrook](#relayoverbrook)
   * [dnsserveroverbrook](#dnsserveroverbrook)
   * [tproxy](#tproxy)
   * [link](#link)
   * [connect](#connect)
   * [relay](#relay)
   * [dnsserver](#dnsserver)
   * [dnsclient](#dnsclient)
   * [dohserver](#dohserver)
   * [dohclient](#dohclient)
   * [dhcpserver](#dhcpserver)
   * [socks5](#socks5)
   * [socks5tohttp](#socks5tohttp)
   * [pac](#pac)
   * [testsocks5](#testsocks5)
   * [testbrook](#testbrook)
   * [echoserver](#echoserver)
   * [echoclient](#echoclient)
   * [completion](#completion)
   * [mdpage](#mdpage)
      * [help, h](#help-h)
   * [manpage](#manpage)
   * [help, h](#help-h-1)
* [GUI Documentation](#gui-documentation)
   * [Software for which this article applies](#software-for-which-this-article-applies)
   * [Windows Proxy mode, Linux Proxy mode](#windows-proxy-mode-linux-proxy-mode)
   * [iOS, Mac, Android, Windows TUN mode, Linux TUN mode](#ios-mac-android-windows-tun-mode-linux-tun-mode)
   * [Configuration Introduction](#configuration-introduction)
   * [Programmable](#programmable)
      * [Introduction to incoming variables](#introduction-to-incoming-variables)
      * [in_brooklinks](#in_brooklinks)
      * [in_dnsquery](#in_dnsquery)
      * [in_address](#in_address)
      * [in_httprequest](#in_httprequest)
      * [in_httpresponse](#in_httpresponse)
   * [Write script](#write-script)
   * [Debug script](#debug-script)
   * [Standalone Script Example](#standalone-script-example)
   * [Brook Script Builder](#brook-script-builder)
   * [Packet Capture](#packet-capture)
   * [Install CA](#install-ca)
      * [iOS](#ios)
      * [Android](#android)
      * [macOS](#macos)
      * [Windows](#windows)
   * [Apple Push Problem](#apple-push-problem)
* [图形客户端文档](#图形客户端文档)
   * [本文适用的软件](#本文适用的软件)
   * [Windows GUI proxy 模式, Linux GUI proxy 模式](#windows-gui-proxy-模式-linux-gui-proxy-模式)
   * [iOS, Mac, Android, Windows TUN 模式, Linux TUN 模式](#ios-mac-android-windows-tun-模式-linux-tun-模式)
   * [配置介绍](#配置介绍)
   * [Programmable](#programmable-1)
      * [传入变量介绍](#传入变量介绍)
      * [in_brooklinks](#in_brooklinks-1)
      * [in_dnsquery](#in_dnsquery-1)
      * [in_address](#in_address-1)
      * [in_httprequest](#in_httprequest-1)
      * [in_httpresponse](#in_httpresponse-1)
   * [写脚本](#写脚本)
   * [调试脚本](#调试脚本)
   * [独立脚本例子](#独立脚本例子)
   * [脚本生成器](#脚本生成器)
   * [抓包](#抓包)
   * [安装 CA](#安装-ca)
      * [iOS](#ios-1)
      * [Android](#android-1)
      * [macOS](#macos-1)
      * [Windows](#windows-1)
   * [Apple 推送问题](#apple-推送问题)
* [Diagram 图解](#diagram-图解)
   * [overview](#overview)
   * [withoutBrookProtocol](#withoutbrookprotocol)
   * [relayoverbrook](#relayoverbrook-1)
   * [dnsserveroverbrook](#dnsserveroverbrook-1)
   * [relay](#relay-1)
   * [dnsserver](#dnsserver-1)
   * [tproxy](#tproxy-1)
   * [gui](#gui)
   * [script](#script)
* [Protocol](#protocol)
* [Blog](#blog)
* [YouTube](#youtube)
* [Telegram](#telegram)
* [brook-mamanger](#brook-mamanger)
* [nico](#nico)
* [Brook Deploy](#brook-deploy)
* [Pastebin](#pastebin)
* [独立脚本例子 | Standalone Script Example](#独立脚本例子--standalone-script-example)
* [脚本生成器 | Brook Script Builder](#脚本生成器--brook-script-builder)

# Brook
<!--SIDEBAR-->
<!--G-R3M673HK5V-->
A cross-platform programmable network tool. 一个跨平台可编程网络工具

# Sponsor
**❤️  [Shiliew - China Optimized VPN](https://www.txthinking.com/shiliew.html)**
# Getting Started 快速上手

## Server

```
bash <(curl https://bash.ooo/nami.sh)
```

```
nami install brook
```

```
brook server -l :9999 -p hello
```

## GUI Client

| iOS | Android      | Mac    |Windows      |Linux        |OpenWrt      |
| --- | --- | --- | --- | --- | --- |
| [![](https://brook.app/images/appstore.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![](https://brook.app/images/android.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.apk) | [![](https://brook.app/images/mac.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![Windows](https://brook.app/images/windows.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.exe) | [![](https://brook.app/images/linux.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.bin) | [![OpenWrt](https://brook.app/images/openwrt.png)](https://github.com/txthinking/brook/releases) |

> Linux: [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)<br/>
> OpenWrt: After installation, you need to refresh the page to see the menu

-   brook server: `1.2.3.4:9999` replace 1.2.3.4 with your server IP
-   password: `hello`

## CLI Client

```
brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```
# Install CLI 安装命令行

## nami

The easy way to download anything from anywhere

```
bash <(curl https://bash.ooo/nami.sh)
```

## brook

A cross-platform network tool

```
nami install brook
```

## joker

Joker can turn process into daemon

```
nami install joker
```

## jinbe

Auto start at boot. thanks to the cute cat

```
nami install jinbe
```

## tun2brook

Proxy all traffic just one line command

```
nami install tun2brook
```

## via pacman

maintained by felixonmars

```
pacman -S brook
```

## via brew

```
brew install brook
```

## via docker

maintained by teddysun

```
docker pull teddysun/brook
```
# Daemon 守护进程

Run the brook daemon with joker

```
joker brook server -l :9999 -p hello
```

Get the last command ID

```
joker last
```

View output and error of a command

```
joker log ID
```

View running commmands

```
joker list
```

Stop a running command

```
joker stop ID
```

# Auto Start at Boot 开机自启

Add one auto-start command at boot

```
jinbe joker brook server -l :9999 -p hello
```

View added commmands

```
jinbe list
```

Remove one added command

```
jinbe remove ID
```

# One Click Script 一键脚本

```
bash <(curl https://bash.ooo/brook.sh)
```
# CLI Documentation 命令行文档
# NAME

Brook - A cross-platform programmable network tool

# SYNOPSIS

Brook

```
[--dialWithDNSPrefer]=[value]
[--dialWithDNS]=[value]
[--dialWithIP4]=[value]
[--dialWithIP6]=[value]
[--dialWithNIC]=[value]
[--dialWithSocks5Password]=[value]
[--dialWithSocks5TCPTimeout]=[value]
[--dialWithSocks5UDPTimeout]=[value]
[--dialWithSocks5Username]=[value]
[--dialWithSocks5]=[value]
[--help|-h]
[--log]=[value]
[--pprof]=[value]
[--prometheusPath]=[value]
[--prometheus]=[value]
[--tag]=[value]
[--version|-v]
```

**Usage**:

```
Brook [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--dialWithDNS**="": When a domain name needs to be resolved, use the specified DNS. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required. Note that for client-side commands, this does not affect the client passing the domain address to the server

**--dialWithDNSPrefer**="": This is used with the dialWithDNS parameter. Prefer A record or AAAA record. Value is A or AAAA

**--dialWithIP4**="": When the current machine establishes a network connection to the outside IPv4, both TCP and UDP, it is used to specify the IPv4 used

**--dialWithIP6**="": When the current machine establishes a network connection to the outside IPv6, both TCP and UDP, it is used to specify the IPv6 used

**--dialWithNIC**="": When the current machine establishes a network connection to the outside, both TCP and UDP, it is used to specify the NIC used

**--dialWithSocks5**="": When the current machine establishes a network connection to the outside, both TCP and UDP, with your socks5 proxy, such as 127.0.0.1:1081

**--dialWithSocks5Password**="": If there is

**--dialWithSocks5TCPTimeout**="": time (s) (default: 0)

**--dialWithSocks5UDPTimeout**="": time (s) (default: 60)

**--dialWithSocks5Username**="": If there is

**--help, -h**: show help

**--log**="": Enable log. A valid value is file path or 'console'. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true

**--pprof**="": go http pprof listen addr, such as :6060

**--prometheus**="": prometheus http listen addr, such as :7070. If it is transmitted on the public network, it is recommended to use it with nico

**--prometheusPath**="": prometheus http path, such as /xxx. If it is transmitted on the public network, a hard-to-guess value is recommended

**--tag**="": Tag can be used to the process, will be append into log, such as: 'key1:value1'

**--version, -v**: print the version


# COMMANDS

## server

Run as brook server, both TCP and UDP

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--listen, -l**="": Listen address, like: ':9999'

**--password, -p**="": Server password

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

## client

Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]

**--http**="": Where to listen for HTTP proxy connections

**--password, -p**="": Brook server password

**--server, -s**="": Brook server address, like: 1.2.3.4:9999

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: UDP over TCP

## wsserver

Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--listen, -l**="": Listen address, like: ':80'

**--password, -p**="": Server password

**--path**="": URL path (default: /ws)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wsclient

Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--http**="": Where to listen for HTTP proxy connections

**--password, -p**="": Brook wsserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

**--wsserver, -s**="": Brook wsserver address, like: ws://1.2.3.4:80, if no path then /ws will be used. Do not omit the port under any circumstances

## wssserver

Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--password, -p**="": Server password

**--path**="": URL path (default: /ws)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wssclient

Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--http**="": Where to listen for HTTP proxy connections

**--insecure**: Client do not verify the server's certificate chain and host name

**--password, -p**="": Brook wssserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

**--wssserver, -s**="": Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances

## quicserver

Run as brook quicserver, both TCP and UDP

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--password, -p**="": Server password

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## quicclient

Run as brook quicclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]. (Note that the global dial parameter is ignored now)

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": Specify ca instead of insecure, such as /path/to/ca.pem

**--http**="": Where to listen for HTTP proxy connections

**--insecure**: Client do not verify the server's certificate chain and host name

**--password, -p**="": Brook quicserver password

**--quicserver, -s**="": Brook quicserver address, like: quic://google.com:443. Do not omit the port under any circumstances

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## relayoverbrook

Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--from, -f, -l**="": Listen address: like ':9999'

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## dnsserveroverbrook

Run as dns server over brook, both TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## tproxy

Run as transparent proxy, a router gateway, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--bypassCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--bypassDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassGeoIP**="": Bypass IP by Geo country code, such as US

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

**--dnsForDefault**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsListen**="": Start a DNS server, like: ':53'. MUST contain IP, like '192.168.1.1:53', if you expect your gateway to accept requests from clients to other public DNS servers at the same time

**--doNotRunScripts**: This will not change iptables and others if you want to do by yourself

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--link**="": brook link. This will ignore server, password, udpovertcp, address, insecure, withoutBrookProtocol, ca

**--listen, -l**="": Listen address, DO NOT contain IP, just like: ':8888'. No need to operate iptables by default! (default: :8888)

**--password, -p**="": Password

**--redirectDNS**="": It is usually the value of dnsListen. If the client has set custom DNS instead of dnsListen, this parameter can be intercepted and forwarded to dnsListen. Usually you don't need to set this, only if you want to control it instead of being proxied directly as normal UDP data.

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--webListen**="": Ignore all other parameters, run web UI, like: ':9999'

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## link

Generate brook link

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver or brook quicserver, specify ca for untrusted cert, such as /path/to/ca.pem

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--name**="": Give this server a name

**--password, -p**="": Password

**--server, -s**="": Support brook server, brook wsserver, brook wssserver, socks5 server, brook quicserver. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080, quic://google.com:443

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpovertcp**: When server is brook server, UDP over TCP

**--username, -u**="": Username, when server is socks5 server

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## connect

Run as client and connect to brook link, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

**--http**="": Where to listen for HTTP proxy connections

**--link, -l**="": brook link, you can get it via $ brook link

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## relay

Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]

**--from, -f, -l**="": Listen address: like ':9999'

**--tcpTimeout**="": time (s) (default: 0)

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": time (s) (default: 60)

## dnsserver

Run as standalone dns server

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## dnsclient

Send a dns query

**--dns, -s**="": DNS server, such as 8.8.8.8:53 (default: 8.8.8.8:53)

**--domain, -d**="": Domain

**--short**: Short for A/AAAA

**--type, -t**="": Type, such as A (default: A)

## dohserver

Run as standalone doh server

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

**--domainaddress**="": Such as: domain.com:443, if you want to create a https server. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--listen**="": listen address, if you want to create a http server behind nico

**--path**="": URL path (default: /dns-query)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## dohclient

Send a dns query

**--doh, -s**="": DOH server, the address is required (default: https://dns.quad9.net/dns-query?address=9.9.9.9%3A443)

**--domain, -d**="": Domain

**--short**: Short for A/AAAA

**--type, -t**="": Type, such as A (default: A)

## dhcpserver

Run as standalone dhcp server. Note that you need to stop other dhcp servers, if there are.

**--cache**="": Cache file, local absolute file path, default is $HOME/.brook.dhcpserver

**--count**="": IP range from the start, which you want to assign to clients (default: 100)

**--dnsserver**="": The dns server which you want to assign to clients, such as: 192.168.1.1 or 8.8.8.8

**--gateway**="": The router gateway which you want to assign to clients, such as: 192.168.1.1

**--interface**="": Select interface on multi interface device. Linux only

**--netmask**="": Subnet netmask which you want to assign to clients (default: 255.255.255.0)

**--serverip**="": DHCP server IP, the IP of the this machine, you shoud set a static IP to this machine before doing this, such as: 192.168.1.10

**--start**="": Start IP which you want to assign to clients, such as: 192.168.1.100

## socks5

Run as standalone standard socks5 server, both TCP and UDP

**--limitUDP**: The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.

**--listen, -l**="": Socks5 server listen address, like: :1080 or 1.2.3.4:1080

**--password**="": Password, optional

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--username**="": User name, optional

## socks5tohttp

Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]

**--listen, -l**="": HTTP proxy which will be create: like: 127.0.0.1:8010

**--socks5, -s**="": Socks5 server address, like: 127.0.0.1:1080

**--socks5password**="": Socks5 password, optional

**--socks5username**="": Socks5 username, optional

**--tcpTimeout**="": Connection tcp timeout (s) (default: 0)

## pac

Run as PAC server or save PAC to file

**--bypassDomainList, -b**="": One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--file, -f**="": Save PAC to file, this will ignore listen address

**--listen, -l**="": Listen address, like: 127.0.0.1:1980

**--proxy, -p**="": Proxy, like: 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' (default: SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT)

## testsocks5

Test UDP and TCP of socks5 server

**--dns**="": DNS server for connecting (default: 8.8.8.8:53)

**--domain**="": Domain for query (default: http3.ooo)

**--password, -p**="": Socks5 password

**--socks5, -s**="": Like: 127.0.0.1:1080

**--username, -u**="": Socks5 username

**-a**="": The A record of domain (default: 137.184.237.95)

## testbrook

Test UDP and TCP of brook server/wsserver/wssserver/quicserver. (Note that the global dial parameter is ignored now)

**--dns**="": DNS server for connecting (default: 8.8.8.8:53)

**--domain**="": Domain for query (default: http3.ooo)

**--link, -l**="": brook link. Get it via $ brook link

**--socks5**="": Temporarily listening socks5 (default: 127.0.0.1:11080)

**-a**="": The A record of domain (default: 137.184.237.95)

## echoserver

Echo server, echo UDP and TCP address of routes

**--listen, -l**="": Listen address, like: ':7777'

## echoclient

Connect to echoserver, echo UDP and TCP address of routes

**--server, -s**="": Echo server address, such as 1.2.3.4:7777

**--times**="": Times of interactions (default: 1)

## completion

Generate shell completions

**--file, -f**="": Write to file (default: brook_autocomplete)

## mdpage

Generate markdown page

**--file, -f**="": Write to file, default print to stdout

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## manpage

Generate man.1 page

**--file, -f**="": Write to file, default print to stdout. You should put to /path/to/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos

## help, h

Shows a list of commands or help for one command

# GUI Documentation

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## Software for which this article applies

-   [Brook](https://github.com/txthinking/brook)
-   [Brook Plus](https://www.txthinking.com/brook.html)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Windows Proxy mode, Linux Proxy mode

This mode is very simple, will create:

-   Socks5 proxy: `socks5://[::1]:1080` or `socks5://127.0.0.1:1080`
-   HTTP proxy: `http://[::1]:8010` or `http://127.0.0.1:8010`
-   PAC: `http://127.0.0.1:1093/proxy.pac` or `http://[::1]:1093/proxy.pac` based on Bypass Domain list
-   Intel Mac GUI, Windows GUI set PAC to system proxy。Linux GUI can work with [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
-   What is socks5 and http proxy? [Article](https://www.txthinking.com/talks/articles/socks5-and-http-proxy-en.article) and [Video](https://www.youtube.com/watch?v=sBCB-X7BoP8)

## iOS, Mac, Android, Windows TUN mode, Linux TUN mode

```
The so-called Internet connection is IP to IP connection, not domain name connection. Therefore, the domain name will be resolved into IP before deciding how to connect.
```

## Configuration Introduction

| Configuration  | Support Systems                   | Conditions                                                                          | Description  |
| -------------- | --------------------------------- | ----------------------------------------------------------------------------------- | --- |
| Import Servers | iOS,Android,Mac,Windows,Linux     | /                                                                                   | brook link list                                                                                                                                                                                                                                                                       |
| System DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System DNS. **Do not bypass this IP**                                                                                                                                                                                                                                                 |
| Fake DNS       | iOS,Android,Mac,Windows,Linux | [How to prevent Brook's Fake DNS from not working](https://www.txthinking.com/talks/articles/brook-fakedns-en.article) | The domain name is resolved to Fake IP, which will be converted to a domain name when a connection is initiated, and then the domain name address will be sent to the server, and the server is responsible for domain name resolution                                                |
| Block          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Block switch                                                                                                                                                                                                                                                                          |
| Block Domain   | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Domain name list, matching domain names will be blocked. **Domain name is suffix matching mode**                                                                                                                                                                                      |
| Bypass         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Bypass switch                                                                                                                                                                                                                                                                         |
| Bypass IP      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | CIDR list, matched IP will be bypassed                                                                                                                                                                                                                                                |
| Bypass Geo IP  | iOS,Android,Mac,Windows,Linux                       | /                                                                                   | The matched IP will be bypassed. Note: Global IP changes frequently, so the Geo library is time-sensitive                                                                                                                                                                             |
| Bypass Apps    | Android,Mac            | /                                                                                   | These apps will be bypassed                                                                                                                                                                                                                                                           |
| Bypass DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Support normal DNS, such as `223.5.5.5:53`, support DoH, but need to specify the address of DoH through the parameter address, such as `https://dns. alidns.com/dns-query?address=223.5.5.5%3A443` is used to resolve Bypass Domain. **The IP of this DNS will automatically Bypass** |
| Bypass Domain  | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | List of domain names, matching domain names will use Bypass DNS resolution to get IP, **whether the final connection will be bypassed depends on the Bypass IP** . **The domain name is a suffix matching pattern**. Of course, you can also use the script to bypass the domain regardless of its IP |
| Hosts          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Hosts switch                                                                                                                                                                                                                                                                          |
| Hosts List     | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Specify IP, v4, v6 for the domain name, if the value is empty, the effect is the same as Block                                                                                                                                                                                        |
| Programmable   | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Programmable switch                                                                                                                                                                                                                                                                   |
| Script         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Script. All functions above can be controlled. And more and more, **The whole process can be controled, see below for details**.                                                                                                                                                      |
| Log            | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log switch                                                                                                                                                                                                                                                                            |
| Activity  | iOS,Android,Mac,Windows,Linux     | /  | Networking activity |
| MITM  | iOS,Android,Mac,Windows,Linux     | / | Packet capture and modify, such as https request response, hexadecimal, JSON, image, etc.                                                                                                                                                                                                         |
| TUN            | iOS,Android,Mac,Windows,Linux                 | / | Choose Proxy/TUN/App mode. [macOS bug](https://www.txthinking.com/talks/articles/macos-bug.article). iOS,Android,Mac default TUN mode mode                                                                                                                                                                                                                                                         |
| Capture Me     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Test your packet capture or proxy software is working as a system proxy or TUN                                                                                                                                                                                                        |
| DNS Client    | iOS,Android,Mac,Windows,Linux | /                                           | DNS client                                                                                   |
| DOH Client    | iOS,Android,Mac,Windows,Linux | /                                           | DOH client                                                                                                           |
| Echo Client    | iOS,Android,Mac,Windows,Linux | /                                           | Echo client                                                                   |
| Test Socks5    | iOS,Android,Mac,Linux | /                                           | Test socks5 server                                                                    |
| Dark Mode      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System / Light / Dark                                                                                                                                                                                                                                                                 |     |
| Shortcut       | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Quickly control the functions in the menu on the home page                                                                                                                                                                                                                            |
| System Tray    | Windows                           | /                                                                                   | Open as systray, then open dashboard from the systray                                                                                                                                                                                                                                 |

## Programmable

```
Brook GUI will pass different global variables to the script at different times, and the script only needs to assign the processing result to the global variable out
```

### Introduction to incoming variables

| variable                       | type | condition   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_brooklinks                  | map  | / | Before connecting  | Predefine multiple brook links, and then programmatically specify which one to connect to | map      |
| in_dnsquery                    | map  | FakeDNS: On | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | /           | When connecting to an address     | script can decide how to connect                  | map      |
| in_httprequest                 | map  | /           | When an HTTP(S) request comes in  | the script can decide how to handle this request  | map      |
| in_httprequest,in_httpresponse | map  | /           | when an HTTP(S) response comes in | the script can decide how to handle this response | map      |

### in_brooklinks

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

### in_dnsquery

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| domain | string | domain name | google.com |
| type   | string | query type  | A          |
| appid   | string | App ID. Mac only | com.google.Chrome.helper          |
| interface   | string | network interface. Mac only | en0          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`. It is an OR relationship with GUI Block Domain                                                | false   |
| ip           | string | Specify IP directly, only valid when `type` is `A`/`AAAA`                                                                     | 1.2.3.4 |
| forcefakedns | bool   | Ignore GUI Bypass Domain, handle with Fake DNS, only valid when `type` is `A`/`AAAA`, default `false`                         | false   |
| system       | bool   | Get IP from system DNS, default `false`                                                                                       | false   |
| bypass       | bool   | whether to Bypass, default `false`, if `true` then use bypass DNS to resolve. It is an OR relationship with GUI Bypass Domain | false   |
| brooklinkkey | string   | When need to connect the Server，instead, connect to the brook link specified by the key in_brooklinks | custom name   |

### in_address

| Key           | Type   | Description                                                                                                         | Example        |
| ------------- | ------ | ------------------------------------------------------------------------------------------------------------------- | -------------- |
| network       | string | Network type, the value `tcp`/`udp`                                                                                 | tcp            |
| ipaddress     | string | IP type address. There is only of ipaddress and domainaddress. Note that there is no relationship between these two | 1.2.3.4:443    |
| domainaddress | string | Domain type address, because of FakeDNS we can get the domain name address here                                     | google.com:443 |
| appid   | string | App ID. Mac only | com.google.Chrome.helper          |
| interface   | string | network interface. Mac only | en0          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key                    | Type   | Description                                                                                                                                                                                             | Example     |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| block                  | bool   | Whether Block, default `false`                                                                                                                                                                          | false       |
| ipaddress              | string | IP type address, rewrite destination                                                                                                                                                                    | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | Use Bypass DNS to obtain `A` or `AAAA` IP and rewrite the destination, only valid when `domainaddress` exists, the value `A`/`AAAA`                                                                     | A           |
| bypass                 | bool   | Bypass, default `false`. If `true` and `domainaddress`, then `ipaddress` or `ipaddressfrombypassdns` must be specified. It is an OR relationship with GUI Bypass IP | false       |
| mitm                   | bool   | Whether to perform MITM, default `false`. Only valid when `network` is `tcp`. Need to install CA, see below                                                                                             | false       |
| mitmprotocol           | string | MITM protocol needs to be specified explicitly, the value is `http`/`https`                                                                                                                             | https       |
| mitmcertdomain         | string | The MITM certificate domain name, which is taken from `domainaddress` by default. If `ipaddress` and `mitm` is `true` and `mitmprotocol` is `https` then must be must be specified explicitly           | example.com |
| mitmwithbody           | bool   | Whether to manipulate the http body, default `false`. will read the body of the request and response into the memory and interact with the script. iOS 50M total memory limit may kill process      | false       |
| mitmautohandlecompress | bool   | Whether to automatically decompress the http body when interacting with the script, default `false`                                                                                                     | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                                                      | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                                                    | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                                                     | 0           |
| brooklinkkey | string   | When need to connect the Server，instead, connect to the brook link specified by the key in_brooklinks | custom name   |

### in_httprequest

| Key    | Type   | Description                   | Example                     |
| ------ | ------ | ----------------------------- | --------------------------- |
| URL    | string | URL                           | `https://example.com/hello` |
| Method | string | HTTP method                   | GET                         |
| Body   | bytes  | HTTP request body             | /                           |
| ...    | string | other fields are HTTP headers | /                           |

`out`, must be set to a request or response

### in_httpresponse

| Key        | Type   | Description                   | Example |
| ---------- | ------ | ----------------------------- | ------- |
| StatusCode | int    | HTTP status code              | 200     |
| Body       | bytes  | HTTP response body            | /       |
| ...        | string | other fields are HTTP headers | /       |

`out`, must be set to a response

## Write script

[Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

Library

-   [text](https://github.com/d5/tengo/blob/master/docs/stdlib-text.md): regular expressions, string conversion, and manipulation
-   [math](https://github.com/d5/tengo/blob/master/docs/stdlib-math.md): mathematical constants and functions
-   [times](https://github.com/d5/tengo/blob/master/docs/stdlib-times.md): time-related functions
-   [rand](https://github.com/d5/tengo/blob/master/docs/stdlib-rand.md): random functions
-   [fmt](https://github.com/d5/tengo/blob/master/docs/stdlib-fmt.md): formatting functions
-   [json](https://github.com/d5/tengo/blob/master/docs/stdlib-json.md): JSON functions
-   [enum](https://github.com/d5/tengo/blob/master/docs/stdlib-enum.md): Enumeration functions
-   [hex](https://github.com/d5/tengo/blob/master/docs/stdlib-hex.md): hex encoding and decoding functions
-   [base64](https://github.com/d5/tengo/blob/master/docs/stdlib-base64.md): base64 encoding and decoding functions
-   `brook`: brook module

    ```
    Constants

    * os: string, linux/darwin/windows/ios/android. If ios app run on mac, it is ios
    * iosapponmac: bool, ios app run on mac

    Functions

    * splithostport(address string) => map/error: splits a network address of the form "host:port" to { "host": "xxx", "port": "xxx" }
    * country(ip string) => string/error: get country code from ip
    * cidrcontainsip(cidr string, ip string) => bool/error: reports whether the network includes ip
    * parseurl(url string) => map/error: parses a raw url into a map, keys: scheme/host/path/rawpath/rawquery
    * parsequery(query string) => map/error: parses a raw query into a kv map
    * map2query(kv map) => string/error: convert map{string:string} into a query string
    * bytes2ints(b bytes) => array/error: convert bytes into [int]
    * ints2bytes(ints array) => bytes/error: convert [int] into bytes
    * bytescompare(a bytes, b bytes) => int/error: returns an integer comparing two bytes lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b
    * bytescontains(b bytes, sub bytes) => bool/error: reports whether sub is within b
    * byteshasprefix(s bytes, prefix bytes) => bool/error: tests whether the bytes s begins with prefix
    * byteshassuffix(s bytes, suffix bytes) => bool/error: tests whether the bytes s ends with suffix
    * bytesindex(s bytes, sep bytes) => int/error: returns the index of the first instance of sep in s, or -1 if sep is not present in s
    * byteslastindex(s bytes, sep bytes) => int/error: returns the index of the last instance of sep in s, or -1 if sep is not present in s
    * bytesreplace(s bytes, old bytes, new bytes, n int) => bytes/error: returns a copy of the s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements
    * pathescape(s string) => string/error: escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed
    * pathunescape(s string) => string/error: does the inverse transformation of pathescape
    * queryescape(s string) => string/error: escapes the string so it can be safely placed inside a URL query
    * queryunescape(s string) => string/error: does the inverse transformation of queryescape
    * hexdecode(s string) => bytes/error: returns the bytes represented by the hexadecimal string s
    * hexencode(s string) => string/error: returns the hexadecimal encoding of src
    ```

## Debug script

It is recommended to use [tun2brook](https://github.com/txthinking/tun2brook) on desktop to debug with `fmt.println`

## Standalone Script Example

https://github.com/txthinking/bypass

## Brook Script Builder

https://modules.brook.app

## Packet Capture

-   [Brook and mitmproxy for mobile app deep packet capture](https://www.txthinking.com/talks/articles/brook-mitmproxy-en.article)
-   [Brook Packet Capture on All Platform](https://www.txthinking.com/talks/articles/brook-packet-capture-en.article)
-   [mitmproxy helper](https://www.txthinking.com/mitmproxy.html)

## Install CA

https://txthinking.github.io/ca/ca.pem

### iOS

https://www.youtube.com/watch?v=HSGPC2vpDGk

### Android

Android has user CA and system CA, must be installed in the system CA after ROOT

### macOS

```
nami install mad ca.txthinking
sudo mad install --ca ~/.nami/bin/ca.pem
```

### Windows

Open GitBash

```
nami install mad ca.txthinking
```

Then open GitBash with admin

```
mad install --ca ~/.nami/bin/ca.pem
```

Note that software such as GitBash or Firefox may not read the system CA, you can use the system Edge browser to test after installation

## Apple Push Problem

To receive push, Apple Server only allows Ethernet, cellular data, Wi-Fi connections. So you need to Bypass the relevant domain name and IP. [Reference link](https://github.com/txthinking/bypass/tree/master/apple)
# 图形客户端文档

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## 本文适用的软件

-   [Brook](https://github.com/txthinking/brook)
-   [Brook Plus](https://www.txthinking.com/brook.html)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Windows GUI proxy 模式, Linux GUI proxy 模式

这个模式比较简单，会创建:

-   Socks5 代理: `socks5://[::1]:1080` 或 `socks5://127.0.0.1:1080`
-   HTTP 代理: `http://[::1]:8010` 或 `http://127.0.0.1:8010`
-   PAC: `http://127.0.0.1:1093/proxy.pac` 或 `http://[::1]:1093/proxy.pac` 基于 Bypass Domain 列表
-   Windows GUI 同时会配置 PAC 到系统代理。Linux GUI 可以配合 [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
-   什么是 socks5 和 http proxy? [文章](https://www.txthinking.com/talks/articles/socks5-and-http-proxy.article) 和 [视频](https://www.youtube.com/watch?v=Tb0_8odTxEI)

## iOS, Mac, Android, Windows TUN 模式, Linux TUN 模式

```
所谓的互联网连接，是 IP 连接 IP，不是连接域名。所以域名会被先解析成IP再决定怎么去连接。
```

## 配置介绍

| 配置    | 支持系统                      | 条件                 | 描述 |
| ------------- | ----------------------------- | ------------------------------------------- |  --- |
| 导入服务器    | iOS,Android,Mac,Windows,Linux | /                                           | brook link 列表                                                                                                                                                                                             |
| 系统 DNS    | iOS,Android,Mac,Windows,Linux | /                                           | 系统 DNS. **不要 bypass 此 IP**                                                                                                                                                                             |
| 虚拟 DNS      | iOS,Android,Mac,Windows,Linux | [如何避免 Brook 的 虚拟 DNS 不生效](https://www.txthinking.com/talks/articles/brook-fakedns.article) | 解析域名为 Fake IP，发起连接时会再转换为域名，然后把域名地址送到服务端进行代理，同时由服务端来负责域名解析                                                                                                  |
| 屏蔽         | iOS,Android,Mac,Windows,Linux | /                                           | Block 开关                                                                                                                                                                                                  |
| 屏蔽域名  | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 域名列表，匹配的域名会被阻断解析. **域名是后缀匹配模式**                                                                                                                                                    |
| 跳过        | iOS,Android,Mac,Windows,Linux | /                                           | Bypass 开关                                                                                                                                                                                                 |
| 跳过 IP     | iOS,Android,Mac,Windows,Linux | /                                           | CIDR 列表，匹配到的 IP 会被 bypass                                                                                                                                                                          |
| 跳过 Geo IP | iOS,Android,Mac,Windows,Linux                   | /                                           | 匹配到的 IP 会被 bypass. 提示: 全球 IP 变动频繁, 所以 Geo 库有时效性                                                                                                                                        |
| 跳过 Apps   | Android, Mac                       | /                                           | 这些 App 会被 bypass                                                                                                                                                                                        |
|  跳过 DNS    | iOS,Android,Mac,Windows,Linux | /                                           | 支持普通 DNS, 比如 `223.5.5.5:53`, 支持 DoH, 但需要通过参数 address 指定 DoH 的地址, 比如 `https://dns.alidns.com/dns-query?address=223.5.5.5%3A443` 用来解析 Bypass Domain. **此 DNS 的 IP 会自动 Bypass** |
| 跳过域名 | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 域名列表，匹配的域名会使用 Bypass DNS 解析来得到 IP, **最终连接是否会被 Bypass，还取决于 Bypass IP**. **域名是后缀匹配模式.** 当然也可以用脚本直接跳过域名而无关其IP               |
| Hosts         | iOS,Android,Mac,Windows,Linux | /                                           | Hosts 开关                                                                                                                                                                                                  |
| Host 列表    | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 给域名指定 IP, v4, v6，如果值为空效果同 Block                                                                                                                                                               |
| 可编程  | iOS,Android,Mac,Windows,Linux | /                                           | 可编程开关                                                                                                                                                                                                  |
| 脚本        | iOS,Android,Mac,Windows,Linux | /                                           | 脚本。可以控制上面所有的功能。以及上面没有的功能，**全流程控制一切，具体看下文**。                                                                                                                          |
| 日志           | iOS,Android,Mac,Windows,Linux | /                                           | 日志开关                                                                                                                                                                                                    |
| 活动 | iOS,Android,Mac,Windows,Linux | / | 网络活动 |
| MITM | iOS,Android,Mac,Windows,Linux | / | 抓包，修改包，比如 https 的请求响应，十六进制，JSON，图片等 |
| TUN           | iOS,Android,Mac,Windows,Linux             | /                                           | 选择 Proxy/TUN/App 模式. [macOS bug](https://www.txthinking.com/talks/articles/macos-bug.article). iOS,Android,Mac 默认 TUN 模式 |
| 抓我    | iOS,Android,Mac,Windows,Linux | /                                           | 测试你的抓包或代理软件工作在系统代理还是 TUN                                                                                                                                                                |
| DNS 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | DNS 客户端                                                                                                                              |
| DOH 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | DOH 客户端                                                                                                                              |
| Echo 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | Echo 客户端                                                                    |
| 测试 Socks5    | iOS,Android,Mac,Linux | /                                           | Test socks5 server                                                                    |
| 深色主题      | iOS,Android,Mac,Windows,Linux | /                                           | 暗黑模式                                                                                                                                                                                                    |     |
| 快捷方式      | iOS,Android,Mac,Windows,Linux | /                                           | 在首页快捷控制打开菜单里的功能                                                                                                                                                                              |
| 系统托盘      | Windows                       | /                                           | 以系统托盘形式打开，然后从系统托盘处打开控制面板                                                                                                                                                            |

## Programmable

```
Brook GUI 会在不同时机向脚本传入不同的全局变量，脚本只需要将处理结果赋值到全局变量 out 即可
```

### 传入变量介绍

| 变量                           | 类型 | 条件          | 时机                   | 描述                       | out 类型 |
| ------------------------------ | ---- | ------------- | ---------------------- | -------------------------- | -------- |
| in_brooklinks                  | map  | / | 连接之前  | 预定义多个 brook link，之后可编程指定连接哪个 | map      |
| in_dnsquery                    | map  | FakeDNS: 开启 | 当 DNS 查询发生时      | 脚本可以决定如何处理此请求 | map      |
| in_address                     | map  | /             | 当要连接某地址时       | 脚本可以决定如何进行连接   | map      |
| in_httprequest                 | map  | /             | 当有 HTTP(S)请求传入时 | 脚本可以决定如何处理此请求 | map      |
| in_httprequest,in_httpresponse | map  | /             | 当有 HTTP(S)响应传入时 | 脚本可以决定如何处理此响应 | map      |

### in_brooklinks

| Key    | 类型   | 描述     | 示例       |
| ------ | ------ | -------- | ---------- |
| _ | bool | 占位，无实际意义    | true |

`out`, 如果不是 `map` 类型则会被忽略

| Key          | 类型   | 描述                                                                                               | 示例    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| 自定义名字    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

### in_dnsquery

| Key    | 类型   | 描述     | 示例       |
| ------ | ------ | -------- | ---------- |
| domain | string | 域名     | google.com |
| type   | string | 查询类型 | A          |
| appid   | string | App ID. 仅 Mac | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key          | 类型   | 描述                                                                                               | 示例    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | 是否 Block, 默认 `false`. 与 GUI Block Domain 是或的关系                                           | false   |
| ip           | string | 直接指定 IP，仅当 `type` 为 `A`/`AAAA`有效                                                         | 1.2.3.4 |
| forcefakedns | bool   | 忽略 GUI Bypass Domain，使用 Fake DNS 来处理，仅当 `type` 为 `A`/`AAAA`有效，默认 `false`          | false   |
| system       | bool   | 使用 System DNS 来解析，默认 `false`                                                               | false   |
| bypass       | bool   | 是否 Bypass, 默认 `false`, 如果为 `true` 则使用 Bypass DNS 来解析. 与 GUI Bypass Domain 是或的关系 | false   |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的 brook link | 自定义名字   |

### in_address

| Key           | 类型   | 描述                                                                     | 示例           |
| ------------- | ------ | ------------------------------------------------------------------------ | -------------- |
| network       | string | 即将发起连接网络，取值 `tcp`/`udp`                                       | tcp            |
| ipaddress     | string | IP 类型的地址，与 domainaddress 只会存在一个。注意这两个之间没有任何关系 | 1.2.3.4:443    |
| domainaddress | string | 域名类型的地址，因为 FakeDNS 我们这里才能拿到域名地址                    | google.com:443 |
| appid   | string | App ID. 仅 Mac | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key                    | 类型   | 描述                                                                                                                                                                           | 示例        |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- |
| block                  | bool   | 是否 Block, 默认 `false`                                                                                                                                                       | false       |
| ipaddress              | string | IP 类型地址，重写目的地                                                                                                                                                        | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | 使用 Bypass DNS 获取`A`或`AAAA` IP 并重写目的地, 仅当 `domainaddress` 存在时有效，取值 `A`/`AAAA`                                                                              | A           |
| bypass                 | bool   | 是否 Bypass, 默认 `false`. 如果为 `true` 并且是 `domainaddress`, 那么必须指定 `ipaddress` 或 `ipaddressfrombypassdns`. 与 GUI Bypass IP 是或的关系 | false       |
| mitm                   | bool   | 是否进行 MITM, 默认 `false`. 仅当 `network` 为 `tcp` 时有效. 需要安装 CA，看下文介绍                                                                                           | false       |
| mitmprotocol           | string | 需要明确指定 MITM 协议, 取值 `http`/`https`                                                                                                                                    | https       |
| mitmcertdomain         | string | MITM 证书域名，默认从`domainaddress`里取。如果是 `ipaddress` 且 `mitm` 为 `true` 且 `mitmprotocol` 为 `https` 那么必须明确指定                                                 | example.com |
| mitmwithbody           | bool   | 是否操作 http body，默认 `false`. 会将请求和响应的 body 读取到内存里和脚本交互。iOS 50M 总内存限制可能会杀进程                                                             | false       |
| mitmautohandlecompress | bool   | 和脚本交互时是否自动解压缩 http body, 默认 `false`                                                                                                                             | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                             | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                           | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                            | 0           |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的 brook link | 自定义名字   |

### in_httprequest

| Key    | 类型   | 描述                     | 示例                        |
| ------ | ------ | ------------------------ | --------------------------- |
| URL    | string | URL                      | `https://example.com/hello` |
| Method | string | HTTP method              | GET                         |
| Body   | bytes  | HTTP request body        | /                           |
| ...    | string | 其他字段均为 HTTP header | /                           |

`out`, 必须设置为一个 request 或 response

### in_httpresponse

| Key        | 类型   | 描述                     | 示例 |
| ---------- | ------ | ------------------------ | ---- |
| StatusCode | int    | HTTP status code         | 200  |
| Body       | bytes  | HTTP response body       | /    |
| ...        | string | 其他字段均为 HTTP header | /    |

`out`, 必须设置为一个 response

## 写脚本

[Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

Library

-   [text](https://github.com/d5/tengo/blob/master/docs/stdlib-text.md): regular expressions, string conversion, and manipulation
-   [math](https://github.com/d5/tengo/blob/master/docs/stdlib-math.md): mathematical constants and functions
-   [times](https://github.com/d5/tengo/blob/master/docs/stdlib-times.md): time-related functions
-   [rand](https://github.com/d5/tengo/blob/master/docs/stdlib-rand.md): random functions
-   [fmt](https://github.com/d5/tengo/blob/master/docs/stdlib-fmt.md): formatting functions
-   [json](https://github.com/d5/tengo/blob/master/docs/stdlib-json.md): JSON functions
-   [enum](https://github.com/d5/tengo/blob/master/docs/stdlib-enum.md): Enumeration functions
-   [hex](https://github.com/d5/tengo/blob/master/docs/stdlib-hex.md): hex encoding and decoding functions
-   [base64](https://github.com/d5/tengo/blob/master/docs/stdlib-base64.md): base64 encoding and decoding functions
-   `brook`: brook module

    ```
    Constants

    * os: string, linux/darwin/windows/ios/android. If ios app run on mac, it is ios
    * iosapponmac: bool, ios app run on mac

    Functions

    * splithostport(address string) => map/error: splits a network address of the form "host:port" to { "host": "xxx", "port": "xxx" }
    * country(ip string) => string/error: get country code from ip
    * cidrcontainsip(cidr string, ip string) => bool/error: reports whether the network includes ip
    * parseurl(url string) => map/error: parses a raw url into a map, keys: scheme/host/path/rawpath/rawquery
    * parsequery(query string) => map/error: parses a raw query into a kv map
    * map2query(kv map) => string/error: convert map{string:string} into a query string
    * bytes2ints(b bytes) => array/error: convert bytes into [int]
    * ints2bytes(ints array) => bytes/error: convert [int] into bytes
    * bytescompare(a bytes, b bytes) => int/error: returns an integer comparing two bytes lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b
    * bytescontains(b bytes, sub bytes) => bool/error: reports whether sub is within b
    * byteshasprefix(s bytes, prefix bytes) => bool/error: tests whether the bytes s begins with prefix
    * byteshassuffix(s bytes, suffix bytes) => bool/error: tests whether the bytes s ends with suffix
    * bytesindex(s bytes, sep bytes) => int/error: returns the index of the first instance of sep in s, or -1 if sep is not present in s
    * byteslastindex(s bytes, sep bytes) => int/error: returns the index of the last instance of sep in s, or -1 if sep is not present in s
    * bytesreplace(s bytes, old bytes, new bytes, n int) => bytes/error: returns a copy of the s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements
    * pathescape(s string) => string/error: escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed
    * pathunescape(s string) => string/error: does the inverse transformation of pathescape
    * queryescape(s string) => string/error: escapes the string so it can be safely placed inside a URL query
    * queryunescape(s string) => string/error: does the inverse transformation of queryescape
    * hexdecode(s string) => bytes/error: returns the bytes represented by the hexadecimal string s
    * hexencode(s string) => string/error: returns the hexadecimal encoding of src
    ```

## 调试脚本

建议使用 [tun2brook](https://github.com/txthinking/tun2brook) 在电脑上`fmt.println`调试

## 独立脚本例子

https://github.com/txthinking/bypass

## 脚本生成器

https://modules.brook.app

## 抓包

-   [Brook 搭配 mitmproxy 进行手机 App 深度抓包](https://www.txthinking.com/talks/articles/brook-mitmproxy.article)
-   [Brook 全平台抓包](https://www.txthinking.com/talks/articles/brook-packet-capture.article)
-   [用 mitmproxy helper 抓包](https://www.txthinking.com/mitmproxy.html)

## 安装 CA

https://txthinking.github.io/ca/ca.pem

### iOS

https://www.youtube.com/watch?v=HSGPC2vpDGk

### Android

Android 分系统 CA 和用户 CA，必须要 ROOT 后安装到系统 CA 里

### macOS

```
nami install mad ca.txthinking
sudo mad install --ca ~/.nami/bin/ca.pem
```

### Windows

打开 GitBash

```
nami install mad ca.txthinking
```

然后用管理员打开 GitBash

```
mad install --ca ~/.nami/bin/ca.pem
```

注意 GitBash 或 Firefox 等软件可能不读取系统 CA，安装后可以用系统 Edge 浏览器测试

## Apple 推送问题

要接收推送，Apple Server 只允许 Ethernet, cellular data, Wi-Fi 连接. 所以你需要 Bypass 掉相关域名和 IP. [参考链接](https://github.com/txthinking/bypass/tree/master/apple)
# Diagram 图解

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## overview

![overview](https://txthinking.github.io/brook/svg/overview.svg)

## withoutBrookProtocol

![wbp](https://txthinking.github.io/brook/svg/wbp.svg)

## relayoverbrook

![relayoverbrook](https://txthinking.github.io/brook/svg/relayoverbrook.svg)

## dnsserveroverbrook

![dnsserveroverbrook](https://txthinking.github.io/brook/svg/dnsserveroverbrook.svg)

## relay

![relay](https://txthinking.github.io/brook/svg/relay.svg)

## dnsserver

![dnsserver](https://txthinking.github.io/brook/svg/dnsserver.svg)

## tproxy

![tproxy](https://txthinking.github.io/brook/svg/tproxy.svg)

## gui

![gui](https://txthinking.github.io/brook/svg/gui.svg)

## script

![script](https://txthinking.github.io/brook/svg/script.svg)

# Protocol
https://github.com/txthinking/brook/tree/master/protocol
# Blog
https://www.txthinking.com/talks/
# YouTube
https://www.youtube.com/txthinking
# Telegram
https://t.me/s/txthinking_news
# brook-mamanger
https://github.com/txthinking/brook-manager
# nico
https://github.com/txthinking/nico
# Brook Deploy
https://www.txthinking.com/deploy.html
# Pastebin
https://paste.brook.app
# 独立脚本例子 | Standalone Script Example
https://github.com/txthinking/bypass
# 脚本生成器 | Brook Script Builder
https://modules.brook.app
