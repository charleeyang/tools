import { api } from '@/api/client'
export function getProductList(params?: Record<string, any>) {
  return api.get('/api/products?' + new URLSearchParams(params).toString())
}
export function createProduct(body: any) {
  return api.post('/api/products', body)
}
export function updateProduct(id: number, body: any) {
  return api.put(`/api/products/${id}`, body)
}
export function deleteProduct(id: number) {
  return api.del(`/api/products/${id}`)
}
