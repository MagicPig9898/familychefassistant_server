# familychefassistant_server

家庭厨师助手后端服务（Go + Gin）。

## 当前引入的包

### 业务直接使用

```bash
go get github.com/gin-gonic/gin@v1.10.1
go get github.com/MagicPig9898/easy_db@v1.0.0-beta.1
```

### 标准库

- `log`

### 项目内核心包

- `config/db_config`：数据库初始化、获取、关闭
- `config/router_config`：统一路由注册入口
- `config/router_config/user_router`：user 模块路由
- `config/router_config/router_utils`：路由层通用工具方法
- `services/user_service`：user 业务服务层

## 启动方式

```bash
go run main.go
```

默认监听地址：`http://localhost:8080`

## 路由清单

基础前缀：`/api/v1/users`

### 1) 健康检查

- 方法：`GET`
- 路径：`/api/v1/users/healthz`
- 响应：`ok`
- 调用示例：

```bash
curl -X GET "http://localhost:8080/api/v1/users/healthz"
```

### 2) 查询用户信息

- 方法：`GET`
- 路径：`/api/v1/users/info`
- Query 参数：
  - `id`（必填，int64）
- 成功响应示例：

```json
{"id":1,"name":"Tom"}
```

- 失败响应示例：
  - 缺少参数：`missing id`
  - 参数错误：`invalid id`
- 调用示例（成功）：

```bash
curl -X GET "http://localhost:8080/api/v1/users/info?id=1"
```

- 调用示例（缺少参数）：

```bash
curl -X GET "http://localhost:8080/api/v1/users/info"
```

- 调用示例（参数非法）：

```bash
curl -X GET "http://localhost:8080/api/v1/users/info?id=abc"
```
