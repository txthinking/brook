## $ brook tproxy

$ brook tproxy 可以创建透明代理在你的Linux路由器, Linux需要有TPROXY内核模块. 它与$ brook server一起工作.

假设你的brook server是 `1.2.3.4:9999`, 密码是 `hello`

## 运行 brook tproxy

> 以下步骤在官方openwrt x86-64位固件, 使用的brook_linux_amd64文件测试通过

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

iptables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1
iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1
iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
```

#### IPv6

> 你的服务器必须支持IPv6

```
echo 1 > /proc/sys/net/ipv6/conf/all/forwarding

ip -6 rule add fwmark 1 table 106
ip -6 route add local ::/0 dev lo table 106

ip6tables -t mangle -F
ip6tables -t mangle -X

for s in `ip address | grep -w inet6 | awk '{print $2}'`; do ip6tables -t mangle -A PREROUTING -d $s -j RETURN; done

ip6tables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1
ip6tables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
ip6tables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1
ip6tables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
```

#### 运行 brook tproxy

```
brook tproxy --server 1.2.3.4:9999 --password hello --listen :1080
```

> 更多参数: $ brook tproxy -h

### 在你的电脑上

* 你需要设置一个信用好点的DNS, 比如8.8.8.8
