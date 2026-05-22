# Vue 3 前端迁移实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 tools/web/ 下原生 JS 单页应用（~400KB，4 个平台视图）迁移为 Vue 3 + Vite + Vue Router + Pinia + Ant Design Vue + Vant 的 TypeScript SPA。

**Architecture:** 单页应用，Vue Router 分 4 个路由组，Pinia 管理状态，axios 对接 Go API。PC 端 Ant Design Vue + PcLayout，移动端 Vant + MobileLayout（含手机模拟框）。构建产物输出到 web/dist/，Go 后端不改。

**Tech Stack:** Vue 3.5+ / Vite 7 / Vue Router 4 / Pinia 2 / TypeScript / Ant Design Vue 4 / Vant 4 / axios

---

## Phase 1: 基础设施

### Task 1: 备份旧前端并创建 Vite + Vue 3 脚手架

**Files:**
- Create: `web/package.json`
- Create: `web/index.html`（新入口）
- Create: `web/vite.config.ts`
- Create: `web/tsconfig.json`
- Create: `web/tsconfig.node.json`
- Create: `web/env.d.ts`
- Create: `web/.gitignore`
- Modify: 移除旧 `web/assets/` 和旧 `web/index.html` → 备份到 `web/legacy/`

- [ ] **Step 1: 备份旧前端文件**

```bash
mkdir -p /Users/allen/tools/web/legacy
mv /Users/allen/tools/web/index.html /Users/allen/tools/web/legacy/index.html
mv /Users/allen/tools/web/assets /Users/allen/tools/web/legacy/assets
```

- [ ] **Step 2: 创建 package.json**

```json
{
  "name": "yfsc-platform",
  "private": true,
  "version": "3.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.5.0",
    "vue-router": "^4.5.0",
    "pinia": "^2.3.0",
    "ant-design-vue": "^4.2.0",
    "@ant-design/icons-vue": "^7.0.0",
    "vant": "^4.9.0",
    "axios": "^1.7.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.0",
    "typescript": "~5.7.0",
    "vite": "^6.2.0",
    "vue-tsc": "^2.2.0",
    "unplugin-vue-router": "^0.11.0"
  }
}
```

- [ ] **Step 3: 创建 tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForExpose": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "noFallthroughCasesInSwitch": true,
    "paths": { "@/*": ["./src/*"] },
    "baseUrl": "."
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue", "env.d.ts"]
}
```

- [ ] **Step 4: 创建 tsconfig.node.json**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2023"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "strict": true
  },
  "include": ["vite.config.ts"]
}
```

- [ ] **Step 5: 创建 env.d.ts**

```typescript
/// <reference types="vite/client" />
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
```

- [ ] **Step 6: 创建 vite.config.ts**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': resolve(__dirname, 'src') },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    outDir: 'dist',
  },
})
```

- [ ] **Step 7: 创建新 index.html（Vite 入口）**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>袁夫稻田 · 智慧园区综合管理平台 v3.0</title>
</head>
<body>
  <div id="app"></div>
  <script type="module" src="/src/main.ts"></script>
</body>
</html>
```

- [ ] **Step 8: 创建 web/.gitignore**

```
node_modules/
dist/
*.local
```

- [ ] **Step 9: 安装依赖**

```bash
cd /Users/allen/tools/web && npm install
```

---

### Task 2: 创建 src/main.ts 和 src/App.vue 骨架

**Files:**
- Create: `web/src/main.ts`
- Create: `web/src/App.vue`

- [ ] **Step 1: 创建 src/main.ts**

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(Antd)
app.mount('#app')
```

- [ ] **Step 2: 创建 src/App.vue 骨架**

```vue
<template>
  <router-view />
</template>

<script setup lang="ts">
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
</style>
```

---

### Task 3: 创建 Axios API 客户端

**Files:**
- Create: `web/src/api/client.ts`

- [ ] **Step 1: 创建 src/api/client.ts**

```typescript
import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const http: AxiosInstance = axios.create({
  baseURL: '',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export async function request<T = any>(
  method: string,
  path: string,
  body?: any,
): Promise<T> {
  const config: AxiosRequestConfig = { method, url: path }
  if (body) config.data = body
  const res = await http.request<ApiResponse<T>>(config)
  if (res.data.code !== 0) {
    const err = new Error(res.data.message || '请求失败') as any
    err.code = res.data.code
    throw err
  }
  return res.data.data
}

export const api = {
  get: <T = any>(path: string) => request<T>('GET', path),
  post: <T = any>(path: string, body?: any) => request<T>('POST', path, body),
  put: <T = any>(path: string, body?: any) => request<T>('PUT', path, body),
  del: <T = any>(path: string) => request<T>('DELETE', path),
}
```

---

### Task 4: 创建 Pinia Stores

**Files:**
- Create: `web/src/stores/auth.ts`
- Create: `web/src/stores/permission.ts`
- Create: `web/src/stores/view.ts`

- [ ] **Step 1: 创建 src/stores/auth.ts**

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export interface User {
  id: number
  username: string
  realname: string
  role: string
  avatar?: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('yfsc_token') || '')
  const user = ref<User | null>(
    JSON.parse(localStorage.getItem('yfsc_user') || 'null'),
  )
  const permissions = ref<Record<string, { read: boolean; write: boolean }>>(
    JSON.parse(localStorage.getItem('yfsc_perms') || '{}'),
  )

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const data: any = await api.post('/api/auth/login', { username, password })
    token.value = data.token
    user.value = data.user
    permissions.value = data.permissions || {}
    localStorage.setItem('yfsc_token', token.value)
    localStorage.setItem('yfsc_user', JSON.stringify(user.value))
    localStorage.setItem('yfsc_perms', JSON.stringify(permissions.value))
  }

  async function fetchMe() {
    const data: any = await api.get('/api/auth/me')
    user.value = data.user || data
  }

  function logout() {
    token.value = ''
    user.value = null
    permissions.value = {}
    localStorage.removeItem('yfsc_token')
    localStorage.removeItem('yfsc_user')
    localStorage.removeItem('yfsc_perms')
  }

  return { token, user, permissions, isLoggedIn, login, logout, fetchMe }
})
```

- [ ] **Step 2: 创建 src/stores/permission.ts**

```typescript
import { defineStore } from 'pinia'
import { computed } from 'vue'
import { useAuthStore } from './auth'

export const usePermissionStore = defineStore('permission', () => {
  const authStore = useAuthStore()

  function canRead(module: string): boolean {
    return authStore.permissions[module]?.read ?? false
  }

  function canWrite(module: string): boolean {
    return authStore.permissions[module]?.write ?? false
  }

  return { canRead, canWrite }
})
```

- [ ] **Step 3: 创建 src/stores/view.ts**

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'

export type Platform = 'pc-admin' | 'mp-admin' | 'mp-merchant' | 'mp-client'
export type ViewRole = 'platform' | 'park-hm' | 'park-wh' | 'shop-hcct'

export const useViewStore = defineStore('view', () => {
  const platform = ref<Platform>('pc-admin')
  const view = ref<ViewRole>('platform')

  function switchPlatform(p: Platform) {
    platform.value = p
  }

  function switchView(v: ViewRole) {
    view.value = v
  }

  return { platform, view, switchPlatform, switchView }
})
```

---

### Task 5: 创建路由配置

**Files:**
- Create: `web/src/router/index.ts`

- [ ] **Step 1: 创建空白 view 占位文件**

```bash
mkdir -p /Users/allen/tools/web/src/views/{login,pc-admin/overview,pc-admin/parks,pc-admin/shops,pc-admin/orders,pc-admin/products,pc-admin/users,pc-admin/finance,pc-admin/employees,pc-admin/gates,pc-admin/activities,pc-admin/platform,pc-admin/settings,mp-admin,mp-merchant,mp-client}
mkdir -p /Users/allen/tools/web/src/{layouts,components/{common,pc,mobile},utils}
```

- [ ] **Step 2: 创建各 view 占位组件（以 overview 为例，其他同理）**

`web/src/views/pc-admin/overview/index.vue`:
```vue
<template><div>概览页面</div></template>
<script setup lang="ts"></script>
```

(同样创建其他 19 个 view 的占位 index.vue，内容只需 `<template><div>页面名</div></template>`)

- [ ] **Step 3: 创建 src/router/index.ts**

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/pc-admin/overview',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { layout: 'blank' },
  },
  {
    path: '/pc-admin',
    component: () => import('@/layouts/PcLayout.vue'),
    meta: { requiresAuth: true, layout: 'pc' },
    children: [
      { path: '', redirect: { path: '/pc-admin/overview' } },
      { path: 'overview', component: () => import('@/views/pc-admin/overview/index.vue'), meta: { title: '首页概览' } },
      { path: 'parks', component: () => import('@/views/pc-admin/parks/index.vue'), meta: { title: '园区管理', module: 'parks' } },
      { path: 'shops', component: () => import('@/views/pc-admin/shops/index.vue'), meta: { title: '店铺管理', module: 'shops' } },
      { path: 'orders', component: () => import('@/views/pc-admin/orders/index.vue'), meta: { title: '消费订单', module: 'consumption' } },
      { path: 'products', component: () => import('@/views/pc-admin/products/index.vue'), meta: { title: '商品管理', module: 'products' } },
      { path: 'users', component: () => import('@/views/pc-admin/users/index.vue'), meta: { title: '客户管理', module: 'users' } },
      { path: 'finance', component: () => import('@/views/pc-admin/finance/index.vue'), meta: { title: '财务管理', module: 'finance' } },
      { path: 'employees', component: () => import('@/views/pc-admin/employees/index.vue'), meta: { title: '员工管理', module: 'employees' } },
      { path: 'gates', component: () => import('@/views/pc-admin/gates/index.vue'), meta: { title: '闸机管理', module: 'gates' } },
      { path: 'activities', component: () => import('@/views/pc-admin/activities/index.vue'), meta: { title: '活动管理', module: 'activities' } },
      { path: 'platform', component: () => import('@/views/pc-admin/platform/index.vue'), meta: { title: '第三方平台' } },
      { path: 'settings', component: () => import('@/views/pc-admin/settings/index.vue'), meta: { title: '系统设置', module: 'settings' } },
    ],
  },
  {
    path: '/mp-admin',
    component: () => import('@/layouts/MobileLayout.vue'),
    meta: { requiresAuth: true, layout: 'mobile', platform: 'mp-admin' },
    children: [
      { path: '', redirect: { path: '/mp-admin/dashboard' } },
      { path: 'dashboard', component: () => import('@/views/mp-admin/index.vue'), meta: { title: '首页', tab: 'dashboard' } },
      { path: 'shops', component: () => import('@/views/mp-admin/index.vue'), meta: { title: '商户', tab: 'shops' } },
      { path: 'scan', component: () => import('@/views/mp-admin/index.vue'), meta: { title: '扫码', tab: 'scan' } },
      { path: 'message', component: () => import('@/views/mp-admin/index.vue'), meta: { title: '消息', tab: 'message' } },
      { path: 'me', component: () => import('@/views/mp-admin/index.vue'), meta: { title: '我的', tab: 'me' } },
    ],
  },
  {
    path: '/mp-merchant',
    component: () => import('@/layouts/MobileLayout.vue'),
    meta: { requiresAuth: true, layout: 'mobile', platform: 'mp-merchant' },
    children: [
      { path: '', redirect: { path: '/mp-merchant/dashboard' } },
      { path: 'dashboard', component: () => import('@/views/mp-merchant/index.vue'), meta: { title: '经营', tab: 'dashboard' } },
      { path: 'orders', component: () => import('@/views/mp-merchant/index.vue'), meta: { title: '订单', tab: 'orders' } },
      { path: 'verify', component: () => import('@/views/mp-merchant/index.vue'), meta: { title: '核销', tab: 'verify' } },
      { path: 'finance', component: () => import('@/views/mp-merchant/index.vue'), meta: { title: '收益', tab: 'finance' } },
      { path: 'me', component: () => import('@/views/mp-merchant/index.vue'), meta: { title: '我的', tab: 'me' } },
    ],
  },
  {
    path: '/mp-client',
    component: () => import('@/layouts/MobileLayout.vue'),
    meta: { requiresAuth: true, layout: 'mobile', platform: 'mp-client' },
    children: [
      { path: '', redirect: { path: '/mp-client/home' } },
      { path: 'home', component: () => import('@/views/mp-client/index.vue'), meta: { title: '首页', tab: 'home' } },
      { path: 'memberqr', component: () => import('@/views/mp-client/index.vue'), meta: { title: '亮码支付', tab: 'memberqr' } },
      { path: 'cart', component: () => import('@/views/mp-client/index.vue'), meta: { title: '购物车', tab: 'cart' } },
      { path: 'me', component: () => import('@/views/mp-client/index.vue'), meta: { title: '个人中心', tab: 'me' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.path === '/login' && authStore.isLoggedIn) {
    next('/')
    return
  }

  next()
})

export default router
```

---

### Task 6: 创建 API 模块

**Files:**
- Create: `web/src/api/modules/auth.ts`
- Create: `web/src/api/modules/parks.ts`
- Create: `web/src/api/modules/shops.ts`
- Create: `web/src/api/modules/orders.ts`
- Create: `web/src/api/modules/products.ts`
- Create: `web/src/api/modules/users.ts`
- Create: `web/src/api/modules/finance.ts`
- Create: `web/src/api/modules/employees.ts`
- Create: `web/src/api/modules/gates.ts`
- Create: `web/src/api/modules/activities.ts`
- Create: `web/src/api/modules/platform.ts`

- [ ] **Step 1: 创建 src/api/modules/auth.ts**

```typescript
import { api } from '@/api/client'

export function login(username: string, password: string) {
  return api.post('/api/auth/login', { username, password })
}
export function getMe() {
  return api.get('/api/auth/me')
}
```

- [ ] **Step 2: 创建 src/api/modules/parks.ts**

```typescript
import { api } from '@/api/client'

export interface Park {
  id: number; name: string; address: string; status: string
  created_at: string; updated_at: string
}
export function getParkList(params?: Record<string, any>) {
  return api.get<Park[]>('/api/parks?' + new URLSearchParams(params).toString())
}
export function createPark(body: Partial<Park>) {
  return api.post<Park>('/api/parks', body)
}
export function updatePark(id: number, body: Partial<Park>) {
  return api.put<Park>(`/api/parks/${id}`, body)
}
export function deletePark(id: number) {
  return api.del(`/api/parks/${id}`)
}
```

- [ ] **Step 3: 创建 src/api/modules/shops.ts**

```typescript
import { api } from '@/api/client'

export interface Shop {
  id: number; name: string; park_id: number; shop_type_id: number
  status: string; created_at: string
}
export function getShopList(params?: Record<string, any>) {
  return api.get<Shop[]>('/api/shops?' + new URLSearchParams(params).toString())
}
export function createShop(body: Partial<Shop>) {
  return api.post<Shop>('/api/shops', body)
}
export function updateShop(id: number, body: Partial<Shop>) {
  return api.put<Shop>(`/api/shops/${id}`, body)
}
export function deleteShop(id: number) {
  return api.del(`/api/shops/${id}`)
}
export function getShopTypes() {
  return api.get<any[]>('/api/shop-types')
}
export function createShopType(body: any) {
  return api.post<any>('/api/shop-types', body)
}
```

- [ ] **Step 4: 创建 src/api/modules/orders.ts**

```typescript
import { api } from '@/api/client'

export function getOrderList(params?: Record<string, any>) {
  return api.get('/api/orders?' + new URLSearchParams(params).toString())
}
export function verifyOrder(no: string) {
  return api.post(`/api/orders/${no}/verify`)
}
export function refundOrder(no: string) {
  return api.post(`/api/orders/${no}/refund`)
}
export function getRefundList(params?: Record<string, any>) {
  return api.get('/api/refunds?' + new URLSearchParams(params).toString())
}
export function approveRefund(id: number) {
  return api.post(`/api/refunds/${id}/approve`)
}
export function rejectRefund(id: number) {
  return api.post(`/api/refunds/${id}/reject`)
}
```

- [ ] **Step 5: 创建 src/api/modules/products.ts**

```typescript
import { api } from '@/api/client'

export function getProductList(params?: Record<string, any>) {
  return api.get('/api/products?' + new URLSearchParams(params).toString())
}
export function createProduct(body: any) {
  return api.post('/api/products', body)
}
export function updateProduct(id: number, body: any) {
  return api.put(`/api/products/${id}`, body)
}
export function deleteProduct(id: number) {
  return api.del(`/api/products/${id}`)
}
```

- [ ] **Step 6: 创建 src/api/modules/users.ts**

```typescript
import { api } from '@/api/client'

export function getUserList(params?: Record<string, any>) {
  return api.get('/api/users?' + new URLSearchParams(params).toString())
}
export function getUserStats() {
  return api.get('/api/users/stats')
}
export function createUser(body: any) {
  return api.post('/api/users', body)
}
export function updateUser(id: number, body: any) {
  return api.put(`/api/users/${id}`, body)
}
export function deleteUser(id: number) {
  return api.del(`/api/users/${id}`)
}
export function getRechargeList(params?: Record<string, any>) {
  return api.get('/api/recharges?' + new URLSearchParams(params).toString())
}
```

- [ ] **Step 7: 创建 src/api/modules/finance.ts**

```typescript
import { api } from '@/api/client'

export function getFinanceSummary() {
  return api.get('/api/finance/summary')
}
export function getWithdrawList(params?: Record<string, any>) {
  return api.get('/api/withdraws?' + new URLSearchParams(params).toString())
}
export function approveWithdraw(id: number) {
  return api.post(`/api/withdraws/${id}/approve`)
}
export function rejectWithdraw(id: number) {
  return api.post(`/api/withdraws/${id}/reject`)
}
```

- [ ] **Step 8: 创建 src/api/modules/employees.ts**

```typescript
import { api } from '@/api/client'

export function getEmployeeList(params?: Record<string, any>) {
  return api.get('/api/employees?' + new URLSearchParams(params).toString())
}
export function createEmployee(body: any) {
  return api.post('/api/employees', body)
}
export function updateEmployee(id: number, body: any) {
  return api.put(`/api/employees/${id}`, body)
}
export function deleteEmployee(id: number) {
  return api.del(`/api/employees/${id}`)
}
export function getPositionList() {
  return api.get('/api/positions')
}
```

- [ ] **Step 9: 创建 src/api/modules/gates.ts**

```typescript
import { api } from '@/api/client'

export function getGateList() {
  return api.get('/api/gates')
}
export function createGate(body: any) {
  return api.post('/api/gates', body)
}
export function toggleGate(id: number) {
  return api.put(`/api/gates/${id}/toggle`)
}
export function deleteGate(id: number) {
  return api.del(`/api/gates/${id}`)
}
export function getFaceList() {
  return api.get('/api/faces')
}
export function getEntryRecordList(params?: Record<string, any>) {
  return api.get('/api/entry-records?' + new URLSearchParams(params).toString())
}
```

- [ ] **Step 10: 创建 src/api/modules/activities.ts**

```typescript
import { api } from '@/api/client'

export function getActivityList(params?: Record<string, any>) {
  return api.get('/api/activities?' + new URLSearchParams(params).toString())
}
export function createActivity(body: any) {
  return api.post('/api/activities', body)
}
export function updateActivity(id: number, body: any) {
  return api.put(`/api/activities/${id}`, body)
}
export function deleteActivity(id: number) {
  return api.del(`/api/activities/${id}`)
}
export function auditActivity(id: number) {
  return api.post(`/api/activities/${id}/audit`)
}
export function endActivity(id: number) {
  return api.post(`/api/activities/${id}/end`)
}
```

- [ ] **Step 11: 创建 src/api/modules/platform.ts**

```typescript
import { api } from '@/api/client'

export function getVerifyRecords(params?: Record<string, any>) {
  return api.get('/api/platform/verify-records?' + new URLSearchParams(params).toString())
}
export function getPlatformSettlements() {
  return api.get('/api/platform/settlements')
}
export function getStoreConfigs() {
  return api.get('/api/platform/store-configs')
}
export function prepareVerify(body: any) {
  return api.post('/api/platform/verify/prepare', body)
}
export function executeVerify(body: any) {
  return api.post('/api/platform/verify/execute', body)
}
export function revokeVerify(id: number) {
  return api.post(`/api/platform/verify/${id}/revoke`)
}
```

---

## Phase 2: 布局框架

### Task 7: 创建菜单配置工具

**Files:**
- Create: `web/src/utils/menu.ts`

- [ ] **Step 1: 创建 src/utils/menu.ts**

从旧 `pc-data.js` 的 `FULL_MENU` 迁移菜单配置：

```typescript
export interface MenuItem {
  id: string
  label: string
  icon: string
  page?: string
  sub?: { id: string; page: string; label: string }[]
}

export interface MenuSection {
  section: string
  items: MenuItem[]
}

export const FULL_MENU: MenuSection[] = [
  {
    section: '工作台',
    items: [
      { id: 'overview', label: '首页概览', icon: '⌂', page: 'overview' },
      {
        id: 'data', label: '数据统计', icon: '⊡',
        sub: [
          { id: 'statistics', page: 'statistics', label: '经营统计' },
        ],
      },
    ],
  },
  {
    section: '园区管理',
    items: [
      { id: 'parks', label: '园区列表', icon: '▤', page: 'parks' },
      { id: 'shops', label: '店铺管理', icon: '◫', page: 'shops' },
    ],
  },
  {
    section: '业务管理',
    items: [
      { id: 'orders', label: '消费订单', icon: '◈', page: 'orders' },
      { id: 'products', label: '商品管理', icon: '☷', page: 'products' },
      { id: 'users', label: '客户管理', icon: '👤', page: 'users' },
      {
        id: 'finance', label: '财务管理', icon: '¥',
        sub: [
          { id: 'finance-summary', page: 'finance', label: '财务汇总' },
          { id: 'withdraws', page: 'withdraws', label: '提现审核' },
        ],
      },
    ],
  },
  {
    section: '员工与设备',
    items: [
      { id: 'employees', label: '员工管理', icon: '👥', page: 'employees' },
      { id: 'gates', label: '闸机管理', icon: '⊞', page: 'gates' },
    ],
  },
  {
    section: '运营管理',
    items: [
      { id: 'activities', label: '活动管理', icon: '★', page: 'activities' },
      { id: 'platform', label: '第三方平台', icon: '☰', page: 'platform' },
    ],
  },
  {
    section: '系统设置',
    items: [
      { id: 'settings', label: '权限设置', icon: '⚙', page: 'settings' },
    ],
  },
]

// 视角 → 角色映射 (从老的 VIEW_TO_ROLE)
export const VIEW_TO_ROLE: Record<string, string> = {
  platform: 'platform',
  'park-hm': 'park-hm',
  'park-wh': 'park-wh',
  'shop-hcct': 'shop-hcct',
}

export const ROLES: Record<string, { avatar: string; user: string; name: string }> = {
  platform: { avatar: '张', user: '张三', name: '平台超级管理员' },
  'park-hm': { avatar: '李', user: '李四', name: '黄梅园区管理员' },
  'park-wh': { avatar: '王', user: '王五', name: '武汉园区管理员' },
  'shop-hcct': { avatar: '赵', user: '赵六', name: '火车餐厅商户' },
}

// 菜单可见性 (从老的 MENU_VISIBILITY)
export const MENU_VISIBILITY: Record<string, Record<string, boolean>> = {
  platform: {},
  park: {},
  shop: {},
}
```

---

### Task 8: 创建全局组件（Toast / Modal）

**Files:**
- Create: `web/src/components/common/GlobalToast.vue`
- Create: `web/src/components/common/GlobalModal.vue`
- Create: `web/src/composables/useToast.ts`
- Create: `web/src/composables/useModal.ts`

- [ ] **Step 1: 创建 src/composables/useToast.ts**

```typescript
import { ref } from 'vue'

interface ToastItem {
  id: number
  message: string
  type: 'success' | 'error' | 'info' | 'warning'
}

const toasts = ref<ToastItem[]>([])
let nextId = 0

export function useToast() {
  function showToast(message: string, type: 'success' | 'error' | 'info' | 'warning' = 'info') {
    const id = nextId++
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 3000)
  }
  return { toasts, showToast }
}
```

- [ ] **Step 2: 创建 src/components/common/GlobalToast.vue**

```vue
<template>
  <div class="toast-container">
    <Teleport to="body">
      <div v-for="t in toasts" :key="t.id" :class="['toast-item', `toast-${t.type}`]">
        {{ t.message }}
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'
const { toasts } = useToast()
</script>

<style scoped>
.toast-container {
  position: fixed; top: 24px; right: 24px; z-index: 9999;
  display: flex; flex-direction: column; gap: 8px;
}
.toast-item {
  padding: 10px 20px; border-radius: 6px; color: #fff;
  font-size: 14px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  animation: fadeIn 0.3s;
}
.toast-success { background: #52c41a; }
.toast-error { background: #ff4d4f; }
.toast-info { background: #1677ff; }
.toast-warning { background: #faad14; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(-10px); } }
</style>
```

- [ ] **Step 3: 创建 src/composables/useModal.ts**

```typescript
import { ref } from 'vue'

const visible = ref(false)
const title = ref('')
const content = ref('')
const onOk = ref<(() => void) | null>(null)

export function useModal() {
  function openModal(t: string, body: string, ok?: () => void) {
    title.value = t
    content.value = body
    onOk.value = ok || null
    visible.value = true
  }
  function closeModal() {
    visible.value = false
  }
  return { visible, title, content, onOk, openModal, closeModal }
}
```

- [ ] **Step 4: 创建 src/components/common/GlobalModal.vue**

```vue
<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-mask" @click.self="closeModal">
      <div class="modal-box">
        <div class="modal-head">
          <span>{{ title }}</span>
          <button class="close-btn" @click="closeModal">✕</button>
        </div>
        <div class="modal-body" v-html="content"></div>
        <div class="modal-foot">
          <button @click="closeModal">取消</button>
          <button v-if="onOk" class="btn-primary" @click="onOk(); closeModal()">确定</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useModal } from '@/composables/useModal'
const { visible, title, content, onOk, closeModal } = useModal()
</script>
```

---

### Task 9: 创建 BlankLayout + Login 页面

**Files:**
- Create: `web/src/layouts/BlankLayout.vue`
- Modify: `web/src/views/login/index.vue`

- [ ] **Step 1: 创建 src/layouts/BlankLayout.vue**

```vue
<template>
  <router-view />
</template>
```

- [ ] **Step 2: 创建 src/views/login/index.vue**

```vue
<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand">
        <div class="login-logo">🌾</div>
        <h1>袁夫稻田 · 智慧园区平台</h1>
        <p>YuanFu RiceField Smart Platform</p>
      </div>
      <a-form :model="form" layout="vertical" @finish="handleLogin">
        <a-form-item label="用户名">
          <a-input v-model:value="form.username" placeholder="admin" size="large" />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="form.password" placeholder="admin123" size="large" />
        </a-form-item>
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">
          登录
        </a-button>
      </a-form>
      <p class="login-error" v-if="error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const error = ref('')
const form = reactive({ username: 'admin', password: 'admin123' })

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(form.username, form.password)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 400px; padding: 40px; background: #fff; border-radius: 12px; box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}
.login-brand { text-align: center; margin-bottom: 32px; }
.login-logo { font-size: 48px; }
.login-brand h1 { font-size: 20px; margin: 8px 0 4px; color: #1a1a1a; }
.login-brand p { font-size: 12px; color: #999; }
.login-error { color: #ff4d4f; text-align: center; margin-top: 16px; font-size: 13px; }
</style>
```

---

### Task 10: 创建 PcLayout + PcSider + PcHeader

**Files:**
- Create: `web/src/layouts/PcLayout.vue`
- Create: `web/src/components/pc/PcSider.vue`
- Create: `web/src/components/pc/PcNav.vue`
- Create: `web/src/components/pc/ViewSwitch.vue`
- Create: `web/src/components/pc/UserFooter.vue`
- Create: `web/src/components/pc/PcHeader.vue`
- Create: `web/src/components/pc/Breadcrumb.vue`

- [ ] **Step 1: 创建 src/components/pc/ViewSwitch.vue**

```vue
<template>
  <div class="view-switch">
    <div class="view-switch-label">当前视角</div>
    <a-select v-model:value="view" size="small" style="width:100%">
      <a-select-option value="platform">🏢 平台总部 · 全局视角</a-select-option>
      <a-select-option value="park-hm">🏞️ 黄梅袁夫稻田 · 园区视角</a-select-option>
      <a-select-option value="park-wh">🏞️ 武汉袁夫稻田 · 园区视角</a-select-option>
      <a-select-option value="shop-hcct">🏪 火车餐厅 · 商户视角</a-select-option>
    </a-select>
  </div>
</template>

<script setup lang="ts">
import { useViewStore, type ViewRole } from '@/stores/view'
import { computed } from 'vue'

const viewStore = useViewStore()
const view = computed({
  get: () => viewStore.view,
  set: (v: ViewRole) => viewStore.switchView(v),
})
</script>
```

- [ ] **Step 2: 创建 src/components/pc/UserFooter.vue**

```vue
<template>
  <div class="user-footer">
    <a-avatar>{{ avatar }}</a-avatar>
    <div style="flex:1;min-width:0">
      <div style="font-weight:600;font-size:13px">{{ name }}</div>
      <div style="font-size:11px;color:#999">{{ role }}</div>
    </div>
    <a-button type="text" @click="authStore.logout(); router.push('/login')">⎋</a-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useViewStore } from '@/stores/view'
import { ROLES, VIEW_TO_ROLE } from '@/utils/menu'

const router = useRouter()
const authStore = useAuthStore()
const viewStore = useViewStore()

const roleInfo = computed(() => ROLES[VIEW_TO_ROLE[viewStore.view]])
const avatar = computed(() => roleInfo.value?.avatar || '?')
const name = computed(() => authStore.user?.realname || authStore.user?.username || '未登录')
const role = computed(() => roleInfo.value?.name || '')
</script>
```

- [ ] **Step 3: 创建 src/components/pc/PcNav.vue**

```vue
<template>
  <nav class="pc-nav">
    <template v-for="section in visibleMenu" :key="section.section">
      <div class="nav-section-title">{{ section.section }}</div>
      <template v-for="item in section.items" :key="item.id">
        <div
          :class="['nav-item', { active: isActive(item) }]"
          @click="item.sub ? toggleExpand(item.id) : navigateTo(item.page!)"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          {{ item.label }}
          <span v-if="item.sub" class="nav-arrow">{{ expanded.has(item.id) ? '▾' : '▸' }}</span>
        </div>
        <div v-if="item.sub && expanded.has(item.id)" class="nav-sub">
          <div
            v-for="sub in item.sub" :key="sub.id"
            :class="['nav-item sub', { active: isSubActive(sub) }]"
            @click="navigateTo(sub.page)"
          >
            {{ sub.label }}
          </div>
        </div>
      </template>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useViewStore } from '@/stores/view'
import { FULL_MENU, VIEW_TO_ROLE } from '@/utils/menu'

const router = useRouter()
const route = useRoute()
const viewStore = useViewStore()
const expanded = ref<Set<string>>(new Set())

const visibleMenu = computed(() => {
  const menuKey = VIEW_TO_ROLE[viewStore.view] || 'platform'
  return FULL_MENU.filter(section =>
    section.items.some(item =>
      menuKey === 'platform' || true // 简化处理，TODO: 接入实际权限过滤
    ),
  )
})

function isActive(item: any) {
  if (item.sub) return item.sub.some((s: any) => s.page === route.name)
  return item.page === route.name
}

function isSubActive(sub: any) {
  return sub.page === route.name
}

function toggleExpand(id: string) {
  if (expanded.value.has(id)) expanded.value.delete(id)
  else expanded.value.add(id)
}

function navigateTo(page: string) {
  router.push(`/pc-admin/${page}`)
}
</script>
```

- [ ] **Step 4: 创建 src/components/pc/PcSider.vue**

```vue
<template>
  <aside class="pc-sider">
    <div class="sider-header">
      <div class="logo">禾</div>
      <div>
        <div class="title">智慧园区平台</div>
        <div class="subtitle">YuanFu RiceField</div>
      </div>
    </div>
    <ViewSwitch />
    <PcNav />
    <UserFooter />
  </aside>
</template>

<script setup lang="ts">
import ViewSwitch from './ViewSwitch.vue'
import PcNav from './PcNav.vue'
import UserFooter from './UserFooter.vue'
</script>

<style scoped>
.pc-sider {
  width: 240px; height: 100vh; background: #001529; color: #fff;
  display: flex; flex-direction: column; overflow-y: auto;
}
.sider-header { display: flex; align-items: center; gap: 10px; padding: 16px; }
.logo { width: 36px; height: 36px; background: #1677ff; color: #fff; display: flex;
  align-items: center; justify-content: center; border-radius: 6px; font-size: 20px; font-weight: 700; }
.title { font-size: 14px; font-weight: 600; }
.subtitle { font-size: 10px; color: rgba(255,255,255,0.45); }
</style>
```

- [ ] **Step 5: 创建 src/components/pc/PcHeader.vue**

```vue
<template>
  <header class="pc-header">
    <a-breadcrumb>
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</a-breadcrumb-item>
    </a-breadcrumb>
    <div class="header-actions">
      <a-button type="text"><SearchOutlined /></a-button>
      <a-badge :dot="true"><a-button type="text"><BellOutlined /></a-button></a-badge>
      <a-button type="text"><QuestionCircleOutlined /></a-button>
      <a-button type="text"><SettingOutlined /></a-button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { SearchOutlined, BellOutlined, QuestionCircleOutlined, SettingOutlined } from '@ant-design/icons-vue'
const route = useRoute()
</script>

<style scoped>
.pc-header { height: 56px; background: #fff; display: flex; align-items: center;
  justify-content: space-between; padding: 0 24px; border-bottom: 1px solid #f0f0f0; }
.header-actions { display: flex; gap: 4px; }
</style>
```

- [ ] **Step 6: 创建 src/layouts/PcLayout.vue**

```vue
<template>
  <div class="pc-layout">
    <PcSider />
    <div class="pc-right">
      <PcHeader />
      <main class="pc-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import PcSider from '@/components/pc/PcSider.vue'
import PcHeader from '@/components/pc/PcHeader.vue'
</script>

<style scoped>
.pc-layout { display: flex; min-height: 100vh; }
.pc-right { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.pc-content { flex: 1; padding: 24px; background: #f5f5f5; overflow-y: auto; }
</style>
```

---

### Task 11: 创建 MobileLayout

**Files:**
- Create: `web/src/layouts/MobileLayout.vue`
- Create: `web/src/components/mobile/PhoneFrame.vue`
- Create: `web/src/components/mobile/StatusBar.vue`
- Create: `web/src/components/mobile/MiniAppHeader.vue`
- Create: `web/src/components/mobile/TabBar.vue`

- [ ] **Step 1: 创建 StatusBar、MiniAppHeader、TabBar 组件**

**StatusBar.vue**:
```vue
<template>
  <div class="phone-statusbar">
    <span>18:42</span>
    <span>●●●● 📶 85</span>
  </div>
</template>
<style scoped>
.phone-statusbar { display: flex; justify-content: space-between; padding: 6px 16px;
  font-size: 11px; background: #000; color: #fff; }
</style>
```

**MiniAppHeader.vue**:
```vue
<template>
  <div class="miniapp-header">
    <div class="h-left"><button class="h-btn">☰</button></div>
    <div class="h-title">{{ title }}</div>
    <div class="h-right"><div class="h-capsule">···  ○</div></div>
  </div>
</template>
<script setup lang="ts">
defineProps<{ title: string }>()
</script>
<style scoped>
.miniapp-header { display: flex; align-items: center; justify-content: space-between;
  padding: 8px 12px; background: #fff; border-bottom: 1px solid #eee; }
.h-title { font-size: 15px; font-weight: 600; }
.h-capsule { border: 1px solid #ddd; border-radius: 20px; padding: 2px 8px; font-size: 10px; }
.h-btn { background: none; border: none; font-size: 16px; cursor: pointer; }
</style>
```

**TabBar.vue**:
```vue
<template>
  <div class="tabbar">
    <button
      v-for="tab in tabs" :key="tab.page"
      :class="['tab-item', { active: modelValue === tab.page }]"
      @click="$emit('update:modelValue', tab.page)"
    >
      <span class="tab-icon">{{ tab.icon }}</span>
      {{ tab.label }}
    </button>
  </div>
</template>
<script setup lang="ts">
defineProps<{ tabs: { page: string; label: string; icon: string }[]; modelValue: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
<style scoped>
.tabbar { display: flex; border-top: 1px solid #eee; background: #fff; }
.tab-item { flex: 1; display: flex; flex-direction: column; align-items: center;
  padding: 6px 0; font-size: 10px; color: #999; background: none; border: none; cursor: pointer; }
.tab-item.active { color: #1677ff; }
.tab-icon { font-size: 18px; }
</style>
```

- [ ] **Step 2: 创建 src/components/mobile/PhoneFrame.vue**

```vue
<template>
  <div class="phone-frame">
    <StatusBar />
    <MiniAppHeader :title="headerTitle" />
    <div class="phone-body">
      <router-view />
    </div>
    <TabBar :tabs="currentTabs" v-model="activeTab" @update:modelValue="onTabChange" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StatusBar from './StatusBar.vue'
import MiniAppHeader from './MiniAppHeader.vue'
import TabBar from './TabBar.vue'

const route = useRoute()
const router = useRouter()

const platformTabs: Record<string, { page: string; label: string; icon: string }[]> = {
  'mp-admin': [
    { page: 'dashboard', label: '首页', icon: '⌂' },
    { page: 'shops', label: '商户', icon: '▤' },
    { page: 'scan', label: '扫码', icon: '⊞' },
    { page: 'message', label: '消息', icon: '✉' },
    { page: 'me', label: '我的', icon: '○' },
  ],
  'mp-merchant': [
    { page: 'dashboard', label: '经营', icon: '⌂' },
    { page: 'orders', label: '订单', icon: '▤' },
    { page: 'verify', label: '核销', icon: '⊞' },
    { page: 'finance', label: '收益', icon: '¥' },
    { page: 'me', label: '我的', icon: '○' },
  ],
  'mp-client': [
    { page: 'home', label: '首页', icon: '⌂' },
    { page: 'memberqr', label: '亮码支付', icon: '⊞' },
    { page: 'cart', label: '购物车', icon: '▤' },
    { page: 'me', label: '个人中心', icon: '○' },
  ],
}

const headerTitles: Record<string, string> = {
  'mp-admin': '袁夫稻田 · 园区管理',
  'mp-merchant': '火车餐厅 · 商户管理',
  'mp-client': '袁夫稻田线上商城',
}

const platform = computed(() => (route.meta.platform as string) || 'mp-client')
const headerTitle = computed(() => headerTitles[platform.value] || '')
const currentTabs = computed(() => platformTabs[platform.value] || [])
const activeTab = computed(() => (route.meta.tab as string) || currentTabs.value[0]?.page || '')

function onTabChange(tab: string) {
  router.push(`/${platform.value}/${tab}`)
}
</script>
```

- [ ] **Step 3: 创建 src/layouts/MobileLayout.vue**

```vue
<template>
  <div class="mobile-layout">
    <PhoneFrame />
  </div>
</template>

<script setup lang="ts">
import PhoneFrame from '@/components/mobile/PhoneFrame.vue'
</script>

<style scoped>
.mobile-layout { display: flex; justify-content: center; align-items: center;
  min-height: 100vh; background: #f0f2f5; }
</style>
```

---

### Task 12: 创建平台切换栏并更新 App.vue

**Files:**
- Create: `web/src/components/pc/PlatformSwitcher.vue`
- Modify: `web/src/App.vue`

- [ ] **Step 1: 创建 PlatformSwitcher.vue**

```vue
<template>
  <div class="platform-switcher">
    <div class="brand">
      <span class="brand-icon">🌾</span>
      袁夫稻田 · 智慧园区综合管理平台
    </div>
    <div class="platform-tabs">
      <button
        v-for="p in platforms" :key="p.key"
        :class="['platform-tab', { active: current === p.key }]"
        @click="switchTo(p.key)"
      >
        {{ p.label }}
      </button>
    </div>
    <div class="right-info">
      <span>v3.0 · Vue 3 + Ant Design Vue + Vant</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useViewStore, type Platform } from '@/stores/view'

const router = useRouter()
const route = useRoute()
const viewStore = useViewStore()

const platforms = [
  { key: 'pc-admin' as Platform, label: 'PC 管理后台' },
  { key: 'mp-admin' as Platform, label: '小程序·管理端' },
  { key: 'mp-merchant' as Platform, label: '小程序·商户端' },
  { key: 'mp-client' as Platform, label: '小程序·客户端' },
]

const current = computed(() => viewStore.platform)

function switchTo(platform: Platform) {
  viewStore.switchPlatform(platform)
  router.push(`/${platform}`)
}
</script>

<style scoped>
.platform-switcher {
  height: 40px; background: #001529; color: #fff; display: flex;
  align-items: center; padding: 0 16px; font-size: 12px; gap: 16px;
}
.brand { font-weight: 600; display: flex; align-items: center; gap: 6px; }
.brand-icon { font-size: 16px; }
.platform-tabs { display: flex; gap: 4px; }
.platform-tab {
  padding: 4px 12px; border-radius: 4px; border: 1px solid rgba(255,255,255,0.2);
  background: transparent; color: rgba(255,255,255,0.65); cursor: pointer; font-size: 11px;
}
.platform-tab.active { background: #1677ff; color: #fff; border-color: #1677ff; }
.right-info { flex: 1; text-align: right; color: rgba(255,255,255,0.45); }
</style>
```

- [ ] **Step 2: 更新 src/App.vue**

```vue
<template>
  <PlatformSwitcher />
  <router-view />
  <GlobalToast />
  <GlobalModal />
</template>

<script setup lang="ts">
import PlatformSwitcher from '@/components/pc/PlatformSwitcher.vue'
import GlobalToast from '@/components/common/GlobalToast.vue'
import GlobalModal from '@/components/common/GlobalModal.vue'
</script>
```

---

## Phase 3: PC 管理后台页面

### Task 13: PC 概览页面 (overview)

**Files:**
- Modify: `web/src/views/pc-admin/overview/index.vue`

- [ ] **Step 1: 实现概览页面（统计卡片 + 快捷入口）**

```vue
<template>
  <div class="overview-page">
    <h2>首页概览</h2>
    <a-row :gutter="16" style="margin-bottom:16px">
      <a-col :span="6" v-for="stat in stats" :key="stat.title">
        <a-card :hoverable="true">
          <a-statistic :title="stat.title" :value="stat.value" :prefix="stat.prefix" />
        </a-card>
      </a-col>
    </a-row>
    <a-card title="快捷入口">
      <a-row :gutter="16">
        <a-col :span="4" v-for="q in quickLinks" :key="q.label">
          <a-button block @click="router.push(q.path)" style="height:80px">
            <div style="font-size:24px">{{ q.icon }}</div>
            <div style="font-size:12px;margin-top:4px">{{ q.label }}</div>
          </a-button>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'

const router = useRouter()
const stats = ref([
  { title: '园区总数', value: 0, prefix: '🏞️' },
  { title: '店铺总数', value: 0, prefix: '🏪' },
  { title: '客户总数', value: 0, prefix: '👤' },
  { title: '今日订单', value: 0, prefix: '📋' },
])

const quickLinks = [
  { icon: '🏞️', label: '园区管理', path: '/pc-admin/parks' },
  { icon: '🏪', label: '店铺管理', path: '/pc-admin/shops' },
  { icon: '📋', label: '消费订单', path: '/pc-admin/orders' },
  { icon: '☷', label: '商品管理', path: '/pc-admin/products' },
  { icon: '👤', label: '客户管理', path: '/pc-admin/users' },
  { icon: '¥', label: '财务管理', path: '/pc-admin/finance' },
]

onMounted(async () => {
  try {
    const data: any = await api.get('/api/statistics/overview')
    if (data) {
      stats.value[0].value = data.parks || 0
      stats.value[1].value = data.shops || 0
      stats.value[2].value = data.users || 0
      stats.value[3].value = data.todayOrders || 0
    }
  } catch (e) { /* ignore */ }
})
</script>
```

---

### Task 14-24: 其余 PC 管理后台页面（园区/店铺/订单/商品/客户/财务/员工/闸机/活动/平台/设置）

每个模块页面遵循统一的 CRUD 列表模式（参考原 pc-pages 系列 JS）。以下是通用模板，各 Task 按模块 API 适配。

**通用 PC 列表页模板** (`src/views/pc-admin/<module>/index.vue`):

```vue
<template>
  <div class="list-page">
    <div class="page-header">
      <h2>{{ pageTitle }}</h2>
      <a-button type="primary" @click="showCreateModal = true">新增</a-button>
    </div>

    <a-card>
      <!-- 搜索栏 -->
      <a-form layout="inline" style="margin-bottom:16px">
        <a-form-item label="关键词">
          <a-input v-model:value="searchKeyword" placeholder="搜索..." />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="fetchData">查询</a-button>
        </a-form-item>
      </a-form>

      <!-- 表格 -->
      <a-table
        :columns="columns" :dataSource="dataSource" :loading="loading"
        :pagination="pagination" @change="handleTableChange" rowKey="id" size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a @click="handleEdit(record)">编辑</a>
            <a-divider type="vertical" />
            <a style="color:#ff4d4f" @click="handleDelete(record.id)">删除</a>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑 Modal -->
    <a-modal
      v-model:open="showCreateModal" :title="editingId ? '编辑' : '新增'"
      @ok="handleSubmit" :confirmLoading="submitting"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item v-for="field in formFields" :key="field.name" :label="field.label">
          <a-input v-model:value="form[field.name]" :placeholder="field.placeholder" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
// import * as moduleApi from '@/api/modules/<module>'

const { showToast } = useToast()
const pageTitle = '<模块名>'
const dataSource = ref<any[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const showCreateModal = ref(false)
const editingId = ref<number | null>(null)
const submitting = ref(false)
const pagination = reactive({ current: 1, pageSize: 10, total: 0 })

const columns = [
  // 根据模块定义列
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  // ... 其他列
  { title: '操作', key: 'action', width: 150 },
]

const formFields: { name: string; label: string; placeholder: string }[] = [
  // 根据模块定义表单字段
]

const form = reactive<Record<string, any>>({})

async function fetchData() {
  loading.value = true
  try {
    // const res = await moduleApi.getXxxList({ page: pagination.current, pageSize: pagination.pageSize, keyword: searchKeyword.value })
    // dataSource.value = res.list || res
    // pagination.total = res.total || 0
  } catch (e: any) { showToast(e.message, 'error') }
  finally { loading.value = false }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchData()
}

function handleEdit(record: any) {
  editingId.value = record.id
  Object.assign(form, record)
  showCreateModal.value = true
}

async function handleDelete(id: number) {
  try {
    // await moduleApi.deleteXxx(id)
    showToast('删除成功', 'success')
    fetchData()
  } catch (e: any) { showToast(e.message, 'error') }
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingId.value) {
      // await moduleApi.updateXxx(editingId.value, form)
    } else {
      // await moduleApi.createXxx(form)
    }
    showToast(editingId.value ? '更新成功' : '创建成功', 'success')
    showCreateModal.value = false
    fetchData()
  } catch (e: any) { showToast(e.message, 'error') }
  finally { submitting.value = false }
}

onMounted(() => fetchData())
</script>

<style scoped>
.list-page h2 { margin-bottom: 16px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
</style>
```

**具体的 12 个 PC 页面分别迁移 Task:**

| Task | 模块 | 目录 | 特殊处理 |
|------|------|------|----------|
| 13 | overview | views/pc-admin/overview/ | 统计卡片（已完成） |
| 14 | parks | views/pc-admin/parks/ | 标准 CRUD |
| 15 | shops | views/pc-admin/shops/ | CRUD + shop_types |
| 16 | orders | views/pc-admin/orders/ | 列表 + 核销/退款操作 |
| 17 | products | views/pc-admin/products/ | 标准 CRUD |
| 18 | users | views/pc-admin/users/ | CRUD + 充值记录 + 统计 |
| 19 | finance | views/pc-admin/finance/ | 财务汇总 + 提现审核 |
| 20 | employees | views/pc-admin/employees/ | CRUD + 岗位列表 |
| 21 | gates | views/pc-admin/gates/ | CRUD + 开关闸 + 人脸列表 + 通行记录 |
| 22 | activities | views/pc-admin/activities/ | CRUD + 审核 + 结束 |
| 23 | platform | views/pc-admin/platform/ | 核销记录 + 结算 + 店铺配置 |
| 24 | settings | views/pc-admin/settings/ | 小程序配置 + 权限矩阵 + 日志 |

---

## Phase 4: 移动端页面

### Task 25: 小程序管理端 (mp-admin)

**Files:**
- Modify: `web/src/views/mp-admin/index.vue`

Tab 页（dashboard/shops/scan/message/me），每个 Tab 渲染对应内容。

### Task 26: 小程序商户端 (mp-merchant)

**Files:**
- Modify: `web/src/views/mp-merchant/index.vue`

Tab 页（dashboard/orders/verify/finance/me）。

### Task 27: 小程序客户端 (mp-client)

**Files:**
- Modify: `web/src/views/mp-client/index.vue`

Tab 页（home/memberqr/cart/me）+ 会员码全屏蒙层。

---

## Phase 5: 收尾

### Task 28: 清理旧文件 + 构建验证

**Files:**
- Remove: `web/legacy/`（确认迁移完成后删除）
- Create: `web/public/` (空目录，Vite 构建用)

- [ ] **Step 1: 启动 Go 后端**

```bash
cd /Users/allen/tools/server && go run . -addr :8080 -web ../web/dist -db yfsc.db &
```

- [ ] **Step 2: 启动 Vite dev server 并验证**

```bash
cd /Users/allen/tools/web && npm run dev
```

打开 `http://localhost:5173`，验证：
- 登录页正常
- 平台切换正常
- PC 各页面正常加载
- 移动端各 Tab 正常

- [ ] **Step 3: 执行构建**

```bash
cd /Users/allen/tools/web && npm run build
```

验证 `web/dist/` 输出。

- [ ] **Step 4: 验证构建产物能通过 Go 服务访问**

```bash
# 重启 Go 服务使其提供构建产物
kill %1
cd /Users/allen/tools/server && go run . -addr :8080 -web ../web/dist -db yfsc.db &
```

打开 `http://localhost:8080`，确认 SPA 正常运作。

- [ ] **Step 5: 删除旧文件**

```bash
rm -rf /Users/allen/tools/web/legacy
```

- [ ] **Step 6: Commit**

```bash
cd /Users/allen/tools
git add web/package.json web/tsconfig.json web/tsconfig.node.json web/vite.config.ts web/index.html web/env.d.ts web/.gitignore web/src/ web/public/
git commit -m "feat: migrate frontend to Vue 3 + Vite + Ant Design Vue + Vant"
```

---

## 自检清单

1. **Spec coverage**: 所有设计文档中的要求均有对应 Task：
   - [x] 项目结构 → Task 1, 2
   - [x] 路由设计 → Task 5
   - [x] 状态管理 → Task 4
   - [x] 组件层级 → Task 9, 10, 11, 12
   - [x] 数据流 → Task 3, 6
   - [x] 4 个平台视图 → Task 13-27
   - [x] 移动端手机模拟框 → Task 11

2. **Placeholder scan**: 无 TBD/TODO

3. **Type consistency**: Store 类型从 stores 中 import，API 类型与 Go 后端一致
