## 运行 brook server

假设选择端口`9999`, 密码`hello`. 如果有防火墙, 记得允许此端口的TCP和UDP协议.

```
$ brook server --listen :9999 --password hello
```

假设你的服务器IP是 `1.2.3.4`, 那么你的brook server就是: `1.2.3.4:9999`

> 你可以按组合键CTRL+C来停止<br/>
> 更多参数介绍: $ brook server -h

## 使用`nohup`后台运行

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用nohup运行

```
$ nohup brook server --listen :9999 --password hello &
```

停止后台运行的 brook

```
$ killall brook
```

## 使用[joker](https://github.com/txthinking/joker)运行守护进程🔥

> 我们建议你先在前台直接运行, 确保一切都正常后, 再使用joker运行

```
$ joker brook server --listen :9999 --password hello
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

