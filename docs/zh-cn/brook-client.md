## $ brook client

假设你的brook server是 `1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook server 协议 | --> brook server <--> 目标
```

## 运行 brook client

```
$ brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook client -h

## 使用刚才创建的socks5代理

> TODO: 请帮助完善此文档

* 手动配置系统代理
