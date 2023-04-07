# Brook GUI Documentation

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## Software for which this article applies

-   [Brook](https://github.com/txthinking/brook)
-   [Brook Plus](https://www.txthinking.com/brook.html)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Intel Mac GUI proxy mode, Windows GUI proxy mode, Linux GUI proxy mode

This mode is very simple, will create:

-   Socks5 proxy: `socks5://[::1]:1080` or `socks5://127.0.0.1:1080`
-   HTTP proxy: `http://[::1]:8010` or `http://127.0.0.1:8010`
-   PAC: `http://127.0.0.1:1093/proxy.pac` or `http://[::1]:1093/proxy.pac` based on Bypass Domain list
-   Intel Mac GUI, Windows GUI set PAC to system proxy。Linux GUI can work with [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
-   What is socks5 and http proxy? [Article](https://www.txthinking.com/talks/articles/socks5-and-http-proxy-en.article) and [Video](https://www.youtube.com/watch?v=sBCB-X7BoP8)

## iOS, M1 Mac GUI, Android GUI, Intel Mac GUI tun mode, Windows GUI tun mode, Linux GUI tun mode

```
The so-called Internet connection is IP to IP connection, not domain name connection. Therefore, the domain name will be resolved into IP before deciding how to connect.
```

## Configuration Introduction

| Configuration  | Support Systems                   | Conditions                                                                          | Description  |
| -------------- | --------------------------------- | ----------------------------------------------------------------------------------- | --- |
| Import Servers | iOS,Android,Mac,Windows,Linux     | /                                                                                   | brook link list                                                                                                                                                                                                                                                                       |
| System DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System DNS. **Do not bypass this IP**                                                                                                                                                                                                                                                 |
| Fake DNS       | iOS, Android, Mac, Windows, Linux | Turn off the security DNS that comes with the system/browser, see below for details | The domain name is resolved to Fake IP, which will be converted to a domain name when a connection is initiated, and then the domain name address will be sent to the server, and the server is responsible for domain name resolution                                                |
| Block          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Block switch                                                                                                                                                                                                                                                                          |
| Block Domain   | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Domain name list, matching domain names will be blocked. **Domain name is suffix matching mode**                                                                                                                                                                                      |
| Bypass         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Bypass switch                                                                                                                                                                                                                                                                         |
| Bypass IP      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | CIDR list, matched IP will be bypassed                                                                                                                                                                                                                                                |
| Bypass Geo IP  | iOS,Android                       | /                                                                                   | The matched IP will be bypassed. Note: Global IP changes frequently, so the Geo library is time-sensitive                                                                                                                                                                             |
| Bypass Apps    | Android                           | /                                                                                   | These apps will be bypassed                                                                                                                                                                                                                                                           |
| Bypass DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Support normal DNS, such as `223.5.5.5:53`, support DoH, but need to specify the address of DoH through the parameter address, such as `https://dns. alidns.com/dns-query?address=223.5.5.5%3A443` is used to resolve Bypass Domain. **The IP of this DNS will automatically Bypass** |
| Bypass Domain  | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | List of domain names, matching domain names will use Bypass DNS resolution to get IP, **whether the final connection will be bypassed depends on the Bypass IP** . **The domain name is a suffix matching pattern**                                                                   |
| Hosts          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Hosts switch                                                                                                                                                                                                                                                                          |
| Hosts List     | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Specify IP, v4, v6 for the domain name, if the value is empty, the effect is the same as Block                                                                                                                                                                                        |
| Programmable   | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Programmable switch                                                                                                                                                                                                                                                                   |
| Script         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Script. All functions above can be controlled. And more and more, **The whole process can be controled, see below for details**.                                                                                                                                                      |
| Log            | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log switch                                                                                                                                                                                                                                                                            |
| Log View       | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log List                                                                                                                                                                                                                                                                              |
| Log View Plus  | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log list, easier to read, filter conditions, etc.                                                                                                                                                                                                                                     |
| MITM Log View  | iOS,Android,Mac,Windows,Linux     | /                                                                                   | MITM log list, such as https request response, hexadecimal, JSON, image, etc.                                                                                                                                                                                                         |
| TUN            | Mac,Windows,Linux                 | /                                                                                   | Choose proxy mode or tun mode                                                                                                                                                                                                                                                         |
| Capture Me     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Test your packet capture or proxy software is working as a system proxy or TUN                                                                                                                                                                                                        |
| Dark Mode      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System / Light / Dark                                                                                                                                                                                                                                                                 |     |
| Shortcut       | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Quickly control the functions in the menu on the home page                                                                                                                                                                                                                            |
| System Tray    | Windows                           | /                                                                                   | Open as systray, then open dashboard from the systray                                                                                                                                                                                                                                 |

## Programmable

```
Brook GUI will pass different global variables to the script at different times, and the script only needs to assign the processing result to the global variable out
```

Take full control of your own network

-   Like turning off IPv6 by blocking AAAA
-   Block system/browser built-in secure DNS
-   Override DST
-   Flexible and finer rules
-   Directly bypass the domain name regardless of whether the resolved IP is in Bypass
-   MITM decrypt HTTPS
-   Packet capture
-   Packet modify
-   Disable HTTP3
-   more and more...

### Introduction to incoming variables

| variable                       | type | condition   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_guiconfig                   | map  | /           | before connected                  | to override GUI configuration                     | map      |
| in_dnsquery                    | map  | FakeDNS: On | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | /           | When connecting to an address     | script can decide how to connect                  | map      |
| in_httprequest                 | map  | /           | When an HTTP(S) request comes in  | the script can decide how to handle this request  | map      |
| in_httprequest,in_httpresponse | map  | /           | when an HTTP(S) response comes in | the script can decide how to handle this response | map      |

### in_guiconfig

| Key | Type | Description                                       |
| --- | ---- | ------------------------------------------------- |
| \_  | bool | For future compatibility, this key can be ignored |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`, if it is `map` then explicitly specify each configuration item.

| Key        | Type   | Description       |
| ---------- | ------ | ----------------- |
| systemdns4 | string | System DNS v4     |
| systemdns6 | string | System DNS v6     |
| fakedns    | bool   | Fake DNS switch   |
| block      | bool   | GUI Block switch  |
| bypass     | bool   | GUI Bypass switch |
| bypassdns4 | string | Bypass DNS v4     |
| bypassdns6 | string | Bypass DNS v6     |
| hosts      | bool   | GUI Hosts switch  |

### in_dnsquery

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| domain | string | domain name | google.com |
| type   | string | query type  | A          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`. It is an OR relationship with GUI Block Domain                                                | false   |
| ip           | string | Specify IP directly, only valid when `type` is `A`/`AAAA`                                                                     | 1.2.3.4 |
| forcefakedns | bool   | Ignore GUI Bypass Domain, handle with Fake DNS, only valid when `type` is `A`/`AAAA`, default `false`                         | false   |
| system       | bool   | Get IP from system DNS, default `false`                                                                                       | false   |
| bypass       | bool   | whether to Bypass, default `false`, if `true` then use bypass DNS to resolve. It is an OR relationship with GUI Bypass Domain | false   |

### in_address

| Key           | Type   | Description                                                                                                         | Example        |
| ------------- | ------ | ------------------------------------------------------------------------------------------------------------------- | -------------- |
| network       | string | Network type, the value `tcp`/`udp`                                                                                 | tcp            |
| ipaddress     | string | IP type address. There is only of ipaddress and domainaddress. Note that there is no relationship between these two | 1.2.3.4:443    |
| domainaddress | string | Domain type address, because of FakeDNS we can get the domain name address here                                     | google.com:443 |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key                    | Type   | Description                                                                                                                                                                                             | Example     |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| block                  | bool   | Whether Block, default `false`                                                                                                                                                                          | false       |
| ipaddress              | string | IP type address, rewrite destination                                                                                                                                                                    | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | Use Bypass DNS to obtain `A` or `AAAA` IP and rewrite the destination, only valid when `domainaddress` exists, the value `A`/`AAAA`                                                                     | A           |
| bypass                 | bool   | **Only available on iOS, Android.** Bypass, default `false`. If `true` and `domainaddress`, then `ipaddress` or `ipaddressfrombypassdns` must be specified. It is an OR relationship with GUI Bypass IP | false       |
| mitm                   | bool   | Whether to perform MITM, default `false`. Only valid when `network` is `tcp`. Need to install CA, see below                                                                                             | false       |
| mitmprotocol           | string | MITM protocol needs to be specified explicitly, the value is `http`/`https`                                                                                                                             | https       |
| mitmcertdomain         | string | The MITM certificate domain name, which is taken from `domainaddress` by default. If `ipaddress` and `mitm` is `true` and `mitmprotocol` is `https` then must be must be specified explicitly           | example.com |
| mitmwithbody           | bool   | Whether to manipulate the http body, default `false`. **will read the body of the request and response into the memory and interact with the script. iOS 50M total memory limit may kill process**      | false       |
| mitmautohandlecompress | bool   | Whether to automatically decompress the http body when interacting with the script, default `false`                                                                                                     | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                                                      | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                                                    | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                                                     | 0           |

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

### How to write Tengo script

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

Example

https://github.com/txthinking/bypass/blob/master/example_script.tengo

### How to debug script

-   It is recommended to use [tun2brook](https://github.com/txthinking/tun2brook) on desktop to debug with print
-   It is recommended to use [mitmproxy helper](https://www.txthinking.com/mitmproxy.html) and [Wireshark Helper](https://www.txthinking.com/wireshark.html) to capture packets to determine what to modify

## Why and How to Turn Off System and Browser Security DNS

Because if Security DNS is turned on, the Fake DNS will not work. So we have to turn it off:

-   Android: Settings -> Network & internet -> Private DNS -> Off
-   Chrome on Mobile: Settings -> Privacy and security -> Use secure DNS -> Off
-   Chrome on Desktop: Settings -> Privacy and security -> Security -> Use secure DNS -> Off
-   Windows: Windows Settings -> Network & Internet -> Your Network -> DNS settings -> Edit -> Preferred DNS -> Unencrypted only -> 8.8.8.8

Other systems and software, please find out whether it exists and how to close it

## Install CA

https://txthinking.github.io/ca/ca.pem

### iOS

https://www.youtube.com/watch?v=HSGPC2vpDGk

### Android

Android has user CA and system CA, must be installed in the system CA after ROOT

### macOS

```
nami install mad ca.txthinking
sudo mad install --ca ~/.nami/bin/ca.pem
```

### Windows

Open GitBash

```
nami install mad ca.txthinking
```

Then open GitBash with admin

```
mad install --ca ~/.nami/bin/ca.pem
```

Note that software such as GitBash or Firefox may not read the system CA, you can use the system Edge browser to test after installation

## Apple Push Problem

To receive push, Apple Server only allows Ethernet, cellular data, Wi-Fi connections. So you need to Bypass the relevant domain name and IP. [Reference link](https://support.apple.com/en-us/HT210060)
