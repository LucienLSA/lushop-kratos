# 快速启动指南

## 🎯 当前状态

✅ **基础设施已启动并正常运行**

### 正在运行的服务
- ✅ MySQL: 127.0.0.1:3306
- ✅ Redis: 127.0.0.1:6379
- ✅ Consul: http://127.0.0.1:8500
- ✅ Jaeger: http://127.0.0.1:16686
- ✅ Nacos: http://127.0.0.1:8848/nacos
- ✅ Prometheus: http://127.0.0.1:9090
- ✅ Grafana: http://127.0.0.1:3000

## 🚀 启动微服务

现在你可以在**多个终端窗口**中分别运行以下命令启动各个微服务：

### 1. User Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/user
kratos run
```

### 2. UserAuth Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/userauth
kratos run
```

### 3. Goods Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/goods
kratos run
```

### 4. Order Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/order
kratos run
```

### 5. Inventory Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/inventory
kratos run
```

### 6. UserOp Service
```bash
cd /home/zzx/GoProject/lushop-kratos-main/service/userop
kratos run
```

### 7. API Gateway (最后启动)
```bash
cd /home/zzx/GoProject/lushop-kratos-main/lushop
kratos run
```

## 📝 管理命令

### 查看状态
```bash
./deploy-local.sh status
```

### 停止基础设施
```bash
./deploy-local.sh stop
```

### 重启基础设施
```bash
./deploy-local.sh restart
```

## 🔍 访问地址

- **API Gateway**: http://127.0.0.1:8001
- **Consul UI**: http://127.0.0.1:8500
- **Jaeger UI**: http://127.0.0.1:16686
- **Nacos**: http://127.0.0.1:8848/nacos (用户名/密码: nacos/nacos)
- **Prometheus**: http://127.0.0.1:9090
- **Grafana**: http://127.0.0.1:3000 (用户名/密码: admin/admin)

## ⚠️ 注意事项

1. **顺序很重要**: 先启动基础服务，再启动微服务
2. **端口占用**: 确保 50051-50056 端口未被占用
3. **依赖检查**: 确保 kratos CLI 已安装
4. **配置文件**: 所有服务的配置已改为使用 127.0.0.1

## 🐛 故障排查

### 端口被占用
```bash
sudo lsof -i :50051  # 查看端口占用
kill -9 <PID>        # 关闭进程
```

### 查看日志
每个服务的日志会在各自的终端中显示

### 停止服务
在运行 kratos 的终端中按 `Ctrl+C`
