#!/bin/bash

# 检查是否传递文件路径参数
if [ -z "$1" ]; then
    echo "Usage: $0 <file_path>"
    exit 1
fi

FILE_PATH="$1"

# 确保文件存在，不存在则创建
if [ ! -f "$FILE_PATH" ]; then
    echo "文件$FILE_PATH 不存在，创建中..."
    touch "$FILE_PATH"
fi

# 检查 ban_ip_network 集合是否存在
if ! ipset list ban_ip_network &>/dev/null; then
    echo "ban_ip_network不存在"
    exit 1
fi

# 保存 ban_ip 集合的内容到指定文件
echo "保存 ban_ip_network 集合到文件 $FILE_PATH..."
ipset save ban_ip_network > "$FILE_PATH"
echo "保存成功！文件路径: $FILE_PATH"
