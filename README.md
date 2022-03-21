# Brook

<!--G-R3M673HK5V-->

[🇨🇳 中文](README_ZH.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[🤝 Telegram](https://t.me/brookgroup)
[💬 Private](https://join.txthinking.com)
[🩸 Youtube](https://www.youtube.com/txthinking)
[❤️ Sponsor](https://github.com/sponsors/txthinking)

Brook is a cross-platform strong encryption and not detectable proxy. Keep it simple, stupid.

[🗣 Subscribe Announcement](https://t.me/txthinking_news)

<!--TOC-->

## Install

### Install brook command

> [nami](https://github.com/txthinking/nami) can automatically download the command corresponding to your system. If on Windows, run in [Git Bash](https://gitforwindows.org)<br/>
> or<br/>
> If your system is not Linux, MacOS, Windows, or don't want nami, you can download it directly on the [releases](https://github.com/txthinking/brook/releases) page<br/>
> or<br/>
> the script but only some parameters are supported: `bash <(curl https://bash.ooo/brook.sh)`<br/>
> or<br/>
> Archlinux: `pacman -S brook`<br/>
> or<br/>
> brew: `brew install brook`<br/><br/>
> recommend run command with root<br/>

Install nami

```
bash <(curl https://bash.ooo/nami.sh)
```

Install brook

```
nami install brook
```

### Install Brook GUI client

[macOS](https://github.com/txthinking/brook/releases/latest/download/Brook.dmg)
[Windows](https://github.com/txthinking/brook/releases/latest/download/Brook.exe)
[Android](https://github.com/txthinking/brook/releases/latest/download/Brook.apk)
[iOS & M1 Mac](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)
[OpenWrt](#gui-for-official-openwrt)

> Windows requires that the latest version of Edge(chromium-based) has been installed<br/>

这里有 [brook GUI 客户端工作原理](https://talks.txthinking.com/articles/brook.article)

## brook `subcommand` and `command line arguments`

-   all `subcoommand`: `brook --help`
-   command line arguments of `subommand`: `brook xxx --help`

## brook rule format

There are three types of rule files

-   domain list: One domain name per line, the suffix matches mode. Can be a local file or an HTTPS URL
-   CIDR v4 list: One CIDR per line, which can be a local file or an HTTPS URL
-   CIDR v6 list: One CIDR per line, which can be a local file or an HTTPS URL

Rules file can be used for

-   Server-side: blocking domain name and IP
-   brook dns: bypass, block domain
-   brook tproxy: bypass, block, domain, ip
-   OpenWrt: bypass, block, domain, ip
-   Brook GUI: bypass, block, domain, ip

## Examples

List some examples of common scene commands, pay attention to replace the parameters such as IP, port, password, domain name, certificate path, etc. in the example by yourself

### Run brook server

```
SRC --TCP--> brook client/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook server --TCP--> DST
SRC --UDP--> brook client/relayoverbrook/dns/tproxy/GUI Client --UDP/TCP(Brook Protocol)--> brook server --UDP--> DST
```

```
brook server --listen :9999 --password hello
```

Get brook link with `--udpovertcp`

```
brook link --server 1.2.3.4:9999 --password hello --udpovertcp --name 'my brook server'
```

or get brook link with udp over udp

> Make sure you have no problem with your local UDP network to your server

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

### Run brook wsserver

```
SRC --TCP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --TCP--> DST
SRC --UDP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --UDP--> DST
```

```
brook wsserver --listen :9999 --password hello
```

Get brook link

```
brook link --server ws://1.2.3.4:9999 --password hello --name 'my brook wsserver'
```

or get brook link with domain, even if that's not your domain

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

### Run brook wssserver: automatically certificate

> Make sure your domain has been resolved to your server IP successfully. Automatic certificate issuance requires the use of port 80

```
brook wssserver --domainaddress domain.com:443 --password hello
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### Run brook wssserver Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### Run brook wssserver issue untrusted certificates yourself, any domain

Install [mad](https://github.com/txthinking/mad)

```
nami install mad
```

Generate root ca

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

Generate domain cert by root ca

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

Get brook link with `--insecure`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --insecure
```

or get brook link with `--ca`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --ca /root/ca.pem
```

### withoutBrookProtocol

Better performance, but data is not strongly encrypted using Brook protocol. So please use certificate encryption, and it is not recommended to use --withoutBrookProtocol and --insecure together

### withoutBrookProtocol automatically certificate

> Make sure your domain has been resolved to your server IP successfully. Automatic certificate issuance requires the use of port 80

```
brook wssserver --domainaddress domain.com:443 --password hello --withoutBrookProtocol
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

### withoutBrookProtocol Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --withoutBrookProtocol
```

### withoutBrookProtocol issue untrusted certificates yourself, any domain

Install [mad](https://github.com/txthinking/mad)

```
nami install mad
```

Generate root ca

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

Generate domain cert by root ca

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

Run brook wssserver

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol --address 1.2.3.4:443 --ca /root/ca.pem
```

### brook server wsserver wssserver forward to another socks5 server on server-side

-   --toSocks5
-   --toSocks5Username
-   --toSocks5Password

### brook server wsserver wssserver block domain and ip on server-side

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

### Run brook socks5, A stand-alone standard socks5 server

```
SRC --TCP--> brook socks5 --TCP--> DST
SRC --UDP--> brook socks5 --UDP--> DST
```

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

Get brook link

```
brook link --server socks5://1.2.3.4:1080
```

### Run brook socks5 with username and password. A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

Get brook link

```
brook link --server socks5://1.2.3.4:1080 --username hello --password world
```

### brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with brook server wsserver wssserver.

```
SRC --TCP--> brook relayoverbrook --TCP(Brook Protocol) --> brook server/wsserver/wssserver --TCP--> DST
SRC --UDP--> brook relayoverbrook --TCP/UDP(Brook Protocol) --> brook server/wsserver/wssserver --UDP--> DST
```

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

### brook dns can create a encrypted DNS server, both TCP and UDP, it works with brook server wsserver wssserver.

```
SRC --TCP--> brook dns --TCP(Brook Protocol) --> brook server/wsserver/wssserver --TCP--> DST
SRC --UDP--> brook dns --TCP/UDP(Brook Protocol) --> brook server/wsserver/wssserver --UDP--> DST
```

```
brook dns ... --listen 127.0.0.1:53
```

Rule

-   --dns
-   --dnsForBypass
-   --bypassDomainList
-   --blockDomainList

### brook tproxy Transparent Proxy Gateway on official OpenWrt

```
opkg install ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy
```

```
brook tproxy --link 'brook://...' --dnsListen :5353
```

1. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
2. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
3. By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

Rule

-   --dnsForDefault
-   --dnsForBypass
-   --bypassDomainList
-   --bypassCIDR4List
-   --bypassCIDR6List
-   --blockDomainList

### brook tproxy Transparent Proxy Gateway on Ubuntu

```
systemctl stop systemd-resolved
```

```
brook tproxy --link 'brook://...' --dnsListen :53
```

1. You may need to manually configure the computer or mobile gateway and DNS.

### GUI for official OpenWrt

> **Dependencies: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy**

port 9999, 1080, 5353 will be used. It work with brook server, brook wsserver and brook wssserver.

1. Download the [ipk](https://github.com/txthinking/brook/releases) file for your router
2. Upload and install: OpenWrt Web -> System -> Software -> Upload Package...
3. Refresh page, the Brook menu will appear at the top
4. OpenWrt Web -> Brook -> type and Connect
5. And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
6. And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

### brook relay can relay a address to a remote address. It can relay any tcp and udp server

```
SRC --TCP--> brook relay --TCP--> DST
SRC --UDP--> brook relay --UDP--> DST
```

```
brook relay --from :9999 --to 1.2.3.4:9999
```

### brook socks5tohttp can convert a socks5 to a http proxy

```
brook socks5tohttp --socks5 127.0.0.1:1080 --listen 127.0.0.1:8010
```

### brook pac creates pac server

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ...
```

### brook pac creates pac file

```
brook pac --file proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ...
```

### IPv6

| Command/Client       | Remark | Support IPv4 | Support IPv6 |
| -------------------- | ------ | ------------ | ------------ |
| brook server         | CLI    | Yes          | Yes          |
| brook client         | CLI    | Yes          | Yes          |
| brook wsserver       | CLI    | Yes          | Yes          |
| brook wsclient       | CLI    | Yes          | Yes          |
| brook wssserver      | CLI    | Yes          | Yes          |
| brook wssclient      | CLI    | Yes          | Yes          |
| brook relayoverbrook | CLI    | Yes          | Yes          |
| brook dns            | CLI    | Yes          | Yes          |
| brook tproxy         | CLI    | Yes          | Yes/2        |
| brook connect        | CLI    | Yes          | Yes          |
| brook relay          | CLI    | Yes          | Yes          |
| brook socks5         | CLI    | Yes          | Yes          |
| brook socks5tohttp   | CLI    | Yes          | Yes          |
| brook hijackhttps    | CLI    | Yes          | Yes          |
| macOS Client         | GUI    | Yes          | Yes          |
| Windows Client       | GUI    | Yes          | Yes/2        |
| iOS Client           | GUI    | Yes          | Yes          |
| Android Client       | GUI    | Yes          | Yes          |
| OpenWrt Client       | GUI    | Yes          | Yes/2        |

### NAT Type

Symmetric

## Run command as daemon via joker

Install [joker](https://github.com/txthinking/joker)

```
nami install joker
```

To run the brook daemon with joker, just prefix the original command with joker

```
joker brook ...
```

Get the last command ID

```
joker last
```

View output and error of a command run via joker

```
joker log <ID>
```

View running commmands via joker

```
joker list
```

Stop a running command via joker

```
joker stop <ID>
```

## Auto start at boot via jinbe

Install [jinbe](https://github.com/txthinking/jinbe)

```
nami install jinbe
```

To use jinbe to add a self-starting command at boot, just add jinbe in front of the original command

```
jinbe joker brook ...
```

View added commmands via jinbe

```
jinbe list
```

Remove a added command via jinbe

```
jinbe remove <ID>
```

## Protocol

### brook server protocol

[brook-server-protocol.md](protocol/brook-server-protocol.md)

### brook wsserver protocol

[brook-wsserver-protocol.md](protocol/brook-wsserver-protocol.md)

### brook wssserver protocol

[brook-wssserver-protocol.md](protocol/brook-wssserver-protocol.md)

### withoutBrookProtocol protocol

[withoutbrookprotocol-protocol.md](protocol/withoutbrookprotocol-protocol.md)

### brook link protocol

[brook-link-protocol.md](protocol/brook-link-protocol.md)

## Resources

-   Brook GUI 工作原理: https://talks.txthinking.com/articles/brook.article
-   brook wsserver and Cloudflare CDN: https://www.youtube.com/watch?v=KFzS55bUk6A
-   用 nico 将 brook wsserver 包装成任意 https 网站: https://talks.txthinking.com/articles/nico-brook-wsserver.article
-   Brook, Shadowsocks, V2ray 协议层面的区别: https://www.youtube.com/watch?v=WZSfZU6rgbQ
-   Blog: https://talks.txthinking.com
-   Youtube: https://www.youtube.com/txthinking
-   Rule list demo: https://github.com/txthinking/bypass
-   https://ipip.ooo
-   Discuss: https://github.com/txthinking/brook/discussions
-   Telegram: https://t.me/brookgroup
-   News: https://t.me/txthinking_news
-   Chrome Extension: [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
