# MySQL 慢查询监控配置指南

## 方案概述

使用 **mysqld_exporter** 监控 MySQL 性能指标，包括慢查询。

## 1. 配置 MySQL 慢查询日志

### 方式1: 修改 MySQL 配置文件

创建 `deploy/mysql/my.cnf`:

```ini
[mysqld]
# 启用慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow-query.log

# 慢查询阈值（秒）
long_query_time = 2

# 记录未使用索引的查询
log_queries_not_using_indexes = 1

# 限制未使用索引的查询记录频率
log_throttle_queries_not_using_indexes = 10

# 记录慢管理语句
log_slow_admin_statements = 1

# 记录慢从库语句
log_slow_slave_statements = 1
```

### 方式2: 运行时动态设置

```sql
-- 启用慢查询日志
SET GLOBAL slow_query_log = 'ON';

-- 设置慢查询阈值为2秒
SET GLOBAL long_query_time = 2;

-- 记录未使用索引的查询
SET GLOBAL log_queries_not_using_indexes = 'ON';

-- 查看当前配置
SHOW VARIABLES LIKE 'slow_query%';
SHOW VARIABLES LIKE 'long_query_time';
```

## 2. 部署 mysqld_exporter

已在 `docker-compose.yml` 中配置：

```yaml
mysqld-exporter:
  image: prom/mysqld-exporter:latest
  container_name: lushop-mysqld-exporter
  restart: always
  ports:
    - "9104:9104"
  environment:
    - DATA_SOURCE_NAME=root:root123456@(mysql:3306)/
  command:
    - '--collect.info_schema.processlist'      # 进程列表
    - '--collect.info_schema.innodb_metrics'   # InnoDB指标
    - '--collect.info_schema.tablestats'       # 表统计
    - '--collect.perf_schema.eventswaits'      # 事件等待
    - '--collect.perf_schema.tableiowaits'     # 表IO等待
    - '--collect.perf_schema.indexiowaits'     # 索引IO等待
```

## 3. 启动服务

```bash
# 启动 MySQL 和 mysqld_exporter
docker-compose up -d mysql mysqld-exporter

# 查看 exporter 日志
docker logs -f lushop-mysqld-exporter

# 测试 metrics 端点
curl http://localhost:9104/metrics
```

## 4. 关键监控指标

### 慢查询相关指标

```promql
# 慢查询总数
mysql_global_status_slow_queries

# 慢查询增长率（每秒）
rate(mysql_global_status_slow_queries[5m])

# 查询执行时间
mysql_perf_schema_events_statements_total

# 表扫描次数（全表扫描）
mysql_global_status_select_scan

# 未使用索引的查询
mysql_global_status_select_full_join
```

### 性能指标

```promql
# 查询响应时间
mysql_perf_schema_table_io_waits_seconds_total

# 锁等待时间
mysql_perf_schema_table_lock_waits_seconds_total

# 连接数
mysql_global_status_threads_connected

# QPS (每秒查询数)
rate(mysql_global_status_questions[1m])

# TPS (每秒事务数)
rate(mysql_global_status_commands_total{command="commit"}[1m])
```

## 5. Prometheus 查询示例

### 查询1: 慢查询增长率

```promql
# 每分钟慢查询数
rate(mysql_global_status_slow_queries[1m]) * 60
```

### 查询2: 全表扫描率

```promql
# 全表扫描占比
rate(mysql_global_status_select_scan[5m]) / 
rate(mysql_global_status_questions[5m]) * 100
```

### 查询3: 平均查询时间

```promql
# 平均查询执行时间（毫秒）
rate(mysql_perf_schema_events_statements_seconds_total[5m]) /
rate(mysql_perf_schema_events_statements_total[5m]) * 1000
```

### 查询4: 锁等待时间

```promql
# 表锁等待时间
rate(mysql_perf_schema_table_lock_waits_seconds_total[5m])
```

### 查询5: 慢查询 Top 表

```promql
# 按表统计的慢查询
topk(10, 
  rate(mysql_perf_schema_table_io_waits_seconds_total[5m])
)
```

## 6. Grafana 仪表板配置

### 导入预制仪表板

1. 访问 Grafana: http://localhost:3000
2. 登录 (admin/admin)
3. 导入仪表板:
   - Dashboard ID: **7362** (MySQL Overview)
   - Dashboard ID: **11323** (MySQL InnoDB Metrics)

### 自定义慢查询面板

#### Panel 1: 慢查询趋势

```promql
# Query
rate(mysql_global_status_slow_queries[5m]) * 60

# Visualization: Graph
# Title: 慢查询数/分钟
```

#### Panel 2: 慢查询占比

```promql
# Query
(rate(mysql_global_status_slow_queries[5m]) / 
 rate(mysql_global_status_questions[5m])) * 100

# Visualization: Gauge
# Title: 慢查询占比 (%)
# Thresholds: 0-1 (green), 1-5 (yellow), 5+ (red)
```

#### Panel 3: Top 慢查询表

```promql
# Query
topk(10, 
  sum by (schema, table) (
    rate(mysql_perf_schema_table_io_waits_seconds_total[5m])
  )
)

# Visualization: Table
# Title: 慢查询 Top 10 表
```

#### Panel 4: 全表扫描率

```promql
# Query
rate(mysql_global_status_select_scan[5m]) * 60

# Visualization: Graph
# Title: 全表扫描次数/分钟
```

## 7. 告警规则配置

创建 `deploy/prometheus/alerts/mysql_alerts.yml`:

```yaml
groups:
  - name: mysql_slow_query_alerts
    interval: 30s
    rules:
      # 慢查询数量告警
      - alert: HighSlowQueries
        expr: rate(mysql_global_status_slow_queries[5m]) * 60 > 10
        for: 5m
        labels:
          severity: warning
          service: mysql
        annotations:
          summary: "MySQL 慢查询过多"
          description: "实例 {{ $labels.instance }} 每分钟慢查询数: {{ $value | printf \"%.2f\" }}"

      # 慢查询占比告警
      - alert: HighSlowQueryRatio
        expr: |
          (rate(mysql_global_status_slow_queries[5m]) / 
           rate(mysql_global_status_questions[5m])) * 100 > 5
        for: 5m
        labels:
          severity: warning
          service: mysql
        annotations:
          summary: "MySQL 慢查询占比过高"
          description: "实例 {{ $labels.instance }} 慢查询占比: {{ $value | printf \"%.2f\" }}%"

      # 全表扫描告警
      - alert: HighTableScans
        expr: rate(mysql_global_status_select_scan[5m]) * 60 > 100
        for: 5m
        labels:
          severity: warning
          service: mysql
        annotations:
          summary: "MySQL 全表扫描过多"
          description: "实例 {{ $labels.instance }} 每分钟全表扫描: {{ $value | printf \"%.2f\" }}"

      # 锁等待告警
      - alert: HighLockWaits
        expr: rate(mysql_perf_schema_table_lock_waits_seconds_total[5m]) > 1
        for: 5m
        labels:
          severity: warning
          service: mysql
        annotations:
          summary: "MySQL 锁等待时间过长"
          description: "实例 {{ $labels.instance }} 锁等待时间: {{ $value | printf \"%.2f\" }}s"

      # 连接数告警
      - alert: HighConnections
        expr: mysql_global_status_threads_connected > 100
        for: 5m
        labels:
          severity: warning
          service: mysql
        annotations:
          summary: "MySQL 连接数过高"
          description: "实例 {{ $labels.instance }} 当前连接数: {{ $value }}"
```

在 `prometheus.yml` 中引用告警规则:

```yaml
rule_files:
  - '/etc/prometheus/alerts/*.yml'
```

## 8. 应用层集成 - GORM 慢查询日志

在你的 Go 应用中配置 GORM 慢查询日志：

```go
// internal/data/data.go
import (
    "time"
    "gorm.io/gorm/logger"
)

func NewDB(c *conf.Data, l log.Logger) *gorm.DB {
    // 自定义日志配置
    newLogger := logger.New(
        log.NewStdLogger(os.Stdout),
        logger.Config{
            SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
            LogLevel:                  logger.Warn,            // 日志级别
            IgnoreRecordNotFoundError: true,                   // 忽略 ErrRecordNotFound
            Colorful:                  false,                  // 禁用彩色打印
        },
    )

    db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
        Logger: newLogger,
    })
    
    if err != nil {
        panic(err)
    }
    
    return db
}
```

## 9. 慢查询分析工具

### 使用 pt-query-digest 分析慢查询日志

```bash
# 安装 percona-toolkit
docker exec -it lushop-mysql bash
apt-get update && apt-get install -y percona-toolkit

# 分析慢查询日志
pt-query-digest /var/log/mysql/slow-query.log

# 输出到文件
pt-query-digest /var/log/mysql/slow-query.log > slow-query-report.txt
```

### 使用 mysqldumpslow

```bash
# 查看最慢的10条查询
docker exec lushop-mysql mysqldumpslow -s t -t 10 /var/log/mysql/slow-query.log

# 查看访问次数最多的10条查询
docker exec lushop-mysql mysqldumpslow -s c -t 10 /var/log/mysql/slow-query.log
```

## 10. 优化建议

### 识别慢查询

1. **查看 Prometheus 指标**
   ```promql
   rate(mysql_global_status_slow_queries[5m]) > 0
   ```

2. **查看慢查询日志**
   ```bash
   docker exec lushop-mysql tail -f /var/log/mysql/slow-query.log
   ```

3. **查看应用日志**
   - 查找 GORM 慢查询警告

### 优化步骤

1. **添加索引**
   ```sql
   -- 分析查询
   EXPLAIN SELECT * FROM goods WHERE category_id = 1;
   
   -- 添加索引
   CREATE INDEX idx_category_id ON goods(category_id);
   ```

2. **优化查询**
   - 避免 SELECT *
   - 使用 LIMIT 限制结果
   - 避免子查询，使用 JOIN

3. **分页优化**
   ```go
   // 使用游标分页代替 OFFSET
   db.Where("id > ?", lastID).Limit(pageSize).Find(&results)
   ```

4. **缓存热点数据**
   ```go
   // 使用 Redis 缓存
   if cached, err := redis.Get(key); err == nil {
       return cached
   }
   ```

## 11. 监控检查清单

- [ ] MySQL 慢查询日志已启用
- [ ] mysqld_exporter 正常运行
- [ ] Prometheus 正常抓取 MySQL 指标
- [ ] Grafana 仪表板已配置
- [ ] 告警规则已设置
- [ ] 应用层慢查询日志已配置
- [ ] 定期分析慢查询日志

## 12. 常见问题

### Q1: mysqld_exporter 无法连接 MySQL

**解决方案**:
```bash
# 检查 MySQL 连接
docker exec lushop-mysql mysql -uroot -proot123456 -e "SELECT 1"

# 检查 exporter 日志
docker logs lushop-mysqld-exporter
```

### Q2: 没有慢查询数据

**解决方案**:
```sql
-- 确认慢查询日志已启用
SHOW VARIABLES LIKE 'slow_query_log';

-- 降低阈值测试
SET GLOBAL long_query_time = 0.1;
```

### Q3: 指标数据不准确

**解决方案**:
```bash
# 重启 exporter
docker restart lushop-mysqld-exporter

# 检查 Prometheus targets
# 访问 http://localhost:9090/targets
```

## 13. 参考资源

- [mysqld_exporter GitHub](https://github.com/prometheus/mysqld_exporter)
- [MySQL Performance Schema](https://dev.mysql.com/doc/refman/8.0/en/performance-schema.html)
- [Grafana MySQL Dashboard](https://grafana.com/grafana/dashboards/7362)
- [Percona Toolkit](https://www.percona.com/software/database-tools/percona-toolkit)

## 14. 快速开始

```bash
# 1. 启动服务
docker-compose up -d mysql mysqld-exporter prometheus grafana

# 2. 配置 MySQL 慢查询
docker exec lushop-mysql mysql -uroot -proot123456 -e "
  SET GLOBAL slow_query_log = 'ON';
  SET GLOBAL long_query_time = 2;
  SET GLOBAL log_queries_not_using_indexes = 'ON';
"

# 3. 验证 exporter
curl http://localhost:9104/metrics | grep mysql_global_status_slow_queries

# 4. 访问 Grafana
open http://localhost:3000

# 5. 导入仪表板 ID: 7362
```

完成！现在你可以在 Prometheus 和 Grafana 中监控 MySQL 慢查询了。
