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
.miniapp-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 12px; background: #fff; border-bottom: 1px solid #eee; flex-shrink: 0;
}
.h-title { font-size: 15px; font-weight: 600; }
.h-capsule { border: 1px solid #ddd; border-radius: 20px; padding: 2px 8px; font-size: 10px; }
.h-btn { background: none; border: none; font-size: 16px; cursor: pointer; }
</style>
