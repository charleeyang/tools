<template>
  <div class="mp-admin-view">
    <!-- ======== dashboard (首页) ======== -->
    <div v-if="currentTab === 'dashboard'" class="tab-content">
      <div class="mp-data-hero">
        <div class="dh-title">{{ parkName }} · 2026-05-14</div>
        <div class="dh-subtitle">今日经营概览</div>
        <div class="dh-row">
          <div><div class="dh-val">¥{{ parkData.sales.toLocaleString() }}</div><div class="dh-label">今日销售额</div></div>
          <div><div class="dh-val">{{ parkData.orders }}</div><div class="dh-label">订单数</div></div>
          <div><div class="dh-val">{{ parkData.flow }}</div><div class="dh-label">客流量</div></div>
        </div>
      </div>

      <div class="mp-data-grid">
        <div class="dg-card"><div class="dg-val">{{ parkData.shops }}</div><div class="dg-label">商户总数</div></div>
        <div class="dg-card"><div class="dg-val">{{ parkData.users.toLocaleString() }}</div><div class="dg-label">累计用户</div></div>
        <div class="dg-card"><div class="dg-val" style="color:#ff4d4f">{{ parkData.refunds }}</div><div class="dg-label">待处理退款</div></div>
        <div class="dg-card"><div class="dg-val">¥{{ parkData.monthSales.toLocaleString() }}</div><div class="dg-label">本月销售额</div></div>
      </div>

      <div class="section-title">快捷功能</div>
      <div class="quick-actions">
        <div class="qa-item" @click="router.push('/mp-admin/scan')">
          <div class="qa-icon">⊞</div>扫码核销
        </div>
        <div class="qa-item" @click="showToast('充值查询 · ' + parkName)">
          <div class="qa-icon">¥</div>充值查询
        </div>
        <div class="qa-item" @click="showToast('订单管理 · ' + parkName)">
          <div class="qa-icon">📋</div>订单管理
        </div>
        <div class="qa-item" @click="showToast('客流查看 · ' + parkName)">
          <div class="qa-icon">📈</div>客流查看
        </div>
        <div class="qa-item" @click="router.push('/mp-admin/shops')">
          <div class="qa-icon">🏪</div>商户切换
        </div>
        <div class="qa-item" @click="showToast('活动审核 · ' + parkName)">
          <div class="qa-icon">🎁</div>活动审核
        </div>
        <div class="qa-item" @click="showToast('退款审核 · ' + parkName)">
          <div class="qa-icon">↩</div>退款审核
        </div>
        <div class="qa-item" @click="switchPark">
          <div class="qa-icon">⇄</div>切换园区
        </div>
      </div>

      <div class="section-title">今日动态</div>
      <div class="record-list">
        <div class="record-row">
          <div class="rr-icon consume">↓</div>
          <div class="rr-info"><div class="rr-title">稻田守望者 在 火车餐厅 消费</div><div class="rr-sub">14:22 · 已支付</div></div>
          <div class="rr-amount minus">-¥86.50</div>
        </div>
        <div class="record-row">
          <div class="rr-icon check">✓</div>
          <div class="rr-info"><div class="rr-title">田园生活家 完成核销</div><div class="rr-sub">13:15 · 树下咖啡</div></div>
          <div class="rr-amount minus">-¥32.00</div>
        </div>
        <div class="record-row">
          <div class="rr-icon refund">↩</div>
          <div class="rr-info"><div class="rr-title">游客小新 申请退款</div><div class="rr-sub">11:40 · 待审核</div></div>
          <div class="rr-amount" style="color:#fa8c16">¥22.00</div>
        </div>
        <div class="record-row">
          <div class="rr-icon plus">+</div>
          <div class="rr-info"><div class="rr-title">稻田守望者 充值</div><div class="rr-sub">10:22 · 赠 ¥30</div></div>
          <div class="rr-amount plus-amount">+¥200.00</div>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== shops (商户) ======== -->
    <div v-else-if="currentTab === 'shops'" class="tab-content">
      <div class="tab-filter-row">
        <button :class="['filter-btn', { active: shopFilter === 'all' }]" @click="shopFilter = 'all'">全部 ({{ shops.length }})</button>
        <button :class="['filter-btn', { active: shopFilter === 'active' }]" @click="shopFilter = 'active'">营业中 ({{ shops.filter(s => s.status === '营业中').length }})</button>
        <button :class="['filter-btn', { active: shopFilter === 'rest' }]" @click="shopFilter = 'rest'">休息中 ({{ shops.filter(s => s.status !== '营业中').length }})</button>
      </div>

      <div class="section-title">{{ parkName }} · 商户列表</div>
      <div class="shop-list">
        <div v-for="s in filteredShops" :key="s.name" class="shop-card" @click="showToast('已切换到' + s.name + '视角')">
          <div class="sc-img" :style="{ background: s.bg }">{{ s.img }}</div>
          <div class="sc-info">
            <div class="sc-name">{{ s.name }}</div>
            <div class="sc-desc">今日 ¥{{ s.sales }} · 订单 {{ s.orders }} · 客流 {{ s.flow }}</div>
            <span :class="['mp-tag', s.status === '营业中' ? 'green' : 'gray']">{{ s.status }}</span>
          </div>
          <span class="sc-arrow">›</span>
        </div>
      </div>

      <div class="section-title">商户排行 (本月)</div>
      <div class="rank-card">
        <div v-for="(s, i) in rankList" :key="s.n" class="rank-row">
          <span class="rank-num" :style="{ background: i === 0 ? '#ff4d4f' : i === 1 ? '#fa8c16' : '#999' }">{{ i + 1 }}</span>
          <span class="rank-name">{{ s.n }}</span>
          <div class="rank-bar-bg"><div class="rank-bar" :style="{ width: s.p + '%' }"></div></div>
          <span class="rank-val">¥{{ s.v.toLocaleString() }}</span>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== scan (扫码) ======== -->
    <div v-else-if="currentTab === 'scan'" class="tab-content">
      <div class="scan-mode-bar">
        <div class="scan-mode-tabs">
          <button :class="['sm-btn', { active: scanMode === 'verify' }]" @click="scanMode = 'verify'">✓ 订单核销</button>
          <button :class="['sm-btn', { active: scanMode === 'collect' }]" @click="scanMode = 'collect'">¥ 收款核销</button>
          <button :class="['sm-btn', { active: scanMode === 'meituan' }]" @click="scanMode = 'meituan'">☷ 美团验券</button>
          <button :class="['sm-btn', { active: scanMode === 'douyin' }]" @click="scanMode = 'douyin'">♫ 抖音验券</button>
        </div>
        <div class="scan-mode-desc">
          {{ scanMode === 'verify' ? '扫描用户出示的订单核销码完成核销' : scanMode === 'collect' ? '扫描用户付款码 + 手动输入金额完成收款' : scanMode === 'meituan' ? '扫描美团团购券二维码' : '扫描抖音团购券二维码' }}
        </div>
      </div>

      <!-- Scanner viewport per mode -->
      <div v-if="scanMode === 'meituan'" class="scan-viewport meituan-bg">
        <div class="scan-frame-wrap"><div class="scan-frame-icon">⊞</div></div>
        <div class="scan-frame-label">扫描美团团购券码</div>
        <div class="scan-frame-sub">支持美团/大众点评团购券</div>
      </div>
      <div v-else-if="scanMode === 'douyin'" class="scan-viewport douyin-bg">
        <div class="scan-frame-wrap"><div class="scan-frame-icon">⊞</div></div>
        <div class="scan-frame-label">扫描抖音团购券码</div>
        <div class="scan-frame-sub">支持套餐/次卡/代金券</div>
      </div>
      <div v-else class="scan-viewport default-bg">
        <div class="scan-frame-wrap scan-animated">
          <div class="scan-frame-icon">{{ scanMode === 'collect' ? '¥' : '⊞' }}</div>
          <div class="scan-line"></div>
        </div>
        <div class="scan-frame-label">{{ scanMode === 'collect' ? '扫描用户付款码' : '扫描用户核销码' }}</div>
      </div>

      <div class="scan-manual">
        <div class="manual-title">手动输入核销</div>
        <div class="manual-row">
          <input v-model="verifyCode" class="manual-input" :placeholder="scanMode === 'meituan' ? '输入美团券码' : scanMode === 'douyin' ? '输入抖音券码' : '输入核销码'" />
          <button class="mp-btn" @click="doVerify">核销</button>
        </div>
      </div>

      <div class="section-title">{{ parkName === '黄梅袁夫稻田' ? '黄梅' : '武汉' }} · 最近核销</div>
      <div class="record-list">
        <div class="record-row">
          <div class="rr-icon check">✓</div>
          <div class="rr-info"><div class="rr-title">{{ verifyRecords[0].user }} · {{ verifyRecords[0].item }}</div><div class="rr-sub">2026-05-14 {{ verifyRecords[0].time }} · 核销人 {{ verifyRecords[0].operator }}</div></div>
          <div class="rr-amount minus">{{ verifyRecords[0].amount }}</div>
        </div>
        <div class="record-row">
          <div class="rr-icon check">✓</div>
          <div class="rr-info"><div class="rr-title">{{ verifyRecords[1].user }} · {{ verifyRecords[1].item }}</div><div class="rr-sub">2026-05-14 {{ verifyRecords[1].time }} · 核销人 {{ verifyRecords[1].operator }}</div></div>
          <div class="rr-amount minus">{{ verifyRecords[1].amount }}</div>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== message (消息) ======== -->
    <div v-else-if="currentTab === 'message'" class="tab-content">
      <div class="section-title">通知中心</div>
      <div class="record-list">
        <div v-for="m in messages" :key="m.id" class="record-row">
          <div class="rr-icon" :style="{ background: m.iconBg, color: m.iconColor }">{{ m.icon }}</div>
          <div class="rr-info">
            <div class="rr-title">{{ m.title }}</div>
            <div class="rr-sub">{{ m.sub }}</div>
          </div>
          <span v-if="m.action === 'review'" class="mp-btn sm" @click="showToast('审核处理')">审核</span>
          <span v-else-if="m.action === 'view'" class="mp-btn sm outline" @click="showToast('查看详情')">查看</span>
          <span v-else-if="m.action === 'pending'" class="mp-tag orange">待处理</span>
          <span v-else class="mp-tag gray">已读</span>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== me (我的) ======== -->
    <div v-else-if="currentTab === 'me'" class="tab-content">
      <div class="profile-hero">
        <div class="pf-row">
          <div class="pf-avatar">{{ adminAvatar }}</div>
          <div class="pf-info">
            <div class="pf-name">{{ adminName }}</div>
            <div class="pf-role">园区管理员 · {{ parkName }}</div>
          </div>
          <button class="pf-action-btn" @click="showToast('权限管理')">权限管理</button>
        </div>
      </div>

      <div class="stat-row">
        <div class="stat-item">
          <div class="stat-val">{{ meStats.shops }}</div>
          <div class="stat-label">管辖商户</div>
        </div>
        <div class="stat-item">
          <div class="stat-val">{{ meStats.employees }}</div>
          <div class="stat-label">在岗员工</div>
        </div>
        <div class="stat-item">
          <div class="stat-val">99.4%</div>
          <div class="stat-label">SLA</div>
        </div>
      </div>

      <div class="menu-card">
        <div class="menu-row" @click="showToast('消息通知')"><span class="mm-icon">✉</span>消息通知<span class="mm-extra">3 条未读</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('角色与权限')"><span class="mm-icon">🔑</span>角色与权限<span class="mm-extra">查看可访问模块</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('移动端入口二维码')"><span class="mm-icon">📱</span>移动端入口二维码<span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="switchPark"><span class="mm-icon">⇄</span>切换园区<span class="mm-extra">{{ parkName }} › 切换</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('系统设置')"><span class="mm-icon">⚙</span>系统设置<span class="mm-arrow">›</span></div>
      </div>

      <div class="menu-card">
        <div class="menu-row" @click="showToast('关于我们')"><span class="mm-icon">ℹ</span>关于我们<span class="mm-arrow">›</span></div>
        <div class="menu-row logout" @click="handleLogout"><span class="mm-icon">🚪</span>退出登录</div>
      </div>
      <div style="height:12px"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()

const currentTab = computed(() => (route.meta.tab as string) || 'dashboard')
const currentPark = ref<'park-hm' | 'park-wh'>('park-hm')
const scanMode = ref<'verify' | 'collect' | 'meituan' | 'douyin'>('verify')
const verifyCode = ref('')
const shopFilter = ref<'all' | 'active' | 'rest'>('all')

const parkMap = {
  'park-hm': { name: '黄梅袁夫稻田', sales: 1286, orders: 33, flow: 128, shops: 3, users: 2837, refunds: 3, monthSales: 12340 },
  'park-wh': { name: '武汉袁夫稻田', sales: 667, orders: 18, flow: 52, shops: 2, users: 1045, refunds: 1, monthSales: 5680 },
} as const

const parkData = computed(() => parkMap[currentPark.value])
const parkName = computed(() => parkData.value.name)

const adminAvatar = computed(() => currentPark.value === 'park-hm' ? '李' : '周')
const adminName = computed(() => currentPark.value === 'park-hm' ? '李四' : '周八')

const meStats = computed(() => ({
  shops: currentPark.value === 'park-hm' ? 3 : 2,
  employees: currentPark.value === 'park-hm' ? 6 : 4,
}))

const hmShops = [
  { name: '火车餐厅', img: '🚂', bg: '#fff3e0', sales: 860, orders: 22, flow: 68, status: '营业中' },
  { name: '树下咖啡', img: '☕', bg: '#e8d5b7', sales: 326, orders: 9, flow: 42, status: '营业中' },
  { name: '稻田手作坊', img: '🧺', bg: '#fff0f6', sales: 100, orders: 2, flow: 18, status: '休息中' },
]
const whShops = [
  { name: '江夏火车厨', img: '🚂', bg: '#fff3e0', sales: 420, orders: 12, flow: 35, status: '营业中' },
  { name: '稻田鲜货铺', img: '🏪', bg: '#e8f5e9', sales: 247, orders: 6, flow: 17, status: '营业中' },
]

const shops = computed(() => currentPark.value === 'park-hm' ? hmShops : whShops)
const filteredShops = computed(() => {
  if (shopFilter.value === 'active') return shops.value.filter(s => s.status === '营业中')
  if (shopFilter.value === 'rest') return shops.value.filter(s => s.status !== '营业中')
  return shops.value
})

const rankList = [
  { n: '火车餐厅', v: 8420, p: 90 },
  { n: '树下咖啡', v: 4280, p: 46 },
  { n: '稻田手作坊', v: 640, p: 7 },
]

const verifyRecords = computed(() => {
  if (currentPark.value === 'park-hm') {
    return [
      { user: '小麦的麦', item: '火车特色套餐', amount: '¥128.00', time: '10:25', operator: '王五' },
      { user: '田园生活家', item: '精品咖啡 ×2', amount: '¥32.00', time: '13:15', operator: '赵六' },
    ]
  }
  return [
    { user: '周八', item: '江夏特色套餐', amount: '¥188.00', time: '10:25', operator: '周八' },
    { user: '李大伟', item: '鲜货礼盒', amount: '¥68.00', time: '13:15', operator: '孙七' },
  ]
})

const messages = [
  { id: 1, icon: '⚠', iconBg: '#fff2f0', iconColor: '#ff4d4f', title: '游客小新 · 申请订单退款', sub: '11:40 · 待您审核 ¥22.00', action: 'review' },
  { id: 2, icon: '🔔', iconBg: '#fff7e6', iconColor: '#fa8c16', title: '总部活动 · 夏日稻田音乐节', sub: '已推送，请审核后下发至商户', action: 'pending' },
  { id: 3, icon: '¥', iconBg: '#f0f9eb', iconColor: '#67c23a', title: '提现申请 · 火车餐厅', sub: '¥1,200.00 · 招商银行', action: 'view' },
  { id: 4, icon: 'ℹ', iconBg: '#e6f4ff', iconColor: '#555', title: '系统消息', sub: '闸机 #3 已离线超过 24 小时', action: 'read' },
  { id: 5, icon: '🎁', iconBg: '#e6fffb', iconColor: '#13c2c2', title: '营销中心', sub: '充值送好礼活动累计参与 420 人', action: 'read' },
]

function switchPark() {
  currentPark.value = currentPark.value === 'park-hm' ? 'park-wh' : 'park-hm'
  showToast('已切换到：' + parkName.value)
}

function doVerify() {
  if (!verifyCode.value) { showToast('请输入核销码'); return }
  showToast('核销成功 ✓ ' + verifyCode.value)
  verifyCode.value = ''
}

function handleLogout() {
  showToast('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.tab-content { background: #f5f5f5; min-height: 100%; }

/* Data Hero */
.mp-data-hero {
  background: linear-gradient(135deg, #1677ff, #0958d9);
  color: #fff; padding: 20px 16px 24px;
}
.dh-title { font-size: 12px; opacity: 0.85; }
.dh-subtitle { font-size: 18px; font-weight: 600; margin-top: 4px; }
.dh-row { display: flex; gap: 0; margin-top: 14px; }
.dh-row > div { flex: 1; text-align: center; border-right: 1px solid rgba(255,255,255,0.2); }
.dh-row > div:last-child { border-right: none; }
.dh-val { font-size: 24px; font-weight: 700; letter-spacing: -0.5px; }
.dh-label { font-size: 11px; opacity: 0.8; margin-top: 4px; }

/* Data Grid */
.mp-data-grid {
  display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px;
  margin: -12px 12px 12px; position: relative; z-index: 1;
}
.dg-card {
  background: #fff; border-radius: 10px; padding: 14px; text-align: center;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.dg-val { font-size: 20px; font-weight: 700; }
.dg-label { font-size: 11px; color: #999; margin-top: 4px; }

/* Section */
.section-title { font-size: 14px; font-weight: 600; padding: 14px 16px 10px; color: #333; }

/* Quick Actions */
.quick-actions {
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px;
  padding: 0 12px;
}
.qa-item {
  background: #fff; border-radius: 10px; padding: 14px 8px; text-align: center;
  font-size: 11px; color: #333; cursor: pointer;
}
.qa-icon { font-size: 24px; margin-bottom: 6px; }

/* Record List */
.record-list { margin: 0 12px; }
.record-row {
  display: flex; align-items: center; gap: 10px; padding: 12px 0;
  border-bottom: 1px solid #f0f0f0; background: #fff; padding: 12px 14px;
}
.record-row:first-child { border-radius: 10px 10px 0 0; }
.record-row:last-child { border-radius: 0 0 10px 10px; border-bottom: none; }
.record-row:only-child { border-radius: 10px; }
.rr-icon {
  width: 36px; height: 36px; border-radius: 50%; display: flex;
  align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0;
}
.rr-icon.consume { background: #e6f4ff; color: #1677ff; }
.rr-icon.check { background: #f0f9eb; color: #67c23a; }
.rr-icon.refund { background: #fff7e6; color: #fa8c16; }
.rr-icon.plus { background: #f0f9eb; color: #67c23a; }
.rr-info { flex: 1; min-width: 0; }
.rr-title { font-size: 13px; color: #333; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rr-sub { font-size: 11px; color: #999; margin-top: 2px; }
.rr-amount { font-size: 15px; font-weight: 600; flex-shrink: 0; }
.rr-amount.minus { color: #333; }
.rr-amount.plus-amount { color: #67c23a; }

/* Tags */
.mp-tag {
  display: inline-block; padding: 2px 8px; border-radius: 4px;
  font-size: 10px; font-weight: 500;
}
.mp-tag.green { background: #f0f9eb; color: #67c23a; }
.mp-tag.gray { background: #f5f5f5; color: #999; }
.mp-tag.red { background: #fff2f0; color: #ff4d4f; }
.mp-tag.orange { background: #fff7e6; color: #fa8c16; }

/* Buttons */
.mp-btn {
  background: #1677ff; color: #fff; border: none; border-radius: 8px;
  padding: 8px 20px; font-size: 13px; cursor: pointer; font-weight: 500;
}
.mp-btn.sm { padding: 4px 12px; font-size: 12px; }
.mp-btn.outline { background: #fff; color: #1677ff; border: 1px solid #1677ff; }
.mp-btn.gray { background: #f5f5f5; color: #666; }

/* Shop Tab */
.tab-filter-row {
  display: flex; gap: 4px; padding: 6px; margin: 10px 12px 0;
  background: #f0f0f0; border-radius: 8px;
}
.filter-btn {
  flex: 1; padding: 6px; border: none; border-radius: 6px;
  font-size: 12px; background: transparent; color: #666; cursor: pointer;
}
.filter-btn.active { background: #fff; color: #1677ff; font-weight: 500; box-shadow: 0 1px 3px rgba(0,0,0,0.06); }
.shop-list { margin: 0 12px; background: #fff; border-radius: 12px; overflow: hidden; }
.shop-card {
  display: flex; align-items: center; gap: 10px; padding: 12px 14px;
  border-bottom: 1px solid #f5f5f5; cursor: pointer;
}
.shop-card:last-child { border-bottom: none; }
.sc-img {
  width: 44px; height: 44px; border-radius: 10px; display: flex;
  align-items: center; justify-content: center; font-size: 24px; flex-shrink: 0;
}
.sc-info { flex: 1; min-width: 0; }
.sc-name { font-size: 14px; font-weight: 500; }
.sc-desc { font-size: 11px; color: #999; margin-top: 2px; }
.sc-arrow { font-size: 18px; color: #ccc; }

/* Rank */
.rank-card { margin: 0 12px; background: #fff; border-radius: 10px; padding: 14px 16px; }
.rank-row { display: flex; align-items: center; gap: 10px; padding: 8px 0; }
.rank-num {
  width: 18px; height: 18px; border-radius: 4px; color: #fff;
  font-size: 11px; display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.rank-name { flex: 1; font-size: 14px; }
.rank-bar-bg {
  width: 80px; height: 6px; background: #f5f5f5; border-radius: 3px; overflow: hidden;
}
.rank-bar { height: 100%; background: linear-gradient(90deg, #ff4d4f, #fa8c16); border-radius: 3px; }
.rank-val { font-weight: 600; font-size: 14px; width: 64px; text-align: right; }

/* Scan Tab */
.scan-mode-bar { background: #fff; padding: 12px 16px; border-bottom: 1px solid #f0f0f0; }
.scan-mode-tabs {
  display: flex; gap: 4px; padding: 3px; background: #f0f0f0; border-radius: 8px;
}
.sm-btn {
  flex: 1; min-width: 72px; padding: 7px 6px; border: none; border-radius: 6px;
  font-size: 12px; cursor: pointer; background: transparent; color: #666;
}
.sm-btn.active { background: #fff; color: #ff4d4f; font-weight: 500; box-shadow: 0 1px 4px rgba(0,0,0,0.06); }
.scan-mode-desc { font-size: 11px; color: #bbb; margin-top: 6px; text-align: center; }

.scan-viewport { padding: 20px 16px; text-align: center; color: #fff; }
.scan-viewport.default-bg { background: linear-gradient(180deg, #1677ff, #0958d9); }
.scan-viewport.meituan-bg { background: linear-gradient(180deg, #fa8c16, #d46b08); }
.scan-viewport.douyin-bg { background: linear-gradient(180deg, #010101, #1a1a1a); }
.scan-frame-wrap {
  width: 180px; height: 180px; border: 3px solid rgba(255,255,255,0.5); border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  background: rgba(255,255,255,0.05); margin: 0 auto; position: relative;
}
.scan-frame-icon { font-size: 64px; opacity: 0.8; }
.scan-frame-label { margin-top: 18px; font-size: 15px; }
.scan-frame-sub { font-size: 12px; opacity: 0.7; margin-top: 4px; }

.scan-manual { padding: 14px 12px; }
.manual-title { font-size: 14px; font-weight: 500; margin-bottom: 10px; background: #fff; padding: 14px; border-radius: 10px 10px 0 0; }
.manual-row { display: flex; gap: 8px; background: #fff; padding: 0 14px 14px; border-radius: 0 0 10px 10px; }
.manual-input {
  flex: 1; height: 38px; padding: 8px 12px; border: 1px solid #e0e0e0;
  border-radius: 8px; font-size: 14px; font-family: monospace; outline: none;
}

/* Profile */
.profile-hero {
  background: linear-gradient(135deg, #1677ff, #0958d9);
  padding: 24px 16px 20px; color: #fff;
}
.pf-row { display: flex; align-items: center; gap: 12px; }
.pf-avatar {
  width: 48px; height: 48px; border-radius: 50%; background: rgba(255,255,255,0.95);
  color: #1677ff; display: flex; align-items: center; justify-content: center;
  font-size: 20px; font-weight: 600;
}
.pf-info { flex: 1; }
.pf-name { font-size: 18px; font-weight: 600; }
.pf-role { color: rgba(255,255,255,0.85); font-size: 13px; margin-top: 4px; }
.pf-action-btn {
  background: rgba(255,255,255,0.15); border: none; color: #fff;
  padding: 6px 14px; border-radius: 14px; font-size: 12px; cursor: pointer;
}

/* Stats Row */
.stat-row {
  display: grid; grid-template-columns: repeat(3, 1fr); background: #fff;
  border-radius: 12px; margin: 12px; padding: 14px 0;
}
.stat-item { text-align: center; border-right: 1px solid #f0f0f0; }
.stat-item:last-child { border-right: none; }
.stat-val { font-size: 18px; font-weight: 600; }
.stat-label { font-size: 12px; color: #999; margin-top: 2px; }

/* Menu */
.menu-card { background: #fff; border-radius: 12px; margin: 12px; overflow: hidden; }
.menu-row {
  display: flex; align-items: center; gap: 10px; padding: 14px 16px;
  font-size: 14px; cursor: pointer; border-bottom: 1px solid #f5f5f5;
}
.menu-row:last-child { border-bottom: none; }
.menu-row.logout { color: #ff4d4f; }
.mm-icon { font-size: 18px; width: 24px; text-align: center; flex-shrink: 0; }
.mm-extra { flex: 1; text-align: right; font-size: 12px; color: #999; }
.mm-arrow { color: #ccc; flex-shrink: 0; }

/* Animated scan line */
.scan-animated .scan-line {
  position: absolute; top: 0; left: 0; right: 0; height: 2px;
  background: linear-gradient(90deg, transparent, #fff, transparent);
  animation: scanline 2s linear infinite;
}
@keyframes scanline { 0%,100%{transform:translateY(0)} 50%{transform:translateY(168px)} }
</style>
