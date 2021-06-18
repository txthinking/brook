### Brook Link

```
brook://KIND?QUERY
```

- **KIND**: `server`, `wsserver`, `wssserver`, `socks5`
- **QUERY**: key=value, key and value should be urlencoded(RFC3986), such as `key0=xxx&key1=xxx`

#### brook server

- **KIND**: `server`
- **QUERY**:
    - `server`: brook server, such as `1.2.3.4:9999`
    - `password`: password
    - Any other custom key

Example

```
brook://server?password=hello&server=1.2.3.4%3A9999
```

#### brook wsserver

- **KIND**: `wsserver`
- **QUERY**:
    - `wsserver`: brook wsserver, such as `ws://1.2.3.4:9999`, `ws://1.2.3.4:9999/ws`
    - `password`: password
    - Any other custom key

Example

```
brook://wsserver?password=hello&wsserver=ws%3A%2F%2F1.2.3.4%3A9999
brook://wsserver?password=hello&wsserver=ws%3A%2F%2F1.2.3.4%3A9999%2Fws
```

#### brook wssserver

- **KIND**: `wssserver`
- **QUERY**:
    - `wssserver`: brook wssserver, such as `wss://domain.com:443`, `wss://domain.com:443/ws`
    - `password`: password
    - Any other custom key

Example

```
brook://wssserver?password=hello&wssserver=wss%3A%2F%2Fdomain.com%3A443
brook://wssserver?password=hello&wssserver=wss%3A%2F%2Fdomain.com%3A443%2Fws
```

#### socks5 server

- **KIND**: `socks5`
- **QUERY**:
    - `socks5`: socks5 server, such as `socks5://1.2.3.4:9999`
    - `username`: username, such as `hello`, optional
    - `password`: password, such as `world`, optional
    - Any other custom key

Example

```
brook://socks5?socks5=socks5%3A%2F%2F1.2.3.4%3A9999
brook://socks5?password=world&socks5=socks5%3A%2F%2F1.2.3.4%3A9999&username=hello
```

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
