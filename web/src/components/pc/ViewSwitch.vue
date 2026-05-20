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

<style scoped>
.view-switch { padding: 12px 16px; border-bottom: 1px solid rgba(255,255,255,0.1); }
.view-switch-label { font-size: 11px; color: rgba(255,255,255,0.45); margin-bottom: 6px; }
</style>
