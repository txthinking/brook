## $ brook client

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook server protocol | --> brook server <--> a remote address
```

## Run brook client

```
$ brook client --server 1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook client -h

## Use the socks5 proxy

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

* In Chrome, install and configure extension [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
