# Brook

[🇨🇳 中文](README_ZH.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Donate](https://img.shields.io/badge/Support-Donate-ff69b4.svg)](https://github.com/sponsors/txthinking)
[![Talks](https://img.shields.io/badge/talks-txthinking-blue.svg)](https://talks.txthinking.com)
[![Youtube](https://img.shields.io/badge/Youtube-@txthinking-blue.svg)](https://www.youtube.com/txthinking)
[![Document](https://img.shields.io/badge/github.io-Document-blue.svg)](https://txthinking.github.io/brook/)
[![Telegram](https://img.shields.io/badge/telegram-@txthinking_talks-blue.svg)](https://t.me/txthinking_talks)
[![Telegram](https://img.shields.io/badge/telegram-@brookgroup-blue.svg)](https://t.me/brookgroup)

---

## What is Brook

Brook is a cross-platform strong encryption and not detectable proxy.<br/>
Brook's goal is to keep it **simple**, **stupid** and **not detectable**.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install CLI via curl, and run `brook server`

```
curl -L https://github.com/txthinking/brook/releases/latest/download/brook_linux_amd64 -o /usr/bin/brook
chmod +x /usr/bin/brook
brook server --listen :9999 --password hello
```

> Then, your `brook server` is `YOUR_SERVER_IP:9999`, password is `hello`

### Install CLI and [joker](https://github.com/txthinking/joker) via curl, and run `brook server` with [joker](https://github.com/txthinking/joker)

```
curl -L https://github.com/txthinking/joker/releases/latest/download/joker_linux_amd64 -o /usr/bin/joker
chmod +x /usr/bin/joker
curl -L https://github.com/txthinking/brook/releases/latest/download/brook_linux_amd64 -o /usr/bin/brook
chmod +x /usr/bin/brook
joker brook server --listen :9999 --password hello
```

> Then, your `brook server` is `YOUR_SERVER_IP:9999`, password is `hello`

### Install CLI via [nami](https://github.com/txthinking/nami), and run `brook server` with [joker](https://github.com/txthinking/joker)

```
bash <(curl -s https://bash.ooo/nami.sh)
nami install joker brook
joker brook server --listen :9999 --password hello
```

> Then, your `brook server` is `YOUR_SERVER_IP:9999`, password is `hello`

### Install via script

```
bash <(curl -s https://bash.ooo/brook.sh)
```

[Read Document](https://txthinking.github.io/brook/#/install-cli)

### Install GUI

[Read Document](https://txthinking.github.io/brook/#/install-gui-client)

## Usage

```
NAME:
   Brook - A cross-platform strong encryption and not detectable proxy

USAGE:
   brook [global options] command [command options] [arguments...]

VERSION:
   20210701

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   server          Run as brook server, both TCP and UDP
   client          Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> brook client <-> brook server <-> dst]
   wsserver        Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server
   wsclient        Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> brook wsclient <-> brook wsserver <-> dst]
   wssserver       Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server
   wssclient       Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> brook wssclient <-> brook wssserver <-> dst]
   relayoverbrook  Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> brook server/wsserver/wssserver <-> to address]
   dns             Run as dns server over brook, both TCP and UDP, [src <-> brook dns <-> brook server/wsserver/wssserver <-> dns] or [src <-> brook dns <-> dnsForBypass]
   tproxy          Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> brook tproxy <-> brook server/wsserver/wssserver <-> dst]
   link            Print brook link
   qr              Print brook server QR code
   connect         Connect via standard sharing link (brook server & brook wsserver & brook wssserver)
   relay           Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]
   socks5          Run as standalone standard socks5 server, both TCP and UDP
   socks5tohttp    Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]
   hijackhttps     Hijack domains and assume is TCP/TLS/443. Requesting these domains from anywhere in the system will be hijacked . [src <-> brook hijackhttps <-> socks5 server] or [src <-> direct]
   pac             Run as PAC server or save PAC to file
   servers         Run as multiple brook servers
   relays          Run as multiple standalone relays
   map             Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> brook <-> to address]
   howto           Print some useful tutorial resources
   help, h         Shows a list of commands or help for one command

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
