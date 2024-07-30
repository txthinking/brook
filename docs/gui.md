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
