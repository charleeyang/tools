import { api } from '@/api/client'
export function getOrderList(params?: Record<string, any>) {
  return api.get('/api/orders?' + new URLSearchParams(params).toString())
}
export function verifyOrder(no: string) {
  return api.post(`/api/orders/${no}/verify`)
}
export function refundOrder(no: string) {
  return api.post(`/api/orders/${no}/refund`)
}
export function getRefundList(params?: Record<string, any>) {
  return api.get('/api/refunds?' + new URLSearchParams(params).toString())
}
export function approveRefund(id: number) {
  return api.post(`/api/refunds/${id}/approve`)
}
export function rejectRefund(id: number) {
  return api.post(`/api/refunds/${id}/reject`)
}
