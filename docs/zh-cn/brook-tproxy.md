## brook tproxy

在官方 x86-64 OpenWrt 和 x86-64 Ubuntu 测试通过. 理论上可以运行在任何 Linux 作为网关

> 如果选择 OpenWrt, 推荐使用官方原版 OpenWrt<br/>
> 如果运行在 openwrt 上, 依赖: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy

brook tproxy 可以创建透明代理在你的 Linux 路由器, **Linux 需要有`TPROXY内核模块`**. 它与 brook server, brook wsserver, brook wssserver 一起工作.

假设你的 brook server 是 `1.2.3.4:9999`, 密码是 `hello`

> 只支持 IPv4 server, 如果你的服務端支持 IPv6, 你可以稍後開啟, 請看下面介紹

## 運行 brook tproxy

```
brook tproxy --server 1.2.3.4:9999 --password hello
```

-   OpenWrt:
    -   默認, OpenWrt 將會下發 router 的 IP 的為電腦或手機的網關
    -   配置你的電腦或手機的 DNS, 比如`8.8.8.8`
-   Other Linux:
    -   You may need to manually configure the computer or mobile gateway.
    -   And configure your computer/mobile's DNS: such as `8.8.8.8`

## 運行 brook tproxy + DNS

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353
```

-   OpenWrt:
    -   OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
    -   OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
    -   默認, OpenWrt 將會下發 router 的 IP 的為電腦或手機的網關和 DNS
-   Other Linux:
    -   You may need to change `:5353` to `:53`
    -   You may need to stop system dns server, such as ubuntu `systemctl stop systemd-resolved`
    -   You may need to manually configure the computer or mobile gateway and DNS.

## 運行 brook tproxy + DNS + 開啟 IPv6

> 需要你本地和服務端同時支持 IPv6

```
brook tproxy --server 1.2.3.4:9999 --password hello --dnsListen :5353 --enableIPv6
```

-   OpenWrt:
    -   OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353
    -   OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
    -   默認, OpenWrt 將會下發 router 的 IP 的為電腦或手機的網關和 DNS
-   Other Linux:
    -   You may need to change `:5353` to `:53`
    -   You may need to stop system dns server, such as ubuntu `systemctl stop systemd-resolved`
    -   You may need to manually configure the computer or mobile gateway and DNS.

## Bypass domain, IP(如分流), and block domain(如屏蔽广告)

Check these parameters

-   --dnsForDefault
-   --dnsForBypass
-   --bypassDomainList
-   --bypassCIDR4List
-   --bypassCIDR6List
-   --blockDomainList

> 更多參數: brook tproxy -h
