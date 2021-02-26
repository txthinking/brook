## $ brook wssclient

Assume your brook wssserver is `wss://domain.com:443` and password is `hello`, and you want to create a socks5 proxy `127.0.0.1:1080` on local.

```
send request <--> local socks5 <-- | brook wssserver protocol | --> brook wssserver <--> a remote address
```

## Run brook wssclient

```
$ brook wssclient -s wss://domain.com:443 -p hello --socks5 127.0.0.1:1080
```

> More parameters: $ brook wssclient -h

## Use the socks5 proxy

Once brook is listening as a SOCKS5 proxy on `127.0.0.1` port `1080`, you need to configure your browser to use the SOCKS5 proxy.

There are two ways to cause your browser to use a SOCKS5 proxy: either system-wide, or browser by browser

To use a SOCKS5 proxy system-wide:

* On Windows, do **Settings** > **Network & Internet** > **Proxy** > **Manual proxy setup**
* On macOS, do Apple > **System Preferences** > **Network** > select network > **Advanced** > **Proxies**

To use a SOCKS5 proxy in just one browser:

* In Firefox, go to **Options** / **Preferences** > **General** > **Network Settings** > **Manual proxy configuration**
* In Chrome, install and configure SwitchyOmega by FelisCatus
* In Safari on macOS, choose **Safari** > **Preferences** > **Advanced** > **Proxies**
