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
