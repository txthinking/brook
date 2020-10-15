## $ brook wsclient

假设你的brook server是 `ws://1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook wsserver 协议 | --> brook wsserver <--> 目标
```

## 运行 brook wsclient

```
$ brook wsclient -s ws://1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook wsclient -h

## 使用刚才创建的socks5代理

> TODO: 请帮助完善此文档

* 手动配置系统代理
