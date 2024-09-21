# Brook
<!--SIDEBAR-->
<!--G-R3M673HK5V-->
A cross-platform programmable network tool.

# Sponsor
**❤️  [Shiliew - A network app designed for those who value their time](https://www.txthinking.com/shiliew.html)**
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

## Client

- [iOS](https://apps.apple.com/us/app/brook-network-tool/id1216002642)
- [Android](https://github.com/txthinking/brook/releases/latest/download/Brook.apk)
- [macOS](https://apps.apple.com/us/app/brook-network-tool/id1216002642)
- [Windows](https://github.com/txthinking/brook/releases/latest/download/Brook.msix)
- [Linux](https://github.com/txthinking/brook/releases/latest/download/Brook.bin)
- [OpenWrt](https://github.com/txthinking/brook/releases)

> You may want to use `brook link` to customize some parameters

- [About App Mode on macOS](https://www.txthinking.com/talks/articles/macos-app-mode-en.article)
- [How to install Brook on Windows?](https://www.txthinking.com/talks/articles/msix-brook-en.article)
- [How to install Brook on Linux](https://www.txthinking.com/talks/articles/linux-app-brook-en.article)
- [How to install Brook on OpenWrt](https://www.txthinking.com/talks/articles/brook-openwrt-en.article)

# Client

Brook GUI will pass different _global variables_ to the script at different times, and the script only needs to assign the processing result to the global variable `out`

## CLI

Before discussing the GUI client, let's first talk about the command line client `brook`. As we know, after you have deployed the server, you can use the command line client `brook` to create a local socks5 proxy or http proxy on your machine, and then configure it in your system proxy settings or in your browser to use this proxy. However:

1. Not all apps will use this proxy, whether they use it is up to the app itself.
2. Generally, all UDP protocols will not go through this proxy, such as http3.

For the specifics of socks5 and http proxy, you can read [this article](https://www.txthinking.com/talks/articles/socks5-and-http-proxy-en.article).

## GUI

The GUI client does not use socks5 and http proxy mode, so there is no issue with some software not using the system proxy. Instead, it uses a virtual network card to take over the entire system's network, including UDP-based http3. Moreover, Brook allows us to control network requests programmatically, so it is necessary to have basic knowledge of network requests.

## Without Brook: Basic Knowledge of Network Requests

> Note: When we talk about addresses, we mean addresses that include the port number, such as a domain address: `google.com:443`, or an IP address: `8.8.8.8:53`

1. When an app requests a domain address, such as `google.com:443`
2. It will first perform a DNS resolution, which means that the app will send a network request to the system-configured DNS, such as `8.8.8.8:53`, to inquire about the IP of `google.com`
   1. The system DNS will return the IP of `google.com`, such as `1.2.3.4`, to the app
3. The app will combine the IP and port into an IP address, such as: `1.2.3.4:443`
4. The app makes a network request to this IP address `1.2.3.4:443`
5. The app receives the response data

In the above process, the app actually makes two network requests: one to the IP address `8.8.8.8:53` and another to the IP address `1.2.3.4:443`. In other words, the domain name is essentially an alias for the IP, and must obtain the domain's IP to establish a connection.

## With Brook: Fake DNS On

Brook has a Fake DNS feature, which can parse the domain name out of the query requests that an app sends to the system DNS and decide how to respond to the app.

1. When an app requests a domain name address, such as `google.com:443`
2. A DNS resolution will be performed first. That is, the app will send a network request to the system-configured DNS, such as `8.8.8.8:53`, to inquire about the IP of `google.com`
3. The Brook client detects that an app is sending a network request to `8.8.8.8:53`. <mark>This will trigger the `in_dnsquery` variable, carrying information such as `domain`</mark>
   1. The Brook client returns a fake IP to the app, such as `240.0.0.1`
4. The app combines the IP and port into an IP address, such as: `240.0.0.1:443`
5. The app makes a network request to the IP address `240.0.0.1:443`
6. The Brook client detects that an app is sending a network request to `240.0.0.1:443`, discovers that this is a fake IP, and will convert the fake IP address back to the domain address `google.com:443`. <mark>This will trigger the `in_address` variable, carrying information such as `domainaddress`</mark>
   1. The Brook client sends `google.com:443` to the Brook Server
   2. The Brook Server first requests its own DNS to resolve the domain name to find out the IP of `google.com`, such as receiving `1.2.3.4`
   3. The Brook Server combines the IP and port into an IP address, such as: `1.2.3.4:443`
   4. The Brook Server sends a network request to `1.2.3.4:443` and returns the data to the Brook client
   5. The Brook client then returns the data to the app
7. The app receives the response data

However, if the following situations occur, the domain name will not/cannot be parsed, meaning that the Brook client will not/cannot know what the domain name is and will treat it as a normal request sent to an IP address:

- Fake DNS not enabled: in this case, the Brook client will not attempt to parse the domain name from the request sent to the system DNS and will treat it as a normal request sent to an IP address.
- Even with Fake DNS enabled, but the app uses the system's secure DNS or the app's own secure DNS: in this case, the Brook client cannot parse the domain name from the request sent to the secure DNS and will treat it as a normal request sent to an IP address.

To avoid the ineffectiveness of Fake DNS, please refer to [this article](https://www.txthinking.com/talks/articles/brook-fakedns-en.article).

## With Brook: Fake DNS Off

1. When an app requests a domain address, such as `google.com:443`
2. A DNS resolution will be performed first. That is, the app will send a network request to the system-configured DNS, such as `8.8.8.8:53`, to inquire about the IP of `google.com`
3. The Brook client detects that an app is sending a network request to `8.8.8.8:53`. <mark>This will trigger the `in_address` variable, carrying information such as `ipaddress`</mark>
   1. The Brook client sends `8.8.8.8:53` to the Brook Server
   2. The Brook Server sends a network request to `8.8.8.8:53` and returns the result, such as `1.2.3.4`, to the Brook client
   3. The Brook client then returns the result to the app
4. The app combines the IP and port into an IP address, such as: `1.2.3.4:443`
5. The app makes a network request to the IP address `1.2.3.4:443`
6. The Brook client detects that an app is sending a network request to `1.2.3.4:443`. <mark>This will trigger the `in_address` variable, carrying information such as `ipaddress`</mark>
   1. The Brook client sends `1.2.3.4:443` to the Brook Server
   2. The Brook Server sends a network request to `1.2.3.4:443` and returns the data to the Brook client
   3. The Brook client then returns the data to the app
7. The app receives the response data

## With Brook: Fake DNS On, But the App Uses the System's Secure DNS or Its Own Secure DNS

1. When an app requests a domain name address, such as `google.com:443`
2. A DNS resolution will be performed first. That is, the app will send a network request to the secure DNS, such as `8.8.8.8:443`, to inquire about the IP of `google.com`
3. The Brook client detects that an app is sending a network request to `8.8.8.8:443`. <mark>This will trigger the `in_address` variable, carrying information such as `ipaddress`</mark>
   1. The Brook client sends `8.8.8.8:443` to the Brook Server
   2. The Brook Server sends a network request to `8.8.8.8:443`, and returns the result, such as `1.2.3.4`, to the Brook client
   3. The Brook client then returns the result to the app
4. The app combines the IP and port into an IP address, such as: `1.2.3.4:443`
5. The app makes a network request to the IP address `1.2.3.4:443`
6. The Brook client detects that an app is sending a network request to `1.2.3.4:443`. <mark>This will trigger the `in_address` variable, carrying information such as `ipaddress`</mark>
   1. The Brook client sends `1.2.3.4:443` to the Brook Server
   2. The Brook Server sends a network request to `1.2.3.4:443` and returns the data to the Brook client
   3. The Brook client then returns the data to the app
7. The app receives the response data

To avoid the ineffectiveness of Fake DNS, please refer to [this article](https://www.txthinking.com/talks/articles/brook-fakedns-en.article).

## Handle Variable Trigger

- When the `in_brooklinks` variable is triggered:
    - This is currently the only variable that gets triggered before the Brook client starts.
    - We know that Brook starts with your choice of a Brook Server, and this variable lets you specify multiple Brook Servers.
    - Then during runtime, you can use one of these Brook Servers as needed.
- When the `in_dnsquery` variable is triggered, you can process as needed, such as:
    - Blocking, such as to prevent ad domain names.
    - Directly specifying the response IP.
    - Letting the system DNS resolve this domain.
    - Letting Bypass DNS resolve this domain.
    - And so on.
- When the `in_address` variable is triggered, you can process as needed, such as:
    - Block this connection.
    - Rewrite the destination.
    - If it's a domain address, you can specify that Bypass DNS is responsible for resolving the IP of this domain.
    - Allow it to connect directly without going through a proxy.
    - If it's HTTP/HTTPS, you can start MITM (Man-In-The-Middle), which will subsequently trigger `in_httprequest` and `in_httpresponse`.
    - And so on.
- When the `in_httprequest` variable is triggered, you can process as needed, such as:
    - Modifying the HTTP request.
    - Returning a custom HTTP response directly.
- When the `in_httpresponse` variable is triggered, you can process as needed, such as:
    - Modifying the HTTP response.

For detailed information on the properties and responses of variables, please refer to the following content.

## Variables

| variable                       | type | condition   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_brooklinks                  | map  | / | Before connecting  | Predefine multiple brook links, and then programmatically specify which one to connect to | map      |
| in_dnsquery                    | map  | FakeDNS: On | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | /           | When connecting to an address     | script can decide how to connect                  | map      |
| in_httprequest                 | map  | /           | When an HTTP(S) request comes in  | the script can decide how to handle this request  | map      |
| in_httprequest,in_httpresponse | map  | /           | when an HTTP(S) response comes in | the script can decide how to handle this response | map      |

## in_brooklinks

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

## in_dnsquery

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

## in_address

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

## in_httprequest

| Key    | Type   | Description                   | Example                     |
| ------ | ------ | ----------------------------- | --------------------------- |
| URL    | string | URL                           | `https://example.com/hello` |
| Method | string | HTTP method                   | GET                         |
| Body   | bytes  | HTTP request body             | /                           |
| ...    | string | other fields are HTTP headers | /                           |

`out`, must be set to a request or response

## in_httpresponse

| Key        | Type   | Description                   | Example |
| ---------- | ------ | ----------------------------- | ------- |
| StatusCode | int    | HTTP status code              | 200     |
| Body       | bytes  | HTTP response body            | /       |
| ...        | string | other fields are HTTP headers | /       |

`out`, must be set to a response

## Modules

In Brook GUI, scripts are abstracted into **Modules**. There are already [some modules](https://github.com/txthinking/brook/blob/master/programmable/modules/), and thre is no magic, it just automatically combine [_header.tengo](https://github.com/txthinking/brook/blob/master/programmable/modules/_header.tengo) and [_footer.tengo](https://github.com/txthinking/brook/blob/master/programmable/modules/_footer.tengo), so you only need to write the module itself.

```
modules = append(modules, {
    // If you want to predefine multiple brook links, and then programmatically specify which one to connect to, then define `brooklinks` key a function
    brooklinks: func(m) {
        // Please refer to the example in `brooklinks.tengo`
    },
    // If you want to intercept and handle a DNS query, then define `dnsquery` key a function, `m` is the `in_dnsquery`
    dnsquery: func(m) {
        // Please refer to the example in `block_aaaa.tengo`
    },
    // If you want to intercept and handle an address, then define `address` key a function, `m` is the `in_address`
    address: func(m) {
        // Please refer to the example in `block_google_secure_dns.tengo`
    },
    // If you want to intercept and handle a http request, then define `httprequest` key a function, `request` is the `in_httprequest`
    httprequest: func(request) {
        // Please refer to the example in `ios_app_downgrade.tengo` or `redirect_google_cn.tengo`
    },
    // If you want to intercept and handle a http response, then define `httpresponse` key a function, `request` is the `in_httprequest`, `response` is the `in_httpresponse`
    httpresponse: func(request, response) {
        // Please refer to the example in `response_sample.tengo`
    }
})
```

## tun2brook

https://github.com/txthinking/tun2brook

If you are using tun2brook, you can manually combine multiple modules into a complete script in the following way. For example:

```
cat _header.tengo > my.tengo

cat block_google_secure_dns.tengo >> my.tengo
cat block_aaaa.tengo >> my.tengo

cat _footer.tengo >> my.tengo
```

## Syntax

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

## Debug

If you are writing complex scripts, the GUI may not be convenient for debugging. It is recommended to use [tun2brook](https://github.com/txthinking/tun2brook) on desktop to debug with `fmt.println`

## CA

https://txthinking.github.io/ca/ca.pem

| OS | How |
| --- | --- |
| iOS | https://www.youtube.com/watch?v=HSGPC2vpDGk |
| Android | Android has user CA and system CA, must be installed in the system CA after ROOT |
| macOS | `nami install mad ca.txthinking`, `sudo mad install --ca ~/.nami/bin/ca.pem` |
| Windows | `nami install mad ca.txthinking`, Admin: `mad install --ca ~/.nami/bin/ca.pem` |

> Some software may not read the system CA，you can use `curl --cacert ~/.nami/bin/ca.pem` to debug

## OpenWrt

[Brook OpenWRT: Perfectly supports IPv4/IPv6/TCP/UDP](https://www.txthinking.com/talks/articles/brook-openwrt-en.article)

## IPv6

Brook's stance on IPv6 is positive, if your server or local environment doesn't have an IPv6 stack, read [this article](https://www.txthinking.com/talks/articles/brook-ipv6-en.article).

## Troubleshooting Steps

1. After adding your Server to the Brook client
2. If your Server uses a domain and has not specified an IP address via `brook link --address`, then Brook client will attempt to resolve the domain's IP using local DNS, preferring AAAA record. For example:
   - domain.com:9999
   - ws://domain.com:9999
   - wss://domain.com:9999
   - quic://domain.com:9999
3. Connectivity check: Go to the Server details page and click `Connectivity Check`. If it works sometimes but not others, this indicates instability.
4. After connected
1. Brook will change your system DNS to the System DNS configured in Brook (by default Google's DNS). In very rare cases, this change may be ignored on Windows, you can confirm this in the system settings.
5. Test IPv4 TCP: Use `Test IPv4 TCP` for testing; this test has hardcoded the IP address, so does not trigger DNS resolution.
5. Test IPv4 UDP: Use `Test IPv4 UDP` for testing; this test has hardcoded the IP address, so does not trigger DNS resolution.
6. Test IPv6 TCP: Use `Test IPv6 TCP` for testing; this test has hardcoded the IP address, so does not trigger DNS resolution.
6. Test IPv6 UDP: Use `Test IPv6 UDP` for testing; this test has hardcoded the IP address, so does not trigger DNS resolution.
7. Test TCP and UDP: Use the `Echo Client` for testing. If the echo server entered is a domain address, it will trigger DNS resolution.
8. Ensure the effectiveness of Fake DNS: Fake DNS is essential to do something with a domain or domain address. Generally, enable the `Block Google Secure DNS` module is sufficient. For other cases, refer to [this article](https://www.txthinking.com/talks/articles/brook-fakedns-en.article).
9. If your local or Server does not support IPv6: Refer to [this article](https://www.txthinking.com/talks/articles/brook-ipv6-en.article).
10. macOS App Mode: Refer to [this article](https://www.txthinking.com/talks/articles/macos-app-mode-en.article).
11. Windows:
    - The client can pass the tests without any special configuration on a brand-new, genuine Windows 11.
    - Be aware that the Windows system time is often incorrect.
    - Do not have other similar network software installed; they can cause conflicting network settings in the system.
    - Try restarting the computer.
    - Windows Defender may ask for permission to connect to the network or present other issues.
    - System DNS may need to be set to 8.8.8.8 and/or 2001:4860:4860::8888
12. Android:
    - The client can pass the tests without any special configuration on the official Google ROM.
    - Different ROMs may have made different modifications to the system.
    - Permission for background running might require separate settings.
    - System DNS may need to be set to 8.8.8.8 and/or 2001:4860:4860::8888
13. Bypass traffic such as China, usually requires the following modules to be activated:
    - `Block Google Secure DNS`
    - `Bypass Geo`
    - `Bypass Apple`: To prevent issues receiving Apple message notifications.
    - `Bypass China domain` or `Bypass China domain A`: The former uses `Bypass DNS` to obtain the IP, then `Bypass Geo` or other modules decide whether to bypass; the latter bypasses directly after obtaining the IP with `Bypass DNS` using A records. The latter is needed if your local does not support IPv6.
    - If you are a [Shiliew](https://www.txthinking.com/shiliew.html) user, some modules are enabled by default, which is usually sufficient.
14. Search [GitHub issues](https://github.com/txthinking/brook/issues?q=is%3Aissue)
15. Read the [blog](https://www.txthinking.com/talks/)
16. Read the [documentation](https://brook.app)
14. Submit [new issue](https://github.com/txthinking/brook/issues?q=is%3Aissue)
17. Seek help in the [group](https://t.me/txthinking)
# CLI Documentation
Each subcommand has a `--example` parameter that can print the minimal example of usage
# NAME

Brook - A cross-platform programmable network tool

# SYNOPSIS

Brook

```
brook --help
```

**Usage**:

```
Brook [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

- **--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_cidr4.txt. Works with server/wsserver/wssserver/quicserver

- **--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_cidr6.txt. Works with server/wsserver/wssserver/quicserver

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt. Works with server/wsserver/wssserver/quicserver

- **--blockGeoIP**="": Block IP by Geo country code, such as US. Works with server/wsserver/wssserver/quicserver

- **--blockListUpdateInterval**="": Update list --blockDomainList,--blockCIDR4List,--blockCIDR6List interval, second. default 0, only read one time on start (default: 0)

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



- **--ipLimitInterval**="": Interval (s) for ipLimitMax (default: 0)

- **--ipLimitMax**="": Limit the number of client IP addresses, be careful when using this parameter, as the client may have dynamic IP. Works with server/wsserver/wssserver/quicserver (default: 0)

- **--ipLimitWait**="": How long (s) to wait for recovery after exceeding ipLimitMax (default: 0)

- **--log**="": Enable log. A valid value is file path or 'console'. Send SIGUSR1 to me to reset the log file on unix system. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true

- **--pid**="": A file path used to store pid. Send SIGUSR1 to me to reset the --serverLog file on unix system

- **--pprof**="": go http pprof listen addr, such as :6060

- **--prometheus**="": prometheus http listen addr, such as :7070. If it is transmitted on the public network, it is recommended to use it with nico

- **--prometheusPath**="": prometheus http path, such as /xxx. If it is transmitted on the public network, a hard-to-guess value is recommended

- **--serverHKDFInfo**="": server HKDF info, most time you don't need to change this, if changed, all and each brook links in client side must be same, I mean each (default: "brook")

- **--serverLog**="": Enable server log, traffic and more. A valid value is file path or 'console'. Send SIGUSR1 to me to reset the log file on unix system. Mutually exclusive with the --log parameter. Works with server/wsserver/wssserver/quicserver with brook protocol

- **--speedLimit**="": Limit speed (b), 500kb/s such as: 500000, works with server/wsserver/wssserver/quicserver (default: 0)

- **--tag**="": Tag can be used to the process, will be append into log or serverLog, such as: 'key1:value1'. All tags will also be appended as query parameters one by one to the userAPI

- **--userAPI**="": When you build your own user system, Brook Server will send GET request to your userAPI to check if token is valid, for example: https://your-api-server.com/a_unpredictable_path. Yes, it is recommended to add an unpredictable path to your https API, of course, you can also use the http api for internal network communication. The request format is https://your-api-server.com/a_unpredictable_path?token=xxx. When the response is 200, the body should be the user's unique identifier, such as user ID; all other status codes are considered to represent an illegitimate user, and in these cases, the body should be a string describing the error. It should be used with --serverLog and server/wsserver/wssserver/quicserver with brook protocol. For more information, please read https://github.com/txthinking/brook/blob/master/protocol/user.md

- **--userAPIInvalidCacheTime**="": Once a token is checked and invalid, the userAPI will not be requested to validate again for a certain period (s). A reasonable value must be set, otherwise it will affect the performance of each incoming connection. Note that this may affect the user experience, when you change the user status from invalid to valid in your user system (default: 1800)

- **--userAPIValidCacheTime**="": Once a token is checked and valid, the userAPI will not be requested to validate again for a certain period (s). A reasonable value must be set, otherwise it will affect the performance of each incoming connection (default: 3600)

- **--version, -v**: print the version


# COMMANDS

## server

Start a brook server that supports tcp and udp

- **--blockCIDR4List**="": This option will be removed in a future version, please use the global option instead

- **--blockCIDR6List**="": This option will be removed in a future version, please use the global option instead

- **--blockDomainList**="": This option will be removed in a future version, please use the global option instead

- **--blockGeoIP**="": This option will be removed in a future version, please use the global option instead

- **--example**: Show a minimal example of usage

- **--listen, -l**="": Listen address, like: ':9999'

- **--password, -p**="": Server password

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": This option will be removed in a future version, please use the global option instead (default: 0)

## client

Start a brook client that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]

- **--example**: Show a minimal example of usage

- **--http**="": Where to listen for HTTP proxy connections

- **--link**="": brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored

- **--password, -p**="": Brook server password

- **--server, -s**="": Brook server address, like: 1.2.3.4:9999

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## wsserver

Start a brook wsserver that supports tcp and udp. It opens a standard http server and a websocket server

- **--blockCIDR4List**="": This option will be removed in a future version, please use the global option instead

- **--blockCIDR6List**="": This option will be removed in a future version, please use the global option instead

- **--blockDomainList**="": This option will be removed in a future version, please use the global option instead

- **--blockGeoIP**="": This option will be removed in a future version, please use the global option instead

- **--example**: Show a minimal example of usage

- **--listen, -l**="": Listen address, like: ':80'

- **--password, -p**="": Server password

- **--path**="": URL path (default: /ws)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": This option will be removed in a future version, please use the global option instead (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

- **--xForwardedFor**: Replace the from field in --log, note that this may be forged

## wsclient

Start a brook wsclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]

- **--example**: Show a minimal example of usage

- **--http**="": Where to listen for HTTP proxy connections

- **--link**="": brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored

- **--password, -p**="": Brook wsserver password

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--wsserver, -s**="": Brook wsserver address, like: ws://1.2.3.4:80, if no path then /ws will be used. Do not omit the port under any circumstances

## wssserver

Start a brook wssserver that supports tcp and udp. It opens a standard https server and a websocket server

- **--blockCIDR4List**="": This option will be removed in a future version, please use the global option instead

- **--blockCIDR6List**="": This option will be removed in a future version, please use the global option instead

- **--blockDomainList**="": This option will be removed in a future version, please use the global option instead

- **--blockGeoIP**="": This option will be removed in a future version, please use the global option instead

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--example**: Show a minimal example of usage

- **--password, -p**="": Server password

- **--path**="": URL path (default: /ws)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": This option will be removed in a future version, please use the global option instead (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wssclient

Start a brook wssclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]

- **--example**: Show a minimal example of usage

- **--http**="": Where to listen for HTTP proxy connections

- **--link**="": brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored

- **--password, -p**="": Brook wssserver password

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--wssserver, -s**="": Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances

## quicserver

Start a brook quicserver that supports tcp and udp.

- **--blockCIDR4List**="": This option will be removed in a future version, please use the global option instead

- **--blockCIDR6List**="": This option will be removed in a future version, please use the global option instead

- **--blockDomainList**="": This option will be removed in a future version, please use the global option instead

- **--blockGeoIP**="": This option will be removed in a future version, please use the global option instead

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--example**: Show a minimal example of usage

- **--password, -p**="": Server password

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

- **--updateListInterval**="": This option will be removed in a future version, please use the global option instead (default: 0)

- **--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## quicclient

Start a brook quicclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]. (The global-dial-parameter is ignored)

- **--example**: Show a minimal example of usage

- **--http**="": Where to listen for HTTP proxy connections

- **--link**="": brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## relayoverbrook

Relay network traffic over brook, which supports TCP and UDP. Accessing [from address] is equal to accessing [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]

- **--example**: Show a minimal example of usage

- **--from, -f, -l**="": Listen address: like ':9999'

- **--link**="": brook link, you can get it via $ brook link. The server and password parameters will be ignored

- **--password, -p**="": Password

- **--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443

- **--tcpTimeout**="": time (s) (default: 0)

- **--to, -t**="": Address which relay to, like: 1.2.3.4:9999

- **--udpTimeout**="": time (s) (default: 0)

## dnsserveroverbrook

Run a dns server over brook, which supports TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt

- **--bypassDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

- **--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

- **--example**: Show a minimal example of usage

- **--link**="": brook link, you can get it via $ brook link. The server and password parameters will be ignored

- **--listen, -l**="": Listen address, like: 127.0.0.1:53

- **--password, -p**="": Password

- **--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## link

Generate a brook link

- **--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

- **--ca**="": When server is brook wssserver or brook quicserver, specify ca for untrusted cert, such as /path/to/ca.pem

- **--clientHKDFInfo**="": client HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is

- **--example**: Show a minimal example of usage

- **--fragment**="": When server is brook wssserver, split the ClientHello into multiple fragments and then send them one by one with delays (millisecond). The format is min_length:max_length:min_delay:max_delay, cannot be zero, such as 50:100:10:50

- **--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

- **--name**="": Give this server a name

- **--password, -p**="": Password

- **--server, -s**="": Support brook server, brook wsserver, brook wssserver, socks5 server, brook quicserver. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080, quic://google.com:443

- **--serverHKDFInfo**="": server HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is

- **--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

- **--token**="": A token represents a user's identity. A string encoded in hexadecimal. Server needs to have --userAPI enabled. Note that: Only supported by the brook GUI(except for OpenWrt) and tun2brook

- **--udpoverstream**: When server is brook quicserver, UDP over Stream. Under normal circumstances, you need this parameter because the max datagram size for QUIC is very small. Note: only brook CLI and tun2brook suppport for now

- **--udpovertcp**: When server is brook server, UDP over TCP

- **--username, -u**="": Username, when server is socks5 server

- **--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## connect

Run a client and connect with a brook link, which supports TCP and UDP. It can start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

- **--example**: Show a minimal example of usage

- **--http**="": Where to listen for HTTP proxy connections

- **--link, -l**="": brook link, you can get it via $ brook link

- **--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## relay

Run a standalone relay, which supports TCP and UDP. Accessing [from address] is equal to accessing [to address], [src <-> from address <-> to address]

- **--example**: Show a minimal example of usage

- **--from, -f, -l**="": Listen address: like ':9999'

- **--tcpTimeout**="": time (s) (default: 0)

- **--to, -t**="": Address which relay to, like: 1.2.3.4:9999

- **--udpTimeout**="": time (s) (default: 0)

## dnsserver

Run a standalone dns server

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

- **--example**: Show a minimal example of usage

- **--listen, -l**="": Listen address, like: 127.0.0.1:53

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## dnsclient

Send a dns query

- **--dns, -s**="": DNS server, such as 8.8.8.8:53 (default: 8.8.8.8:53)

- **--domain, -d**="": Domain

- **--example**: Show a minimal example of usage

- **--short**: Short for A/AAAA

- **--type, -t**="": Type, such as A (default: A)

## dohserver

Run a standalone doh server

- **--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt

- **--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

- **--disableA**: Disable A query

- **--disableAAAA**: Disable AAAA query

- **--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

- **--domainaddress**="": Such as: domain.com:443, if you want to create a https server. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

- **--example**: Show a minimal example of usage

- **--listen**="": listen address, if you want to create a http server behind nico

- **--path**="": URL path (default: /dns-query)

- **--tcpTimeout**="": time (s) (default: 0)

- **--udpTimeout**="": time (s) (default: 0)

## dohclient

Send a dns query

- **--doh, -s**="": DOH server, the address is required (default: https://dns.quad9.net/dns-query?address=9.9.9.9%3A443)

- **--domain, -d**="": Domain

- **--example**: Show a minimal example of usage

- **--short**: Short for A/AAAA

- **--type, -t**="": Type, such as A (default: A)

## dhcpserver

Run a standalone dhcp server. Other running dhcp servers need to be stopped.

- **--cache**="": Cache file, local absolute file path, default is $HOME/.brook.dhcpserver

- **--count**="": IP range from the start, which you want to assign to clients (default: 0)

- **--dnsserver**="": The dns server which you want to assign to clients, such as: 192.168.1.1 or 8.8.8.8

- **--example**: Show a minimal example of usage

- **--gateway**="": The router gateway which you want to assign to clients, such as: 192.168.1.1

- **--interface**="": Select interface on multi interface device. Linux only

- **--netmask**="": Subnet netmask which you want to assign to clients (default: 255.255.255.0)

- **--serverip**="": DHCP server IP, the IP of the this machine, you shoud set a static IP to this machine before doing this, such as: 192.168.1.10

- **--start**="": Start IP which you want to assign to clients, such as: 192.168.1.100

## socks5

Run a standalone standard socks5 server, which supports TCP and UDP

- **--example**: Show a minimal example of usage

- **--limitUDP**: The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.

- **--listen, -l**="": Socks5 server listen address, like: :1080 or 1.2.3.4:1080

- **--password**="": Password, optional

- **--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

- **--tcpTimeout**="": Connection deadline time (s) (default: 0)

- **--udpTimeout**="": Connection deadline time (s) (default: 0)

- **--username**="": User name, optional

## socks5tohttp

Convert a socks5 proxy to a http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]

- **--example**: Show a minimal example of usage

- **--listen, -l**="": HTTP proxy which will be create: like: 127.0.0.1:8010

- **--socks5, -s**="": Socks5 server address, like: 127.0.0.1:1080

- **--socks5password**="": Socks5 password, optional

- **--socks5username**="": Socks5 username, optional

- **--tcpTimeout**="": Connection tcp timeout (s) (default: 0)

## pac

Run a PAC server or save PAC to a file

- **--bypassDomainList, -b**="": One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt

- **--example**: Show a minimal example of usage

- **--file, -f**="": Save PAC to file, this will ignore listen address

- **--listen, -l**="": Listen address, like: 127.0.0.1:1980

- **--proxy, -p**="": Proxy, like: 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' (default: SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT)

## testsocks5

Test a socks5 server to see if it works properly

- **--dns**="": DNS server for connecting (default: 8.8.8.8:53)

- **--domain**="": Domain for query (default: http3.ooo)

- **--example**: Show a minimal example of usage

- **--password, -p**="": Socks5 password

- **--socks5, -s**="": Like: 127.0.0.1:1080

- **--username, -u**="": Socks5 username

- **-a**="": The A record of domain (default: 137.184.237.95)

## testbrook

Test UDP and TCP of a brook server/wsserver/wssserver/quicserver connection. (The global-dial-parameter is ignored)

- **--dns**="": DNS server for connecting (default: 8.8.8.8:53)

- **--domain**="": Domain for query (default: http3.ooo)

- **--example**: Show a minimal example of usage

- **--link, -l**="": brook link. Get it via $ brook link

- **--socks5**="": Temporarily listening socks5 (default: 127.0.0.1:11080)

- **-a**="": The A record of domain (default: 137.184.237.95)

## echoserver

Echo server, echo UDP and TCP address of routes

- **--example**: Show a minimal example of usage

- **--listen, -l**="": Listen address, like: ':7777'

## echoclient

Connect to echoserver, echo UDP and TCP address of routes

- **--example**: Show a minimal example of usage

- **--server, -s**="": Echo server address, such as 1.2.3.4:7777

- **--times**="": Times of interactions (default: 0)

## ipcountry

Get country of IP

- **--example**: Show a minimal example of usage

- **--ip**="": 1.1.1.1

## completion

Generate shell completions

- **--example**: Show a minimal example of usage

- **--file, -f**="": Write to file (default: brook_autocomplete)

## mdpage

Generate markdown page

- **--example**: Show a minimal example of usage

- **--file, -f**="": Write to file, default print to stdout







## manpage

Generate man.1 page

- **--example**: Show a minimal example of usage

- **--file, -f**="": Write to file, default print to stdout. You should put to /path/to/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos

## help, h

Shows a list of commands or help for one command

# Examples

List some examples of common scene commands, pay attention to replace the parameters such as IP, port, password, domain name, certificate path, etc. in the example by yourself

## Run brook server

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

## Run brook wsserver

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

## Run brook wssserver: automatically certificate

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

## Run brook wssserver Use a certificate issued by an existing trust authority

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

## Run brook wssserver issue untrusted certificates yourself, any domain

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

## withoutBrookProtocol

Better performance, but data is not strongly encrypted using Brook protocol. So please use certificate encryption, and it is not recommended to use --withoutBrookProtocol and --insecure together

## withoutBrookProtocol automatically certificate

> Make sure your domain has been resolved to your server IP successfully. Automatic certificate issuance requires the use of port 80

```
brook wssserver --domainaddress domain.com:443 --password hello --withoutBrookProtocol
```

get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

## withoutBrookProtocol Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --withoutBrookProtocol
```

## withoutBrookProtocol issue untrusted certificates yourself, any domain

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

## Run brook socks5, A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

then

-   server: `1.2.3.4:1080`

or get brook link

```
brook link --server socks5://1.2.3.4:1080
```

## Run brook socks5 with username and password. A stand-alone standard socks5 server

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

## brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

## brook dnsserveroverbrook can create a encrypted DNS server, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook dnsserveroverbrook ... --listen 127.0.0.1:53
```

## Brook OpenWRT Router: Perfectly supports IPv4/IPv6/TCP/UDP. Native IPv6

https://www.txthinking.com/talks/articles/brook-openwrt-en.article

## Turn macOS into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-macos-gateway-en.article

## Turn Windows into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-windows-gateway-en.article

## Turn Linux into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-linux-gateway-en.article

## brook relay can relay a address to a remote address. It can relay any tcp and udp server

```
brook relay --from :9999 --to 1.2.3.4:9999
```

## brook socks5tohttp can convert a socks5 to a http proxy

```
brook socks5tohttp --socks5 127.0.0.1:1080 --listen 127.0.0.1:8010
```

## There are countless examples; for more feature suggestions, it's best to look at the commands and parameters in the CLI documentation one by one, and blog, YouTube...
# Resources

| CLI | Description |
| --- | --- |
| [nami](https://github.com/txthinking/nami) | A clean and tidy decentralized package manager |
| [joker](https://github.com/txthinking/joker) | Joker can turn process into daemon. Zero-Configuration |
| [nico](https://github.com/txthinking/nico) | Nico can work with brook wsserver together |
| [z](https://github.com/txthinking/z) | z - process manager |
| [tun2brook](https://github.com/txthinking/tun2brook) | Proxy all traffic just one line command |
| [mad](https://github.com/txthinking/mad) | Generate root CA and derivative certificate for any domains and any IPs |
| [hancock](https://github.com/txthinking/hancock) | Manage multiple remote servers and execute commands remotely |
| [sshexec](https://github.com/txthinking/sshexec) | A command-line tool to execute remote command through ssh |
| [jb](https://github.com/txthinking/jb) | write script in an easier way than bash |
| [bash](https://github.com/txthinking/bash) | Many one-click scripts |
| [docker](https://hub.docker.com/r/txthinking/brook) | `docker run txthinking/brook` | 

| Resources | Description |
| --- | --- |
| [Protocol](https://github.com/txthinking/brook/tree/master/protocol) | Brook Protocol |
| [Blog](https://www.txthinking.com/talks/) | Some articles you should read |
| [YouTube](https://www.youtube.com/txthinking) | Some videos you should watch |
| [Telegram](https://t.me/txthinking) | Ask questions here |
| [Announce](https://t.me/s/txthinking_news) | All news you should care |
| [GitHub](https://github.com/txthinking) | Other useful repos |
| [Socks5 Configurator](https://chromewebstore.google.com/detail/socks5-configurator/hnpgnjkeaobghpjjhaiemlahikgmnghb) | If you prefer CLI brook client | 
| [IPvBar](https://chromewebstore.google.com/detail/ipvbar/nepjlegfiihpkcdhlmaebfdfppckonlj) | See domain, IP and country in browser | 
| [TxThinking SSH](https://www.txthinking.com/ssh.html) | A SSH Terminal |
| [brook-user-system](https://github.com/txthinkinginc/brook-user-system) | A Brook User System |
| [TxThinking](https://www.txthinking.com) | Everything |

