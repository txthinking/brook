### Brook Link

#### brook server & brook wsserver & brook wssserver

```
brook://urlencode(SERVER PASSWORD)
```

> urlencode() is a 虚拟 RFC3986 函数

SERVER 格式可以是:

* brook server: `server_ip:port`
* brook wsserver: `ws://wsserver_ip:port` 或 `ws://wsserver_ip:port/path`
* brook wssserver: `wss://wssserver_domain:port` 或 `wss://wssserver_domain:port/path`

#### socks5 server

```
brook://urlencode(SERVER)
brook://urlencode(SERVER USERNAME PASSWORD)
```

SERVER 格式可以是:

* socks5 server: `socks5://server_ip:port`

### $ brook link/qr

```
$ brook link -s server_address:port -p password
$ brook link -s ws://wsserver_address:port -p password
$ brook link -s wss://wsserver_domain:port -p password
$ brook link -s socks5://server_address:port
$ brook link -s socks5://server_address:port -u username -p password

$ brook qr -s server_address:port -p password
$ brook qr -s ws://wsserver_address:port -p password
$ brook qr -s wss://wsserver_domain:port
$ brook qr -s socks5://server_address:port -u username -p password
```

### 举例

* brook server: `1.2.3.4:9999`
* password: `password`

那么生成的 Brook Link:

```
brook://1.2.3.4%3A9999%20password
```

