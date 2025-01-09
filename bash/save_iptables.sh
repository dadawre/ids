#!/bin/bash

# 检查是否提供了文件路径
if [ -z "$1" ]; then
    echo "用法: $0 <文件路径>"
    exit 1
fi

FILE_PATH="$1"


# 保存 iptables配置到指定文件
echo "保存 ban_ip 集合到文件 $FILE_PATH"
iptables-save > "$FILE_PATH"
echo "保存成功！文件路径: $FILE_PATH"
