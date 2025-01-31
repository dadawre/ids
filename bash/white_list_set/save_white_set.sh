#!/bin/bash

# 检查是否提供了文件路径
if [ -z "$1" ]; then
    echo "用法: $0 <文件路径>"
    exit 1
fi

FILE_PATH="$1"

# 检查 white_set 集合是否存在
if ! ipset list white_set &>/dev/null; then
    echo "ipset 集合 white_set 不存在，请先运行初始化脚本."
    exit 1
fi

# 保存 white_set 集合的内容到指定文件
echo "保存 white_set 集合到文件 $FILE_PATH..."
ipset save white_set > "$FILE_PATH"
echo "保存成功！文件路径: $FILE_PATH"
