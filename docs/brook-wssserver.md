## brook wssserver

```
SRC --TCP--> brook wssclient/relayoverbrook/dns/tproxy/GUI Client --TCP(TLS(Brook Protocol))--> brook wssserver --TCP--> DST
SRC --UDP--> brook wssclient/relayoverbrook/dns/tproxy/GUI Client --TCP(TLS(Brook Protocol))--> brook wssserver --UDP--> DST
```

## Case 1️: Run brook wssserver with automatically certificate with [your own domain]

-   Assume your domain is `domain.com`, with port `443`, with password `hello`
-   If there is a firewall, remember to open **TCP on port 80, 443**
-   Make sure your `domain.com` has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello
```

> You can stop it with CTRL+C

#### How to connect on the client side

-   brook wssserver: `wss://domain.com:443`
-   password: `hello`

> Connect with CLI: `brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080`. More parameters: `brook wssclient -h`<br/>
> Connect with GUI: add info as above

**or get brook link**

```
brook link --server wss://domain.com:443 --password hello
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**or get brook link with `name`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

## Case 2: Run brook wssserver with custom certificate with [your own domain]

-   Assume your domain is `domain.com`, with port `443`, with password `hello`
-   If there is a firewall, remember to open **TCP on port 443**
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [How to issue a certificate yourself](https://github.com/txthinking/mad)
-   Make sure your `domain.com` has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

> You can stop it with CTRL+C

#### How to connect on the client side

**if your certificate is issued by a trusted authority**

-   brook wssserver: `wss://domain.com:443`
-   password: `hello`

> Connect with CLI: `brook wssclient --wssserver wss://domain.com:443 --password hello --socks5 127.0.0.1:1080`. More parameters: `brook wssclient -h`<br/>
> Connect with GUI: add info as above

**if your certificate is issued by a trusted authority, get brook link**

```
brook link --server wss://domain.com:443 --password hello
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**if your certificate is issued by a trusted authority, get brook link with `name`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

**if your certificate is issued by yourself, get brook link with `insecure`**

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --insecure
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link

## Case 3: Run brook wssserver with custom certificate with [not your own domain]

-   Assume the domain is `domain.com`, with port `443`, with password `hello`
-   The cert is `/root/cert.pem`, your cert key is `/root/certkey.pem`. [How to issue a certificate yourself](https://github.com/txthinking/mad)
-   If there is a firewall, remember to open **TCP on port 443**

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

> You can stop it with CTRL+C

#### How to connect on the client side

Assume your server IP is `1.2.3.4`

**get brook link**

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --insecure
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link result

**or get brook link with `name`**

```
brook link --server wss://domain.com:443 --password hello --address 1.2.3.4:443 --insecure --name 'my brook wssserver'
```

> Connect with CLI: `brook connect --link 'brook://...' --socks5 127.0.0.1:1080`. More parameters: `brook connect -h`<br>
> Connect with GUI: add the brook link result

#### Block domain and IP in server-side

Check these parameters

-   --blockDomainList
-   --blockCIDR4List
-   --blockCIDR6List
-   --updateListInterval

> More parameters: brook wssserver -h

---

## Run as daemon via [`joker`](https://github.com/txthinking/joker) 🔥

> We recommend running the command directly to make sure there are no errors

```
joker brook wssserver --domainaddress domain.com:443 --password hello
```

Get the last command ID

```
joker last
```

View output and error of a command run via joker

```
joker log <ID>
```

View running commmands via joker

```
joker list
```

Stop a running command via joker

```
joker stop <ID>
```

---

## Auto start at boot via [`jinbe`](https://github.com/txthinking/jinbe)

> We recommend running the command directly to make sure there are no errors

```
jinbe joker brook wssserver --domainaddress domain.com:443 --password hello
```

View added commmands via jinbe

```
jinbe list
```

Remove a added command via jinbe

```
jinbe remove <ID>
```
