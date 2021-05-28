# Brook

[中文](README_ZH.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![Docs](https://img.shields.io/badge/Tutorial-docs-yellow.svg)](https://txthinking.github.io/brook/)
[![Blog](https://img.shields.io/badge/Tutorial-Slides-blueviolet.svg)](https://talks.txthinking.com)
[![Youtube](https://img.shields.io/badge/Tutorial-Youtube-red.svg)](https://www.youtube.com/txthinking)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Donate](https://img.shields.io/badge/Support-Donate-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)

<p align="center">
    <img style="float:right;" src="https://txthinking.github.io/brook/_static/brook.png" alt="Brook"/>
</p>

---

**v20210601**

-   [CLI] \$ brook map supports brook server/wsserver/wssserver
-   [CLI] \$ brook dns supports brook server/wsserver/wssserver
-   [CLI] \$ brook tproxy supports brook server/wsserver/wssserver
-   [GUI] OpenWrt supports brook server/wsserver/wssserver
-   [Document](https://txthinking.github.io/brook/)
-   [Community(ask here)](https://github.com/txthinking/brook/discussions)

---

## What is Brook

Brook is a cross-platform strong encryption and not detectable proxy.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

[Read Document](https://txthinking.github.io/brook/#/?id=cli-and-gui)

### Install CLI

```
$ curl -L https://github.com/txthinking/brook/releases/latest/download/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

[Read Document](https://txthinking.github.io/brook/#/install-cli)

### Install GUI

[Read Document](https://txthinking.github.io/brook/#/install-gui-client)

## Usage

[Read Document](https://txthinking.github.io/brook/)

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20210601

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   server        Run as brook server, both TCP and UDP
   servers       Run as multiple brook servers
   client        Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst], [works with
$ brook server]
   map           Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook <-> to addres
s], works with $ brook server/wsserver/wssserver
   dns           Run as DNS server, both TCP and UDP, [src <-> $ brook dns <-> $ brook <-> dns server] or [src <-> $ brook dns <-> dns server for bypass], wo
rks with $ brook server/wsserver/wssserver
   tproxy        Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook <-> dst], works with $ brook server/wss
erver/wssserver
   wsserver      Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server
   wssserver     Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server
   wsclient      Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst], [works
 with $ brook wsserver]
   wssclient     Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst], [wo
rks with $ brook wssserver]
   link          Print brook link
   qr            Print brook server QR code
   connect       Connect via standard sharing link (brook server & brook wsserver & brook wssserver)
   relay         Run as standalone relay, both TCP and UDP, this means access [listen address] is equal to access [to address], [src <-> listen address <-> t
o address]
   relays        Run as multiple standalone relays
   socks5        Run as standalone standard socks5 server, both TCP and UDP
   socks5tohttp  Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]
   hijackhttps   Hijack domains and assume is TCP/TLS/443. Requesting these domains from anywhere in the system will be hijacked . [src <-> $ brook hijackhtt
ps <-> socks5 server] or [src <-> direct]
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

[Read Document](https://txthinking.github.io/brook/)

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
