#!/bin/bash

# 检查 ipset 集合是否存在，如果不存在则创建
if ! ipset list ban_ip_network &>/dev/null; then
    echo "Creating ipset collection 'ban_ip_network'..."
    ipset create ban_ip_network hash:net
fi

# 提示用户输入 IP 地址段
read -p "Enter the IP network to block (e.g., 192.168.1.0/24): " ip_network

# 验证输入的 IP 地址段格式
if [[ ! $ip_network =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+$ ]]; then
    echo "Invalid IP network format. Exiting."
    exit 1
fi

# 将 IP 地址段添加到 ban_ip_network 集合
ipset add ban_ip_network "$ip_network" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "Added $ip_network to ban_ip_network."
else
    echo "$ip_network is already in the ban_ip_network set or failed to add."
fi

# 配置 iptables 使用 ban_ip_network 集合
if ! iptables -C INPUT -m set --match-set ban_ip_network src -j DROP &>/dev/null; then
    echo "Adding iptables rule for ban_ip_network..."
    iptables -A INPUT -m set --match-set ban_ip_network src -j DROP
else
    echo "iptables rule for ban_ip_network already exists."
fi
