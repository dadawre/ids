#!/bin/bash

# 检查是否存在名为 ban_ip 的 ipset 集合
if ! ipset list ban_ip_network &>/dev/null; then
    echo "创建 ipset 集合 ban_ip_network..."
    ipset create ban_ip_network hash:net
else
    echo "ipset 集合 ban_ip_network 已存在."
fi

# 配置 iptables 使用 ban_ip_network 集合
if ! iptables -C INPUT -m set --match-set ban_ip_network src -j DROP &>/dev/null; then
    echo "将 ban_ip_network 集合配置到 iptables 中..."
    iptables -A INPUT -m set --match-set ban_ip_network src -j DROP
else
    echo "ban_ip_network 已配置到 iptables 中."
fi