## \$ brook tproxy

> **依赖: ca-certificates openssl-util ca-bundle coreutils-nohup iptables-mod-tproxy**

$ brook tproxy 可以创建透明代理在你的Linux路由器, **端口9999, 1080, 5353将会被使用**. 它与$ brook server, $ brook wsserver, $ brook wssserver 一起工作.

> 只支持 IPv4 server, 如果你的服務端支持 IPv6, 你可以稍後開啟, 請看下面介紹

## Install ipk

1. 下載適合你系統的[ipk](https://github.com/txthinking/brook/releases)文件
2. 上傳並安裝: OpenWrt Web -> System -> Software -> Upload Package...
3. 刷新頁面, 頂部菜單會出現 Brook 按鈕
4. OpenWrt Web -> Brook -> 輸入後點擊 Connect
5. OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353 或 Brook 創建的其他端口
6. OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
7. 默認, OpenWrt 將會下發 router 的 IP 為電腦或手機的網關和 DNS

## Error log files

* `/root/.brook.web.err`
* `/root/.brook.tproxy.err`
