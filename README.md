# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![Docs](https://img.shields.io/badge/Tutorial-docs-yellow.svg)](https://txthinking.github.io/brook/)
[![Slides](https://img.shields.io/badge/Tutorial-Slides-blueviolet.svg)](https://talks.txthinking.com)
[![Youtube](https://img.shields.io/badge/Tutorial-Youtube-red.svg)](https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)

<p align="center">
    <img style="float:right;" src="https://txthinking.github.io/brook/_static/brook.png" alt="Brook"/>
</p>

---

**v20200901**

- **❗️Breaking change, you should upgrade both server and client**
- New [Docs](https://txthinking.github.io/brook/)


**v20200909**

- Import server list from URL. [ref here](https://txthinking.github.io/brook/#/brook-link) and [here](https://gist.githubusercontent.com/txthinking/7ecdb282982e14cc95714141c0ce2581/raw/350363229d1ce123b87b7cb0789e459969620cb3/brooklink.list)

---

## What is Brook

Brook is a cross-platform strong encryption and not detectable proxy.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

### Install CLI

> The CLI file has both server and client functions

Download from [releases](https://github.com/txthinking/brook/releases)

```
# For example, on linux amd64, v20200909

$ curl -L https://github.com/txthinking/brook/releases/download/v20200909/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

Install via [nami](https://github.com/txthinking/nami)

```
nami install github.com/txthinking/brook
```

### Install GUI

> The GUI file has only client function

Download from [releases](https://github.com/txthinking/brook/releases): [macOS](https://github.com/txthinking/brook/releases/download/v20200909/Brook.dmg), [Windows](https://github.com/txthinking/brook/releases/download/v20200909/Brook.exe), [Android](https://github.com/txthinking/brook/releases/download/v20200909/Brook.apk), [iOS](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)

Install via brew

```
brew cask install brook
```

## Usage

[Docs](https://txthinking.github.io/brook/)

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20200909

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

[Docs](https://txthinking.github.io/brook/)

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
