<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">财务管理</div>
      <div class="page-subtitle">商户提现审核、资金管理</div>
    </div>

    <!-- 财务概览 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic title="总收入" :value="financeSummary.totalRevenue ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#1677ff' }">
            <template #prefix><DollarOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="账户余额" :value="financeSummary.balance ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#52c41a' }">
            <template #prefix><WalletOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="已提现" :value="financeSummary.withdrawn ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#722ed1' }">
            <template #prefix><ArrowUpOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="待提现" :value="financeSummary.pending ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#fa8c16' }">
            <template #prefix><ClockCircleOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="退款单号 / 店铺名" style="width: 200px" allow-clear />
        <a-select v-model:value="filterStatus" placeholder="状态" style="width: 120px" allow-clear>
          <a-select-option value="pending">待审核</a-select-option>
          <a-select-option value="approved">已通过</a-select-option>
          <a-select-option value="rejected">已驳回</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
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
          <template v-if="column.key === 'amount'">
            <span style="font-weight: 600">¥{{ Number(record.amount).toFixed(2) }}</span>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space v-if="record.status === 'pending'">
              <a @click="handleApprove(record)">通过</a>
              <a style="color: #ff4d4f" @click="handleReject(record)">驳回</a>
            </a-space>
            <span v-else-if="record.status === 'approved' || record.status === '已通过'" style="color: #999">已处理</span>
            <span v-else-if="record.status === 'rejected' || record.status === '已驳回'" style="color: #ff4d4f">已驳回</span>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 驳回弹窗 -->
    <a-modal
      v-model:open="rejectOpen"
      title="驳回提现申请"
      :confirm-loading="rejecting"
      @ok="confirmReject"
      @cancel="rejectOpen = false"
    >
      <a-form layout="vertical">
        <a-form-item label="驳回原因">
          <a-textarea v-model:value="rejectReason" :rows="3" placeholder="请输入驳回原因" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, DollarOutlined, WalletOutlined, ArrowUpOutlined, ClockCircleOutlined } from '@ant-design/icons-vue'
import { getFinanceSummary, getWithdrawList, approveWithdraw, rejectWithdraw } from '@/api/modules/finance'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)
const list = ref<any[]>([])
const keyword = ref('')
const filterStatus = ref<string | undefined>(undefined)
const pagination = ref({ current: 1, pageSize: 10, total: 0 })
const financeSummary = ref({ totalRevenue: 0, balance: 0, withdrawn: 0, pending: 0 })

const rejectOpen = ref(false)
const rejecting = ref(false)
const rejectTarget = ref<any>(null)
const rejectReason = ref('')

const columns = [
  { title: '申请编号', dataIndex: 'id', key: 'id', width: 140 },
  { title: '店铺名称', dataIndex: 'shop_name', key: 'shop_name', width: 130 },
  { title: '提现金额', key: 'amount', width: 120, align: 'right' as const },
  { title: '开户银行', dataIndex: 'bank_name', key: 'bank_name', width: 120 },
  { title: '银行账号', dataIndex: 'bank_account', key: 'bank_account', width: 150 },
  { title: '开户姓名', dataIndex: 'bank_holder', key: 'bank_holder', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '申请时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 120 },
]

function statusColor(s: string): string {
  if (s === 'approved' || s === '已通过') return 'green'
  if (s === 'rejected' || s === '已驳回') return 'red'
  if (s === 'pending') return 'orange'
  return 'default'
}

function statusLabel(s: string): string {
  if (s === 'approved') return '已通过'
  if (s === 'rejected') return '已驳回'
  if (s === 'pending') return '待审核'
  return s
}

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (filterStatus.value) params.status = filterStatus.value
    const res: any = await getWithdrawList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if (res?.list) {
      list.value = res.list
      pagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载提现列表失败', 'error')
  } finally {
    loading.value = false
  }
}

function onPageChange(pag: any) {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchList()
}

async function fetchSummary() {
  try {
    const res: any = await getFinanceSummary()
    if (res) financeSummary.value = res
  } catch (_) { /* ignore */ }
}

async function handleApprove(record: any) {
  try {
    await approveWithdraw(record.id)
    showToast('提现已通过')
    fetchList()
    fetchSummary()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

function handleReject(record: any) {
  rejectTarget.value = record
  rejectReason.value = ''
  rejectOpen.value = true
}

async function confirmReject() {
  if (!rejectTarget.value) return
  rejecting.value = true
  try {
    await rejectWithdraw(rejectTarget.value.id)
    showToast('提现已驳回')
    rejectOpen.value = false
    fetchList()
    fetchSummary()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  } finally {
    rejecting.value = false
  }
}

onMounted(() => {
  fetchList()
  fetchSummary()
})
</script>
