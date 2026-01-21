# Go DevOps 指南

## 目录
- [1. Docker 容器化](#1-docker-容器化)
- [2. Docker Compose](#2-docker-compose)
- [3. Kubernetes 编排](#3-kubernetes-编排)
- [4. CI/CD 流水线](#4-cicd-流水线)
- [5. 监控与日志](#5-监控与日志)

---

## 1. Docker 容器化

### 1.1 Docker 基础概念

```
┌─────────────────────────────────────────────────────────────┐
│                      Docker 架构                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────┐      ┌─────────────────────────────┐     │
│   │   Client    │─────▶│          Daemon              │     │
│   └─────────────┘      │  ┌─────────────────────┐    │     │
│                        │  │     Containers      │    │     │
│                        │  │  ┌─────┐ ┌─────┐    │    │     │
│                        │  │  │App1 │ │App2 │ ... │    │     │
│                        │  │  └─────┘ └─────┘    │    │     │
│                        │  └─────────────────────┘    │     │
│                        │  ┌─────────────────────┐    │     │
│                        │  │     Images          │    │     │
│                        │  │  ┌─────┐ ┌─────┐    │    │     │
│                        │  │  │Img1 │ │Img2 │ ... │    │     │
│                        │  │  └─────┘ └─────┘    │    │     │
│                        │  └─────────────────────┘    │     │
│                        └─────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Dockerfile 最佳实践

```dockerfile
# 1. 使用多阶段构建减小镜像体积

# ========== 构建阶段 ==========
# 使用官方 Go 镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
# CGO_ENABLED=0 创建静态链接的二进制文件
# -ldflags="-s -w" 减小二进制文件大小
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o main \
    ./cmd/main.go

# ========== 运行阶段 ==========
# 使用轻量级 Alpine 镜像
FROM alpine:3.19 AS runner

# 创建非 root 用户
RUN addgroup --system --gid 1001 appgroup && \
    adduser --system --uid 1001 appuser

# 设置工作目录
WORKDIR /home/appuser

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 更改文件所有者
RUN chown -R appuser:appgroup /home/appuser

# 切换用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["./main"]
```

### 1.3 .dockerignore

```
# Git
.git
.gitignore

# IDE
.idea
.vscode
*.swp

# 文档
README.md
LICENSE
docs/

# 测试
*_test.go
coverage.out

# 本地配置
local.yaml
.env

# 构建产物
bin/
dist/

# 临时文件
*.log
tmp/
```

### 1.4 常用 Docker 命令

```bash
# ========== 镜像操作 ==========

# 构建镜像
docker build -t myapp:latest .

# 标记镜像
docker tag myapp:latest myregistry.com/myapp:v1.0.0

# 推送镜像
docker push myregistry.com/myapp:v1.0.0

# 拉取镜像
docker pull myregistry.com/myapp:v1.0.0

# 列出镜像
docker images

# 删除镜像
docker rmi myapp:latest

# 清理未使用镜像
docker image prune -a

# ========== 容器操作 ==========

# 运行容器
docker run -d \
    --name myapp-container \
    -p 8080:8080 \
    -e ENV=production \
    -v /data/myapp:/data \
    myapp:latest

# 列出运行中容器
docker ps

# 列出所有容器
docker ps -a

# 查看容器日志
docker logs -f myapp-container

# 查看容器资源使用
docker stats myapp-container

# 进入容器
docker exec -it myapp-container /bin/sh

# 停止容器
docker stop myapp-container

# 删除容器
docker rm myapp-container

# 强制删除运行中容器
docker rm -f myapp-container

# 查看容器详细信息
docker inspect myapp-container

# ========== Docker Hub ==========

# 登录
docker login

# 登出
docker logout
```

### 1.5 多阶段构建示例

```dockerfile
# ========== 开发阶段 ==========
FROM golang:1.21-alpine AS development

WORKDIR /app

# 安装开发工具
RUN apk add --no-cache git delve

# 复制依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 调试版本（支持热重载需要配合 air）
RUN go build -gcflags="all=-N -l" -o main_debug ./cmd/main.go

# ========== 构建阶段 ==========
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 生产版本
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o main \
    ./cmd/main.go

# ========== 运行阶段 ==========
FROM alpine:3.19 AS production

# 安装必要的运行时库
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /home/appuser

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./main"]
```

### 1.6 安全最佳实践

```dockerfile
# 1. 使用特定版本标签，不要使用 latest
FROM golang:1.21.5-alpine3.19

# 2. 使用非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 3. 只复制必要文件
COPY --chown=appuser:appgroup ./bin/myapp /home/appuser/

# 4. 使用健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 5. 使用只读文件系统（可选）
# VOLUME ["/data"]
```

---

## 2. Docker Compose

### 2.1 docker-compose.yml 示例

```yaml
version: '3.8'

services:
  # 应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: myapp
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - DB_HOST=db
      - REDIS_HOST=redis
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./config.yaml:/home/appuser/config.yaml:ro
      - app_data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  # PostgreSQL 数据库
  db:
    image: postgres:15-alpine
    container_name: myapp-db
    environment:
      POSTGRES_USER: appuser
      POSTGRES_PASSWORD: securepassword
      POSTGRES_DB: myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    ports:
      - "5432:5432"
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U appuser -d myapp"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: myapp-redis
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - app-network
    restart: unless-stopped

  # Nginx 反向代理
  nginx:
    image: nginx:1.25-alpine
    container_name: myapp-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
  app_data:
```

### 2.2 docker-compose.override.yml（开发环境）

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: myapp-dev
    environment:
      - ENV=development
      - DB_HOST=localhost
      - REDIS_HOST=localhost
    volumes:
      - .:/app
      - go_cache:/go/pkg/mod
    ports:
      - "8080:8080"
      - "40000:40000" # Delve 调试端口
    command: ["air", "-b", "0.0.0.0"]

  # 开发时使用本地数据库
  db:
    ports:
      - "5432:5432"

  redis:
    ports:
      - "6379:6379"

volumes:
  go_cache:
```

### 2.3 常用 Compose 命令

```bash
# 启动服务（后台）
docker-compose up -d

# 启动服务并构建镜像
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v

# 重启服务
docker-compose restart

# 进入服务终端
docker-compose exec app /bin/sh

# 查看服务资源使用
docker-compose top

# 扩展服务
docker-compose up -d --scale app=3
```

---

## 3. Kubernetes 编排

### 3.1 Kubernetes 架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes 集群                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Control Plane                    │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐ │   │
│  │  │   API   │  │  etcd   │  │Scheduler│  │Controller│ │   │
│  │  │ Server  │  │         │  │         │  │ Manager  │ │   │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Data Plane                       │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐              │   │
│  │  │ Node 1  │  │ Node 2  │  │ Node 3  │              │   │
│  │  │ ┌─────┐ │  │ ┌─────┐ │  │ ┌─────┐ │              │   │
│  │  │ │Kubelet│ │  │ │Kubelet│ │  │ │Kubelet│ │              │   │
│  │  │ └─────┘ │  │ └─────┘ │  │ └─────┘ │              │   │
│  │  │ ┌─────┐ │  │ ┌─────┐ │  │ ┌─────┐ │              │   │
│  │  │ │ Pod  │ │  │ │ Pod  │ │  │ │ Pod  │ │              │   │
│  │  │ └─────┘ │  │ └─────┘ │  │ └─────┘ │              │   │
│  │  └─────────┘  └─────────┘  └─────────┘              │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: myapp
        version: v1
    spec:
      # 服务账户
      serviceAccountName: myapp-sa
      
      # 安全上下文
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      
      # 容器定义
      containers:
        - name: myapp
          image: myregistry.com/myapp:v1.0.0
          imagePullPolicy: Always
          
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          
          # 环境变量
          env:
            - name: ENV
              value: "production"
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: myapp-config
                  key: db_host
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: myapp-secrets
                  key: db_password
          
          # 资源配置
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          
          # 健康检查
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 3
            failureThreshold: 3
          
          readinessProbe:
            httpGet:
              path: /ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          
          # 挂载卷
          volumeMounts:
            - name: config
              mountPath: /home/appuser/config.yaml
              subPath: config.yaml
            - name: tmp
              mountPath: /tmp
      
      # 数据卷
      volumes:
        - name: config
          configMap:
            name: myapp-config
        - name: tmp
          emptyDir: {}
```

### 3.3 Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  type: ClusterIP  # ClusterIP, NodePort, LoadBalancer, ExternalName
  selector:
    app: myapp
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 8080
  
  # 会话亲和性
  sessionAffinity: None
  
  # 发布外部端口
  # externalTrafficPolicy: Cluster
```

### 3.4 Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "30"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "60"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - myapp.example.com
      secretName: myapp-tls
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: myapp
                port:
                  number: 80
          - path: /static
            pathType: Prefix
            backend:
              service:
                name: static-content
                port:
                  number: 80
```

### 3.5 ConfigMap 和 Secret

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  ENV: "production"
  DB_HOST: "myapp-db"
  DB_PORT: "5432"
  DB_NAME: "myapp"
  REDIS_HOST: "myapp-redis"
  LOG_LEVEL: "info"
```

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secrets
type: Opaque
data:
  # base64 编码
  DB_USERNAME: YXBwdXNlcg==
  DB_PASSWORD: c2VjdXJlcGFzc3dvcmQ=
  API_KEY: YXBpLWtleS1oZXJl
```

### 3.6 Horizontal Pod Autoscaler

```yaml
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
```

### 3.7 常用 kubectl 命令

```bash
# ========== 部署管理 ==========

# 应用配置
kubectl apply -f k8s/

# 查看部署
kubectl get deployment myapp

# 查看 Pod
kubectl get pods -l app=myapp

# 查看 Pod 日志
kubectl logs -f myapp-7d8f9c4d8-xk2mn

# 进入 Pod
kubectl exec -it myapp-7d8f9c4d8-xk2mn -- /bin/sh

# 缩放 Deployment
kubectl scale deployment myapp --replicas=5

# 查看部署状态
kubectl rollout status deployment myapp

# 回滚部署
kubectl rollout undo deployment myapp

# 查看部署历史
kubectl rollout history deployment myapp

# 回滚到指定版本
kubectl rollout undo deployment myapp --to-revision=2

# ========== 服务管理 ==========

# 查看 Service
kubectl get service myapp

# 端口转发
kubectl port-forward service/myapp 8080:80

# ========== 配置管理 ==========

# 应用 ConfigMap
kubectl apply -f k8s/configmap.yaml

# 查看 ConfigMap
kubectl get configmap myapp-config

# 应用 Secret
kubectl apply -f k8s/secret.yaml

# ========== 调试 ==========

# 查看 Pod 详情
kubectl describe pod myapp-7d8f9c4d8-xk2mn

# 查看资源使用
kubectl top pod -l app=myapp

# 应用 YAML
kubectl apply -f k8s/deployment.yaml

# 删除资源
kubectl delete -f k8s/
```

---

## 4. CI/CD 流水线

### 4.1 GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata for Docker
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=raw,value={{sha}}
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to Kubernetes
        uses: azure/k8s-set-context@v4
        with:
          kubeconfig: ${{ secrets.KUBE_CONFIG }}

      - name: Deploy to cluster
        run: |
          kubectl set image deployment/myapp myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} -n myapp
          kubectl rollout status deployment/myapp -n myapp
```

### 4.2 GitLab CI/CD

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: ""
  IMAGE_TAG: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

test:
  stage: test
  image: golang:1.21-alpine
  script:
    - go test -v -race ./...
  coverage: '/coverage: \d+\.\d+%/'
  artifacts:
    reports:
      junit: report.xml
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

build:
  stage: build
  image: docker:24-dind
  services:
    - docker:24-dind
  script:
    - docker build -t $IMAGE_TAG .
    - docker push $IMAGE_TAG
  only:
    - main

deploy:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/myapp myapp=$IMAGE_TAG -n myapp
    - kubectl rollout status deployment/myapp -n myapp
  only:
    - main
  environment:
    name: production
    url: https://myapp.example.com
```

---

## 5. 监控与日志

### 5.1 Prometheus 配置

```yaml
# prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: []

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'myapp'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
```

### 5.2 Go 应用指标

```go
package main

import (
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    httpRequestsInFlight = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "http_requests_in_flight",
            Help: "Current number of HTTP requests being processed",
        },
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(httpRequestsInFlight)
}

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        httpRequestsInFlight.Inc()
        defer httpRequestsInFlight.Dec()

        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            http.StatusText(http.StatusOK),
        ).Inc()

        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    })
}

func main() {
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8081", nil)
}
```

### 5.3 日志配置

```go
package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level string) *zap.Logger {
    config := zap.NewProductionConfig()
    
    switch level {
    case "debug":
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.DebugLevelColor
    case "info":
        config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
    case "error":
        config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
    }

    logger, _ := config.Build()
    return logger
}

func main() {
    logger := NewLogger("production")
    defer logger.Sync()

    logger.Info("Application started",
        zap.String("version", "1.0.0"),
        zap.Int("port", 8080),
    )

    logger.Error("Failed to connect to database",
        zap.Error(err),
        zap.String("host", "localhost"),
    )
}
```

---

## 最佳实践

### 1. Docker 最佳实践

- 使用多阶段构建减小镜像体积
- 使用非 root 用户运行
- 固定基础镜像版本
- 利用构建缓存
- 扫描镜像漏洞

### 2. Kubernetes 最佳实践

- 使用命名空间隔离环境
- 设置资源请求和限制
- 使用探针保证可用性
- 使用 ConfigMap 管理配置
- 使用 HPA 自动扩缩容

### 3. CI/CD 最佳实践

- 自动化所有测试
- 使用阶段门禁
- 实现渐进式发布
- 保留回滚能力
- 监控部署质量

### 4. 监控最佳实践

- 定义清晰的 SLO
- 收集关键指标
- 设置告警阈值
- 保留足够的历史数据
- 可视化关键指标
