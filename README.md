# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![Wiki](https://img.shields.io/badge/docs-wiki-yellow.svg)](https://github.com/txthinking/brook/wiki)
[![Slides](https://img.shields.io/badge/Tutorial-Slides-blueviolet.svg)](https://talks.txthinking.com)
[![Youtube](https://img.shields.io/badge/Tutorial-Youtube-red.svg)](https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A)
[![Telegram Group](https://img.shields.io/badge/Telegram%20Group-brookgroup-blue.svg)](https://t.me/brookgroup)
[![Telegram Channel](https://img.shields.io/badge/Telegram%20Channel-brookchannel-blue.svg)](https://t.me/brookchannel)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](http://www.gnu.org/licenses/gpl-3.0)

<p align="center">
    <img style="float:right;" src="https://storage.googleapis.com/txthinking/_/brook.png" alt="Brook"/>
</p>

---

**v20200502**

* CLI: Add `$ brook dns`
* CLI: Add `$ brook pac`
* GUI: Support multiple servers

Some useful tutorial resources

* Slides: https://talks.txthinking.com
* Youtube: https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A
* Telegram Group: https://t.me/brookgroup
* Telegram Channel: https://t.me/brookchannel

---

### Table of Contents

- [What is Brook](#what-is-brook)
- [Install](#install-via-nami)
- [**Server**](#server)
- [**Client**](#client)
- [Tunnel](#tunnel)
- [DNS Server](#dns-server)
- [Transparent Proxy](#transparent-proxy)
- [Tun](#tun)
- [WebSocket Server](#websocket-server)
- [WebSocket Client](#websocket-client)
- [Link](#link)
- [QR](#qr)
- [Relay](#relay)
- [Socks5 Server](#socks5-server)
- [Socks5 to HTTP](#socks5-to-http)
- [PAC](#pac)
- [System Proxy](#system-proxy)
- [Shadowsocks](#shadowsocks)
- [How to](#how-to)
- [Contributing](#contributing)
- [License](#license)

## What is Brook

Brook is a cross-platform proxy/vpn software.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

### Install via [nami](https://github.com/txthinking/nami)

install CLI using nami on Linux/BSD/macOS

```
nami install github.com/txthinking/brook
```

or install CLI on Archlinux

```
pacman -S brook
```

**or download CLI from [releases](https://github.com/txthinking/brook/releases)**

install GUI on macOS

```
brew cask install brook
```

**or download GUI: [macOS](https://github.com/txthinking/brook/releases/download/v20200502/Brook.pkg), [Windows](https://github.com/txthinking/brook/releases/download/v20200502/Brook.msi), [Android](https://github.com/txthinking/brook/releases/download/v20200502/Brook.apk), [iOS](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)**

> CLI contains server and client, GUI only contains client. iOS client only supports non-China AppStore.

## Brook

```
NAME:
   Brook - A Cross-Platform Proxy/VPN Software

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20200502

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   server        Run as brook server
   servers       Run as multiple brook servers
   client        Run as brook client, with brook server, to start a socks5 proxy or a http proxy
   tunnel        Run as tunnel, with brook server
   dns           Run as DNS server, with brook server
   tproxy        Run as transparent proxy, with brook server, only works on Linux
   tun           Run as tun, with brook server
   wsserver      Run as brook wsserver, it will start a standard http server and websocket server
   wsclient      Run as brook wsclient, with brook wsserver, to start a socks5 proxy or a http proxy
   link          Print brook link
   qr            Print brook server QR code
   relay         Run as standalone relay
   relays        Run as multiple standalone relays
   socks5        Run as standalone standard socks5 server
   socks5tohttp  Convert socks5 to http proxy
   pac           Run as PAC server or save PAC to file
   systemproxy   Set system proxy with pac url, or remove, only works on macOS/Windows
   ssserver      Run as shadowsocks server, fixed method is aes-256-cfb
   ssservers     Run as shadowsocks multiple servers, fixed method is aes-256-cfb
   ssclient      Run as shadowsocks client, with shadowsocks server, to start a socks5 proxy or a http proxy, fixed method is aes-256-cfb
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
$ brook server -l listen_address:port -p password
```

> More parameters: $ brook server -h

### Client

```
# Run as brook client, connect to brook server, start a socks5 proxy server
$ brook client -s server_address:port -p password -l listen_address:port -i socks5_server_ip
```

> More parameters: $ brook client -h

### Tunnel

```
# Run as tunnel, connect to brook server
$ brook tunnel -s server_address:port -p password -l listen_address:port -t tunnel_to_address:port
```

> More parameters: $ brook tunnel -h

### DNS Server

```
# Run as DNS server, connect to brook server
$ brook dns -s server_address:port -p password -l listen_address:port
```

> More parameters: $ brook dns -h

### Transparent Proxy

See [wiki](https://github.com/txthinking/brook/wiki/How-to-run-transparent-proxy-on-Linux)

### Tun

```
# Run as tun, connect to brook server, proxy all TCP/UDP. [ROOT privileges required].
$ sudo brook tun -s server_address:port -p password -l 127.0.0.1:1080

# Must exit by Ctrl+C
```

> More parameters: $ brook tun -h

**See [wiki](https://github.com/txthinking/brook/wiki/How-to-run-tun-on-Linux,-macOS-and-Windows) for more tutorials**

### WebSocket Server

```
# Run as a brook wsserver
$ brook wsserver -l listen_address:port -p password
```

```
# Run as a brook wsserver with domain, make sure your domain name has been successfully resolved, 80 and 443 are open, brook will automatically issue certificate for you
$ brook wsserver --domain yourdomain.com -p password
```

> More parameters: $ brook wsserver -h

### WebSocket Client

```
# Run as brook wsclient, connect to brook wsserver, start a socks5 proxy server
$ brook wsclient -s ws://wsserver_address:port -p password -l listen_address:port -i socks5_server_ip
```

```
# Run as brook wsclient, connect to brook wsserver with domain, start a http proxy
$ brook wsclient -s wss://wsserver_domain:port -p password -l listen_address:port --http
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
$ brook relay -l listen_address:port -r relay_to_address:port
```

> More parameters: $ brook relay -h

### Socks5 Server

```
# Run as standard socks5 server
$ brook socks5 -l listen_address:port -i server_ip
```

> More parameters: $ brook socks5 -h

### Socks5 to http

```
# Convert socks5 proxy to http proxy
$ brook socks5tohttp -l listen_address:port -s socks5_server_address:port
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

### System Proxy

```
# Set system PAC proxy
$ brook systemproxy -u pac_url

# Clear system PAC proxy
$ brook systemproxy -r
```

> More parameters: $ brook systemproxy -h

### Shadowsocks

```
# Run as shadowsocks server
$ brook ssserver -l listen_address:port -p password
```

> More parameters: $ brook ssserver -h

```
# Run as shadowsocks client, connect to shadowsocks server, start a socks5 proxy server
$ brook ssclient -s ssserver_address:port -p password -l listen_address:port -i socks5_server_ip
```

> More parameters: $ brook ssclient -h

> Fixed method is aes-256-cfb

### How to

Some useful tutorial resources

* Brook Wiki: https://github.com/txthinking/brook/wiki
* Brook Issues: https://github.com/txthinking/brook/issues
* Slides: https://talks.txthinking.com
* Youtube: https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A
* Telegram Group: https://t.me/brookgroup
* Telegram Channel: https://t.me/brookchannel
* Nami: https://github.com/txthinking/nami
* Joker: https://github.com/txthinking/joker

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
