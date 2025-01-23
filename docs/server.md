# Server

brook dnsserver, dohserver, dnsserveroverbrook, server, wsserver, wssserver, quicserver can use script to do more complex thing. brook will pass different _global variables_ to the script at different times, and the script only needs to assign the processing result to the global variable `out`

## Brook DNS Server

![x](https://brook.app/images/brook-dns-server.svg)

Script can do more:

- There are [examples](https://github.com/txthinking/brook/blob/master/programmable/dnsserver/) for dns server
- In the `script: in_dnsquery` step, script can do more, read more below

## Brook Server

![x](https://brook.app/images/brook-server.svg)

Script can do more:

- There are [examples](https://github.com/txthinking/brook/blob/master/programmable/server/) for server
- In the `script: in_address` step, script can do more, read more below

## Variables

| variable                       | type | command   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_dnsservers                  | map  | dnsserver/dnsserveroverbrook/dohserver/server/wsserver/wssserver/quicserver | When just running | Predefine multiple dns servers, and then programmatically specify which one to use | map      |
| in_dohservers                  | map  | dnsserver/dnsserveroverbrook/dohserver/server/wsserver/wssserver/quicserver | When just running | Predefine multiple doh servers, and then programmatically specify which one to use | map      |
| in_brooklinks                  | map  | server/wsserver/wssserver/quicserver | When just running | Predefine multiple brook links, and then programmatically specify which one to use | map      |
| in_dnsquery                    | map  | dnsserver/dnsserveroverbrook/dohserver | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | server/wsserver/wssserver/quicserver           | When the Server connects the proxied address  | Script can decide how to handle this request                  | map      |

## in_dnsservers

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | dns server | 8.8.8.8:53                           |
| ...    | ... | ... | ... |


## in_dohservers

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | dohserver | https://dns.quad9.net/dns-query?address=9.9.9.9%3A443                           |
| ...    | ... | ... | ... |


## in_brooklinks

| Key    | Type   | Description | Example    |
| ------ | ------ | -------- | ---------- |
| _ | bool | meaningless    | true |

`out`, ignored if not of type `map`

| Key    | Type   | Description | Example    |
| ------------ | ------ | -------------------------------------------------------------------------------------------------- | ------- |
| ...    | ... | ... | ... |
| custom name    | string | brook link | brook://...                           |
| ...    | ... | ... | ... |

## in_dnsquery

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| fromipaddress | string | client address which send this request | 1.2.3.4:5 |
| domain | string | domain name | google.com |
| type   | string | query type  | A          |
| ...   | ... | ...  | ... |
| tag_key   | string | --tag specifies the key value | tag_value |
| ...   | ... | ...  | ... |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`                                                | false   |
| ip           | string | Specify IP directly, only valid when `type` is `A`/`AAAA`                                                                     | 1.2.3.4 |
| dnsserverkey       | string   | Use the dnsserver specified by key to resolve | custom name |
| dohserverkey       | string   | Use the dohserver specified by key to resolve | custom name |

## in_address

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| network | string | `tcp` or `udp` | tcp |
| fromipaddress | string | client address which send this request | 1.2.3.4:5 |
| ipaddress   | string | ip address to be proxied  | 1.2.3.4:443          |
| domainaddress   | string | domain address to be proxied  | google.com:443          |
| user   | string | user ID, only available when used with --userAPI  | 9         |
| ...   | ... | ...  | ... |
| tag_key   | string | --tag specifies the key value | tag_value |
| ...   | ... | ...  | ... |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`                                                | false   |
| address           | string | Rewrite destination to an address                                                                     | 1.2.3.4 |
| ipaddressfromdnsserverkey       | string   | If the destination is domain address, use the dnsserver specified by key to resolve | custom name |
| ipaddressfromdnsserverkey       | string   | If the destination is domain address, use the dohserver specified by key to resolve | custom name |
| aoraaaa       | string   | Must be used with ipaddressfromdnsserverkey or ipaddressfromdnsserverkey. Valid value is `A`/`AAAA` | A |
| speedlimit       | int   | Set a rate limit for this request, for example `1000000` means 1000 kb/s | 1000000 |
| brooklinkkey       | string   | Use the brook link specified by key to proxy | custom name |
| dialwith       | string   | If your server has multiple IPs or network interfaces, you can specify the IP or network interface name to initiate this request | 192.168.1.2 or 2606:4700:3030::ac43:a86a or en1 |
