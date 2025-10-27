# Lushop 压力测试

简单的压力测试工具，使用 Apache Bench 进行基本压测。

## 快速开始

### 1. 启动服务
```bash
cd /home/zzx/GoProject/lushop-kratos-main
./deploy.sh start
```

### 2. 运行压力测试
```bash
cd stress-test

# 快速测试（约1分钟）
./quick_stress_test.sh
```

这个脚本会测试3个核心接口：
- 商品列表 (100并发, 1000请求)
- 商品详情 (200并发, 2000请求)  
- 库存查询 (300并发, 3000请求)

## 性能指标

- **QPS**: 每秒请求数
- **响应时间**: 平均、最小、最大
- **成功率**: 目标 > 99%

## 测试目标

- QPS ≥ 1000
- 平均响应时间 < 300ms
- 错误率 < 1%

## 报告

测试报告保存在 `reports/logs/` 目录下。

## 需要安装

```bash
sudo apt-get install apache2-utils
```

## 问题排查

```bash
# 检查服务状态
curl http://localhost:8001/metrics

# 查看服务日志
docker logs lushop-api-gateway
```

