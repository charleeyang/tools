<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">店铺管理</div>
      <div class="page-subtitle">管理所有商户/店铺</div>
    </div>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="店铺名称" style="width: 200px" allow-clear @press-enter="fetchList" />
        <a-select v-model:value="filterStatus" placeholder="状态" style="width: 120px" allow-clear>
          <a-select-option value="营业中">营业中</a-select-option>
          <a-select-option value="休息中">休息中</a-select-option>
          <a-select-option value="已关闭">已关闭</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
        <a-button @click="resetSearch">重置</a-button>
        <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 添加店铺</a-button>
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
          <template v-else-if="column.key === 'shop_type_id'">
            {{ getShopTypeName(record.shop_type_id) }}
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === '营业中' ? 'green' : 'default'">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确定删除该店铺？" @confirm="handleDelete(record.id)">
                <a style="color: #ff4d4f">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑店铺' : '添加店铺'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="店铺名称" required>
          <a-input v-model:value="form.name" placeholder="如：火车餐厅" />
        </a-form-item>
        <a-form-item label="所属园区" required>
          <a-select v-model:value="form.park_id" placeholder="选择园区">
            <a-select-option v-for="p in parkOptions" :key="p.id" :value="p.id">{{ p.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="店铺类型" required>
          <a-select v-model:value="form.shop_type_id" placeholder="选择类型">
            <a-select-option v-for="t in shopTypeOptions" :key="t.id" :value="t.id">{{ t.name }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="营业中">营业中</a-select-option>
            <a-select-option value="休息中">休息中</a-select-option>
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
import { getShopList, createShop, updateShop, deleteShop, getShopTypes } from '@/api/modules/shops'
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
const filterStatus = ref<string | undefined>(undefined)
const pagination = ref({ current: 1, pageSize: 10, total: 0 })
const parkOptions = ref<any[]>([])
const shopTypeOptions = ref<any[]>([])

const form = ref({
  name: '',
  park_id: undefined as number | undefined,
  shop_type_id: undefined as number | undefined,
  status: '营业中',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '店铺名称', dataIndex: 'name', key: 'name' },
  { title: '所属园区', key: 'park_id', width: 130 },
  { title: '店铺类型', key: 'shop_type_id', width: 120 },
  { title: '状态', key: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 150 },
]

function getParkName(id: number): string {
  const p = parkOptions.value.find((x: any) => x.id === id)
  return p?.name ?? '-'
}

function getShopTypeName(id: number): string {
  const t = shopTypeOptions.value.find((x: any) => x.id === id)
  return t?.name ?? '-'
}

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (filterStatus.value) params.status = filterStatus.value
    const res: any = await getShopList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if (res?.list) {
      list.value = res.list
      pagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载店铺列表失败', 'error')
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
  filterStatus.value = undefined
  pagination.value.current = 1
  fetchList()
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { name: '', park_id: undefined, shop_type_id: undefined, status: '营业中' }
  modalOpen.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editId.value = record.id
  form.value = {
    name: record.name,
    park_id: record.park_id,
    shop_type_id: record.shop_type_id,
    status: record.status,
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function handleSubmit() {
  if (!form.value.name || !form.value.park_id || !form.value.shop_type_id) {
    showToast('请填写必填项', 'error')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateShop(editId.value, form.value)
      showToast('店铺更新成功')
    } else {
      await createShop(form.value)
      showToast('店铺创建成功')
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
    await deleteShop(id)
    showToast('店铺已删除')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

async function loadSelectOptions() {
  try {
    const [parks, types] = await Promise.all([
      getParkList({ page_size: 100 }),
      getShopTypes(),
    ])
    parkOptions.value = Array.isArray(parks) ? parks : (parks as any)?.list ?? []
    shopTypeOptions.value = Array.isArray(types) ? types : []
  } catch (_) { /* ignore */ }
}

onMounted(() => {
  loadSelectOptions()
  fetchList()
})
</script>
