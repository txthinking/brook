## --withoutBrookProtocol

Brook protocol itself uses strong encryption algorithm, when brook wsserver and brook wssserver are enabled --withoutBrookProtocol means:<br/>
**`Better performance, but data is not strongly encrypted using Brook protocol. So please use certificate encryption, and it is not recommended to use --withoutBrookProtocol and --insecure together `**

## Case 1: Run brook wsserver --withoutBrookProtocol + [nico](https://github.com/txthinking/nico) with trusted certificate

-   Assume your domain is `domain.com`, nico default requires port `443` and `80`, `80` for issuing certificates, with password `hello`
-   If there is a firewall, remember to open **TCP on port 80, 443**
-   Make sure your `domain.com` has been resolved to your server IP successfully

```
brook wsserver --listen 127.0.0.1:9999 --password hello --withoutBrookProtocol
```
```
nico domain.com http://127.0.0.1:9999
```

#### How to connect on the client side

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

## Case 2: Run brook wssserver --withoutBrookProtocol with custom certificate with [your own domain]

-   Assume your domain is `domain.com`, with port `443`, with password `hello`
-   If there is a firewall, remember to open **TCP on port 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [How to issue a certificate yourself](https://github.com/txthinking/mad)
-   Make sure your `domain.com` has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

#### How to connect on the client side

**if your certificate is issued by a trusted authority**

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**if your certificate is issued by yourself**

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol --ca /path/to/ca.pem
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

## Case 3: Run brook wssserver --withoutBrookProtocol with custom certificate with [not your own domain]

-   Assume the domain is `domain.com`, with port `443`, with password `hello`
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [How to issue a certificate yourself](https://github.com/txthinking/mad)
-   If there is a firewall, remember to open **TCP on port 443**

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

#### How to connect on the client side

Assume your server IP is `1.2.3.4`

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --withoutBrookProtocol --ca /path/to/ca.pem
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link result

