FROM golang:1.25 AS builder

# 复制 order 服务代码
COPY order/ /src/
# 复制 goods 服务代码（order 的依赖）
COPY goods/ /workspace/goods/

WORKDIR /src

# 设置环境变量（全局生效）
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
ENV GOPRIVATE=""
ENV GOSUMDB=sum.golang.org

# 设置 go.mod 以便使用本地 goods 模块
RUN go mod edit -replace goods=/workspace/goods

# 设置超时和重试
ENV GOTIMEOUT=300

# 下载依赖（带重试机制）
RUN go mod download || sleep 5 && go mod download || sleep 5 && go mod download

# 更新依赖并构建
RUN go mod tidy && make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./order", "-conf", "/data/conf"]
