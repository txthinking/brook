# Examples

List some examples of common scene commands, pay attention to replace the parameters such as IP, port, password, domain name, certificate path, etc. in the example by yourself

## Run brook server

```
brook server --listen :9999 --password hello
```

then

-   server: `1.2.3.4:9999`
-   password: `hello`

or get brook link

```
brook link --server 1.2.3.4:9999 --password hello --name 'my brook server'
```

or get brook link with `--udpovertcp`

```
brook link --server 1.2.3.4:9999 --password hello --udpovertcp --name 'my brook server'
```

## Run brook wsserver

```
brook wsserver --listen :9999 --password hello
```

then

-   server: `ws://1.2.3.4:9999`
-   password: `hello`

or get brook link

```
brook link --server ws://1.2.3.4:9999 --password hello --name 'my brook wsserver'
```

or get brook link with domain, even if that's not your domain

```
brook link --server ws://hello.com:9999 --password hello --address 1.2.3.4:9999 --name 'my brook wsserver'
```

## Run brook wssserver: automatically certificate

> Make sure your domain has been resolved to your server IP successfully. Automatic certificate issuance requires the use of port 80

```
brook wssserver --domainaddress domain.com:443 --password hello
```

then

-   server: `wss://domain.com:443`
-   password: `hello`

or get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

## Run brook wssserver Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

then

-   server: `wss://domain.com:443`
-   password: `hello`

or get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver'
```

## Run brook wssserver issue untrusted certificates yourself, any domain

Install [mad](https://github.com/txthinking/mad)

```
nami install mad
```

Generate root ca

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

Generate domain cert by root ca

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

Run brook

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem
```

get brook link with `--insecure`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --insecure
```

or get brook link with `--ca`

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --address 1.2.3.4:443 --ca /root/ca.pem
```

## withoutBrookProtocol

Better performance, but data is not strongly encrypted using Brook protocol. So please use certificate encryption, and it is not recommended to use --withoutBrookProtocol and --insecure together

## withoutBrookProtocol automatically certificate

> Make sure your domain has been resolved to your server IP successfully. Automatic certificate issuance requires the use of port 80

```
brook wssserver --domainaddress domain.com:443 --password hello --withoutBrookProtocol
```

get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol
```

## withoutBrookProtocol Use a certificate issued by an existing trust authority

> Make sure your domain has been resolved to your server IP successfully

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

get brook link

```
brook link --server wss://domain.com:443 --password hello --name 'my brook wssserver' --withoutBrookProtocol
```

## withoutBrookProtocol issue untrusted certificates yourself, any domain

Install [mad](https://github.com/txthinking/mad)

```
nami install mad
```

Generate root ca

```
mad ca --ca /root/ca.pem --key /root/cakey.pem
```

Generate domain cert by root ca

```
mad cert --ca /root/ca.pem --ca_key /root/cakey.pem --cert /root/cert.pem --key /root/certkey.pem --domain domain.com
```

Run brook wssserver

```
brook wssserver --domainaddress domain.com:443 --password hello --cert /root/cert.pem --certkey /root/certkey.pem --withoutBrookProtocol
```

Get brook link

```
brook link --server wss://domain.com:443 --password hello --withoutBrookProtocol --address 1.2.3.4:443 --ca /root/ca.pem
```

## Run brook socks5, A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4
```

then

-   server: `1.2.3.4:1080`

or get brook link

```
brook link --server socks5://1.2.3.4:1080
```

## Run brook socks5 with username and password. A stand-alone standard socks5 server

```
brook socks5 --listen :1080 --socks5ServerIP 1.2.3.4 --username hello --password world
```

then

-   server: `1.2.3.4:1080`
-   username: `hello`
-   password: `world`

or get brook link

```
brook link --server socks5://1.2.3.4:1080 --username hello --password world
```

## brook relayoverbrook can relay a local address to a remote address over brook, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook relayoverbrook ... --from 127.0.0.1:5353 --to 8.8.8.8:53
```

## brook dnsserveroverbrook can create a encrypted DNS server, both TCP and UDP, it works with brook server wsserver wssserver.

```
brook dnsserveroverbrook ... --listen 127.0.0.1:53
```

## Brook OpenWRT Router: Perfectly supports IPv4/IPv6/TCP/UDP. Native IPv6

https://www.txthinking.com/talks/articles/brook-openwrt-en.article

## Turn macOS into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-macos-gateway-en.article

## Turn Windows into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-windows-gateway-en.article

## Turn Linux into a Gateway with Brook

https://www.txthinking.com/talks/articles/brook-linux-gateway-en.article

## brook relay can relay a address to a remote address. It can relay any tcp and udp server

```
brook relay --from :9999 --to 1.2.3.4:9999
```

## brook socks5tohttp can convert a socks5 to a http proxy

```
brook socks5tohttp --socks5 127.0.0.1:1080 --listen 127.0.0.1:8010
```

## brook pac creates pac server

```
brook pac --listen 127.0.0.1:8080 --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ...
```

## brook pac creates pac file

```
brook pac --file proxy.pac --proxy 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' --bypassDomainList ...
```

## There are countless examples; for more feature suggestions, it's best to look at the commands and parameters in the CLI documentation one by one, and blog, YouTube...
