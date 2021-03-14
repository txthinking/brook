## $ brook relay

$ brook relay 可以中继任何TCP和UDP server, 这是一个独立的功能, 它不依赖$ brook server和$ brook wsserver等. 如果有防火墙, 记得允许此端口的TCP和UDP协议.

```
请求 <--> relay server <--> 被中继的server
```

假设你要中继的server地址是 `1.2.3.4:9999`, 你想中继服务器监听端口 `9999` 中继到 `1.2.3.4:9999`

## Run brook relay

```
$ brook relay --from :9999 --to 1.2.3.4:9999
```

假设你的中继服务器IP是 `5.6.7.8`, 那么你就可以访问 `5.6.7.8:9999`等于访问`1.2.3.4:9999`

> 更多参数: $ brook relay -h

