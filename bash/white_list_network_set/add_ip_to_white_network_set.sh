#!/bin/bash


# 检查 white_network_set 集合是否存在，如果不存在则创建
if ! ipset list white_network_set &>/dev/null; then
    echo "Creating ipset collection 'white_network_set'..."
    ipset create white_network_set hash:net
fi

# 检查是否提供了 IP 地址
if [ -z "$1" ]; then
    echo "用法: $0 <IP地址>"
    exit 1
fi

IP="$1"

# 检查 white_network_set 集合是否存在
if ! ipset list white_network_set &>/dev/null; then
    echo "ipset 集合 white_network_set 不存在，请先运行初始化脚本."
    exit 1
fi

# 插入 IP 到 white_network_set 集合
if ipset test white_network_set "$IP" &>/dev/null; then
    echo "IP $IP 已存在于 white_network_set 集合中."
else
    echo "将 IP $IP 插入到 white_network_set 集合中..."
    ipset add white_network_set "$IP"
    echo "IP $IP 已成功插入."
fi