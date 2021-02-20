## Run brook wsserver

Assume with port `9999` and password `hello`

```
$ brook wsserver -l :9999 -p hello
```

Assume your server public IP is `1.2.3.4`, then your brook wsserver is: `ws://1.2.3.4:9999`

> More parameters: $ brook wsserver -h

## Run in background or daemon

* Reference [Background](brook-server.md)
* Reference [Daemon](joker.md)
