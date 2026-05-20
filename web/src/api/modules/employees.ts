import { api } from '@/api/client'
export function getEmployeeList(params?: Record<string, any>) {
  return api.get('/api/employees?' + new URLSearchParams(params).toString())
}
export function createEmployee(body: any) {
  return api.post('/api/employees', body)
}
export function updateEmployee(id: number, body: any) {
  return api.put(`/api/employees/${id}`, body)
}
export function deleteEmployee(id: number) {
  return api.del(`/api/employees/${id}`)
}
export function getPositionList() {
  return api.get('/api/positions')
}
