## $ brook relayoverbrook

$ brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with $ brook server/wsserver/wssserver.

```
send request <--> a local address <-- | brook protocol | --> brook <--> a remote address
```

Assume your brook server is `1.2.3.4:9999` and password is `hello`, and you want to relay `127.0.0.1:5353` to `8.8.8.8:53`

## Run brook relayoverbrook

```
$ brook relayoverbrook --server 1.2.3.4:9999 --password hello --from 127.0.0.1:5353 --to 8.8.8.8:53
```

> More parameters: $ brook relayoverbrook -h

