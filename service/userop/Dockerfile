FROM golang:1.25 AS builder

COPY . /src
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
ENV GOPRIVATE=""
ENV GOSUMDB=sum.golang.org
ENV GOTIMEOUT=300

# 先下载依赖，带重试机制
RUN go mod download || sleep 10 && go mod download || sleep 10 && go mod download

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
CMD ["./userop", "-conf", "/data/conf"]
