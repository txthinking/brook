## 运行 brook wsserver

假设选择端口`9999`, 密码`hello`

```
$ brook wsserver -l :9999 -p hello
```

假设你的服务器IP是 `1.2.3.4`, 那么你的brook wsserver是: `ws://1.2.3.4:9999`

> 更多参数介绍: $ brook wsserver -h

## 运行 brook wsserver 和域名

确保你的域名已经解析到你的服务器IP, 并且防火墙已开放80和443端口. brook 会自动签发域名证书. 假设你的域名是 `domain.com`

```
$ brook wsserver --domain domain.com -p hello
```

那么你的 brook wsserver是: `wss://domain.com:443`

## 后台运行和守护进程

* 参考 [后台运行](brook-server.md)
* 参考 [守护进程](joker.md)
