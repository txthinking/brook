## brook wsserver

```
SRC --TCP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --TCP--> DST
SRC --UDP--> brook wsclient/relayoverbrook/dns/tproxy/GUI Client --TCP(Brook Protocol)--> brook wsserver --UDP--> DST
```

## 运行 brook wsserver

-   假设选择端口 `9999`, 密码 `hello`
-   如果有防火墙, 记得允许端口 `9999` 的 **TCP 协议**

```
brook wsserver --listen :9999 --password hello
```

> 你可以按组合键 CTRL+C 来停止

#### 在客户端如何连接

**假设你的服务器 IP 是 `1.2.3.4`**

-   brook wsserver: `ws://1.2.3.4:9999`
-   password: `hello`

> 用 CLI 连接: `brook wsclient --wsserver ws://1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080`. 更多参数: `brook wsclient -h`<br/>
> 用 GUI 连接: 添加如上信息

**或 获取 brook link**

```
brook link --server ws://1.2.3.4:9999 --password hello
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**或 获取 brook link 指定个自定义域名, 任何域名, 甚至不是你自己的域名也可以**

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

**或 获取 brook link 指定个 `name`**

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

> 用 CLI 连接: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. 更多参数: `brook connect -h`<br>
> 用 GUI 连接: 添加 brook link

#### 在服务端屏蔽域名和 IP 列表

查看这些参数

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

> 更多参数: brook wsserver -h

---

## 使用[`joker`](https://github.com/txthinking/joker)运行守护进程 🔥

> 我们建议你先在前台直接运行, 确保一切都正常

```
joker brook wsserver --listen :9999 --password hello
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
jinbe joker brook wsserver --listen :9999 --password hello
```

查看添加的开机命令

```
jinbe list
```

移除某个开机命令

```
jinbe remove <ID>
```
