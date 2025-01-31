#!/bin/bash

# 检查是否存在名为 white_network_set 的 ipset 集合
if ! ipset list white_network_set &>/dev/null; then
    echo "创建 ipset 集合 white_network_set..."
    ipset create white_network_set hash:net
else
    echo "ipset 集合 white_network_set 已存在."
fi

# 检查是否已经在 iptables 中配置
if ! iptables -L INPUT -v -n | grep -q "match-set white_network_set"; then
    echo "将 white_network_set 集合配置到 iptables 中..."
    iptables -I INPUT -m set --match-set white_network_set src -j ACCEPT
else
    echo "white_network_set 已配置到 iptables 中."
fi
