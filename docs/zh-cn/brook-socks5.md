## $ brook socks5

$ brook socks5 运行一个独立的socks5 server, 支持 TCP and UDP, 假设你的服务器IP是 `1.2.3.4`, 你想创建一个 socks5 server `1.2.3.4:1080`. 如果有防火墙, 记得允许此端口的TCP和UDP协议.

## Run brook socks5

```
$ brook socks5 --socks5 1.2.3.4:1080
```

> 更多参数: $ brook socks5 -h

## 使用刚才创建的socks5代理

* 安裝並配置Chrome擴展[Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
