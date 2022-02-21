## brook tproxy

Test passed on the official x86-64 OpenWrt and x86_64 Ubuntu. Can theoretically run any Linux as a gateway

> If you use OpenWrt, it is recommended to use the official original OpenWrt<br/>
> If running on OpenWrt, dependencies: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy

brook tproxy can create Transparent Proxy on your linux router with `TPROXY mod`, it works with brook server, brook wsserver and brook wssserver.

Assume your brook server is `1.2.3.4:9999` and password is `hello`

> Only support IPv4 server, but if your server support IPv6 you can enable later, please see the introduction below

## Run brook tproxy

```
brook tproxy --server 1.2.3.4:9999 --password hello
```

-   OpenWrt:
    -   By default, OpenWrt will automatically issue the IP of the router as gateway for your computers and mobiles
    -   And configure your computer/mobile's DNS: such as `8.8.8.8`
-   Other Linux:
    -   You may need to manually configure the computer or mobile gateway.
    -   And configure your computer/mobile's DNS: such as `8.8.8.8`

## Run brook tproxy + DNS

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353
```

-   OpenWrt:
    -   And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
    -   And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
    -   By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles
-   Other Linux:
    -   You may need to change `:5353` to `:53`
    -   You may need to stop system dns server, such as ubuntu `systemctl stop systemd-resolved`
    -   You may need to manually configure the computer or mobile gateway and DNS.

## Run brook tproxy + DNS + enable IPv6

> Need to support IPv6 both local and server

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --enableIPv6
```

-   OpenWrt:
    -   And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
    -   And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
    -   By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles
-   Other Linux:
    -   You may need to change `:5353` to `:53`
    -   You may need to stop system dns server, such as ubuntu `systemctl stop systemd-resolved`
    -   You may need to manually configure the computer or mobile gateway and DNS.

## Bypass domain, IP and block domain

Check these parameters

-   --dnsForDefault
-   --dnsForBypass
-   --bypassDomainList
-   --bypassCIDR4List
-   --bypassCIDR6List
-   --blockDomainList

> More parameters: brook tproxy -h
