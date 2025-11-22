.PHONY: help build build-images build-service build-gateway images-list images-clean start stop restart logs status clean test

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

## help: 显示帮助信息
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "  ${YELLOW}%-15s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## build: 构建所有 Docker 镜像（使用 docker-compose）
build:
	@echo "${GREEN}构建 Docker 镜像...${RESET}"
	@docker-compose build

## build-images: 构建所有服务镜像（使用构建脚本）
build-images:
	@./build-images.sh all

## build-service: 构建指定服务镜像
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "${YELLOW}用法: make build-service SERVICE=user${RESET}"; \
		exit 1; \
	fi
	@./build-images.sh $(SERVICE)

## build-gateway: 构建网关镜像
build-gateway:
	@./build-images.sh gateway

## images-list: 列出所有构建的镜像
images-list:
	@./build-images.sh list

## images-clean: 清理未使用的镜像
images-clean:
	@./build-images.sh clean

## start: 启动所有服务
start:
	@echo "${GREEN}启动所有服务...${RESET}"
	@./deploy.sh start

## stop: 停止所有服务
stop:
	@echo "${YELLOW}停止所有服务...${RESET}"
	@./deploy.sh stop

## restart: 重启所有服务
restart:
	@echo "${YELLOW}重启所有服务...${RESET}"
	@./deploy.sh restart

## logs: 查看服务日志
logs:
	@docker-compose logs -f

## status: 查看服务状态
status:
	@./deploy.sh status

## clean: 清理所有容器和数据
clean:
	@echo "${RED}清理所有容器和数据...${RESET}"
	@./deploy.sh clean

## test: 运行测试
test:
	@echo "${GREEN}运行测试...${RESET}"
	@cd test/integration && ./run_integration_tests.sh

## ps: 查看容器状态
ps:
	@docker-compose ps

## down: 停止并删除容器
down:
	@docker-compose down

## up: 启动服务（不重新构建）
up:
	@docker-compose up -d

## pull: 拉取最新镜像
pull:
	@docker-compose pull
