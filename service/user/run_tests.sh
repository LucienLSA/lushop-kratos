#!/bin/bash

# User Service 单元测试运行脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}User Service 单元测试${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 1. 运行 Biz 层测试
echo -e "${GREEN}[1/3] 运行 Biz 层测试...${NC}"
go test -v ./internal/biz/... -coverprofile=coverage_biz.out
echo -e "${GREEN}✓ Biz 层测试完成${NC}\n"

# 2. 运行 Service 层测试
echo -e "${GREEN}[2/3] 运行 Service 层测试...${NC}"
go test -v ./internal/service/... -coverprofile=coverage_service.out
echo -e "${GREEN}✓ Service 层测试完成${NC}\n"

# 3. 运行 Data 层测试
echo -e "${GREEN}[3/3] 运行 Data 层测试...${NC}"
go test -v ./internal/data/... -coverprofile=coverage_data.out
echo -e "${GREEN}✓ Data 层测试完成${NC}\n"

# 4. 生成总体覆盖率报告
echo -e "${YELLOW}生成覆盖率报告...${NC}"
go test -cover ./...

# 5. 生成详细覆盖率报告
echo -e "${YELLOW}生成详细覆盖率报告...${NC}"
go test -coverprofile=coverage_all.out ./...
go tool cover -html=coverage_all.out -o coverage.html

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "覆盖率报告已生成: ${BLUE}coverage.html${NC}"
echo -e "使用浏览器打开查看详细覆盖率信息"
