#!/bin/bash


# 检查 white_set 集合是否存在，如果不存在则创建
if ! ipset list white_set &>/dev/null; then
    echo "Creating ipset collection 'white_set'..."
    ipset create white_set hash:ip
fi

# 检查是否提供了 IP 地址
if [ -z "$1" ]; then
    echo "用法: $0 <IP地址>"
    exit 1
fi

IP="$1"

# 检查 white_set 集合是否存在
if ! ipset list white_set &>/dev/null; then
    echo "ipset 集合 white_set 不存在，请先运行初始化脚本."
    exit 1
fi

# 插入 IP 到 white_set 集合
if ipset test white_set "$IP" &>/dev/null; then
    echo "IP $IP 已存在于 white_set 集合中."
else
    echo "将 IP $IP 插入到 white_set 集合中..."
    ipset add white_set "$IP"
    echo "IP $IP 已成功插入."
fi