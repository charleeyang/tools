<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">活动管理</div>
      <div class="page-subtitle">创建和审核营销活动</div>
    </div>

    <!-- 活动统计 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card><a-statistic title="进行中" :value="stats.active" suffix="个" :value-style="{ color: '#52c41a' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="待审核" :value="stats.pending" suffix="个" :value-style="{ color: '#faad14' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="已结束" :value="stats.ended" suffix="个" :value-style="{ color: '#999' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="累计参与" :value="stats.totalJoins" suffix="人次" :value-style="{ color: '#1677ff' }" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-select v-model:value="filterStatus" placeholder="状态" style="width: 120px" allow-clear>
          <a-select-option value="进行中">进行中</a-select-option>
          <a-select-option value="待审核">待审核</a-select-option>
          <a-select-option value="已结束">已结束</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
        <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 创建活动</a-button>
      </a-space>

      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @change="onPageChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === '进行中' ? 'green' : record.status === '待审核' ? 'orange' : 'default'">
              {{ record.status }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'time'">
            {{ record.start_time }} ~ {{ record.end_time }}
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确定删除该活动？" @confirm="handleDelete(record.id)">
                <a style="color: #ff4d4f">删除</a>
              </a-popconfirm>
              <template v-if="record.status === '待审核'">
                <a @click="handleAudit(record.id)">审核通过</a>
                <a style="color: #ff4d4f" @click="handleReject(record.id)">驳回</a>
              </template>
              <a v-if="record.status === '进行中'" style="color: #ff4d4f" @click="handleEnd(record.id)">结束</a>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑活动' : '创建活动'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
      width="640px"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="活动标题" required>
          <a-input v-model:value="form.title" placeholder="如：夏日稻田音乐节" />
        </a-form-item>
        <a-form-item label="活动描述">
          <a-textarea v-model:value="form.description" :rows="3" placeholder="活动详情描述" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="开始时间" required>
              <a-date-picker v-model:value="form.start_time" style="width: 100%" placeholder="开始时间" show-time />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结束时间" required>
              <a-date-picker v-model:value="form.end_time" style="width: 100%" placeholder="结束时间" show-time />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="适用园区">
          <a-select v-model:value="form.park_id" placeholder="选择园区" allow-clear>
            <a-select-option v-for="p in parkOptions" :key="p.id" :value="p.id">{{ p.name }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Dayjs } from 'dayjs'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import {
  getActivityList, createActivity, updateActivity, deleteActivity, auditActivity, endActivity,
} from '@/api/modules/activities'
import { getParkList } from '@/api/modules/parks'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)
const submitting = ref(false)
const modalOpen = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const list = ref<any[]>([])
const filterStatus = ref<string | undefined>(undefined)
const pagination = ref({ current: 1, pageSize: 10, total: 0 })
const parkOptions = ref<any[]>([])

const stats = ref({ active: 0, pending: 0, ended: 0, totalJoins: 0 })

const form = ref({
  title: '',
  description: '',
  start_time: undefined as Dayjs | undefined,
  end_time: undefined as Dayjs | undefined,
  park_id: undefined as number | undefined,
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '活动标题', dataIndex: 'title', key: 'title' },
  { title: '活动时间', key: 'time', width: 280 },
  { title: '状态', key: 'status', width: 90 },
  { title: '参与', dataIndex: 'joins', key: 'joins', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 280 },
]

function fmtDate(d: Dayjs | undefined): string {
  if (!d) return ''
  return d.format('YYYY-MM-DD HH:mm:ss')
}

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (filterStatus.value) params.status = filterStatus.value
    const res: any = await getActivityList(params)
    let data: any[] = []
    if (Array.isArray(res)) data = res
    else if (res?.list) {
      data = res.list
      pagination.value.total = res.total ?? res.list.length
    }
    list.value = data
    stats.value = {
      active: data.filter((r: any) => r.status === '进行中').length,
      pending: data.filter((r: any) => r.status === '待审核').length,
      ended: data.filter((r: any) => r.status === '已结束').length,
      totalJoins: data.reduce((s: number, r: any) => s + Number(r.joins ?? 0), 0),
    }
  } catch (err: any) {
    showToast(err.message || '加载活动列表失败', 'error')
  } finally {
    loading.value = false
  }
}

function onPageChange(pag: any) {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchList()
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { title: '', description: '', start_time: undefined, end_time: undefined, park_id: undefined }
  modalOpen.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editId.value = record.id
  form.value = {
    title: record.title,
    description: record.description ?? '',
    start_time: record.start_time ? undefined : undefined,
    end_time: record.end_time ? undefined : undefined,
    park_id: record.park_id,
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function handleSubmit() {
  if (!form.value.title) {
    showToast('请填写活动标题', 'error')
    return
  }
  submitting.value = true
  try {
    const body: any = {
      title: form.value.title,
      description: form.value.description,
      start_time: fmtDate(form.value.start_time),
      end_time: fmtDate(form.value.end_time),
      park_id: form.value.park_id,
    }
    if (isEdit.value && editId.value) {
      await updateActivity(editId.value, body)
      showToast('活动更新成功')
    } else {
      await createActivity(body)
      showToast('活动创建成功')
    }
    closeModal()
    fetchList()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteActivity(id)
    showToast('活动已删除')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

async function handleAudit(id: number) {
  try {
    await auditActivity(id)
    showToast('审核通过')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

async function handleReject(id: number) {
  try {
    // Reject by using update to change status
    await updateActivity(id, { status: '已驳回' })
    showToast('已驳回')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

async function handleEnd(id: number) {
  try {
    await endActivity(id)
    showToast('活动已结束')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

onMounted(async () => {
  try {
    const parks = await getParkList({ page_size: 100 })
    parkOptions.value = Array.isArray(parks) ? parks : (parks as any)?.list ?? []
  } catch (_) { /* ignore */ }
  fetchList()
})
</script>
