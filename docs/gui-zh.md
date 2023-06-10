# 图形客户端文档

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## 本文适用的软件

-   [Brook](https://github.com/txthinking/brook)
-   [Brook Plus](https://www.txthinking.com/brook.html)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Windows GUI proxy 模式, Linux GUI proxy 模式

这个模式比较简单，会创建:

-   Socks5 代理: `socks5://[::1]:1080` 或 `socks5://127.0.0.1:1080`
-   HTTP 代理: `http://[::1]:8010` 或 `http://127.0.0.1:8010`
-   PAC: `http://127.0.0.1:1093/proxy.pac` 或 `http://[::1]:1093/proxy.pac` 基于 Bypass Domain 列表
-   Windows GUI 同时会配置 PAC 到系统代理。Linux GUI 可以配合 [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
-   什么是 socks5 和 http proxy? [文章](https://www.txthinking.com/talks/articles/socks5-and-http-proxy.article) 和 [视频](https://www.youtube.com/watch?v=Tb0_8odTxEI)

## iOS, Mac, Android, Windows TUN 模式, Linux TUN 模式

```
所谓的互联网连接，是 IP 连接 IP，不是连接域名。所以域名会被先解析成IP再决定怎么去连接。
```

## 配置介绍

| 配置    | 支持系统                      | 条件                 | 描述 |
| ------------- | ----------------------------- | ------------------------------------------- |  --- |
| 导入服务器    | iOS,Android,Mac,Windows,Linux | /                                           | brook link 列表                                                                                                                                                                                             |
| 系统 DNS    | iOS,Android,Mac,Windows,Linux | /                                           | 系统 DNS. **不要 bypass 此 IP**                                                                                                                                                                             |
| 虚拟 DNS      | iOS,Android,Mac,Windows,Linux | [如何避免 Brook 的 虚拟 DNS 不生效](https://www.txthinking.com/talks/articles/brook-fakedns.article) | 解析域名为 Fake IP，发起连接时会再转换为域名，然后把域名地址送到服务端进行代理，同时由服务端来负责域名解析                                                                                                  |
| 屏蔽         | iOS,Android,Mac,Windows,Linux | /                                           | Block 开关                                                                                                                                                                                                  |
| 屏蔽域名  | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 域名列表，匹配的域名会被阻断解析. **域名是后缀匹配模式**                                                                                                                                                    |
| 跳过        | iOS,Android,Mac,Windows,Linux | /                                           | Bypass 开关                                                                                                                                                                                                 |
| 跳过 IP     | iOS,Android,Mac,Windows,Linux | /                                           | CIDR 列表，匹配到的 IP 会被 bypass                                                                                                                                                                          |
| 跳过 Geo IP | iOS,Android,Mac,Windows,Linux                   | /                                           | 匹配到的 IP 会被 bypass. 提示: 全球 IP 变动频繁, 所以 Geo 库有时效性                                                                                                                                        |
| 跳过 Apps   | Android, Mac                       | /                                           | 这些 App 会被 bypass                                                                                                                                                                                        |
|  跳过 DNS    | iOS,Android,Mac,Windows,Linux | /                                           | 支持普通 DNS, 比如 `223.5.5.5:53`, 支持 DoH, 但需要通过参数 address 指定 DoH 的地址, 比如 `https://dns.alidns.com/dns-query?address=223.5.5.5%3A443` 用来解析 Bypass Domain. **此 DNS 的 IP 会自动 Bypass** |
| 跳过域名 | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 域名列表，匹配的域名会使用 Bypass DNS 解析来得到 IP, **最终连接是否会被 Bypass，还取决于 Bypass IP**. **域名是后缀匹配模式.** 当然也可以用脚本直接跳过域名而无关其IP               |
| Hosts         | iOS,Android,Mac,Windows,Linux | /                                           | Hosts 开关                                                                                                                                                                                                  |
| Host 列表    | iOS,Android,Mac,Windows,Linux | Fake DNS: 开启                              | 给域名指定 IP, v4, v6，如果值为空效果同 Block                                                                                                                                                               |
| 可编程  | iOS,Android,Mac,Windows,Linux | /                                           | 可编程开关                                                                                                                                                                                                  |
| 脚本        | iOS,Android,Mac,Windows,Linux | /                                           | 脚本。可以控制上面所有的功能。以及上面没有的功能，**全流程控制一切，具体看下文**。                                                                                                                          |
| 日志           | iOS,Android,Mac,Windows,Linux | /                                           | 日志开关                                                                                                                                                                                                    |
| 活动 | iOS,Android,Mac,Windows,Linux | / | 网络活动 |
| MITM | iOS,Android,Mac,Windows,Linux | / | 抓包，修改包，比如 https 的请求响应，十六进制，JSON，图片等 |
| TUN           | iOS,Android,Mac,Windows,Linux             | /                                           | 选择 Proxy/TUN/App 模式. [macOS bug](https://www.txthinking.com/talks/articles/macos-bug.article). iOS,Android,Mac 默认 TUN 模式 |
| 抓我    | iOS,Android,Mac,Windows,Linux | /                                           | 测试你的抓包或代理软件工作在系统代理还是 TUN                                                                                                                                                                |
| DNS 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | DNS 客户端                                                                                                                              |
| DOH 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | DOH 客户端                                                                                                                              |
| Echo 客户端    | iOS,Android,Mac,Windows,Linux | /                                           | Echo 客户端                                                                    |
| 测试 Socks5    | iOS,Android,Mac,Linux | /                                           | Test socks5 server                                                                    |
| 深色主题      | iOS,Android,Mac,Windows,Linux | /                                           | 暗黑模式                                                                                                                                                                                                    |     |
| 快捷方式      | iOS,Android,Mac,Windows,Linux | /                                           | 在首页快捷控制打开菜单里的功能                                                                                                                                                                              |
| 系统托盘      | Windows                       | /                                           | 以系统托盘形式打开，然后从系统托盘处打开控制面板                                                                                                                                                            |

## Programmable

```
Brook GUI 会在不同时机向脚本传入不同的全局变量，脚本只需要将处理结果赋值到全局变量 out 即可
```

### 传入变量介绍

| 变量                           | 类型 | 条件          | 时机                   | 描述                       | out 类型 |
| ------------------------------ | ---- | ------------- | ---------------------- | -------------------------- | -------- |
| in_brooklinks                  | map  | / | 连接之前  | 预定义多个 brook link，之后可编程指定连接哪个 | map      |
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
| appid   | string | App ID. 仅 Mac | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key          | 类型   | 描述                                                                                               | 示例    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | 是否 Block, 默认 `false`. 与 GUI Block Domain 是或的关系                                           | false   |
| ip           | string | 直接指定 IP，仅当 `type` 为 `A`/`AAAA`有效                                                         | 1.2.3.4 |
| forcefakedns | bool   | 忽略 GUI Bypass Domain，使用 Fake DNS 来处理，仅当 `type` 为 `A`/`AAAA`有效，默认 `false`          | false   |
| system       | bool   | 使用 System DNS 来解析，默认 `false`                                                               | false   |
| bypass       | bool   | 是否 Bypass, 默认 `false`, 如果为 `true` 则使用 Bypass DNS 来解析. 与 GUI Bypass Domain 是或的关系 | false   |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的 brook link | 自定义名字   |

### in_address

| Key           | 类型   | 描述                                                                     | 示例           |
| ------------- | ------ | ------------------------------------------------------------------------ | -------------- |
| network       | string | 即将发起连接网络，取值 `tcp`/`udp`                                       | tcp            |
| ipaddress     | string | IP 类型的地址，与 domainaddress 只会存在一个。注意这两个之间没有任何关系 | 1.2.3.4:443    |
| domainaddress | string | 域名类型的地址，因为 FakeDNS 我们这里才能拿到域名地址                    | google.com:443 |
| appid   | string | App ID. 仅 Mac | com.google.Chrome.helper          |
| interface   | string | 网络接口. 仅 Mac | en0          |

`out`, 如果是 `error` 类型会被记录在日志。如果不是 `map` 类型则会被忽略

| Key                    | 类型   | 描述                                                                                                                                                                           | 示例        |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------- |
| block                  | bool   | 是否 Block, 默认 `false`                                                                                                                                                       | false       |
| ipaddress              | string | IP 类型地址，重写目的地                                                                                                                                                        | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | 使用 Bypass DNS 获取`A`或`AAAA` IP 并重写目的地, 仅当 `domainaddress` 存在时有效，取值 `A`/`AAAA`                                                                              | A           |
| bypass                 | bool   | 是否 Bypass, 默认 `false`. 如果为 `true` 并且是 `domainaddress`, 那么必须指定 `ipaddress` 或 `ipaddressfrombypassdns`. 与 GUI Bypass IP 是或的关系 | false       |
| mitm                   | bool   | 是否进行 MITM, 默认 `false`. 仅当 `network` 为 `tcp` 时有效. 需要安装 CA，看下文介绍                                                                                           | false       |
| mitmprotocol           | string | 需要明确指定 MITM 协议, 取值 `http`/`https`                                                                                                                                    | https       |
| mitmcertdomain         | string | MITM 证书域名，默认从`domainaddress`里取。如果是 `ipaddress` 且 `mitm` 为 `true` 且 `mitmprotocol` 为 `https` 那么必须明确指定                                                 | example.com |
| mitmwithbody           | bool   | 是否操作 http body，默认 `false`. 会将请求和响应的 body 读取到内存里和脚本交互。iOS 50M 总内存限制可能会杀进程                                                             | false       |
| mitmautohandlecompress | bool   | 和脚本交互时是否自动解压缩 http body, 默认 `false`                                                                                                                             | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                             | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                           | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                            | 0           |
| brooklinkkey | string   | 当需要连接代理服务器时，转而连接 通过 in_brooklinks 的 key 指定的 brook link | 自定义名字   |

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

    * os: string, linux/darwin/windows/ios/android. If ios app run on mac, it is ios
    * iosapponmac: bool, ios app run on mac

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

## 独立脚本例子

https://github.com/txthinking/bypass

## 脚本生成器

https://modules.brook.app

## 抓包

-   [Brook 搭配 mitmproxy 进行手机 App 深度抓包](https://www.txthinking.com/talks/articles/brook-mitmproxy.article)
-   [Brook 全平台抓包](https://www.txthinking.com/talks/articles/brook-packet-capture.article)
-   [用 mitmproxy helper 抓包](https://www.txthinking.com/mitmproxy.html)

## 安装 CA

https://txthinking.github.io/ca/ca.pem

### iOS

https://www.youtube.com/watch?v=HSGPC2vpDGk

### Android

Android 分系统 CA 和用户 CA，必须要 ROOT 后安装到系统 CA 里

### macOS

```
nami install mad ca.txthinking
sudo mad install --ca ~/.nami/bin/ca.pem
```

### Windows

打开 GitBash

```
nami install mad ca.txthinking
```

然后用管理员打开 GitBash

```
mad install --ca ~/.nami/bin/ca.pem
```

注意 GitBash 或 Firefox 等软件可能不读取系统 CA，安装后可以用系统 Edge 浏览器测试

## Apple 推送问题

要接收推送，Apple Server 只允许 Ethernet, cellular data, Wi-Fi 连接. 所以你需要 Bypass 掉相关域名和 IP. [参考链接](https://github.com/txthinking/bypass/tree/master/apple)
