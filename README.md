# 分布式链路追踪服务（distributed-tracing）

一个基于 Go 标准库实现的轻量级分布式链路追踪后端服务，用于采集、存储与检索服务调用轨迹（Trace/Span），并支持采样策略、慢轨迹检测与全局统计。

## 技术栈

- 纯 Go 标准库（`net/http` + 标准库），零第三方依赖
- 内存存储（`sync.RWMutex` 保证并发安全）
- 标准分层架构：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`

## 运行

```bash
go run ./cmd/server
# 或
go build -o tracing-server ./cmd/server && ./tracing-server
```

默认监听 `:8080`，可通过环境变量配置：

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| `PORT` | `8080` | 监听端口 |
| `ADDR` | `:8080` | 完整监听地址（优先级高于 PORT） |
| `MAX_PAGE_SIZE` | `100` | 分页最大页大小 |
| `DEFAULT_SAMPLE` | `100` | 未命中采样规则时的默认采样率（0-100） |
| `LOG_LEVEL` | `info` | 日志级别：debug/info/warn/error |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误码：400 校验失败、404 不存在、409 冲突、500 服务器错误。

### 服务 Service

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/services` | 创建服务 |
| GET | `/api/services` | 列表（status/keyword 筛选 + 分页） |
| GET | `/api/services/{id}` | 详情 |
| PUT | `/api/services/{id}` | 更新 |
| DELETE | `/api/services/{id}` | 删除 |

### 轨迹 Trace

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/traces` | 创建轨迹（service_id/operation） |
| GET | `/api/traces` | 列表（service_id/status/min_duration/keyword + 分页） |
| GET | `/api/traces/{id}` | 详情 |
| PUT | `/api/traces/{id}` | 更新 operation |
| POST | `/api/traces/{id}/finish` | 结束轨迹（status/duration_ms，状态机校验） |
| DELETE | `/api/traces/{id}` | 删除 |

### 跨度 Span

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/spans` | 创建跨度（校验 trace 与 service 存在） |
| GET | `/api/spans` | 列表（trace_id/service_id/status/operation + 分页） |
| GET | `/api/spans/{id}` | 详情 |
| PUT | `/api/spans/{id}` | 更新 operation/tags |
| POST | `/api/spans/{id}/finish` | 结束跨度（状态机校验） |
| DELETE | `/api/spans/{id}` | 删除 |

### 标注 Annotation

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/annotations` | 创建标注（校验 span 存在且 trace 一致） |
| GET | `/api/annotations` | 列表（span_id/trace_id/keyword + 分页） |
| GET | `/api/annotations/{id}` | 详情 |
| PUT | `/api/annotations/{id}` | 更新 name/message |
| DELETE | `/api/annotations/{id}` | 删除 |

### 采样规则 SamplingRule

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/sampling-rules` | 创建采样规则 |
| GET | `/api/sampling-rules` | 列表（status/keyword + 分页） |
| GET | `/api/sampling-rules/{id}` | 详情 |
| PUT | `/api/sampling-rules/{id}` | 更新 |
| DELETE | `/api/sampling-rules/{id}` | 删除 |
| GET | `/api/sampling/evaluate?service_name=x` | 采样决策（返回命中规则与是否采样） |

### 慢轨迹规则 SlowTraceRule

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/slow-rules` | 创建慢轨迹规则 |
| GET | `/api/slow-rules` | 列表（service_id/status + 分页） |
| GET | `/api/slow-rules/{id}` | 详情 |
| PUT | `/api/slow-rules/{id}` | 更新 |
| DELETE | `/api/slow-rules/{id}` | 删除 |
| GET | `/api/slow-rules/detect` | 检测命中的慢轨迹 |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 全局观测统计 |

## 核心实体

1. **Service**：被观测服务节点（name/host/port/status）。
2. **Trace**：一次请求的完整调用轨迹，状态机 `running → completed/failed`。
3. **Span**：轨迹中的单个操作片段，状态机 `running → completed/failed`。
4. **Annotation**：挂在 Span 上的关键事件标注。
5. **SamplingRule**：采样策略，支持 `*` 通配服务模式，按哈希比例采样。
6. **SlowTraceRule**：慢轨迹检测规则，按耗时阈值匹配。

## 测试

```bash
go test ./...
```
