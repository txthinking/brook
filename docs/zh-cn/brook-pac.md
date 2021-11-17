## brook pac

brook pac 可以创建一个 pac server 或 pac 文件

## Run brook pac

创建一个 pac server, 使用在线域名列表(规则: 后缀匹配模式)

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList https://txthinking.github.io/bypass/china_domain.txt
```
创建一个 pac server, 使用本地域名列表

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ./china_domain.txt
```

创建一个 pac 文件, 使用在线域名列表

```
brook pac --file ./proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList https://txthinking.github.io/bypass/china_domain.txt
```

创建一个 pac 文件, 使用本地域名列表

```
brook pac --file ./proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ./china_domain.txt
```

> 更多参数: brook pac -h

