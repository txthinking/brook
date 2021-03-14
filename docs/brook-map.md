## $ brook map

$ brook map can map a local address to a remote address with brook server protocol, it must work with $ brook server.

```
send request <--> a local address <-- | brook server protocol | --> brook server <--> a remote address
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to map `127.0.0.1:5353` to `8.8.8.8:53`

## Run brook map

```
$ brook map --server 1.2.3.4:9999 --password hello --from 127.0.0.1:5353 --to 8.8.8.8:53
```

> More parameters: $ brook map -h

