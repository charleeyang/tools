<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">系统设置</div>
      <div class="page-subtitle">小程序配置、权限矩阵、操作日志</div>
    </div>

    <a-tabs v-model:activeKey="activeTab">
      <!-- 小程序配置 -->
      <a-tab-pane key="miniprograms" tab="小程序配置">
        <a-card>
          <a-space style="margin-bottom: 16px; justify-content: space-between; width: 100%">
            <span style="color: #999">已接入小程序 ({{ miniprograms.length }})</span>
            <a-button type="primary" @click="openMpCreate"><PlusOutlined /> 添加小程序</a-button>
          </a-space>

          <a-table
            :columns="mpColumns"
            :data-source="miniprograms"
            :loading="mpLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === 'active' ? 'green' : 'default'">
                  {{ record.status === 'active' ? '运行中' : '已停用' }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'type'">
                <a-tag :color="record.type === 'client' ? 'blue' : record.type === 'merchant' ? 'purple' : 'cyan'">
                  {{ record.type_label ?? record.type }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'app_secret'">
                ************
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a @click="openMpEdit(record)">编辑</a>
                  <a @click="handleToggleMpStatus(record)">
                    {{ record.status === 'active' ? '停用' : '启用' }}
                  </a>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 权限矩阵 -->
      <a-tab-pane key="permissions" tab="权限矩阵">
        <a-card>
          <a-space style="margin-bottom: 16px; justify-content: space-between; width: 100%">
            <span style="color: #999">按角色配置功能模块访问权限</span>
            <a-button type="primary" @click="savePermissions">保存权限</a-button>
          </a-space>

          <a-table
            :columns="permColumns"
            :data-source="permList"
            :loading="permLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'visibility'">
                <a-select v-model:value="record.visibility" size="small" style="width: 100px">
                  <a-select-option value="可操作">可操作</a-select-option>
                  <a-select-option value="只读">只读</a-select-option>
                  <a-select-option value="隐藏">隐藏</a-select-option>
                </a-select>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 操作日志 -->
      <a-tab-pane key="logs" tab="操作日志">
        <a-card>
          <a-space style="margin-bottom: 16px" wrap>
            <a-input v-model:value="logKeyword" placeholder="操作人" style="width: 160px" allow-clear @press-enter="fetchLogs" />
            <a-select v-model:value="logType" placeholder="操作类型" style="width: 120px" allow-clear>
              <a-select-option value="新增">新增</a-select-option>
              <a-select-option value="修改">修改</a-select-option>
              <a-select-option value="删除">删除</a-select-option>
              <a-select-option value="查询">查询</a-select-option>
            </a-select>
            <a-button type="primary" @click="fetchLogs"><SearchOutlined /> 查询</a-button>
            <a-button @click="resetLogSearch">重置</a-button>
          </a-space>

          <a-table
            :columns="logColumns"
            :data-source="logs"
            :loading="logLoading"
            row-key="id"
            :pagination="logPagination"
            @change="onLogPageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'action'">
                <a-tag :color="record.action === '删除' ? 'red' : record.action === '新增' ? 'green' : record.action === '查询' ? 'blue' : 'orange'">
                  {{ record.action }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'result'">
                <a-tag :color="record.result === '成功' ? 'green' : 'red'">{{ record.result }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 小程序编辑弹窗 -->
    <a-modal
      v-model:open="mpModalOpen"
      :title="mpIsEdit ? '编辑小程序' : '添加小程序'"
      :confirm-loading="mpSubmitting"
      @ok="handleMpSubmit"
      @cancel="mpModalOpen = false"
      width="560px"
    >
      <a-form :model="mpForm" layout="vertical">
        <a-form-item label="小程序名称" required>
          <a-input v-model:value="mpForm.name" placeholder="如：袁夫稻田·客户端" />
        </a-form-item>
        <a-form-item label="AppID" required>
          <a-input v-model:value="mpForm.app_id" placeholder="wx..." />
        </a-form-item>
        <a-form-item label="AppSecret" required>
          <a-input-password v-model:value="mpForm.app_secret" placeholder="AppSecret" />
        </a-form-item>
        <a-form-item label="端类型">
          <a-select v-model:value="mpForm.type">
            <a-select-option value="client">客户端</a-select-option>
            <a-select-option value="merchant">商户端</a-select-option>
            <a-select-option value="admin">管理端</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="版本号">
          <a-input v-model:value="mpForm.version" placeholder="v1.0.0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { api } from '@/api/client'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const activeTab = ref('miniprograms')

// Mini programs
const mpLoading = ref(false)
const miniprograms = ref<any[]>([])
const mpModalOpen = ref(false)
const mpIsEdit = ref(false)
const mpEditId = ref<number | null>(null)
const mpSubmitting = ref(false)
const mpForm = ref({
  name: '',
  app_id: '',
  app_secret: '',
  type: 'client' as string,
  version: 'v1.0.0',
})

const mpColumns = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: 'AppID', dataIndex: 'app_id', key: 'app_id', width: 200 },
  { title: 'AppSecret', key: 'app_secret', width: 140 },
  { title: '端类型', key: 'type', width: 100 },
  { title: '版本', dataIndex: 'version', key: 'version', width: 90 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 130 },
]

// Permissions
const permLoading = ref(false)
const permList = ref<any[]>([])
const permColumns = [
  { title: '角色', dataIndex: 'role_name', key: 'role_name', width: 150 },
  { title: '角色编码', dataIndex: 'role_code', key: 'role_code', width: 180 },
  { title: '功能模块', dataIndex: 'module_key', key: 'module_key', width: 180 },
  { title: '访问权限', key: 'visibility', width: 130 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
]

// Logs
const logLoading = ref(false)
const logs = ref<any[]>([])
const logKeyword = ref('')
const logType = ref<string | undefined>(undefined)
const logPagination = ref({ current: 1, pageSize: 10, total: 0 })

const logColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '操作人', dataIndex: 'operator', key: 'operator', width: 120 },
  { title: '角色', dataIndex: 'role_code', key: 'role_code', width: 140 },
  { title: '操作类型', key: 'action', width: 90 },
  { title: '操作模块', dataIndex: 'module', key: 'module', width: 120 },
  { title: '操作描述', dataIndex: 'detail', key: 'detail' },
  { title: '状态', key: 'result', width: 80 },
  { title: '操作时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
]

// Mini programs
async function fetchMiniPrograms() {
  mpLoading.value = true
  try {
    const res: any = await api.get('/api/system/mini-programs')
    miniprograms.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载小程序配置失败', 'error')
  } finally {
    mpLoading.value = false
  }
}

function openMpCreate() {
  mpIsEdit.value = false
  mpEditId.value = null
  mpForm.value = { name: '', app_id: '', app_secret: '', type: 'client', version: 'v1.0.0' }
  mpModalOpen.value = true
}

function openMpEdit(record: any) {
  mpIsEdit.value = true
  mpEditId.value = record.id
  mpForm.value = {
    name: record.name ?? '',
    app_id: record.app_id ?? '',
    app_secret: '',
    type: record.type ?? 'client',
    version: record.version ?? 'v1.0.0',
  }
  mpModalOpen.value = true
}

async function handleMpSubmit() {
  if (!mpForm.value.name || !mpForm.value.app_id) {
    showToast('请填写必填项', 'error')
    return
  }
  mpSubmitting.value = true
  try {
    if (mpIsEdit.value && mpEditId.value) {
      await api.put(`/api/system/mini-programs/${mpEditId.value}`, mpForm.value)
      showToast('小程序配置已更新')
    } else {
      // Create via put to a new path, or post
      await api.put('/api/system/mini-programs/0', mpForm.value)
      showToast('小程序已添加')
    }
    mpModalOpen.value = false
    fetchMiniPrograms()
  } catch (err: any) {
    showToast(err.message || '保存失败', 'error')
  } finally {
    mpSubmitting.value = false
  }
}

async function handleToggleMpStatus(record: any) {
  try {
    const newStatus = record.status === 'active' ? 'inactive' : 'active'
    await api.put(`/api/system/mini-programs/${record.id}`, { ...record, status: newStatus })
    showToast(`小程序已${newStatus === 'active' ? '启用' : '停用'}`)
    fetchMiniPrograms()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

// Permissions
async function fetchPermissions() {
  permLoading.value = true
  try {
    const res: any = await api.get('/api/permissions')
    permList.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载权限配置失败', 'error')
  } finally {
    permLoading.value = false
  }
}

async function savePermissions() {
  try {
    await api.put('/api/permissions', permList.value)
    showToast('权限已保存')
  } catch (err: any) {
    showToast(err.message || '保存权限失败', 'error')
  }
}

// Logs
async function fetchLogs() {
  logLoading.value = true
  try {
    const params: Record<string, any> = { page: logPagination.value.current, page_size: logPagination.value.pageSize }
    if (logKeyword.value) params.keyword = logKeyword.value
    if (logType.value) params.type = logType.value
    const res: any = await api.get('/api/system/logs?' + new URLSearchParams(params).toString())
    if (Array.isArray(res)) {
      logs.value = res
      logPagination.value.total = res.length
    } else if (res?.list) {
      logs.value = res.list
      logPagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载操作日志失败', 'error')
  } finally {
    logLoading.value = false
  }
}

function onLogPageChange(pag: any) {
  logPagination.value.current = pag.current
  logPagination.value.pageSize = pag.pageSize
  fetchLogs()
}

function resetLogSearch() {
  logKeyword.value = ''
  logType.value = undefined
  logPagination.value.current = 1
  fetchLogs()
}

onMounted(() => {
  fetchMiniPrograms()
  fetchPermissions()
  fetchLogs()
})
</script>
