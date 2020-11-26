# Brook

[English](README.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![文档](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E6%96%87%E6%A1%A3-yellow.svg)](https://txthinking.github.io/brook/#/zh-cn/)
[![幻灯片](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E5%B9%BB%E7%81%AF%E7%89%87-blueviolet.svg)](https://talks.txthinking.com)
[![视频](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E8%A7%86%E9%A2%91-red.svg)](https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A)
[![开源协议: GPL v3](https://img.shields.io/badge/%E5%BC%80%E6%BA%90%E5%8D%8F%E8%AE%AE-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)
[![交流群](https://img.shields.io/badge/%E7%94%B3%E8%AF%B7%E5%8A%A0%E5%85%A5-%E4%BA%A4%E6%B5%81%E7%BE%A4-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

<p align="center">
    <img style="float:right;" src="https://txthinking.github.io/brook/_static/brook.png" alt="Brook"/>
</p>

---

**v20200901**

- **此版本不兼容之前的版本, 建议一起升级服务端和客户端**
- [新的文档站点](https://txthinking.github.io/brook/#/zh-cn/)


**v20210101**

- 支持从HTTP URL导入服务器列表. [参考这里](https://txthinking.github.io/brook/#/brook-link) 和 [这里](https://gist.githubusercontent.com/txthinking/7ecdb282982e14cc95714141c0ce2581/raw/350363229d1ce123b87b7cb0789e459969620cb3/brooklink.list)
- 恢复支持iOS 13

---

## 什么是 Brook

Brook 是一个跨平台的强加密无特征的代理软件. 偏爱 KISS 哲学.

### 安装 CLI (命令行版本)

> CLI 版本同时具有服务端和客户端等很多功能

从 [releases](https://github.com/txthinking/brook/releases) 页面下载

```
# 举例, linux amd64, v20210101

curl -L https://github.com/txthinking/brook/releases/download/v20210101/brook_linux_amd64 -o /usr/bin/brook
chmod +x /usr/bin/brook
```

通过 [nami](https://github.com/txthinking/nami) 安装

```
nami install github.com/txthinking/brook
```

### 安装 GUI (图形客户端)

从 [releases](https://github.com/txthinking/brook/releases) 页面下载: [macOS](https://github.com/txthinking/brook/releases/download/v20210101/Brook.dmg), [Windows](https://github.com/txthinking/brook/releases/download/v20210101/Brook.exe), [Android](https://github.com/txthinking/brook/releases/download/v20210101/Brook.apk), [iOS](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)

通过 brew 安装

```
brew cask install brook
```

## 使用

[文档](https://txthinking.github.io/brook/#/zh-cn/)

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20210101

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   server        Run as brook server, both TCP and UDP
   servers       Run as multiple brook servers
   client        Run as brook client, both TCP and UDP, to start a socks5 proxy or a http proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst], [works with $ brook server]
   map           Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server <-> to address], [works with $ brook server]
   dns           Run as DNS server, both TCP and UDP, [src <-> $ brook dns <-> $ brook server <-> dns server] or [src <-> $ brook dns <-> dns server for bypass], [works with $ brook server]
   tproxy        Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server <-> dst], [works with $ brook server]
   wsserver      Run as brook wsserver, both TCP and UDP, it will start a standard http(s) server and websocket server
   wsclient      Run as brook wsclient, both TCP and UDP, to start a socks5 proxy or a http proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst], [works with $ brook wsserver]
   link          Print brook link
   qr            Print brook server QR code
   relay         Run as standalone relay, both TCP and UDP, this means access [listen address] is equal to access [to address], [src <-> listen address <-> to address]
   relays        Run as multiple standalone relays
   socks5        Run as standalone standard socks5 server, both TCP and UDP
   socks5tohttp  Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]
   hijackhttps   Hijack domains and assume is TCP/TLS/443. Requesting these domains from anywhere in the system will be hijacked . [src <-> $ brook hijackhttps <-> socks5 server] or [src <-> direct]
   pac           Run as PAC server or save PAC to file
   howto         Print some useful tutorial resources
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug (default: false)
   --listen value, -l value  Listen address for debug (default: ":6060")
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)

COPYRIGHT:
   https://github.com/txthinking/brook
```

[文档](https://txthinking.github.io/brook/#/zh-cn/)

## 贡献

请先阅读 [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md)

## 开源协议

基于 GPLv3 协议开源
