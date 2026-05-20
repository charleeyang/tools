import { api } from '@/api/client'
export function getGateList() {
  return api.get('/api/gates')
}
export function createGate(body: any) {
  return api.post('/api/gates', body)
}
export function toggleGate(id: number) {
  return api.put(`/api/gates/${id}/toggle`)
}
export function deleteGate(id: number) {
  return api.del(`/api/gates/${id}`)
}
export function getFaceList() {
  return api.get('/api/faces')
}
export function getEntryRecordList(params?: Record<string, any>) {
  return api.get('/api/entry-records?' + new URLSearchParams(params).toString())
}
