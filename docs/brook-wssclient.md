## $ brook wssclient

Assume your brook wssserver is `wss://domain.com:443` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook wssserver protocol | --> brook wssserver <--> a remote address
```

## Run brook wssclient

```
$ brook wssclient -s wss://domain.com:443 -p hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook wssclient -h

## Use the socks5 proxy

> TODO: Please help improve the documentation here

* Configure it on your system network settings
* Configure it on your browser
