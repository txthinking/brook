# Brook

A cross-platform network tool designed for developers.

[❤️  *A txthinking project*](https://www.txthinking.com)


Table of Contents
=================

* [Documentation](#documentation)
* [Getting Started](#getting-started)
   * [Server](#server)
   * [Client](#client)
* [Brook CLI Documentation](#brook-cli-documentation)
* [NAME](#name)
* [SYNOPSIS](#synopsis)
* [GLOBAL OPTIONS](#global-options)
* [COMMANDS](#commands)
   * [server](#server-1)
   * [client](#client-1)
   * [wsserver](#wsserver)
   * [wsclient](#wsclient)
   * [wssserver](#wssserver)
   * [wssclient](#wssclient)
   * [quicserver](#quicserver)
   * [quicclient](#quicclient)
   * [relayoverbrook](#relayoverbrook)
   * [dnsserveroverbrook](#dnsserveroverbrook)
   * [tproxy](#tproxy)
   * [link](#link)
   * [connect](#connect)
   * [relay](#relay)
   * [dnsserver](#dnsserver)
   * [dnsclient](#dnsclient)
   * [dohserver](#dohserver)
   * [dohclient](#dohclient)
   * [dhcpserver](#dhcpserver)
   * [socks5](#socks5)
   * [socks5tohttp](#socks5tohttp)
   * [pac](#pac)
   * [testsocks5](#testsocks5)
   * [testbrook](#testbrook)
   * [echoserver](#echoserver)
   * [echoclient](#echoclient)
   * [completion](#completion)
   * [mdpage](#mdpage)
      * [help, h](#help-h)
   * [manpage](#manpage)
   * [help, h](#help-h-1)
* [Brook GUI Documentation](#brook-gui-documentation)
   * [Software for which this article applies](#software-for-which-this-article-applies)
   * [Intel Mac GUI proxy mode, Windows GUI proxy mode, Linux GUI proxy mode](#intel-mac-gui-proxy-mode-windows-gui-proxy-mode-linux-gui-proxy-mode)
   * [iOS, M1 Mac GUI, Android GUI, Intel Mac GUI tun mode, Windows GUI tun mode, Linux GUI tun mode](#ios-m1-mac-gui-android-gui-intel-mac-gui-tun-mode-windows-gui-tun-mode-linux-gui-tun-mode)
   * [Configuration Introduction](#configuration-introduction)
   * [Programmable](#programmable)
      * [Introduction to incoming variables](#introduction-to-incoming-variables)
      * [in_guiconfig](#in_guiconfig)
      * [in_dnsquery](#in_dnsquery)
      * [in_address](#in_address)
      * [in_httprequest](#in_httprequest)
      * [in_httpresponse](#in_httpresponse)
      * [How to write Tengo script](#how-to-write-tengo-script)
      * [How to debug script](#how-to-debug-script)
   * [Why and How to Turn Off System and Browser Security DNS](#why-and-how-to-turn-off-system-and-browser-security-dns)
   * [Install CA](#install-ca)
      * [iOS](#ios)
      * [Android](#android)
      * [macOS](#macos)
      * [Windows](#windows)
   * [Apple Push Problem](#apple-push-problem)
* [Brook Diagram](#brook-diagram)
   * [overview](#overview)
   * [withoutBrookProtocol](#withoutbrookprotocol)
   * [relayoverbrook](#relayoverbrook-1)
   * [dnsserveroverbrook](#dnsserveroverbrook-1)
   * [relay](#relay-1)
   * [dnsserver](#dnsserver-1)
   * [tproxy](#tproxy-1)
   * [gui](#gui)
   * [script](#script)

# Documentation

👉 [**Documentation**](https://txthinking.github.io/brook/)
# Getting Started

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

## Client

[GUI Client](https://txthinking.github.io/brook/)

> replace 1.2.3.4 with your server IP

-   brook server: `1.2.3.4:9999`
-   password: `hello`

[CLI Client](https://txthinking.github.io/brook/)

> create socks5://127.0.0.1:1080

`brook client -s 1.2.3.4:9999 -p hello`
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

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be chrome or firefox

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

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be chrome or firefox

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

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be chrome or firefox

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

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be chrome or firefox

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

**--tlsfingerprint**="": When server is brook wssserver, select tls fingerprint, value can be chrome or firefox

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

# Brook GUI Documentation

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## Software for which this article applies

-   [Brook](https://github.com/txthinking/brook)
-   [Brook Plus](https://www.txthinking.com/brook.html)
-   [Shiliew](https://www.txthinking.com/shiliew.html)
-   [tun2brook](https://github.com/txthinking/tun2brook)

## Intel Mac GUI proxy mode, Windows GUI proxy mode, Linux GUI proxy mode

This mode is very simple, will create:

-   Socks5 proxy: `socks5://[::1]:1080` or `socks5://127.0.0.1:1080`
-   HTTP proxy: `http://[::1]:8010` or `http://127.0.0.1:8010`
-   PAC: `http://127.0.0.1:1093/proxy.pac` or `http://[::1]:1093/proxy.pac` based on Bypass Domain list
-   Intel Mac GUI, Windows GUI set PAC to system proxy。Linux GUI can work with [Socks5 Configurator](https://chrome.google.com/webstore/detail/hnpgnjkeaobghpjjhaiemlahikgmnghb)
-   What is socks5 and http proxy? [Article](https://www.txthinking.com/talks/articles/socks5-and-http-proxy-en.article) and [Video](https://www.youtube.com/watch?v=sBCB-X7BoP8)

## iOS, M1 Mac GUI, Android GUI, Intel Mac GUI tun mode, Windows GUI tun mode, Linux GUI tun mode

```
The so-called Internet connection is IP to IP connection, not domain name connection. Therefore, the domain name will be resolved into IP before deciding how to connect.
```

## Configuration Introduction

| Configuration  | Support Systems                   | Conditions                                                                          | Description  |
| -------------- | --------------------------------- | ----------------------------------------------------------------------------------- | --- |
| Import Servers | iOS,Android,Mac,Windows,Linux     | /                                                                                   | brook link list                                                                                                                                                                                                                                                                       |
| System DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System DNS. **Do not bypass this IP**                                                                                                                                                                                                                                                 |
| Fake DNS       | iOS,Android,Mac,Windows,Linux | **Turn off or block the security DNS that comes with the system/browser/etc, see below for details** | The domain name is resolved to Fake IP, which will be converted to a domain name when a connection is initiated, and then the domain name address will be sent to the server, and the server is responsible for domain name resolution                                                |
| Block          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Block switch                                                                                                                                                                                                                                                                          |
| Block Domain   | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Domain name list, matching domain names will be blocked. **Domain name is suffix matching mode**                                                                                                                                                                                      |
| Bypass         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Bypass switch                                                                                                                                                                                                                                                                         |
| Bypass IP      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | CIDR list, matched IP will be bypassed                                                                                                                                                                                                                                                |
| Bypass Geo IP  | iOS,Android,Mac,Windows,Linux                       | /                                                                                   | The matched IP will be bypassed. Note: Global IP changes frequently, so the Geo library is time-sensitive                                                                                                                                                                             |
| Bypass Apps    | Android                           | /                                                                                   | These apps will be bypassed                                                                                                                                                                                                                                                           |
| Bypass DNS     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Support normal DNS, such as `223.5.5.5:53`, support DoH, but need to specify the address of DoH through the parameter address, such as `https://dns. alidns.com/dns-query?address=223.5.5.5%3A443` is used to resolve Bypass Domain. **The IP of this DNS will automatically Bypass** |
| Bypass Domain  | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | List of domain names, matching domain names will use Bypass DNS resolution to get IP, **whether the final connection will be bypassed depends on the Bypass IP** . **The domain name is a suffix matching pattern**                                                                   |
| Hosts          | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Hosts switch                                                                                                                                                                                                                                                                          |
| Hosts List     | iOS,Android,Mac,Windows,Linux     | Fake DNS: On                                                                        | Specify IP, v4, v6 for the domain name, if the value is empty, the effect is the same as Block                                                                                                                                                                                        |
| Programmable   | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Programmable switch                                                                                                                                                                                                                                                                   |
| Script         | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Script. All functions above can be controlled. And more and more, **The whole process can be controled, see below for details**.                                                                                                                                                      |
| Log            | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log switch                                                                                                                                                                                                                                                                            |
| Log View       | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log List                                                                                                                                                                                                                                                                              |
| Log View Plus  | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Log list, easier to read, filter conditions, etc.                                                                                                                                                                                                                                     |
| MITM Log View  | iOS,Android,Mac,Windows,Linux     | /                                                                                   | MITM log list, such as https request response, hexadecimal, JSON, image, etc.                                                                                                                                                                                                         |
| TUN            | iOS,Android,Mac,Windows,Linux                 | /                                                                                   | Choose proxy mode or tun. iOS and Android force TUN mode mode                                                                                                                                                                                                                                                         |
| Capture Me     | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Test your packet capture or proxy software is working as a system proxy or TUN                                                                                                                                                                                                        |
| DNS Client    | iOS,Android,Mac,Windows,Linux | /                                           | DNS client                                                                                   |
| DOH Client    | iOS,Android,Mac,Windows,Linux | /                                           | DOH client                                                                                                           |
| Echo Client    | iOS,Android,Mac,Windows,Linux | /                                           | Echo client                                                                   |
| Test Socks5    | iOS,Android,Mac,Linux | /                                           | Test socks5 server                                                                    |
| Dark Mode      | iOS,Android,Mac,Windows,Linux     | /                                                                                   | System / Light / Dark                                                                                                                                                                                                                                                                 |     |
| Shortcut       | iOS,Android,Mac,Windows,Linux     | /                                                                                   | Quickly control the functions in the menu on the home page                                                                                                                                                                                                                            |
| System Tray    | Windows                           | /                                                                                   | Open as systray, then open dashboard from the systray                                                                                                                                                                                                                                 |

## Programmable

```
Brook GUI will pass different global variables to the script at different times, and the script only needs to assign the processing result to the global variable out
```

Take full control of your own network

-   Like turning off IPv6 by blocking AAAA
-   Block system/browser built-in secure DNS
-   Override DST
-   Flexible and finer rules
-   Directly bypass the domain name regardless of whether the resolved IP is in Bypass
-   MITM decrypt HTTPS
-   Packet capture
-   Packet modify
-   Disable HTTP3
-   more and more...

### Introduction to incoming variables

| variable                       | type | condition   | timing                            | description                                       | out type |
| ------------------------------ | ---- | ----------- | --------------------------------- | ------------------------------------------------- | -------- |
| in_guiconfig                   | map  | /           | before connected                  | to override GUI configuration                     | map      |
| in_dnsquery                    | map  | FakeDNS: On | When a DNS query occurs           | Script can decide how to handle this request      | map      |
| in_address                     | map  | /           | When connecting to an address     | script can decide how to connect                  | map      |
| in_httprequest                 | map  | /           | When an HTTP(S) request comes in  | the script can decide how to handle this request  | map      |
| in_httprequest,in_httpresponse | map  | /           | when an HTTP(S) response comes in | the script can decide how to handle this response | map      |

### in_guiconfig

| Key | Type | Description                                       |
| --- | ---- | ------------------------------------------------- |
| \_  | bool | For future compatibility, this key can be ignored |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`, if it is `map` then explicitly specify each configuration item.

| Key        | Type   | Description       |
| ---------- | ------ | ----------------- |
| systemdns4 | string | System DNS v4     |
| systemdns6 | string | System DNS v6     |
| fakedns    | bool   | Fake DNS switch   |
| block      | bool   | GUI Block switch  |
| bypass     | bool   | GUI Bypass switch |
| bypassdns4 | string | Bypass DNS v4     |
| bypassdns6 | string | Bypass DNS v6     |
| hosts      | bool   | GUI Hosts switch  |

### in_dnsquery

| Key    | Type   | Description | Example    |
| ------ | ------ | ----------- | ---------- |
| domain | string | domain name | google.com |
| type   | string | query type  | A          |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key          | Type   | Description                                                                                                                   | Example |
| ------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- | ------- |
| block        | bool   | Whether Block, default `false`. It is an OR relationship with GUI Block Domain                                                | false   |
| ip           | string | Specify IP directly, only valid when `type` is `A`/`AAAA`                                                                     | 1.2.3.4 |
| forcefakedns | bool   | Ignore GUI Bypass Domain, handle with Fake DNS, only valid when `type` is `A`/`AAAA`, default `false`                         | false   |
| system       | bool   | Get IP from system DNS, default `false`                                                                                       | false   |
| bypass       | bool   | whether to Bypass, default `false`, if `true` then use bypass DNS to resolve. It is an OR relationship with GUI Bypass Domain | false   |

### in_address

| Key           | Type   | Description                                                                                                         | Example        |
| ------------- | ------ | ------------------------------------------------------------------------------------------------------------------- | -------------- |
| network       | string | Network type, the value `tcp`/`udp`                                                                                 | tcp            |
| ipaddress     | string | IP type address. There is only of ipaddress and domainaddress. Note that there is no relationship between these two | 1.2.3.4:443    |
| domainaddress | string | Domain type address, because of FakeDNS we can get the domain name address here                                     | google.com:443 |

`out`, if it is `error` type will be recorded in the log. Ignored if not of type `map`

| Key                    | Type   | Description                                                                                                                                                                                             | Example     |
| ---------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| block                  | bool   | Whether Block, default `false`                                                                                                                                                                          | false       |
| ipaddress              | string | IP type address, rewrite destination                                                                                                                                                                    | 1.2.3.4:443 |
| ipaddressfrombypassdns | string | Use Bypass DNS to obtain `A` or `AAAA` IP and rewrite the destination, only valid when `domainaddress` exists, the value `A`/`AAAA`                                                                     | A           |
| bypass                 | bool   | Bypass, default `false`. If `true` and `domainaddress`, then `ipaddress` or `ipaddressfrombypassdns` must be specified. It is an OR relationship with GUI Bypass IP | false       |
| mitm                   | bool   | Whether to perform MITM, default `false`. Only valid when `network` is `tcp`. Need to install CA, see below                                                                                             | false       |
| mitmprotocol           | string | MITM protocol needs to be specified explicitly, the value is `http`/`https`                                                                                                                             | https       |
| mitmcertdomain         | string | The MITM certificate domain name, which is taken from `domainaddress` by default. If `ipaddress` and `mitm` is `true` and `mitmprotocol` is `https` then must be must be specified explicitly           | example.com |
| mitmwithbody           | bool   | Whether to manipulate the http body, default `false`. **will read the body of the request and response into the memory and interact with the script. iOS 50M total memory limit may kill process**      | false       |
| mitmautohandlecompress | bool   | Whether to automatically decompress the http body when interacting with the script, default `false`                                                                                                     | false       |
| mitmclienttimeout      | int    | Timeout for MITM talk to server, second, default 0                                                                                                                                                      | 0           |
| mitmserverreadtimeout  | int    | Timeout for MITM read from client, second, default 0                                                                                                                                                    | 0           |
| mitmserverwritetimeout | int    | Timeout for MITM write to client, second, default 0                                                                                                                                                     | 0           |

### in_httprequest

| Key    | Type   | Description                   | Example                     |
| ------ | ------ | ----------------------------- | --------------------------- |
| URL    | string | URL                           | `https://example.com/hello` |
| Method | string | HTTP method                   | GET                         |
| Body   | bytes  | HTTP request body             | /                           |
| ...    | string | other fields are HTTP headers | /                           |

`out`, must be set to a request or response

### in_httpresponse

| Key        | Type   | Description                   | Example |
| ---------- | ------ | ----------------------------- | ------- |
| StatusCode | int    | HTTP status code              | 200     |
| Body       | bytes  | HTTP response body            | /       |
| ...        | string | other fields are HTTP headers | /       |

`out`, must be set to a response

### How to write Tengo script

[Tengo Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)

Library

-   [text](https://github.com/d5/tengo/blob/master/docs/stdlib-text.md): regular expressions, string conversion, and manipulation
-   [math](https://github.com/d5/tengo/blob/master/docs/stdlib-math.md): mathematical constants and functions
-   [times](https://github.com/d5/tengo/blob/master/docs/stdlib-times.md): time-related functions
-   [rand](https://github.com/d5/tengo/blob/master/docs/stdlib-rand.md): random functions
-   [fmt](https://github.com/d5/tengo/blob/master/docs/stdlib-fmt.md): formatting functions
-   [json](https://github.com/d5/tengo/blob/master/docs/stdlib-json.md): JSON functions
-   [enum](https://github.com/d5/tengo/blob/master/docs/stdlib-enum.md): Enumeration functions
-   [hex](https://github.com/d5/tengo/blob/master/docs/stdlib-hex.md): hex encoding and decoding functions
-   [base64](https://github.com/d5/tengo/blob/master/docs/stdlib-base64.md): base64 encoding and decoding functions
-   `brook`: brook module

    ```
    Constants

    * os: string, linux/darwin/windows/ios/android. If ios app run on mac, it is ios
    * iosapponmac: bool, ios app run on mac

    Functions

    * splithostport(address string) => map/error: splits a network address of the form "host:port" to { "host": "xxx", "port": "xxx" }
    * country(ip string) => string/error: get country code from ip
    * cidrcontainsip(cidr string, ip string) => bool/error: reports whether the network includes ip
    * parseurl(url string) => map/error: parses a raw url into a map, keys: scheme/host/path/rawpath/rawquery
    * parsequery(query string) => map/error: parses a raw query into a kv map
    * map2query(kv map) => string/error: convert map{string:string} into a query string
    * bytes2ints(b bytes) => array/error: convert bytes into [int]
    * ints2bytes(ints array) => bytes/error: convert [int] into bytes
    * bytescompare(a bytes, b bytes) => int/error: returns an integer comparing two bytes lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b
    * bytescontains(b bytes, sub bytes) => bool/error: reports whether sub is within b
    * byteshasprefix(s bytes, prefix bytes) => bool/error: tests whether the bytes s begins with prefix
    * byteshassuffix(s bytes, suffix bytes) => bool/error: tests whether the bytes s ends with suffix
    * bytesindex(s bytes, sep bytes) => int/error: returns the index of the first instance of sep in s, or -1 if sep is not present in s
    * byteslastindex(s bytes, sep bytes) => int/error: returns the index of the last instance of sep in s, or -1 if sep is not present in s
    * bytesreplace(s bytes, old bytes, new bytes, n int) => bytes/error: returns a copy of the s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements
    * pathescape(s string) => string/error: escapes the string so it can be safely placed inside a URL path segment, replacing special characters (including /) with %XX sequences as needed
    * pathunescape(s string) => string/error: does the inverse transformation of pathescape
    * queryescape(s string) => string/error: escapes the string so it can be safely placed inside a URL query
    * queryunescape(s string) => string/error: does the inverse transformation of queryescape
    * hexdecode(s string) => bytes/error: returns the bytes represented by the hexadecimal string s
    * hexencode(s string) => string/error: returns the hexadecimal encoding of src
    ```

Example

https://github.com/txthinking/bypass/blob/master/example_script.tengo

### How to debug script

-   It is recommended to use [tun2brook](https://github.com/txthinking/tun2brook) on desktop to debug with print
-   It is recommended to use [mitmproxy helper](https://www.txthinking.com/mitmproxy.html) and [Wireshark Helper](https://www.txthinking.com/wireshark.html) to capture packets to determine what to modify

## Why and How to Turn Off System and Browser Security DNS

Because if Security DNS is turned on, the Fake DNS will not work. So we have to turn it off:

-   Android: Settings -> Network & internet -> Private DNS -> Off
-   Chrome on Mobile: Settings -> Privacy and security -> Use secure DNS -> Off
-   Chrome on Desktop: Settings -> Privacy and security -> Security -> Use secure DNS -> Off
-   Windows: Windows Settings -> Network & Internet -> Your Network -> DNS settings -> Edit -> Preferred DNS -> Unencrypted only -> 8.8.8.8
-   iOS/Mac avoid upgrading to secure DNS: related DST can be blocked by script. You can also create a DNS by yourself: `brook dnsserver --listen :53`

Other systems and software, please find out whether it exists and how to close it

## Install CA

https://txthinking.github.io/ca/ca.pem

### iOS

https://www.youtube.com/watch?v=HSGPC2vpDGk

### Android

Android has user CA and system CA, must be installed in the system CA after ROOT

### macOS

```
nami install mad ca.txthinking
sudo mad install --ca ~/.nami/bin/ca.pem
```

### Windows

Open GitBash

```
nami install mad ca.txthinking
```

Then open GitBash with admin

```
mad install --ca ~/.nami/bin/ca.pem
```

Note that software such as GitBash or Firefox may not read the system CA, you can use the system Edge browser to test after installation

## Apple Push Problem

To receive push, Apple Server only allows Ethernet, cellular data, Wi-Fi connections. So you need to Bypass the relevant domain name and IP. [Reference link](https://support.apple.com/en-us/HT210060)
# Brook Diagram

<!--SIDEBAR-->
<!--G-R3M673HK5V-->

## overview

![overview](https://txthinking.github.io/brook/svg/overview.svg)

## withoutBrookProtocol

![wbp](https://txthinking.github.io/brook/svg/wbp.svg)

## relayoverbrook

![relayoverbrook](https://txthinking.github.io/brook/svg/relayoverbrook.svg)

## dnsserveroverbrook

![dnsserveroverbrook](https://txthinking.github.io/brook/svg/dnsserveroverbrook.svg)

## relay

![relay](https://txthinking.github.io/brook/svg/relay.svg)

## dnsserver

![dnsserver](https://txthinking.github.io/brook/svg/dnsserver.svg)

## tproxy

![tproxy](https://txthinking.github.io/brook/svg/tproxy.svg)

## gui

![gui](https://txthinking.github.io/brook/svg/gui.svg)

## script

![script](https://txthinking.github.io/brook/svg/script.svg)

