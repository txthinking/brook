# Brook CLI Documentation
# NAME

Brook - A cross-platform network tool designed for developers

# SYNOPSIS

Brook

```
[--dialWithDNSPrefer]=[value]
[--dialWithDNS]=[value]
[--dialWithIP4]=[value]
[--dialWithIP6]=[value]
[--dialWithNIC]=[value]
[--dialWithSocks5Password]=[value]
[--dialWithSocks5TCPTimeout]=[value]
[--dialWithSocks5UDPTimeout]=[value]
[--dialWithSocks5Username]=[value]
[--dialWithSocks5]=[value]
[--help|-h]
[--log]=[value]
[--pprof]=[value]
[--prometheusPath]=[value]
[--prometheus]=[value]
[--tag]=[value]
[--version|-v]
```

**Usage**:

```
Brook [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--dialWithDNS**="": When a domain name needs to be resolved, use the specified DNS. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required. Note that for client-side commands, this does not affect the client passing the domain address to the server

**--dialWithDNSPrefer**="": This is used with the dialWithDNS parameter. Prefer A record or AAAA record. Value is A or AAAA

**--dialWithIP4**="": When the current machine establishes a network connection to the outside IPv4, both TCP and UDP, it is used to specify the IPv4 used

**--dialWithIP6**="": When the current machine establishes a network connection to the outside IPv6, both TCP and UDP, it is used to specify the IPv6 used

**--dialWithNIC**="": When the current machine establishes a network connection to the outside, both TCP and UDP, it is used to specify the NIC used

**--dialWithSocks5**="": When the current machine establishes a network connection to the outside, both TCP and UDP, with your socks5 proxy, such as 127.0.0.1:1081

**--dialWithSocks5Password**="": If there is

**--dialWithSocks5TCPTimeout**="": time (s) (default: 0)

**--dialWithSocks5UDPTimeout**="": time (s) (default: 60)

**--dialWithSocks5Username**="": If there is

**--help, -h**: show help

**--log**="": Enable log. A valid value is file path or 'console'. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true

**--pprof**="": go http pprof listen addr, such as :6060

**--prometheus**="": prometheus http listen addr, such as :7070. If it is transmitted on the public network, it is recommended to use it with nico

**--prometheusPath**="": prometheus http path, such as /xxx. If it is transmitted on the public network, a hard-to-guess value is recommended

**--tag**="": Tag can be used to the process, will be append into log, such as: 'key1:value1'

**--version, -v**: print the version


# COMMANDS

## server

Run as brook server, both TCP and UDP

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--listen, -l**="": Listen address, like: ':9999'

**--password, -p**="": Server password

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

## client

Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]

**--http**="": Where to listen for HTTP proxy connections

**--password, -p**="": Brook server password

**--server, -s**="": Brook server address, like: 1.2.3.4:9999

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: UDP over TCP

## wsserver

Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--listen, -l**="": Listen address, like: ':80'

**--password, -p**="": Server password

**--path**="": URL path (default: /ws)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wsclient

Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--http**="": Where to listen for HTTP proxy connections

**--password, -p**="": Brook wsserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

**--wsserver, -s**="": Brook wsserver address, like: ws://1.2.3.4:80, if no path then /ws will be used. Do not omit the port under any circumstances

## wssserver

Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--password, -p**="": Server password

**--path**="": URL path (default: /ws)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wssclient

Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--http**="": Where to listen for HTTP proxy connections

**--insecure**: Client do not verify the server's certificate chain and host name

**--password, -p**="": Brook wssserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

**--wssserver, -s**="": Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances

## quicserver

Run as brook quicserver, both TCP and UDP

**--blockCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--blockCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--domainaddress**="": Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--password, -p**="": Server password

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## quicclient

Run as brook quicclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]. (Note that the global dial parameter is ignored now)

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": Specify ca instead of insecure, such as /path/to/ca.pem

**--http**="": Where to listen for HTTP proxy connections

**--insecure**: Client do not verify the server's certificate chain and host name

**--password, -p**="": Brook quicserver password

**--quicserver, -s**="": Brook quicserver address, like: quic://google.com:443. Do not omit the port under any circumstances

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## relayoverbrook

Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--from, -f, -l**="": Listen address: like ':9999'

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## dnsserveroverbrook

Run as dns server over brook, both TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## tproxy

Run as transparent proxy, a router gateway, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--bypassCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--bypassDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassGeoIP**="": Bypass IP by Geo country code, such as US

**--ca**="": When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dnsForBypass**="": DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required (default: 223.5.5.5:53)

**--dnsForDefault**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsListen**="": Start a DNS server, like: ':53'. MUST contain IP, like '192.168.1.1:53', if you expect your gateway to accept requests from clients to other public DNS servers at the same time

**--doNotRunScripts**: This will not change iptables and others if you want to do by yourself

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--link**="": brook link. This will ignore server, password, udpovertcp, address, insecure, withoutBrookProtocol, ca

**--listen, -l**="": Listen address, DO NOT contain IP, just like: ':8888'. No need to operate iptables by default! (default: :8888)

**--password, -p**="": Password

**--redirectDNS**="": It is usually the value of dnsListen. If the client has set custom DNS instead of dnsListen, this parameter can be intercepted and forwarded to dnsListen. Usually you don't need to set this, only if you want to control it instead of being proxied directly as normal UDP data.

**--server, -s**="": brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443

**--tcpTimeout**="": time (s) (default: 0)

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpTimeout**="": time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--webListen**="": Ignore all other parameters, run web UI, like: ':9999'

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## link

Generate brook link

**--address**="": When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver or brook quicserver, specify ca for untrusted cert, such as /path/to/ca.pem

**--insecure**: When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name

**--name**="": Give this server a name

**--password, -p**="": Password

**--server, -s**="": Support brook server, brook wsserver, brook wssserver, socks5 server, brook quicserver. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080, quic://google.com:443

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be: chrome

**--udpovertcp**: When server is brook server, UDP over TCP

**--username, -u**="": Username, when server is socks5 server

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol

## connect

Run as client and connect to brook link, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]

**--http**="": Where to listen for HTTP proxy connections

**--link, -l**="": brook link, you can get it via $ brook link

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## relay

Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]

**--from, -f, -l**="": Listen address: like ':9999'

**--tcpTimeout**="": time (s) (default: 0)

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": time (s) (default: 60)

## dnsserver

Run as standalone dns server

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## dnsclient

Send a dns query

**--dns, -s**="": DNS server, such as 8.8.8.8:53 (default: 8.8.8.8:53)

**--domain, -d**="": Domain

**--short**: Short for A/AAAA

**--type, -t**="": Type, such as A (default: A)

## dohserver

Run as standalone doh server

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--cert**="": The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically

**--certkey**="": The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically

**--disableA**: Disable A query

**--disableAAAA**: Disable AAAA query

**--dns**="": DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required (default: 8.8.8.8:53)

**--domainaddress**="": Such as: domain.com:443, if you want to create a https server. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used

**--listen**="": listen address, if you want to create a http server behind nico

**--path**="": URL path (default: /dns-query)

**--tcpTimeout**="": time (s) (default: 0)

**--udpTimeout**="": time (s) (default: 60)

## dohclient

Send a dns query

**--doh, -s**="": DOH server, the address is required (default: https://dns.google/dns-query?address=8.8.8.8%3A443)

**--domain, -d**="": Domain

**--short**: Short for A/AAAA

**--type, -t**="": Type, such as A (default: A)

## dhcpserver

Run as standalone dhcp server. Note that you need to stop other dhcp servers, if there are.

**--cache**="": Cache file, local absolute file path, default is $HOME/.brook.dhcpserver

**--count**="": IP range from the start, which you want to assign to clients (default: 100)

**--dnsserver**="": The dns server which you want to assign to clients, such as: 192.168.1.1 or 8.8.8.8

**--gateway**="": The router gateway which you want to assign to clients, such as: 192.168.1.1

**--interface**="": Select interface on multi interface device. Linux only

**--netmask**="": Subnet netmask which you want to assign to clients (default: 255.255.255.0)

**--serverip**="": DHCP server IP, the IP of the this machine, you shoud set a static IP to this machine before doing this, such as: 192.168.1.10

**--start**="": Start IP which you want to assign to clients, such as: 192.168.1.100

## socks5

Run as standalone standard socks5 server, both TCP and UDP

**--limitUDP**: The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.

**--listen, -l**="": Socks5 server listen address, like: :1080 or 1.2.3.4:1080

**--password**="": Password, optional

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--username**="": User name, optional

## socks5tohttp

Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]

**--listen, -l**="": HTTP proxy which will be create: like: 127.0.0.1:8010

**--socks5, -s**="": Socks5 server address, like: 127.0.0.1:1080

**--socks5password**="": Socks5 password, optional

**--socks5username**="": Socks5 username, optional

**--tcpTimeout**="": Connection tcp timeout (s) (default: 0)

## pac

Run as PAC server or save PAC to file

**--bypassDomainList, -b**="": One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--file, -f**="": Save PAC to file, this will ignore listen address

**--listen, -l**="": Listen address, like: 127.0.0.1:1980

**--proxy, -p**="": Proxy, like: 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT' (default: SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT)

## testsocks5

Test UDP and TCP of socks5 server

**--dns**="": DNS server for connecting (default: 8.8.8.8:53)

**--domain**="": Domain for query (default: http3.ooo)

**--password, -p**="": Socks5 password

**--socks5, -s**="": Like: 127.0.0.1:1080

**--username, -u**="": Socks5 username

**-a**="": The A record of domain (default: 137.184.237.95)

## testbrook

Test UDP and TCP of brook server/wsserver/wssserver/quicserver. (Note that the global dial parameter is ignored now)

**--dns**="": DNS server for connecting (default: 8.8.8.8:53)

**--domain**="": Domain for query (default: http3.ooo)

**--link, -l**="": brook link. Get it via $ brook link

**--socks5**="": Temporarily listening socks5 (default: 127.0.0.1:11080)

**-a**="": The A record of domain (default: 137.184.237.95)

## echoserver

Echo server, echo UDP and TCP address of routes

**--listen, -l**="": Listen address, like: ':7777'

## echoclient

Connect to echoserver, echo UDP and TCP address of routes

**--server, -s**="": Echo server address, such as 1.2.3.4:7777

**--times**="": Times of interactions (default: 1)

## completion

Generate shell completions

**--file, -f**="": Write to file (default: brook_autocomplete)

## mdpage

Generate markdown page

**--file, -f**="": Write to file, default print to stdout

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## manpage

Generate man.1 page

**--file, -f**="": Write to file, default print to stdout. You should put to /path/to/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos

## help, h

Shows a list of commands or help for one command

<!--SIDEBAR-->
<!--G-R3M673HK5V-->
