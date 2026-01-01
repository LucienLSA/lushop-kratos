#!/usr/bin/env bash
# 修复 .env 文件格式的工具脚本

set -e

ENV_FILE=".env"

if [ ! -f "$ENV_FILE" ]; then
    echo "错误: 未找到 .env 文件"
    echo "请先运行: bash deploy/env-template.sh > .env"
    exit 1
fi

echo "修复 .env 文件格式..."

# 创建临时文件
TEMP_FILE=$(mktemp)

# 处理每一行
while IFS= read -r line || [ -n "$line" ]; do
    # 跳过空行
    [[ -z $line ]] && echo "" >> "$TEMP_FILE" && continue

    # 保留注释行
    if [[ $line =~ ^[[:space:]]*# ]]; then
        echo "$line" >> "$TEMP_FILE"
        continue
    fi

    # 处理变量赋值
    if [[ $line =~ ^[[:space:]]*[A-Z_][A-Z0-9_]*= ]]; then
        key=$(echo "$line" | cut -d'=' -f1)
        value=$(echo "$line" | cut -d'=' -f2-)

        # 如果值包含空格或特殊字符，用引号包围
        if [[ $value =~ [[:space:]] || $value =~ [\(\)\[\]\{\}\&\;\|] ]]; then
            # 如果还没有被引号包围
            if [[ ! $value =~ ^\".*\"$ && ! $value =~ ^\'.*\'$ ]]; then
                value="\"$value\""
            fi
        fi

        echo "$key=$value" >> "$TEMP_FILE"
    else
        echo "$line" >> "$TEMP_FILE"
    fi
done < "$ENV_FILE"

# 备份原文件
cp "$ENV_FILE" "${ENV_FILE}.backup.$(date +%Y%m%d_%H%M%S)"

# 替换原文件
mv "$TEMP_FILE" "$ENV_FILE"

echo "✅ .env 文件修复完成"
echo "原文件已备份为: ${ENV_FILE}.backup.*"
echo ""
echo "现在可以运行:"
echo "  ./deploy/deploy-all.sh start"
