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
.pc-header {
  height: 56px; background: #fff; display: flex; align-items: center;
  justify-content: space-between; padding: 0 24px; border-bottom: 1px solid #f0f0f0;
}
.header-actions { display: flex; gap: 4px; align-items: center; }
</style>
