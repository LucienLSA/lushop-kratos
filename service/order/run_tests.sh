#!/bin/bash

echo "========================================"
echo "Order 服务测试"
echo "========================================"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo ""
echo "📦 运行 Biz 层测试..."
go test ./internal/biz/... -v -count=1 -coverprofile=biz_coverage.out
BIZ_RESULT=$?

echo ""
echo "📦 运行 Service 层测试..."
go test ./internal/service/... -v -count=1 -coverprofile=service_coverage.out 2>&1 | grep -v "no test files"
SERVICE_RESULT=$?

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

echo ""
echo "========================================"
echo "生成覆盖率报告"
echo "========================================"

echo "mode: set" > coverage.out
grep -h -v "^mode:" *_coverage.out >> coverage.out 2>/dev/null
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | tail -1

echo ""
echo "========================================"
echo "测试完成！"
echo "========================================"
echo "覆盖率报告已生成: coverage.html"

if [ $BIZ_RESULT -eq 0 ] && [ $SERVICE_RESULT -eq 0 ]; then
    exit 0
else
    exit 1
fi
