# Kubernetes 部署指引

该目录提供了将 `lushop` 项目部署到 Kubernetes 集群的基础清单。整体结构：

```
k8s/
├── base/            # 组件基础清单，可直接 kubectl apply
│   ├── namespace.yaml
│   ├── redis/
│   ├── mysql/
│   ├── rocketmq/
│   ├── prometheus/
│   └── grafana/
└── overlays/
    └── dev/         # 示例环境，可在此叠加镜像、资源等差异
```

## 部署步骤（详细）

1. **准备命名空间与存储**
   - 所有清单默认部署到 `lushop` 命名空间，若需修改请调整 `k8s/base/namespace.yaml` 并在 overlay 中覆盖。
   - Redis、MySQL、RocketMQ、Prometheus、Grafana 的 PVC/`volumeClaimTemplates` 默认使用集群 `default` StorageClass。如需指定（例如 `local-path`、`rook-ceph`），请在各 YAML 中显式设置 `storageClassName`。

2. **准备 Secret 与配置**
   - 仓库中的 Secret 仅为演示值，部署前务必替换，可通过：
     - `kustomize edit set secret` 在 base/overlay 中更新；
     - SealedSecret、External Secrets 或 CI/CD 管道注入真实凭据。
   - 需要重点替换的 Secret 文件：
     - `redis-auth`（`k8s/base/redis/secret.yaml`）：Redis 密码。
     - `mysql-auth`（`k8s/base/mysql/secret.yaml`）：MySQL root/业务账号。
     - `rocketmq-credentials`（`k8s/base/rocketmq/secret.yaml`）：ACL access/secret key 及 dashboard 登录。
     - `grafana-admin`（`k8s/base/grafana/secret.yaml`）：Grafana 管理员账号。
   - Prometheus ConfigMap 来自 `deploy/prometheus/conf/prometheus.yaml`；若需要追加抓取目标，请直接修改该文件或在 overlay 中覆盖。

3. **应用基础设施组件**
   ```bash
   kubectl apply -k k8s/base
   ```
   - 使用 `kubectl get pods -n lushop` 观察 Redis、MySQL、RocketMQ（namesrv/broker/proxy/dashboard）、Prometheus、Grafana 均进入 `Ready`。
   - 也可以在 `k8s/base/<component>` 目录中按组件分步 `kubectl apply -k`。

4. **部署业务服务**
   - 建议在 `k8s/base` 下新增 `services/` 目录或在 overlay 中维护所有微服务的 `Deployment + Service`。
   - 通过 ConfigMap/Secret 注入数据库、Redis、RocketMQ 等连接信息；将 `deploy/` 目录中的环境变量迁移到 Kubernetes 配置中以保持一致。

5. **使用 overlays（可选）**
   - `k8s/overlays/dev` 为示例，可复制成 `stage`、`prod`。利用 overlay patches 覆盖镜像 tag、副本数、资源限制、Service 类型、NodeSelector 等。
   - 渲染并部署 overlay：
     ```bash
     kustomize build k8s/overlays/dev | kubectl apply -f -
     ```

6. **监控、可观测与网络暴露**
   - Prometheus 已挂载抓取配置，可在 ConfigMap 中扩展 ServiceMonitor 或静态 targets。
   - Grafana 仅提供基础 Deployment 与存储卷，建议引入 ConfigMap/provisioning 目录自动导入数据源与仪表盘。
   - 所有 Service 默认为 `ClusterIP`，如需外部访问可新增 Ingress、Gateway、LoadBalancer，或直接使用 `kubectl port-forward`。

7. **后续运维建议**
   - 默认副本为 1，适用于本地/测试。生产环境需结合 overlay 调整副本并配置 HPA/VPA、PDB。
   - 为 Grafana、RocketMQ Dashboard、Prometheus 等敏感组件配置 NetworkPolicy、防火墙或额外认证。
   - 针对 MySQL、RocketMQ Store、Prometheus 等数据卷规划备份（Velero、快照、CronJob）。
   - 按需补充 Nacos/Consul、Elasticsearch/Kibana 以及业务服务清单，使其与 `deploy/` 脚本保持一致。

## 自定义建议

- **资源与副本**：默认副本为 1，适用于 PoC/测试；生产环境请在 overlay 中扩容并配置 HPA/VPA、PDB。
- **网络暴露**：默认 `ClusterIP`，根据需要配置 Ingress/Gateway/LoadBalancer 并启用 TLS/认证。
- **安全**：借助 Secret 管理敏感信息，为 Grafana、RocketMQ Dashboard、Prometheus 等增加 NetworkPolicy 或统一认证机制。
- **备份**：MySQL、RocketMQ Store、Prometheus 等数据卷需制定备份策略（Velero、快照、CronJob）。
- **持续扩展**：在 `k8s/` 中持续补充业务服务、Nacos/Consul、Elasticsearch/Kibana 等，使配置与 `deploy/` 脚本对齐。

## 单机部署示例（从零开始）

> 假设你有一台全新 Linux 云主机（>=4 vCPU / 8 GB RAM / 50 GB SSD）。以下步骤以 k3s 为例，Minikube/kind 同理。

1. **准备环境**
   1. 安装依赖：
      ```bash
      curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--write-kubeconfig-mode 644" sh -
      sudo ln -sf /usr/local/bin/kubectl /usr/bin/kubectl
      curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
      sudo mv kustomize /usr/local/bin/
      ```
   2. 验证集群：
      ```bash
      kubectl get nodes
      ```
      输出 Ready 即可继续。
   3. 拉取仓库：
      ```bash
      git clone https://github.com/<your-org>/lushop-kratos.git
      cd lushop-kratos
      ```

2. **命名空间与存储**
   - 默认使用 `lushop` 命名空间。若想变更（例如 `prod`），需同步修改 `k8s/base/namespace.yaml` 与 overlay。
   - k3s 自动提供 `local-path` StorageClass，不需要额外配置。如果你用 Minikube/kind，请提前启用对应 CSI 或把各组件的 `storageClassName` 改为你实际可用的名称。

3. **配置 Secret 与 ConfigMap**
   1. 替换示例 Secret，可直接编辑对应文件或使用命令：
      ```bash
      cd k8s/base
      kustomize edit set secret redis-auth --from-literal=password=<redis_pwd>
      kustomize edit set secret mysql-auth --from-literal=rootPassword=<mysql_root> --from-literal=appPassword=<mysql_app>
      kustomize edit set secret rocketmq-credentials --from-literal=accessKey=<ak> --from-literal=secretKey=<sk>
      kustomize edit set secret grafana-admin --from-literal=admin-user=admin --from-literal=admin-password=<grafana_pwd>
      cd ../../
      ```
   2. 将 `deploy/` 目录中的 `.env` 或脚本内环境变量抄写进 ConfigMap/Secret，确保微服务在集群上能读取（推荐放在 `k8s/base/services/config/`）。
   3. 如需自定义 Prometheus 抓取目标，编辑 `deploy/prometheus/conf/prometheus.yaml` 或在 overlay 中提供 patches。

4. **部署基础组件**
   ```bash
   kubectl apply -k k8s/base
   kubectl get pods -n lushop
   ```
   - 等待所有 Pod 进入 `Running`/`Ready`。如 PVC Pending，可执行 `kubectl describe pvc <name> -n lushop` 查看 storageClass 是否匹配。
   - RocketMQ 依赖多组件，若其中一个 CrashLoop，使用 `kubectl logs <pod> -n lushop` 排查端口/配置。

5. **部署业务服务**
   1. 在 `k8s/base/services/` 下为每个微服务创建 `Deployment` 与 `Service`（可参考 `deploy/` 中 docker-compose 的镜像与环境变量）。
   2. 通过 `envFrom`/`volumeMounts` 引用第 3 步准备好的 ConfigMap/Secret。
   3. 如果需要按环境差异调整镜像 tag、资源、副本，复制 `k8s/overlays/dev` 为 `k8s/overlays/local-single-node`（或 `prod`），并在 `kustomization.yaml` 中引用对应 patches。
   4. 应用服务：
      ```bash
      kustomize build k8s/overlays/local-single-node | kubectl apply -f -
      ```

6. **暴露与验证**
   - 查看服务状态：
     ```bash
     kubectl get svc -n lushop
     ```
   - 本地调试可 `kubectl port-forward svc/grafana 3000:3000 -n lushop`、`svc/api-gateway 8080:80` 等。
   - 需要外部访问时创建 Ingress（k3s 自带 Traefik）：
     ```bash
     kubectl apply -f k8s/base/ingress/lushop-dashboard.yaml
     ```
   - 登录 Grafana（默认 admin/<你的密码>），添加 Prometheus 数据源 `http://prometheus:9090` 并导入仪表盘；访问 RocketMQ Dashboard 验证 ACL。

7. **运维加固**
   - 资源：在 overlay 中为关键服务设置 `resources.requests/limits` 并开启 HPA：
     ```bash
     kubectl autoscale deployment api-gateway -n lushop --cpu-percent=80 --min=1 --max=3
     ```
   - 网络：为 Grafana、Prometheus、RocketMQ Dashboard 编写 NetworkPolicy 限制来源 IP；单机可使用 `namespaceSelector` + `podSelector`。
   - 备份：使用 Velero 或 `kubectl cp` + CronJob 对 MySQL 数据卷和 RocketMQ Store 做快照；Prometheus tsdb 同理。
   - 监控告警：在 Prometheus 中加入 Exporter（node-exporter、kube-state-metrics），并在 Grafana 配置仪表盘/告警规则。

完成以上 7 个步骤，即可在单机环境从零拉起 Kubernetes、部署 `lushop` 所需中间件与业务服务，并具备基本的访问、监控与运维能力。

