## $ brook dns

$ brook dns can create a encrypted DNS server, both TCP and UDP, it works with $ brook server/wsserver/wssserver.

```
send request <--> encrypted DNS server <-- | brook protocol | --> brook <--> DNS server
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to create a encrypted DNS server `127.0.0.1:53`

## Run brook dns

```
$ brook dns --server 1.2.3.4:9999 --password hello --listen 127.0.0.1:53
```

> More parameters: $ brook dns -h

