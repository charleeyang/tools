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
defineProps<{
  tabs: { page: string; label: string; icon: string }[]
  modelValue: string
}>()
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
<style scoped>
.tabbar {
  display: flex; border-top: 1px solid #eee; background: #fff; flex-shrink: 0;
}
.tab-item {
  flex: 1; display: flex; flex-direction: column; align-items: center;
  padding: 6px 0; font-size: 10px; color: #999;
  background: none; border: none; cursor: pointer;
}
.tab-item.active { color: #1677ff; }
.tab-icon { font-size: 18px; margin-bottom: 1px; }
</style>
