#!/bin/bash

sudo sysctl -w net.ipv4.ip_forward=1
sudo sysctl -w net.ipv6.conf.all.forwarding=1

sudo iptables -t mangle -F
sudo iptables -t mangle -X

sudo iptables -t mangle -A PREROUTING -d 0.0.0.0/8 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 10.0.0.0/8 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 127.0.0.0/8 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 169.254.0.0/16 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 172.16.0.0/12 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 192.168.0.0/16 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 224.0.0.0/4 -j RETURN
sudo iptables -t mangle -A PREROUTING -d 240.0.0.0/4 -j RETURN
#sudo iptables -t mangle -A PREROUTING -d ::/128 -j RETURN
#sudo iptables -t mangle -A PREROUTING -d ::1/128 -j RETURN
#sudo iptables -t mangle -A PREROUTING -d fc00::/7 -j RETURN
#sudo iptables -t mangle -A PREROUTING -d fe80::/10 -j RETURN
#sudo iptables -t mangle -A PREROUTING -d ff00::/8 -j RETURN
sudo iptables -t mangle -A PREROUTING -d BROOK_SERVER_IP -j RETURN

sudo iptables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1
sudo iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
sudo iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1
sudo iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port 1080
