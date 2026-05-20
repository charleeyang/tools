<template>
  <div class="mp-client-view">
    <!-- ======== home (首页) ======== -->
    <div v-if="currentTab === 'home'" class="tab-content">
      <!-- Customer service banner -->
      <div class="cs-banner">
        <div class="cs-avatars">
          <div class="cs-avatar">🐰</div>
          <div class="cs-avatar">🥚</div>
        </div>
        <div class="cs-body">
          <div class="cs-line1">点击添加企微客服 <span class="cs-tag">@袁夫稻田</span></div>
          <div class="cs-line2">在线时间 9:00 - 18:00</div>
        </div>
        <div class="cs-arrow">›</div>
      </div>

      <!-- Category grid -->
      <div class="cat-grid">
        <div v-for="cat in categories" :key="cat.icon" class="cat-item" @click="showToast(cat.name)">
          <div class="cat-icon" v-html="cat.iconSvg"></div>
          <div class="cat-label">{{ cat.name }}</div>
        </div>
      </div>

      <!-- New products -->
      <div class="section-header">
        <div class="sh-title"><span class="sh-bar"></span>新品专区</div>
        <div class="sh-more"><span class="cs-avatar" style="width:24px;height:18px;font-size:14px;border-radius:9px">🐰</span> ›</div>
      </div>
      <div class="product-row">
        <div v-for="p in newProducts" :key="p.label" class="product-card">
          <div class="pc-img" :style="{ background: p.bg }">
            <div class="pc-brand">袁夫稻田</div>
            <div v-html="p.illustration"></div>
            <div v-if="p.presale" class="pc-presale">预售</div>
          </div>
          <div class="pc-meta">
            <span class="pc-price"><span style="font-size:12px">¥</span>{{ p.price }}</span>
            <span class="pc-cart-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#333" stroke-width="1.5"><circle cx="9" cy="20" r="1.5"/><circle cx="18" cy="20" r="1.5"/><path d="M3 4 L6 4 L8 16 L20 16 L22 8 L7 8"/></svg>
            </span>
          </div>
        </div>
      </div>

      <!-- Yuanfu Rice Section -->
      <div class="section-header">
        <div class="sh-title"><span class="sh-bar"></span>袁夫米面</div>
        <div class="sh-more">更多 ›</div>
      </div>
      <div class="product-row">
        <div v-for="p in riceProducts" :key="p.label" class="product-card">
          <div class="pc-img" :style="{ background: p.bg }">
            <div class="pc-brand">袁夫稻田</div>
            <div v-html="p.illustration"></div>
          </div>
          <div class="pc-meta">
            <span class="pc-price"><span style="font-size:12px">¥</span>{{ p.price }}</span>
            <span class="pc-cart-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#333" stroke-width="1.5"><circle cx="9" cy="20" r="1.5"/><circle cx="18" cy="20" r="1.5"/><path d="M3 4 L6 4 L8 16 L20 16 L22 8 L7 8"/></svg>
            </span>
          </div>
        </div>
      </div>

      <!-- Yuanfu Condiments -->
      <div class="section-header">
        <div class="sh-title"><span class="sh-bar"></span>袁夫味品</div>
        <div class="sh-more">更多 ›</div>
      </div>
      <div class="product-row" style="padding-bottom:16px">
        <div v-for="p in condimentProducts" :key="p.label" class="product-card">
          <div class="pc-img" :style="{ background: p.bg }">
            <div class="pc-brand">袁夫稻田</div>
            <div v-html="p.illustration"></div>
          </div>
          <div class="pc-meta">
            <span class="pc-price"><span style="font-size:12px">¥</span>{{ p.price }}</span>
            <span class="pc-cart-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#333" stroke-width="1.5"><circle cx="9" cy="20" r="1.5"/><circle cx="18" cy="20" r="1.5"/><path d="M3 4 L6 4 L8 16 L20 16 L22 8 L7 8"/></svg>
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- ======== memberqr (亮码支付) ======== -->
    <div v-else-if="currentTab === 'memberqr'" class="tab-content memberqr-tab">
      <div class="mqr-panel">
        <div class="mqr-head">
          <div class="mqr-avatar">稻</div>
          <div class="mqr-user">
            <div class="mqr-name">稻田守望者</div>
            <div class="mqr-member">会员</div>
          </div>
          <div class="mqr-balance-col">
            <div class="mqr-balance">¥{{ balance.toFixed(2) }}</div>
            <div class="mqr-balance-label">可用余额</div>
          </div>
        </div>

        <div class="mqr-card">
          <div class="mqr-tabs">
            <button :class="['mqr-tab', { active: qrMode === 'pay' }]" @click="qrMode = 'pay'">付款码</button>
            <button :class="['mqr-tab', { active: qrMode === 'member' }]" @click="qrMode = 'member'">会员码</button>
            <button :class="['mqr-tab', { active: qrMode === 'points' }]" @click="qrMode = 'points'">积分核销</button>
          </div>
          <div class="mqr-codebox">
            <div class="mqr-qrcode" style="display:flex;align-items:center;justify-content:center">
              <svg viewBox="0 0 100 100" width="200" height="200">
                <rect width="100" height="100" fill="#fff"/>
                <g v-html="qrPattern"></g>
                <rect x="0" y="0" width="35" height="35" fill="#fff"/>
                <rect x="0" y="0" width="35" height="35" fill="none" stroke="#000" stroke-width="4"/>
                <rect x="10" y="10" width="15" height="15" fill="#000" rx="1"/>
                <rect x="65" y="0" width="35" height="35" fill="#fff"/>
                <rect x="65" y="0" width="35" height="35" fill="none" stroke="#000" stroke-width="4"/>
                <rect x="75" y="10" width="15" height="15" fill="#000" rx="1"/>
                <rect x="0" y="65" width="35" height="35" fill="#fff"/>
                <rect x="0" y="65" width="35" height="35" fill="none" stroke="#000" stroke-width="4"/>
                <rect x="10" y="75" width="15" height="15" fill="#000" rx="1"/>
                <rect x="42" y="42" width="16" height="16" fill="#fff"/>
                <rect x="44" y="44" width="12" height="12" rx="2" fill="#1677ff"/>
              </svg>
            </div>
          </div>
          <div class="mqr-tip">
            {{ qrMode === 'pay' ? '出示付款码，商家扫码即付' : qrMode === 'member' ? '出示会员码，享受会员优惠' : '出示积分码，使用积分抵扣' }}
          </div>

          <div class="mqr-actions">
            <button class="mqr-action-btn" @click="refreshQr">🔄 刷新</button>
            <button class="mqr-action-btn" @click="showToast('余额 ¥' + balance.toFixed(2) + ' · 本金 ¥80.00 · 赠送 ¥50.00')">💰 余额</button>
            <button class="mqr-action-btn" @click="subPage = 'recharge'">💳 充值</button>
          </div>
        </div>

        <div class="mqr-footer-tip">安全支付 · 袁夫稻田智慧园区</div>
      </div>
    </div>

    <!-- ======== cart (购物车) ======== -->
    <div v-else-if="currentTab === 'cart'" class="tab-content cart-tab">
      <div class="cart-header">
        <span>共 {{ cartItems.length }} 件商品</span>
        <span style="color:#ff4d4f" @click="showToast('编辑模式')">编辑</span>
      </div>

      <div v-for="(it, i) in cartItems" :key="i" class="cart-item">
        <div :class="['cart-check', { checked: it.sel }]" @click="it.sel = !it.sel">{{ it.sel ? '✓' : '' }}</div>
        <div class="cart-img">{{ it.img }}</div>
        <div class="cart-info">
          <div class="cart-name">{{ it.n }}</div>
          <div class="cart-shop">袁夫稻田 · 自营</div>
          <div class="cart-bottom">
            <span class="cart-price">¥{{ it.price.toFixed(2) }}</span>
            <div class="cart-qty">
              <button class="cq-btn" @click="it.qty = Math.max(1, it.qty - 1)">−</button>
              <span class="cq-val">{{ it.qty }}</span>
              <button class="cq-btn" @click="it.qty++">+</button>
            </div>
          </div>
        </div>
      </div>

      <div class="cart-spacer"></div>
      <div class="cart-bar">
        <div :class="['cart-check', { checked: cartAllSel }]" @click="toggleAllCart">{{ cartAllSel ? '✓' : '' }}</div>
        <span class="cart-bar-label">全选</span>
        <div class="cart-bar-right">
          <div class="cart-bar-total">合计：<span class="cart-total-price">¥{{ cartTotal.toFixed(2) }}</span></div>
          <button class="mp-btn" style="padding:10px 24px" @click="showToast('结算成功')">结算 ({{ cartSelCount }})</button>
        </div>
      </div>
    </div>

    <!-- ======== me (个人中心) ======== -->
    <div v-else-if="currentTab === 'me'" class="tab-content me-tab">
      <!-- Sub pages inside me -->
      <template v-if="subPage === 'recharge'">
        <div class="recharge-hero">
          <div class="rh-label">当前余额</div>
          <div class="rh-balance">¥ {{ balance.toFixed(2) }}</div>
          <div class="rh-sub">累计消费 ¥256.50</div>
        </div>
        <div class="recharge-card">
          <div class="rc-header">
            <span class="rc-title">充值金额</span>
            <span style="font-size:11px;color:#fa8c16">🎁 充值有礼，多充多送</span>
          </div>
          <div class="rc-grid">
            <div v-for="t in rechargeTiles" :key="t.amount" :class="['rc-tile', { active: rechargeAmount === t.amount }]" @click="rechargeAmount = t.amount">
              <div v-if="t.tag" class="rc-ribbon">{{ t.tag }}</div>
              <template v-if="t.amount === 0">
                <div class="rc-amount-label">自定义</div>
                <div class="rc-gift">输入金额</div>
              </template>
              <template v-else>
                <div class="rc-amount">¥{{ t.amount }}</div>
                <div class="rc-gift">{{ t.gift > 0 ? '送 ¥' + t.gift : '' }}</div>
              </template>
            </div>
          </div>
          <div v-if="rechargeAmount === 0" class="rc-custom">
            <span style="font-size:18px;font-weight:600">¥</span>
            <input v-model="rechargeCustom" class="rc-custom-input" placeholder="输入充值金额" />
          </div>
        </div>
        <div class="recharge-payment">
          <span>支付方式</span>
          <div class="rp-wx"><span style="color:#07c160;font-size:16px">💚</span><span style="font-weight:500">微信支付</span></div>
        </div>
        <div class="recharge-rules">
          <div class="rrules-title">充值说明</div>
          <div>· 充值金额到账后可在园区内消费抵扣</div>
          <div>· 消费时<b>优先扣除本金</b>，本金用完后再使用赠送金额</div>
          <div>· 赠送金额仅可消费使用，不支持提现</div>
          <div>· <b style="color:#ff4d4f">⚠ 退款规则：</b>申请退款时仅退还剩余本金，赠送金额将<b style="color:#ff4d4f">自动归零且不可恢复</b></div>
          <div>· 可通过<a @click="subPage = 'refund'" style="color:#ff4d4f">「余额退款」</a>原路退回剩余本金</div>
          <div>· 充值后可在「充值记录」中查询交易</div>
        </div>
        <div class="bottom-bar">
          <button class="mp-btn" style="width:100%;padding:14px" @click="doRecharge">立即充值 ¥{{ rechargeAmount === 0 ? (rechargeCustom || '0') : rechargeAmount }}</button>
        </div>
      </template>

      <template v-else-if="subPage === 'recharge-records'">
        <div class="section-title">充值记录</div>
        <div class="record-list">
          <div v-for="r in rechargeRecords" :key="r.id" class="record-row">
            <div class="rr-icon" :style="{ background: r.iconBg, color: r.iconColor }">{{ r.icon }}</div>
            <div class="rr-info"><div class="rr-title">{{ r.title }}</div><div class="rr-sub">{{ r.sub }}</div></div>
            <div class="rr-amount" :style="{ color: r.color }">{{ r.amount }}</div>
          </div>
        </div>
        <div style="height:12px"></div>
        <button class="mp-btn gray" style="margin:12px;width:calc(100% - 24px)" @click="subPage = ''">返回个人中心</button>
      </template>

      <template v-else-if="subPage === 'refund'">
        <template v-if="refundStep === 1">
          <div class="refund-banner">
            <div class="rb-title">⤺ 余额退款</div>
            <div class="rb-desc">将余额原路退回到充值时的微信账户。退款仅支持本金部分，赠送金额不可退。</div>
          </div>
          <div class="refund-card">
            <div class="rfc-row"><span>当前余额</span><span class="rfc-val-lg">¥{{ balance.toFixed(2) }}</span></div>
            <div class="rfc-row"><span>本金 (可退)</span><span class="rfc-val-red">¥80.00</span></div>
            <div class="rfc-row"><span>赠送 (不可退)</span><span class="rfc-val-gray">¥50.00</span></div>
          </div>
          <div class="refund-card">
            <div class="rfc-section-title">退款金额</div>
            <div class="rfc-amount-input">
              <span style="font-size:22px;font-weight:600">¥</span>
              <input v-model.number="refundAmount" class="rfc-amount-field" placeholder="0.00" />
              <button class="mp-btn sm gray" @click="refundAmount = 80">全部退款</button>
            </div>
            <div class="rfc-hint">最高可退 ¥80.00 · 根据充值订单原路退回</div>
          </div>
          <div class="refund-card">
            <div class="rfc-section-title">退款原因 <span style="color:#ff4d4f">*</span></div>
            <div v-for="r in refundReasons" :key="r" :class="['rfc-reason', { selected: refundReason === r }]" @click="refundReason = r">
              <span :class="['rfc-radio', { checked: refundReason === r }]"><span v-if="refundReason === r" class="rfc-dot"></span></span>
              <span>{{ r }}</span>
            </div>
            <textarea v-if="refundReason === '其他'" class="rfc-textarea" placeholder="请描述具体原因 (最多 200 字)"></textarea>
          </div>
          <div class="recharge-rules">
            <div class="rrules-title">退款说明</div>
            <div>· 退款将原路返回到充值时使用的微信账户</div>
            <div>· 多次充值的，按 <b>"后充先退"</b> 原则匹配充值单</div>
            <div>· 园区管理员将在 1 个工作日内审核</div>
            <div>· 审核通过后，款项 1-3 个工作日内到账</div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="subPage = ''">取消</button>
            <button class="mp-btn" style="flex:2" @click="refundNext">提交申请</button>
          </div>
        </template>
        <template v-else-if="refundStep === 2">
          <div class="refund-banner"><div class="rb-title">⚠️ 请确认退款信息</div><div class="rb-desc">提交后不可撤销，请确认无误</div></div>
          <div class="refund-card">
            <div style="text-align:center;padding:14px 0">
              <div style="font-size:13px;color:#999">退款金额</div>
              <div style="font-size:36px;font-weight:700;color:#ff4d4f;margin-top:6px">¥ {{ refundAmount.toFixed(2) }}</div>
            </div>
            <div class="rfc-divider"></div>
            <div class="rfc-row"><span>退款方式</span><span>原路退回微信</span></div>
            <div class="rfc-row"><span>微信账户</span><span>稻***者 (138****5678)</span></div>
            <div class="rfc-row"><span>原充值单</span><span class="mono">CZ2026051410...</span></div>
            <div class="rfc-row"><span>退款原因</span><span>{{ refundReason || '不再使用' }}</span></div>
            <div class="rfc-row"><span>预计到账</span><span>1-3 个工作日</span></div>
          </div>
          <div class="refund-agreement">点击"确认提交"即表示您同意将本次退款金额按原路退回至微信账户，并知悉<a style="color:#ff4d4f">《退款服务协议》</a></div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="refundStep = 1">返回修改</button>
            <button class="mp-btn" style="flex:2" @click="refundSubmit">确认提交</button>
          </div>
        </template>
        <template v-else-if="refundStep === 3">
          <div class="done-view">
            <div class="done-icon refund-done">✓</div>
            <div class="done-title">退款申请已提交</div>
            <div class="done-sub">园区管理员将在 1 个工作日内审核</div>
            <div class="done-card" style="margin-top:30px">
              <div class="dok-row"><span class="dok-label">退款单号</span><span class="mono">{{ 'RF20260514' + Date.now().toString().slice(-6) }}</span></div>
              <div class="dok-row"><span class="dok-label">退款金额</span><span class="big-price" style="color:#ff4d4f">¥{{ refundAmount.toFixed(2) }}</span></div>
              <div class="dok-row"><span class="dok-label">当前状态</span><span class="mp-tag orange">待审核</span></div>
            </div>
            <div class="refund-progress">
              <b>退款进度</b><br>
              · 第 1 步 (已完成)：提交申请<br>
              · 第 2 步 (进行中)：园区审核 (预计 1 工作日)<br>
              · 第 3 步：微信原路退回 (审核通过后 1-3 工作日)
            </div>
          </div>
          <div class="bottom-bar">
            <button class="mp-btn gray" style="flex:1" @click="subPage = ''">返回首页</button>
            <button class="mp-btn" style="flex:1" @click="subPage = 'refund-records'">查看记录</button>
          </div>
        </template>
      </template>

      <template v-else-if="subPage === 'refund-records'">
        <div class="section-title">退款记录</div>
        <div v-for="r in refundRecords" :key="r.id" class="order-card">
          <div class="oc-header">
            <span class="oc-user">余额退款</span>
            <span :class="['mp-tag', r.tagClass]">{{ r.status }}</span>
          </div>
          <div class="oc-footer" style="border:none;padding-top:0"><span class="oc-no">{{ r.id }}</span><span class="oc-amount" style="color:#ff4d4f">¥{{ r.amt.toFixed(2) }}</span></div>
          <div style="font-size:12px;color:#999">{{ r.sub }}</div>
          <div class="rr-actions"><button class="mp-btn sm outline">查看详情</button><button v-if="r.status === '待审核'" class="mp-btn sm gray">取消申请</button></div>
        </div>
        <div style="height:12px"></div>
        <button class="mp-btn gray" style="margin:12px;width:calc(100% - 24px)" @click="subPage = ''">返回个人中心</button>
      </template>

      <template v-else-if="subPage === 'orders'">
        <div class="tab-filter-row">
          <button v-for="f in clientOrderFilters" :key="f" :class="['filter-btn', { active: clientOrderFilter === f }]" @click="clientOrderFilter = f">{{ f }}</button>
        </div>
        <div v-for="o in clientOrders" :key="o.shop" class="order-card">
          <div class="oc-header"><span class="oc-user">🏪 {{ o.shop }}</span><span :class="['mp-tag', o.statusTag]">{{ o.status }}</span></div>
          <div v-for="it in o.items" :key="it.n" style="display:flex;gap:10px;padding:8px 0">
            <div class="client-oi-img">{{ it.img }}</div>
            <div style="flex:1"><div style="font-size:14px">{{ it.n }}</div><div style="font-size:13px;color:#ff4d4f;margin-top:6px">¥{{ it.price.toFixed(2) }}</div></div>
          </div>
          <div class="oc-footer">
            <span style="font-size:13px;color:#666">合计 <span style="color:#ff4d4f;font-size:16px;font-weight:600">¥{{ o.total.toFixed(2) }}</span></span>
            <div class="oc-actions">
              <button v-if="o.status === '待付款'" class="mp-btn sm">立即支付</button>
              <button v-if="o.status === '待发货'" class="mp-btn sm outline">申请退款</button>
              <button v-if="o.status === '待收货'" class="mp-btn sm">确认收货</button>
            </div>
          </div>
        </div>
        <div style="height:12px"></div>
        <button class="mp-btn gray" style="margin:12px;width:calc(100% - 24px)" @click="subPage = ''">返回个人中心</button>
      </template>

      <template v-else-if="subPage === 'address'">
        <div v-for="a in addresses" :key="a.addr" class="order-card address-card">
          <div class="addr-header">
            <span style="font-size:15px;font-weight:600">{{ a.n }}</span>
            <span style="font-size:13px;color:#666;margin:0 8px">{{ a.phone }}</span>
            <span :class="['mp-tag', a.def ? 'red' : 'gray']">{{ a.def ? '默认' : a.tag }}</span>
          </div>
          <div style="font-size:13px;color:#666;line-height:1.5">{{ a.addr }}</div>
          <div class="rr-actions"><button class="mp-btn sm outline">编辑</button><button class="mp-btn sm gray">删除</button></div>
        </div>
        <div class="bottom-bar"><button class="mp-btn" style="width:100%">+ 新建地址</button></div>
      </template>

      <!-- Main me tab -->
      <template v-else>
        <div class="profile-hero">
          <div class="pf-row">
            <div class="pf-avatar">👤</div>
            <div class="pf-info">
              <div class="pf-name">点击显示微信头像</div>
              <span class="pf-grow">成长值 0</span>
            </div>
            <button class="pf-action-btn" @click="showToast('登录')">登录</button>
          </div>
        </div>

        <div class="member-strip">
          <span style="font-size:14px">◆</span>
          <span>成为会员，领 100 积分</span>
          <a class="m-link" @click="showToast('会员注册')">立即注册 ›</a>
        </div>

        <div class="stat-row">
          <div class="stat-item" @click="showToast('余额 ¥' + balance.toFixed(2) + ' · 本金 ¥80.00 · 赠送 ¥50.00')">
            <div class="stat-val" style="font-size:22px;font-weight:700;color:#ff4d4f">{{ balance.toFixed(2) }}</div>
            <div class="stat-label" style="font-weight:500">余额</div>
          </div>
          <div class="stat-item" @click="router.push('/mp-client/memberqr')">
            <div class="stat-val" style="font-size:22px;font-weight:700;color:#fa8c16">128</div>
            <div class="stat-label" style="font-weight:500">积分</div>
          </div>
          <div class="stat-item" @click="subPage = 'recharge'">
            <div class="stat-val"><span style="color:#1677ff;font-size:20px">¥</span></div>
            <div class="stat-label">钱包</div>
          </div>
        </div>

        <div class="menu-card">
          <div class="menu-row" @click="subPage = 'recharge'">
            <span class="mm-icon">💰</span>账户充值
            <span class="mm-extra" style="color:#ff4d4f">充 100 送 10</span>
            <span class="mm-arrow">›</span>
          </div>
          <div class="menu-row" @click="subPage = 'recharge-records'">
            <span class="mm-icon">📊</span>充值记录<span class="mm-arrow">›</span>
          </div>
          <div class="menu-row" @click="subPage = 'refund'">
            <span class="mm-icon">↩</span>余额退款
            <span class="mm-extra">可退至原支付方式</span>
            <span class="mm-arrow">›</span>
          </div>
          <div class="menu-row" @click="subPage = 'refund-records'">
            <span class="mm-icon">📋</span>退款记录<span class="mm-arrow">›</span>
          </div>
          <div class="menu-row" @click="subPage = 'orders'">
            <span class="mm-icon">📦</span>我的订单<span class="mm-arrow">›</span>
          </div>
          <div class="menu-row" @click="subPage = 'address'">
            <span class="mm-icon">📍</span>地址管理<span class="mm-arrow">›</span>
          </div>
        </div>

        <div class="menu-card">
          <div class="menu-row" @click="showToast('消息通知')"><span class="mm-icon">✉</span>消息通知<span class="mm-arrow">›</span></div>
          <div class="menu-row" @click="showToast('设置')"><span class="mm-icon">⚙</span>设置<span class="mm-arrow">›</span></div>
          <div class="menu-row" @click="showToast('关于我们')"><span class="mm-icon">ℹ</span>关于我们<span class="mm-arrow">›</span></div>
          <div class="menu-row logout" @click="handleLogout"><span class="mm-icon">🚪</span>退出登录</div>
        </div>
        <div style="height:12px"></div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()

const currentTab = computed(() => (route.meta.tab as string) || 'home')

// State
const balance = ref(130.00)
const qrMode = ref<'pay' | 'member' | 'points'>('pay')
const subPage = ref('')
const rechargeAmount = ref(100)
const rechargeCustom = ref('')
const refundStep = ref(1)
const refundAmount = ref(50)
const refundReason = ref('')

// Reset sub-page when tab changes
watch(currentTab, () => {
  if (currentTab.value !== 'me') subPage.value = ''
})

// QR pattern
const qrPattern = computed(() => {
  let html = ''
  const seed = Date.now()
  for (let r = 0; r < 20; r++) {
    for (let c = 0; c < 20; c++) {
      if (r < 7 && c < 7) continue
      if (r < 7 && c > 12) continue
      if (r > 12 && c < 7) continue
      const s = (r * 31 + c * 17 + r * c + seed % 100) % 7
      if (s < 3) html += '<rect x="' + (c * 5) + '" y="' + (r * 5) + '" width="5" height="5" fill="#000"/>'
    }
  }
  return html
})

function refreshQr() {
  showToast('二维码已刷新')
}

// Category SVG icons
const categories = [
  { name: '袁夫米面', icon: 'rice', iconSvg: '<svg viewBox="0 0 48 48" fill="none" stroke="#333" stroke-width="1.5"><path d="M8 24 Q8 38 24 38 Q40 38 40 24"/><ellipse cx="24" cy="24" rx="16" ry="4"/><path d="M14 20 Q14 16 18 16 Q22 14 24 16 Q26 14 30 16 Q34 16 34 20"/><ellipse cx="20" cy="18" rx="1" ry="1.5"/><ellipse cx="24" cy="16.5" rx="1" ry="1.5"/><ellipse cx="28" cy="18" rx="1" ry="1.5"/><line x1="6" y1="40" x2="42" y2="40"/></svg>' },
  { name: '袁夫有酒', icon: 'wine', iconSvg: '<svg viewBox="0 0 48 48" fill="none" stroke="#333" stroke-width="1.5"><rect x="13" y="9" width="6" height="5"/><path d="M11 14 L11 38 Q11 41 14 41 L18 41 Q21 41 21 38 L21 14 Z"/><ellipse cx="16" cy="22" rx="3" ry="1"/><line x1="13" y1="26" x2="19" y2="26"/><path d="M28 14 L30 26 Q30 30 34 30 Q38 30 38 26 L40 14 Z"/><line x1="34" y1="30" x2="34" y2="40"/><line x1="30" y1="40" x2="38" y2="40"/></svg>' },
  { name: '袁夫味品', icon: 'can', iconSvg: '<svg viewBox="0 0 48 48" fill="none" stroke="#333" stroke-width="1.5"><ellipse cx="24" cy="13" rx="13" ry="3"/><path d="M11 13 L11 37 Q11 40 24 40 Q37 40 37 37 L37 13"/><path d="M14 19 Q14 17 16 17 L32 17 Q34 17 34 19 L34 33 Q34 35 32 35 L16 35 Q14 35 14 33 Z"/><circle cx="24" cy="26" r="3"/></svg>' },
  { name: '袁夫零食', icon: 'snack', iconSvg: '<svg viewBox="0 0 48 48" fill="none" stroke="#333" stroke-width="1.5"><circle cx="14" cy="18" r="6"/><circle cx="14" cy="18" r="1" fill="#333"/><circle cx="16" cy="15" r="0.8" fill="#333"/><circle cx="12" cy="20" r="0.8" fill="#333"/><path d="M22 14 L34 14 Q36 14 36 16 L36 36 Q36 38 34 38 L22 38 Q20 38 20 36 L20 16 Q20 14 22 14 Z"/><line x1="20" y1="20" x2="36" y2="20"/><line x1="20" y1="28" x2="36" y2="28"/></svg>' },
]

// Product illustrations
const noodleSvg = '<svg viewBox="0 0 168 130" style="position:absolute;inset:0;width:100%;height:100%"><rect x="0" y="80" width="168" height="50" fill="#e8d4b0" opacity="0.5"/><ellipse cx="84" cy="80" rx="50" ry="28" fill="#f5f0e8"/><ellipse cx="84" cy="78" rx="46" ry="20" fill="#f9e4b2"/><path d="M50 68 Q60 70 70 68 Q80 70 90 68 Q100 70 110 68 Q118 72 122 76" stroke="#f4d575" stroke-width="2" fill="none"/><path d="M52 72 Q62 74 72 72 Q82 74 92 72 Q102 74 112 72 Q120 76 124 80" stroke="#f4d575" stroke-width="2" fill="none"/><ellipse cx="110" cy="65" rx="10" ry="6" fill="#7fb069" transform="rotate(-20 110 65)"/><ellipse cx="115" cy="60" rx="7" ry="5" fill="#9bc168" transform="rotate(-25 115 60)"/><ellipse cx="70" cy="72" rx="8" ry="6" fill="#ffd16b"/><ellipse cx="70" cy="72" rx="3" ry="2" fill="#ff9f4f"/></svg>'
const jarSvg = '<svg viewBox="0 0 168 130" style="position:absolute;inset:0;width:100%;height:100%"><ellipse cx="84" cy="38" rx="34" ry="6" fill="#a87a4b"/><path d="M50 38 L50 102 Q50 110 84 110 Q118 110 118 102 L118 38" fill="#f5f0e6" stroke="#d0c2a8" stroke-width="1"/><ellipse cx="84" cy="38" rx="34" ry="6" fill="#c79a6e"/><rect x="58" y="60" width="52" height="28" fill="#fff" stroke="#c8a878" stroke-width="0.5"/><text x="84" y="78" text-anchor="middle" font-size="11" font-weight="600" fill="#5b4a30" font-family="serif">袁夫稻田</text><text x="84" y="86" text-anchor="middle" font-size="5" fill="#8a7860">净含量 450g</text><path d="M120 42 Q138 30 145 42 Q138 50 132 52 L130 60" stroke="#a87a4b" stroke-width="3" fill="none"/></svg>'
const ricebagSvg = '<svg viewBox="0 0 168 130" style="position:absolute;inset:0;width:100%;height:100%"><ellipse cx="20" cy="10" rx="40" ry="20" fill="#7fb069" opacity="0.4"/><ellipse cx="150" cy="20" rx="35" ry="25" fill="#5b8c2a" opacity="0.5"/><ellipse cx="160" cy="100" rx="30" ry="40" fill="#7fb069" opacity="0.3"/><path d="M50 50 L120 50 L122 110 L48 110 Z" fill="#fffbf0"/><path d="M50 50 Q60 42 70 50 Q80 42 90 50 Q100 42 110 50 Q115 46 120 50" fill="#fffbf0" stroke="#d8c8a0" stroke-width="0.5"/><rect x="62" y="68" width="46" height="32" fill="#fff" stroke="#bf4040" stroke-width="0.5"/><text x="85" y="84" text-anchor="middle" font-size="11" font-weight="600" fill="#bf4040" font-family="serif">袁夫</text><text x="85" y="94" text-anchor="middle" font-size="8" fill="#5b8c2a">稻田</text><ellipse cx="85" cy="50" rx="14" ry="3" fill="none" stroke="#d8c8a0" stroke-width="1"/></svg>'
const ricepackSvg = '<svg viewBox="0 0 168 130" style="position:absolute;inset:0;width:100%;height:100%"><rect x="0" y="60" width="168" height="70" fill="#f5e8c2" opacity="0.6"/><path d="M42 30 L96 30 L98 110 L40 110 Z" fill="#fffbf0"/><rect x="44" y="34" width="50" height="4" fill="#d0c2a8" rx="2"/><text x="69" y="68" text-anchor="middle" font-size="11" font-weight="600" fill="#bf4040" font-family="serif">袁夫</text><text x="69" y="80" text-anchor="middle" font-size="9" fill="#7a8a4a">稻田米</text><path d="M88 38 L138 38 L140 110 L86 110 Z" fill="#fff5db"/><rect x="90" y="42" width="46" height="3" fill="#c8b890" rx="1.5"/><text x="113" y="74" text-anchor="middle" font-size="9" font-weight="600" fill="#bf4040" font-family="serif">袁夫</text></svg>'
const porridgeSvg = '<svg viewBox="0 0 168 130" style="position:absolute;inset:0;width:100%;height:100%"><ellipse cx="84" cy="85" rx="55" ry="20" fill="#a87a4b" opacity="0.3"/><ellipse cx="84" cy="78" rx="48" ry="22" fill="#c79a6e"/><ellipse cx="84" cy="76" rx="44" ry="16" fill="#fffbf0"/><ellipse cx="84" cy="76" rx="38" ry="12" fill="#f5e8c2"/><ellipse cx="75" cy="72" rx="2" ry="1.5" fill="#fff"/><ellipse cx="92" cy="74" rx="2" ry="1.5" fill="#fff"/><ellipse cx="85" cy="78" rx="2" ry="1.5" fill="#fff"/><ellipse cx="100" cy="76" rx="2" ry="1.5" fill="#fff"/></svg>'

const newProducts = [
  { label: 'noodle', price: '12.8', bg: '#fef0e0', illustration: noodleSvg, presale: true },
  { label: 'jar', price: '28.8', bg: '#f5ede0', illustration: jarSvg, presale: false },
  { label: 'ricebag', price: '39.9', bg: '#e8efd8', illustration: ricebagSvg, presale: true },
  { label: 'porridge', price: '18.8', bg: '#fff0e0', illustration: porridgeSvg, presale: true },
]

const riceProducts = [
  { label: 'ricebag2', price: '68.0', bg: '#f5e8c2', illustration: ricebagSvg },
  { label: 'ricepack', price: '45.0', bg: '#fff5db', illustration: ricepackSvg },
  { label: 'noodle2', price: '22.8', bg: '#fef0e0', illustration: noodleSvg },
]

const condimentProducts = [
  { label: 'jar2', price: '28.0', bg: '#fde8d8', illustration: jarSvg },
  { label: 'jar3', price: '32.8', bg: '#f5ede0', illustration: jarSvg },
  { label: 'jar4', price: '18.0', bg: '#fde0d8', illustration: jarSvg },
]

// Cart
const cartItems = ref([
  { n: '稻田大米 5kg', img: '🌾', price: 68, qty: 1, sel: true },
  { n: '稻田酒 500ml', img: '🍶', price: 128, qty: 2, sel: true },
  { n: '自制蜂蜜 ×500g', img: '🍯', price: 38, qty: 1, sel: false },
])

const cartAllSel = computed(() => cartItems.value.every(it => it.sel))
const cartSelCount = computed(() => cartItems.value.filter(it => it.sel).length)
const cartTotal = computed(() => cartItems.value.filter(it => it.sel).reduce((s, it) => s + it.price * it.qty, 0))

function toggleAllCart() {
  const all = !cartAllSel.value
  cartItems.value.forEach(it => it.sel = all)
}

// Recharge
const rechargeTiles = [
  { amount: 50, gift: 0, tag: '' },
  { amount: 100, gift: 10, tag: '热门' },
  { amount: 200, gift: 30, tag: '' },
  { amount: 500, gift: 100, tag: '超值' },
  { amount: 1000, gift: 250, tag: '最划算' },
  { amount: 0, gift: 0, tag: '自定义' },
]

function doRecharge() {
  const amt = rechargeAmount.value === 0 ? parseFloat(rechargeCustom.value || '0') : rechargeAmount.value
  if (!amt || amt <= 0) { showToast('请输入充值金额'); return }
  balance.value += amt
  showToast('充值 ¥' + amt + ' 成功 ✓')
  subPage.value = 'recharge-records'
}

// Recharge records
const rechargeRecords = [
  { id: 1, icon: '+', iconBg: '#fff5e6', iconColor: '#fa8c16', title: '微信充值 (送 ¥30)', sub: '2026-05-14 10:22 · CZ20260514102230', amount: '+¥200.00', color: '#fa8c16' },
  { id: 2, icon: '+', iconBg: '#fff5e6', iconColor: '#fa8c16', title: '微信充值 (送 ¥10)', sub: '2026-05-13 16:30 · TF20260513163019', amount: '+¥100.00', color: '#fa8c16' },
  { id: 3, icon: '⤺', iconBg: '#fff2f0', iconColor: '#ff4d4f', title: '余额退款 (原路退回)', sub: '2026-05-12 09:12 · TK20260512091215', amount: '-¥50.00', color: '#ff4d4f' },
  { id: 4, icon: '+', iconBg: '#fff5e6', iconColor: '#fa8c16', title: '微信充值 (送 ¥10)', sub: '2026-05-10 09:15 · CZ20260510091523', amount: '+¥100.00', color: '#fa8c16' },
]

// Refund
const refundReasons = ['不再使用该账户', '重复充值', '金额输入错误', '其他']

function refundNext() {
  if (!refundReason.value) { showToast('请选择退款原因'); return }
  if (!refundAmount.value || refundAmount.value <= 0) { showToast('请输入退款金额'); return }
  if (refundAmount.value > 80) { showToast('超过可退金额'); return }
  refundStep.value = 2
}

function refundSubmit() {
  refundStep.value = 3
  balance.value -= refundAmount.value
}

// Refund records
const refundRecords = [
  { id: 'RF20260514142255', amt: 50, status: '待审核', tagClass: 'orange', sub: '已提交·等待园区审核' },
  { id: 'RF20260512091215', amt: 50, status: '已退款', tagClass: 'green', sub: '原路退回 · 微信账户' },
  { id: 'RF20260501163022', amt: 30, status: '已驳回', tagClass: 'gray', sub: '原因：未达退款门槛' },
]

// Client orders
const clientOrderFilter = ref('全部 (6)')
const clientOrderFilters = ['全部 (6)', '待付款 (1)', '待发货 (2)', '待收货 (1)', '已完成 (2)']

const clientOrders = [
  { shop: '袁夫米面专卖', items: [{ n: '稻田大米 5kg', img: '🌾', price: 68 }], status: '待发货', statusTag: 'red', total: 68 },
  { shop: '袁夫味品', items: [{ n: '自制蜂蜜 500g', img: '🍯', price: 38 }, { n: '稻田咸菜', img: '🫙', price: 18 }], status: '待收货', statusTag: 'red', total: 56 },
  { shop: '门票预订', items: [{ n: '园区门票 ×2', img: '🎫', price: 80 }], status: '待付款', statusTag: 'orange', total: 80 },
]

// Addresses
const addresses = [
  { n: '张三', phone: '138****5678', addr: '湖北省武汉市江夏区袁夫稻田园区民宿 B 栋 201', tag: '家', def: true },
  { n: '张三', phone: '138****5678', addr: '湖北省黄冈市黄梅县大河镇袁夫稻田旁', tag: '公司', def: false },
]

function handleLogout() {
  showToast('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.tab-content { background: #f5f5f5; min-height: 100%; }

/* Customer Service Banner */
.cs-banner { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: #fff; cursor: pointer; }
.cs-avatars { display: flex; }
.cs-avatar { width: 36px; height: 36px; border-radius: 50%; background: #fbe9e7; display: flex; align-items: center; justify-content: center; font-size: 20px; margin-right: -8px; border: 2px solid #fff; }
.cs-body { flex: 1; }
.cs-line1 { font-size: 13px; font-weight: 500; }
.cs-tag { display: inline-block; background: #fff2f0; color: #ff4d4f; font-size: 10px; padding: 1px 6px; border-radius: 8px; border: 1px solid #ffccc7; }
.cs-line2 { font-size: 11px; color: #999; margin-top: 4px; }
.cs-arrow { font-size: 18px; color: #ccc; }

/* Category Grid */
.cat-grid { display: grid; grid-template-columns: repeat(4, 1fr); padding: 14px 16px; background: #fff; gap: 8px; }
.cat-item { text-align: center; cursor: pointer; }
.cat-icon { width: 48px; height: 48px; margin: 0 auto 6px; display: flex; align-items: center; justify-content: center; }
.cat-icon svg { width: 40px; height: 40px; }
.cat-label { font-size: 11px; color: #333; }

/* Section Header */
.section-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px 10px; background: #f5f5f5; }
.sh-title { font-size: 16px; font-weight: 600; color: #333; display: flex; align-items: center; gap: 8px; }
.sh-bar { width: 3px; height: 14px; background: #ff4d4f; border-radius: 2px; display: inline-block; }
.sh-more { font-size: 12px; color: #999; cursor: pointer; display: flex; align-items: center; gap: 4px; }

/* Product Row */
.product-row { display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; padding: 0 16px; background: #f5f5f5; }
.product-card { background: #fff; border-radius: 10px; overflow: hidden; cursor: pointer; }
.pc-img { height: 140px; position: relative; overflow: hidden; }
.pc-brand { position: absolute; top: 8px; left: 8px; background: rgba(0,0,0,0.45); color: #fff; font-size: 10px; padding: 2px 6px; border-radius: 4px; }
.pc-presale { position: absolute; top: 8px; right: 8px; background: #ff4d4f; color: #fff; font-size: 10px; padding: 2px 6px; border-radius: 4px; }
.pc-meta { display: flex; justify-content: space-between; align-items: center; padding: 8px 10px 10px; }
.pc-price { font-size: 16px; font-weight: 600; color: #ff4d4f; }
.pc-cart-icon { color: #333; cursor: pointer; }

/* Member QR */
.memberqr-tab { background: linear-gradient(180deg, #1677ff 0%, #0958d9 100%); min-height: 100%; }
.mqr-panel { padding: 16px; color: #fff; }
.mqr-head { display: flex; align-items: center; gap: 12px; padding-bottom: 16px; }
.mqr-avatar { width: 44px; height: 44px; border-radius: 50%; background: rgba(255,255,255,0.95); color: #1677ff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 600; border: 2px solid rgba(255,255,255,0.5); }
.mqr-user { flex: 1; }
.mqr-name { font-size: 17px; font-weight: 600; }
.mqr-member { font-size: 12px; opacity: 0.75; margin-top: 2px; }
.mqr-balance-col { text-align: right; }
.mqr-balance { font-size: 26px; font-weight: 700; letter-spacing: -0.5px; }
.mqr-balance-label { font-size: 11px; opacity: 0.7; }

.mqr-card { background: #fff; border-radius: 14px; padding: 16px; color: #333; box-shadow: 0 8px 24px rgba(0,0,0,0.18); }
.mqr-tabs { display: flex; gap: 4px; padding: 3px; background: #f5f5f5; border-radius: 8px; margin-bottom: 16px; }
.mqr-tab { flex: 1; padding: 8px; border-radius: 6px; border: none; background: transparent; font-size: 14px; color: #666; font-weight: 500; cursor: pointer; }
.mqr-tab.active { background: #fff; color: #ff4d4f; box-shadow: 0 1px 4px rgba(0,0,0,0.08); }
.mqr-codebox { display: flex; flex-direction: column; align-items: center; }
.mqr-qrcode { padding: 12px; background: #fff; border-radius: 10px; border: 1px solid #f0f0f0; }

.mqr-tip { margin-top: 14px; font-size: 12px; color: #999; text-align: center; }
.mqr-actions { display: flex; gap: 10px; margin-top: 16px; padding-top: 14px; border-top: 0.5px dashed #f0f0f0; }
.mqr-action-btn { flex: 1; padding: 8px; border-radius: 8px; border: 0.5px solid #ddd; background: #fff; color: #333; font-size: 13px; cursor: pointer; }
.mqr-footer-tip { margin-top: 20px; font-size: 11px; opacity: 0.7; text-align: center; }

/* Cart */
.cart-tab { display: flex; flex-direction: column; }
.cart-header { background: #fff; padding: 10px 16px; font-size: 13px; display: flex; justify-content: space-between; color: #666; border-bottom: 1px solid #f0f0f0; }
.cart-item { background: #fff; padding: 14px 16px; display: flex; align-items: center; gap: 10px; border-bottom: 1px solid #f5f5f5; }
.cart-check { width: 22px; height: 22px; border-radius: 50%; border: 2px solid #d9d9d9; background: #fff; display: flex; align-items: center; justify-content: center; font-size: 12px; flex-shrink: 0; cursor: pointer; }
.cart-check.checked { border-color: #ff4d4f; background: #ff4d4f; color: #fff; }
.cart-img { width: 72px; height: 72px; border-radius: 8px; background: linear-gradient(135deg, #fff5e6, #f0d8a0); display: flex; align-items: center; justify-content: center; font-size: 36px; flex-shrink: 0; }
.cart-info { flex: 1; min-width: 0; }
.cart-name { font-size: 14px; font-weight: 500; color: #333; }
.cart-shop { font-size: 11px; color: #999; margin-top: 4px; }
.cart-bottom { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; }
.cart-price { font-size: 16px; font-weight: 600; color: #ff4d4f; }
.cart-qty { display: flex; align-items: center; gap: 8px; border: 1px solid #f0f0f0; border-radius: 14px; padding: 2px 8px; }
.cq-btn { border: none; background: transparent; width: 18px; height: 18px; font-size: 14px; color: #666; cursor: pointer; padding: 0; }
.cq-val { font-size: 13px; min-width: 14px; text-align: center; }
.cart-spacer { flex: 1; }
.cart-bar { background: #fff; border-top: 1px solid #f0f0f0; padding: 10px 16px; display: flex; align-items: center; gap: 10px; }
.cart-bar-label { font-size: 13px; color: #333; }
.cart-bar-right { flex: 1; display: flex; align-items: center; justify-content: flex-end; gap: 8px; }
.cart-bar-total { font-size: 13px; color: #666; }
.cart-total-price { color: #ff4d4f; font-size: 18px; font-weight: 600; }

/* Me */
.me-tab { background: #f5f5f5; }

/* Profile */
.profile-hero { background: linear-gradient(135deg, #1677ff, #0958d9); padding: 24px 16px 20px; color: #fff; }
.pf-row { display: flex; align-items: center; gap: 12px; }
.pf-avatar { width: 48px; height: 48px; border-radius: 50%; background: rgba(255,255,255,0.95); color: #1677ff; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.pf-info { flex: 1; }
.pf-name { font-size: 18px; font-weight: 600; }
.pf-grow { font-size: 12px; opacity: 0.75; }
.pf-action-btn { background: rgba(255,255,255,0.15); border: none; color: #fff; padding: 5px 14px; border-radius: 14px; font-size: 12px; cursor: pointer; }

/* Member Strip */
.member-strip { background: #fff; padding: 12px 16px; display: flex; align-items: center; gap: 6px; font-size: 14px; color: #ff4d4f; }
.m-link { color: #ff4d4f; font-size: 13px; cursor: pointer; }

/* Shared: Tags */
.mp-tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 10px; font-weight: 500; }
.mp-tag.green { background: #f0f9eb; color: #67c23a; }
.mp-tag.gray { background: #f5f5f5; color: #999; }
.mp-tag.red { background: #fff2f0; color: #ff4d4f; }
.mp-tag.orange { background: #fff7e6; color: #fa8c16; }

/* Shared: Buttons */
.mp-btn { background: #1677ff; color: #fff; border: none; border-radius: 8px; padding: 8px 20px; font-size: 13px; cursor: pointer; font-weight: 500; }
.mp-btn.sm { padding: 4px 12px; font-size: 12px; }
.mp-btn.outline { background: #fff; color: #1677ff; border: 1px solid #1677ff; }
.mp-btn.gray { background: #f5f5f5; color: #666; }

/* Shared: Stats Row */
.stat-row { display: grid; grid-template-columns: repeat(3, 1fr); background: #fff; margin: 12px; border-radius: 12px; padding: 14px 0; }
.stat-item { text-align: center; cursor: pointer; }
.stat-val { font-size: 18px; font-weight: 600; }
.stat-label { font-size: 12px; color: #999; margin-top: 4px; }

/* Shared: Menu */
.menu-card { background: #fff; border-radius: 12px; margin: 12px; overflow: hidden; }
.menu-row { display: flex; align-items: center; gap: 10px; padding: 14px 16px; font-size: 14px; cursor: pointer; border-bottom: 1px solid #f5f5f5; }
.menu-row:last-child { border-bottom: none; }
.menu-row.logout { color: #ff4d4f; }
.mm-icon { font-size: 18px; width: 24px; text-align: center; flex-shrink: 0; }
.mm-extra { flex: 1; text-align: right; font-size: 12px; color: #999; }
.mm-arrow { color: #ccc; flex-shrink: 0; }

/* Shared: Section */
.section-title { font-size: 14px; font-weight: 600; padding: 14px 16px 10px; color: #333; }

/* Shared: Record List */
.record-list { margin: 0 12px; }
.record-row { display: flex; align-items: center; gap: 10px; padding: 12px 14px; background: #fff; border-bottom: 1px solid #f0f0f0; }
.record-row:first-child { border-radius: 10px 10px 0 0; }
.record-row:last-child { border-radius: 0 0 10px 10px; border-bottom: none; }
.record-row:only-child { border-radius: 10px; }
.rr-icon { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; }
.rr-info { flex: 1; min-width: 0; }
.rr-title { font-size: 13px; color: #333; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.rr-sub { font-size: 11px; color: #999; margin-top: 2px; }
.rr-amount { font-size: 15px; font-weight: 600; flex-shrink: 0; }
.rr-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 10px; }

/* Recharge */
.recharge-hero { background: linear-gradient(135deg, #1677ff, #0958d9); padding: 24px 16px; color: #fff; }
.rh-label { font-size: 13px; opacity: 0.85; }
.rh-balance { font-size: 36px; font-weight: 600; margin-top: 4px; letter-spacing: -1px; }
.rh-sub { font-size: 12px; opacity: 0.85; margin-top: 4px; }
.recharge-card { background: #fff; margin: 14px 12px 12px; border-radius: 12px; padding: 16px; }
.rc-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.rc-title { font-size: 15px; font-weight: 600; }
.rc-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.rc-tile {
  position: relative; border: 1.5px solid #f0f0f0; border-radius: 10px;
  padding: 14px 8px; text-align: center; cursor: pointer; overflow: hidden;
}
.rc-tile.active { border-color: #ff4d4f; background: #fff2f0; }
.rc-ribbon { position: absolute; top: 0; right: 0; background: #ff4d4f; color: #fff; font-size: 9px; padding: 2px 8px; border-radius: 0 8px 0 8px; }
.rc-amount { font-size: 18px; font-weight: 700; }
.rc-amount-label { font-size: 14px; font-weight: 500; }
.rc-gift { font-size: 11px; color: #fa8c16; margin-top: 2px; min-height: 16px; }
.rc-custom { margin-top: 14px; padding: 12px; background: #fafafa; border-radius: 8px; display: flex; align-items: center; gap: 8px; }
.rc-custom-input { flex: 1; border: none; background: transparent; font-size: 18px; font-weight: 500; padding: 4px; outline: none; }
.recharge-payment { background: #fff; margin: 0 12px 12px; border-radius: 12px; padding: 14px 16px; display: flex; justify-content: space-between; align-items: center; font-size: 14px; }
.rp-wx { display: flex; align-items: center; gap: 8px; }
.recharge-rules { background: #fff; margin: 0 12px 12px; border-radius: 12px; padding: 14px 16px; font-size: 12px; color: #666; line-height: 1.7; }
.rrules-title { font-weight: 600; color: #333; margin-bottom: 6px; }

/* Bottom Bar */
.bottom-bar { background: #fff; padding: 12px 16px 20px; border-top: 1px solid #f0f0f0; display: flex; gap: 10px; }

/* Refund */
.refund-banner { background: #fff; padding: 16px; }
.rb-title { font-size: 16px; font-weight: 600; }
.rb-desc { font-size: 12px; color: #999; margin-top: 6px; line-height: 1.6; }
.refund-card { background: #fff; margin: 0 12px 12px; border-radius: 12px; padding: 16px; }
.rfc-row { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid #f5f5f5; font-size: 13px; }
.rfc-row:last-child { border-bottom: none; }
.rfc-val-lg { font-size: 18px; font-weight: 600; }
.rfc-val-red { font-size: 16px; color: #ff4d4f; font-weight: 500; }
.rfc-val-gray { font-size: 14px; color: #999; }
.rfc-divider { border-top: 1px dashed #ebebeb; padding-top: 14px; margin-top: 14px; }
.rfc-section-title { font-size: 14px; font-weight: 600; margin-bottom: 12px; }
.rfc-amount-input { display: flex; align-items: center; gap: 8px; padding: 10px; background: #fafafa; border-radius: 8px; }
.rfc-amount-field { flex: 1; border: none; background: transparent; font-size: 22px; font-weight: 600; padding: 0; outline: none; }
.rfc-hint { font-size: 11px; color: #999; margin-top: 8px; }
.rfc-reason { display: flex; align-items: center; gap: 10px; padding: 12px 4px; border-bottom: 1px solid #f5f5f5; cursor: pointer; }
.rfc-reason:last-child { border-bottom: none; }
.rfc-radio { width: 18px; height: 18px; border-radius: 50%; border: 2px solid #d9d9d9; display: flex; align-items: center; justify-content: center; }
.rfc-radio.checked { border-color: #ff4d4f; }
.rfc-dot { width: 8px; height: 8px; background: #ff4d4f; border-radius: 50%; }
.rfc-textarea { margin-top: 10px; width: 100%; border: 1px solid #ebebeb; border-radius: 8px; padding: 8px; font-size: 13px; min-height: 60px; outline: none; resize: none; box-sizing: border-box; }
.refund-agreement { background: #fff; margin: 0 12px 12px; border-radius: 12px; padding: 16px; font-size: 13px; color: #666; line-height: 1.7; }
.refund-progress { margin-top: 24px; padding: 14px; background: #fff7e6; border-radius: 10px; font-size: 12px; color: #d48806; line-height: 1.7; text-align: left; }

/* Done view */
.done-view { padding: 40px 20px 20px; text-align: center; }
.done-icon { width: 80px; height: 80px; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 42px; }
.done-icon.refund-done { background: linear-gradient(135deg, #67c23a, #52c41a); box-shadow: 0 8px 20px rgba(82,196,26,0.3); }
.done-title { font-size: 20px; font-weight: 600; margin-bottom: 8px; }
.done-sub { font-size: 13px; color: #999; }
.done-card { background: #fff; border-radius: 12px; padding: 20px; text-align: left; }
.dok-row { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px dashed #ebebeb; }
.dok-row:last-child { border-bottom: none; }
.dok-label { font-size: 13px; color: #999; width: 80px; flex-shrink: 0; }
.mono { font-family: monospace; font-size: 13px; }
.big-price { font-size: 18px; font-weight: 600; }

/* Shared: Order Card */
.order-card { margin: 12px; background: #fff; border-radius: 10px; padding: 14px; }
.oc-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.oc-user { font-size: 14px; font-weight: 500; }
.oc-footer { display: flex; align-items: center; justify-content: space-between; border-top: 1px dashed #f0f0f0; padding-top: 10px; }
.oc-no { font-size: 11px; color: #999; font-family: monospace; }
.oc-amount { font-size: 18px; font-weight: 600; }
.oc-actions { display: flex; align-items: center; gap: 10px; }

/* Client order item image */
.client-oi-img { width: 60px; height: 60px; background: linear-gradient(135deg, #fff5e6, #f0d8a0); border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 30px; flex-shrink: 0; }

/* Address */
.address-card { margin: 12px; }
.addr-header { display: flex; align-items: center; margin-bottom: 8px; }

/* Tab filter */
.tab-filter-row { display: flex; gap: 4px; padding: 6px; margin: 10px 12px 0; background: #f0f0f0; border-radius: 8px; }
.filter-btn { flex: 1; padding: 6px; border: none; border-radius: 6px; font-size: 12px; background: transparent; color: #666; cursor: pointer; }
.filter-btn.active { background: #fff; color: #1677ff; font-weight: 500; box-shadow: 0 1px 3px rgba(0,0,0,0.06); }
</style>
