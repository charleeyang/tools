import { api } from '@/api/client'
export function getUserList(params?: Record<string, any>) {
  return api.get('/api/users?' + new URLSearchParams(params).toString())
}
export function getUserStats() {
  return api.get('/api/users/stats')
}
export function createUser(body: any) {
  return api.post('/api/users', body)
}
export function updateUser(id: number, body: any) {
  return api.put(`/api/users/${id}`, body)
}
export function deleteUser(id: number) {
  return api.del(`/api/users/${id}`)
}
export function getRechargeList(params?: Record<string, any>) {
  return api.get('/api/recharges?' + new URLSearchParams(params).toString())
}
