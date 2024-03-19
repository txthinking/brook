# brook link protocol

```
brook://KIND?QUERY
```

-   **KIND**: `server`, `wsserver`, `wssserver`, `socks5`, `quicserver`
-   **QUERY**: key=value, key and value should be urlencoded(RFC3986), such as `key0=xxx&key1=xxx`

Checkout `brook link --help`

Example

```
brook://server?password=hello&server=1.2.3.4%3A9999
```
