#!/bin/bash

sudo sysctl -w net.ipv4.ip_forward=1
sudo sysctl -w net.ipv6.conf.all.forwarding=1

sudo iptables -t mangle -F
sudo iptables -t mangle -X
sudo iptables -P INPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
sudo iptables -P OUTPUT ACCEPT

sudo iptables -t mangle -N BROOK

sudo iptables -t mangle -A BROOK -d 0.0.0.0/8 -j RETURN
sudo iptables -t mangle -A BROOK -d 10.0.0.0/8 -j RETURN
sudo iptables -t mangle -A BROOK -d 127.0.0.0/8 -j RETURN
sudo iptables -t mangle -A BROOK -d 169.254.0.0/16 -j RETURN
sudo iptables -t mangle -A BROOK -d 172.16.0.0/12 -j RETURN
sudo iptables -t mangle -A BROOK -d 192.168.0.0/16 -j RETURN
sudo iptables -t mangle -A BROOK -d 224.0.0.0/4 -j RETURN
sudo iptables -t mangle -A BROOK -d 240.0.0.0/4 -j RETURN
sudo iptables -t mangle -A BROOK -d BROOK_SERVER_IP -j RETURN

sudo iptables -t mangle -A BROOK -j MARK --set-mark 1
sudo iptables -t mangle -A BROOK -j ACCEPT

sudo iptables -t mangle -A PREROUTING -p tcp -m socket -j BROOK
sudo iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port BROOK_TPROXY_PORT

sudo iptables -t mangle -A PREROUTING -p udp -m socket -j BROOK
sudo iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port BROOK_TPROXY_PORT
