<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">商品管理</div>
      <div class="page-subtitle">管理商品价格、库存、上下架</div>
    </div>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="商品名称" style="width: 200px" allow-clear @press-enter="fetchList" />
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
        <a-button @click="resetSearch">重置</a-button>
        <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 新增商品</a-button>
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
          <template v-if="column.key === 'price'">¥{{ Number(record.price).toFixed(2) }}</template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === '已上架' ? 'green' : 'default'">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确定删除该商品？" @confirm="handleDelete(record.id)">
                <a style="color: #ff4d4f">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑商品' : '新增商品'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="商品名称" required>
          <a-input v-model:value="form.name" placeholder="如：稻田双人套餐" />
        </a-form-item>
        <a-form-item label="售价" required>
          <a-input-number v-model:value="form.price" :min="0" :precision="2" style="width: 100%" placeholder="0.00" />
        </a-form-item>
        <a-form-item label="原价">
          <a-input-number v-model:value="form.original_price" :min="0" :precision="2" style="width: 100%" placeholder="0.00" />
        </a-form-item>
        <a-form-item label="库存" required>
          <a-input-number v-model:value="form.stock" :min="0" style="width: 100%" placeholder="0" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" :rows="3" placeholder="商品描述" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="已上架">已上架</a-select-option>
            <a-select-option value="已下架">已下架</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { getProductList, createProduct, updateProduct, deleteProduct } from '@/api/modules/products'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)
const submitting = ref(false)
const modalOpen = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const list = ref<any[]>([])
const keyword = ref('')
const pagination = ref({ current: 1, pageSize: 10, total: 0 })

const form = ref({
  name: '',
  price: 0,
  original_price: undefined as number | undefined,
  stock: 0,
  description: '',
  status: '已上架',
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '商品名称', dataIndex: 'name', key: 'name' },
  { title: '售价', key: 'price', width: 100, align: 'right' as const },
  { title: '库存', dataIndex: 'stock', key: 'stock', width: 80, align: 'right' as const },
  { title: '状态', key: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 150 },
]

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    const res: any = await getProductList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if (res?.list) {
      list.value = res.list
      pagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载商品列表失败', 'error')
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
  form.value = { name: '', price: 0, original_price: undefined, stock: 0, description: '', status: '已上架' }
  modalOpen.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editId.value = record.id
  form.value = {
    name: record.name,
    price: record.price,
    original_price: record.original_price,
    stock: record.stock,
    description: record.description ?? '',
    status: record.status,
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function handleSubmit() {
  if (!form.value.name) {
    showToast('请填写商品名称', 'error')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateProduct(editId.value, form.value)
      showToast('商品更新成功')
    } else {
      await createProduct(form.value)
      showToast('商品创建成功')
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
    await deleteProduct(id)
    showToast('商品已删除')
    fetchList()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

onMounted(() => { fetchList() })
</script>
