## $ brook wsclient

Assume your brook wsserver is `ws://1.2.3.4:9999` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook wsserver protocol | --> brook wsserver <--> a remote address
```

## Run brook wsclient

```
$ brook wsclient --wsserver ws://1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook wsclient -h

## Use the socks5 proxy

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

* In Chrome, install and configure extension SwitchyOmega by FelisCatus
