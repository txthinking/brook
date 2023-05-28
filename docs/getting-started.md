# Getting Started 快速上手

## Server

```
bash <(curl https://bash.ooo/nami.sh)
```

```
nami install brook
```

```
brook server -l :9999 -p hello
```

## GUI Client

| iOS | Android      | Mac    |Windows      |Linux        |OpenWrt      |
| --- | --- | --- | --- | --- | --- |
| [![](https://brook.app/images/appstore.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![](https://brook.app/images/android.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.apk) | [![](https://brook.app/images/mac.png)](https://apps.apple.com/us/app/brook-network-tool/id1216002642) | [![Windows](https://brook.app/images/windows.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.exe) | [![](https://brook.app/images/linux.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.bin) | [![OpenWrt](https://brook.app/images/openwrt.png)](https://github.com/txthinking/brook/releases) |

> Linux: [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)<br/>
> OpenWrt: After installation, you need to refresh the page to see the menu

-   brook server: `1.2.3.4:9999` replace 1.2.3.4 with your server IP
-   password: `hello`

## CLI Client

```
brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```
