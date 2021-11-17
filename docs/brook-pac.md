## brook pac

brook pac can create a pac server or a pac file.

## Run brook pac

Create a pac server with online domain list(suffix match mode)

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList https://txthinking.github.io/bypass/china_domain.txt
```

Create a pac server with local domain list

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ./china_domain.txt
```

Create a pac file with online domain list

```
brook pac --file ./proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList https://txthinking.github.io/bypass/china_domain.txt
```

Create a pac file with local domain list

```
brook pac --file ./proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ./china_domain.txt
```

> More parameters: brook pac -h

