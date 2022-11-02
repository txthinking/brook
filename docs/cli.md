# NAME

Brook - A cross-platform network tool designed for developers

# SYNOPSIS

Brook

```
[--debug|-d]
[--help|-h]
[--listen|-l]=[value]
[--version|-v]
```

**Usage**:

```
Brook [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--debug, -d**: Enable debug

**--help, -h**: show help

**--listen, -l**="": Listen address for debug (default: :6060)

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

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--toSocks5**="": Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080

**--toSocks5Password**="": Forward to socks5 server, password

**--toSocks5Username**="": Forward to socks5 server, username

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

## client

Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]

**--http**="": Where to listen for HTTP connections

**--password, -p**="": Brook server password

**--server, -s**="": Brook server address, like: 1.2.3.4:9999

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

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

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--toSocks5**="": Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080

**--toSocks5Password**="": Forward to socks5 server, password

**--toSocks5Username**="": Forward to socks5 server, username

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wsclient

Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--http**="": Where to listen for HTTP connections

**--password, -p**="": Brook wsserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

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

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--toSocks5**="": Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080

**--toSocks5Password**="": Forward to socks5 server, password

**--toSocks5Username**="": Forward to socks5 server, username

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--updateListInterval**="": Update list interval, second. default 0, only read one time on start (default: 0)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

## wssclient

Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]

**--address**="": Specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--http**="": Where to listen for HTTP connections

**--insecure**: Client do not verify the server's certificate chain and host name

**--password, -p**="": Brook wssserver password

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--withoutBrookProtocol**: The data will not be encrypted with brook protocol

**--wssserver, -s**="": Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances

## relayoverbrook

Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver <-> to address]

**--address**="": When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--from, -f, -l**="": Listen address: like ':9999'

**--insecure**: When server is brook wssserver, client do not verify the server's certificate chain and host name

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver, the data will not be encrypted with brook protocol

## dnsserveroverbrook

Run as dns server over brook, both TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]

**--address**="": When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--dns**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsForBypass**="": DNS server for resolving domains in bypass list (default: 223.5.5.5:53)

**--insecure**: When server is brook wssserver, client do not verify the server's certificate chain and host name

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver, the data will not be encrypted with brook protocol

## tproxy

Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver <-> dst]

**--address**="": When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--blockDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--bypassCIDR4List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt

**--bypassCIDR6List**="": One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt

**--bypassDomainList**="": One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--dnsForBypass**="": DNS server for resolving domains in bypass list (default: 223.5.5.5:53)

**--dnsForDefault**="": DNS server for resolving domains NOT in list (default: 8.8.8.8:53)

**--dnsListen**="": Start a smart DNS server, like: ':53'

**--doNotRunScripts**: This will not change iptables and others if you want to do by yourself

**--enableIPv6**: Your local and server must support IPv6 both

**--insecure**: When server is brook wssserver, client do not verify the server's certificate chain and host name

**--link**="": brook link. This will ignore server, password, udpovertcp, address, insecure, withoutBrookProtocol, ca

**--listen, -l**="": Listen address, DO NOT contain IP, just like: ':1080'. No need to operate iptables by default! (default: :1080)

**--password, -p**="": Password

**--server, -s**="": brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

**--udpovertcp**: When server is brook server, UDP over TCP

**--webListen**="": Ignore all other parameters, run web UI, like: ':9999'

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver, the data will not be encrypted with brook protocol

## link

Generate brook link

**--address**="": When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443

**--ca**="": When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem

**--insecure**: When server is brook wssserver, client do not verify the server's certificate chain and host name

**--name**="": Give this server a name

**--password, -p**="": Password

**--server, -s**="": Support brook server, brook wsserver, brook wssserver, socks5 server. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080

**--udpovertcp**: When server is brook server, UDP over TCP

**--username, -u**="": Username, when server is socks5 server

**--withoutBrookProtocol**: When server is brook wsserver or brook wssserver, the data will not be encrypted with brook protocol

## connect

Run as client and connect to brook link, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver <-> dst]

**--dialSocks5**="": If you already have a socks5, such as 127.0.0.1:1081, and want [src <-> listen socks5 <-> $ brook connect <-> dialSocks5 <-> $ brook server/wsserver/wssserver <-> dst]

**--dialSocks5Password**="": Optional

**--dialSocks5Username**="": Optional

**--http**="": Where to listen for HTTP connections

**--link, -l**="": brook link, you can get it via $ brook link

**--socks5**="": Where to listen for SOCKS5 connections (default: 127.0.0.1:1080)

**--socks5ServerIP**="": Only if your socks5 server IP is different from listen IP

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

## relay

Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]

**--from, -f, -l**="": Listen address: like ':9999'

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--to, -t**="": Address which relay to, like: 1.2.3.4:9999

**--udpTimeout**="": Connection deadline time (s) (default: 60)

## dnsserver

Run as standalone dns server, both TCP and UDP

**--blockDomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--blockGeoIP**="": Block IP by Geo country code, such as US

**--disableIPv4DomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--disableIPv6DomainList**="": One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt

**--dns**="": DNS server which forward to (default: 8.8.8.8:53)

**--listen, -l**="": Listen address, like: 127.0.0.1:53

**--tcpTimeout**="": Connection deadline time (s) (default: 0)

**--udpTimeout**="": Connection deadline time (s) (default: 60)

## socks5

Run as standalone standard socks5 server, both TCP and UDP

**--limitUDP**: The server MAY use this information to limit access to the UDP association

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

Test UDP and TCP of brook server/wsserver/wssserver

**--dns**="": DNS server for connecting (default: 8.8.8.8:53)

**--domain**="": Domain for query (default: http3.ooo)

**--link, -l**="": brook link. Get it via $ brook link

**--socks5**="": Temporarily listening socks5 (default: 127.0.0.1:11080)

**-a**="": The A record of domain (default: 137.184.237.95)

## completion

Generate shell completions

**--file, -f**="": Write to file (default: brook_autocomplete)

## markdown

Generate markdown page

**--file, -f**="": Write to file, default print to stdout

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## manpage

Generate man.1 page

**--file, -f**="": Write to file, default print to stdout. You should put to /usr/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos

## help, h

Shows a list of commands or help for one command

