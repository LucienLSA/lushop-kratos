# Build stage
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /src

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/userauth ./cmd/userauth

# Runtime stage
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/userauth .

# 复制配置文件
COPY configs /app/configs

# 暴露端口
EXPOSE 50056

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 50056 || exit 1

# 运行服务
CMD ["./userauth", "-conf", "/app/configs"]
