## $ brook tproxy

$ brook tproxy can create Transparent Proxy on your linux router with TPROXY mod, it must work with $ brook server.

Assume your brook server is `1.2.3.4:9999` and password is `hello`

## Run brook tproxy

> The following steps are for reference only

#### IPv4

```
echo 1 > /proc/sys/net/ipv4/ip_forward

ip rule add fwmark 1 lookup 100
ip route add local 0.0.0.0/0 dev lo table 100

iptables -t mangle -F
iptables -t mangle -X

iptables -t mangle -A PREROUTING -d 0.0.0.0/8 -j RETURN
iptables -t mangle -A PREROUTING -d 10.0.0.0/8 -j RETURN
iptables -t mangle -A PREROUTING -d 127.0.0.0/8 -j RETURN
iptables -t mangle -A PREROUTING -d 169.254.0.0/16 -j RETURN
iptables -t mangle -A PREROUTING -d 172.16.0.0/12 -j RETURN
iptables -t mangle -A PREROUTING -d 192.168.0.0/16 -j RETURN
iptables -t mangle -A PREROUTING -d 224.0.0.0/4 -j RETURN
iptables -t mangle -A PREROUTING -d 240.0.0.0/4 -j RETURN

# IMPORTANT
iptables -t mangle -A PREROUTING -d 1.2.3.4 -j RETURN

iptables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1
iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1
iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
```

#### IPv6

> Your server must support IPv6

```
echo 1 > /proc/sys/net/ipv6/conf/all/forwarding

ip -6 rule add fwmark 1 table 106
ip -6 route add local ::/0 dev lo table 106

ip6tables -t mangle -F
ip6tables -t mangle -X

# This command with print some CIDRs, remember them
ip address | grep -w inet6 | awk '{print $2}'
# replace REPLACE_ME_WITH_CIDR with the CIDR you remembered, one CIDR one time
ip6tables -t mangle -A PREROUTING -d REPLACE_ME_WITH_CIDR -j RETURN

ip6tables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1
ip6tables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
ip6tables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1
ip6tables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
```

#### Run brook

```
brook tproxy -s 1.2.3.4:9999 -p hello -l :1080
```

> More parameters: $ brook tproxy -h

### On your computer

* Set the gateway to your Linux box IP
* Set the DNS server to 8.8.8.8(or any other working DNS server)
