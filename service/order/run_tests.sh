#!/bin/bash

echo "========================================"
echo "🛒 Order 服务完整测试套件"
echo "========================================"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo ""
echo -e "${BLUE}📦 运行 Data 层测试...${NC}"
echo "----------------------------------------"
go test ./internal/data/... -v -count=1 -coverprofile=data_coverage.out
DATA_RESULT=$?

if [ $DATA_RESULT -eq 0 ]; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo -e "${BLUE}📦 运行 Biz 层测试...${NC}"
echo "----------------------------------------"
go test ./internal/biz/... -v -count=1 -coverprofile=biz_coverage.out
BIZ_RESULT=$?

if [ $BIZ_RESULT -eq 0 ]; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo -e "${BLUE}📦 运行 Service 层测试...${NC}"
echo "----------------------------------------"
go test ./internal/service/... -v -count=1 -coverprofile=service_coverage.out 2>&1 | grep -v "no test files"
SERVICE_RESULT=$?

if [ $SERVICE_RESULT -eq 0 ]; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "========================================"
echo "📊 测试结果汇总"
echo "========================================"

if [ $DATA_RESULT -eq 0 ]; then
    echo -e "${GREEN}✅ Data 层测试通过${NC}"
else
    echo -e "${RED}❌ Data 层测试失败${NC}"
fi

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
echo "📈 生成覆盖率报告"
echo "========================================"

# 合并覆盖率报告
echo "mode: set" > coverage.out
grep -h -v "^mode:" *_coverage.out >> coverage.out 2>/dev/null

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html

# 显示总体覆盖率
echo ""
echo -e "${BLUE}总体覆盖率:${NC}"
go tool cover -func=coverage.out | tail -1

# 显示各层覆盖率
echo ""
echo -e "${BLUE}各层覆盖率详情:${NC}"
if [ -f data_coverage.out ]; then
    echo -n "Data 层: "
    go tool cover -func=data_coverage.out | tail -1 | awk '{print $NF}'
fi
if [ -f biz_coverage.out ]; then
    echo -n "Biz 层:  "
    go tool cover -func=biz_coverage.out | tail -1 | awk '{print $NF}'
fi
if [ -f service_coverage.out ]; then
    echo -n "Service 层: "
    go tool cover -func=service_coverage.out | tail -1 | awk '{print $NF}'
fi

echo ""
echo "========================================"
echo "🎯 测试统计"
echo "========================================"
echo -e "总测试模块: ${BLUE}${TOTAL_TESTS}${NC}"
echo -e "通过: ${GREEN}${PASSED_TESTS}${NC}"
echo -e "失败: ${RED}${FAILED_TESTS}${NC}"

echo ""
echo "========================================"
echo "✨ 测试完成！"
echo "========================================"
echo -e "📄 覆盖率报告: ${YELLOW}coverage.html${NC}"
echo -e "📊 Data 层覆盖率: ${YELLOW}data_coverage.out${NC}"
echo -e "📊 Biz 层覆盖率: ${YELLOW}biz_coverage.out${NC}"
echo -e "📊 Service 层覆盖率: ${YELLOW}service_coverage.out${NC}"

echo ""
if [ $DATA_RESULT -eq 0 ] && [ $BIZ_RESULT -eq 0 ] && [ $SERVICE_RESULT -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}⚠️  部分测试失败，请检查上面的错误信息${NC}"
    exit 1
fi
