## $ brook dns

$ brook dns can create a encrypted DNS server, it must work with $ brook server.

```
send request <--> encrypted DNS server <-- | brook server protocol | --> brook server <--> DNS server
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to create a encrypted DNS server `127.0.0.1:53`

## Run brook dns

```
$ brook dns -s 1.2.3.4:9999 -p hello -l 127.0.0.1:53
```

> More parameters: $ brook dns -h

