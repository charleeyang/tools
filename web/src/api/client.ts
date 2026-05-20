import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const http: AxiosInstance = axios.create({
  baseURL: '',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export async function request<T = any>(
  method: string,
  path: string,
  body?: any,
): Promise<T> {
  const config: AxiosRequestConfig = { method, url: path }
  if (body) config.data = body
  const res = await http.request<ApiResponse<T>>(config)
  if (res.data.code !== 0) {
    const err = new Error(res.data.message || '请求失败') as any
    err.code = res.data.code
    throw err
  }
  return res.data.data
}

export const api = {
  get: <T = any>(path: string) => request<T>('GET', path),
  post: <T = any>(path: string, body?: any) => request<T>('POST', path, body),
  put: <T = any>(path: string, body?: any) => request<T>('PUT', path, body),
  del: <T = any>(path: string) => request<T>('DELETE', path),
}
