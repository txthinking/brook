## Run brook wsserver

Assume with port `9999` and password `hello`

```
$ brook wsserver -l :9999 -p hello
```

Assume your server public IP is `1.2.3.4`, then your brook wsserver is: `ws://1.2.3.4:9999`

> More parameters: $ brook wsserver -h

## Run brook wsserver with domain

Make sure your domain name has been successfully resolved, 80 and 443 are open, brook will automatically issue certificate for you, assume your domain is `domain.com`

```
$ brook wsserver --domain domain.com -p hello
```

Then your brook wsserver is: `wss://domain.com:443`

## Run in background or daemon

* Reference [Background](brook-server.md)
* Reference [Daemon](joker.md)
