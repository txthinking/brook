## $ brook tproxy

$ brook tproxy 可以创建透明代理在你的Linux路由器, Linux需要有`TPROXY内核模块`. 它与$ brook server一起工作.

假设你的brook server是 `1.2.3.4:9999`, 密码是 `hello`

> 只支持IPv4 server, 如果你的服務端支持IPv6, 你可以稍後開啟, 請看下面介紹

## 運行 brook tproxy

```
brook tproxy --server 1.2.3.4:9999 --password hello
```

* 默認, OpenWrt 將會下發router的IP的為電腦或手機的網關
* 配置你的電腦或手機的DNS, 比如`8.8.8.8`

> 更多參數: $ brook tproxy -h

## 運行 brook tproxy + 智能DNS + bypass list

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --dnsForDefault 8.8.8.8:53 --dnsForBypass 223.5.5.5:53 --bypassDomainList https://txthinking.github.io/bypass/chinadomain.txt --bypassCIDR4List https://txthinking.github.io/bypass/chinacidr4.txt --bypassCIDR6List https://txthinking.github.io/bypass/chinacidr6.txt
```

* OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
* OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
* 默認, OpenWrt 將會下發router的IP的為電腦或手機的網關和DNS

> 更多參數: $ brook tproxy -h

## 運行 brook tproxy + 智能DNS + bypass list + 開啟IPv6

> 需要你本地和服務端同時支持IPv6

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --dnsForDefault 8.8.8.8:53 --dnsForBypass 223.5.5.5:53 --bypassDomainList https://txthinking.github.io/bypass/chinadomain.txt --bypassCIDR4List https://txthinking.github.io/bypass/chinacidr4.txt --bypassCIDR6List https://txthinking.github.io/bypass/chinacidr6.txt --enableIPv6
```

* OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
* OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
* 默認, OpenWrt 將會下發router的IP的為電腦或手機的網關和DNS

> 更多參數: $ brook tproxy -h
