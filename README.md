# Brook

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook) [![Go Report Card](https://goreportcard.com/badge/github.com/txthinking/brook)](https://goreportcard.com/report/github.com/txthinking/brook) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0) [![Slack](https://img.shields.io/badge/join-slack-red.svg)](https://brook-proxy.herokuapp.com)

<p align="center">
    <img style="float: right;" src="https://dn-txthinking.qbox.me/tmp/logo200.png" alt="Brook"/>
</p>

### Table of Content

* [What is Brook](#what-is-brook)
* [Server](#server)
    * [Brook Server](#brook-server)
    * [Shadowsocks Server](#shadowsocks-server)
* [Client](#client)
    * [Linux Client](#linux-client)
    * [MacOS Client](#macos-client)
    * [Windows Client](#windows-client)
    * [Android Client](#android-client)
    * [iOS Client](#ios-client)
* [Advanced Usage](#advanced-usage)
    * [Relay Server](#relay-server)
* [Developer](#developer)
* [License](#license)

## What is Brook

Brook is a cross-platform(Linux/MacOS/Windows/Android/iOS) proxy/vpn software

## Server

### [Download Server for Linux (amd64)](https://github.com/txthinking/brook/releases/download/v20170316/brook) [version: 20170316]

```
NAME:
   Brook - A Cross-Platform Proxy Software

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20170316

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
     relays     Run as multi relays mode
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug, more logs
   --listen value, -l value  Listen address for debug (default: "0.0.0.0:6060")
   --help, -h                show help
   --version, -v             print the version
```

### Brook Server

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

### Shadowsocks Server

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

## Client

### Linux Client

#### Download (Same as the above server link, server and client are in one on linux)

#### Brook Client

```
# Run as brook client
$ brook bkclient -l 127.0.0.1:1080 -s server_address:port -p password
```

```
# Run as brook client with music, music must be same as server's
$ brook bkclient -l 127.0.0.1:1080 -s server_address:port -p password -m muisc_name
```

#### Shadowsocks Client

```
# Run as shadowsocks client
$ brook ssclient -l 127.0.0.1:1080 -s server_address:port -p password
```

### MacOS Client

#### [Download Client for MacOS (amd64)](https://github.com/txthinking/brook/releases/download/v20170316/Brook.app.zip) [version: 20170316]

* Need MacOS version > 10.12
* If MacOS notice it is from an unidentified developer, then go `System Preferences` -> `Security & Privacy`, click Open Anyway
* You may prefer to copy Brook.app to Application folder
* This client use this [pac white list](https://github.com/txthinking/pac) auto proxy rule

### Windows Client

#### [Download Client for Windows (amd64)](https://github.com/txthinking/brook/releases/download/v20170316/Brook.exe) [version: 20170316]

* Need Windows version > 7
* Please like to use Chrome browser
* This client use this [pac white list](https://github.com/txthinking/pac) auto proxy rule

### Android Client

#### [Download Client for Android on Google Play](https://play.google.com/store/apps/details?id=com.txthinking.brook)
#### [Download Client for Android (apk)](https://github.com/txthinking/brook/releases/download/v20170316/Brook.apk)

* This client use this [app white list](https://github.com/txthinking/pac/blob/master/white_apps.list)  auto proxy rule

### iOS Client

#### [Download Client for iOS on AppStore](https://play.google.com/store/apps/details?id=com.txthinking.brook)

* This client use this [pac white list](https://github.com/txthinking/pac) auto proxy rule

## Advanced Usage

### Relay Server

#### What is Relay Server

```
client <---> relay server <---> server
```

#### Relay Server

```
# Run as a relay server
$ brook relay -l :9999 -s server_address:port -t 10
```

```
# Run as multiple relay servers
$ brook relays \
        -l ":9999 -s server1_address:port" \
        -l ":8888 -s server2_address:port" \
        -t 10
```

## Developer

```
$ go get github.com/txthinking/brook/cli/brook
$ brook -h
```

### Contributing

* Please create PR on `develop` branch

## License

Licensed under The GPLv3 License
