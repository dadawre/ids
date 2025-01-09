#!/bin/bash

# 检查是否存在名为 ban_ip 的 ipset 集合
if ! ipset list ban_ip &>/dev/null; then
    echo "创建 ipset 集合 ban_ip..."
    ipset create ban_ip hash:ip
else
    echo "ipset 集合 ban_ip 已存在."
fi

# 检查是否已经在 iptables 中配置
if ! iptables -L INPUT -v -n | grep -q "match-set ban_ip"; then
    echo "将 ban_ip 集合配置到 iptables 中..."
    iptables -I INPUT -m set --match-set ban_ip src -j DROP
else
    echo "ban_ip 已配置到 iptables 中."
fi
