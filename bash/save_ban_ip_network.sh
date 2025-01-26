#!/bin/bash

# 检查是否传递文件路径参数
if [ -z "$1" ]; then
    echo "Usage: $0 <file_path>"
    exit 1
fi

FILE="$1"

# 确保文件存在，不存在则创建
if [ ! -f "$FILE" ]; then
    echo "文件$FILE不存在，创建中..."
    touch "$FILE"
fi

# 检查 ban_ip_network 集合是否存在
if ! ipset list ban_ip_network &>/dev/null; then
    echo "ban_ip_network不存在"
    exit 1
fi

# 保存 ban_ip_network 集合的 IP 地址段到文件
echo "保存ban_ip_network到$FILE中..."
ipset list ban_ip_network | grep -oP "^\d+\.\d+\.\d+\.\d+/\d+" > "$FILE"

if [ $? -eq 0 ]; then
    echo "IP networks saved successfully to $FILE."
else
    echo "Failed to save IP networks to $FILE."
fi
