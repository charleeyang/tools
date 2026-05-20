<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">客户管理</div>
      <div class="page-subtitle">管理 C 端注册客户、查看余额与消费</div>
    </div>

    <!-- 客户统计 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card><a-statistic title="总用户数" :value="userStats.totalUsers" :value-style="{ color: '#1677ff' }"><template #prefix><UserOutlined /></template></a-statistic></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日新增" :value="userStats.todayNew" :value-style="{ color: '#52c41a' }"><template #prefix><PlusOutlined /></template></a-statistic></a-card>
      </a-col>
      <a-col :span="6">
        <a-card><a-statistic title="今日活跃" :value="userStats.todayActive" :value-style="{ color: '#722ed1' }"><template #prefix><AimOutlined /></template></a-statistic></a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic title="储值余额池" :value="userStats.totalBalance" prefix="¥" :precision="2" :value-style="{ color: '#fa8c16' }">
            <template #prefix><WalletOutlined /></template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-tabs v-model:activeKey="activeTab">
      <!-- 客户列表 -->
      <a-tab-pane key="users" tab="客户列表">
        <a-card>
          <a-space style="margin-bottom: 16px" wrap>
            <a-input v-model:value="keyword" placeholder="客户昵称" style="width: 180px" allow-clear @press-enter="fetchList" />
            <a-input v-model:value="phone" placeholder="手机号" style="width: 160px" allow-clear @press-enter="fetchList" />
            <a-button type="primary" @click="fetchList"><SearchOutlined /> 查询</a-button>
            <a-button @click="resetSearch">重置</a-button>
            <a-button type="primary" style="margin-left: auto" @click="openCreate"><PlusOutlined /> 手动添加</a-button>
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
              <template v-if="column.key === 'balance'">
                <span :style="{ color: Number(record.balance) > 0 ? '#52c41a' : undefined, fontWeight: '600' }">
                  ¥{{ Number(record.balance).toFixed(2) }}
                </span>
              </template>
              <template v-else-if="column.key === 'member_level'">
                <a-tag :color="record.member_level === 'VIP3' ? 'red' : record.member_level === 'VIP2' ? 'blue' : record.member_level === 'VIP1' ? 'green' : 'default'">
                  {{ record.member_level }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'action'">
                <a-space>
                  <a @click="openEdit(record)">编辑</a>
                  <a @click="openRecharge(record)">充值</a>
                  <a-popconfirm title="确定删除该用户？" @confirm="handleDelete(record.id)">
                    <a style="color: #ff4d4f">删除</a>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <!-- 充值记录 -->
      <a-tab-pane key="recharges" tab="充值记录">
        <a-card>
          <a-table
            :columns="rechargeColumns"
            :data-source="rechargeList"
            :loading="rechargeLoading"
            row-key="id"
            :pagination="rechargePagination"
            @change="onRechargePageChange"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'amount'">
                <span :style="{ color: Number(record.amount) >= 0 ? '#52c41a' : '#ff4d4f', fontWeight: '600' }">
                  {{ Number(record.amount) >= 0 ? '+' : '' }}¥{{ Math.abs(Number(record.amount)).toFixed(2) }}
                </span>
              </template>
              <template v-else-if="column.key === 'status'">
                <a-tag :color="record.status === '成功' ? 'green' : record.status === '已退款' ? 'red' : 'default'">
                  {{ record.status }}
                </a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>
    </a-tabs>

    <!-- 用户编辑弹窗 -->
    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑用户' : '添加用户'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="closeModal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="用户昵称" required>
          <a-input v-model:value="form.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model:value="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="会员等级">
          <a-select v-model:value="form.member_level">
            <a-select-option value="普通用户">普通用户</a-select-option>
            <a-select-option value="VIP1">VIP1</a-select-option>
            <a-select-option value="VIP2">VIP2</a-select-option>
            <a-select-option value="VIP3">VIP3</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 充值弹窗 -->
    <a-modal
      v-model:open="rechargeOpen"
      title="余额充值"
      :confirm-loading="rechargeSubmitting"
      @ok="handleRecharge"
      @cancel="rechargeOpen = false"
    >
      <a-form :model="rechargeForm" layout="vertical">
        <a-form-item label="用户">
          <span style="font-weight: 600">{{ rechargeTarget?.nickname ?? rechargeTarget?.name }}</span>
          <span style="color: #999; margin-left: 8px">当前余额：¥{{ Number(rechargeTarget?.balance ?? 0).toFixed(2) }}</span>
        </a-form-item>
        <a-form-item label="充值金额" required>
          <a-input-number v-model:value="rechargeForm.amount" :min="0" :precision="2" style="width: 100%" placeholder="请输入充值金额" />
        </a-form-item>
        <a-form-item label="赠送金额">
          <a-input-number v-model:value="rechargeForm.gift" :min="0" :precision="2" style="width: 100%" placeholder="赠送金额（可选）" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SearchOutlined, PlusOutlined, UserOutlined, AimOutlined, WalletOutlined } from '@ant-design/icons-vue'
import {
  getUserList, getUserStats, createUser, updateUser, deleteUser, getRechargeList,
} from '@/api/modules/users'
import { useToast } from '@/composables/useToast'

const { showToast } = useToast()
const activeTab = ref('users')
const loading = ref(false)
const list = ref<any[]>([])
const keyword = ref('')
const phone = ref('')
const pagination = ref({ current: 1, pageSize: 10, total: 0 })

const userStats = ref({ totalUsers: 0, todayNew: 0, todayActive: 0, totalBalance: 0 })

// Edit modal
const modalOpen = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitting = ref(false)
const form = ref({ nickname: '', phone: '', member_level: '普通用户' })

// Recharge modal
const rechargeOpen = ref(false)
const rechargeSubmitting = ref(false)
const rechargeTarget = ref<any>(null)
const rechargeForm = ref({ amount: 0, gift: 0 })

// Recharge list
const rechargeLoading = ref(false)
const rechargeList = ref<any[]>([])
const rechargePagination = ref({ current: 1, pageSize: 10, total: 0 })

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '昵称', dataIndex: 'nickname', key: 'nickname', width: 130 },
  { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
  { title: '余额', key: 'balance', width: 110, align: 'right' as const },
  { title: '会员等级', key: 'member_level', width: 100 },
  { title: '注册时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 180 },
]

const rechargeColumns = [
  { title: '充值单号', dataIndex: 'id', key: 'id', width: 160 },
  { title: '手机号', dataIndex: 'phone', key: 'phone', width: 130 },
  { title: '金额', key: 'amount', width: 110, align: 'right' as const },
  { title: '赠送', dataIndex: 'gift', key: 'gift', width: 80, align: 'right' as const },
  { title: '支付方式', dataIndex: 'pay_method', key: 'pay_method', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '备注', dataIndex: 'remark', key: 'remark' },
  { title: '时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
]

async function fetchList() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: pagination.value.current, page_size: pagination.value.pageSize }
    if (keyword.value) params.keyword = keyword.value
    if (phone.value) params.phone = phone.value
    const res: any = await getUserList(params)
    if (Array.isArray(res)) {
      list.value = res
      pagination.value.total = res.length
    } else if (res?.list) {
      list.value = res.list
      pagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载用户列表失败', 'error')
  } finally {
    loading.value = false
  }
}

function onPageChange(pag: any) {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchList()
}

async function fetchStats() {
  try {
    const res: any = await getUserStats()
    if (res) userStats.value = res
  } catch (_) { /* ignore */ }
}

async function fetchRecharges() {
  rechargeLoading.value = true
  try {
    const params: Record<string, any> = { page: rechargePagination.value.current, page_size: rechargePagination.value.pageSize }
    const res: any = await getRechargeList(params)
    if (Array.isArray(res)) {
      rechargeList.value = res
      rechargePagination.value.total = res.length
    } else if (res?.list) {
      rechargeList.value = res.list
      rechargePagination.value.total = res.total ?? res.list.length
    }
  } catch (err: any) {
    showToast(err.message || '加载充值记录失败', 'error')
  } finally {
    rechargeLoading.value = false
  }
}

function onRechargePageChange(pag: any) {
  rechargePagination.value.current = pag.current
  rechargePagination.value.pageSize = pag.pageSize
  fetchRecharges()
}

function resetSearch() {
  keyword.value = ''
  phone.value = ''
  pagination.value.current = 1
  fetchList()
}

function openCreate() {
  isEdit.value = false
  editId.value = null
  form.value = { nickname: '', phone: '', member_level: '普通用户' }
  modalOpen.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editId.value = record.id
  form.value = {
    nickname: record.nickname ?? record.name ?? '',
    phone: record.phone ?? '',
    member_level: record.member_level ?? '普通用户',
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

async function handleSubmit() {
  if (!form.value.nickname) {
    showToast('请填写昵称', 'error')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value && editId.value) {
      await updateUser(editId.value, form.value)
      showToast('用户更新成功')
    } else {
      await createUser(form.value)
      showToast('用户创建成功')
    }
    closeModal()
    fetchList()
    fetchStats()
  } catch (err: any) {
    showToast(err.message || '操作失败', 'error')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteUser(id)
    showToast('用户已删除')
    fetchList()
    fetchStats()
  } catch (err: any) {
    showToast(err.message || '删除失败', 'error')
  }
}

function openRecharge(record: any) {
  rechargeTarget.value = record
  rechargeForm.value = { amount: 0, gift: 0 }
  rechargeOpen.value = true
}

async function handleRecharge() {
  if (rechargeForm.value.amount <= 0) {
    showToast('请输入充值金额', 'error')
    return
  }
  rechargeSubmitting.value = true
  try {
    const { api } = await import('@/api/client')
    await api.post('/api/recharges', {
      user_id: rechargeTarget.value?.id,
      amount: rechargeForm.value.amount,
      gift: rechargeForm.value.gift,
    })
    showToast('充值成功')
    rechargeOpen.value = false
    fetchList()
    fetchRecharges()
  } catch (err: any) {
    showToast(err.message || '充值失败', 'error')
  } finally {
    rechargeSubmitting.value = false
  }
}

onMounted(() => {
  fetchList()
  fetchStats()
  fetchRecharges()
})
</script>
