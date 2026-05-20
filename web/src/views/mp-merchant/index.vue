<template>
  <div class="mp-merchant-view">
    <!-- ======== dashboard (经营) ======== -->
    <div v-if="currentTab === 'dashboard'" class="tab-content">
      <div class="mp-data-hero">
        <div class="dh-title">🚂 火车餐厅 · 2026-05-14</div>
        <div class="dh-subtitle">今日经营</div>
        <div class="dh-row">
          <div><div class="dh-val">¥860</div><div class="dh-label">营业额</div></div>
          <div><div class="dh-val">22</div><div class="dh-label">订单数</div></div>
          <div><div class="dh-val">68</div><div class="dh-label">客流量</div></div>
        </div>
      </div>

      <div class="overview-card">
        <div class="ov-title">今日订单概况</div>
        <div class="ov-grid">
          <div><div class="ov-val" style="color:#fa8c16">3</div><div class="ov-label">待支付</div></div>
          <div><div class="ov-val" style="color:#ff4d4f">5</div><div class="ov-label">待核销</div></div>
          <div><div class="ov-val" style="color:#67c23a">13</div><div class="ov-label">已完成</div></div>
          <div><div class="ov-val" style="color:#666">1</div><div class="ov-label">已退款</div></div>
        </div>
      </div>

      <div class="section-title">快捷功能</div>
      <div class="quick-actions">
        <div class="qa-item" @click="router.push('/mp-merchant/verify')">
          <div class="qa-icon">⊞</div>扫码核销
        </div>
        <div class="qa-item" @click="router.push('/mp-merchant/orders')">
          <div class="qa-icon">📋</div>订单管理
        </div>
        <div class="qa-item" @click="router.push('/mp-merchant/finance')">
          <div class="qa-icon">¥</div>提现申请
        </div>
        <div class="qa-item" @click="showMenuMode = true">
          <div class="qa-icon">🏪</div>菜单上架
        </div>
        <div class="qa-item" @click="showToast('营业时间已更新')">
          <div class="qa-icon">🕐</div>营业时间
        </div>
        <div class="qa-item" @click="showToast('活动管理')">
          <div class="qa-icon">🎁</div>活动
        </div>
        <div class="qa-item" @click="showBusinessMode = true">
          <div class="qa-icon">📊</div>经营数据
        </div>
        <div class="qa-item" @click="showToast('店员管理')">
          <div class="qa-icon">👥</div>店员管理
        </div>
      </div>

      <div class="section-title">本月趋势</div>
      <div class="trend-card">
        <div class="trend-bars">
          <div v-for="(v, i) in trendData" :key="i" class="trend-bar" :style="{ height: (v / 900 * 100) + '%' }" :title="'¥' + v"></div>
        </div>
        <div class="trend-labels">
          <span>5/1</span><span>5/7</span><span>5/14</span>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== orders (订单) ======== -->
    <div v-else-if="currentTab === 'orders'" class="tab-content">
      <div class="tab-filter-row">
        <button v-for="f in orderFilters" :key="f.k" :class="['filter-btn', { active: orderFilter === f.k }]" @click="orderFilter = f.k">{{ f.t }} ({{ f.n }})</button>
      </div>

      <div v-for="o in orderList" :key="o.no" class="order-card">
        <div class="oc-header">
          <div class="oc-avatar">{{ o.avatar }}</div>
          <div class="oc-user-info">
            <div class="oc-user">{{ o.user }}</div>
            <div class="oc-time">{{ o.time }} · {{ o.pay }}</div>
          </div>
          <span :class="['mp-tag', statusClass(o.status)]">{{ o.status }}</span>
        </div>
        <div class="oc-items">{{ o.items }}</div>
        <div class="oc-footer">
          <span class="oc-no">{{ o.no }}</span>
          <div class="oc-actions">
            <span class="oc-amount">¥{{ o.amount.toFixed(2) }}</span>
            <button v-if="o.status === '待核销'" class="mp-btn sm" @click="showToast('核销成功 ✓ ' + o.code)">核销</button>
            <button v-else-if="o.status === '已核销'" class="mp-btn sm outline" @click="showToast('订单详情')">详情</button>
            <button v-else-if="o.status === '退款处理中'" class="mp-btn sm outline" @click="showToast('退款详情')">查看</button>
            <button v-else-if="o.status === '待支付'" class="mp-btn sm outline" @click="showToast('已取消')">取消</button>
          </div>
        </div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- ======== verify (核销) ======== -->
    <div v-else-if="currentTab === 'verify'" class="tab-content">
      <div class="scan-mode-bar">
        <div class="scan-mode-tabs">
          <button :class="['sm-btn', { active: verifyMode === 'verify' }]" @click="verifyMode = 'verify'">✓ 订单核销</button>
          <button :class="['sm-btn', { active: verifyMode === 'collect' }]" @click="verifyMode = 'collect'">¥ 收款核销</button>
          <button :class="['sm-btn', { active: verifyMode === 'meituan' }]" @click="verifyMode = 'meituan'">☷ 美团验券</button>
          <button :class="['sm-btn', { active: verifyMode === 'douyin' }]" @click="verifyMode = 'douyin'">♫ 抖音验券</button>
        </div>
        <div class="scan-mode-desc">{{ verifyModeDesc }}</div>
      </div>

      <!-- Meituan flow -->
      <template v-if="verifyMode === 'meituan'">
        <template v-if="meituanStep === 'scan'">
          <div class="scan-viewport meituan-bg">
            <div class="sf-corners">
              <div class="sf-corner tl"></div><div class="sf-corner tr"></div>
              <div class="sf-corner bl"></div><div class="sf-corner br"></div>
            </div>
            <div class="scan-frame-wrap"><div class="scan-frame-icon">⊞</div><div class="scan-line-meituan"></div></div>
            <div style="display:flex;align-items:center;justify-content:center;gap:8px;margin-bottom:20px"><span style="font-size:28px">☷</span><span style="font-size:18px;font-weight:600">美团 · 大众点评</span></div>
            <div class="scan-frame-label">扫描美团/大众点评团购券码</div>
            <div class="scan-frame-sub">对准用户出示的团购券二维码 / 条形码</div>
          </div>
          <div class="scan-manual">
            <div class="manual-title">手动输入券码</div>
            <div class="manual-row">
              <input v-model="meituanCode" class="manual-input" placeholder="输入 12 位美团券码" />
              <button class="mp-btn" style="background:#fa8c16;color:#fff" @click="meituanPreVerify">验券</button>
            </div>
          </div>
          <!-- Recent meituan verifications -->
          <div class="recent-card">
            <div class="recent-title">最近美团团购核销</div>
            <div v-for="r in meituanRecent" :key="r.code" class="recent-row">
              <div class="recent-info"><div class="recent-name">{{ r.product }}</div><div class="recent-sub">券码 {{ r.code }} · {{ r.time }}</div></div>
              <span class="recent-price">¥{{ r.price }}</span>
              <span class="mp-tag" :class="r.tagClass">{{ r.status }}</span>
            </div>
          </div>
          <div class="info-notice">
            <b>美团验券说明</b><br>
            · 扫描用户展示的美团/大众点评团购券二维码<br>
            · 系统自动调用「验券准备」接口预检查券码状态<br>
            · 确认无误后「执行验券」完成核销 · 核销后可在 24 小时内发起撤销
          </div>
        </template>
        <template v-else-if="meituanStep === 'confirm'">
          <div class="verify-confirm-header">
            <span style="font-size:24px">☷</span><span style="font-size:16px;font-weight:600">美团 · 验券确认</span>
          </div>
          <div class="verify-detail-card">
            <div class="vdc-header">
              <div class="vdc-icon meituan">券</div>
              <div class="vdc-info"><div class="vdc-name">{{ meituanData.product }}</div><div class="vdc-sub">美团 · 大众点评团购券</div></div>
              <span class="mp-tag green">{{ meituanData.status }}</span>
            </div>
            <div class="vdc-details">
              <div class="vdc-row"><span>券码</span><span class="mono">{{ meituanData.code }}</span></div>
              <div class="vdc-row"><span>原价</span><span class="strike">¥{{ meituanData.originPrice.toFixed(2) }}</span></div>
              <div class="vdc-row"><span>团购价</span><span class="big-price" style="color:#fa8c16">¥{{ meituanData.salePrice.toFixed(2) }}</span></div>
              <div class="vdc-row"><span>有效期</span><span>{{ meituanData.validFrom }} ~ {{ meituanData.validTo }}</span></div>
            </div>
          </div>
          <div class="verify-detail-card">
            <div class="vdc-section-title">核销信息</div>
            <div class="vdc-row"><span>用户</span><span>{{ meituanData.user }} ({{ meituanData.phone }})</span></div>
            <div class="vdc-row"><span>核销店铺</span><span>火车餐厅</span></div>
            <div class="vdc-row"><span>核销人</span><span>王五</span></div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="meituanBack">返回</button>
            <button class="mp-btn" style="flex:2;background:#fa8c16;color:#fff" @click="meituanExecute">确认核销 · ¥{{ meituanData.salePrice.toFixed(2) }}</button>
          </div>
        </template>
        <template v-else-if="meituanStep === 'done'">
          <div class="done-view">
            <div class="done-icon meituan-done">✓</div>
            <div class="done-title">美团验券成功</div>
            <div class="done-sub">券码 {{ meituanData.code }} 已完成核销</div>
            <div class="done-card">
              <div class="dok-row"><span class="dok-label">核销单号</span><span class="mono">{{ doneOrderNo('MT') }}</span></div>
              <div class="dok-row"><span class="dok-label">商品</span><span>{{ meituanData.product }}</span></div>
              <div class="dok-row"><span class="dok-label">核销金额</span><span class="big-price" style="color:#fa8c16">¥{{ meituanData.salePrice.toFixed(2) }}</span></div>
              <div class="dok-row"><span class="dok-label">状态</span><span class="mp-tag green">已核销</span></div>
            </div>
            <div class="info-notice">
              <b>温馨提示</b><br>· 已核销的券码如需撤销，请在 24 小时内操作<br>· 核销数据将同步至美团平台进行结算
            </div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="meituanReset">继续验券</button>
            <button class="mp-btn" style="flex:1;background:#fa8c16;color:#fff" @click="router.push('/mp-merchant/orders')">查看订单</button>
          </div>
        </template>
      </template>

      <!-- Douyin flow -->
      <template v-else-if="verifyMode === 'douyin'">
        <template v-if="douyinStep === 'scan'">
          <div class="scan-viewport douyin-bg">
            <div class="sf-corners dy-corners">
              <div class="sf-corner dy tl"></div><div class="sf-corner dy tr"></div>
              <div class="sf-corner dy bl"></div><div class="sf-corner dy br"></div>
            </div>
            <div class="scan-frame-wrap"><div class="scan-frame-icon">⊞</div><div class="scan-line-douyin"></div></div>
            <div style="display:flex;align-items:center;justify-content:center;gap:8px;margin-bottom:20px"><span style="font-size:24px">♫</span><span style="font-size:16px;font-weight:600">抖音 · 生活服务</span></div>
            <div class="scan-frame-label">扫描抖音团购券码</div>
            <div class="scan-frame-sub">支持团购套餐 · 次卡 · 代金券 · 组合券包</div>
          </div>
          <div class="scan-manual">
            <div class="manual-title">手动输入券码 / 扫码宣券</div>
            <div class="manual-row">
              <input v-model="douyinCode" class="manual-input" placeholder="输入抖音券码" />
              <button class="mp-btn" style="background:#010101;color:#fff" @click="douyinPreVerify">宣券</button>
            </div>
            <div class="manual-hint">宣券：调用「验券准备」接口，查询券码可验状态与可验张数</div>
          </div>
          <div class="recent-card">
            <div class="recent-title">最近抖音团购核销</div>
            <div v-for="r in douyinRecent" :key="r.code" class="recent-row">
              <div class="recent-info"><div class="recent-name">{{ r.product }}</div><div class="recent-sub">券码 {{ r.code }} · {{ r.time }}</div></div>
              <span class="recent-price">¥{{ r.price }}</span>
              <span class="mp-tag" :class="r.tagClass">{{ r.status }}</span>
            </div>
          </div>
          <div class="info-notice dy-notice">
            <b>抖音验券说明</b><br>
            · 扫描用户展示的抖音团购券二维码<br>
            · 系统自动调用「验券准备(宣券)」接口查询券码信息<br>
            · 确认可验张数后，调用「执行验券」完成核销<br>
            · 次卡/组合券包支持多次核销<br>
            · 核销后可在 24 小时内发起撤销
          </div>
        </template>
        <template v-else-if="douyinStep === 'confirm'">
          <div class="verify-confirm-header">
            <span style="font-size:24px">♫</span><span style="font-size:16px;font-weight:600">抖音 · 验券确认</span>
            <span :class="['mp-tag', douyinData.couponType === '团购套餐' ? 'red' : douyinData.couponType === '代金券' ? 'orange' : 'green']" style="margin-left:8px">{{ douyinData.couponType }}</span>
          </div>
          <div class="verify-detail-card">
            <div class="vdc-header">
              <div class="vdc-icon douyin">♫</div>
              <div class="vdc-info"><div class="vdc-name">{{ douyinData.product }}</div><div class="vdc-sub">抖音 · 生活服务</div></div>
              <span class="mp-tag green">{{ douyinData.status }}</span>
            </div>
            <div class="vdc-details">
              <div class="vdc-row"><span>券码</span><span class="mono">{{ douyinData.code }}</span></div>
              <div class="vdc-row"><span>原价</span><span class="strike">¥{{ douyinData.originPrice.toFixed(2) }}</span></div>
              <div class="vdc-row"><span>团购价</span><span class="big-price">¥{{ douyinData.salePrice.toFixed(2) }}</span></div>
              <div class="vdc-row"><span>有效期</span><span>{{ douyinData.validFrom }} ~ {{ douyinData.validTo }}</span></div>
              <div v-if="douyinData.couponType === '次卡' || douyinData.couponType === '组合券包'" class="vdc-row"><span>核销进度</span><span class="fw600">{{ douyinData.verifyCount }} / {{ douyinData.maxCount }} 次</span></div>
            </div>
          </div>
          <div class="verify-detail-card">
            <div class="vdc-section-title">核销信息</div>
            <div v-if="douyinData.couponType === '次卡' || douyinData.couponType === '组合券包'" class="verify-count-row">
              <div class="vcr-label">本次核销张数</div>
              <div class="vcr-controls">
                <button class="vcr-btn" @click="douyinData.verifyCount = Math.max(1, douyinData.verifyCount - 1)">−</button>
                <span class="vcr-val">{{ douyinData.verifyCount }}</span>
                <button class="vcr-btn" @click="douyinData.verifyCount = Math.min(douyinData.maxCount, douyinData.verifyCount + 1)">+</button>
                <span class="vcr-remain">剩余可验 {{ douyinData.maxCount - douyinData.verifyCount }} 次</span>
              </div>
            </div>
            <div class="vdc-row"><span>用户</span><span>{{ douyinData.user }} ({{ douyinData.phone }})</span></div>
            <div class="vdc-row"><span>核销店铺</span><span>火车餐厅</span></div>
            <div class="vdc-row"><span>核销人</span><span>王五</span></div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="douyinBack">返回</button>
            <button class="mp-btn" style="flex:2;background:#010101;color:#fff" @click="douyinExecute">确认核销 · ¥{{ (douyinData.salePrice * douyinData.verifyCount).toFixed(2) }}</button>
          </div>
        </template>
        <template v-else-if="douyinStep === 'done'">
          <div class="done-view">
            <div class="done-icon douyin-done">✓</div>
            <div class="done-title">抖音验券成功</div>
            <div class="done-sub">券码 {{ douyinData.code }} 已完成核销 · {{ douyinData.couponType }}</div>
            <div class="done-card">
              <div class="dok-row"><span class="dok-label">核销单号</span><span class="mono">{{ doneOrderNo('DY') }}</span></div>
              <div class="dok-row"><span class="dok-label">商品</span><span>{{ douyinData.product }}</span></div>
              <div class="dok-row"><span class="dok-label">核销金额</span><span class="big-price">¥{{ (douyinData.salePrice * douyinData.verifyCount).toFixed(2) }}</span></div>
              <div class="dok-row"><span class="dok-label">状态</span><span class="mp-tag green">已核销</span></div>
            </div>
            <div class="info-notice dy-notice">
              <b>温馨提示</b><br>· 已核销的券码如需撤销，请在 24 小时内操作<br>· 核销数据将同步至抖音生活服务平台进行结算<br>· 用户可在抖音 App「我的订单」查看核销状态
            </div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="douyinReset">继续验券</button>
            <button class="mp-btn" style="flex:1;background:#010101;color:#fff" @click="router.push('/mp-merchant/orders')">查看订单</button>
          </div>
        </template>
      </template>

      <!-- Standard verify / collect -->
      <template v-else>
        <div class="scan-viewport default-bg">
          <div class="scan-frame-wrap scan-animated">
            <div class="scan-frame-icon">{{ verifyMode === 'collect' ? '¥' : '⊞' }}</div>
            <div class="scan-line"></div>
          </div>
          <div class="scan-frame-label">{{ verifyMode === 'collect' ? '扫描用户付款码' : '扫描用户核销码' }}</div>
          <div class="scan-frame-sub">{{ verifyMode === 'collect' ? '用户在「会员码」展示付款码' : '对准用户出示的核销二维码' }}</div>
        </div>

        <!-- Collect panel -->
        <template v-if="verifyMode === 'collect'">
          <div class="collect-panel">
            <div class="cp-user-row">
              <div class="cp-avatar">稻</div>
              <div class="cp-info">
                <div class="cp-name">稻田守望者 <span class="cp-phone">(187****9908)</span></div>
                <div class="cp-meta">VIP2 · 余额 ¥130.00 · 已识别</div>
              </div>
              <span class="mp-tag green">已识别</span>
            </div>
            <div class="cp-amount-section">
              <div class="cp-label">收款金额</div>
              <div class="cp-amount-input">
                <span class="cp-yuan">¥</span>
                <input v-model="collectAmount" type="text" inputmode="decimal" placeholder="0.00" class="cp-input" />
              </div>
              <div class="cp-presets">
                <button v-for="v in [10, 20, 50, 100, 200]" :key="v" class="cp-preset-btn" @click="collectAmount = String(v)">+¥{{ v }}</button>
                <button class="cp-preset-btn clear" @click="collectAmount = ''">清空</button>
              </div>
            </div>
            <div class="cp-remark">
              <div class="cp-label">备注（可选）</div>
              <input type="text" placeholder="如：堂食 / 外带 / 餐位号" class="cp-remark-input" />
            </div>
            <div class="cp-payment">
              <div class="cp-label">扣款方式</div>
              <label class="cp-radio checked">
                <span class="cp-radio-dot"></span>
                <span class="cp-radio-text">优先从余额扣款</span>
                <span class="cp-radio-extra">余额 ¥130.00</span>
              </label>
              <label class="cp-radio">
                <span class="cp-radio-ring"></span>
                <span class="cp-radio-text">余额不足时使用微信支付</span>
              </label>
            </div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="collectAmount = ''; verifyMode = 'verify'">取消</button>
            <button class="mp-btn" style="flex:2" @click="confirmCollect">确认收款{{ collectAmount ? ' ¥' + parseFloat(collectAmount).toFixed(2) : '' }}</button>
          </div>
        </template>

        <!-- Standard verify panel -->
        <template v-else>
          <div class="scan-manual">
            <div class="manual-title">手动核销</div>
            <div class="manual-row">
              <input v-model="verifyCode" class="manual-input" placeholder="输入核销码 (HX...)" />
              <button class="mp-btn" @click="doVerify">核销</button>
            </div>
          </div>
          <div class="section-title">今日已核销 (5)</div>
          <div class="record-list">
            <div v-for="r in verifiedToday" :key="r.code" class="record-row">
              <div class="rr-icon check">✓</div>
              <div class="rr-info"><div class="rr-title">{{ r.user }}</div><div class="rr-sub">{{ r.item }} · {{ r.time }} · {{ r.code }}</div></div>
              <div class="rr-amount minus">¥{{ r.amount }}</div>
            </div>
          </div>
        </template>
      </template>
      <div style="height:12px"></div>
    </div>

    <!-- ======== finance (收益) ======== -->
    <div v-else-if="currentTab === 'finance'" class="tab-content">
      <div class="finance-hero">
        <div class="fh-label">可用余额</div>
        <div class="fh-balance">¥ 2,180.50</div>
        <div class="fh-sub-row">
          <span>冻结金额 ¥0.00</span>
          <span>累计收益 ¥6,420.50</span>
        </div>
        <div class="fh-actions">
          <button class="fh-btn primary" @click="showWithdrawAccount = true">💰 申请提现</button>
          <button class="fh-btn" @click="showToast('收益明细')">📊 收益明细</button>
        </div>
      </div>

      <div class="mp-data-grid">
        <div class="dg-card"><div class="dg-val">¥860</div><div class="dg-label">今日营业额</div></div>
        <div class="dg-card"><div class="dg-val">¥240</div><div class="dg-label">线下收入</div></div>
        <div class="dg-card"><div class="dg-val">¥4,260</div><div class="dg-label">已提现</div></div>
        <div class="dg-card"><div class="dg-val" style="color:#fa8c16">¥0.00</div><div class="dg-label">待入账</div></div>
      </div>

      <div class="section-title">提现历史</div>
      <div class="record-list">
        <div class="record-row">
          <div class="rr-icon check">↑</div>
          <div class="rr-info"><div class="rr-title">提现至招商银行 6214****6908</div><div class="rr-sub">2026-05-12 09:30 · 已打款</div></div>
          <div class="rr-amount minus">¥2,400.00</div>
        </div>
        <div class="record-row">
          <div class="rr-icon check">↑</div>
          <div class="rr-info"><div class="rr-title">提现至招商银行 6214****6908</div><div class="rr-sub">2026-04-30 16:00 · 已打款</div></div>
          <div class="rr-amount minus">¥1,860.00</div>
        </div>
        <div class="record-row">
          <div class="rr-icon" style="background:#fff7e6;color:#fa8c16">🕐</div>
          <div class="rr-info"><div class="rr-title">提现申请审核中</div><div class="rr-sub">2026-05-14 15:07 · 待园区审核</div></div>
          <div class="rr-amount" style="color:#fa8c16">¥1,200.00</div>
        </div>
      </div>
      <div style="height:12px"></div>

      <!-- Withdraw account popup -->
      <div v-if="showWithdrawAccount" class="popup-overlay" @click.self="showWithdrawAccount = false">
        <div class="popup-panel">
          <div class="popup-header">
            <span class="popup-title">提现账户</span>
            <button class="popup-close" @click="showWithdrawAccount = false">✕</button>
          </div>
          <div class="popup-body">
            <div class="wd-bank-info">
              <span style="font-size:24px">💳</span>
              <div><div style="font-weight:500">当前绑定账户</div><div style="font-size:12px;color:#666;margin-top:2px">招商银行 · 6214****6908 · 黄轩辉</div></div>
            </div>
            <div class="wd-field"><label>开户银行</label><select class="wd-select"><option>招商银行</option><option>工商银行</option><option>建设银行</option></select></div>
            <div class="wd-field"><label>银行账号</label><input class="wd-input" value="6214830169088888" /></div>
            <div class="wd-row">
              <div class="wd-field"><label>收款人姓名</label><input class="wd-input" value="黄轩辉" /></div>
              <div class="wd-field"><label>收款人手机</label><input class="wd-input" value="13912345678" /></div>
            </div>
            <button class="mp-btn" style="width:100%;margin-top:12px" @click="showToast('提现申请已提交，待园区审核'); showWithdrawAccount = false">申请提现</button>
          </div>
        </div>
      </div>
    </div>

    <!-- ======== me (我的) ======== -->
    <div v-else-if="currentTab === 'me'" class="tab-content">
      <div class="profile-hero">
        <div class="pf-row">
          <div class="pf-avatar">王</div>
          <div class="pf-info">
            <div class="pf-name">王五</div>
            <div class="pf-role">商户管理员 · 火车餐厅</div>
          </div>
        </div>
      </div>

      <div class="stat-row">
        <div class="stat-item"><div class="stat-val">22</div><div class="stat-label">今日订单</div></div>
        <div class="stat-item"><div class="stat-val">¥860</div><div class="stat-label">今日营业额</div></div>
        <div class="stat-item"><div class="stat-val">2</div><div class="stat-label">在岗员工</div></div>
      </div>

      <div class="menu-card">
        <div class="menu-row" @click="showToast('店铺信息')"><span class="mm-icon">🏪</span>店铺信息<span class="mm-extra">营业中</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('营业时间')"><span class="mm-icon">🕐</span>营业时间<span class="mm-extra">10:00 - 21:00</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showMenuMode = true"><span class="mm-icon">🍱</span>菜单/商品管理<span class="mm-extra">18 项</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('员工权限')"><span class="mm-icon">👥</span>员工权限<span class="mm-extra">2 名收银员</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showWithdrawAccount = true"><span class="mm-icon">💳</span>提现账户<span class="mm-extra">招商银行</span><span class="mm-arrow">›</span></div>
      </div>

      <div class="menu-card">
        <div class="menu-row" @click="showToast('消息通知')"><span class="mm-icon">✉</span>消息通知<span class="mm-extra">2 条未读</span><span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('账户安全')"><span class="mm-icon">🔒</span>账户安全<span class="mm-arrow">›</span></div>
        <div class="menu-row" @click="showToast('关于我们')"><span class="mm-icon">ℹ</span>关于我们<span class="mm-arrow">›</span></div>
        <div class="menu-row logout" @click="handleLogout"><span class="mm-icon">🚪</span>退出登录</div>
      </div>
      <div style="height:12px"></div>
    </div>

    <!-- Menu management popup -->
    <div v-if="showMenuMode" class="popup-overlay" @click.self="showMenuMode = false">
      <div class="popup-panel popup-full">
        <div class="popup-header">
          <span class="popup-title">🍱 菜单上架管理</span>
          <button class="popup-close" @click="showMenuMode = false">✕</button>
        </div>
        <div class="popup-body" style="max-height:400px;overflow-y:auto">
          <div class="tab-filter-row">
            <button v-for="(c, i) in menuCategories" :key="c" :class="['filter-btn', { active: menuCatActive === i }]" @click="menuCatActive = i">{{ c }}</button>
          </div>
          <div v-for="m in menuItems" :key="m.name" class="menu-item-card">
            <div class="mi-img" style="font-size:32px">{{ m.img }}</div>
            <div class="mi-info">
              <div class="mi-name-row">
                <span class="mi-name">{{ m.name }}</span>
                <span :class="['mp-tag', m.status === 'on' ? 'green' : 'gray']" style="font-size:10px">{{ m.status === 'on' ? '上架' : '下架' }}</span>
              </div>
              <div class="mi-desc">{{ m.desc }}</div>
              <span class="mi-price">¥{{ m.price.toFixed(2) }}</span>
            </div>
            <button :class="['mp-btn', 'sm', m.status === 'on' ? 'outline' : 'gray']" @click="toggleMenuItem(m)">{{ m.status === 'on' ? '下架' : '上架' }}</button>
          </div>
        </div>
        <div class="popup-footer">
          <button class="mp-btn" style="width:100%" @click="showToast('新增菜品表单已打开')">+ 新增菜品</button>
        </div>
      </div>
    </div>

    <!-- Business data popup -->
    <div v-if="showBusinessMode" class="popup-overlay" @click.self="showBusinessMode = false">
      <div class="popup-panel popup-full">
        <div class="popup-header">
          <span class="popup-title">📊 经营数据清单</span>
          <button class="popup-close" @click="showBusinessMode = false">✕</button>
        </div>
        <div class="popup-body" style="max-height:400px;overflow-y:auto">
          <div class="mp-data-hero">
            <div class="dh-title">今日经营概览</div>
            <div class="dh-row">
              <div><div class="dh-val">¥860</div><div class="dh-label">营业额</div></div>
              <div><div class="dh-val">22</div><div class="dh-label">订单数</div></div>
              <div><div class="dh-val">68</div><div class="dh-label">客流量</div></div>
            </div>
          </div>
          <div class="mp-data-grid">
            <div class="dg-card"><div class="dg-val">¥620</div><div class="dg-label">线上 (美团+抖音)</div></div>
            <div class="dg-card"><div class="dg-val">¥240</div><div class="dg-label">线下收入</div></div>
            <div class="dg-card"><div class="dg-val" style="color:#fa8c16">¥346</div><div class="dg-label">☷ 美团团购</div></div>
            <div class="dg-card"><div class="dg-val" style="color:#010101">¥207</div><div class="dg-label">♫ 抖音团购</div></div>
          </div>
          <div class="section-title">今日订单明细</div>
          <div class="record-list">
            <div v-for="r in businessOrders" :key="r.title" class="record-row">
              <div class="rr-icon" :style="{ background: r.iconBg, color: r.iconColor }">{{ r.icon }}</div>
              <div class="rr-info"><div class="rr-title">{{ r.title }}</div><div class="rr-sub">{{ r.sub }}</div></div>
              <div class="rr-amount" :style="{ color: r.color }">{{ r.amount }}</div>
            </div>
          </div>
        </div>
      </div>
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

// Orders
const orderFilter = ref('all')
const orderFilters = [
  { k: 'all', t: '全部', n: 22 },
  { k: 'pending', t: '待支付', n: 3 },
  { k: 'paid', t: '待核销', n: 5 },
  { k: 'done', t: '已完成', n: 13 },
  { k: 'refunded', t: '已退款', n: 1 },
]

const orderList = [
  { no: '202605141422001', user: '稻田守望者', avatar: '稻', items: '火车特色套餐 ×1, 稻田米浆 ×1', amount: 86.50, status: '待核销', time: '14:22', pay: '余额', code: 'HX26051414220' },
  { no: '202605141315002', user: '田园生活家', avatar: '田', items: '精品咖啡 ×2', amount: 32.00, status: '已核销', time: '13:15', pay: '微信', code: '' },
  { no: '202605141025004', user: '小麦的麦', avatar: '麦', items: '火车特色套餐 ×2', amount: 128.00, status: '已核销', time: '10:25', pay: '余额', code: '' },
  { no: '202605140950006', user: '稻香一缕', avatar: '稻', items: '红烧肉饭 ×1', amount: 46.00, status: '退款处理中', time: '09:50', pay: '余额', code: '' },
  { no: '202605140915007', user: '小麦的麦', avatar: '麦', items: '稻田米浆 ×3', amount: 36.00, status: '待支付', time: '09:15', pay: '-', code: '' },
]

function statusClass(s: string) {
  const m: Record<string, string> = { '已核销': 'green', '待核销': 'red', '待支付': 'orange', '退款处理中': 'gray' }
  return m[s] || 'gray'
}

// Verify
const verifyMode = ref<'verify' | 'collect' | 'meituan' | 'douyin'>('verify')
const verifyCode = ref('')
const collectAmount = ref('')

const verifyModeDesc = computed(() => {
  const m: Record<string, string> = {
    'verify': '扫描用户出示的订单核销码完成核销',
    'collect': '扫描用户付款码 + 手动输入金额完成收款',
    'meituan': '扫描美团团购券二维码 · 验证券码有效性并完成核销',
    'douyin': '扫描抖音团购券二维码 · 宣券验证券码 · 支持套餐/次卡/代金券',
  }
  return m[verifyMode.value]
})

const trendData = [280, 320, 290, 380, 420, 380, 560, 480, 520, 420, 580, 620, 540, 860]

const verifiedToday = [
  { user: '田园生活家', item: '精品咖啡 ×2', amount: '32.00', time: '13:15', code: 'HX26051413150' },
  { user: '小麦的麦', item: '火车特色套餐 ×2', amount: '128.00', time: '10:25', code: 'HX26051410250' },
  { user: '稻香一缕', item: '美式咖啡', amount: '46.00', time: '09:15', code: 'HX26051409150' },
]

function doVerify() {
  if (!verifyCode.value) { showToast('请输入核销码'); return }
  showToast('核销成功 ✓ ' + verifyCode.value)
  verifyCode.value = ''
}

function confirmCollect() {
  const amt = parseFloat(collectAmount.value)
  if (!amt || amt <= 0) { showToast('请输入收款金额'); return }
  showToast('收款成功 ¥' + amt.toFixed(2) + ' · 从余额扣款')
  collectAmount.value = ''
  router.push('/mp-merchant/orders')
}

// Meituan flow
const meituanStep = ref<'scan' | 'confirm' | 'done'>('scan')
const meituanCode = ref('')
const meituanData = ref({
  code: 'MT-YF-8821-5023', product: '稻田双人套餐', originPrice: 198.00, salePrice: 128.00,
  user: '稻田守望者', phone: '187****9908', validFrom: '2026-05-01', validTo: '2026-06-30', status: '可核销',
})

const meituanRecent = [
  { code: 'MT-YF-8821-5023', product: '稻田双人套餐', price: '128.00', time: '14:22', status: '已核销', tagClass: 'green' },
  { code: 'MT-YF-9942-7163', product: '树下咖啡体验券', price: '39.90', time: '13:15', status: '已核销', tagClass: 'green' },
  { code: 'MT-YF-7710-5532', product: '火车特色套餐', price: '99.00', time: '10:25', status: '已撤销', tagClass: 'gray' },
]

function meituanPreVerify() {
  meituanData.value = {
    code: meituanCode.value || 'MT-YF-8821-5023', product: '稻田双人套餐', originPrice: 198.00, salePrice: 128.00,
    user: '稻田守望者', phone: '187****9908', validFrom: '2026-05-01', validTo: '2026-06-30', status: '可核销',
  }
  meituanStep.value = 'confirm'
}
function meituanExecute() { meituanStep.value = 'done'; showToast('美团验券成功 ✓ 券码 ' + meituanData.value.code) }
function meituanBack() { meituanStep.value = 'scan' }
function meituanReset() { meituanStep.value = 'scan'; meituanCode.value = '' }

// Douyin flow
const douyinStep = ref<'scan' | 'confirm' | 'done'>('scan')
const douyinCode = ref('')
const douyinData = ref({
  code: 'DY-GROUP-A8K2-M9P1', product: '稻田丰收双人餐', couponType: '团购套餐',
  originPrice: 228.00, salePrice: 158.00, user: '稻田守望者', phone: '187****9908',
  validFrom: '2026-05-01', validTo: '2026-06-30', maxCount: 1, verifyCount: 1, status: '可核销',
})

const douyinRecent = [
  { code: 'DY-GROUP-A8K2-M9P1', product: '稻田丰收双人餐', price: '158.00', time: '18:20', status: '已核销', tagClass: 'green' },
  { code: 'DY-VOUCHER-B1N6-S3D8', product: '满100减20代金券', price: '9.90', time: '14:20', status: '已核销', tagClass: 'green' },
  { code: 'DY-CARD-T9W4-L8H6', product: '稻田米浆次卡 (4/10)', price: '199.00', time: '15:00', status: '部分核销', tagClass: 'blue' },
]

function douyinPreVerify() {
  douyinData.value = {
    code: douyinCode.value || 'DY-GROUP-A8K2-M9P1', product: '稻田丰收双人餐', couponType: '团购套餐',
    originPrice: 228.00, salePrice: 158.00, user: '稻田守望者', phone: '187****9908',
    validFrom: '2026-05-01', validTo: '2026-06-30', maxCount: 1, verifyCount: 1, status: '可核销',
  }
  douyinStep.value = 'confirm'
}
function douyinExecute() { douyinStep.value = 'done'; showToast('抖音验券成功 ♫ 券码 ' + douyinData.value.code) }
function douyinBack() { douyinStep.value = 'scan' }
function douyinReset() { douyinStep.value = 'scan'; douyinCode.value = '' }

function doneOrderNo(prefix: string) {
  return prefix + '20260514' + Date.now().toString().slice(-6)
}

// Menu management
const showMenuMode = ref(false)
const menuCatActive = ref(0)
const menuCategories = ['招牌推荐', '主食', '小吃', '饮品', '甜品', '套餐']
const menuItemsRaw = [
  { name: '火车特色套餐', cat: '招牌推荐', price: 128.00, status: 'on' as const, img: '🍱', desc: '含红烧肉饭+稻田米浆+时蔬' },
  { name: '稻田双人餐', cat: '套餐', price: 198.00, status: 'on' as const, img: '🍱', desc: '双人份·火车特色菜品' },
  { name: '红烧肉饭', cat: '主食', price: 38.00, status: 'on' as const, img: '🍚', desc: '招牌红烧肉·配米饭' },
  { name: '稻田米浆', cat: '饮品', price: 18.00, status: 'on' as const, img: '🥤', desc: '现磨鲜制·冷/热可选' },
  { name: '手工豆腐花', cat: '甜品', price: 12.00, status: 'on' as const, img: '🍮', desc: '传统手工·甜/咸可选' },
  { name: '精品咖啡', cat: '饮品', price: 28.00, status: 'off' as const, img: '☕', desc: '手冲单品·树下特供' },
  { name: '时令蔬菜沙拉', cat: '小吃', price: 22.00, status: 'on' as const, img: '🥗', desc: '当季时蔬·现拌' },
  { name: '糯米糍粑', cat: '小吃', price: 15.00, status: 'on' as const, img: '🍡', desc: '手工糯米·红糖浇汁' },
]
const menuItems = ref([...menuItemsRaw])
function toggleMenuItem(m: typeof menuItemsRaw[0]) {
  menuItems.value = menuItems.value.map(it => it.name === m.name ? { ...it, status: it.status === 'on' ? 'off' : 'on' } as typeof it : it)
  showToast((m.status === 'on' ? '已下架' : '已上架') + ': ' + m.name)
}

// Business data
const showBusinessMode = ref(false)
const showWithdrawAccount = ref(false)
const businessOrders = [
  { icon: '↓', iconBg: '#e6f4ff', iconColor: '#1677ff', title: '稻田守望者 · 火车特色套餐', sub: '14:22 · 余额 · 已支付', amount: '-¥86.50', color: '#333' },
  { icon: '☷', iconBg: '#fff7e6', iconColor: '#fa8c16', title: '稻田守望者 · 美团·稻田双人套餐', sub: '14:22 · 美团券 MT-YF-8821', amount: '-¥128.00', color: '#333' },
  { icon: '✓', iconBg: '#f0f9eb', iconColor: '#67c23a', title: '田园生活家 · 精品咖啡 ×2', sub: '13:15 · 微信支付 · 已核销', amount: '-¥32.00', color: '#333' },
  { icon: '♫', iconBg: '#f5f5f5', iconColor: '#010101', title: '小麦的麦 · 抖音·丰收双人餐', sub: '18:20 · 抖音券 DY-GROUP-A8K2', amount: '-¥158.00', color: '#333' },
  { icon: '↓', iconBg: '#e6f4ff', iconColor: '#1677ff', title: '稻香一缕 · 红烧肉饭', sub: '09:50 · 余额 · 退款处理中', amount: '¥46.00', color: '#fa8c16' },
]

function handleLogout() {
  showToast('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.tab-content { background: #f5f5f5; min-height: 100%; }

/* Shared: Data Hero */
.mp-data-hero { background: linear-gradient(135deg, #1677ff, #0958d9); color: #fff; padding: 20px 16px 24px; }
.dh-title { font-size: 12px; opacity: 0.85; }
.dh-subtitle { font-size: 18px; font-weight: 600; margin-top: 4px; }
.dh-row { display: flex; margin-top: 14px; }
.dh-row > div { flex: 1; text-align: center; border-right: 1px solid rgba(255,255,255,0.2); }
.dh-row > div:last-child { border-right: none; }
.dh-val { font-size: 24px; font-weight: 700; letter-spacing: -0.5px; }
.dh-label { font-size: 11px; opacity: 0.8; margin-top: 4px; }

/* Shared: Data Grid */
.mp-data-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; margin: -12px 12px 12px; position: relative; z-index: 1; }
.dg-card { background: #fff; border-radius: 10px; padding: 14px; text-align: center; box-shadow: 0 1px 4px rgba(0,0,0,0.04); }
.dg-val { font-size: 20px; font-weight: 700; }
.dg-label { font-size: 11px; color: #999; margin-top: 4px; }

/* Shared: Section */
.section-title { font-size: 14px; font-weight: 600; padding: 14px 16px 10px; color: #333; }

/* Shared: Quick Actions */
.quick-actions { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; padding: 0 12px; }
.qa-item { background: #fff; border-radius: 10px; padding: 14px 8px; text-align: center; font-size: 11px; color: #333; cursor: pointer; }
.qa-icon { font-size: 24px; margin-bottom: 6px; }

/* Shared: Record List */
.record-list { margin: 0 12px; }
.record-row { display: flex; align-items: center; gap: 10px; padding: 12px 14px; background: #fff; border-bottom: 1px solid #f0f0f0; }
.record-row:first-child { border-radius: 10px 10px 0 0; }
.record-row:nth-last-child(2) { border-radius: 0 0 10px 10px; border-bottom: none; }
.rr-icon { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; background: #f0f9eb; color: #67c23a; }
.rr-icon.check { background: #f0f9eb; color: #67c23a; }
.rr-info { flex: 1; min-width: 0; }
.rr-title { font-size: 13px; color: #333; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rr-sub { font-size: 11px; color: #999; margin-top: 2px; }
.rr-amount { font-size: 15px; font-weight: 600; flex-shrink: 0; }
.rr-amount.minus { color: #333; }

/* Tags */
.mp-tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 10px; font-weight: 500; }
.mp-tag.green { background: #f0f9eb; color: #67c23a; }
.mp-tag.gray { background: #f5f5f5; color: #999; }
.mp-tag.red { background: #fff2f0; color: #ff4d4f; }
.mp-tag.orange { background: #fff7e6; color: #fa8c16; }
.mp-tag.blue { background: #e6f4ff; color: #1677ff; }

/* Buttons */
.mp-btn { background: #1677ff; color: #fff; border: none; border-radius: 8px; padding: 8px 20px; font-size: 13px; cursor: pointer; font-weight: 500; }
.mp-btn.sm { padding: 4px 12px; font-size: 12px; }
.mp-btn.outline { background: #fff; color: #1677ff; border: 1px solid #1677ff; }
.mp-btn.gray { background: #f5f5f5; color: #666; }

/* Order Card */
.tab-filter-row { display: flex; gap: 4px; padding: 6px; margin: 10px 12px 0; background: #f0f0f0; border-radius: 8px; }
.filter-btn { flex: 1; padding: 6px; border: none; border-radius: 6px; font-size: 12px; background: transparent; color: #666; cursor: pointer; }
.filter-btn.active { background: #fff; color: #1677ff; font-weight: 500; box-shadow: 0 1px 3px rgba(0,0,0,0.06); }

.order-card { margin: 12px; background: #fff; border-radius: 10px; padding: 14px; }
.oc-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.oc-avatar { width: 36px; height: 36px; border-radius: 50%; background: #f5e8d8; color: #8b6443; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 600; }
.oc-user-info { flex: 1; }
.oc-user { font-size: 14px; font-weight: 600; }
.oc-time { font-size: 11px; color: #999; }
.oc-items { font-size: 13px; color: #666; margin-bottom: 10px; padding-left: 46px; }
.oc-footer { display: flex; align-items: center; justify-content: space-between; border-top: 1px dashed #f0f0f0; padding-top: 10px; }
.oc-no { font-size: 11px; color: #999; font-family: monospace; }
.oc-actions { display: flex; align-items: center; gap: 10px; }
.oc-amount { font-size: 18px; font-weight: 600; color: #ff4d4f; }

/* Overview */
.overview-card { background: #fff; margin: -12px 12px 12px; border-radius: 10px; padding: 14px 16px; position: relative; z-index: 1; }
.ov-title { font-size: 14px; font-weight: 600; margin-bottom: 10px; }
.ov-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 6px; text-align: center; }
.ov-val { font-size: 20px; font-weight: 600; }
.ov-label { font-size: 11px; color: #999; margin-top: 2px; }

/* Trend */
.trend-card { margin: 0 12px; background: #fff; border-radius: 10px; padding: 14px; }
.trend-bars { display: flex; align-items: flex-end; gap: 4px; height: 120px; }
.trend-bar { flex: 1; background: linear-gradient(180deg, #1677ff, #0958d9); border-radius: 3px 3px 0 0; min-width: 0; }
.trend-labels { display: flex; justify-content: space-between; font-size: 11px; color: #999; margin-top: 6px; }

/* Scan */
.scan-mode-bar { background: #fff; padding: 12px 16px; border-bottom: 1px solid #f0f0f0; }
.scan-mode-tabs { display: flex; gap: 4px; padding: 3px; background: #f0f0f0; border-radius: 8px; }
.sm-btn { flex: 1; min-width: 72px; padding: 7px 6px; border: none; border-radius: 6px; font-size: 12px; cursor: pointer; background: transparent; color: #666; }
.sm-btn.active { background: #fff; color: #ff4d4f; font-weight: 500; box-shadow: 0 1px 4px rgba(0,0,0,0.06); }
.scan-mode-desc { font-size: 11px; color: #bbb; margin-top: 6px; text-align: center; }
.scan-viewport { padding: 20px 16px; text-align: center; color: #fff; min-height: 280px; }
.scan-viewport.default-bg { background: linear-gradient(180deg, #1677ff, #0958d9); }
.scan-viewport.meituan-bg { background: linear-gradient(180deg, #fa8c16, #d46b08); }
.scan-viewport.douyin-bg { background: linear-gradient(180deg, #010101, #1a1a1a); }
.scan-frame-wrap { width: 180px; height: 180px; border: 3px solid rgba(255,255,255,0.5); border-radius: 16px; display: flex; align-items: center; justify-content: center; background: rgba(255,255,255,0.05); margin: 0 auto 16px; position: relative; }
.scan-frame-icon { font-size: 64px; opacity: 0.8; }
.scan-frame-label { font-size: 15px; font-weight: 500; }
.scan-frame-sub { font-size: 12px; opacity: 0.7; margin-top: 4px; }
.scan-manual { padding: 14px 12px; }
.manual-title { font-size: 14px; font-weight: 500; margin-bottom: 10px; background: #fff; padding: 14px; border-radius: 10px 10px 0 0; }
.manual-row { display: flex; gap: 8px; background: #fff; padding: 0 14px 14px; border-radius: 0 0 10px 10px; }
.manual-input { flex: 1; height: 38px; padding: 8px 12px; border: 1px solid #e0e0e0; border-radius: 8px; font-size: 14px; font-family: monospace; outline: none; }
.manual-hint { font-size: 11px; color: #bbb; margin-top: 6px; padding: 0 14px; }
.scan-animated .scan-line { position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, transparent, #fff, transparent); animation: scanline 2s linear infinite; }
@keyframes scanline { 0%,100%{transform:translateY(0)} 50%{transform:translateY(168px)} }
.sf-corners { position: relative; }
.sf-corner { position: absolute; width: 28px; height: 28px; }
.sf-corner.tl { top: -3px; left: -3px; border-top: 3px solid #fff; border-left: 3px solid #fff; border-radius: 6px 0 0 0; }
.sf-corner.tr { top: -3px; right: -3px; border-top: 3px solid #fff; border-right: 3px solid #fff; border-radius: 0 6px 0 0; }
.sf-corner.bl { bottom: -3px; left: -3px; border-bottom: 3px solid #fff; border-left: 3px solid #fff; border-radius: 0 0 0 6px; }
.sf-corner.br { bottom: -3px; right: -3px; border-bottom: 3px solid #fff; border-right: 3px solid #fff; border-radius: 0 0 6px 0; }
.sf-corner.dy.tl { border-color: #ff4d4f; }
.sf-corner.dy.tr { border-color: #00f2fe; }
.sf-corner.dy.bl { border-color: #00f2fe; }
.sf-corner.dy.br { border-color: #ff4d4f; }
.scan-line-meituan { position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, transparent, #fff, transparent); animation: scanlinemt 2s linear infinite; }
.scan-line-douyin { position: absolute; top: 0; left: 0; right: 0; height: 2px; background: linear-gradient(90deg, transparent, #ff4d4f, transparent); animation: scanlinemt 2s linear infinite; }
@keyframes scanlinemt { 0%,100%{transform:translateY(0)} 50%{transform:translateY(188px)} }

/* Meituan/Douyin recent */
.recent-card { margin: 8px 12px; background: #fff; border-radius: 10px; padding: 14px; }
.recent-title { font-size: 14px; font-weight: 500; margin-bottom: 10px; }
.recent-row { display: flex; align-items: center; padding: 10px 0; border-bottom: 1px solid #f5f5f5; }
.recent-row:last-child { border-bottom: none; }
.recent-info { flex: 1; }
.recent-name { font-size: 13px; font-weight: 500; }
.recent-sub { font-size: 11px; color: #999; margin-top: 2px; }
.recent-price { font-size: 16px; font-weight: 600; color: #fa8c16; margin: 0 8px; }

/* Info notice */
.info-notice { margin: 0 12px; padding: 12px 14px; background: #fff7e6; border-radius: 8px; border: 1px solid #ffd591; font-size: 12px; color: #d48806; line-height: 1.7; }
.info-notice b { color: #fa8c16; }
.info-notice.dy-notice { background: #f5f5f5; border-color: #e5e5e5; color: #666; }
.info-notice.dy-notice b { color: #010101; }

/* Verify confirm */
.verify-confirm-header { display: flex; align-items: center; gap: 8px; padding: 16px; background: linear-gradient(180deg, #fff7e6, #fff 100px); }
.verify-detail-card { margin: 12px; background: #fff; border-radius: 10px; padding: 16px; }
.vdc-header { display: flex; align-items: center; gap: 10px; padding-bottom: 14px; border-bottom: 1px dashed #f0f0f0; }
.vdc-icon { width: 48px; height: 48px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 600; }
.vdc-icon.meituan { background: #fff7e6; color: #fa8c16; }
.vdc-icon.douyin { background: #f5f5f5; color: #010101; }
.vdc-info { flex: 1; }
.vdc-name { font-size: 16px; font-weight: 600; }
.vdc-sub { font-size: 12px; color: #999; margin-top: 2px; }
.vdc-details { padding-top: 14px; }
.vdc-row { display: flex; justify-content: space-between; padding: 8px 0; font-size: 13px; }
.vdc-row span:first-child { color: #999; }
.mono { font-family: monospace; font-weight: 500; }
.strike { text-decoration: line-through; color: #999; }
.big-price { font-size: 18px; font-weight: 600; }
.fw600 { font-weight: 600; }
.vdc-section-title { font-size: 14px; font-weight: 500; margin-bottom: 10px; }

.verify-count-row { padding: 10px 0; border-bottom: 1px solid #f5f5f5; }
.vcr-label { font-size: 13px; color: #999; margin-bottom: 8px; }
.vcr-controls { display: flex; align-items: center; gap: 12px; }
.vcr-btn { width: 36px; height: 36px; border: 1px solid #ddd; border-radius: 50%; background: #fff; font-size: 18px; display: flex; align-items: center; justify-content: center; cursor: pointer; }
.vcr-val { font-size: 22px; font-weight: 600; min-width: 40px; text-align: center; }
.vcr-remain { font-size: 12px; color: #999; }

/* Done view */
.done-view { padding: 40px 20px 20px; text-align: center; }
.done-icon { width: 80px; height: 80px; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 42px; }
.done-icon.meituan-done { background: linear-gradient(135deg, #fa8c16, #d46b08); box-shadow: 0 8px 20px rgba(250,140,22,0.3); }
.done-icon.douyin-done { background: linear-gradient(135deg, #010101, #333); box-shadow: 0 8px 20px rgba(0,0,0,0.2); }
.done-title { font-size: 20px; font-weight: 600; margin-bottom: 8px; }
.done-sub { font-size: 13px; color: #999; margin-bottom: 24px; }
.done-card { background: #fff; border-radius: 12px; padding: 20px; text-align: left; }
.dok-row { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px dashed #ebebeb; }
.dok-row:last-child { border-bottom: none; }
.dok-label { font-size: 13px; color: #999; width: 80px; flex-shrink: 0; }

/* Bottom bar */
.bottom-bar { background: #fff; padding: 12px 16px 20px; border-top: 1px solid #f0f0f0; display: flex; gap: 10px; }

/* Collect panel */
.collect-panel { padding: 14px 12px; }
.cp-user-row { display: flex; align-items: center; gap: 10px; background: #fff; border-radius: 10px; padding: 14px; margin-bottom: 10px; }
.cp-avatar { width: 40px; height: 40px; border-radius: 50%; background: #fff2f0; color: #ff4d4f; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 600; }
.cp-info { flex: 1; }
.cp-name { font-size: 14px; font-weight: 500; }
.cp-phone { font-size: 11px; color: #999; font-weight: 400; }
.cp-meta { font-size: 11px; color: #999; margin-top: 2px; }
.cp-amount-section { background: #fff; border-radius: 10px; padding: 14px; margin-bottom: 10px; }
.cp-label { font-size: 13px; color: #666; margin-bottom: 8px; }
.cp-amount-input { display: flex; align-items: center; background: #f5f5f5; border-radius: 8px; padding: 14px; gap: 8px; }
.cp-yuan { font-size: 30px; font-weight: 300; color: #333; line-height: 1; }
.cp-input { flex: 1; border: none; background: transparent; outline: none; font-size: 30px; font-weight: 500; color: #333; letter-spacing: -1px; }
.cp-presets { display: flex; gap: 6px; margin-top: 10px; flex-wrap: wrap; }
.cp-preset-btn { padding: 6px 14px; border: 1px solid #ddd; background: #fff; border-radius: 14px; font-size: 13px; color: #333; cursor: pointer; }
.cp-preset-btn.clear { color: #999; }
.cp-remark { background: #fff; border-radius: 10px; padding: 14px; margin-bottom: 10px; }
.cp-remark-input { width: 100%; height: 36px; padding: 6px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 13px; outline: none; box-sizing: border-box; }
.cp-payment { background: #fff; border-radius: 10px; padding: 14px; }
.cp-radio { display: flex; align-items: center; gap: 10px; padding: 10px; border: 1px solid #ddd; border-radius: 8px; margin-bottom: 6px; cursor: pointer; }
.cp-radio.checked { border-color: #ff4d4f; background: #fff2f0; }
.cp-radio-dot { width: 16px; height: 16px; border-radius: 50%; border: 2px solid #ff4d4f; display: flex; align-items: center; justify-content: center; }
.cp-radio-dot::after { content: ''; width: 8px; height: 8px; background: #ff4d4f; border-radius: 50%; }
.cp-radio-ring { width: 16px; height: 16px; border-radius: 50%; border: 2px solid #d9d9d9; }
.cp-radio-text { flex: 1; font-size: 14px; }
.cp-radio-extra { font-size: 13px; color: #999; }

/* Finance */
.finance-hero { background: linear-gradient(135deg, #1677ff, #0958d9); padding: 20px 16px 22px; color: #fff; }
.fh-label { font-size: 13px; opacity: 0.85; }
.fh-balance { font-size: 36px; font-weight: 600; margin-top: 6px; letter-spacing: -1px; }
.fh-sub-row { display: flex; gap: 24px; font-size: 11px; opacity: 0.85; margin-top: 8px; }
.fh-actions { display: flex; gap: 10px; margin-top: 18px; }
.fh-btn { flex: 1; padding: 10px; border-radius: 8px; border: 1px solid rgba(255,255,255,0.4); background: rgba(255,255,255,0.18); color: #fff; font-size: 13px; cursor: pointer; }
.fh-btn.primary { background: #fff; color: #1677ff; border: none; }

/* Profile */
.profile-hero { background: linear-gradient(135deg, #1677ff, #0958d9); padding: 24px 16px 20px; color: #fff; }
.pf-row { display: flex; align-items: center; gap: 12px; }
.pf-avatar { width: 48px; height: 48px; border-radius: 50%; background: rgba(255,255,255,0.95); color: #1677ff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 600; }
.pf-info { flex: 1; }
.pf-name { font-size: 18px; font-weight: 600; }
.pf-role { color: rgba(255,255,255,0.85); font-size: 13px; margin-top: 4px; }

/* Stats Row */
.stat-row { display: grid; grid-template-columns: repeat(3, 1fr); background: #fff; border-radius: 12px; margin: 12px; padding: 14px 0; }
.stat-item { text-align: center; border-right: 1px solid #f0f0f0; }
.stat-item:last-child { border-right: none; }
.stat-val { font-size: 18px; font-weight: 600; }
.stat-label { font-size: 12px; color: #999; margin-top: 2px; }

/* Menu */
.menu-card { background: #fff; border-radius: 12px; margin: 12px; overflow: hidden; }
.menu-row { display: flex; align-items: center; gap: 10px; padding: 14px 16px; font-size: 14px; cursor: pointer; border-bottom: 1px solid #f5f5f5; }
.menu-row:last-child { border-bottom: none; }
.menu-row.logout { color: #ff4d4f; }
.mm-icon { font-size: 18px; width: 24px; text-align: center; flex-shrink: 0; }
.mm-extra { flex: 1; text-align: right; font-size: 12px; color: #999; }
.mm-arrow { color: #ccc; flex-shrink: 0; }

/* Popup */
.popup-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 10; }
.popup-panel { width: 340px; max-height: 80%; background: #fff; border-radius: 14px; overflow: hidden; }
.popup-panel.popup-full { width: 340px; }
.popup-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #f0f0f0; }
.popup-title { font-size: 15px; font-weight: 600; }
.popup-close { width: 28px; height: 28px; border-radius: 50%; background: #f5f5f5; border: none; font-size: 14px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.popup-body { padding: 12px; }
.popup-footer { padding: 12px 16px 20px; border-top: 1px solid #f0f0f0; background: #fff; }

/* Menu items in popup */
.menu-item-card { display: flex; align-items: center; gap: 12px; padding: 12px; background: #fff; border-radius: 10px; margin-bottom: 8px; border: 1px solid #f5f5f5; }
.mi-img { width: 56px; height: 56px; border-radius: 10px; background: #f5f0e8; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.mi-info { flex: 1; min-width: 0; }
.mi-name-row { display: flex; align-items: center; gap: 6px; margin-bottom: 4px; }
.mi-name { font-size: 14px; font-weight: 500; }
.mi-desc { font-size: 11px; color: #999; margin-bottom: 4px; }
.mi-price { font-size: 16px; font-weight: 600; color: #ff4d4f; }

/* Withdraw field */
.wd-bank-info { display: flex; align-items: center; gap: 8px; padding: 12px; background: #fff2f0; border-radius: 8px; margin-bottom: 12px; }
.wd-field { margin-bottom: 12px; }
.wd-field label { display: block; font-size: 13px; color: #666; margin-bottom: 6px; }
.wd-input, .wd-select { width: 100%; height: 40px; padding: 6px 12px; border: 1px solid #ddd; border-radius: 8px; font-size: 14px; outline: none; box-sizing: border-box; }
.wd-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
</style>
