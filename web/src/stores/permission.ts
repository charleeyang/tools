import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

export const usePermissionStore = defineStore('permission', () => {
  const authStore = useAuthStore()

  function canRead(module: string): boolean {
    return authStore.permissions[module]?.read ?? false
  }

  function canWrite(module: string): boolean {
    return authStore.permissions[module]?.write ?? false
  }

  return { canRead, canWrite }
})
