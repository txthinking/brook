## brook socks5

brook socks5 运行一个独立的 socks5 server, 支持 TCP and UDP, 假设你的服务器 IP 是 `1.2.3.4`, 你想创建一个 socks5 server `1.2.3.4:1080`. 如果有防火墙, 记得允许此端口的 TCP 和 UDP 协议.

```
SRC --TCP--> brook socks5 --TCP--> DST
SRC --UDP--> brook socks5 --UDP--> DST
```

## Run brook socks5

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

or with username and password

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

> 更多参数: brook socks5 -h

## 使用刚才创建的 socks5 代理

-   安裝並配置 Chrome 擴展[Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
