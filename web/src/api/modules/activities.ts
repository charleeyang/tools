import { api } from '@/api/client'
export function getActivityList(params?: Record<string, any>) {
  return api.get('/api/activities?' + new URLSearchParams(params).toString())
}
export function createActivity(body: any) {
  return api.post('/api/activities', body)
}
export function updateActivity(id: number, body: any) {
  return api.put(`/api/activities/${id}`, body)
}
export function deleteActivity(id: number) {
  return api.del(`/api/activities/${id}`)
}
export function auditActivity(id: number) {
  return api.post(`/api/activities/${id}/audit`)
}
export function endActivity(id: number) {
  return api.post(`/api/activities/${id}/end`)
}
