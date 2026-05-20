import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api/client'

export interface User {
  id: number
  username: string
  realname: string
  role: string
  avatar?: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('yfsc_token') || '')
  const user = ref<User | null>(
    JSON.parse(localStorage.getItem('yfsc_user') || 'null'),
  )
  const permissions = ref<Record<string, { read: boolean; write: boolean }>>(
    JSON.parse(localStorage.getItem('yfsc_perms') || '{}'),
  )

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const data: any = await api.post('/api/auth/login', { username, password })
    token.value = data.token
    user.value = data.user
    permissions.value = data.permissions || {}
    localStorage.setItem('yfsc_token', token.value)
    localStorage.setItem('yfsc_user', JSON.stringify(user.value))
    localStorage.setItem('yfsc_perms', JSON.stringify(permissions.value))
  }

  async function fetchMe() {
    const data: any = await api.get('/api/auth/me')
    user.value = data.user || data
  }

  function logout() {
    token.value = ''
    user.value = null
    permissions.value = {}
    localStorage.removeItem('yfsc_token')
    localStorage.removeItem('yfsc_user')
    localStorage.removeItem('yfsc_perms')
  }

  return { token, user, permissions, isLoggedIn, login, logout, fetchMe }
})
