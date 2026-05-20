// ============================================================
// 启动引导 — 登录门禁 + 实时数据加载 (对接 Go 后端)
// 失败时回退到原型内置数据，保证演示永不空白
// ============================================================
window.LIVE = { overview: null, loaded: false };

function roleLabel(code) {
  return ({
    ADMIN_PLATFORM: '平台超级管理员', PLATFORM_OPER: '总部运营',
    PARK_ADMIN: '园区管理员', PARK_MANAGER: '园区经理',
    SHOP_ADMIN: '商户管理员', SHOP_CASHIER: '商户收银员',
  })[code] || code;
}

// ---- 登录门禁 ----
function showLogin(errMsg) {
  let el = document.getElementById('yfscLogin');
  if (!el) {
    el = document.createElement('div');
    el.id = 'yfscLogin';
    el.innerHTML = `
      <div class="yfsc-login-mask">
        <div class="yfsc-login-card">
          <div class="yfsc-login-logo">🌾 袁夫稻田</div>
          <div class="yfsc-login-sub">智慧园区综合管理平台 · 登录</div>
          <input id="yfscUser" class="yfsc-login-input" placeholder="用户名" value="admin" autocomplete="username">
          <input id="yfscPass" class="yfsc-login-input" type="password" placeholder="密码" value="admin123" autocomplete="current-password">
          <div id="yfscErr" class="yfsc-login-err"></div>
          <button id="yfscLoginBtn" class="yfsc-login-btn">登 录</button>
          <div class="yfsc-login-tip">演示账号：admin/admin123 · park-hm/123456 · hcct-mgr/123456</div>
        </div>
      </div>`;
    document.body.appendChild(el);
    const submit = async () => {
      const u = document.getElementById('yfscUser').value.trim();
      const p = document.getElementById('yfscPass').value;
      const errEl = document.getElementById('yfscErr');
      const btn = document.getElementById('yfscLoginBtn');
      errEl.textContent = '';
      btn.disabled = true; btn.textContent = '登录中...';
      try {
        await YFSC.login(u, p);
        el.remove();
        await startApp();
      } catch (e) {
        errEl.textContent = e.message || '登录失败';
        btn.disabled = false; btn.textContent = '登 录';
      }
    };
    document.getElementById('yfscLoginBtn').addEventListener('click', submit);
    el.querySelectorAll('.yfsc-login-input').forEach(i =>
      i.addEventListener('keydown', e => { if (e.key === 'Enter') submit(); }));
  }
  if (errMsg) document.getElementById('yfscErr').textContent = errMsg;
}

// ---- 用户徽标 + 退出 ----
function renderUserBadge() {
  if (!YFSC.user) return;
  // 同步原型左下角用户信息
  const av = document.getElementById('userAvatar');
  const nm = document.getElementById('userName');
  const rl = document.getElementById('userRole');
  if (av) av.textContent = (YFSC.user.name || 'U')[0];
  if (nm) nm.textContent = YFSC.user.name;
  if (rl) rl.textContent = roleLabel(YFSC.user.roleCode);

  let b = document.getElementById('yfscBadge');
  if (!b) {
    b = document.createElement('div');
    b.id = 'yfscBadge';
    b.className = 'yfsc-badge';
    document.body.appendChild(b);
  }
  b.innerHTML = `<span class="yfsc-dot"></span>实时数据 · ${YFSC.user.name}（${roleLabel(YFSC.user.roleCode)}）
    <button onclick="yfscLogout()">退出</button>`;
}
function yfscLogout() {
  YFSC.logout();
  location.reload();
}

// ---- 数据映射: API -> 原型字段 ----
const num = v => (v == null ? 0 : Number(v));
function mapParks(list) {
  return list.map(p => ({ id: p.id, name: p.name, code: p.code, addr: p.address || '',
    shopCount: num(p.shopCount), status: p.status, createdAt: p.createdAt }));
}
function mapShops(list) {
  return list.map(s => ({ id: s.id, name: s.name, addr: s.address || '', type: s.typeName || '其他',
    status: s.status, park: s.parkName || '', createdAt: s.createdAt }));
}
function mapUsers(list) {
  return list.map(u => ({ id: u.id, nick: u.nickname, phone: u.phone, level: u.memberLevel,
    balance: num(u.balance), consumption: num(u.consumptionTotal), faceState: u.faceStatus,
    regAt: u.registeredAt, openid: u.openid || '' }));
}
function mapRecharges(list) {
  return list.map(r => ({ id: r.id, phone: r.phone, amount: num(r.amount), gift: num(r.giftAmount),
    pay: r.payMethod, status: r.status, txn: r.txnNo || '', remark: r.remark || '', at: r.at }));
}
function mapEmployees(list) {
  return list.map(e => ({ name: e.name, username: e.username, phone: e.phone, park: e.park, shop: e.shop,
    role: e.role, position: e.position, online: !!e.online, createdAt: e.createdAt }));
}

// 通用: 直接可用 (字段已对齐)
const idFn = x => x;

async function loadLiveData() {
  // 并行拉取核心数据集；任一失败仅记录，不阻断
  const tasks = [
    ['parks', '/api/parks', mapParks, v => v.list],
    ['shops', '/api/shops?size=100', mapShops, v => v.list],
    ['users', '/api/users?size=100', mapUsers, v => v.list],
    ['orders', '/api/orders?size=100', idFn, v => v.list],
    ['recharges', '/api/recharges?size=100', mapRecharges, v => v.list],
    ['refunds', '/api/refunds', idFn, v => v.list],
    ['employees', '/api/employees', mapEmployees, v => v.list],
    ['activities', '/api/activities', idFn, v => v.list],
    ['mtRec', '/api/platform/verify-records?platform=meituan', idFn, v => v.list],
    ['mtStl', '/api/platform/settlements?platform=meituan', idFn, v => v.list],
    ['mtCfg', '/api/platform/store-configs?platform=meituan', idFn, v => v.list],
    ['dyRec', '/api/platform/verify-records?platform=douyin', idFn, v => v.list],
    ['dyStl', '/api/platform/settlements?platform=douyin', idFn, v => v.list],
    ['dyCfg', '/api/platform/store-configs?platform=douyin', idFn, v => v.list],
    ['mp', '/api/system/mini-programs', idFn, v => v.list],
    ['withdraws', '/api/withdraws', idFn, v => v.list],
    ['overview', '/api/statistics/overview', idFn, v => v],
  ];
  const results = await Promise.allSettled(tasks.map(t => YFSC.get(t[1])));
  const got = {};
  results.forEach((r, i) => {
    const [key, , mapper, pick] = tasks[i];
    if (r.status === 'fulfilled') {
      try { got[key] = mapper(pick(r.value)); } catch (e) { console.warn('map fail', key, e); }
    } else {
      console.warn('load fail', key, r.reason && r.reason.message);
    }
  });

  // 覆盖原型全局数据 (这些变量在 pc-data.js 中以 var 声明)
  if (got.parks) { PARKS = got.parks; }
  if (got.shops) { SHOPS = got.shops; }
  if (got.users) { USERS = got.users; }
  if (got.orders) { ORDERS = got.orders; }
  if (got.recharges) { RECHARGES = got.recharges; }
  if (got.refunds) { REFUNDS = got.refunds; }
  if (got.employees) { EMPLOYEES = got.employees; }
  if (got.activities) { ACTIVITIES = got.activities; }
  if (got.mtRec) { MEITUAN_VERIFY_RECORDS = got.mtRec; }
  if (got.mtStl) { MEITUAN_SETTLEMENTS = got.mtStl; }
  if (got.mtCfg) { MEITUAN_STORE_CONFIGS = got.mtCfg; }
  if (got.dyRec) { DOUYIN_VERIFY_RECORDS = got.dyRec; }
  if (got.dyStl) { DOUYIN_SETTLEMENTS = got.dyStl; }
  if (got.dyCfg) { DOUYIN_STORE_CONFIGS = got.dyCfg; }
  if (got.mp) { MINI_PROGRAMS = got.mp; }
  if (got.withdraws) { applyWithdraws(got.withdraws); }
  if (got.overview) { window.LIVE.overview = got.overview; }
  window.LIVE.loaded = true;
}

// 提现: 中文状态 -> 原型英文状态, 并按待审核/历史拆分
const WD_STATUS = { '待审核': 'pending', '已支付': 'paid', '已通过': 'paid', '已驳回': 'rejected' };
function applyWithdraws(list) {
  const mapped = list.map(w => ({
    id: w.id, shop: w.shop, shopPark: w.shopPark, amount: num(w.amount),
    bankName: w.bankName, bankAccount: w.bankAccount, bankHolder: w.bankHolder,
    status: WD_STATUS[w.status] || 'pending', applyAt: w.applyAt,
    reconciliationId: w.reconciliationId || '', reviewer: w.reviewer || '',
    reviewedAt: w.reviewedAt || '', paidAt: w.paidAt || '',
  }));
  WITHDRAW_APPLICATIONS = mapped.filter(w => w.status === 'pending');
  WITHDRAW_HISTORY = mapped.filter(w => w.status !== 'pending');
}

// 真实提现审核 (覆盖原型本地实现)
async function approveWithdraw(id) {
  try {
    await YFSC.post('/api/withdraws/' + id + '/approve', { remark: '审核通过' });
    const w = await YFSC.get('/api/withdraws');
    applyWithdraws(w.list);
    showToast('✅ 提现已通过 · 对账单已生成 · 已通知商户 ' + id);
    if (typeof renderPCContent === 'function') renderPCContent('withdraw');
  } catch (e) { showToast('操作失败：' + e.message); }
}
async function rejectWithdraw(id) {
  try {
    await YFSC.post('/api/withdraws/' + id + '/reject', { remark: '审核驳回' });
    const w = await YFSC.get('/api/withdraws');
    applyWithdraws(w.list);
    showToast('已驳回提现申请 ' + id);
    if (typeof renderPCContent === 'function') renderPCContent('withdraw');
  } catch (e) { showToast('操作失败：' + e.message); }
}

// ---- 视角切换时按角色重新登录拉取作用域数据 ----
const VIEW_ACCOUNTS = {
  'platform': ['admin', 'admin123'],
  'park-hm': ['park-hm', '123456'],
  'park-wh': ['park-wh', '123456'],
  'shop-hcct': ['hcct-mgr', '123456'],
};
async function switchViewLive(view) {
  const acct = VIEW_ACCOUNTS[view];
  if (!acct) return;
  try {
    await YFSC.login(acct[0], acct[1]);
    await loadLiveData();
    renderUserBadge();
    if (typeof renderPCNav === 'function') renderPCNav();
    if (typeof renderPCContent === 'function') renderPCContent(APP.page);
    showToast && showToast('已切换并加载实时数据：' + roleLabel(YFSC.user.roleCode));
  } catch (e) {
    showToast && showToast('切换失败：' + e.message);
  }
}

async function startApp() {
  try {
    await loadLiveData();
  } catch (e) {
    console.warn('实时数据加载失败，使用内置演示数据', e);
  }
  renderUserBadge();
  if (typeof init === 'function') init();
  // 接管视角切换：登录对应角色账号后重载作用域数据
  const sw = document.getElementById('viewSwitch');
  if (sw) {
    sw.addEventListener('change', (e) => { switchViewLive(e.target.value); });
  }
}

// bootApp 由 app.js 末尾调用
async function bootApp() {
  if (YFSC.token && YFSC.user) {
    await startApp();
  } else {
    showLogin();
  }
}
