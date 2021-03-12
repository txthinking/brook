# Brook

[English](README.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![文档](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E6%96%87%E6%A1%A3-yellow.svg)](https://txthinking.github.io/brook/#/zh-cn/)
[![博客](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E5%B9%BB%E7%81%AF%E7%89%87-blueviolet.svg)](https://talks.txthinking.com)
[![视频](https://img.shields.io/badge/%E6%95%99%E7%A8%8B-%E8%A7%86%E9%A2%91-red.svg)](https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A)
[![开源协议: GPL v3](https://img.shields.io/badge/%E5%BC%80%E6%BA%90%E5%8D%8F%E8%AE%AE-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)

<p align="center">
    <img style="float:right;" src="https://txthinking.github.io/brook/_static/brook.png" alt="Brook"/>
</p>

---

**v20210214**

- 只是更新了客户端, 优化link, QR, 分享, 导入
- [官方文档](https://txthinking.github.io/brook/#/zh-cn/)
- [官方论坛](https://github.com/txthinking/brook/discussions)

---

## 什么是 Brook

[官方文档](https://txthinking.github.io/brook/#/zh-cn/)

Brook 是一个跨平台的强加密无特征的代理软件. 偏爱 KISS 哲学.

### 安装 CLI (命令行版本)

> CLI 版本同时具有服务端和客户端等很多功能

从 [releases](https://github.com/txthinking/brook/releases) 页面下载

```
# 举例, linux amd64, v20210214

curl -L https://github.com/txthinking/brook/releases/download/v20210214/brook_linux_amd64 -o /usr/bin/brook
chmod +x /usr/bin/brook
```

通过 [nami](https://github.com/txthinking/nami) 安装

```
nami install github.com/txthinking/brook
```

### 安装 GUI (图形客户端)

从 [releases](https://github.com/txthinking/brook/releases) 页面下载: [macOS](https://github.com/txthinking/brook/releases/download/v20210214/Brook.dmg), [Windows](https://github.com/txthinking/brook/releases/download/v20210214/Brook.exe), [Android](https://github.com/txthinking/brook/releases/download/v20210214/Brook.apk), [iOS](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)

通过 brew 安装

```
brew install --cask brook
```

## 使用

[文档](https://txthinking.github.io/brook/#/zh-cn/)

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20210214

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   server        Run as brook server, both TCP and UDP
   servers       Run as multiple brook servers
   client        Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst], [works with $ brook server]
   map           Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server <-> to address], [works with $ brook server]
   dns           Run as DNS server, both TCP and UDP, [src <-> $ brook dns <-> $ brook server <-> dns server] or [src <-> $ brook dns <-> dns server for bypass], [works with $ brook server]
   tproxy        Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server <-> dst], [works with $ brook server]
   wsserver      Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server
   wssserver     Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server
   wsclient      Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst], [works with $ brook wsserver]
   wssclient     Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst], [works with $ brook wssserver]
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
