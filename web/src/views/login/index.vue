<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand">
        <div class="login-logo">🌾</div>
        <h1>袁夫稻田 · 智慧园区平台</h1>
        <p>YuanFu RiceField Smart Platform</p>
      </div>
      <a-form :model="form" layout="vertical" @finish="handleLogin">
        <a-form-item label="用户名">
          <a-input v-model:value="form.username" placeholder="admin" size="large" />
        </a-form-item>
        <a-form-item label="密码">
          <a-input-password v-model:value="form.password" placeholder="admin123" size="large" />
        </a-form-item>
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">
          登录
        </a-button>
      </a-form>
      <p class="login-error" v-if="error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const error = ref('')
const form = reactive({ username: 'admin', password: 'admin123' })

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(form.username, form.password)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 400px; padding: 40px; background: #fff; border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}
.login-brand { text-align: center; margin-bottom: 32px; }
.login-logo { font-size: 48px; }
.login-brand h1 { font-size: 20px; margin: 8px 0 4px; color: #1a1a1a; }
.login-brand p { font-size: 12px; color: #999; }
.login-error { color: #ff4d4f; text-align: center; margin-top: 16px; font-size: 13px; }
</style>
