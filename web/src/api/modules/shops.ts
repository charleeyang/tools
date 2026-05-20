import { api } from '@/api/client'
export interface Shop {
  id: number; name: string; park_id: number; shop_type_id: number
  status: string; created_at: string
}
export function getShopList(params?: Record<string, any>) {
  return api.get<Shop[]>('/api/shops?' + new URLSearchParams(params).toString())
}
export function createShop(body: Partial<Shop>) {
  return api.post<Shop>('/api/shops', body)
}
export function updateShop(id: number, body: Partial<Shop>) {
  return api.put<Shop>(`/api/shops/${id}`, body)
}
export function deleteShop(id: number) {
  return api.del(`/api/shops/${id}`)
}
export function getShopTypes() {
  return api.get<any[]>('/api/shop-types')
}
export function createShopType(body: any) {
  return api.post<any>('/api/shop-types', body)
}
