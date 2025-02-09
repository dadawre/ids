#!/bin/bash

# 检查 ban_port 集合是否存在，如果不存在则创建
if ! ipset list ban_port &>/dev/null; then
    echo "Creating ipset collection 'ban_port'..."
    ipset create ban_port hash:port
fi

# 检查是否提供了 port 地址
if [ -z "$1" ]; then
    echo "用法: $0 <port地址> [超时时间]"
    exit 1
fi

port="$1"
UnbanTimeOut="$2"

# 检查 ban_port 集合是否存在
if ! ipset list ban_port &>/dev/null; then
    echo "ipset 集合 ban_port 不存在，请先运行初始化脚本."
    exit 1
fi

# 插入 port 到 ban_port 集合
if ipset test ban_port "$port" &>/dev/null; then
    echo "port $port 已存在于 ban_port 集合中."
else
    echo "将 port $port 插入到 ban_port 集合中..."

    # 根据是否存在超时参数执行不同命令
    if [ -n "$UnbanTimeOut" ]; then
        ipset add ban_port "$port" timeout "$UnbanTimeOut"
    else
        ipset add ban_port "$port"
    fi

    echo "port $port 已成功插入."
fi