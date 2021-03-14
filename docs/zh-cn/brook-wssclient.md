## $ brook wssclient

假设你的brook wssserver是 `wss://domain.com:443`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook wssserver 协议 | --> brook wssserver <--> 目标
```

## 运行 brook wssclient

```
$ brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook wssclient -h

## 使用刚才创建的socks5代理

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

* In Chrome, install and configure extension SwitchyOmega by FelisCatus
