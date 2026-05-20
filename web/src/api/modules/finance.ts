import { api } from '@/api/client'
export function getFinanceSummary() {
  return api.get('/api/finance/summary')
}
export function getWithdrawList(params?: Record<string, any>) {
  return api.get('/api/withdraws?' + new URLSearchParams(params).toString())
}
export function approveWithdraw(id: number) {
  return api.post(`/api/withdraws/${id}/approve`)
}
export function rejectWithdraw(id: number) {
  return api.post(`/api/withdraws/${id}/reject`)
}
