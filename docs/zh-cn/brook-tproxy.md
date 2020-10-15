## $ brook tproxy

$ brook tproxy 可以在创建透明代理在你的路由器, 它与brook server一起工作.

假设你的brook server是 `1.2.3.4:9999`, 密码是 `hello`

## 运行 brook tproxy

> 因为不同路由器可能有差异, 所以以下仅供参考

#### 1. Route table

```
ip rule add fwmark 1 lookup 100
ip route add local 0.0.0.0/0 dev lo table 100
```

#### 2. Enable forward

```
echo 1 > /proc/sys/net/ipv4/ip_forward
echo 1 > /proc/sys/net/ipv6/conf/all/forwarding
```

#### 3 Load mod

```
modprobe xt_socket
modprobe xt_TPROXY
```

#### 4. iptables

```
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

#### 5. 运行 brook tproxy

```
brook tproxy -s 1.2.3.4:9999 -p hello -l :1080
```

> 更多参数: $ brook tproxy -h

### 在你的电脑上

* 你需要设置一个信用好点的DNS, 比如8.8.8.8
