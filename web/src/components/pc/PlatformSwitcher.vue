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
  position: sticky; top: 0; z-index: 100;
}
.brand { font-weight: 600; display: flex; align-items: center; gap: 6px; white-space: nowrap; }
.brand-icon { font-size: 16px; }
.platform-tabs { display: flex; gap: 4px; }
.platform-tab {
  padding: 4px 12px; border-radius: 4px; border: 1px solid rgba(255,255,255,0.2);
  background: transparent; color: rgba(255,255,255,0.65); cursor: pointer; font-size: 11px;
  white-space: nowrap; transition: all 0.2s;
}
.platform-tab:hover { border-color: rgba(255,255,255,0.4); color: #fff; }
.platform-tab.active { background: #1677ff; color: #fff; border-color: #1677ff; }
.right-info { flex: 1; text-align: right; color: rgba(255,255,255,0.45); }
</style>
