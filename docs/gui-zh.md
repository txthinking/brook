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

