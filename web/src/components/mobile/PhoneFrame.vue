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

<style scoped>
.phone-frame {
  width: 375px; height: 750px; background: #f5f5f5;
  border-radius: 32px; overflow: hidden; display: flex; flex-direction: column;
  box-shadow: 0 0 0 4px #1a1a1a, 0 0 0 6px #333, 0 20px 60px rgba(0,0,0,0.3);
}
.phone-body { flex: 1; overflow-y: auto; background: #fff; }
</style>
