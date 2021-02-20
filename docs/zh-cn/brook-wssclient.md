## $ brook wssclient

假设你的brook wssserver是 `wss://domain.com:443`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook wssserver 协议 | --> brook wssserver <--> 目标
```

## 运行 brook wssclient

```
$ brook wssclient -s wss://domain.com:443 -p hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook wssclient -h

## 使用刚才创建的socks5代理

> TODO: 请帮助完善此文档

* 手动配置系统代理
