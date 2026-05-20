<template>
  <PlatformSwitcher />
  <router-view />
  <GlobalToast />
  <GlobalModal />
</template>

<script setup lang="ts">
import PlatformSwitcher from '@/components/pc/PlatformSwitcher.vue'
import GlobalToast from '@/components/common/GlobalToast.vue'
import GlobalModal from '@/components/common/GlobalModal.vue'
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
</style>
