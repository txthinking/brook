# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/brook)](https://goreportcard.com/report/github.com/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/brook/wiki)

<p align="center">
    <img style="float: right;" src="https://storage.googleapis.com/txthinking-file/_/brook_200x200.png" alt="Brook"/>
</p>

---

### New features (v20180227)

* Transparent proxy: $ brook tproxy. [#117](https://github.com/txthinking/brook/issues/117)
* New MacOS/Windows GUI

---

### Table of Contents

* [What is Brook](#what-is-brook)
* [Download](#download)
* [Server](#server)
    * [Brook Server](#brook-server)
    * [Shadowsocks Server](#shadowsocks-server)
* [Client (CLI)](#client-cli)
    * [Brook Client](#brook-client)
    * [Shadowsocks Client](#shadowsocks-client)
* [Contributing](#contributing)
* [License](#license)

## What is Brook

Brook is a cross-platform(Linux/MacOS/Windows/Android/iOS) proxy/vpn software.<br/>
Brook's goal is to reduce the configuration steps. Keep it simple, stupid.

## Download

| Download | Server/Client | OS | Arch | Remark |
| --- | --- | --- | --- | --- |
| [brook](https://github.com/txthinking/brook/releases/download/v20180227/brook) | Server & Client | Linux | amd64 | CLI |
| [brook_linux_386](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_386) | Server & Client | Linux | 386 | CLI |
| [brook_linux_arm64](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_arm64) | Server & Client | Linux | arm64 | CLI |
| [brook_linux_arm5](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_arm5) | Server & Client | Linux | arm5 | CLI |
| [brook_linux_arm6](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_arm6) | Server & Client | Linux | arm6 | CLI |
| [brook_linux_arm7](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_arm7) | Server & Client | Linux | arm7 | CLI |
| [brook_linux_mips](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_mips) | Server & Client | Linux | mips | CLI |
| [brook_linux_mipsle](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_mipsle) | Server & Client | Linux | mipsle | CLI |
| [brook_linux_mips64](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_mips64) | Server & Client | Linux | mips64 | CLI |
| [brook_linux_mips64le](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_mips64le) | Server & Client | Linux | mips64le | CLI |
| [brook_linux_ppc64](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_ppc64) | Server & Client | Linux | ppc64 | CLI |
| [brook_linux_ppc64le](https://github.com/txthinking/brook/releases/download/v20180227/brook_linux_ppc64le) | Server & Client | Linux | ppc64le | CLI |
| [brook_macos_amd64](https://github.com/txthinking/brook/releases/download/v20180227/brook_macos_amd64) | Server & Client | MacOS | amd64 | CLI |
| [brook_windows_amd64.exe](https://github.com/txthinking/brook/releases/download/v20180227/brook_windows_amd64.exe) | Server & Client | Windows | amd64 | CLI |
| [brook_windows_386.exe](https://github.com/txthinking/brook/releases/download/v20180227/brook_windows_386.exe) | Server & Client | Windows | 386 | CLI |
| [Brook.dmg](https://github.com/txthinking/brook/releases/download/v20180227/Brook.dmg) | Client | MacOS | amd64 | GUI |
| [Brook Setup.exe](https://github.com/txthinking/brook/releases/download/v20180227/Brook%20Setup.exe) | Client | Windows | amd64 | GUI |
| [App Store](https://itunes.apple.com/us/app/brook-brook-shadowsocks-vpn-proxy/id1216002642) | Client | iOS | - | GUI |
| [Google Play](https://play.google.com/store/apps/details?id=com.txthinking.brook) / [Brook.apk](https://github.com/txthinking/brook/releases/download/v20180227/Brook.apk) | Client | Android | - | GUI |

**See [wiki](https://github.com/txthinking/brook/wiki) for more tutorials**

## Server

```
NAME:
   Brook - A Cross-Platform Proxy Software

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20180227

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
     server         Run as server mode
     servers        Run as multiple servers mode
     client         Run as client mode
     tunnel         Run as tunnel mode
     tproxy         Run as tproxy mode, transparent proxy, only works on Linux
     streamserver   Run as server mode
     streamservers  Run as multiple servers mode
     streamclient   Run as client mode
     ssserver       Run as shadowsocks server mode, fixed method is aes-256-cfb
     ssservers      Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb
     ssclient       Run as shadowsocks client mode, fixed method is aes-256-cfb
     socks5         Run as raw socks5 server
     relay          Run as relay mode
     relays         Run as multiple relays mode
     qr             Print brook server QR code
     socks5tohttp   Convert socks5 to http proxy
     systemproxy    Set system proxy with pac url, or remove, only works on MacOS/Windows
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug
   --listen value, -l value  Listen address for debug (default: ":6060")
   --help, -h                show help
   --version, -v             print the version
```

#### Brook Server

```
# Run as a brook server
$ brook server -l :9999 -p password
```

```
# Run as a brook stream server
$ brook streamserver -l :9999 -p password
```

```
# Run as multiple brook servers
$ brook servers \
        -l ":9999 password" \
        -l ":8888 password"
```

#### Shadowsocks Server

```
# Run as a shadowsocks server
$ brook ssserver -l :9999 -p password
```

```
# Run as multiple shadowsocks servers
$ brook ssservers \
        -l ":9999 password" \
        -l ":8888 password"
```

Fixed method is aes-256-cfb

> If you run a public/shared server, do not forget this parameter --tcpDeadline

## Client (CLI)

#### Brook Client

```
# Run as brook client, start a socks5 proxy
$ brook client -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as brook client, start a http(s) proxy
$ brook client -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password --http
```

```
# Run as brook stream client, start a socks5 proxy
$ brook streamclient -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as brook stream client, start a http(s) proxy
$ brook streamclient -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password --http
```


#### Shadowsocks Client

```
# Run as shadowsocks client, start a socks5 proxy
$ brook ssclient -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as shadowsocks client, start a http(s) proxy
$ brook ssclient -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password --http
```

**See [wiki](https://github.com/txthinking/brook/wiki) for more tutorials**

#### Contributing

* Please create PR on `develop` branch

## License

Licensed under The GPLv3 License
