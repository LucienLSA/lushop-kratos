#!/bin/bash

echo "========================================"
echo "UserOp 服务测试"
echo "========================================"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 运行测试
echo ""
echo "📦 运行 Biz 层测试..."
go test ./internal/biz/... -v -count=1 -coverprofile=biz_coverage.out
BIZ_RESULT=$?

echo ""
echo "📦 运行 Service 层测试..."
go test ./internal/service/... -v -count=1 -coverprofile=service_coverage.out
SERVICE_RESULT=$?

echo ""
echo "📦 运行 Data 层测试..."
go test ./internal/data/... -v -count=1 -coverprofile=data_coverage.out
DATA_RESULT=$?

echo ""
echo "========================================"
echo "测试结果汇总"
echo "========================================"

if [ $BIZ_RESULT -eq 0 ]; then
    echo -e "${GREEN}✅ Biz 层测试通过${NC}"
else
    echo -e "${RED}❌ Biz 层测试失败${NC}"
fi

if [ $SERVICE_RESULT -eq 0 ]; then
    echo -e "${GREEN}✅ Service 层测试通过${NC}"
else
    echo -e "${RED}❌ Service 层测试失败${NC}"
fi

if [ $DATA_RESULT -eq 0 ]; then
    echo -e "${GREEN}✅ Data 层测试通过${NC}"
else
    echo -e "${RED}❌ Data 层测试失败${NC}"
fi

echo ""
echo "========================================"
echo "生成覆盖率报告"
echo "========================================"

# 合并覆盖率
echo "mode: set" > coverage.out
grep -h -v "^mode:" *_coverage.out >> coverage.out 2>/dev/null

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html

# 显示覆盖率
go tool cover -func=coverage.out | tail -1

echo ""
echo "========================================"
echo "测试完成！"
echo "========================================"
echo "覆盖率报告已生成: coverage.html"
echo "使用浏览器打开查看详细覆盖率信息"

# 返回测试结果
if [ $BIZ_RESULT -eq 0 ] && [ $SERVICE_RESULT -eq 0 ] && [ $DATA_RESULT -eq 0 ]; then
    exit 0
else
    exit 1
fi
