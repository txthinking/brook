## $ brook wsclient

假设你的brook wsserver是 `ws://1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook wsserver 协议 | --> brook wsserver <--> 目标
```

## 运行 brook wsclient

```
$ brook wsclient --wsserver ws://1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook wsclient -h

## 使用刚才创建的socks5代理

* 安裝並配置Chrome擴展[Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
