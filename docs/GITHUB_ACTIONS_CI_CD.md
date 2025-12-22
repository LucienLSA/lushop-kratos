# GitHub Actions CI/CD 使用说明

本文档介绍仓库中新增的 GitHub Actions workflow、helper 脚本、所需的 secrets，以及如何触发和在本地/CI 上复现构建与部署流程。

## 增加的文件
- `.github/workflows/ci.yml` — CI：格式化、静态检查、单元测试、构建二进制（不推镜像），并上传 artifact。触发器：`push` / `pull_request` 到 `main`。  
- `.github/workflows/release.yml` — Release：使用 buildx 构建并推送镜像到 registry（触发：打 tag 或手动）  
- `.github/workflows/deploy.yml` — Deploy：使用 `kubectl` 从 `KUBE_CONFIG_DATA` 部署镜像（触发：手动）  
- `scripts/ci/build-for-service.sh` — Helper：在本地或 runner 上复用的单服务 build/test 脚本

## 目标服务列表（matrix）
- `lushop`  
- `service/goods`  
- `service/inventory`  
- `service/order`  
- `service/user`  
- `service/userauth`  
- `service/userop`

## 必要的 GitHub Secrets
- `DOCKER_REGISTRY`：容器注册表地址（例如 `docker.io`、`ghcr.io` 或私有 registry）  
- `DOCKER_USERNAME`：registry 用户名或 PAT  
- `DOCKER_PASSWORD`：registry 密码或 PAT  
- `IMAGE_REPOSITORY_PREFIX`：镜像前缀/组织名，例如 `myorg`，构成最终镜像为 `${{ secrets.DOCKER_REGISTRY }}/${{ secrets.IMAGE_REPOSITORY_PREFIX }}/${service}:${tag}`  
- `KUBE_CONFIG_DATA`：base64 编码的 kubeconfig 内容（用于 `deploy.yml`）  
- `K8S_NAMESPACE`：目标 Kubernetes 命名空间（例如 `default`）

可选：`KUBECTL_VERSION`、`DOCKER_BUILDX` 相关配置或额外的 secrets（例如私有 registry 的 CA 或凭证）。

## 如何生成 `KUBE_CONFIG_DATA`
在有访问权限的机器上（例如你的本地）。示例（Linux/macOS）：

```bash
cat ~/.kube/config | base64 | tr -d '\n' | pbcopy   # macOS
cat ~/.kube/config | base64 | tr -d '\n' | xclip    # Linux + xclip
```

在 Windows PowerShell：

```powershell
[Convert]::ToBase64String([IO.File]::ReadAllBytes("$env:USERPROFILE\.kube\config")) | clip
```

把编码后的字符串粘贴到仓库 Settings → Secrets → New repository secret，名称为 `KUBE_CONFIG_DATA`。

## 使用说明与示例

- 在 CI 中运行（自动）：当你 push 到 `main` 或发起 PR，会自动执行 `ci.yml`（格式化、lint、test、构建并上传 artifact）。  
- 手动触发 Release：在仓库中打 tag（例如 `git tag v1.2.3 && git push origin v1.2.3`），`release.yml` 会触发并把镜像推到你配置的 registry。也可在 Actions 界面手动触发 workflow_dispatch。  
- 手动触发 Deploy：在 Actions → Deploy to Kubernetes → Run workflow，传入 `service` 和 `image_tag`（例如 image_tag 为 release 构建得到的 sha 或 tag）。

示例：本地复现单服务构建

```bash
# 在仓库根目录下运行（需要已安装 go）
bash scripts/ci/build-for-service.sh service/goods
# 构建产物输出到 build/goods
```

示例：手动用 kubectl 更新 deployment（与 deploy.yml 同逻辑）

```bash
# 假设 deployment 名称与服务名相同
kubectl -n <namespace> set image deployment/goods goods=<REGISTRY>/<PREFIX>/goods:<TAG> --record
kubectl -n <namespace> rollout status deployment/goods --timeout=2m
```

## 注意事项与建议
- `golangci-lint` 当前在 workflow 中以非阻断方式运行（lint 失败不会终止 job）。如果需要强制 lint 成功请移除 `|| true`。  
- `ci.yml` 里的 `go-version` 目前固定为 `1.20`，可根据项目 `go.mod` 修改。  
- 如果仓库较大或 CI 时长需优化，建议改为“仅构建变更过的 service”（通过 `git diff --name-only` 找到受影响的 top-level service）。我可以为你添加该优化步骤。  
- 如果你使用 Helm 或 kustomize 管理 manifests，我可以把 `deploy.yml` 改为 `helm upgrade --install` 或 `kubectl apply -k` 的形式，并加入镜像替换/values 参数化。

## 我可以接着帮你做（可选）
- 把 `deploy.yml` 改为基于 `kustomize`（overlay）或 `helm`（chart）部署。  
- 把 `release.yml` 的触发从 tag 改为 `push to main`。  
- 添加“仅构建变更服务”的 CI 优化。  

---  
如果你需要我现在把该文档提交到其它位置或调整为英文版，请告诉我。  


