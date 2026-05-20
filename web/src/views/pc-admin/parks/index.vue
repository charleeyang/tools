<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">园区管理</div>
      <div class="page-subtitle">管理平台旗下所有园区</div>
    </div>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="园区名称" style="width: 200px" allow-clear @press-enter="fetchList" />
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
        <a-button @click="resetSearch">重置</a-button>
        <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 新增园区</a-button>
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
            <a-tag :color="record.status === '营业中' ? 'green' : record.status === '筹备中' ? 'orange' : 'default'">
              {{ record.status }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确定删除该园区？" @confirm="handleDelete(record.id)">
                <a style="color: #ff4d4f">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑园区' : '新增园区'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="园区名称" required>
          <a-input v-model:value="form.name" placeholder="如：黄梅袁夫稻田" />
        </a-form-item>
        <a-form-item label="地址" required>
          <a-input v-model:value="form.address" placeholder="请输入详细地址" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="营业中">营业中</a-select-option>
            <a-select-option value="筹备中">筹备中</a-select-option>
            <a-select-option value="已关闭">已关闭</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { getParkList, createPark, updatePark, deletePark, type Park } from '@/api/modules/parks'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)
const submitting = ref(false)
const modalOpen = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const list = ref<Park[]>([])
const keyword = ref('')
const pagination = ref({ current: 1, pageSize: 10, total: 0 })

const form = ref<Partial<Park>>({ name: '', address: '', status: '营业中' })

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '园区名称', dataIndex: 'name', key: 'name' },
  { title: '地址', dataIndex: 'address', key: 'address' },
  { title: '状态', key: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    const res = await getParkList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if ((res as any)?.list) {
      list.value = (res as any).list
      pagination.value.total = (res as any).total ?? (res as any).list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载园区列表失败', 'error')
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
  pagination.value.current = 1
  fetchList()
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { name: '', address: '', status: '营业中' }
  modalOpen.value = true
}

function openEdit(record: Park) {
  isEdit.value = true
  editId.value = record.id
  form.value = { name: record.name, address: record.address, status: record.status }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  form.value = { name: '', address: '', status: '营业中' }
}

async function handleSubmit() {
  if (!form.value.name || !form.value.address) {
    showToast('请填写园区名称和地址', 'error')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updatePark(editId.value, form.value)
      showToast('园区更新成功')
    } else {
      await createPark(form.value)
      showToast('园区创建成功')
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
    await deletePark(id)
    showToast('园区已删除')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

onMounted(() => { fetchList() })
</script>
