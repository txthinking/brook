## $ brook wsclient

假设你的brook wsserver是 `ws://1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建一个socks5代理 `127.0.0.1:1080`

```
请求 <--> 本地socks5 <-- | brook wsserver 协议 | --> brook wsserver <--> 目标
```

## 运行 brook wsclient

```
$ brook wsclient -s ws://1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```

> 更多参数: $ brook wsclient -h

## 使用刚才创建的socks5代理

> TODO: 请帮助完善此文档

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

There are two ways to cause your browser to use a SOCKS5 proxy: either system-wide, or browser by browser

To use a SOCKS5 proxy system-wide:

* On Windows, do **Settings** > **Network & Internet** > **Proxy** > **Manual proxy setup**
* On macOS, do Apple > **System Preferences** > **Network** > select network > **Advanced** > **Proxies**

To use a SOCKS5 proxy in just one browser:

* In Firefox, go to **Options** / **Preferences** > **General** > **Network Settings** > **Manual proxy configuration**
* In Chrome, install and configure SwitchyOmega by FelisCatus
* In Safari on macOS, choose **Safari** > **Preferences** > **Advanced** > **Proxies**
