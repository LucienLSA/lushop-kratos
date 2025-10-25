#!/bin/bash

# Goods Service 测试运行脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}开始运行 Goods Service 测试${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 清理之前的覆盖率文件
rm -f coverage*.out coverage*.html

# 运行 Biz 层测试
echo -e "${BLUE}运行 Biz 层测试...${NC}"
go test -v ./internal/biz/... -coverprofile=coverage_biz.out
echo -e "${GREEN}✓ Biz 层测试完成${NC}\n"

# 运行 Service 层测试
echo -e "${BLUE}运行 Service 层测试...${NC}"
go test -v ./internal/service/... -coverprofile=coverage_service.out
echo -e "${GREEN}✓ Service 层测试完成${NC}\n"

# 运行 Data 层测试
echo -e "${BLUE}运行 Data 层测试...${NC}"
go test -v ./internal/data/... -coverprofile=coverage_data.out
echo -e "${GREEN}✓ Data 层测试完成${NC}\n"

# 生成覆盖率报告
echo -e "${BLUE}生成覆盖率报告...${NC}"
go test -coverprofile=coverage_all.out ./...
go tool cover -html=coverage_all.out -o coverage.html

# 生成详细覆盖率报告
echo -e "${BLUE}生成详细覆盖率报告...${NC}"
go test -coverprofile=coverage_detail.out ./internal/...
go tool cover -func=coverage_detail.out

echo -e "\n${BLUE}========================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}覆盖率报告已生成: coverage.html${NC}"
echo -e "${GREEN}使用浏览器打开查看详细覆盖率信息${NC}\n"
