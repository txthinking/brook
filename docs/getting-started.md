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

| iOS / M1 Mac | Android      | Intel Mac    |Windows      |Linux        |OpenWrt      |
| --- | --- | --- | --- | --- | --- |
| [![](./images/appstore.png)](https://apps.apple.com/us/app/brook-a-cross-platform-proxy/id1216002642) | [![](./images/android.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.apk) | [![](./images/mac.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.dmg) | [![Windows](./images/windows.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.exe) | [![](./images/linux.png)](https://github.com/txthinking/brook/releases/latest/download/Brook.bin) | [![OpenWrt](./images/openwrt.png)](https://github.com/txthinking/brook/releases) |

> Linux: [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)<br/>
> OpenWrt: After installation, you need to refresh the page to see the menu

-   brook server: `1.2.3.4:9999` replace 1.2.3.4 with your server IP
-   password: `hello`

## CLI Client

```
brook client -s 1.2.3.4:9999 -p hello --socks5 127.0.0.1:1080
```
