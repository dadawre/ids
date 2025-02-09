#!/bin/bash

# 检查 ipset 集合是否存在，如果不存在则创建
if ! ipset list ban_ip_network &>/dev/null; then
    echo "创建ban_ip_network集合中..."
    ipset create ban_ip_network hash:net
fi

# 检查是否提供了 IPNetwork 地址
if [ -z "$1" ]; then
    echo "用法: $0 <IP地址>"
    exit 1
fi

IPNetwork="$1"
UnbanTimeOut="$2"
# 检查 ban_ip_network 集合是否存在
if ! ipset list ban_ip_network &>/dev/null; then
    echo "ipset 集合 ban_ip_network 不存在，请先运行初始化脚本."
    exit 1
fi

# 插入 IPNetwork 到 ban_ip_network 集合
if ipset test ban_ip_network "$IPNetwork" &>/dev/null; then
    echo "IPNetwork $IPNetwork 已存在于 ban_ip_network 集合中."
else
    echo "将 IPNetwork $IPNetwork 插入到 ban_ip_network 集合中..."

    # 根据是否存在超时参数执行不同命令
    if [ -n "$UnbanTimeOut" ]; then
        ipset add ban_ip_network "$IPNetwork" timeout "$UnbanTimeOut"
    else
        ipset add ban_ip_network "$IPNetwork"
    fi

    echo "IPNetwork $IPNetwork 已成功插入."
fi
