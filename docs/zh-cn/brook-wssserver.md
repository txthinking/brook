## 运行 brook wssserver

确保你的域名已经解析到你的服务器IP, brook 会自动签发域名证书. 假设你的域名是 `domain.com`. 如果有防火墙, 记得允许80, 443端口的TCP协议.

```
$ brook wssserver --domain domain.com --password hello
```

> 更多参数介绍: $ brook wssserver -h

那么你的 brook wssserver是: `wss://domain.com:443`

## 使用`nohup`后台运行

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用nohup运行

```
$ nohup brook wssserver --domain domain.com --password hello &
```

停止后台运行的 brook

```
$ killall brook
```

## 使用[joker](https://github.com/txthinking/joker)运行守护进程🔥

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用joker运行

```
$ joker brook wssserver --domain domain.com --password hello
```

> 可以看得出来, 这条命令相比之前的命令只是前面多个joker. 用joker守护某个进程就是这样简单

查看joker守护的所有进程

```
$ joker list
```

停止joker守护某个进程

> $ joker list 会输出所有进程ID

```
$ joker stop <ID>
```

查看某个进程的日志

> $ joker list 会输出所有进程ID

```
$ joker log <ID>
```

---

## 使用[boa](https://github.com/brook-community/boa)开机自动启动命令

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用boa运行

```
$ boa brook wssserver --domain domain.com --password hello
```

或者同时用上joker

```
$ boa joker brook wssserver --domain domain.com --password hello
```
