#!/usr/bin/env bash
set -euo pipefail

# 一键部署所有服务脚本
# 使用方法: ./deploy-all.sh [start|stop|restart|status]

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务列表
SERVICES=(
    "mysql:mysqld"
    "redis:redis"
    "consul:consul"
    "nacos:nacos"
    "elasticsearch:elasticsearch"
    "kibana:kibana"
    "grafana:grafana"
    "prometheus:prometheus"
    "jaeger:jaeger"
)

# RocketMQ 特殊处理
ROCKETMQ_SERVICES=("rocketmq-namesrv" "rocketmq-broker" "rocketmq-proxy" "rocketmq-dashboard")

# 加载环境变量
if [ -f ".env" ]; then
    echo -e "${BLUE}加载环境变量文件 .env${NC}"
    # 使用更安全的方式加载环境变量
    while IFS='=' read -r key value; do
        # 跳过注释行和空行
        [[ $key =~ ^[[:space:]]*# ]] && continue
        [[ -z $key ]] && continue

        # 移除变量名后的空格，移除值后的注释
        key=$(echo "$key" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        value=$(echo "$value" | sed 's/[[:space:]]*#.*$//;s/^[[:space:]]*//;s/[[:space:]]*$//')

        # 如果值被引号包围，移除引号
        if [[ $value =~ ^\".*\"$ ]]; then
            value=$(echo "$value" | sed 's/^"//;s/"$//')
        elif [[ $value =~ ^\'.*\'$ ]]; then
            value=$(echo "$value" | sed "s/^'//;s/'$//")
        fi

        # 导出变量
        export "$key=$value"
    done < .env
fi

# 工具函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装或不在 PATH 中"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行"
        exit 1
    fi
}

wait_for_service() {
    local service=$1
    local max_attempts=30
    local attempt=1

    log_info "等待服务 $service 启动..."

    while [ $attempt -le $max_attempts ]; do
        if docker ps --filter "name=$service" --filter "status=running" | grep -q "$service"; then
            log_info "服务 $service 已启动"
            return 0
        fi

        echo -n "."
        sleep 2
        ((attempt++))
    done

    log_error "服务 $service 启动超时"
    return 1
}

start_mysql() {
    log_info "启动 MySQL..."
    bash "$SCRIPT_DIR/mysql/mysql.sh"

    # 等待 MySQL 就绪
    wait_for_service "lushop-mysql"

    # 初始化数据库（如果存在初始化脚本）
    if [ -f "$PROJECT_ROOT/scripts/init_db.sql" ]; then
        log_info "初始化数据库..."
        sleep 5  # 等待 MySQL 完全就绪

        # 使用容器执行初始化
        docker exec -i lushop-mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < "$PROJECT_ROOT/scripts/init_db.sql"
        log_info "数据库初始化完成"
    fi
}

start_service() {
    local service=$1
    local script_path="$SCRIPT_DIR/$service/$service.sh"

    if [ -f "$script_path" ]; then
        log_info "启动 $service..."
        bash "$script_path"
        wait_for_service "lushop-$service"
    else
        log_error "脚本不存在: $script_path"
        return 1
    fi
}

start_rocketmq() {
    log_info "启动 RocketMQ..."

    # 确保数据目录存在
    mkdir -p "${ROCKETMQ_DATA_DIR:-$SCRIPT_DIR/rocketmq/data}"

    # 生成配置文件
    bash "$SCRIPT_DIR/rocketmq/pre.sh"

    # 启动服务
    cd "$SCRIPT_DIR/rocketmq"
    docker-compose up -d

    # 等待服务启动
    for service in "${ROCKETMQ_SERVICES[@]}"; do
        wait_for_service "$service"
    done

    cd - > /dev/null
}

stop_service() {
    local service=$1
    local container_name="lushop-$service"

    if docker ps -a --format "table {{.Names}}" | grep -q "^${container_name}$"; then
        log_info "停止 $service..."
        docker stop "$container_name" 2>/dev/null || true
        docker rm "$container_name" 2>/dev/null || true
    fi
}

stop_rocketmq() {
    log_info "停止 RocketMQ..."
    cd "$SCRIPT_DIR/rocketmq"
    docker-compose down 2>/dev/null || true
    cd - > /dev/null
}

show_status() {
    echo -e "\n${BLUE}=== 服务状态 ===${NC}"
    docker ps --filter "name=lushop" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

    echo -e "\n${BLUE}=== RocketMQ 服务状态 ===${NC}"
    cd "$SCRIPT_DIR/rocketmq"
    docker-compose ps
    cd - > /dev/null
}

show_usage() {
    echo "用法: $0 [start|stop|restart|status]"
    echo ""
    echo "命令:"
    echo "  start   - 启动所有服务"
    echo "  stop    - 停止所有服务"
    echo "  restart - 重启所有服务"
    echo "  status  - 显示服务状态"
    echo ""
    echo "注意:"
    echo "  - 首次运行前请配置 .env 文件"
    echo "  - 确保 Docker 服务正在运行"
    echo "  - MySQL 初始化需要时间，请耐心等待"
}

# 主逻辑
main() {
    local action=${1:-status}

    check_dependencies

    case "$action" in
        start)
            log_info "开始启动所有服务..."

            # 基础设施服务
            start_mysql
            start_service "redis"
            start_service "consul"
            start_service "nacos"

            # 搜索和监控服务
            start_service "elasticsearch"
            start_service "kibana"
            start_service "grafana"
            start_service "prometheus"
            start_service "jaeger"

            # RocketMQ
            start_rocketmq

            log_info "所有服务启动完成！"
            echo -e "\n${GREEN}访问地址:${NC}"
            echo "  MySQL: localhost:3306"
            echo "  Redis: localhost:6379"
            echo "  Consul: http://localhost:8500"
            echo "  Nacos: http://localhost:8848"
            echo "  Elasticsearch: http://localhost:9200"
            echo "  Kibana: http://localhost:5601"
            echo "  Grafana: http://localhost:3000"
            echo "  Prometheus: http://localhost:9090"
            echo "  Jaeger: http://localhost:16686"
            echo "  RocketMQ Dashboard: http://localhost:8682"
            ;;

        stop)
            log_info "停止所有服务..."

            # 停止 RocketMQ
            stop_rocketmq

            # 停止其他服务（反序）
            for service_info in "${SERVICES[@]}"; do
                local service=$(echo "$service_info" | cut -d: -f1)
                stop_service "$service"
            done

            log_info "所有服务已停止"
            ;;

        restart)
            log_info "重启所有服务..."
            "$0" stop
            sleep 3
            "$0" start
            ;;

        status)
            show_status
            ;;

        *)
            show_usage
            exit 1
            ;;
    esac
}

main "$@"
