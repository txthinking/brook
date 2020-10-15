## $ brook map

$ brook map 可以映射任何TCP和UDP server, 让其走brook server协议. 它与 brook server 一起工作.

```
请求 <--> 映射出来的地址 <-- | brook server 协议 | --> brook server <--> 被映射的server
```

假设你的 brook server 是 `1.2.3.4:9999`, 密码是 `hello`, 你想映射 `127.0.0.1:5353` 到 `8.8.8.8:53`

## 运行 brook map

```
$ brook map -s 1.2.3.4:9999 -p hello -f 127.0.0.1:5353 -t 8.8.8.8:53
```

> 更多参数: $ brook map -h

