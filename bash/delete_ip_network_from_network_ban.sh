#!/bin/bash
set -e

# 参数检查
if [ -z "$1" ]; then
    echo "使用方法: $0 <要解除封禁的IP网段 (CIDR格式)>"
    echo "示例: $0 192.168.1.0/24"
    exit 1
fi

CIDR=$1

# CIDR格式验证
if [[ ! $CIDR =~ ^([0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,2}$ ]]; then
    echo "[-] 错误：无效的CIDR格式，请使用类似 192.168.1.0/24 的格式"
    exit 1
fi

# 检查集合是否存在
if ! ipset list ban_ip_network &>/dev/null; then
    echo "[-] 错误：ipset 集合 ban_ip_network 不存在"
    exit 1
fi

# 检查网段是否存在集合中
if sudo ipset test ban_ip_network "$CIDR" &>/dev/null; then
    echo "[+] 正在从 ban_ip_network 删除 $CIDR"
    sudo ipset del ban_ip_network "$CIDR"
    echo "[√] 已成功删除 $CIDR"
else
    echo "[!] 注意：$CIDR 不存在于封禁列表"
fi