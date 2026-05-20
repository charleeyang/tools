<template>
  <aside class="pc-sider">
    <div class="sider-header">
      <div class="logo">禾</div>
      <div>
        <div class="title">智慧园区平台</div>
        <div class="subtitle">YuanFu RiceField</div>
      </div>
    </div>
    <ViewSwitch />
    <PcNav />
    <UserFooter />
  </aside>
</template>

<script setup lang="ts">
import ViewSwitch from './ViewSwitch.vue'
import PcNav from './PcNav.vue'
import UserFooter from './UserFooter.vue'
</script>

<style scoped>
.pc-sider {
  width: 240px; min-height: 100vh; background: #001529; color: #fff;
  display: flex; flex-direction: column; flex-shrink: 0;
}
.sider-header { display: flex; align-items: center; gap: 10px; padding: 16px; }
.logo {
  width: 36px; height: 36px; background: #1677ff; color: #fff;
  display: flex; align-items: center; justify-content: center;
  border-radius: 6px; font-size: 20px; font-weight: 700;
}
.title { font-size: 14px; font-weight: 600; }
.subtitle { font-size: 10px; color: rgba(255,255,255,0.45); }
</style>
