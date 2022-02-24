## brook wssserver

```
SRC --TCP--> brook wssclient/relayoverbrook/dns/tproxy/GUI Client --TCP(TLS(Brook Protocol))--> brook wssserver --TCP--> DST
SRC --UDP--> brook wssclient/relayoverbrook/dns/tproxy/GUI Client --TCP(TLS(Brook Protocol))--> brook wssserver --UDP--> DST
```

## 第一种场景: 运行 brook wssserver 自动签发证书[你自己拥有的域名]

-   假设你的域名是 `domain.com`, 选择端口 `443`, 密码 `hello`
-   防火墙记得开放 **TCP 80, 443**
-   确保你的域名 `domain.com` 已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello
```

> 你可以按组合键 CTRL+C 来停止

#### 在客户端如何连接

-   brook wssserver: `wss://domain.com:443`
-   password: `hello`

> 用 CLI 连接: `brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080`. 更多参数: `brook wssclient -h`<br/>
> 用 GUI 连接: 添加如上信息

**或 获取 brook link**

```
brook link --server wss://domain.com:443 --password hello
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**或 获取 brook link 指定个 `name`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

## 第二种场景: 运行 brook wssserver 使用指定证书 [你自己拥有的域名]

-   假设你的域名是 `domain.com`, 选择端口 `443`, 密码 `hello`
-   防火墙记得开放 **TCP 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [如何自己签发证书](https://github.com/txthinking/mad/blob/master/readme_zh.md)
-   确保你的域名 `domain.com` 已成功解析到你服务器的 IP

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

> 你可以按组合键 CTRL+C 来停止

#### 在客户端如何连接

**如果你的证书是信任机构签发**

-   brook wssserver: `wss://domain.com:443`
-   password: `hello`

> 用 CLI 连接: `brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080`. 更多参数: `brook wssclient -h`<br/>
> 用 GUI 连接: 添加如上信息

**如果你的证书是信任机构签发, 获取 brook link**

```
brook link --server wss://domain.com:443 --password hello
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**如果你的证书是信任机构签发, 获取 brook link 指定个 `name`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**如果你的证书是你自己签发的, 获取 brook link 并指定 `insecure`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --insecure
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

## 第三种场景: 运行 brook wssserver 使用指定证书 [你自己不拥有的域名]

-   假设那个域名是 `domain.com`, 选择端口 `443`, 密码 `hello`
-   防火墙记得开放 **TCP 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [如何自己签发证书](https://github.com/txthinking/mad/blob/master/readme_zh.md)

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

> 你可以按组合键 CTRL+C 来停止

#### 在客户端如何连接

假设你的服务器的 IP 是 `1.2.3.4`

**获取 brook link**

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --insecure
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**或 获取 brook link 指定个 `name`**

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --insecure --name 'my brook wssserver'
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

#### 在服务端屏蔽域名和 IP 列表

查看这些参数

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

> 更多参数: brook wssserver -h

---

## 使用[`joker`](https://github.com/txthinking/joker)运行守护进程 🔥

> 我们建议你先在前台直接运行, 确保一切都正常

```
joker brook wssserver --domainaddress domain.com:443 --password hello
```

查看最后一个命令的 ID

```
joker last
```

查看某个命令的输出和错误

```
joker log <ID>
```

查看运行的命令列表

```
joker list
```

停止某个命令

```
joker stop <ID>
```

---

## 使用[`jinbe`](https://github.com/txthinking/jinbe)开机自动启动命令

> 我们建议你先在前台直接运行, 确保一切都正常

```
jinbe joker brook wssserver --domainaddress domain.com:443 --password hello
```

查看添加的开机命令

```
jinbe list
```

移除某个开机命令

```
jinbe remove <ID>
```
