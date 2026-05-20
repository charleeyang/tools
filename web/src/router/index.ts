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
