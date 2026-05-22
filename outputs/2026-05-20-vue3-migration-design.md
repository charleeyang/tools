# Vue 3 前端迁移设计方案

**日期**: 2026-05-20
**项目**: 袁夫稻田智慧园区平台 (yuanfu-ricefield/tools)
**目标**: 将 `web/` 下原生 JS 单页应用迁移至 Vue 3 框架

## 1. 决策摘要

| 决策项 | 选择 |
|--------|------|
| 迁移范围 | 4 个平台视图一次性全部迁移 |
| 框架 | Vue 3 SPA (Vite + Vue Router + Pinia + TypeScript) |
| PC 组件库 | Ant Design Vue 4.x |
| 移动端组件库 | Vant 4.x |
| 前后端关系 | 分离部署：Go 只做 API，Vue 构建产物独立部署 |
| 移动端预览 | PC 后台桌面布局，三端小程序保持手机框模拟 |
| 项目位置 | `web/` 原地替换旧前端 |

## 2. 项目结构

```
web/
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── public/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/index.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── permission.ts
│   │   └── view.ts
│   ├── api/
│   │   ├── client.ts
│   │   └── modules/
│   │       ├── auth.ts
│   │       ├── parks.ts
│   │       ├── shops.ts
│   │       ├── orders.ts
│   │       ├── products.ts
│   │       ├── finance.ts
│   │       ├── users.ts
│   │       ├── employees.ts
│   │       ├── gates.ts
│   │       ├── activities.ts
│   │       └── platform.ts
│   ├── layouts/
│   │   ├── PcLayout.vue
│   │   ├── MobileLayout.vue
│   │   └── BlankLayout.vue
│   ├── views/
│   │   ├── login/
│   │   ├── pc-admin/
│   │   │   ├── overview/
│   │   │   ├── parks/
│   │   │   ├── shops/
│   │   │   ├── orders/
│   │   │   ├── products/
│   │   │   ├── users/
│   │   │   ├── finance/
│   │   │   ├── employees/
│   │   │   ├── gates/
│   │   │   ├── activities/
│   │   │   ├── platform/
│   │   │   └── settings/
│   │   ├── mp-admin/
│   │   ├── mp-merchant/
│   │   └── mp-client/
│   ├── components/
│   │   ├── common/
│   │   ├── pc/
│   │   └── mobile/
│   └── utils/
│       ├── auth.ts
│       └── menu.ts
```

## 3. 路由设计

```
/ → redirect → /pc-admin/overview
/login                    layout: BlankLayout

/pc-admin/*               layout: PcLayout
  ├── overview
  ├── parks
  ├── shops
  ├── orders
  ├── products
  ├── users
  ├── finance
  ├── employees
  ├── gates
  ├── activities
  ├── platform
  └── settings

/mp-admin/*               layout: MobileLayout
  ├── dashboard
  ├── shops
  ├── scan
  ├── message
  └── me

/mp-merchant/*            layout: MobileLayout
  ├── dashboard
  ├── orders
  ├── verify
  ├── finance
  └── me

/mp-client/*              layout: MobileLayout
  ├── home
  ├── memberqr
  ├── cart
  └── me
```

**路由守卫 (beforeEach)**:
1. 无 token → `/login`
2. 有 token 但路由需要特定权限 → 检查 permissionStore
3. 根据路由 meta.layout 自动匹配布局组件

**平台切换**: 顶栏按钮点击 → `router.push('/pc-admin')` 或 `window.open` 保持多 Tab 独立状态

## 4. 状态管理 (Pinia)

### authStore
- state: `token`, `user`, `permissions`
- actions: `login()`, `logout()`, `fetchMe()`
- 持久化: `localStorage`

### permissionStore
- state: `permissions: Record<string, {read: boolean, write: boolean}>`
- getters: `canRead(module)`, `canWrite(module)`

### viewStore
- state: `platform`, `view`
- actions: `switchPlatform()`, `switchView()`

## 5. 组件层级

```
App.vue
├── PlatformSwitcher.vue
├── <router-view> (layout)
│   ├── PcLayout.vue
│   │   ├── PcSider.vue
│   │   │   ├── ViewSwitch.vue
│   │   │   ├── PcNav.vue
│   │   │   └── UserFooter.vue
│   │   ├── PcHeader.vue
│   │   │   ├── Breadcrumb.vue
│   │   │   └── HeaderActions.vue
│   │   └── <router-view> (page)
│   ├── MobileLayout.vue
│   │   ├── PhoneFrame.vue
│   │   │   ├── StatusBar.vue
│   │   │   ├── MiniAppHeader.vue
│   │   │   ├── <router-view>
│   │   │   └── TabBar.vue
│   │   └── MemberQrOverlay.vue
│   └── BlankLayout.vue
└── GlobalModal.vue / GlobalToast.vue
```

## 6. 数据流

```
组件 → API 模块 → axios 实例 → Go :8080
                    ├── request 拦截器: 注入 Authorization
                    └── response 拦截器: 401 → authStore.logout() → /login
```

- 组件内部数据: `ref` / `reactive`
- 跨组件共享: Pinia stores
- 持久化数据: localStorage (token, user, permissions)

## 7. 迁移顺序

| Phase | 内容 | 说明 |
|-------|------|------|
| 1 | 基础设施 | Vite 脚手架、Pinia stores、axios client、Router |
| 2 | 布局框架 | App.vue、PcLayout、MobileLayout |
| 3 | PC 管理后台 | 12 个模块页面 |
| 4 | 移动端三端 | mp-admin/merchant/client |
| 5 | 收尾 | 登录页、全局组件、构建配置、删除旧文件 |

## 8. 约束

- Go 后端不修改，API 接口不变
- 旧版 `api.js` 的调用模式（base/token/request/get/post/put/del）在 axios client 中保持一致
- 开发时 Vite dev server 端口 5173，proxy `/api` 到 `localhost:8080`
- 构建产物 `web/dist/`，供 Go 或 CDN 使用
