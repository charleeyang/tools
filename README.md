# 袁夫稻田智慧园区综合管理平台

依据《PRD V3.0》《开发实施计划》与原型图实现的全栈系统。后端为 Go + 内嵌
SQLite（纯 Go，零外部依赖），前端为原型四端（PC 管理后台 + 管理端/商户端/客户端
小程序），由同一可执行文件提供服务并对接真实 API。

```
┌──────────────────────────────────────────────┐
│  浏览器 http://localhost:8080                   │
│  ┌──────────────┐  登录门禁 + 实时数据          │
│  │ web/ 前端原型 │  (api.js / bootstrap.js)      │
│  └──────┬───────┘                              │
│         │ fetch /api/*  (JWT)                   │
│  ┌──────▼───────────────────────────────────┐ │
│  │ server/ Go (net/http)                      │ │
│  │  认证 · RBAC 三级权限 · 视角数据过滤        │ │
│  │  园区/店铺/客户/订单/财务/核销/活动/系统 …  │ │
│  └──────┬───────────────────────────────────┘ │
│         │                                      │
│   ┌─────▼─────┐  27 张表 (PRD 第十章)           │
│   │  SQLite   │  首次启动自动建表 + 灌种子数据    │
│   └───────────┘                                │
└──────────────────────────────────────────────┘
```

## 快速启动（上线本地）

需要 Go 1.22+（已在 1.24 验证）。

```bash
# 方式一：Makefile
make run            # → http://localhost:8080

# 方式二：直接运行
cd server && go run . -addr :8080 -web ../web -db yfsc.db

# 方式三：编译为单文件
make build && ./bin/yfsc-server -web web

# 方式四：Docker
docker build -t yfsc . && docker run -p 8080:8080 -v yfsc-data:/app/data yfsc

# 方式五：发布版自包含单文件（前端已内嵌，无需 Go/Node/web 目录）
make release            # 生成 dist/ 下各平台单文件
./dist/yfsc-server-linux-amd64          # Linux
./dist/yfsc-server-macos-arm64          # macOS (Apple Silicon)
# Windows: 双击 dist\yfsc-server-windows-amd64.exe
```

> 发布版为单个可执行文件，内置前端与 SQLite，运行后浏览器打开
> `http://localhost:8080` 即可。首次运行在同目录生成 `yfsc.db`。

打开 `http://localhost:8080`，使用下方账号登录。首次启动会自动建表并写入种子数据
（`yfsc.db`）。删除该文件即可重置（`make fresh`）。

## 演示账号

| 用户名 | 密码 | 角色 | 数据视角 |
|--------|------|------|----------|
| `admin` | `admin123` | 平台超级管理员 | 全平台 |
| `park-hm` | `123456` | 园区管理员 | 黄梅园区 |
| `park-wh` | `123456` | 园区管理员 | 武汉园区 |
| `hcct-mgr` | `123456` | 商户管理员 | 火车餐厅 |
| `hcct-cashier` | `123456` | 商户收银员 | 火车餐厅（不可提现）|

左侧「当前视角」下拉切换时会以对应角色重新登录并加载受限数据。

## 功能与验证

- **后端 API 冒烟测试**（43 项）：覆盖认证、三级 RBAC、视角过滤、核销/退款/提现
  审核、美团/抖音验券全流程、各模块 CRUD。

  ```bash
  make run &              # 或后台启动
  make smoke             # bash server/smoke_test.sh
  ```

- **前端端到端测试**（20 项，JSDOM）：登录门禁、实时数据注入、各页面以真实数据库
  数据渲染。脚本见 `server/`（开发期使用 jsdom 驱动）。

## API 概览（统一返回 `{code,data,message}`）

| 模块 | 端点 |
|------|------|
| 认证 | `POST /api/auth/login`、`GET /api/auth/me` |
| 统计 | `GET /api/statistics/overview` |
| 财务 | `GET /api/finance/summary`（园区/平台收入聚合：线上/美团/抖音/余额/冻结）|
| 园区 | `GET/POST /api/parks`、`PUT/DELETE /api/parks/{id}` |
| 店铺 | `GET/POST /api/shops`、`PUT /api/shops/{id}`、`/api/shop-types` |
| 客户 | `GET/POST/PUT/DELETE /api/users`、`/api/users/stats`、`/api/recharges` |
| 商品 | `GET/POST/PUT/DELETE /api/products` |
| 消费 | `GET /api/orders`、`POST /api/orders/{no}/verify`、`/refund` |
| 退款 | `GET /api/refunds`、`POST /api/refunds/{id}/approve`、`/reject` |
| 财务 | `GET /api/withdraws`、`POST /api/withdraws/{id}/approve`、`/reject` |
| 员工 | `GET/POST /api/employees`、`/api/positions` |
| 闸机 | `GET/POST /api/gates`、`PUT /api/gates/{id}/toggle`、`DELETE /api/gates/{id}`、`/api/faces`、`/api/entry-records` |
| 活动 | `GET/POST /api/activities`、`POST /api/activities/{id}/audit`、`/end` |
| 第三方 | `/api/platform/verify-records`、`/settlements`、`/store-configs`、`verify/prepare\|execute\|{id}/revoke` |
| 系统 | `GET/PUT /api/system/mini-programs`、`/api/system/logs`、`GET/PUT /api/permissions`（角色权限矩阵实时调整）|

## 目录结构

```
server/        Go 后端 (单包多文件)
  schema.sql   27 张表 DDL (SQLite)
  seed.go      角色权限 + 演示数据
  auth.go      JWT + bcrypt
  middleware.go  认证 / RBAC / 视角 / CORS
  handlers*.go   各模块业务逻辑
  smoke_test.sh  API 冒烟测试
web/           前端原型 (四端) + api.js / bootstrap.js 实时对接层
docs/          PRD V3.0 + 开发实施计划
```

## 实现范围说明

实施计划评估为 4–5 人 16–20 周的工程量（Go 微服务 + Vue + 3 个 uni-app 小程序 +
数据迁移 + 真实三方对接）。本仓库在保持架构与 API 契约一致的前提下做了务实取舍，
交付一个**可本地一键运行、端到端可验证**的系统：

- 数据库：以 SQLite 落地 PRD 全部 27 张表（生产可平滑迁移至 MySQL 8）。
- 后端：以标准库 `net/http` 实现计划中的认证 / RBAC / 视角中间件与全部模块端点
  （等价于 go-zero 自动生成的 handler+logic 分层）。
- 前端：复用原型四端 UI，新增登录门禁与实时数据层对接真实后端。园区总览、店铺/
  客户/订单/充值/退款/员工/活动、美团抖音核销、园区账户财务汇总、提现审核、闸机
  增删启停均跑实时数据与真实写操作；仅对账单等少量页仍用原型内置数据。
- 三方对接（美团/抖音）：实现验券 prepare/execute/revoke 状态机与门店绑定校验，
  外部 HTTP 调用为可替换的模拟实现（签名算法见实施计划 G 章）。
