#!/bin/bash

# 检查是否存在名为 ban_port 的 ipset 集合
if ! ipset list ban_port &>/dev/null; then
    echo "创建 ipset 集合 ban_port..."
    ipset create ban_port hash:port
else
    echo "ipset 集合 ban_port 已存在."
fi

# 检查是否已经在 iptables 中配置
if ! iptables -L INPUT -v -n | grep -q "match-set ban_port"; then
    echo "将 ban_port 集合配置到 iptables 中..."
    iptables -I INPUT -m set --match-set ban_port src -j DROP
else
    echo "ban_port 已配置到 iptables 中."
fi
