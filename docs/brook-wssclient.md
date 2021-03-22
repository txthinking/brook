## $ brook wssclient

Assume your brook wssserver is `wss://domain.com:443` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook wssserver protocol | --> brook wssserver <--> a remote address
```

## Run brook wssclient

```
$ brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook wssclient -h

## Use the socks5 proxy

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

* In Chrome, install and configure extension [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
