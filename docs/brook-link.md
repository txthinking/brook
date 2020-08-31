### Brook Link

#### brook server & brook wsserver

```
brook://urlencode(SERVER PASSWORD)
```

> urlencode() is a virtual RFC3986 function that means encoding string which in the parentheses

SERVER format:

* brook server: `server_ip:port`
* brook wsserver: `ws://wsserver_ip:port`, `wss://wsserver_domain:port`

#### socks5 server

```
brook://urlencode(SERVER)
brook://urlencode(SERVER USERNAME PASSWORD)
```

SERVER format:

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

### Example

* brook server: `1.2.3.4:9999`
* password: `password`

Link:

```
brook://1.2.3.4%3A9999%20password
```

