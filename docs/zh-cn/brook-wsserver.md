## 运行 brook wsserver

假设选择端口`9999`, 密码`hello`

```
$ brook wsserver -l :9999 -p hello
```

假设你的服务器IP是 `1.2.3.4`, 那么你的brook wsserver是: `ws://1.2.3.4:9999`

> 更多参数介绍: $ brook wsserver -h

## 后台运行和守护进程

* 参考 [后台运行](brook-server.md)
* 参考 [守护进程](joker.md)
