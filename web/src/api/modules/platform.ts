import { api } from '@/api/client'
export function getVerifyRecords(params?: Record<string, any>) {
  return api.get('/api/platform/verify-records?' + new URLSearchParams(params).toString())
}
export function getPlatformSettlements() {
  return api.get('/api/platform/settlements')
}
export function getStoreConfigs() {
  return api.get('/api/platform/store-configs')
}
export function prepareVerify(body: any) {
  return api.post('/api/platform/verify/prepare', body)
}
export function executeVerify(body: any) {
  return api.post('/api/platform/verify/execute', body)
}
export function revokeVerify(id: number) {
  return api.post(`/api/platform/verify/${id}/revoke`)
}
