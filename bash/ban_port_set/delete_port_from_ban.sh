#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "使用方法: $0 <要解除封禁的port>"
    exit 1
fi

port=$1

# 验证IP格式
if [[ ! $port =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "[-] 错误：无效的IP地址格式"
    exit 1
fi

# 检查IP是否存在集合中
if sudo ipset test ban_port "$port" &>/dev/null; then
    echo "[+] 正在从 ban_port 删除 $port"
    sudo ipset del ban_port "$port"
    echo "[√] 已成功删除 $port"
else
    echo "[!] 注意：$port 不存在于封禁列表"
fi