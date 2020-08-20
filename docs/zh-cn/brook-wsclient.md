## $ brook wsclient

Assume your brook wsserver is `ws://1.2.3.4:9999` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook wsserver protocol | --> brook wsserver <--> a remote address
```

## Run brook wsclient

```
$ brook wsclient -s ws://1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook wsclient -h

## Use the socks5 proxy created by brook client

> TODO: Please help improve the documentation here

* Configure it on your system network settings
* Configure it on your browser
