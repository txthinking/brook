# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/brook)](https://goreportcard.com/report/github.com/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/brook/wiki)

<p align="center">
    <img style="float: right;" src="https://storage.googleapis.com/txthinking-file/_/brook_200x200.png" alt="Brook"/>
</p>

---

### v20190601

* New macOS/Windows GUI client.

---

### Table of Contents

* [What is Brook](#what-is-brook)
* [Download](#download)
* [Packages](#packages)
* [**Server**](#server)
* [**Client (CLI)**](#client-cli)
* [**Client (GUI)**](#client-gui)
* [Tunnel](#tunnel)
* [Tproxy](#tproxy)
* [VPN](#vpn)
* [Relay](#relay)
* [Socks5](#socks5)
* [Socks5 to HTTP](#socks5-to-http)
* [Shadowsocks](#shadowsocks)
* [Contributing](#contributing)
* [License](#license)

## What is Brook

Brook is a cross-platform proxy/vpn software.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

## Download

| Download | Server/Client | OS | Arch | Remark |
| --- | --- | --- | --- | --- |
| [brook](https://github.com/txthinking/brook/releases/download/v20190601/brook) | Server & Client | Linux | amd64 | CLI |
| [brook_linux_386](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_386) | Server & Client | Linux | 386 | CLI |
| [brook_linux_arm64](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_arm64) | Server & Client | Linux | arm64 | CLI |
| [brook_linux_arm5](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_arm5) | Server & Client | Linux | arm5 | CLI |
| [brook_linux_arm6](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_arm6) | Server & Client | Linux | arm6 | CLI |
| [brook_linux_arm7](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_arm7) | Server & Client | Linux | arm7 | CLI |
| [brook_linux_mips](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_mips) | Server & Client | Linux | mips | CLI |
| [brook_linux_mipsle](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_mipsle) | Server & Client | Linux | mipsle | CLI |
| [brook_linux_mips64](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_mips64) | Server & Client | Linux | mips64 | CLI |
| [brook_linux_mips64le](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_mips64le) | Server & Client | Linux | mips64le | CLI |
| [brook_linux_ppc64](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_ppc64) | Server & Client | Linux | ppc64 | CLI |
| [brook_linux_ppc64le](https://github.com/txthinking/brook/releases/download/v20190601/brook_linux_ppc64le) | Server & Client | Linux | ppc64le | CLI |
| [brook_darwin_amd64](https://github.com/txthinking/brook/releases/download/v20190601/brook_darwin_amd64) | Server & Client | MacOS | amd64 | CLI |
| [brook_windows_amd64.exe](https://github.com/txthinking/brook/releases/download/v20190601/brook_windows_amd64.exe) | Server & Client | Windows | amd64 | CLI |
| [brook_windows_386.exe](https://github.com/txthinking/brook/releases/download/v20190601/brook_windows_386.exe) | Server & Client | Windows | 386 | CLI |
| [Brook.dmg](https://github.com/txthinking/brook/releases/download/v20190601/Brook.dmg) | Client | MacOS | amd64 | GUI |
| [Brook.exe](https://github.com/txthinking/brook/releases/download/v20190601/Brook.exe) | Client | Windows | amd64 | GUI |
| [App Store](https://itunes.apple.com/us/app/brook-brook-shadowsocks-vpn-proxy/id1216002642) | Client | iOS | - | GUI |
| [Brook.apk](https://github.com/txthinking/brook/releases/download/v20190601/Brook.apk)(No Google Play) | Client | Android | - | GUI |

**See [wiki](https://github.com/txthinking/brook/wiki) for more tutorials**

## Packages

### ArchLinux

```
sudo pacman -S brook
```

### MacOS(GUI)

```
brew cask install brook
```

## Brook

```
NAME:
   Brook - A Cross-Platform Proxy/VPN Software

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20190601

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
     server        Run as server mode
     servers       Run as multiple servers mode
     client        Run as client mode
     tunnel        Run as tunnel mode on client-side
     tproxy        Run as tproxy mode on client-side, transparent proxy, only works on Linux
     vpn           Run as VPN mode on client-side
     ssserver      Run as shadowsocks server mode, fixed method is aes-256-cfb
     ssservers     Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb
     ssclient      Run as shadowsocks client mode, fixed method is aes-256-cfb
     socks5        Run as raw socks5 server
     relay         Run as relay mode
     relays        Run as multiple relays mode
     link          Print brook link
     qr            Print brook server QR code
     socks5tohttp  Convert socks5 to http proxy
     systemproxy   Set system proxy with pac url, or remove, only works on MacOS/Windows
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug
   --listen value, -l value  Listen address for debug (default: ":6060")
   --help, -h                show help
   --version, -v             print the version
```

### Server

```
# Run as a brook server
$ brook server -l :9999 -p password
```

```
# Run as multiple brook servers
$ brook servers -l ":9999 password" -l ":8888 password"
```

> If you run a public/shared server, do not forget this parameter --tcpDeadline

### Client (CLI)

```
# Run as brook client, start a socks5 proxy socks5://127.0.0.1:1080
$ brook client -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as brook client, start a http(s) proxy http(s)://127.0.0.1:8080
$ brook client -l 127.0.0.1:8080 -i 127.0.0.1 -s server_address:port -p password --http
```

### Client (GUI)

See [wiki](https://github.com/txthinking/brook/wiki)

#### Tunnel

```
# Run as tunnel 127.0.0.1:5 to 1.2.3.4:5
$ brook tunnel -l 127.0.0.1:5 -t 1.2.3.4:5 -s server_address:port -p password
```

#### Tproxy (usually used on Linux router box)

See [wiki](https://github.com/txthinking/brook/wiki/How-to-run-transparent-proxy-on-Linux%3F)

#### VPN

```
# Run as VPN to proxy all TCP/UDP. [ROOT privileges required].
$ sudo brook vpn -l 127.0.0.1:1080 -s server_address:port -p password
```

**See [wiki](https://github.com/txthinking/brook/wiki/How-to-run-VPN-on-Linux,-MacOS-and-Windows%3F) for more tutorials**

#### Relay

```
# Run as relay to 1.2.3.4:5
$ brook relay -l :5 -r 1.2.3.4:5
```

#### Socks5

```
# Run as a raw socks5 server 1.2.3.4:1080
$ brook socks5 -l :1080 -i 1.2.3.4
```

#### Socks5 to HTTP

```
# Convert socks5://127.0.0.1:1080 to http(s)://127.0.0.1:8080 proxy
$ brook socks5tohttp -l 127.0.0.1:8080 -s 127.0.0.1:1080
```

#### Shadowsocks

```
# Run as a shadowsocks server
$ brook ssserver -l :9999 -p password
```

```
# Run as multiple shadowsocks servers
$ brook ssservers -l ":9999 password" -l ":8888 password"
```

> If you run a public/shared server, do not forget this parameter --tcpDeadline

```
# Run as shadowsocks client, start a socks5 proxy socks5://127.0.0.1:1080
$ brook ssclient -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as shadowsocks client, start a http(s) proxy http(s)://127.0.0.1:8080
$ brook ssclient -l 127.0.0.1:8080 -i 127.0.0.1 -s server_address:port -p password --http
```

> Fixed method is aes-256-cfb

**See [wiki](https://github.com/txthinking/brook/wiki) for more tutorials**

## Contributing

Please read [CONTRIBUTING.md](https://github.com/txthinking/brook/blob/master/.github/CONTRIBUTING.md) first

## License

Licensed under The GPLv3 License
