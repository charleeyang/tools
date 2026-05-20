import { api } from '@/api/client'
export interface Park {
  id: number; name: string; address: string; status: string
  created_at: string; updated_at: string
}
export function getParkList(params?: Record<string, any>) {
  return api.get<Park[]>('/api/parks?' + new URLSearchParams(params).toString())
}
export function createPark(body: Partial<Park>) {
  return api.post<Park>('/api/parks', body)
}
export function updatePark(id: number, body: Partial<Park>) {
  return api.put<Park>(`/api/parks/${id}`, body)
}
export function deletePark(id: number) {
  return api.del(`/api/parks/${id}`)
}
