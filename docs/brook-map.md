## $ brook map

$ brook map can map a local address to a remote address with brook, both TCP and UDP, it works with $ brook server/wsserver/wssserver.

```
send request <--> a local address <-- | brook protocol | --> brook <--> a remote address
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to map `127.0.0.1:5353` to `8.8.8.8:53`

## Run brook map

```
$ brook map --server 1.2.3.4:9999 --password hello --from 127.0.0.1:5353 --to 8.8.8.8:53
```

> More parameters: $ brook map -h

