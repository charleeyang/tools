<template>
  <div class="mobile-layout">
    <PhoneFrame />
  </div>
</template>

<script setup lang="ts">
import PhoneFrame from '@/components/mobile/PhoneFrame.vue'
</script>

<style scoped>
.mobile-layout {
  display: flex; justify-content: center; align-items: center;
  min-height: 100vh; background: #f0f2f5; padding: 40px 0;
}
</style>
