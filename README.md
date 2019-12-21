# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/brook/wiki)

<p align="center">
    <img style="float:right;" src="https://storage.googleapis.com/txthinking-file/_/brook.png" alt="Brook"/>
</p>

---

### v20200102

- **🔊 Please uninstall or delete the old GUI client first❗️**
- Add Brook WebSocket mode, with or without TLS. If with TLS, Brook will automatically request/issue certificate for your domain.
- GUI Client supports QR scanning
- GUI Client supports custom rules
- macOS Client renamed to Brook.pkg (Because many users don't know that needed to copy or drag from dmg)
- Windows Client renamed to Brook.msi
- develop branch is deleted, PR to master and keep master stable
- Keep it simple, stupid

---

### Table of Contents

- [What is Brook](#what-is-brook)
- [Download](#download)
- [Packages](#packages)
- [**Server**](#server)
- [**Client**](#client)
- [WSServer](#wsserver)
- [WSClient](#wsclient)
- [**GUI Client**](#gui-client)
- [Tunnel](#tunnel)
- [Tproxy](#tproxy)
- [VPN](#vpn)
- [Relay](#relay)
- [Socks5](#socks5)
- [Socks5 to HTTP](#socks5-to-http)
- [Shadowsocks](#shadowsocks)
- [Contributing](#contributing)
- [License](#license)

## What is Brook

Brook is a cross-platform proxy/vpn software.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

## Download

| Download                                                                                                                         | Server/Client   | OS      | Arch               | Remark |
| -------------------------------------------------------------------------------------------------------------------------------- | --------------- | ------- | ------------------ | ------ |
| [brook](https://github.com/txthinking/brook/releases/download/v20200102/brook)                                                   | Server & Client | Linux   | amd64              | CLI    |
| [brook_linux_386](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_386)                               | Server & Client | Linux   | 386                | CLI    |
| [brook_linux_arm64](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_arm64)                           | Server & Client | Linux   | arm64              | CLI    |
| [brook_linux_arm5](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_arm5)                             | Server & Client | Linux   | arm5               | CLI    |
| [brook_linux_arm6](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_arm6)                             | Server & Client | Linux   | arm6               | CLI    |
| [brook_linux_arm7](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_arm7)                             | Server & Client | Linux   | arm7               | CLI    |
| [brook_linux_mips](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips)                             | Server & Client | Linux   | mips               | CLI    |
| [brook_linux_mipsle](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mipsle)                         | Server & Client | Linux   | mipsle             | CLI    |
| [brook_linux_mips_softfloat](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips_softfloat)         | Server & Client | Linux   | mips_softfloat     | CLI    |
| [brook_linux_mipsle_softfloat](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mipsle_softfloat)     | Server & Client | Linux   | mipsle_softfloat   | CLI    |
| [brook_linux_mips64](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips64)                         | Server & Client | Linux   | mips64             | CLI    |
| [brook_linux_mips64le](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips64le)                     | Server & Client | Linux   | mips64le           | CLI    |
| [brook_linux_mips64_softfloat](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips64_softfloat)     | Server & Client | Linux   | mips64_softfloat   | CLI    |
| [brook_linux_mips64le_softfloat](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_mips64le_softfloat) | Server & Client | Linux   | mips64le_softfloat | CLI    |
| [brook_linux_ppc64](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_ppc64)                           | Server & Client | Linux   | ppc64              | CLI    |
| [brook_linux_ppc64le](https://github.com/txthinking/brook/releases/download/v20200102/brook_linux_ppc64le)                       | Server & Client | Linux   | ppc64le            | CLI    |
| [brook_darwin_amd64](https://github.com/txthinking/brook/releases/download/v20200102/brook_darwin_amd64)                         | Server & Client | macOS   | amd64              | CLI    |
| [brook_windows_amd64.exe](https://github.com/txthinking/brook/releases/download/v20200102/brook_windows_amd64.exe)               | Server & Client | Windows | amd64              | CLI    |
| [brook_windows_386.exe](https://github.com/txthinking/brook/releases/download/v20200102/brook_windows_386.exe)                   | Server & Client | Windows | 386                | CLI    |
| [Brook.pkg](https://github.com/txthinking/brook/releases/download/v20200102/Brook.pkg)                                           | Client          | macOS   | amd64              | GUI    |
| [Brook.msi](https://github.com/txthinking/brook/releases/download/v20200102/Brook.msi)                                           | Client          | Windows | amd64              | GUI    |
| [App Store](https://itunes.apple.com/us/app/brook-brook-shadowsocks-vpn-proxy/id1216002642)                                      | Client          | iOS     | -                  | GUI    |
| [Brook.apk](https://github.com/txthinking/brook/releases/download/v20200102/Brook.apk)(No Google Play)                           | Client          | Android | -                  | GUI    |

**See [wiki](https://github.com/txthinking/brook/wiki) for more tutorials**

## Packages

### ArchLinux

```
sudo pacman -S brook
```

### macOS(GUI)

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
   20200101

COMMANDS:
   server        Run as server mode
   servers       Run as multiple servers mode
   client        Run as client mode
   wsserver      Run as websocket server mode
   wsclient      Run as websocket client mode
   tunnel        Run as tunnel mode on client-site
   tproxy        Run as tproxy mode on client-site, transparent proxy, only works on Linux
   vpn           Run as VPN mode on client-site
   ssserver      Run as shadowsocks server mode, fixed method is aes-256-cfb
   ssservers     Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb
   ssclient      Run as shadowsocks client mode, fixed method is aes-256-cfb
   socks5        Run as raw socks5 server
   relay         Run as relay mode
   relays        Run as multiple relays mode
   link          Print brook link
   qr            Print brook server QR code
   socks5tohttp  Convert socks5 to http proxy
   systemproxy   Set system proxy with pac url, or remove, only works on macOS/Windows
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
$ brook server -l :9999 -p password
```

```
# Run as multiple brook servers
$ brook servers -l ":9999 password" -l ":8888 password"
```

> If you run a public/shared server, do not forget this parameter --tcpDeadline

### Client

```
# Run as brook client, start a socks5 proxy socks5://127.0.0.1:1080
$ brook client -l 127.0.0.1:1080 -i 127.0.0.1 -s server_address:port -p password
```

```
# Run as brook client, start a http(s) proxy http(s)://127.0.0.1:8080
$ brook client -l 127.0.0.1:8080 -i 127.0.0.1 -s server_address:port -p password --http
```

### WSServer

```
# Run as a brook wsserver
$ brook wsserver -l :9999 -p password
```

```
# Run as a brook wsserver with domain
# Make sure your domain name has been successfully resolved, 80 and 443 are open, brook will automatically issue certificate for you
$ brook wsserver --domain txthinking.com -p password
```

> If you run a public/shared server, do not forget this parameter --tcpDeadline

### WSClient

```
# Run as brook wsclient, connect brook wsserver
$ brook wsclient -l 127.0.0.1:1080 -i 127.0.0.1 -s ws://1.2.3.4:5 -p password
```

```
# Run as brook wsclient, connect brook wsserver with tls
$ brook wsclient -l 127.0.0.1:1080 -i 127.0.0.1 -s wss://txthinking.com:443 -p password
```

### GUI Client

See [download](https://github.com/txthinking/brook#download)

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

# Must exit by Ctrl+C
```

**See [wiki](https://github.com/txthinking/brook/wiki/How-to-run-VPN-on-Linux,-macOS-and-Windows%3F) for more tutorials**

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
