# Other

## Script Syntax

I think just reading this one page is enough: [Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

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

## Example

Each `subcommand` has a `--example`, such as:

```
brook server --example
```

## Resources

| CLI | Description |
| --- | --- |
| [nami](https://github.com/txthinking/nami) | A clean and tidy decentralized package manager |
| [joker](https://github.com/txthinking/joker) | Joker can turn process into daemon. Zero-Configuration |
| [nico](https://github.com/txthinking/nico) | Nico can work with brook wsserver together |
| [z](https://github.com/txthinking/z) | z - process manager |
| [ipio](https://github.com/txthinking/ipio) | Proxy all traffic just one line command |
| [mad](https://github.com/txthinking/mad) | Generate root CA and derivative certificate for any domains and any IPs |
| [hancock](https://github.com/txthinking/hancock) | Manage multiple remote servers and execute commands remotely |
| [sshexec](https://github.com/txthinking/sshexec) | A command-line tool to execute remote command through ssh |
| [bash](https://github.com/txthinking/bash) | Many one-click scripts |
| [docker](https://hub.docker.com/r/txthinking/brook) | `docker run txthinking/brook` | 

| Resources | Description |
| --- | --- |
| [Protocol](https://github.com/txthinking/brook/tree/master/protocol) | Brook Protocol |
| [Blog](https://www.txthinking.com/talks/) | Some articles you should read |
| [YouTube](https://www.youtube.com/txthinking) | Some videos you should watch |
| [Telegram](https://t.me/txthinking) | Ask questions here |
| [Announce](https://t.me/s/txthinking_talks) | All news you should care |
| [GitHub](https://github.com/txthinking) | Other useful repos |
| [Socks5 Configurator](https://chromewebstore.google.com/detail/socks5-configurator/hnpgnjkeaobghpjjhaiemlahikgmnghb) | If you prefer CLI brook client | 
| [IPvBar](https://chromewebstore.google.com/detail/ipvbar/nepjlegfiihpkcdhlmaebfdfppckonlj) | See domain, IP and country in browser | 
| [TxThinking SSH](https://www.txthinking.com/ssh.html) | A SSH Terminal |
| [brook-store](https://github.com/txthinkinginc/brook-store) | A Brook User System |
| [TxThinking](https://www.txthinking.com) | Everything |

