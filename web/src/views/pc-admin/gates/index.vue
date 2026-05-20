<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">闸机管理</div>
      <div class="page-subtitle">管理园区闸机设备、人脸录入、进出记录</div>
    </div>

    <a-tabs v-model:activeKey="activeTab">
      <!-- 闸机列表 -->
      <a-tab-pane key="gates" tab="闸机列表">
        <a-card>
          <a-space style="margin-bottom: 16px" wrap>
            <a-input v-model:value="gateKeyword" placeholder="闸机名称 / SN" style="width: 200px" allow-clear @press-enter="fetchGates" />
            <a-button type="primary" @click="fetchGates"><SearchOutlined /> 查询</a-button>
            <a-button type="primary" style="margin-left: auto" @click="openGateCreate"><PlusOutlined /> 添加闸机</a-button>
          </a-space>

          <a-table
            :columns="gateColumns"
            :data-source="gateList"
            :loading="gateLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'direction'">
                <a-tag :color="record.direction === '出' ? 'red' : 'blue'">{{ record.direction === '出' ? '出口' : '入口' }}</a-tag>
              </template>
              <template v-else-if="column.key === 'enabled'">
                <a-tag :color="record.enabled ? 'green' : 'default'">{{ record.enabled ? '启用' : '禁用' }}</a-tag>
              </template>
              <template v-else-if="column.key === 'onlineStatus'">
                <a-tag :color="record.online_status || record.onlineStatus ? 'green' : 'default'">
                  {{ record.online_status || record.onlineStatus ? '在线' : '离线' }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a @click="handleToggleGate(record)">{{ record.enabled ? '禁用' : '启用' }}</a>
                  <a-popconfirm title="确定删除该闸机？" @confirm="handleDeleteGate(record.id)">
                    <a style="color: #ff4d4f">删除</a>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 人脸管理 -->
      <a-tab-pane key="faces" tab="人脸管理">
        <a-card>
          <a-space style="margin-bottom: 16px">
            <a-button type="primary"><PlusOutlined /> 录入人脸</a-button>
          </a-space>
          <a-table
            :columns="faceColumns"
            :data-source="faceList"
            :loading="faceLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === '启用' ? 'green' : 'default'">{{ record.status }}</a-tag>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a>配置闸机</a>
                  <a>查看</a>
                  <a style="color: #ff4d4f">删除</a>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 进出记录 -->
      <a-tab-pane key="records" tab="进出记录">
        <a-card>
          <a-table
            :columns="recordColumns"
            :data-source="entryRecords"
            :loading="recordLoading"
            row-key="id"
            :pagination="recordPagination"
            @change="onRecordPageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'direction'">
                <a-tag :color="record.direction === '进' ? 'green' : 'orange'">{{ record.direction }}</a-tag>
              </template>
              <template v-else-if="column.key === 'verify_method'">
                <a-tag :color="record.verify_method === '人脸识别' ? 'blue' : record.verify_method === '扫码通行' ? 'cyan' : 'orange'">
                  {{ record.verify_method }}
                </a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 添加闸机弹窗 -->
    <a-modal
      v-model:open="gateModalOpen"
      title="添加闸机"
      :confirm-loading="gateSubmitting"
      @ok="handleGateSubmit"
      @cancel="gateModalOpen = false"
    >
      <a-form :model="gateForm" layout="vertical">
        <a-form-item label="闸机名称" required>
          <a-input v-model:value="gateForm.name" placeholder="如：主入口 1#" />
        </a-form-item>
        <a-form-item label="方向">
          <a-select v-model:value="gateForm.direction">
            <a-select-option value="进">进</a-select-option>
            <a-select-option value="出">出</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="设备 SN">
          <a-input v-model:value="gateForm.device_sn" placeholder="设备序列号" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import {
  getGateList, createGate, toggleGate, deleteGate,
  getFaceList, getEntryRecordList,
} from '@/api/modules/gates'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const activeTab = ref('gates')

// Gates
const gateLoading = ref(false)
const gateList = ref<any[]>([])
const gateKeyword = ref('')
const gateModalOpen = ref(false)
const gateSubmitting = ref(false)
const gateForm = ref({ name: '', direction: '进', device_sn: '' })

const gateColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '闸机名称', dataIndex: 'name', key: 'name' },
  { title: '方向', key: 'direction', width: 80 },
  { title: '设备SN', dataIndex: 'device_sn', key: 'device_sn', width: 160 },
  { title: '启用状态', key: 'enabled', width: 100 },
  { title: '在线状态', key: 'onlineStatus', width: 100 },
  { title: '今日通行', dataIndex: 'today_pass', key: 'today_pass', width: 100 },
  { title: '操作', key: 'action', width: 130 },
]

// Faces
const faceLoading = ref(false)
const faceList = ref<any[]>([])
const faceColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '用户ID', dataIndex: 'user_id', key: 'user_id', width: 100 },
  { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
  { title: '人脸特征ID', dataIndex: 'feature_id', key: 'feature_id', width: 280 },
  { title: '状态', key: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 160 },
]

// Entry records
const recordLoading = ref(false)
const entryRecords = ref<any[]>([])
const recordPagination = ref({ current: 1, pageSize: 10, total: 0 })
const recordColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '用户ID', dataIndex: 'user_id', key: 'user_id', width: 90 },
  { title: '闸机', dataIndex: 'gate_name', key: 'gate_name', width: 120 },
  { title: '方向', key: 'direction', width: 70 },
  { title: '验证方式', key: 'verify_method', width: 100 },
  { title: '进出时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
]

async function fetchGates() {
  gateLoading.value = true
  try {
    const res: any = await getGateList()
    gateList.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载闸机列表失败', 'error')
  } finally {
    gateLoading.value = false
  }
}

async function fetchFaces() {
  faceLoading.value = true
  try {
    const res: any = await getFaceList()
    faceList.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载人脸列表失败', 'error')
  } finally {
    faceLoading.value = false
  }
}

async function fetchRecords() {
  recordLoading.value = true
  try {
    const params = { page: recordPagination.value.current, page_size: recordPagination.value.pageSize }
    const res: any = await getEntryRecordList(params)
    if (Array.isArray(res)) {
      entryRecords.value = res
      recordPagination.value.total = res.length
    } else if (res?.list) {
      entryRecords.value = res.list
      recordPagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载进出记录失败', 'error')
  } finally {
    recordLoading.value = false
  }
}

function onRecordPageChange(pag: any) {
  recordPagination.value.current = pag.current
  recordPagination.value.pageSize = pag.pageSize
  fetchRecords()
}

function openGateCreate() {
  gateForm.value = { name: '', direction: '进', device_sn: '' }
  gateModalOpen.value = true
}

async function handleGateSubmit() {
  if (!gateForm.value.name) {
    showToast('请填写闸机名称', 'error')
    return
  }
  gateSubmitting.value = true
  try {
    await createGate(gateForm.value)
    showToast('闸机添加成功')
    gateModalOpen.value = false
    fetchGates()
  } catch (err: any) {
    showToast(err.message || '添加失败', 'error')
  } finally {
    gateSubmitting.value = false
  }
}

async function handleToggleGate(record: any) {
  try {
    await toggleGate(record.id)
    showToast(`闸机已${record.enabled ? '禁用' : '启用'}`)
    fetchGates()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

async function handleDeleteGate(id: number) {
  try {
    await deleteGate(id)
    showToast('闸机已删除')
    fetchGates()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

onMounted(() => {
  fetchGates()
  fetchFaces()
  fetchRecords()
})
</script>
