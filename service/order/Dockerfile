FROM golang:1.25 AS builder

# 复制 order 服务代码
COPY order/ /src/
# 复制 goods 服务代码（order 的依赖）
COPY goods/ /workspace/goods/

WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
ENV GOPRIVATE=""
ENV GOSUMDB=sum.golang.org
ENV GOTIMEOUT=300

# 设置 go.mod 以便使用本地 goods 模块
RUN go mod edit -replace goods=/workspace/goods

# 下载依赖（带重试机制）
RUN go mod download || sleep 5 && go mod download || sleep 5 && go mod download

# 整理依赖并构建二进制（使用仓库的 make build）
RUN go mod tidy && make build

FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get autoremove -y && apt-get autoclean -y

# 非 root 用户
RUN groupadd --system app && useradd --system --gid app --create-home --home-dir /nonexistent app || true
RUN mkdir -p /app && chown -R app:app /app

COPY --from=builder /src/bin /app

USER app
WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./order", "-conf", "/data/conf"]
