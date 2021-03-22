## $ brook tproxy

$ brook tproxy 可以创建透明代理在你的Linux路由器, Linux需要有`TPROXY内核模块`. 它与$ brook server一起工作.

> 只支持IPv4 server, 如果你的服務端支持IPv6, 你可以稍後開啟, 請看下面介紹

## Install ipk

1. 下載適合你系統的[ipk](https://github.com/txthinking/brook/releases)文件
2. 上傳並安裝: OpenWrt Web -> System -> Software -> Upload Package...
3. 刷新頁面, 頂部菜單會出現Brook按鈕
4. OpenWrt Web -> Brook -> 輸入後點擊Connect

* OpenWrt DNS forwardings: OpenWrt Web -> Network -> DHCP and DNS -> General Settings -> DNS forwardings -> 127.0.0.1#5353 或Brook創建的其他端口
* OpenWrt Ignore resolve file: OpenWrt Web -> Network -> DHCP and DNS -> Resolv and Hosts Files -> Ignore resolve file
* 默認, OpenWrt 將會下發router的IP的為電腦或手機的網關和DNS
