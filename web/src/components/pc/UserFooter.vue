<template>
  <div class="user-footer">
    <a-avatar size="small">{{ avatar }}</a-avatar>
    <div style="flex:1;min-width:0">
      <div style="font-weight:600;font-size:13px;color:#fff">{{ name }}</div>
      <div style="font-size:11px;color:rgba(255,255,255,0.45)">{{ role }}</div>
    </div>
    <a-button type="text" style="color:rgba(255,255,255,0.65)" @click="handleLogout">⎋</a-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useViewStore } from '@/stores/view'
import { ROLES, VIEW_TO_ROLE } from '@/utils/menu'

const router = useRouter()
const authStore = useAuthStore()
const viewStore = useViewStore()

const roleInfo = computed(() => ROLES[VIEW_TO_ROLE[viewStore.view]])
const avatar = computed(() => roleInfo.value?.avatar || '?')
const name = computed(() => authStore.user?.realname || authStore.user?.username || '未登录')
const role = computed(() => roleInfo.value?.name || '')

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.user-footer {
  display: flex; align-items: center; gap: 8px; padding: 12px 16px;
  border-top: 1px solid rgba(255,255,255,0.1); margin-top: auto;
}
</style>
