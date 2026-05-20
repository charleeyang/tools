<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">园区总览</div>
      <div class="page-subtitle">全局经营概况</div>
    </div>

    <a-spin :spinning="loading">
      <!-- 第一行统计卡片 -->
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="总园区数" :value="overview.totalParks ?? 0" :value-style="{ color: '#1677ff' }">
              <template #prefix><HomeOutlined style="font-size:18px;color:#1677ff" /></template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="总商户数" :value="overview.totalShops ?? 0" :value-style="{ color: '#52c41a' }">
              <template #prefix><ShopOutlined style="font-size:18px;color:#52c41a" /></template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="总用户数" :value="overview.totalUsers ?? 0" :value-style="{ color: '#722ed1' }">
              <template #prefix><UserOutlined style="font-size:18px;color:#722ed1" /></template>
            </a-statistic>
            <div style="font-size:12px;color:#52c41a;margin-top:4px">
              <ArrowUpOutlined /> +{{ overview.todayNewUsers ?? 0 }} 今日新增
            </div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日总销售额" :value="overview.todaySales ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#fa8c16' }">
              <template #prefix><DollarOutlined style="font-size:18px;color:#fa8c16" /></template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <!-- 第二行统计卡片 -->
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日订单数" :value="overview.todayOrders ?? 0" :value-style="{ color: '#13c2c2' }">
              <template #prefix><FileTextOutlined style="font-size:18px;color:#13c2c2" /></template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="今日客流量" :value="overview.todayFlow ?? 0" :value-style="{ color: '#eb2f96' }">
              <template #prefix><LineChartOutlined style="font-size:18px;color:#eb2f96" /></template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="储值余额池" :value="overview.totalBalance ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#1677ff' }">
              <template #prefix><WalletOutlined style="font-size:18px;color:#1677ff" /></template>
            </a-statistic>
            <div style="font-size:12px;color:#999">累计在用户账户</div>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic title="待退款金额" :value="overview.pendingRefund ?? 0" prefix="¥" :precision="2" :value-style="{ color: '#ff4d4f' }">
              <template #prefix><UndoOutlined style="font-size:18px;color:#ff4d4f" /></template>
            </a-statistic>
            <div style="font-size:12px;color:#999">{{ overview.pendingRefundCount ?? 0 }} 笔待处理</div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 快捷入口 -->
      <a-card title="快捷入口" style="margin-bottom: 16px">
        <a-space wrap>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/parks')"><BuildOutlined /> 园区管理</a-button>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/shops')"><ShopOutlined /> 店铺管理</a-button>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/orders')"><FileTextOutlined /> 订单管理</a-button>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/users')"><UserOutlined /> 客户管理</a-button>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/finance')"><DollarOutlined /> 财务管理</a-button>
          <a-button type="primary" ghost @click="$router.push('/pc-admin/products')"><TagOutlined /> 商品管理</a-button>
        </a-space>
      </a-card>

      <!-- 实时业务动态 -->
      <a-card title="实时业务动态">
        <a-table :columns="recentColumns" :data-source="recentRecords" row-key="id" :pagination="false" size="small" />
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  HomeOutlined, ShopOutlined, UserOutlined, DollarOutlined,
  FileTextOutlined, LineChartOutlined, WalletOutlined, UndoOutlined,
  BuildOutlined, TagOutlined, ArrowUpOutlined,
} from '@ant-design/icons-vue'
import { api } from '@/api/client'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const loading = ref(false)

interface OverviewData {
  totalParks?: number
  totalShops?: number
  totalUsers?: number
  todaySales?: number
  todayOrders?: number
  todayFlow?: number
  totalBalance?: number
  pendingRefund?: number
  pendingRefundCount?: number
  todayNewUsers?: number
}

const overview = ref<OverviewData>({})

const recentColumns = [
  { title: '时间', dataIndex: 'time', width: 80 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '详情', dataIndex: 'detail' },
  { title: '金额', dataIndex: 'amount', width: 120, align: 'right' as const },
  { title: '状态', dataIndex: 'status', width: 100 },
]

const recentRecords = ref<any[]>([])

async function fetchOverview() {
  loading.value = true
  try {
    const data = await api.get<OverviewData>('/api/statistics/overview')
    if (data) overview.value = data
  } catch (err: any) {
    showToast(err.message || '加载概览数据失败', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOverview()
  recentRecords.value = [
    { id: 1, time: '14:22', type: '订单', detail: '稻田守望者 在 火车餐厅 下单', amount: '¥86.50', status: '已支付' },
    { id: 2, time: '13:15', type: '核销', detail: '田园生活家 在 树下咖啡 核销', amount: '¥32.00', status: '已核销' },
    { id: 3, time: '11:40', type: '退款', detail: '游客小新 申请订单退款', amount: '¥22.00', status: '待审核' },
    { id: 4, time: '10:22', type: '充值', detail: '稻田守望者 充值（赠 ¥30）', amount: '+¥200.00', status: '成功' },
    { id: 5, time: '09:12', type: '余额退款', detail: '稻田守望者 申请余额退款', amount: '¥50.00', status: '已退回' },
  ]
})
</script>
