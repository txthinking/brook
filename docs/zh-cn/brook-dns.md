## $ brook dns

$ brook dns can create a encrypted DNS server, it must work with $ brook server.

```
send request <--> encrypted DNS server <-- | brook server protocol | --> brook server <--> DNS server
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to map `127.0.0.1:5353` to `8.8.8.8:53`

## Run brook map

```
$ brook map -s 1.2.3.4:9999 -p hello -f 127.0.0.1:5353 -t 8.8.8.8:53
```

> More parameters: $ brook map -h

