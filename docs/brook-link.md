### Brook Link

```
brook://urlencode(SERVER PASSWORD)
```

> urlencode() is a virtual RFC3986 function that means encoding string which in the parentheses

SERVER format:

* brook server: `server_ip:port`
* brook wsserver: `ws://wsserver_ip:port`, `wss://wsserver_domain:port`
* socks5: `socks5://server_ip:port`

### Example

* brook server: `1.2.3.4:9999`
* password: `password`

Link/QR:

```
brook://1.2.3.4%3A9999%20password
```

### $ brook link/qr

```
$ brook link -s server_address:port -p password
$ brook link -s ws://wsserver_address:port -p password
$ brook link -s wss://wsserver_domain:port -p password
$ brook link -s socks5://server_address:port -p password

$ brook qr -s server_address:port -p password
$ brook qr -s ws://wsserver_address:port -p password
$ brook qr -s wss://wsserver_domain:port -p password
$ brook qr -s socks5://server_address:port -p password
```
