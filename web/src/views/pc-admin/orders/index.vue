<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">订单管理</div>
      <div class="page-subtitle">管理所有消费订单、查看支付状态与核销情况</div>
    </div>

    <a-card>
      <a-space style="margin-bottom: 16px" wrap>
        <a-input v-model:value="keyword" placeholder="订单号 / 用户昵称" style="width: 200px" allow-clear @press-enter="fetchOrders" />
        <a-select v-model:value="filterStatus" placeholder="状态" style="width: 130px" allow-clear>
          <a-select-option value="待支付">待支付</a-select-option>
          <a-select-option value="已支付">已支付</a-select-option>
          <a-select-option value="已核销">已核销</a-select-option>
          <a-select-option value="已退款">已退款</a-select-option>
          <a-select-option value="退款处理中">退款处理中</a-select-option>
        </a-select>
        <a-button type="primary" @click="fetchOrders"><SearchOutlined /> 查询</a-button>
        <a-button @click="resetSearch">重置</a-button>
      </a-space>

      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="orders" tab="订单列表">
          <a-table
            :columns="orderColumns"
            :data-source="orderList"
            :loading="loading"
            row-key="order_no"
            :pagination="orderPagination"
            @change="onOrderPageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="statusColor(record.status)">{{ record.status }}</a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                <span style="font-weight: 600">¥{{ Number(record.amount).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a v-if="record.status === '已支付'" @click="handleVerify(record.order_no)">核销</a>
                  <a
                    v-if="record.status === '已支付' || record.status === '已核销'"
                    style="color: #ff4d4f"
                    @click="handleRefund(record.order_no)"
                  >退款</a>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-tab-pane>

        <a-tab-pane key="refunds" tab="退款处理">
          <!-- 退款统计 -->
          <a-row :gutter="16" style="margin-bottom: 16px">
            <a-col :span="6">
              <a-card size="small">
                <a-statistic title="待审核" :value="refundStats.pending" suffix="笔" :value-style="{ color: '#faad14' }" />
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card size="small">
                <a-statistic title="已通过" :value="refundStats.approved" suffix="笔" :value-style="{ color: '#52c41a' }" />
              </a-card>
            </a-col>
            <a-col :span="6">
              <a-card size="small">
                <a-statistic title="已驳回" :value="refundStats.rejected" suffix="笔" :value-style="{ color: '#ff4d4f' }" />
              </a-card>
            </a-col>
          </a-row>

          <a-table
            :columns="refundColumns"
            :data-source="refundList"
            :loading="refundLoading"
            row-key="id"
            :pagination="refundPagination"
            @change="onRefundPageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === '已通过' ? 'green' : record.status === '已驳回' ? 'red' : 'orange'">
                  {{ record.status === '已通过' ? '已通过' : record.status === '已驳回' ? '已驳回' : '待审核' }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                <span style="font-weight: 600">¥{{ Number(record.amount).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'type'">
                <a-tag :color="record.type === '余额退款' ? 'purple' : 'blue'">{{ record.type }}</a-tag>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space v-if="record.status === 'pending' || record.status === '待审核'">
                  <a @click="handleApproveRefund(record.id)">通过</a>
                  <a style="color: #ff4d4f" @click="handleRejectRefund(record.id)">驳回</a>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import {
  getOrderList, verifyOrder, refundOrder,
  getRefundList, approveRefund, rejectRefund,
} from '@/api/modules/orders'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const activeTab = ref('orders')
const loading = ref(false)
const refundLoading = ref(false)

// Orders
const orderList = ref<any[]>([])
const keyword = ref('')
const filterStatus = ref<string | undefined>(undefined)
const orderPagination = ref({ current: 1, pageSize: 10, total: 0 })

const orderColumns = [
  { title: '订单号', dataIndex: 'order_no', key: 'order_no', width: 180 },
  { title: '用户', dataIndex: 'user_name', key: 'user_name', width: 120 },
  { title: '商户', dataIndex: 'shop_name', key: 'shop_name', width: 120 },
  { title: '金额', key: 'amount', width: 100, align: 'right' as const },
  { title: '支付方式', dataIndex: 'pay_method', key: 'pay_method', width: 100 },
  { title: '状态', key: 'status', width: 110 },
  { title: '下单时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 120 },
]

function statusColor(s: string): string {
  const map: Record<string, string> = {
    '已支付': 'blue', '已核销': 'green', '已退款': 'default',
    '退款处理中': 'orange', '待支付': 'orange',
  }
  return map[s] ?? 'default'
}

// Refunds
const refundList = ref<any[]>([])
const refundPagination = ref({ current: 1, pageSize: 10, total: 0 })
const refundStats = ref({ pending: 0, approved: 0, rejected: 0 })

const refundColumns = [
  { title: '退款单号', dataIndex: 'id', key: 'id', width: 140 },
  { title: '关联订单', dataIndex: 'order_no', key: 'order_no', width: 180 },
  { title: '用户', dataIndex: 'user_name', key: 'user_name', width: 120 },
  { title: '类型', key: 'type', width: 100 },
  { title: '退款金额', key: 'amount', width: 110, align: 'right' as const },
  { title: '退款原因', dataIndex: 'reason', key: 'reason' },
  { title: '状态', key: 'status', width: 100 },
  { title: '申请时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 120 },
]

async function fetchOrders() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: orderPagination.value.current, page_size: orderPagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (filterStatus.value) params.status = filterStatus.value
    const res: any = await getOrderList(params)
    if (Array.isArray(res)) {
      orderList.value = res
      orderPagination.value.total = res.length
    } else if (res?.list) {
      orderList.value = res.list
      orderPagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载订单列表失败', 'error')
  } finally {
    loading.value = false
  }
}

function onOrderPageChange(pag: any) {
  orderPagination.value.current = pag.current
  orderPagination.value.pageSize = pag.pageSize
  fetchOrders()
}

async function fetchRefunds() {
  refundLoading.value = true
  try {
    const params: Record<string, any> = { page: refundPagination.value.current, page_size: refundPagination.value.pageSize }
    const res: any = await getRefundList(params)
    let data: any[] = []
    if (Array.isArray(res)) data = res
    else if (res?.list) data = res.list
    refundList.value = data
    refundPagination.value.total = res?.total ?? data.length
    refundStats.value = {
      pending: data.filter((r: any) => r.status === 'pending' || r.status === '待审核').length,
      approved: data.filter((r: any) => r.status === 'approved' || r.status === '已通过').length,
      rejected: data.filter((r: any) => r.status === 'rejected' || r.status === '已驳回').length,
    }
  } catch (err: any) {
    showToast(err.message || '加载退款列表失败', 'error')
  } finally {
    refundLoading.value = false
  }
}

function onRefundPageChange(pag: any) {
  refundPagination.value.current = pag.current
  refundPagination.value.pageSize = pag.pageSize
  fetchRefunds()
}

function resetSearch() {
  keyword.value = ''
  filterStatus.value = undefined
  orderPagination.value.current = 1
  fetchOrders()
}

async function handleVerify(orderNo: string) {
  try {
    await verifyOrder(orderNo)
    showToast('核销成功')
    fetchOrders()
  } catch (err: any) {
    showToast(err.message || '核销失败', 'error')
  }
}

async function handleRefund(orderNo: string) {
  try {
    await refundOrder(orderNo)
    showToast('退款已发起')
    fetchOrders()
    fetchRefunds()
  } catch (err: any) {
    showToast(err.message || '退款失败', 'error')
  }
}

async function handleApproveRefund(id: number) {
  try {
    await approveRefund(id)
    showToast('退款已通过')
    fetchRefunds()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

async function handleRejectRefund(id: number) {
  try {
    await rejectRefund(id)
    showToast('退款已驳回')
    fetchRefunds()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  }
}

onMounted(() => {
  fetchOrders()
  fetchRefunds()
})
</script>
