## $ brook dns

$ brook dns 用来创建一个加密DNS Server, 它与 $ brook server 一起工作

```
请求 <--> 加密DNS Server <-- | brook server 协议 | --> brook server <--> DNS Server
```

假设你的brook server是 `1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建加密DNS Server `127.0.0.1:53`

## Run brook dns

```
$ brook dns -s 1.2.3.4:9999 -p hello -l 127.0.0.1:53
```

> 你可能需要用sudo运行<br/>
> 更多参数: $ brook dns -h

