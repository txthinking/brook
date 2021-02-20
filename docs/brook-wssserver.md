## Run brook wssserver

Make sure your domain name has been successfully resolved, 80 and 443 are open, brook will automatically issue certificate for you, assume your domain is `domain.com`

```
$ brook wssserver --domain domain.com -p hello
```

> More parameters: $ brook wssserver -h

Then your brook wsserver is: `wss://domain.com:443`

## Run in background or daemon

* Reference [Background](brook-server.md)
* Reference [Daemon](joker.md)
