## $ brook hijackhttps

$ brook hijackhttps can create a magic DNS and force domains is https protocol. It works with a socks5 server. Assume your socks5 proxy is `127.0.0.1:1080`.

## Run brook hijackhttps

```
$ brook hijackhttps --socks5 127.0.0.1:1080
```

Then configure your system DNS server with `127.0.0.1`

> More parameters: $ brook hijackhttps -h

