## brook dns

brook dns 用来创建一个加密 DNS Server, TCP and UDP, 它与 brook server/wsserver/wssserver 一起工作

```
请求 <--> 加密DNS Server <-- | brook 协议 | --> brook <--> DNS Server
```

假设你的 brook server 是 `1.2.3.4:9999`, 密码是 `hello`, 你要在本地创建加密 DNS Server `127.0.0.1:53`

## Run brook dns

```
brook dns --server 1.2.3.4:9999 --password hello --listen 127.0.0.1:53
```

#### Bypass and Block domain

Check these parameters

-   --dns
-   --dnsForBypass
-   --bypassDomainList
-   --blockDomainList

> 你可能需要用 sudo 运行<br/>
> 更多参数: brook dns -h
