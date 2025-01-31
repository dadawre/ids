#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "使用方法: $0 <要解除白名单的IP>"
    exit 1
fi

IP=$1

# 验证IP格式
if [[ ! $IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "[-] 错误：无效的IP地址格式"
    exit 1
fi

# 检查IP是否存在集合中
if sudo ipset test white_network_set "$IP" &>/dev/null; then
    echo "[+] 正在从 white_network_set 删除 $IP"
    sudo ipset del white_network_set "$IP"
    echo "[√] 已成功删除 $IP"
else
    echo "[!] 注意：$IP 不存在于封禁列表"
fi