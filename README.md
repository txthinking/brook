# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![Docs](https://img.shields.io/badge/Tutorial-docs-yellow.svg)](https://txthinking.github.io/brook/)
[![Slides](https://img.shields.io/badge/Tutorial-Slides-blueviolet.svg)](https://talks.txthinking.com)
[![Youtube](https://img.shields.io/badge/Tutorial-Youtube-red.svg)](https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A)
[![Google Chat](https://img.shields.io/badge/Google-Chat-blue.svg)](https://docs.google.com/forms/d/e/1FAIpQLSd61-WE__WYiDee2UWjhDKNcb-A6KQW9xwzFkeYiAmQ3dpEcA/viewform)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)

<p align="center">
    <img style="float:right;" src="https://txthinking.github.io/brook/_static/brook.png" alt="Brook"/>
</p>

---

**v20200901**

- **❗️Breaking change, you should upgrade both server and client**
- New [Docs](https://txthinking.github.io/brook/)

---

### Table of Contents

- [What is Brook](#what-is-brook)
- [Install CLI](#install-cli)
- [Install GUI](#install-gui)
- [**Server**](#server)
- [**Client**](#client)
- [Map](#map)
- [DNS Server](#dns-server)
- [Transparent Proxy](#transparent-proxy)
- [WebSocket Server](#websocket-server)
- [WebSocket Client](#websocket-client)
- [Link](#link)
- [QR](#qr)
- [Relay](#relay)
- [Socks5 Server](#socks5-server)
- [Socks5 to HTTP](#socks5-to-http)
- [PAC](#pac)
- [How to](#how-to)
- [Contributing](#contributing)
- [License](#license)

## What is Brook

Brook is a cross-platform strong encryption and not detectable proxy.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

### Install CLI

> The CLI file has both server and client functions

Install via [nami](https://github.com/txthinking/nami)

```
nami install github.com/txthinking/brook
```

Install on Archlinux

```
pacman -S brook
```

Download from [releases](https://github.com/txthinking/brook/releases)

```
# For example, on linux amd64

$ curl -L https://github.com/txthinking/brook/releases/download/v20200901/brook_linux_amd64 -o /usr/bin/brook
$ chmod +x /usr/bin/brook
```

### Install GUI

> The GUI file has only client function

Install via brew

```
brew cask install brook
```

Download from [releases](https://github.com/txthinking/brook/releases)

[macOS](https://github.com/txthinking/brook/releases/download/v20200901/Brook.dmg), [Windows](https://github.com/txthinking/brook/releases/download/v20200901/Brook.exe), [Android](https://github.com/txthinking/brook/releases/download/v20200901/Brook.apk), [iOS](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)

## Brook

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20200901

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

### Server

```
# Run as a brook server
$ brook server -l :port -p password

# Run in background
$ nohup brook server -l :port -p password &

# Stop background brook
$ killall brook
```

> More parameters: $ brook server -h

### Client

```
# Run as brook client, connect to brook server, start a socks5 proxy server socks5://127.0.0.1:1080
$ brook client -s server_address:port -p password --socks5 127.0.0.1:1080
```

> More parameters: $ brook client -h

### Map

```
# Run as map, connect to brook server, map 127.0.0.1:5353 to 8.8.8.8:53
$ brook map -s server_address:port -p password -f 127.0.0.1:5353 -t 8.8.8.8:53
```

> More parameters: $ brook map -h

### DNS Server

```
# Run as DNS server, connect to brook server
$ brook dns -s server_address:port -p password -l 127.0.0.1:5353
```

> More parameters: $ brook dns -h

### Transparent Proxy

See [Docs](https://txthinking.github.io/brook/#/brook-tproxy)

### WebSocket Server

```
# Run as a brook wsserver
$ brook wsserver -l :port -p password
```

```
# Run as a brook wsserver with domain, make sure your domain name has been successfully resolved, 80 and 443 are open, brook will automatically issue certificate for you
$ brook wsserver --domain yourdomain.com -p password
```

> More parameters: $ brook wsserver -h

### WebSocket Client

```
# Run as brook wsclient, connect to brook wsserver, start a socks5 proxy server socks5://127.0.0.1:1080
$ brook wsclient -s ws://wsserver_address:port -p password --socks5 127.0.0.1:1080
```

```
# Run as brook wsclient, connect to brook wsserver with domain, start a http proxy
$ brook wsclient -s wss://wsserver_domain:port -p password --socks5 127.0.0.1:1080
```

> More parameters: $ brook wsclient -h

### Link

```
$ brook link -s server_address:port -p password
$ brook link -s ws://wsserver_address:port -p password
$ brook link -s wss://wsserver_domain:port -p password
```

> More parameters: $ brook link -h

### QR

```
$ brook qr -s server_address:port -p password
$ brook qr -s ws://wsserver_address:port -p password
$ brook qr -s wss://wsserver_domain:port -p password
```

> More parameters: $ brook qr -h

### Relay

```
# Run as relay
$ brook relay -f :port -t relay_to_address:port
```

> More parameters: $ brook relay -h

### Socks5 Server

```
# Run as standard socks5 server, assume your server public IP is 1.2.3.4
$ brook socks5 --socks5://1.2.3.4:1080
```

> More parameters: $ brook socks5 -h

### Socks5 to http

```
# Convert socks5 proxy socks5://127.0.0.1:1080 to http proxy http://127.0.0.1:8010
$ brook socks5tohttp -s 127.0.0.1:1080 -l 127.0.0.1:8010
```

> More parameters: $ brook socks5tohttp -h

### PAC

```
# Create PAC server
$ brook pac -l listen_address_port

# Save PAC to local file
$ brook pac -f /path/to/file.pac
```

> More parameters: $ brook pac -h

### How to

Some useful tutorial resources

* Brook Wiki: https://github.com/txthinking/brook/wiki
* Brook Issues: https://github.com/txthinking/brook/issues
* Slides: https://talks.txthinking.com
* Youtube: https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A
* Nami: https://github.com/txthinking/nami
* Joker: https://github.com/txthinking/joker

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
