FROM golang:1.25 AS builder

COPY . /src
WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
ENV GOPRIVATE=""
ENV GOSUMDB=sum.golang.org
ENV GOTIMEOUT=300

# 下载依赖（带重试机制）
RUN go mod download || sleep 10 && go mod download || sleep 10 && go mod download

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

CMD ["./inventory", "-conf", "/data/conf"]
