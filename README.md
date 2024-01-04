# Brook
A cross-platform programmable network tool

# Sponsor
**❤️  [Shiliew - China Optimized Network App](https://www.txthinking.com/shiliew.html)**

Table of Contents
=================

* [Brook](#brook)
* [Sponsor](#sponsor)
* [Getting Started](#getting-started)
   * [Server](#server)
   * [GUI Client](#gui-client)
   * [CLI Client](#cli-client)
* [GUI Documentation](#gui-documentation)
   * [Software for which this article applies](#software-for-which-this-article-applies)
   * [Programmable](#programmable)
      * [Introduction to incoming variables](#introduction-to-incoming-variables)
      * [in_brooklinks](#in_brooklinks)
      * [in_dnsquery](#in_dnsquery)
      * [in_address](#in_address)
      * [in_httprequest](#in_httprequest)
      * [in_httpresponse](#in_httpresponse)
   * [Write script](#write-script)
   * [Debug script](#debug-script)
   * [Install CA](#install-ca)
* [图形客户端文档](#图形客户端文档)
   * [本文适用的软件](#本文适用的软件)
   * [编程](#编程)
      * [传入变量介绍](#传入变量介绍)
      * [in_brooklinks](#in_brooklinks-1)
      * [in_dnsquery](#in_dnsquery-1)
      * [in_address](#in_address-1)
      * [in_httprequest](#in_httprequest-1)
      * [in_httpresponse](#in_httpresponse-1)
   * [写脚本](#写脚本)
   * [调试脚本](#调试脚本)
   * [安装 CA](#安装-ca)
* [Resources](#resources)
* [CLI Documentation](#cli-documentation)
* [NAME](#name)
* [SYNOPSIS](#synopsis)
* [GLOBAL OPTIONS](#global-options)
* [COMMANDS](#commands)
   * [server](#server-1)
   * [client](#client)
   * [wsserver](#wsserver)
   * [wsclient](#wsclient)
   * [wssserver](#wssserver)
   * [wssclient](#wssclient)
   * [quicserver](#quicserver)
   * [quicclient](#quicclient)
   * [relayoverbrook](#relayoverbrook)
   * [dnsserveroverbrook](#dnsserveroverbrook)
   * [tproxy](#tproxy)
   * [link](#link)
   * [connect](#connect)
   * [relay](#relay)
   * [dnsserver](#dnsserver)
   * [dnsclient](#dnsclient)
   * [dohserver](#dohserver)
   * [dohclient](#dohclient)
   * [dhcpserver](#dhcpserver)
   * [socks5](#socks5)
   * [socks5tohttp](#socks5tohttp)
   * [pac](#pac)
   * [testsocks5](#testsocks5)
   * [testbrook](#testbrook)
   * [echoserver](#echoserver)
   * [echoclient](#echoclient)
   * [ipcountry](#ipcountry)
   * [completion](#completion)
   * [mdpage](#mdpage)
      * [help, h](#help-h)
   * [manpage](#manpage)
   * [help, h](#help-h-1)
* [Diagram](#diagram)
   * [overview](#overview)
   * [withoutBrookProtocol](#withoutbrookprotocol)
   * [relayoverbrook](#relayoverbrook-1)
   * [dnsserveroverbrook](#dnsserveroverbrook-1)
   * [relay](#relay-1)
   * [dnsserver](#dnsserver-1)
   * [tproxy](#tproxy-1)
   * [gui](#gui)
   * [script](#script)
* [Examples](#examples)
      * [Run brook server](#run-brook-server)
      * [Run brook wsserver](#run-brook-wsserver)
      * [Run brook wssserver: automatically certificate](#run-brook-wssserver-automatically-certificate)
      * [Run brook wssserver Use a certificate issued by an existing trust authority](#run-brook-wssserver-use-a-certificate-issued-by-an-existing-trust-authority)
      * [Run brook wssserver issue untrusted certificates yourself, any domain](#run-brook-wssserver-issue-untrusted-certificates-yourself-any-domain)
      * [withoutBrookProtocol](#withoutbrookprotocol-1)
      * [withoutBrookProtocol automatically certificate](#withoutbrookprotocol-automatically-certificate)
      * [withoutBrookProtocol Use a certificate issued by an existing trust authority](#withoutbrookprotocol-use-a-certificate-issued-by-an-existing-trust-authority)
      * [withoutBrookProtocol issue untrusted certificates yourself, any domain](#withoutbrookprotocol-issue-untrusted-certificates-yourself-any-domain)
      * [Run brook socks5, A stand-alone standard socks5 server](#run-brook-socks5-a-stand-alone-standard-socks5-server)
      * [Run brook socks5 with username and password. A stand-alone standard socks5 server](#run-brook-socks5-with-username-and-password-a-stand-alone-standard-socks5-server)
      * [brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with brook server wsserver wssserver.](#brook-relayoverbrook-can-relay-a-local-address-to-a-remote-address-over-brook-both-tcp-and-udp-it-works-with-brook-server-wsserver-wssserver)
      * [brook dnsserveroverbrook can create a encrypted DNS server, both TCP and UDP, it works with brook server wsserver wssserver.](#brook-dnsserveroverbrook-can-create-a-encrypted-dns-server-both-tcp-and-udp-it-works-with-brook-server-wsserver-wssserver)
      * [brook tproxy Transparent Proxy Gateway on official OpenWrt](#brook-tproxy-transparent-proxy-gateway-on-official-openwrt)
      * [brook tproxy Transparent Proxy Gateway on any Linux (wired)](#brook-tproxy-transparent-proxy-gateway-on-any-linux-wired)
      * [GUI for official OpenWrt](#gui-for-official-openwrt)
      * [brook relay can relay a address to a remote address. It can relay any tcp and udp server](#brook-relay-can-relay-a-address-to-a-remote-address-it-can-relay-any-tcp-and-udp-server)
      * [brook socks5tohttp can convert a socks5 to a http proxy](#brook-socks5tohttp-can-convert-a-socks5-to-a-http-proxy)
      * [brook pac creates pac server](#brook-pac-creates-pac-server)
      * [brook pac creates pac file](#brook-pac-creates-pac-file)
      * [There are countless examples; for more feature suggestions, it's best to look at the commands and parameters in the CLI documentation one by one, and blog, YouTube...](#there-are-countless-examples-for-more-feature-suggestions-its-best-to-look-at-the-commands-and-parameters-in-the-cli-documentation-one-by-one-and-blog-youtube)
* [例子](#例子)
      * [运行 brook server](#运行-brook-server)
      * [运行 brook wsserver](#运行-brook-wsserver)
      * [运行 brook wssserver: 自动签发信任证书](#运行-brook-wssserver-自动签发信任证书)
      * [运行 brook wssserver 使用已有的信任机构签发的证书](#运行-brook-wssserver-使用已有的信任机构签发的证书)
      * [运行 brook wssserver 自己签发非信任证书, 甚至不是你自己的域名也可以](#运行-brook-wssserver-自己签发非信任证书-甚至不是你自己的域名也可以)
      * [withoutBrookProtocol](#withoutbrookprotocol-2)
      * [withoutBrookProtocol 自动签发信任证书](#withoutbrookprotocol-自动签发信任证书)
      * [withoutBrookProtocol 使用已有的信任机构签发的证书](#withoutbrookprotocol-使用已有的信任机构签发的证书)
      * [withoutBrookProtocol 自己签发非信任证书, 甚至不是你自己的域名也可以](#withoutbrookprotocol-自己签发非信任证书-甚至不是你自己的域名也可以)
      * [运行 brook socks5, 一个独立的标准 socks5 server](#运行-brook-socks5-一个独立的标准-socks5-server)
      * [运行 brook socks5, 一个独立的标准 socks5 server, 指定用户名和密码](#运行-brook-socks5-一个独立的标准-socks5-server-指定用户名和密码)
      * [brook relayoverbrook 中继任何 TCP 和 UDP server, 让其走 brook 协议. 它与 brook server wsserver wssserver 一起工作](#brook-relayoverbrook-中继任何-tcp-和-udp-server-让其走-brook-协议-它与-brook-server-wsserver-wssserver-一起工作)
      * [brook dnsserveroverbrook 用来创建一个加密 DNS Server, TCP and UDP, 它与 brook server wsserver wssserver 一起工作](#brook-dnsserveroverbrook-用来创建一个加密-dns-server-tcp-and-udp-它与-brook-server-wsserver-wssserver-一起工作)
      * [brook tproxy 透明代理网关在官网原版 OpenWrt](#brook-tproxy-透明代理网关在官网原版-openwrt)
      * [brook tproxy 透明代理网关在任意 Linux(有线)](#brook-tproxy-透明代理网关在任意-linux有线)
      * [官网原版 OpenWrt 图形客户端](#官网原版-openwrt-图形客户端)
      * [brook relay 可以中继任何 TCP 和 UDP server, 这是一个独立的功能, 它不依赖 brook server wsserver wssserver](#brook-relay-可以中继任何-tcp-和-udp-server-这是一个独立的功能-它不依赖-brook-server-wsserver-wssserver)
      * [brook socks5tohttp 将 socks5 proxy 转换为 http proxy](#brook-socks5tohttp-将-socks5-proxy-转换为-http-proxy)
      * [brook pac 创建一个 pac server](#brook-pac-创建一个-pac-server)
      * [brook pac 创建一个 pac 文件](#brook-pac-创建一个-pac-文件)
      * [例子不胜枚举，更多功能建议挨个看 CLI 文档的命令和参数吧，还有博客，YouTube 等...](#例子不胜枚举更多功能建议挨个看-cli-文档的命令和参数吧还有博客youtube-等)

# Brook
<!--SIDEBAR-->
<!--G-R3M673HK5V-->
A cross-platform programmable network tool.

# Sponsor
**❤️  [Shiliew - China Optimized Network App](https://www.txthinking.com/shiliew.html)**
# Getting Started

## Server

```
bash <(curl https://bash.ooo/nami.sh)
```

```
nami install brook
```

```
brook server -l :9999 -p hello
```

## GUI Client

| iOS | Android      | Mac    |Windows      |Linux        |OpenWrt      |
| --- | --- | --- | --- | --- | --- |
| [![](https://brook.app/images/appstore.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![](https://brook.app/images/android.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.apk) | [![](https://brook.app/images/mac.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![Windows](https://brook.app/images/windows.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.exe) | [![](https://brook.app/images/linux.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.bin) | [![OpenWrt](https://brook.app/images/openwrt.png)](https://github.com/txthinking/brook/releases) |

> OpenWrt: After installation, you need to refresh the page to see the menu

## CLI Client

```
brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```
# GUI Documentation

## Software for which this article applies

-   [Brook](https://github.com/txthinking/brook)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Programmable

```
Brook GUI will pass different global variables to the script at different times, and the script only needs to assign the processing result to the global variable out
```

### Introduction to incoming variables

| variable                       | type | condition   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_brooklinks                  | map  | / | Before connecting  | Predefine multiple brook links, and then programmatically specify which one to connect to | map      |
| in_dnsquery                    | map  | FakeDNS: On | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | /           | When connecting to an address     | script can decide how to connect                  | map      |
| in_httprequest                 | map  | /           | When an HTTP(S) request comes in  | the script can decide how to handle this request  | map      |
| in_httprequest,in_httpresponse | map  | /           | when an HTTP(S) response comes in | the script can decide how to handle this response | map      |

### in_brooklinks

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

### in_dnsquery

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| domain | string | domain name | google.com |
| type   | string | query type  | A          |
| appid   | string | App ID or path | com.google.Chrome.helper          |
| interface   | string | network interface. Mac only | en0          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`                                                | false   |
| ip           | string | Specify IP directly, only valid when `type` is `A`/`AAAA`                                                                     | 1.2.3.4 |
| system       | bool   | Resolve by System DNS, default `false`                                                                                       | false   |
| bypass       | bool   | Resolve by Bypass DNS, default `false` | false   |
| brooklinkkey | string   | When need to connect the Server，instead, connect to the Server specified by the key in_brooklinks | custom name   |

### in_address

| Key           | Type   | Description                                                                                                         | Example        |
| ------------- | ------ | ------------------------------------------------------------------------------------------------------------------- | -------------- |
| network       | string | Network type, the value `tcp`/`udp`                                                                                 | tcp            |
| ipaddress     | string | IP type address. There is only of ipaddress and domainaddress. Note that there is no relationship between these two | 1.2.3.4:443    |
| domainaddress | string | Domain type address, because of FakeDNS we can get the domain name address here                                     | google.com:443 |
| appid   | string | App ID or path | com.google.Chrome.helper          |
| interface   | string | network interface. Mac only | en0          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key                    | Type   | Description                                                                                                                                                                                             | Example     |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| block                  | bool   | Whether Block, default `false`                                                                                                                                                                          | false       |
| ipaddress              | string | IP type address, rewrite destination                                                                                                                                                                    | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | Use Bypass DNS to obtain `A` or `AAAA` IP and rewrite the destination, only valid when `domainaddress` exists, the value `A`/`AAAA`                                                                     | A           |
| bypass                 | bool   | Bypass, default `false`. If `true` and `domainaddress`, then `ipaddress` or `ipaddressfrombypassdns` must be specified | false       |
| mitm                   | bool   | Whether to perform MITM, default `false`. Only valid when `network` is `tcp`. Need to install CA, see below                                                                                             | false       |
| mitmprotocol           | string | MITM protocol needs to be specified explicitly, the value is `http`/`https`                                                                                                                             | https       |
| mitmcertdomain         | string | The MITM certificate domain name, which is taken from `domainaddress` by default. If `ipaddress` and `mitm` is `true` and `mitmprotocol` is `https` then must be must be specified explicitly           | example.com |
| mitmwithbody           | bool   | Whether to manipulate the http body, default `false`. will read the body of the request and response into the memory and interact with the script. iOS 50M total memory limit may kill process      | false       |
| mitmautohandlecompress | bool   | Whether to automatically decompress the http body when interacting with the script, default `false`                                                                                                     | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                                                      | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                                                    | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                                                     | 0           |
| brooklinkkey | string   | When need to connect the Server，instead, connect to the Server specified by the key in_brooklinks | custom name   |

### in_httprequest

| Key    | Type   | Description                   | Example                     |
| ------ | ------ | ----------------------------- | --------------------------- |
| URL    | string | URL                           | `https://example.com/hello` |
| Method | string | HTTP method                   | GET                         |
| Body   | bytes  | HTTP request body             | /                           |
| ...    | string | other fields are HTTP headers | /                           |

`out`, must be set to a request or response

### in_httpresponse

| Key        | Type   | Description                   | Example |
| ---------- | ------ | ----------------------------- | ------- |
| StatusCode | int    | HTTP status code              | 200     |
| Body       | bytes  | HTTP response body            | /       |
| ...        | string | other fields are HTTP headers | /       |

`out`, must be set to a response

## Write script

[Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

Library

-   [text](https://github.com/d5/tengo/blob/master/docs/stdlib-text.md): regular expressions, string conversion, and manipulation
-   [math](https://github.com/d5/tengo/blob/master/docs/stdlib-math.md): mathematical constants and functions
-   [times](https://github.com/d5/tengo/blob/master/docs/stdlib-times.md): time-related functions
-   [rand](https://github.com/d5/tengo/blob/master/docs/stdlib-rand.md): random functions
-   [fmt](https://github.com/d5/tengo/blob/master/docs/stdlib-fmt.md): formatting functions
-   [json](https://github.com/d5/tengo/blob/master/docs/stdlib-json.md): JSON functions
-   [enum](https://github.com/d5/tengo/blob/master/docs/stdlib-enum.md): Enumeration functions
-   [hex](https://github.com/d5/tengo/blob/master/docs/stdlib-hex.md): hex encoding and decoding functions
-   [base64](https://github.com/d5/tengo/blob/master/docs/stdlib-base64.md): base64 encoding and decoding functions
-   `brook`: brook module

    ```
    Constants

    * os: string, linux/darwin/windows/ios/android

    Functions

    * splithostport(address string) => map/error: splits a network address of the form "host:port" to { "host": "xxx", "port": "xxx" }
    * country(ip string) => string/error: get country code from ip
    * cidrcontainsip(cidr string, ip string) => bool/error: reports whether the network includes ip
    * parseurl(url string) => map/error: parses a raw url into a map, keys: scheme/host/path/rawpath/rawquery
    * parsequery(query string) => map/error: parses a raw query into a kv map
    * map2query(kv map) => string/error: convert map{string:string} into a query string
    * bytes2ints(b bytes) => array/error: convert bytes into [int]
    * ints2bytes(ints array) => bytes/error: convert [int] into bytes
    * bytescompare(a bytes, b bytes) => int/error: returns an integer comparing two bytes lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b
    * bytescontains(b bytes, sub bytes) => bool/error: reports whether sub is within b
    * byteshasprefix(s bytes, prefix bytes) => bool/error: tests whether the bytes s begins with prefix
    * byteshassuffix(s bytes, suffix bytes) => bool/error: tests whether the bytes s ends with suffix
    * bytesindex(s bytes, sep bytes) => int/error: returns the index of the first instance of sep in s, or -1 if sep is not present in s
    * byteslastindex(s bytes, sep bytes) => int/error: returns the index of the last instance of sep in s, or -1 if sep is not present in s
    * bytesreplace(s bytes, old bytes, new bytes, n int) => bytes/error: returns a copy of the s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements
    * pathescape(s string) => string/error: escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed
    * pathunescape(s string) => string/error: does the inverse transformation of pathescape
    * queryescape(s string) => string/error: escapes the string so it can be safely placed inside a URL query
    * queryunescape(s string) => string/error: does the inverse transformation of queryescape
    * hexdecode(s string) => bytes/error: returns the bytes represented by the hexadecimal string s
    * hexencode(s string) => string/error: returns the hexadecimal encoding of src
    ```

## Debug script

It is recommended to use [tun2brook](https://github.com/txthinking/tun2brook) on desktop to debug with `fmt.println`

## Install CA

https://txthinking.github.io/ca/ca.pem

| OS | How |
| --- | --- |
| iOS | https://www.youtube.com/watch?v=HSGPC2vpDGk |
| Android | Android has user CA and system CA, must be installed in the system CA after ROOT |
| macOS | `nami install mad ca.txthinking`, `sudo mad install --ca ~/.nami/bin/ca.pem` |
| Windows | `nami install mad ca.txthinking`, Admin: `mad install --ca ~/.nami/bin/ca.pem` |

> Some software may not read the system CA，you can use `curl --cacert ~/.nami/bin/ca.pem` to debug

# 图形客户端文档

## 本文适用的软件

-   [Brook](https://github.com/txthinking/brook)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## 编程

```
Brook GUI 会在不同时机向脚本传入不同的全局变量，脚本只需要将处理结果赋值到全局变量 out 即可
```

### 传入变量介绍

| 变量                           | 类型 | 条件          | 时机                   | 描述                       | out 类型 |
| ------------------------------ | ---- | ------------- | ---------------------- | -------------------------- | -------- |
| in_brooklinks                  | map  | / | 连接之前  | 预定义多个 brook link，之后可编程指定连接哪个       | map      |
| in_dnsquery                    | map  | FakeDNS: 开启 | 当 DNS 查询发生时      | 脚本可以决定如何处理此请求 | map      |
| in_address                     | map  | /             | 当要连接某地址时       | 脚本可以决定如何进行连接   | map      |
| in_httprequest                 | map  | /             | 当有 HTTP(S)请求传入时 | 脚本可以决定如何处理此请求 | map      |
| in_httprequest,in_httpresponse | map  | /             | 当有 HTTP(S)响应传入时 | 脚本可以决定如何处理此响应 | map      |

### in_brooklinks

| Key    | 类型   | 描述     | 示例       |
| ------ | ------ | -------- | ---------- |
| _ | bool | 占位，无实际意义    | true |

`out`, 如果不是 `map` 类型则会被忽略

| Key          | 类型   | 描述                                                                                               | 示例    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| 自定义名字    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

### in_dnsquery

| Key    | 类型   | 描述     | 示例       |
| ------ | ------ | -------- | ---------- |
| domain | string | 域名     | google.com |
| type   | string | 查询类型 | A          |
| appid   | string | App ID 或路径 | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key          | 类型   | 描述                                                                                               | 示例    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | 是否 Block, 默认 `false`                                           | false   |
| ip           | string | 直接指定 IP，仅当 `type` 为 `A`/`AAAA`有效                                                         | 1.2.3.4 |
| system       | bool   | 使用 System DNS 来解析，默认 `false`                                                               | false   |
| bypass       | bool   | 使用 Bypass DNS 来解析，默认 `false` | false   |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的代理服务器 | 自定义名字   |

### in_address

| Key           | 类型   | 描述                                                                     | 示例           |
| ------------- | ------ | ------------------------------------------------------------------------ | -------------- |
| network       | string | 即将发起连接网络，取值 `tcp`/`udp`                                       | tcp            |
| ipaddress     | string | IP 类型的地址，与 domainaddress 只会存在一个。注意这两个之间没有任何关系 | 1.2.3.4:443    |
| domainaddress | string | 域名类型的地址，因为 FakeDNS 我们这里才能拿到域名地址                    | google.com:443 |
| appid   | string | App ID 或路径 | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key                    | 类型   | 描述                                                                                                                                                                           | 示例        |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- |
| block                  | bool   | 是否 Block, 默认 `false`                                                                                                                                                       | false       |
| ipaddress              | string | IP 类型地址，重写目的地                                                                                                                                                        | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | 使用 Bypass DNS 获取`A`或`AAAA` IP 并重写目的地, 仅当 `domainaddress` 存在时有效，取值 `A`/`AAAA`                                                                              | A           |
| bypass                 | bool   | 是否 Bypass, 默认 `false`. 如果为 `true` 并且是 `domainaddress`, 那么必须指定 `ipaddress` 或 `ipaddressfrombypassdns` | false       |
| mitm                   | bool   | 是否进行 MITM, 默认 `false`. 仅当 `network` 为 `tcp` 时有效. 需要安装 CA，看下文介绍                                                                                           | false       |
| mitmprotocol           | string | 需要明确指定 MITM 协议, 取值 `http`/`https`                                                                                                                                    | https       |
| mitmcertdomain         | string | MITM 证书域名，默认从`domainaddress`里取。如果是 `ipaddress` 且 `mitm` 为 `true` 且 `mitmprotocol` 为 `https` 那么必须明确指定                                                 | example.com |
| mitmwithbody           | bool   | 是否操作 http body，默认 `false`. 会将请求和响应的 body 读取到内存里和脚本交互。iOS 50M 总内存限制可能会杀进程                                                             | false       |
| mitmautohandlecompress | bool   | 和脚本交互时是否自动解压缩 http body, 默认 `false`                                                                                                                             | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                             | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                           | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                            | 0           |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的代理服务器 | 自定义名字   |

### in_httprequest

| Key    | 类型   | 描述                     | 示例                        |
| ------ | ------ | ------------------------ | --------------------------- |
| URL    | string | URL                      | `https://example.com/hello` |
| Method | string | HTTP method              | GET                         |
| Body   | bytes  | HTTP request body        | /                           |
| ...    | string | 其他字段均为 HTTP header | /                           |

`out`, 必须设置为一个 request 或 response

### in_httpresponse

| Key        | 类型   | 描述                     | 示例 |
| ---------- | ------ | ------------------------ | ---- |
| StatusCode | int    | HTTP status code         | 200  |
| Body       | bytes  | HTTP response body       | /    |
| ...        | string | 其他字段均为 HTTP header | /    |

`out`, 必须设置为一个 response

## 写脚本

[Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

Library

-   [text](https://github.com/d5/tengo/blob/master/docs/stdlib-text.md): regular expressions, string conversion, and manipulation
-   [math](https://github.com/d5/tengo/blob/master/docs/stdlib-math.md): mathematical constants and functions
-   [times](https://github.com/d5/tengo/blob/master/docs/stdlib-times.md): time-related functions
-   [rand](https://github.com/d5/tengo/blob/master/docs/stdlib-rand.md): random functions
-   [fmt](https://github.com/d5/tengo/blob/master/docs/stdlib-fmt.md): formatting functions
-   [json](https://github.com/d5/tengo/blob/master/docs/stdlib-json.md): JSON functions
-   [enum](https://github.com/d5/tengo/blob/master/docs/stdlib-enum.md): Enumeration functions
-   [hex](https://github.com/d5/tengo/blob/master/docs/stdlib-hex.md): hex encoding and decoding functions
-   [base64](https://github.com/d5/tengo/blob/master/docs/stdlib-base64.md): base64 encoding and decoding functions
-   `brook`: brook module

    ```
    Constants

    * os: string, linux/darwin/windows/ios/android

    Functions

    * splithostport(address string) => map/error: splits a network address of the form "host:port" to { "host": "xxx", "port": "xxx" }
    * country(ip string) => string/error: get country code from ip
    * cidrcontainsip(cidr string, ip string) => bool/error: reports whether the network includes ip
    * parseurl(url string) => map/error: parses a raw url into a map, keys: scheme/host/path/rawpath/rawquery
    * parsequery(query string) => map/error: parses a raw query into a kv map
    * map2query(kv map) => string/error: convert map{string:string} into a query string
    * bytes2ints(b bytes) => array/error: convert bytes into [int]
    * ints2bytes(ints array) => bytes/error: convert [int] into bytes
    * bytescompare(a bytes, b bytes) => int/error: returns an integer comparing two bytes lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b
    * bytescontains(b bytes, sub bytes) => bool/error: reports whether sub is within b
    * byteshasprefix(s bytes, prefix bytes) => bool/error: tests whether the bytes s begins with prefix
    * byteshassuffix(s bytes, suffix bytes) => bool/error: tests whether the bytes s ends with suffix
    * bytesindex(s bytes, sep bytes) => int/error: returns the index of the first instance of sep in s, or -1 if sep is not present in s
    * byteslastindex(s bytes, sep bytes) => int/error: returns the index of the last instance of sep in s, or -1 if sep is not present in s
    * bytesreplace(s bytes, old bytes, new bytes, n int) => bytes/error: returns a copy of the s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements
    * pathescape(s string) => string/error: escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed
    * pathunescape(s string) => string/error: does the inverse transformation of pathescape
    * queryescape(s string) => string/error: escapes the string so it can be safely placed inside a URL query
    * queryunescape(s string) => string/error: does the inverse transformation of queryescape
    * hexdecode(s string) => bytes/error: returns the bytes represented by the hexadecimal string s
    * hexencode(s string) => string/error: returns the hexadecimal encoding of src
    ```

## 调试脚本

建议使用 [tun2brook](https://github.com/txthinking/tun2brook) 在电脑上`fmt.println`调试

## 安装 CA

https://txthinking.github.io/ca/ca.pem

| OS | 如何 |
| --- | --- |
| iOS | https://www.youtube.com/watch?v=HSGPC2vpDGk |
| Android | Android 分系统 CA 和用户 CA，必须要 ROOT 后安装到系统 CA 里 |
| macOS | `nami install mad ca.txthinking`, `sudo mad install --ca ~/.nami/bin/ca.pem` |
| Windows | `nami install mad ca.txthinking`, 管理员: `mad install --ca ~/.nami/bin/ca.pem` |

> 注意有些软件可能不读取系统 CA，可以使用 `curl --cacert ~/.nami/bin/ca.pem` 调试

# Resources

| CLI | Description |
| --- | --- |
| [nami](https://github.com/txthinking/nami) | A clean and tidy decentralized package manager |
| [joker](https://github.com/txthinking/joker) | Joker can turn process into daemon. Zero-Configuration |
| [nico](https://github.com/txthinking/nico) | Nico can work with brook wsserver together |
| [zhen](https://github.com/txthinking/zhen) | zhen - process and cron manager |
| [tun2brook](https://github.com/txthinking/tun2brook) | Proxy all traffic just one line command |
| [mad](https://github.com/txthinking/mad) | Generate root CA and derivative certificate for any domains and any IPs |
| [hancock](https://github.com/txthinking/hancock) | Manage multiple remote servers and execute commands remotely |
| [sshexec](https://github.com/txthinking/sshexec) | A command-line tool to execute remote command through ssh |
| [jb](https://github.com/txthinking/jb) | write script in an easier way than bash |
| [bash](https://github.com/txthinking/bash) | One-click script 一键脚本 |
| [pacman](https://archlinux.org/packages/extra/x86_64/brook/) | `pacman -S brook` |
| [brew](https://formulae.brew.sh/formula/brook) | `brew install brook` |
| [docker](https://hub.docker.com/r/txthinking/brook) | `docker run txthinking/brook` | 

| Example | 举例 |
| --- | --- |
| [Example](https://github.com/txthinking/brook/blob/master/docs/example.md) | [例子](https://github.com/txthinking/brook/blob/master/docs/example-zh.md) |

| Resources | Description |
| --- | --- |
| [Protocol](https://github.com/txthinking/brook/tree/master/protocol) | Brook Protocol |
| [Blog](https://www.txthinking.com/talks/) | Some articles you should read |
| [YouTube](https://www.youtube.com/txthinking) | Some videos you should watch |
| [Telegram](https://t.me/txthinking) | Ask questions here |
| [Announce](https://t.me/s/txthinking_news) | All news you should care |
| [GitHub](https://github.com/txthinking) | Other useful repos |
| [Socks5 Configurator](https://chromewebstore.google.com/detail/socks5-configurator/hnpgnjkeaobghpjjhaiemlahikgmnghb) | If you prefer CLI brook client | 
| [Pastebin](https://paste.brook.app) | Create importable many brook links |
| [Brook Deploy](https://www.txthinking.com/deploy.html) | Deploy brook with GUI |
| [TxThinking](https://www.txthinking.com) | Everything |

# CLI Documentation
# NAME

Brook - A cross-platform programmable network tool

# SYNOPSIS

Brook

```
brook [全局参数] 子命令 [子命令参数]
```

**Usage**:

```
Brook [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

- **--clientHKDFInfo**="": client HKDF info, most time you don't need to change this, if changed, all and each brook links in client side must be same, I mean each (default: "brook")

- **--dialWithDNS**="": When a domain name needs to be resolved, use the specified DNS. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required. Note that for client-side commands, this does not affect the client passing the domain address to the server

- **--dialWithDNSPrefer**="": This is used with the dialWithDNS parameter. Prefer A record or AAAA record. Value is A or AAAA

- **--dialWithIP4**="": When the current machine establishes a network connection to the outside IPv4, both TCP and UDP, it is used to specify the IPv4 used

- **--dialWithIP6**="": When the current machine establishes a network connection to the outside IPv6, both TCP and UDP, it is used to specify the IPv6 used

- **--dialWithNIC**="": When the current machine establishes a network connection to the outside, both TCP and UDP, it is used to specify the NIC used

- **--dialWithSocks5**="": When the current machine establishes a network connection to the outside, both TCP and UDP, with your socks5 proxy, such as 127.0.0.1:1081

- **--dialWithSocks5Password**="": If there is

- **--dialWithSocks5TCPTimeout**="": time (s) (default: 0)

- **--dialWithSocks5UDPTimeout**="": time (s) (default: 60)

- **--dialWithSocks5Username**="": If there is

- **--help, -h**: show help

- **--log**="": Enable log. A valid value is file path or 'console'. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true

- **--pprof**="": go http pprof listen addr, such as :6060

- **--prometheus**="": prometheus http listen addr, such as :7070. If it is transmitted on the public network, it is recommended to use it with nico

- **--prometheusPath**="": prometheus http path, such as /xxx. If it is transmitted on the public network, a hard-to-guess value is recommended

- **--serverHKDFInfo**="": server HKDF info, most time you don't need to change this, if changed, all and each brook links in client side must be same, I mean each (default: "brook")

- **--tag**="": Tag can be used to the process, will be append into log, such as: 'key1:value1'

- **--version, -v**: print the version


# COMMANDS

## server

Run as brook server, both TCP and UDP

- **--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

- **--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--blockGeoIP**="": Block IP by Geo country code, such as US

- **--listen, -l**="": Listen address, like: ':9999'

- **--password, -p**="": Server password

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

## client

Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]

- **--http**="": Where to listen for HTTP proxy connections

- **--password, -p**="": Brook server password

- **--server, -s**="": Brook server address, like: 1.2.3.4:9999

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--udpovertcp**: UDP over TCP

## wsserver

Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server

- **--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

- **--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--blockGeoIP**="": Block IP by Geo country code, such as US

- **--listen, -l**="": Listen address, like: ':80'

- **--password, -p**="": Server password

- **--path**="": URL path (default: /ws)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

- **--xForwardedFor**: Replace the from field in --log, note that this may be forged

## wsclient

Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]

- **--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--http**="": Where to listen for HTTP proxy connections

- **--password, -p**="": Brook wsserver password

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

- **--wsserver, -s**="": Brook wsserver address, like: ws://1.2.3.4:80, if no path then /ws will be used. Do not omit the port under any circumstances

## wssserver

Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server

- **--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

- **--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--blockGeoIP**="": Block IP by Geo country code, such as US

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--password, -p**="": Server password

- **--path**="": URL path (default: /ws)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wssclient

Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]

- **--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

- **--http**="": Where to listen for HTTP proxy connections

- **--insecure**: Client do not verify the server's certificate chain and host name

- **--password, -p**="": Brook wssserver password

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--udpTimeout**="": time (s) (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

- **--wssserver, -s**="": Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances

## quicserver

Run as brook quicserver, both TCP and UDP

- **--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

- **--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--blockGeoIP**="": Block IP by Geo country code, such as US

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--password, -p**="": Server password

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## quicclient

Run as brook quicclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]. (Note that the global dial parameter is ignored now)

- **--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--ca**="": Specify ca instead of insecure, such as /path/to/ca.pem

- **--http**="": Where to listen for HTTP proxy connections

- **--insecure**: Client do not verify the server's certificate chain and host name

- **--password, -p**="": Brook quicserver password

- **--quicserver, -s**="": Brook quicserver address, like: quic://google.com:443. Do not omit the port under any circumstances

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## relayoverbrook

Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]

- **--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

- **--from, -f, -l**="": Listen address: like ':9999'

- **--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

- **--password, -p**="": Password

- **--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443

- **--tcpTimeout**="": time (s) (default: 0)

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--to, -t**="": Address which relay to, like: 1.2.3.4:9999

- **--udpTimeout**="": time (s) (default: 0)

- **--udpovertcp**: When server is brook server, UDP over TCP

- **--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## dnsserveroverbrook

Run as dns server over brook, both TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]

- **--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--bypassDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

- **--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

- **--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

- **--listen, -l**="": Listen address, like: 127.0.0.1:53

- **--password, -p**="": Password

- **--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

- **--tcpTimeout**="": time (s) (default: 0)

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--udpTimeout**="": time (s) (default: 0)

- **--udpovertcp**: When server is brook server, UDP over TCP

- **--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## tproxy

Run as transparent proxy, a router gateway, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

- **--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--blockDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--bypassCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

- **--bypassCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

- **--bypassDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--bypassGeoIP**="": Bypass IP by Geo country code, such as US

- **--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

- **--dnsForDefault**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

- **--dnsListen**="": Start a DNS server, like: ':53'. MUST contain IP, like '192.168.1.1:53', if you expect your gateway to accept requests from clients to other public DNS servers at the same time

- **--doNotRunScripts**: This will not change iptables and others if you want to do by yourself

- **--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

- **--link**="": brook link. This will ignore server, password, udpovertcp, address, insecure, withoutBrookProtocol, ca, tlsfingerprint

- **--listen, -l**="": Listen address, DO NOT contain IP, just like: ':8888'. No need to operate iptables by default! (default: :8888)

- **--password, -p**="": Password

- **--redirectDNS**="": It is usually the value of dnsListen. If the client has set custom DNS instead of dnsListen, this parameter can be intercepted and forwarded to dnsListen. Usually you don't need to set this, only if you want to control it instead of being proxied directly as normal UDP data.

- **--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

- **--tcpTimeout**="": time (s) (default: 0)

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--udpTimeout**="": time (s) (default: 0)

- **--udpovertcp**: When server is brook server, UDP over TCP

- **--webListen**="": Ignore all other parameters, run web UI, like: ':9999'

- **--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## link

Generate brook link

- **--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--ca**="": When server is brook wssserver or brook quicserver, specify ca for untrusted cert, such as /path/to/ca.pem

- **--clientHKDFInfo**="": client HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is

- **--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

- **--name**="": Give this server a name

- **--password, -p**="": Password

- **--server, -s**="": Support brook server, brook wsserver, brook wssserver, socks5 server, brook quicserver. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080, quic://google.com:443

- **--serverHKDFInfo**="": server HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--udpovertcp**: When server is brook server, UDP over TCP

- **--username, -u**="": Username, when server is socks5 server

- **--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## connect

Run as client and connect to brook link, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

- **--http**="": Where to listen for HTTP proxy connections

- **--link, -l**="": brook link, you can get it via $ brook link

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## relay

Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]

- **--from, -f, -l**="": Listen address: like ':9999'

- **--tcpTimeout**="": time (s) (default: 0)

- **--to, -t**="": Address which relay to, like: 1.2.3.4:9999

- **--udpTimeout**="": time (s) (default: 0)

## dnsserver

Run as standalone dns server

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

- **--listen, -l**="": Listen address, like: 127.0.0.1:53

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## dnsclient

Send a dns query

- **--dns, -s**="": DNS server, such as 8.8.8.8:53 (default: 8.8.8.8:53)

- **--domain, -d**="": Domain

- **--short**: Short for A/AAAA

- **--type, -t**="": Type, such as A (default: A)

## dohserver

Run as standalone doh server

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

- **--domainaddress**="": Such as: domain.com:443, if you want to create a https server. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--listen**="": listen address, if you want to create a http server behind nico

- **--path**="": URL path (default: /dns-query)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## dohclient

Send a dns query

- **--doh, -s**="": DOH server, the address is required (default: https://dns.quad9.net/dns-query?address=9.9.9.9%3A443)

- **--domain, -d**="": Domain

- **--short**: Short for A/AAAA

- **--type, -t**="": Type, such as A (default: A)

## dhcpserver

Run as standalone dhcp server. Note that you need to stop other dhcp servers, if there are.

- **--cache**="": Cache file, local absolute file path, default is $HOME/.brook.dhcpserver

- **--count**="": IP range from the start, which you want to assign to clients (default: 0)

- **--dnsserver**="": The dns server which you want to assign to clients, such as: 192.168.1.1 or 8.8.8.8

- **--gateway**="": The router gateway which you want to assign to clients, such as: 192.168.1.1

- **--interface**="": Select interface on multi interface device. Linux only

- **--netmask**="": Subnet netmask which you want to assign to clients (default: 255.255.255.0)

- **--serverip**="": DHCP server IP, the IP of the this machine, you shoud set a static IP to this machine before doing this, such as: 192.168.1.10

- **--start**="": Start IP which you want to assign to clients, such as: 192.168.1.100

## socks5

Run as standalone standard socks5 server, both TCP and UDP

- **--limitUDP**: The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.

- **--listen, -l**="": Socks5 server listen address, like: :1080 or 1.2.3.4:1080

- **--password**="": Password, optional

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": Connection deadline time (s) (default: 0)

- **--udpTimeout**="": Connection deadline time (s) (default: 0)

- **--username**="": User name, optional

## socks5tohttp

Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]

- **--listen, -l**="": HTTP proxy which will be create: like: 127.0.0.1:8010

- **--socks5, -s**="": Socks5 server address, like: 127.0.0.1:1080

- **--socks5password**="": Socks5 password, optional

- **--socks5username**="": Socks5 username, optional

- **--tcpTimeout**="": Connection tcp timeout (s) (default: 0)

## pac

Run as PAC server or save PAC to file

- **--bypassDomainList, -b**="": One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

- **--file, -f**="": Save PAC to file, this will ignore listen address

- **--listen, -l**="": Listen address, like: 127.0.0.1:1980

- **--proxy, -p**="": Proxy, like: 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' (default: SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT)

## testsocks5

Test UDP and TCP of socks5 server

- **--dns**="": DNS server for connecting (default: 8.8.8.8:53)

- **--domain**="": Domain for query (default: http3.ooo)

- **--password, -p**="": Socks5 password

- **--socks5, -s**="": Like: 127.0.0.1:1080

- **--username, -u**="": Socks5 username

- **-a**="": The A record of domain (default: 137.184.237.95)

## testbrook

Test UDP and TCP of brook server/wsserver/wssserver/quicserver. (Note that the global dial parameter is ignored now)

- **--dns**="": DNS server for connecting (default: 8.8.8.8:53)

- **--domain**="": Domain for query (default: http3.ooo)

- **--link, -l**="": brook link. Get it via $ brook link

- **--socks5**="": Temporarily listening socks5 (default: 127.0.0.1:11080)

- **-a**="": The A record of domain (default: 137.184.237.95)

## echoserver

Echo server, echo UDP and TCP address of routes

- **--listen, -l**="": Listen address, like: ':7777'

## echoclient

Connect to echoserver, echo UDP and TCP address of routes

- **--server, -s**="": Echo server address, such as 1.2.3.4:7777

- **--times**="": Times of interactions (default: 0)

## ipcountry

Get country of IP

- **--ip**="": 1.1.1.1

## completion

Generate shell completions

- **--file, -f**="": Write to file (default: brook_autocomplete)

## mdpage

Generate markdown page

- **--file, -f**="": Write to file, default print to stdout

- **--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## manpage

Generate man.1 page

- **--file, -f**="": Write to file, default print to stdout. You should put to /path/to/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos

## help, h

Shows a list of commands or help for one command
# Diagram

> Maybe outdated

## overview

![overview](https://txthinking.github.io/brook/svg/overview.svg)

## withoutBrookProtocol

![wbp](https://txthinking.github.io/brook/svg/wbp.svg)

## relayoverbrook

![relayoverbrook](https://txthinking.github.io/brook/svg/relayoverbrook.svg)

## dnsserveroverbrook

![dnsserveroverbrook](https://txthinking.github.io/brook/svg/dnsserveroverbrook.svg)

## relay

![relay](https://txthinking.github.io/brook/svg/relay.svg)

## dnsserver

![dnsserver](https://txthinking.github.io/brook/svg/dnsserver.svg)

## tproxy

![tproxy](https://txthinking.github.io/brook/svg/tproxy.svg)

## gui

![gui](https://txthinking.github.io/brook/svg/gui.svg)

## script

![script](https://txthinking.github.io/brook/svg/script.svg)

# Examples

List some examples of common scene commands, pay attention to replace the parameters such as IP, port, password, domain name, certificate path, etc. in the example by yourself

### Run brook server

```
brook server --listen :9999 --password hello
```

then

-   server: `1.2.3.4:9999`
-   password: `hello`

or get brook link

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

or get brook link with `--udpovertcp`

```
brook link --server 1.2.3.4:9999 --password hello --udpovertcp --name 'my brook server'
```

### Run brook wsserver

```
brook wsserver --listen :9999 --password hello
```

then

-   server: `ws://1.2.3.4:9999`
-   password: `hello`

or get brook link

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

then

-   server: `wss://domain.com:443`
-   password: `hello`

or get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### Run brook wssserver Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

then

-   server: `wss://domain.com:443`
-   password: `hello`

or get brook link

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

Run brook

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

get brook link with `--insecure`

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

get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

### withoutBrookProtocol Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

get brook link

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

### Run brook socks5, A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

then

-   server: `1.2.3.4:1080`

or get brook link

```
brook link --server socks5://1.2.3.4:1080
```

### Run brook socks5 with username and password. A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

then

-   server: `1.2.3.4:1080`
-   username: `hello`
-   password: `world`

or get brook link

```
brook link --server socks5://1.2.3.4:1080 --username hello --password world
```

### brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

### brook dnsserveroverbrook can create a encrypted DNS server, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook dnsserveroverbrook ... --listen 127.0.0.1:53
```

### brook tproxy Transparent Proxy Gateway on official OpenWrt

**No need to manipulate iptables!**

```
opkg update
opkg install ca-certificates openssl-util ca-bundle coreutils-nohup iptables iptables-mod-tproxy iptables-mod-socket ip6tables
```

```
brook tproxy --link 'brook://...' --dnsListen :5353
```

1. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
2. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
3. By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

### brook tproxy Transparent Proxy Gateway on any Linux (wired)

**No need to manipulate iptables!**

```
systemctl stop systemd-resolved
systemctl disable systemd-resolved
echo nameserver 8.8.8.8 > /etc/resolv.conf
```

```
brook tproxy --link 'brook://...' --dnsListen 192.168.1.2:53 --disableAAAA
```

Replace 192.168.1.2 with your Linux IP. You may need to manually configure the computer or mobile gateway and DNS.

### GUI for official OpenWrt

**No need to manipulate iptables!**

port 9999, 8888, 5353 will be used. It work with brook server, brook wsserver, brook wssserver and brook quicserver.

1. Download the [ipk](https://github.com/txthinking/brook/releases) file for your router
2. Upload and install: OpenWrt Web -> System -> Software -> Upload Package...
3. Refresh page, the Brook menu will appear at the top
4. OpenWrt Web -> Brook -> type and Connect
5. And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
6. And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

### brook relay can relay a address to a remote address. It can relay any tcp and udp server

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

### There are countless examples; for more feature suggestions, it's best to look at the commands and parameters in the CLI documentation one by one, and blog, YouTube...
# 例子

下面列举一些常用场景命令的例子, 注意自己替换示例中的 IP，端口，密码，域名，证书路径等参数

### 运行 brook server

```
brook server --listen :9999 --password hello
```

然后

-   server: `1.2.3.4:9999`
-   password: `hello`

或 获取 brook link

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

或 获取 brook link 让 udp 走 tcp `--udpovertcp`

```
brook link --server 1.2.3.4:9999 --password hello --udpovertcp --name 'my brook server'
```

### 运行 brook wsserver

```
brook wsserver --listen :9999 --password hello
```

然后

-   server: `ws://1.2.3.4:9999`
-   password: `hello`

或 获取 brook link

```
brook link --server ws://1.2.3.4:9999 --password hello --name 'my brook wsserver'
```

或 获取 brook link 指定个域名, 甚至不是你自己的域名也可以

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

### 运行 brook wssserver: 自动签发信任证书

> 注意：确保你的域名已成功解析到你服务器的 IP, 自动签发证书需要额外监听 80 端口

```
brook wssserver --domainaddress domain.com:443 --password hello
```

然后

-   server: `wss://domain.com:443`
-   password: `hello`

或 获取 brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

### 运行 brook wssserver 使用已有的信任机构签发的证书

> 注意：确保你的域名已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

然后

-   server: `wss://domain.com:443`
-   password: `hello`

或 获取 brook link

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

运行 brook

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
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

### 运行 brook socks5, 一个独立的标准 socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

然后

-   server: `1.2.3.4:1080`

或 获取 brook link

```
brook link --server socks5://1.2.3.4:1080
```

### 运行 brook socks5, 一个独立的标准 socks5 server, 指定用户名和密码

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

然后

-   server: `1.2.3.4:1080`
-   username: `hello`
-   password: `world`

或 获取 brook link

```
brook link --server socks5://1.2.3.4:1080 --username hello --password world
```

### brook relayoverbrook 中继任何 TCP 和 UDP server, 让其走 brook 协议. 它与 brook server wsserver wssserver 一起工作

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

### brook dnsserveroverbrook 用来创建一个加密 DNS Server, TCP and UDP, 它与 brook server wsserver wssserver 一起工作

```
brook dnsserveroverbrook ... --listen 127.0.0.1:53
```

### brook tproxy 透明代理网关在官网原版 OpenWrt

**无需操作 iptables！**

```
opkg update
opkg install ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy iptables-mod-socket ip6tables iptables
```

```
brook tproxy --link 'brook://...' --dnsListen :5353
```

1. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
2. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
3. 默認, OpenWrt 將會下發 router 的 IP 的為電腦或手機的網關和 DNS

### brook tproxy 透明代理网关在任意 Linux(有线)

**无需操作 iptables！**

```
systemctl stop systemd-resolved
systemctl disable systemd-resolved
echo nameserver 8.8.8.8 > /etc/resolv.conf
```

```
brook tproxy --link 'brook://...' --dnsListen 192.168.1.2:53 --disableAAAA
```

替换 192.168.1.2 为你的 Linux 的IP. 配置其他机器的网关和 DNS 为这台机器的 IP 即可

### 官网原版 OpenWrt 图形客户端

**无需操作 iptables！**

**端口 9999, 8888, 5353 将会被使用**. 它与 brook server, brook wsserver, brook wssserver, brook quicserver 一起工作.

1. 下載適合你系統的[ipk](https://github.com/txthinking/brook/releases)文件
2. 上傳並安裝: OpenWrt Web -> System -> Software -> Upload Package...
3. 刷新頁面, 頂部菜單會出現 Brook 按鈕
4. OpenWrt Web -> Brook -> 輸入後點擊 Connect
5. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
6. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. 默認, OpenWrt 將會下發 router 的 IP 為電腦或手機的網關和 DNS

### brook relay 可以中继任何 TCP 和 UDP server, 这是一个独立的功能, 它不依赖 brook server wsserver wssserver

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

### 例子不胜枚举，更多功能建议挨个看 CLI 文档的命令和参数吧，还有博客，YouTube 等...
