## $ brook tproxy

$ brook tproxy can create Transparent Proxy on your linux router with `TPROXY mod`, it must work with $ brook server.

Assume your brook server is `1.2.3.4:9999` and password is `hello`

> Only support IPv4 server, but if your server support IPv6 you can enable later, please see the introduction below

## Run brook tproxy

```
brook tproxy --server 1.2.3.4:9999 --password hello
```

* By default, OpenWrt will automatically issue the IP of the router as gateway for your computers and mobiles
* And configure your computer/mobile's DNS: such as `8.8.8.8`

> More parameters: $ brook tproxy -h

## Run brook tproxy + smart DNS + bypass list

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --dnsForDefault 8.8.8.8:53 --dnsForBypass 223.5.5.5:53 --bypassDomainList https://txthinking.github.io/bypass/chinadomain.txt --bypassCIDR4List https://txthinking.github.io/bypass/chinacidr4.txt --bypassCIDR6List https://txthinking.github.io/bypass/chinacidr6.txt
```

* And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
* And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
* By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

> More parameters: $ brook tproxy -h

## Run brook tproxy + smart DNS + bypass list + enable IPv6

> Need to support IPv6 both local and server

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --dnsForDefault 8.8.8.8:53 --dnsForBypass 223.5.5.5:53 --bypassDomainList https://txthinking.github.io/bypass/chinadomain.txt --bypassCIDR4List https://txthinking.github.io/bypass/chinacidr4.txt --bypassCIDR6List https://txthinking.github.io/bypass/chinacidr6.txt --enableIPv6
```

* And OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
* And OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
* By default, OpenWrt will automatically issue the IP of the router as gateway and DNS for your computers and mobiles

> More parameters: $ brook tproxy -h
