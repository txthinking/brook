## $ brook client

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook server protocol | --> brook server <--> a remote address
```

## Run brook client

```
$ brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook client -h

## Use the socks5 proxy created by brook client

> TODO: Please help improve the documentation here

* Configure it on your system network settings
* Configure it on your browser
