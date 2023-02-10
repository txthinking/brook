# brook link protocol

<!--THEME:github-->
<!--G-R3M673HK5V-->

```
brook://KIND?QUERY
```

-   **KIND**: `server`, `wsserver`, `wssserver`, `socks5`, `quicserver`
-   **QUERY**: key=value, key and value should be urlencoded(RFC3986), such as `key0=xxx&key1=xxx`

#### brook server

-   **KIND**: `server`
-   **QUERY**:
    -   `server`: brook server, such as `1.2.3.4:9999`
    -   `password`: password
    -   `udpovertcp`: `true` [Optional]
    -   `name`: such as `my server` [Optional]

Example

```
brook://server?password=hello&server=1.2.3.4%3A9999
```

#### brook wsserver

-   **KIND**: `wsserver`
-   **QUERY**:
    -   `wsserver`: brook wsserver, such as `ws://1.2.3.4:9999`, `ws://1.2.3.4:9999/ws`
    -   `password`: password
    -   `name`: such as `my wsserver` [Optional]
    -   `address`: such as `1.2.3.4:9999` [Optional]
    -   `withoutBrookProtocol`: `true` [Optional]
    -   Any other custom key

Example

```
brook://wsserver?password=hello&wsserver=ws%3A%2F%2F1.2.3.4%3A9999
brook://wsserver?password=hello&wsserver=ws%3A%2F%2F1.2.3.4%3A9999%2Fws
```

#### brook wssserver

-   **KIND**: `wssserver`
-   **QUERY**:
    -   `wssserver`: brook wssserver, such as `wss://domain.com:443`, `wss://domain.com:443/ws`
    -   `password`: password
    -   `name`: such as `my wssserver` [Optional]
    -   `address`: such as `1.2.3.4:9999` [Optional]
    -   `insecure`: `true` [Optional]
    -   `withoutBrookProtocol`: `true` [Optional]
    -   `ca`: CA content [Optional]

Example

```
brook://wssserver?password=hello&wssserver=wss%3A%2F%2Fdomain.com%3A443
brook://wssserver?password=hello&wssserver=wss%3A%2F%2Fdomain.com%3A443%2Fws
```

#### socks5 server

-   **KIND**: `socks5`
-   **QUERY**:
    -   `socks5`: socks5 server, such as `socks5://1.2.3.4:9999`
    -   `username`: username, such as `hello`, [Optional]
    -   `password`: password, such as `world`, [Optional]
    -   `name`: such as `my socks5 server` [Optional]

Example

```
brook://socks5?socks5=socks5%3A%2F%2F1.2.3.4%3A9999
brook://socks5?password=world&socks5=socks5%3A%2F%2F1.2.3.4%3A9999&username=hello
```

#### brook quicserver

-   **KIND**: `quicserver`
-   **QUERY**:
    -   `quicserver`: brook quicserver, such as `quic://domain.com:443`
    -   `password`: password
    -   `name`: such as `my wssserver` [Optional]
    -   `address`: such as `1.2.3.4:9999` [Optional]
    -   `insecure`: `true` [Optional]
    -   `withoutBrookProtocol`: `true` [Optional]
    -   `ca`: CA content [Optional]

Example

```
brook://quicserver?password=hello&quicserver=quic%3A%2F%2Fdomain.com%3A443
```
