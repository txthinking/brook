## 使用[joker](https://github.com/txthinking/joker)运行brook server守护进程

安装joker, 如果你想看更多信息可以去 [joker github 页面](https://github.com/txthinking/joker)

```
$ nami install github.com/txthinking/joker
```

> 我们建议你先在前台直接运行brook server, 确保一切都正常后再结合joker使用

```
$ joker brook server -l :9999 -p hello
```

> 可以看得出来, 这条命令相比之前的命令只是前面多个joker. 用joker守护某个进程就是这样简单

## 查看joker守护的所有进程

```
$ joker list
```

## 停止joker守护某个进程

> $ joker list 会输出所有进程ID

```
$ joker stop <ID>
```

## 查看某个进程的日志

> $ joker list 会输出所有进程ID

```
$ joker log <ID>
```

