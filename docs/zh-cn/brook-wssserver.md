## 运行 brook wssserver

确保你的域名已经解析到你的服务器IP, 并且防火墙已开放80和443端口. brook 会自动签发域名证书. 假设你的域名是 `domain.com`. 如果有防火墙, 记得允许80, 443端口的TCP协议.

```
$ brook wssserver --domain domain.com -p hello
```

> 更多参数介绍: $ brook wssserver -h

那么你的 brook wssserver是: `wss://domain.com:443`

## 后台运行和守护进程

* 参考 [后台运行](brook-server.md)
* 参考 [守护进程](joker.md)
