import { api } from '@/api/client'
export function login(username: string, password: string) {
  return api.post('/api/auth/login', { username, password })
}
export function getMe() {
  return api.get('/api/auth/me')
}
