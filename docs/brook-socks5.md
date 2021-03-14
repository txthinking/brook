## $ brook socks5

$ brook socks5 run a standalone standard socks5 server, both TCP and UDP, assume your server public IP is `1.2.3.4` and you want to run a socks5 server `1.2.3.4:1080`. If there is a firewall, remember to allow TCP and UDP on this port.

## Run brook socks5

```
$ brook socks5 --socks5 1.2.3.4:1080
```

> More parameters: $ brook socks5 -h

## Use the socks5 proxy

Once brook is listening as a SOCKS5 proxy on `1.2.3.4` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

* In Chrome, install and configure extension SwitchyOmega by FelisCatus
