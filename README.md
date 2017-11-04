# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/brook)](https://goreportcard.com/report/github.com/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/brook/wiki)

<p align="center">
    <img style="float: right;" src="https://storage.googleapis.com/txthinking-file/_/brook_200x200.png" alt="Brook"/>
</p>

---

### Table of Contents

* [What is Brook](#what-is-brook)
* [Download](#download)
* [Server](#server)
    * [Brook Server](#brook-server)
    * [Shadowsocks Server](#shadowsocks-server)
    * [Run as Daemon](#run-as-daemon)
    * [Relay Server](#relay-server)
* [Client (CLI)](#client-cli)
    * [Brook Client](#brook-client)
    * [Shadowsocks Client](#shadowsocks-client)
* [Developer](#developer)
* [License](#license)

## What is Brook

Brook is a cross-platform(Linux/MacOS/Windows/Android/iOS) proxy/vpn software.<br/>
Brook's goal is to reduce the configuration steps. Keep it simple, stupid.

## Download

| Download | Server/Client | OS | Arch | Remark |
| --- | --- | --- | --- | --- |
| [brook](https://github.com/txthinking/brook/releases/download/v20171122/brook) | Server & Client | Linux | amd64 | CLI |
| [brook_linux_386](https://github.com/txthinking/brook/releases/download/v20171122/brook_linux_386) | Server & Client | Linux | 386 | CLI |
| [brook_linux_arm64](https://github.com/txthinking/brook/releases/download/v20171122/brook_linux_arm64) | Server & Client | Linux | arm64 | CLI |
| [brook_linux_arm5](https://github.com/txthinking/brook/releases/download/v20171122/brook_linux_arm5) | Server & Client | Linux | arm5 | CLI |
| [brook_linux_arm6](https://github.com/txthinking/brook/releases/download/v20171122/brook_linux_arm6) | Server & Client | Linux | arm6 | CLI |
| [brook_linux_arm7](https://github.com/txthinking/brook/releases/download/v20171122/brook_linux_arm7) | Server & Client | Linux | arm7 | CLI |
| [brook_macos_amd64](https://github.com/txthinking/brook/releases/download/v20171122/brook_macos_amd64) | Server & Client | MacOS | amd64 | CLI |
| [brook_windows_amd64.exe](https://github.com/txthinking/brook/releases/download/v20171122/brook_windows_amd64.exe) | Server & Client | Windows | amd64 | CLI |
| [brook_windows_386.exe](https://github.com/txthinking/brook/releases/download/v20171122/brook_windows_386.exe) | Server & Client | Windows | 386 | CLI |
| [Brook.app.zip](https://github.com/txthinking/brook/releases/download/v20171122/Brook.app.zip) | Client | MacOS | amd64 | GUI |
| [Brook.exe](https://github.com/txthinking/brook/releases/download/v20171122/Brook.exe) | Client | Windows | amd64 | GUI |
| [Brook.386.exe](https://github.com/txthinking/brook/releases/download/v20171122/Brook.386.exe) | Client | Windows | 386 | GUI |
| [App Store](https://itunes.apple.com/us/app/brook-brook-shadowsocks-vpn-proxy/id1216002642) | Client | iOS | - | GUI |
| [Google Play](https://play.google.com/store/apps/details?id=com.txthinking.brook) / [Brook.apk](https://github.com/txthinking/brook/releases/download/v20171122/Brook.apk) | Client | Android | - | GUI |

MacOS GUI Client

* Need MacOS version >= 10.12
* If MacOS prompts it is from an unidentified developer, then go `System Preferences` -> `Security & Privacy`, click Open Anyway
* You may prefer to copy Brook.app to Application folder
* Follow this [pac white list](https://github.com/txthinking/pac) auto proxy rule

Windows GUI Client

* Need Windows version >= 7
* Please set chrome as your default browser
* You may need to run as an administrator
* Follow this [pac white list](https://github.com/txthinking/pac) auto proxy rule

Android Client

* Follow this [pac white list](https://github.com/txthinking/pac) auto proxy rule
* Not tested on IPv6

iOS Client

* Need iOS version >= 10.0
* Follow this [pac white list](https://github.com/txthinking/pac) auto proxy rule

## Server

```
NAME:
   Brook - A Cross-Platform Proxy Software

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20171122

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
     server         Run as server mode
     servers        Run as multiple servers mode
     client         Run as client mode
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

#### Run as Daemon

With nohup

```
# Start
$ nohup brook server -l :9999 -p password &

# Stop
$ killall brook
```

With systemd

If your linux run with systemd, like Ubuntu 16.04, Archlinux, etc:

```
# Install
$ curl -L git.io/getbrook | sudo bash
$ sudo systemctl daemon-reload

# Config command options
$ sudo vim /etc/default/brook

# Start
$ sudo systemctl start brook.service

# Stop
$ sudo systemctl stop brook.service

# Start on bootup
$ sudo systemctl enable brook.service
```

#### Relay Server

What is Relay Server

```
client <---> relay server <---> server
```

Relay Server

```
# Run as a relay server
$ brook relay -l :9999 -r server_address:port
```

```
# Run as multiple relay servers
$ brook relays \
        -l ":9999 server1_address:port" \
        -l ":8888 server2_address:port"
```

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

#### Contributing

* Please create PR on `develop` branch

## License

Licensed under The GPLv3 License
