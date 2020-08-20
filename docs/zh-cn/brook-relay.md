## $ brook relay

$ brook relay can relay a address to a remote address. It can relay any tcp and udp server

```
send request <--> relay server <--> a remote address
```

Assume your (any) server is `1.2.3.4:9999`, and you want to relay port `9999` on your relay server to `1.2.3.4:9999`

## Run brook relay

```
$ brook relay -f :9999 -t 1.2.3.4:9999
```

Assume your relay server IP is `5.6.7.8`, then send request to `5.6.7.8:9999` is equal with send request to `1.2.3.4:9999` now

> More parameters: $ brook relay -h

