# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/brook)](https://goreportcard.com/report/github.com/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/txthinking/brook/wiki)

<p align="center">
    <img style="float: right;" src="https://storage.googleapis.com/txthinking-init/img/logo200.png" alt="Brook"/>
</p>

### Table of Content

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
| [brook](https://github.com/txthinking/brook/releases/download/v20170826/brook) | Server & Client | Linux | amd64 | CLI |
| [brook_linux_386](https://github.com/txthinking/brook/releases/download/v20170826/brook_linux_386) | Server & Client | Linux | 386 | CLI |
| [brook_linux_arm64](https://github.com/txthinking/brook/releases/download/v20170826/brook_linux_arm64) | Server & Client | Linux | arm64 | CLI |
| [brook_linux_arm5](https://github.com/txthinking/brook/releases/download/v20170826/brook_linux_arm5) | Server & Client | Linux | arm5 | CLI |
| [brook_linux_arm6](https://github.com/txthinking/brook/releases/download/v20170826/brook_linux_arm6) | Server & Client | Linux | arm6 | CLI |
| [brook_linux_arm7](https://github.com/txthinking/brook/releases/download/v20170826/brook_linux_arm7) | Server & Client | Linux | arm7 | CLI |
| [brook_macos_amd64](https://github.com/txthinking/brook/releases/download/v20170826/brook_macos_amd64) | Server & Client | MacOS | amd64 | CLI |
| [brook_windows_amd64.exe](https://github.com/txthinking/brook/releases/download/v20170826/brook_windows_amd64.exe) | Server & Client | Windows | amd64 | CLI |
| [brook_windows_386.exe](https://github.com/txthinking/brook/releases/download/v20170826/brook_windows_386.exe) | Server & Client | Windows | 386 | CLI |
| [Brook.app.zip](https://github.com/txthinking/brook/releases/download/v20170826/Brook.app.zip) | Client | MacOS | amd64 | GUI |
| [Brook.exe](https://github.com/txthinking/brook/releases/download/v20170826/Brook.exe) | Client | Windows | amd64 | GUI |
| [Brook.386.exe](https://github.com/txthinking/brook/releases/download/v20170826/Brook.386.exe) | Client | Windows | 386 | GUI |
| [App Store](https://itunes.apple.com/us/app/brook-brook-shadowsocks-vpn-proxy/id1216002642) | Client | iOS | - | GUI |
| [Google Play](https://play.google.com/store/apps/details?id=com.txthinking.brook) / [Brook.apk](https://github.com/txthinking/brook/releases/download/v20170826/Brook.apk) | Client | Android | - | GUI |

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

* Need Android version >= 5.0
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
   20170826

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
     bkserver   Run as brook protocol server mode
     bkservers  Run as brook protocol multiple servers mode
     bkclient   Run as brook protocol client mode
     s5server   Run as socks5 encrypt protocol server mode
     s5servers  Run as socks5 encrypt protocol multiple servers mode
     s5client   Run as socks5 encrypt protocol client mode
     ssserver   Run as shadowsocks protocol server mode, fixed method is aes-256-cfb
     ssservers  Run as shadowsocks protocol multiple servers mode, fixed method is aes-256-cfb
     ssclient   Run as shadowsocks protocol client mode, fixed method is aes-256-cfb
     relay      Run as relay mode
     relays     Run as multiple relays mode
     qr         Print brook server QR code
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug, more logs
   --listen value, -l value  Listen address for debug (default: ":6060")
   --help, -h                show help
   --version, -v             print the version
```

#### Brook Server

```
# Run as a brook server
$ brook bkserver -l :9999 -p password -t 10
```

```
# Run as a brook server with music
$ brook bkserver -l :9999 -p password -t 10 -m music_name
```

```
# Run as multiple brook servers
$ brook bkservers \
        -l ":9999 password" \
        -l ":8888 password" \
        -l ":7777 password music_name" \
        -l ":6666 password music_name" \
        -t 10
```

[More about Brook Music](https://github.com/txthinking/brook/wiki/Music-List)

> If you run a public/shared server, do not forget this parameter --deadline 60 or -d 60

#### Shadowsocks Server

```
# Run as a shadowsocks server
$ brook ssserver -l :9999 -p password -t 10
```

```
# Run as multiple shadowsocks servers
$ brook ssservers \
        -l ":9999 password" \
        -l ":8888 password" \
        -t 10
```

Fixed method is aes-256-cfb

> If you run a public/shared server, do not forget this parameter --deadline 60 or -d 60

#### Run as Daemon

With nohup

```
# Start
$ nohup brook bkserver -l :9999 -p password -t 10 &

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
$ brook relay -l :9999 -s server_address:port -t 10
```

```
# Run as multiple relay servers
$ brook relays \
        -l ":9999 server1_address:port" \
        -l ":8888 server2_address:port" \
        -t 10
```

## Client (CLI)

#### Brook Client

```
# Run as brook client, start a socks5 proxy
$ brook bkclient -l 127.0.0.1:1080 -s server_address:port -p password
```

```
# Run as brook client, start a http(s) proxy
$ brook bkclient -l 127.0.0.1:8080 -s server_address:port -p password --http
```

```
# Run as brook client with music, music must be same as server's
$ brook bkclient -l 127.0.0.1:1080 -s server_address:port -p password -m muisc_name
```

#### Shadowsocks Client

```
# Run as shadowsocks client, start a socks5 proxy
$ brook ssclient -l 127.0.0.1:1080 -s server_address:port -p password
```

```
# Run as shadowsocks client, start a http(s) proxy
$ brook ssclient -l 127.0.0.1:8080 -s server_address:port -p password --http
```

## Developer

```
$ go get github.com/txthinking/brook/cli/brook
$ brook -h
```

#### Contributing

* Please create PR on `develop` branch

## License

Licensed under The GPLv3 License
