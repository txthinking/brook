#!/bin/bash

sudo ip rule add fwmark 1 lookup 100
sudo ip route add local 0.0.0.0/0 dev lo table 100
#sudo ip route add local 0:0:0:0:0:0:0:0/0 dev lo table 100
