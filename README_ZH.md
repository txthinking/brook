# Brook

<!--G-R3M673HK5V-->

[🇬🇧 English](README.md)

[![Build Status](https://travis-ci.org/txthinking/brook.svg?branch=master)](https://travis-ci.org/txthinking/brook)
[![开源协议: GPL v3](https://img.shields.io/badge/%E5%BC%80%E6%BA%90%E5%8D%8F%E8%AE%AE-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

[🤝 Telegram](https://t.me/brookgroup)
[💬 Chat](https://join.txthinking.com)
[🩸 Youtube](https://www.youtube.com/txthinking)
[❤️ Sponsor](https://github.com/sponsors/txthinking)

Brook 是一个跨平台的强加密无特征的代理软件. Keep it simple, stupid.

[🗣 订阅公告](https://t.me/txthinking_news)

<!--TOC-->

## 安装

### 安装 brook 命令

> [nami](https://github.com/txthinking/nami) 她会自动帮你下载适用你系统的最新版 brook 命令文件. 如果你的系统是 Windows, 你需要在 [Git Bash](https://gitforwindows.org) 里面运行<br/>
> 或<br/>
> 如果你的系统不是 Linux, MacOS, Windows, 或不想用 nami, 你可以去 [releases](https://github.com/txthinking/brook/releases) 自己下载命令文件<br/>
> 或<br/>
> 一键脚本，目前仅支持部分命令参数: `bash <(curl https://bash.ooo/brook.sh)`<br/>
> 或<br/>
> Archlinux: `pacman -S brook`<br/>
> 或<br/>
> brew: `brew install brook`<br/><br/>
> 建议使用 root 用户运行命令<br/>

安装 nami

```
bash <(curl https://bash.ooo/nami.sh)
```

使用 nami 安装 brook

```
nami install brook
```

### 安装 Brook 图形客户端

[macOS](https://github.com/txthinking/brook/releases/latest/download/Brook.dmg)
[Windows](https://github.com/txthinking/brook/releases/latest/download/Brook.exe)
[Android](https://github.com/txthinking/brook/releases/latest/download/Brook.apk)
[iOS & M1 Mac](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642)
[OpenWrt](#官网原版-openwrt-图形客户端)

> Windows 需要你已经安装了最新版(基于 chromium 的那款)的 Edge 浏览器<br/>
> 请使用非中国大陆 Apple ID 下载.

这里有 [brook 图形客户端工作原理](https://talks.txthinking.com/articles/brook.article)

## brook `子命令` 和 `命令参数`

-   查看所有的`子命令`: `brook --help`
-   查看某个`子命令`的参数: `brook xxx --help`

## brook 规则格式

有三种规则文件

-   域名列表: 一行一个域名，后缀匹配模式。可以是本地的文件，也可以是 HTTPS URL
-   CIDR v4 列表: 一行一个 CIDR。可以是本地的文件，也可以是 HTTPS URL
-   CIDR v6 列表: 一行一个 CIDR。可以是本地的文件，也可以是 HTTPS URL

规则文件可用于

-   服务端屏蔽域名和 IP
-   brook dns 分流，屏蔽域名
-   brook tproxy 分流，屏蔽域名
-   OpenWrt 分流，屏蔽域名
-   brook 图形客户端分流，屏蔽域名

## 例子

下面列举一些常用场景命令的例子, 注意自己替换示例中的 IP，端口，密码，域名，证书路径等参数

### 运行 brook server

```
SRC --TCP--> brook client/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook server --TCP--> DST
SRC --UDP--> brook client/relayoverbrook/dns/tproxy/GUI Client --UDP/TCP(Brook Protocol)--> brook server --UDP--> DST
```

```
brook server --listen :9999 --password hello
```

获取 brook link 让 udp 走 tcp `--udpovertcp`

```
brook link --server 1.2.3.4:9999 --password hello --udpovertcp --name 'my brook server'
```

或 获取 brook link 让 udp 走 udp

> 确保你本地到你服务器的 UDP 网络没有问题

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

### 运行 brook wsserver

```
SRC --TCP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --TCP--> DST
SRC --UDP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --UDP--> DST
```

```
brook wsserver --listen :9999 --password hello
```

获取 brook link

```
brook link --server ws://1.2.3.4:9999 --password hello --name 'my brook wsserver'
```

或获取 brook link 指定个域名, 甚至不是你自己的域名也可以

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

### 运行 brook wssserver: 自动签发信任证书

> 注意：确保你的域名已成功解析到你服务器的 IP, 自动签发证书需要额外监听 80 端口

```
brook wssserver --domainaddress domain.com:443 --password hello
```

获取 brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### 运行 brook wssserver 使用已有的信任机构签发的证书

> 注意：确保你的域名已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

获取 brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### 运行 brook wssserver 自己签发非信任证书, 甚至不是你自己的域名也可以

安装 [mad](https://github.com/txthinking/mad)

```
nami install mad
```

使用 mad 生成根证书

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

使用 mad 由根证书派发 domain.com 证书

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

获取 brook link 使用 `--insecure`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --insecure
```

或 获取 brook link 使用 `--ca`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --ca /root/ca.pem
```

### withoutBrookProtocol

性能更好，但数据不使用 Brook 协议进行强加密。所以请使用证书加密，并且不建议--withoutBrookProtocol 和--insecure 一起使用

### withoutBrookProtocol 自动签发信任证书

> 注意：确保你的域名已成功解析到你服务器的 IP, 自动签发证书需要额外监听 80 端口

```
brook wssserver --domainaddress domain.com:443 --password hello --withoutBrookProtocol
```

获取 brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

### withoutBrookProtocol 使用已有的信任机构签发的证书

> 注意：确保你的域名已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

获取 brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --withoutBrookProtocol
```

### withoutBrookProtocol 自己签发非信任证书, 甚至不是你自己的域名也可以

安装 [mad](https://github.com/txthinking/mad)

```
nami install mad
```

使用 mad 生成根证书

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

使用 mad 由根证书派发 domain.com 证书

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

运行 brook wssserver

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

获取 brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol --address 1.2.3.4:443 --ca /root/ca.pem
```

### brook server wsserver wssserver 服务端转发给另外的 socks5 server

-   --toSocks5
-   --toSocks5Username
-   --toSocks5Password

### brook server wsserver wssserver 在服务端屏蔽域名和 IP 列表

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

### 运行 brook socks5, 一个独立的标准 socks5 server

```
SRC --TCP--> brook socks5 --TCP--> DST
SRC --UDP--> brook socks5 --UDP--> DST
```

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

获取 brook link

```
brook link --server socks5://1.2.3.4:1080
```

### 运行 brook socks5, 一个独立的标准 socks5 server, 指定用户名和密码

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

获取 brook link

```
brook link --server socks5://1.2.3.4:1080 --username hello --password world
```

### brook relayoverbrook 中继任何 TCP 和 UDP server, 让其走 brook 协议. 它与 brook server wsserver wssserver 一起工作

```
SRC --TCP--> brook relayoverbrook --TCP(Brook Protocol) --> brook server/wsserver/wssserver --TCP--> DST
SRC --UDP--> brook relayoverbrook --TCP/UDP(Brook Protocol) --> brook server/wsserver/wssserver --UDP--> DST
```

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

### brook dns 用来创建一个加密 DNS Server, TCP and UDP, 它与 brook server wsserver wssserver 一起工作

```
SRC --TCP--> brook dns --TCP(Brook Protocol) --> brook server/wsserver/wssserver --TCP--> DST
SRC --UDP--> brook dns --TCP/UDP(Brook Protocol) --> brook server/wsserver/wssserver --UDP--> DST
```

```
brook dns ... --listen 127.0.0.1:53
```

规则

-   --dns
-   --dnsForBypass
-   --bypassDomainList
-   --blockDomainList

### brook tproxy 透明代理网关在官网原版 OpenWrt

```
opkg install ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy
```

```
brook tproxy --link 'brook://...' --dnsListen :5353
```

1. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
2. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
3. 默認, OpenWrt 將會下發 router 的 IP 的為電腦或手機的網關和 DNS

规则

-   --dnsForDefault
-   --dnsForBypass
-   --bypassDomainList
-   --bypassCIDR4List
-   --bypassCIDR6List
-   --blockDomainList

### brook tproxy 透明代理网关在 Ubuntu

```
systemctl stop systemd-resolved
```

```
brook tproxy --link 'brook://...' --dnsListen :53
```

1. You may need to manually configure the computer or mobile gateway and DNS.

### 官网原版 OpenWrt 图形客户端

> **依赖: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy**

**端口 9999, 1080, 5353 将会被使用**. 它与 brook server, brook wsserver, brook wssserver 一起工作.

1. 下載適合你系統的[ipk](https://github.com/txthinking/brook/releases)文件
2. 上傳並安裝: OpenWrt Web -> System -> Software -> Upload Package...
3. 刷新頁面, 頂部菜單會出現 Brook 按鈕
4. OpenWrt Web -> Brook -> 輸入後點擊 Connect
5. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
6. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. 默認, OpenWrt 將會下發 router 的 IP 為電腦或手機的網關和 DNS

### brook relay 可以中继任何 TCP 和 UDP server, 这是一个独立的功能, 它不依赖 brook server wsserver wssserver

```
SRC --TCP--> brook relay --TCP--> DST
SRC --UDP--> brook relay --UDP--> DST
```

```
brook relay --from :9999 --to 1.2.3.4:9999
```

### brook socks5tohttp 将 socks5 proxy 转换为 http proxy

```
brook socks5tohttp --socks5 127.0.0.1:1080 --listen 127.0.0.1:8010
```

### brook pac 创建一个 pac server

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ...
```

### brook pac 创建一个 pac 文件

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

## 使用 joker 运行守护进程

使用 nami 安装 [joker](https://github.com/txthinking/joker)

```
nami install joker
```

使用 joker 运行 brook 守护进程, 只需要在原来命令前面加上 joker

```
joker brook ...
```

查看最后一个命令的 ID

```
joker last
```

查看某个命令的输出和错误

```
joker log <ID>
```

查看运行的命令列表

```
joker list
```

停止某个命令

```
joker stop <ID>
```

## 使用 jinbe 添加开机自启动命令

使用 nami 安装 [jinbe](https://github.com/txthinking/jinbe)

```
nami install jinbe
```

使用 jinbe 添加开机自启动命令，只需要在原来命令前面加上 jinbe

```
jinbe joker brook ...
```

查看添加的开机命令

```
jinbe list
```

移除某个开机命令

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

## 其他资源

-   Brook 工作原理: https://talks.txthinking.com/articles/brook.article
-   brook wsserver 搭配 Cloudflare CDN: https://www.youtube.com/watch?v=KFzS55bUk6A
-   用 nico 将 brook wsserver 包装成任意 https 网站: https://talks.txthinking.com/articles/nico-brook-wsserver.article
-   Brook, Shadowsocks, V2ray 协议层面的区别: https://www.youtube.com/watch?v=WZSfZU6rgbQ
-   博客: https://talks.txthinking.com
-   Youtube 频道: https://www.youtube.com/txthinking
-   规则列表: https://github.com/txthinking/bypass
-   https://ipip.ooo
-   论坛: https://github.com/txthinking/brook/discussions
-   Telegram 群: https://t.me/brookgroup
-   News: https://t.me/txthinking_news
-   Chrome Extension: [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
