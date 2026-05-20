<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">第三方平台</div>
      <div class="page-subtitle">美团/抖音等平台核销记录、结算明细、店铺配置管理</div>
    </div>

    <a-tabs v-model:activeKey="activeTab">
      <!-- 核销记录 -->
      <a-tab-pane key="verify" tab="核销记录">
        <a-card>
          <a-space style="margin-bottom: 16px" wrap>
            <a-input v-model:value="verifyKeyword" placeholder="核销单号 / 券码" style="width: 200px" allow-clear @press-enter="fetchVerifyRecords" />
            <a-select v-model:value="verifyStatus" placeholder="状态" style="width: 120px" allow-clear>
              <a-select-option value="已核销">已核销</a-select-option>
              <a-select-option value="已撤销">已撤销</a-select-option>
            </a-select>
            <a-button type="primary" @click="fetchVerifyRecords"><SearchOutlined /> 查询</a-button>
            <a-space style="margin-left: auto">
              <a-button type="primary" ghost @click="handlePrepareVerify">验券准备</a-button>
              <a-button type="primary" ghost @click="handleExecuteVerify">执行验券</a-button>
            </a-space>
          </a-space>

          <a-table
            :columns="verifyColumns"
            :data-source="verifyRecords"
            :loading="verifyLoading"
            row-key="id"
            :pagination="verifyPagination"
            @change="onVerifyPageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === '已核销' ? 'green' : 'red'">{{ record.status }}</a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                <span style="font-weight: 600">¥{{ Number(record.amount ?? record.sale_price ?? 0).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'action'">
                <a v-if="record.status === '已核销'" style="color: #ff4d4f" @click="handleRevoke(record.id)">撤销</a>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 结算记录 -->
      <a-tab-pane key="settlement" tab="结算记录">
        <a-card>
          <a-table
            :columns="settlementColumns"
            :data-source="settlements"
            :loading="settlementLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'total'">
                <span style="font-weight: 600">¥{{ Number(record.total_amount ?? 0).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'commission'">
                <span style="color: #ff4d4f">-¥{{ Number(record.commission ?? 0).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'net'">
                <span style="font-weight: 600; color: #52c41a">¥{{ Number(record.net_amount ?? 0).toFixed(2) }}</span>
              </template>
              <template v-else-if="column.key === 'status'">
                <a-tag :color="record.status === '已结算' ? 'green' : 'orange'">{{ record.status }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 店铺配置 -->
      <a-tab-pane key="store" tab="店铺配置">
        <a-card>
          <a-table
            :columns="storeColumns"
            :data-source="storeConfigs"
            :loading="storeLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === '已绑定' ? 'green' : 'orange'">{{ record.status }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import {
  getVerifyRecords, getPlatformSettlements, getStoreConfigs,
  prepareVerify, executeVerify, revokeVerify,
} from '@/api/modules/platform'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const activeTab = ref('verify')

// Verify records
const verifyLoading = ref(false)
const verifyRecords = ref<any[]>([])
const verifyKeyword = ref('')
const verifyStatus = ref<string | undefined>(undefined)
const verifyPagination = ref({ current: 1, pageSize: 10, total: 0 })

const verifyColumns = [
  { title: '核销单号', dataIndex: 'id', key: 'id', width: 160 },
  { title: '券码', dataIndex: 'coupon_code', key: 'coupon_code', width: 160 },
  { title: '商品', dataIndex: 'product', key: 'product' },
  { title: '金额', key: 'amount', width: 100, align: 'right' as const },
  { title: '用户', dataIndex: 'user_name', key: 'user_name', width: 120 },
  { title: '店铺', dataIndex: 'shop_name', key: 'shop_name', width: 120 },
  { title: '状态', key: 'status', width: 90 },
  { title: '核销时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 80 },
]

// Settlements
const settlementLoading = ref(false)
const settlements = ref<any[]>([])
const settlementColumns = [
  { title: '结算单号', dataIndex: 'id', key: 'id', width: 160 },
  { title: '店铺', dataIndex: 'shop_name', key: 'shop_name', width: 120 },
  { title: '结算日期', dataIndex: 'settle_date', key: 'settle_date', width: 120 },
  { title: '订单数', dataIndex: 'order_count', key: 'order_count', width: 80 },
  { title: '总额', key: 'total', width: 110, align: 'right' as const },
  { title: '佣金', key: 'commission', width: 110, align: 'right' as const },
  { title: '净额', key: 'net', width: 110, align: 'right' as const },
  { title: '状态', key: 'status', width: 90 },
]

// Store configs
const storeLoading = ref(false)
const storeConfigs = ref<any[]>([])
const storeColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '店铺', dataIndex: 'shop_name', key: 'shop_name', width: 130 },
  { title: 'AppAuthToken', dataIndex: 'app_auth_token', key: 'app_auth_token', width: 200 },
  { title: '状态', key: 'status', width: 90 },
  { title: '绑定时间', dataIndex: 'bind_at', key: 'bind_at', width: 180 },
]

async function fetchVerifyRecords() {
  verifyLoading.value = true
  try {
    const params: Record<string, any> = { page: verifyPagination.value.current, page_size: verifyPagination.value.pageSize }
    if (verifyKeyword.value) params.keyword = verifyKeyword.value
    if (verifyStatus.value) params.status = verifyStatus.value
    const res: any = await getVerifyRecords(params)
    if (Array.isArray(res)) {
      verifyRecords.value = res
      verifyPagination.value.total = res.length
    } else if (res?.list) {
      verifyRecords.value = res.list
      verifyPagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载核销记录失败', 'error')
  } finally {
    verifyLoading.value = false
  }
}

function onVerifyPageChange(pag: any) {
  verifyPagination.value.current = pag.current
  verifyPagination.value.pageSize = pag.pageSize
  fetchVerifyRecords()
}

async function fetchSettlements() {
  settlementLoading.value = true
  try {
    const res: any = await getPlatformSettlements()
    settlements.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载结算记录失败', 'error')
  } finally {
    settlementLoading.value = false
  }
}

async function fetchStores() {
  storeLoading.value = true
  try {
    const res: any = await getStoreConfigs()
    storeConfigs.value = Array.isArray(res) ? res : res?.list ?? []
  } catch (err: any) {
    showToast(err.message || '加载店铺配置失败', 'error')
  } finally {
    storeLoading.value = false
  }
}

async function handlePrepareVerify() {
  try {
    await prepareVerify({})
    showToast('验券准备完成')
  } catch (err: any) {
    showToast(err.message || '验券准备失败', 'error')
  }
}

async function handleExecuteVerify() {
  try {
    await executeVerify({})
    showToast('验券执行完成')
    fetchVerifyRecords()
  } catch (err: any) {
    showToast(err.message || '执行验券失败', 'error')
  }
}

async function handleRevoke(id: number) {
  try {
    await revokeVerify(id)
    showToast('验券已撤销')
    fetchVerifyRecords()
  } catch (err: any) {
    showToast(err.message || '撤销失败', 'error')
  }
}

onMounted(() => {
  fetchVerifyRecords()
  fetchSettlements()
  fetchStores()
})
</script>
