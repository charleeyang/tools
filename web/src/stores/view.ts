import { defineStore } from 'pinia'
import { ref } from 'vue'

export type Platform = 'pc-admin' | 'mp-admin' | 'mp-merchant' | 'mp-client'
export type ViewRole = 'platform' | 'park-hm' | 'park-wh' | 'shop-hcct'

export const useViewStore = defineStore('view', () => {
  const platform = ref<Platform>('pc-admin')
  const view = ref<ViewRole>('platform')

  function switchPlatform(p: Platform) {
    platform.value = p
  }

  function switchView(v: ViewRole) {
    view.value = v
  }

  return { platform, view, switchPlatform, switchView }
})
