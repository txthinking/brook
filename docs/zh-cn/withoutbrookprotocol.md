## --withoutBrookProtocol

Brook 协议本身使用强加密算法，当 brook wsserver 和 brook wssserver 开启了 --withoutBrookProtocol 意味着：<br/>
**`性能更好，但数据不使用 Brook 协议进行强加密。所以请使用证书加密，并且不建议--withoutBrookProtocol和--insecure一起使用`**

## 第一种场景: 运行 brook wsserver --withoutBrookProtocol 和 [nico](https://github.com/txthinking/nico) 自动签发证书

-   假设你的域名是 `domain.com`, nico 需要 `443` and `80`, `80` 用于签发证书, 密码 `hello`
-   防火墙记得开放 **TCP 80, 443**
-   确保你的域名 `domain.com` 已成功解析到你服务器的 IP

```
brook wsserver --listen 127.0.0.1:9999 --password hello --withoutBrookProtocol
```
```
nico domain.com http://127.0.0.1:9999
```

#### 在客户端如何连接

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

## 第二种场景: 运行 brook wssserver --withoutBrookProtocol 使用指定证书 [你自己拥有的域名]

-   假设你的域名是 `domain.com`, 选择端口 `443`, 密码 `hello`
-   防火墙记得开放 **TCP 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [如何自己签发证书](https://github.com/txthinking/mad)
-   确保你的域名 `domain.com` 已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

#### 在客户端如何连接

**如果你的证书是信任机构签发**

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> 用 GUI 连接: add the brook link

**如果你的证书是你自己签发的**

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol --ca /path/to/ca.pem
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> 用 GUI 连接: add the brook link

## 第三种场景: 运行 brook wssserver --withoutBrookProtocol 使用指定证书 [你自己不拥有的域名]

-   假设域名是 `domain.com`, 选择端口 `443`, 密码 `hello`
-   防火墙记得开放 **TCP 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [如何自己签发证书](https://github.com/txthinking/mad)

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

#### 在客户端如何连接

假设你的服务器的 IP 是 `1.2.3.4`

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --withoutBrookProtocol --ca /path/to/ca.pem
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> 用 GUI 连接: add the brook link

