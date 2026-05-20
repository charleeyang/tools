<template>
  <Teleport to="body">
    <div class="toast-container">
      <div v-for="t in toasts" :key="t.id" :class="['toast-item', `toast-${t.type}`]">
        {{ t.message }}
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'
const { toasts } = useToast()
</script>

<style scoped>
.toast-container {
  position: fixed; top: 24px; right: 24px; z-index: 9999;
  display: flex; flex-direction: column; gap: 8px; pointer-events: none;
}
.toast-item {
  padding: 10px 20px; border-radius: 6px; color: #fff;
  font-size: 14px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  animation: toastFadeIn 0.3s;
}
.toast-success { background: #52c41a; }
.toast-error { background: #ff4d4f; }
.toast-info { background: #1677ff; }
.toast-warning { background: #faad14; }
@keyframes toastFadeIn { from { opacity: 0; transform: translateY(-10px); } }
</style>
