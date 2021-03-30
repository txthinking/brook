## 运行 brook wsserver

假设选择端口`9999`, 密码`hello`. 如果有防火墙, 记得允许此端口的 TCP 协议.

```
$ brook wsserver --listen :9999 --password hello
```

假设你的服务器 IP 是 `1.2.3.4`, 那么你的 brook wsserver 是: `ws://1.2.3.4:9999`

> 更多参数介绍: \$ brook wsserver -h

## 使用`nohup`后台运行

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用 nohup 运行

```
$ nohup brook wsserver --listen :9999 --password hello &
```

停止后台运行的 brook

```
$ killall brook
```

## 使用[joker](https://github.com/txthinking/joker)运行守护进程 🔥

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用 joker 运行

```
$ joker brook wsserver --listen :9999 --password hello
```

> 可以看得出来, 这条命令相比之前的命令只是前面多个 joker. 用 joker 守护某个进程就是这样简单

查看 joker 守护的所有进程

```
$ joker list
```

停止 joker 守护某个进程

> \$ joker list 会输出所有进程 ID

```
$ joker stop <ID>
```

查看某个进程的日志

> \$ joker list 会输出所有进程 ID

```
$ joker log <ID>
```

---

## 使用[jinbe](https://github.com/txthinking/jinbe)开机自动启动命令

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用 jinbe 运行

```
$ jinbe brook wsserver --listen :9999 --password hello
```

或者同时用上 joker

```
$ jinbe joker brook wsserver --listen :9999 --password hello
```
