<template>
  <nav class="pc-nav">
    <template v-for="section in FULL_MENU" :key="section.section">
      <div class="nav-section-title">{{ section.section }}</div>
      <template v-for="item in section.items" :key="item.id">
        <div
          :class="['nav-item', { active: isActive(item) }]"
          @click="item.sub ? toggleExpand(item.id) : navigateTo(item.page!)"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          {{ item.label }}
          <span v-if="item.sub" class="nav-arrow">{{ expanded.has(item.id) ? '▾' : '▸' }}</span>
        </div>
        <div v-if="item.sub && expanded.has(item.id)" class="nav-sub">
          <div
            v-for="sub in item.sub" :key="sub.id"
            :class="['nav-item sub', { active: isSubActive(sub) }]"
            @click="navigateTo(sub.page)"
          >
            {{ sub.label }}
          </div>
        </div>
      </template>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { FULL_MENU } from '@/utils/menu'

const router = useRouter()
const route = useRoute()
const expanded = ref<Set<string>>(new Set())

function isActive(item: any) {
  if (item.sub) return item.sub.some((s: any) => s.page === route.name)
  return item.page === route.name
}

function isSubActive(sub: any) {
  return sub.page === route.name
}

function toggleExpand(id: string) {
  const next = new Set(expanded.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  expanded.value = next
}

function navigateTo(page: string) {
  router.push(`/pc-admin/${page}`)
}
</script>

<style scoped>
.pc-nav { flex: 1; overflow-y: auto; padding: 8px 0; }
.nav-section-title { padding: 8px 16px 4px; font-size: 11px; color: rgba(255,255,255,0.35); text-transform: uppercase; }
.nav-item {
  display: flex; align-items: center; gap: 8px; padding: 8px 16px;
  font-size: 13px; color: rgba(255,255,255,0.65); cursor: pointer;
  transition: all 0.2s;
}
.nav-item:hover { color: #fff; background: rgba(255,255,255,0.08); }
.nav-item.active { color: #fff; background: #1677ff; }
.nav-item.sub { padding-left: 48px; font-size: 12px; }
.nav-icon { font-size: 14px; width: 20px; text-align: center; }
.nav-arrow { margin-left: auto; font-size: 10px; }
</style>
