## 运行 brook server

假设选择端口`9999`, 密码`hello`

```
$ brook server -l :9999 -p hello
```

假设你的服务器IP是 `1.2.3.4`, 那么你的brook server就是: `1.2.3.4:9999`

> 你可以按组合键CTRL+C来停止<br/>
> 更多参数介绍: $ brook server -h

## 后台运行

```
$ nohup brook server -l :9999 -p hello &
```

## 停止后台运行的 brook

```
$ killall brook
```
