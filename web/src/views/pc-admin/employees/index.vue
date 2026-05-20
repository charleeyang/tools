<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">员工管理</div>
      <div class="page-subtitle">管理平台员工账号与岗位</div>
    </div>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="姓名" style="width: 160px" allow-clear @press-enter="fetchList" />
        <a-input v-model:value="phoneFilter" placeholder="手机号" style="width: 160px" allow-clear @press-enter="fetchList" />
        <a-select v-model:value="filterPark" placeholder="所属园区" style="width: 150px" allow-clear>
          <a-select-option v-for="p in parkOptions" :key="p.id" :value="p.id">{{ p.name }}</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
        <a-button @click="resetSearch">重置</a-button>
        <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 新增员工</a-button>
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
          <template v-if="column.key === 'park_id'">
            {{ getParkName(record.park_id) }}
          </template>
          <template v-else-if="column.key === 'position_id'">
            {{ getPositionName(record.position_id) }}
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === '在职' ? 'green' : 'default'">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确定删除该员工？" @confirm="handleDelete(record.id)">
                <a style="color: #ff4d4f">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑员工' : '新增员工'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="姓名" required>
          <a-input v-model:value="form.name" placeholder="请输入姓名" />
        </a-form-item>
        <a-form-item label="手机号" required>
          <a-input v-model:value="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="岗位" required>
          <a-select v-model:value="form.position_id" placeholder="选择岗位">
            <a-select-option v-for="p in positionOptions" :key="p.id" :value="p.id">{{ p.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="所属园区">
          <a-select v-model:value="form.park_id" placeholder="选择园区">
            <a-select-option v-for="p in parkOptions" :key="p.id" :value="p.id">{{ p.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="在职">在职</a-select-option>
            <a-select-option value="离职">离职</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { getEmployeeList, createEmployee, updateEmployee, deleteEmployee, getPositionList } from '@/api/modules/employees'
import { getParkList } from '@/api/modules/parks'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)
const submitting = ref(false)
const modalOpen = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const list = ref<any[]>([])
const keyword = ref('')
const phoneFilter = ref('')
const filterPark = ref<number | undefined>(undefined)
const pagination = ref({ current: 1, pageSize: 10, total: 0 })
const parkOptions = ref<any[]>([])
const positionOptions = ref<any[]>([])

const form = ref({
  name: '',
  phone: '',
  position_id: undefined as number | undefined,
  park_id: undefined as number | undefined,
  status: '在职',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '姓名', dataIndex: 'name', key: 'name', width: 110 },
  { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
  { title: '岗位', key: 'position_id', width: 130 },
  { title: '所属园区', key: 'park_id', width: 140 },
  { title: '状态', key: 'status', width: 90 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 120 },
]

function getParkName(id: number): string {
  const p = parkOptions.value.find((x: any) => x.id === id)
  return p?.name ?? '-'
}

function getPositionName(id: number): string {
  const p = positionOptions.value.find((x: any) => x.id === id)
  return p?.name ?? '-'
}

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (phoneFilter.value) params.phone = phoneFilter.value
    if (filterPark.value) params.park_id = filterPark.value
    const res: any = await getEmployeeList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if (res?.list) {
      list.value = res.list
      pagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载员工列表失败', 'error')
  } finally {
    loading.value = false
  }
}

function onPageChange(pag: any) {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchList()
}

function resetSearch() {
  keyword.value = ''
  phoneFilter.value = ''
  filterPark.value = undefined
  pagination.value.current = 1
  fetchList()
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { name: '', phone: '', position_id: undefined, park_id: undefined, status: '在职' }
  modalOpen.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editId.value = record.id
  form.value = {
    name: record.name,
    phone: record.phone,
    position_id: record.position_id,
    park_id: record.park_id,
    status: record.status,
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function handleSubmit() {
  if (!form.value.name || !form.value.phone || !form.value.position_id) {
    showToast('请填写必填项', 'error')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateEmployee(editId.value, form.value)
      showToast('员工更新成功')
    } else {
      await createEmployee(form.value)
      showToast('员工创建成功')
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
    await deleteEmployee(id)
    showToast('员工已删除')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

async function loadSelectOptions() {
  try {
    const [parks, positions] = await Promise.all([
      getParkList({ page_size: 100 }),
      getPositionList(),
    ])
    parkOptions.value = Array.isArray(parks) ? parks : (parks as any)?.list ?? []
    positionOptions.value = Array.isArray(positions) ? positions : []
  } catch (_) { /* ignore */ }
}

onMounted(() => {
  loadSelectOptions()
  fetchList()
})
</script>
