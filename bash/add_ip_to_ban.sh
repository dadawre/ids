#!/bin/bash


# 检查 ban_ip_network 集合是否存在，如果不存在则创建
if ! ipset list ban_ip &>/dev/null; then
    echo "Creating ipset collection 'ban_ip'..."
    ipset create ban_ip hash:ip
fi

# 检查是否提供了 IP 地址
if [ -z "$1" ]; then
    echo "用法: $0 <IP地址>"
    exit 1
fi

IP="$1"
UnbanTimeOut="$2"

# 检查 ban_ip 集合是否存在
if ! ipset list ban_ip &>/dev/null; then
    echo "ipset 集合 ban_ip 不存在，请先运行初始化脚本."
    exit 1
fi

# 插入 IP 到 ban_ip 集合
if ipset test ban_ip "$IP" &>/dev/null; then
    echo "IP $IP 已存在于 ban_ip 集合中."
else
    echo "将 IP $IP 插入到 ban_ip 集合中..."
    ipset add ban_ip "$IP" timeout "$UnbanTimeOut"
    echo "IP $IP 已成功插入."
fi