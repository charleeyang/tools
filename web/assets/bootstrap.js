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
  return list.map(e => ({ id: e.id, name: e.name, username: e.username, phone: e.phone, park: e.park, shop: e.shop,
    role: e.role, roleCode: e.roleCode, position: e.position, online: !!e.online, status: e.status, createdAt: e.createdAt }));
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
    ['gates', '/api/gates', idFn, v => v.list],
    ['products', '/api/products', idFn, v => v.list],
    ['logs', '/api/system/logs', idFn, v => v.list],
    ['perms', '/api/permissions', idFn, v => v.list],
    ['finance', '/api/finance/summary', idFn, v => v],
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
  if (got.gates) { GATES = got.gates; }
  if (got.products) { PRODUCTS = got.products; }
  if (got.logs) { LOGS = got.logs; }
  if (got.perms) { PERMISSIONS = got.perms; }
  if (got.finance) { PARK_FINANCE = got.finance; }
  if (got.overview) { window.LIVE.overview = got.overview; }
  window.LIVE.loaded = true;
}

// ---- 闸机写操作 (真实 API, 覆盖原型) ----
async function reloadGates() {
  try { const r = await YFSC.get('/api/gates'); GATES = r.list; } catch (e) {}
}
async function submitGateCreate() {
  const name = (document.getElementById('gateName') || {}).value;
  const sn = (document.getElementById('gateSn') || {}).value;
  const type = (document.getElementById('gateType') || {}).value;
  const direction = (document.getElementById('gateDir') || {}).value;
  if (!name || !sn) { showToast('请填写名称与设备SN'); return; }
  try {
    await YFSC.post('/api/gates', { name, deviceSn: sn, type, direction });
    await reloadGates();
    closeModal();
    showToast('✅ 闸机已添加：' + name);
    if (typeof renderPCContent === 'function') renderPCContent('gate-list');
  } catch (e) { showToast('添加失败：' + e.message); }
}
async function toggleGate(id) {
  try {
    const d = await YFSC.put('/api/gates/' + id + '/toggle', {});
    await reloadGates();
    showToast(d.enabled ? '闸机已启用' : '闸机已禁用');
    if (typeof renderPCContent === 'function') renderPCContent('gate-list');
  } catch (e) { showToast('操作失败：' + e.message); }
}
async function deleteGate(id, name) {
  try {
    await YFSC.del('/api/gates/' + id);
    await reloadGates();
    showToast('已删除闸机：' + (name || id));
    if (typeof renderPCContent === 'function') renderPCContent('gate-list');
  } catch (e) { showToast('删除失败：' + e.message); }
}

// ---- 活动写操作 (真实 API) ----
async function reloadActivities() {
  try { const r = await YFSC.get('/api/activities'); ACTIVITIES = r.list; } catch (e) {}
}
async function auditActivity(id, approve) {
  try {
    const d = await YFSC.post('/api/activities/' + id + '/audit', { approve: !!approve });
    await reloadActivities();
    showToast('活动审核：' + d.status);
    if (typeof renderPCContent === 'function') renderPCContent('activity-list');
  } catch (e) { showToast('审核失败：' + e.message); }
}
async function endActivity(id) {
  try {
    await YFSC.post('/api/activities/' + id + '/end', {});
    await reloadActivities();
    showToast('活动已结束');
    if (typeof renderPCContent === 'function') renderPCContent('activity-list');
  } catch (e) { showToast('操作失败：' + e.message); }
}

// ---- 创建表单 (园区/店铺/客户/员工/活动) ----
const fval = id => { const e = document.getElementById(id); return e ? (e.value || '').trim() : ''; };

async function submitParkCreate() {
  const name = fval('parkName'), code = fval('parkCode');
  if (!name || !code) { showToast('请填写园区名称与编码'); return; }
  try {
    await YFSC.post('/api/parks', { name, code, address: fval('parkAddr'),
      contactPerson: fval('parkContact'), contactPhone: fval('parkPhone'), businessHours: fval('parkHours') });
    const r = await YFSC.get('/api/parks'); PARKS = mapParks(r.list);
    showToast('✅ 园区已创建：' + name);
    if (typeof goPage === 'function') goPage('park-list');
  } catch (e) { showToast('创建失败：' + e.message); }
}

async function openShopCreate() {
  let types = [];
  try { types = (await YFSC.get('/api/shop-types')).list; } catch (e) {}
  const parkOpts = (PARKS || []).map(p => `<option value="${p.id}">${p.name}</option>`).join('');
  const typeOpts = types.map(t => `<option value="${t.id}">${t.name}</option>`).join('');
  showModal({
    title: '添加店铺',
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">店铺名称 *</label><input id="shopName" class="input" placeholder="如 火车餐厅"></div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">所属园区 *</label><select id="shopPark" class="select" style="width:100%">${parkOpts}</select></div>
        <div class="form-item"><label class="form-label">店铺类型</label><select id="shopType" class="select" style="width:100%">${typeOpts}</select></div>
      </div>
      <div class="form-item"><label class="form-label">店铺地址</label><input id="shopAddr" class="input" placeholder="详细地址"></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button>
      <button class="btn btn-primary" onclick="submitShopCreate()">确定添加</button>`,
  });
}
async function submitShopCreate() {
  const name = fval('shopName');
  if (!name) { showToast('请填写店铺名称'); return; }
  try {
    await YFSC.post('/api/shops', { name, parkId: Number(fval('shopPark')),
      typeId: Number(fval('shopType')), address: fval('shopAddr') });
    const r = await YFSC.get('/api/shops?size=100'); SHOPS = mapShops(r.list);
    closeModal(); showToast('✅ 店铺已添加：' + name);
    if (typeof renderPCContent === 'function') renderPCContent('shop-list');
  } catch (e) { showToast('添加失败：' + e.message); }
}

function openUserCreate() {
  showModal({
    title: '手动添加客户',
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">昵称</label><input id="uNick" class="input" placeholder="客户昵称"></div>
      <div class="form-item"><label class="form-label">手机号 *</label><input id="uPhone" class="input" placeholder="11 位手机号"></div>
      <div class="form-item"><label class="form-label">会员等级</label>
        <select id="uLevel" class="select" style="width:100%"><option>普通用户</option><option>VIP1</option><option>VIP2</option><option>VIP3</option></select></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button>
      <button class="btn btn-primary" onclick="submitUserCreate()">确定添加</button>`,
  });
}
async function submitUserCreate() {
  const phone = fval('uPhone');
  if (!phone) { showToast('请填写手机号'); return; }
  try {
    await YFSC.post('/api/users', { nickname: fval('uNick'), phone, memberLevel: fval('uLevel') });
    const r = await YFSC.get('/api/users?size=100'); USERS = mapUsers(r.list);
    closeModal(); showToast('✅ 客户已添加');
    if (typeof renderPCContent === 'function') renderPCContent('user-list');
  } catch (e) { showToast('添加失败：' + e.message); }
}

async function openEmployeeCreate() {
  let positions = [], shops = [];
  try { positions = (await YFSC.get('/api/positions')).list; } catch (e) {}
  try { shops = (await YFSC.get('/api/shops?size=100')).list; } catch (e) {}
  const parkOpts = '<option value="">（不限/总部）</option>' + (PARKS || []).map(p => `<option value="${p.id}">${p.name}</option>`).join('');
  const shopOpts = '<option value="">（不限）</option>' + shops.map(s => `<option value="${s.id}">${s.name}</option>`).join('');
  const posOpts = '<option value="">（不限）</option>' + positions.map(p => `<option value="${p.id}">${p.name}</option>`).join('');
  showModal({
    title: '新增员工',
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-row">
        <div class="form-item"><label class="form-label">姓名 *</label><input id="eName" class="input"></div>
        <div class="form-item"><label class="form-label">登录账号 *</label><input id="eUser" class="input"></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">手机号</label><input id="ePhone" class="input"></div>
        <div class="form-item"><label class="form-label">角色 *</label>
          <select id="eRole" class="select" style="width:100%">
            <option value="PARK_ADMIN">园区管理员</option><option value="PARK_MANAGER">园区经理</option>
            <option value="SHOP_ADMIN">商户管理员</option><option value="SHOP_CASHIER">商户收银员</option>
            <option value="PLATFORM_OPER">总部运营</option><option value="ADMIN_PLATFORM">平台超级管理员</option>
          </select></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">所属园区</label><select id="ePark" class="select" style="width:100%">${parkOpts}</select></div>
        <div class="form-item"><label class="form-label">所属店铺</label><select id="eShop" class="select" style="width:100%">${shopOpts}</select></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">岗位</label><select id="ePos" class="select" style="width:100%">${posOpts}</select></div>
        <div class="form-item"><label class="form-label">初始密码</label><input id="ePwd" class="input" placeholder="默认 123456"></div>
      </div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button>
      <button class="btn btn-primary" onclick="submitEmployeeCreate()">确定添加</button>`,
  });
}
async function submitEmployeeCreate() {
  const name = fval('eName'), username = fval('eUser'), roleCode = fval('eRole');
  if (!name || !username) { showToast('请填写姓名与登录账号'); return; }
  try {
    await YFSC.post('/api/employees', {
      name, username, roleCode, phone: fval('ePhone'),
      parkId: Number(fval('ePark')) || 0, shopId: Number(fval('eShop')) || 0,
      positionId: Number(fval('ePos')) || 0, password: fval('ePwd') || '123456',
    });
    const r = await YFSC.get('/api/employees'); EMPLOYEES = mapEmployees(r.list);
    closeModal(); showToast('✅ 员工已创建：' + name);
    if (typeof renderPCContent === 'function') renderPCContent('employee-list');
  } catch (e) { showToast('创建失败：' + e.message); }
}

async function submitActivityCreate() {
  const title = fval('actTitle');
  if (!title) { showToast('请填写活动标题'); return; }
  try {
    await YFSC.post('/api/activities', { title, type: fval('actType'),
      parkScope: '全部园区', shopScope: '全部商户', startAt: fval('actStart'), endAt: fval('actEnd') });
    await reloadActivities();
    showToast('✅ 活动已创建并发起审核');
    if (typeof goPage === 'function') goPage('activity-list');
  } catch (e) { showToast('创建失败：' + e.message); }
}

// ---- 编辑 / 删除 (园区/店铺/客户/员工/活动) ----
async function reloadParks() { try { PARKS = mapParks((await YFSC.get('/api/parks')).list); } catch (e) {} }
async function reloadShops() { try { SHOPS = mapShops((await YFSC.get('/api/shops?size=100')).list); } catch (e) {} }
async function reloadUsers() { try { USERS = mapUsers((await YFSC.get('/api/users?size=100')).list); } catch (e) {} }
async function reloadEmployees() { try { EMPLOYEES = mapEmployees((await YFSC.get('/api/employees')).list); } catch (e) {} }
function rerender(page) { if (typeof renderPCContent === 'function') renderPCContent(page); }
let _pendingDel = null;
function runPendingDel() { closeModal(); if (_pendingDel) { const f = _pendingDel; _pendingDel = null; f(); } }
function confirmDel(msg, fn) {
  _pendingDel = fn;
  showModal({ title: '确认删除', body: `<div style="padding:10px 0;font-size:14px">${msg}</div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button>
      <button class="btn btn-primary" style="background:#ff4d4f;border-color:#ff4d4f" onclick="runPendingDel()">确认删除</button>` });
}

// 园区
function openParkEdit(id) {
  const p = (PARKS || []).find(x => x.id === id); if (!p) return;
  showModal({ title: '编辑园区 · ' + p.name,
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">园区名称</label><input id="epName" class="input" value="${p.name}"></div>
      <div class="form-item"><label class="form-label">地址</label><input id="epAddr" class="input" value="${p.addr || ''}"></div>
      <div class="form-item"><label class="form-label">状态</label><select id="epStatus" class="select" style="width:100%">
        <option ${p.status === '营业中' ? 'selected' : ''}>营业中</option><option ${p.status === '筹备中' ? 'selected' : ''}>筹备中</option><option ${p.status === '已关闭' ? 'selected' : ''}>已关闭</option></select></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitParkEdit(${id})">保存</button>` });
}
async function submitParkEdit(id) {
  try {
    await YFSC.put('/api/parks/' + id, { name: fval('epName'), address: fval('epAddr'), status: fval('epStatus') });
    await reloadParks(); closeModal(); showToast('✅ 园区已更新'); rerender('park-list');
  } catch (e) { showToast('更新失败：' + e.message); }
}
function deletePark(id, name) {
  confirmDel(`确认删除园区「${name}」？`, async () => {
    try { await YFSC.del('/api/parks/' + id); await reloadParks(); showToast('已删除园区：' + name); rerender('park-list'); }
    catch (e) { showToast('删除失败：' + e.message); }
  });
}

// 店铺
async function submitShopEdit(id) {
  const name = fval('editShopName'), addr = fval('editShopAddr'), typeName = fval('editShopType');
  let typeId = 0;
  try { const t = (await YFSC.get('/api/shop-types')).list.find(x => x.name === typeName); if (t) typeId = t.id; } catch (e) {}
  try {
    await YFSC.put('/api/shops/' + id, { name, address: addr, typeId });
    await reloadShops(); closeModal(); showToast('✅ 店铺「' + name + '」已保存'); rerender('shop-list');
  } catch (e) { showToast('保存失败：' + e.message); }
}
async function toggleShopStatus(id) {
  const s = (SHOPS || []).find(x => x.id === id); if (!s) return;
  const ns = s.status === '营业中' ? '休息中' : '营业中';
  try {
    await YFSC.put('/api/shops/' + id, { status: ns });
    await reloadShops(); showToast(ns === '营业中' ? '已上架' : '已下架'); rerender('shop-list');
  } catch (e) { showToast('操作失败：' + e.message); }
}

// 客户
function openUserEdit(id) {
  const u = (USERS || []).find(x => x.id === id); if (!u) return;
  showModal({ title: '编辑客户 · ' + (u.nick || u.phone),
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">昵称</label><input id="euNick" class="input" value="${u.nick || ''}"></div>
      <div class="form-item"><label class="form-label">手机号</label><input id="euPhone" class="input" value="${u.phone || ''}"></div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">会员等级</label><select id="euLevel" class="select" style="width:100%">
          ${['普通用户', 'VIP1', 'VIP2', 'VIP3'].map(l => `<option ${u.level === l ? 'selected' : ''}>${l}</option>`).join('')}</select></div>
        <div class="form-item"><label class="form-label">人脸状态</label><select id="euFace" class="select" style="width:100%">
          <option ${u.faceState === '已录入' ? 'selected' : ''}>已录入</option><option ${u.faceState === '未录入' ? 'selected' : ''}>未录入</option></select></div>
      </div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitUserEdit(${id})">保存</button>` });
}
async function submitUserEdit(id) {
  try {
    await YFSC.put('/api/users/' + id, { nickname: fval('euNick'), phone: fval('euPhone'), memberLevel: fval('euLevel'), faceStatus: fval('euFace') });
    await reloadUsers(); closeModal(); showToast('✅ 客户已更新'); rerender('user-list');
  } catch (e) { showToast('更新失败：' + e.message); }
}
function deleteUser(id, name) {
  confirmDel(`确认删除客户「${name}」？`, async () => {
    try { await YFSC.del('/api/users/' + id); await reloadUsers(); showToast('已删除客户：' + name); rerender('user-list'); }
    catch (e) { showToast('删除失败：' + e.message); }
  });
}

// 员工
async function openEmployeeEdit(id) {
  const e0 = (EMPLOYEES || []).find(x => x.id === id); if (!e0) return;
  let positions = [], shops = [];
  try { positions = (await YFSC.get('/api/positions')).list; } catch (e) {}
  try { shops = (await YFSC.get('/api/shops?size=100')).list; } catch (e) {}
  const roleOpts = [['ADMIN_PLATFORM', '平台超级管理员'], ['PLATFORM_OPER', '总部运营'], ['PARK_ADMIN', '园区管理员'], ['PARK_MANAGER', '园区经理'], ['SHOP_ADMIN', '商户管理员'], ['SHOP_CASHIER', '商户收银员']]
    .map(([c, n]) => `<option value="${c}" ${e0.roleCode === c ? 'selected' : ''}>${n}</option>`).join('');
  const parkOpts = '<option value="">（不限/总部）</option>' + (PARKS || []).map(p => `<option value="${p.id}">${p.name}</option>`).join('');
  const shopOpts = '<option value="">（不限）</option>' + shops.map(s => `<option value="${s.id}">${s.name}</option>`).join('');
  const posOpts = '<option value="">（不限）</option>' + positions.map(p => `<option value="${p.id}">${p.name}</option>`).join('');
  showModal({ title: '编辑员工 · ' + e0.name,
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-row">
        <div class="form-item"><label class="form-label">姓名</label><input id="eeName" class="input" value="${e0.name}"></div>
        <div class="form-item"><label class="form-label">手机号</label><input id="eePhone" class="input" value="${e0.phone || ''}"></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">角色</label><select id="eeRole" class="select" style="width:100%">${roleOpts}</select></div>
        <div class="form-item"><label class="form-label">状态</label><select id="eeStatus" class="select" style="width:100%">
          <option ${e0.status === '启用' ? 'selected' : ''}>启用</option><option ${e0.status === '禁用' ? 'selected' : ''}>禁用</option></select></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">所属园区</label><select id="eePark" class="select" style="width:100%">${parkOpts}</select></div>
        <div class="form-item"><label class="form-label">所属店铺</label><select id="eeShop" class="select" style="width:100%">${shopOpts}</select></div>
      </div>
      <div class="form-item"><label class="form-label">岗位</label><select id="eePos" class="select" style="width:100%">${posOpts}</select></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitEmployeeEdit(${id})">保存</button>` });
}
async function submitEmployeeEdit(id) {
  try {
    await YFSC.put('/api/employees/' + id, {
      name: fval('eeName'), phone: fval('eePhone'), roleCode: fval('eeRole'), status: fval('eeStatus'),
      parkId: Number(fval('eePark')) || 0, shopId: Number(fval('eeShop')) || 0, positionId: Number(fval('eePos')) || 0,
    });
    await reloadEmployees(); closeModal(); showToast('✅ 员工已更新'); rerender('employee-list');
  } catch (e) { showToast('更新失败：' + e.message); }
}
function deleteEmployee(id, name) {
  confirmDel(`确认删除员工「${name}」？`, async () => {
    try { await YFSC.del('/api/employees/' + id); await reloadEmployees(); showToast('已删除员工：' + name); rerender('employee-list'); }
    catch (e) { showToast('删除失败：' + e.message); }
  });
}

// 活动
function openActivityEdit(id) {
  const a = (ACTIVITIES || []).find(x => x.id === id); if (!a) return;
  showModal({ title: '编辑活动 · ' + a.title,
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">活动标题</label><input id="eaTitle" class="input" value="${a.title}"></div>
      <div class="form-item"><label class="form-label">活动类型</label><select id="eaType" class="select" style="width:100%">
        ${['优惠券活动', '满减活动', '限时折扣', '节日活动', '会员专享'].map(t => `<option ${a.type === t ? 'selected' : ''}>${t}</option>`).join('')}</select></div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">开始时间</label><input id="eaStart" class="input" value="${a.startAt || ''}"></div>
        <div class="form-item"><label class="form-label">结束时间</label><input id="eaEnd" class="input" value="${a.endAt || ''}"></div>
      </div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitActivityEdit(${id})">保存</button>` });
}
async function submitActivityEdit(id) {
  try {
    await YFSC.put('/api/activities/' + id, { title: fval('eaTitle'), type: fval('eaType'), startAt: fval('eaStart'), endAt: fval('eaEnd') });
    await reloadActivities(); closeModal(); showToast('✅ 活动已更新'); rerender('activity-list');
  } catch (e) { showToast('更新失败：' + e.message); }
}
function deleteActivity(id, title) {
  confirmDel(`确认删除活动「${title}」？`, async () => {
    try { await YFSC.del('/api/activities/' + id); await reloadActivities(); showToast('已删除活动：' + title); rerender('activity-list'); }
    catch (e) { showToast('删除失败：' + e.message); }
  });
}

// ---- 商品管理 ----
async function reloadProducts() { try { PRODUCTS = (await YFSC.get('/api/products')).list; } catch (e) {} }
async function openProductCreate() {
  let shops = [];
  try { shops = (await YFSC.get('/api/shops?size=100')).list; } catch (e) {}
  const shopOpts = shops.map(s => `<option value="${s.id}">${s.name}</option>`).join('');
  showModal({
    title: '新增商品',
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">商品名称 *</label><input id="prName" class="input"></div>
      <div class="form-item"><label class="form-label">所属店铺 *</label><select id="prShop" class="select" style="width:100%">${shopOpts}</select></div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">售价 *</label><input id="prPrice" class="input" type="number" placeholder="0.00"></div>
        <div class="form-item"><label class="form-label">原价</label><input id="prOrig" class="input" type="number" placeholder="0.00"></div>
      </div>
      <div class="form-item"><label class="form-label">库存</label><input id="prStock" class="input" type="number" value="0"></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitProductCreate()">确定添加</button>`,
  });
}
async function submitProductCreate() {
  const name = fval('prName');
  if (!name) { showToast('请填写商品名称'); return; }
  try {
    await YFSC.post('/api/products', { name, shopId: Number(fval('prShop')),
      price: Number(fval('prPrice')) || 0, originalPrice: Number(fval('prOrig')) || 0, stock: Number(fval('prStock')) || 0 });
    await reloadProducts(); closeModal(); showToast('✅ 商品已添加：' + name); rerender('product-list');
  } catch (e) { showToast('添加失败：' + e.message); }
}
function openProductEdit(id) {
  const p = (PRODUCTS || []).find(x => x.id === id); if (!p) return;
  showModal({
    title: '编辑商品 · ' + p.name,
    body: `<div style="display:flex;flex-direction:column;gap:14px;padding:6px 0">
      <div class="form-item"><label class="form-label">商品名称</label><input id="eprName" class="input" value="${p.name}"></div>
      <div class="form-row">
        <div class="form-item"><label class="form-label">售价</label><input id="eprPrice" class="input" type="number" value="${p.price}"></div>
        <div class="form-item"><label class="form-label">原价</label><input id="eprOrig" class="input" type="number" value="${p.originalPrice || 0}"></div>
      </div>
      <div class="form-item"><label class="form-label">库存</label><input id="eprStock" class="input" type="number" value="${p.stock}"></div>
    </div>`,
    footer: `<button class="btn" onclick="closeModal()">取消</button><button class="btn btn-primary" onclick="submitProductEdit(${id})">保存</button>`,
  });
}
async function submitProductEdit(id) {
  try {
    await YFSC.put('/api/products/' + id, { name: fval('eprName'),
      price: Number(fval('eprPrice')) || 0, originalPrice: Number(fval('eprOrig')) || 0, stock: Number(fval('eprStock')) });
    await reloadProducts(); closeModal(); showToast('✅ 商品已更新'); rerender('product-list');
  } catch (e) { showToast('更新失败：' + e.message); }
}
async function toggleProductStatus(id) {
  const p = (PRODUCTS || []).find(x => x.id === id); if (!p) return;
  const ns = p.status === '已上架' ? '已下架' : '已上架';
  try { await YFSC.put('/api/products/' + id, { status: ns }); await reloadProducts(); showToast(ns); rerender('product-list'); }
  catch (e) { showToast('操作失败：' + e.message); }
}
function deleteProduct(id, name) {
  confirmDel(`确认删除商品「${name}」？`, async () => {
    try { await YFSC.del('/api/products/' + id); await reloadProducts(); showToast('已删除商品：' + name); rerender('product-list'); }
    catch (e) { showToast('删除失败：' + e.message); }
  });
}

// ---- 小程序配置编辑 ----
async function submitMpEdit(id) {
  try {
    await YFSC.put('/api/system/mini-programs/' + id, { callbackUrl: fval('mpCallback'), mchId: fval('mpMch'), version: fval('mpVer') });
    try { MINI_PROGRAMS = (await YFSC.get('/api/system/mini-programs')).list; } catch (e) {}
    closeModal(); showToast('✅ 配置已保存'); rerender('mp-config');
  } catch (e) { showToast('保存失败：' + e.message); }
}

// ---- 权限矩阵: 点击单元格循环 可操作->只读->隐藏 ----
async function cyclePerm(roleCode, moduleKey) {
  const cur = (PERMISSIONS || []).find(p => p.roleCode === roleCode && p.moduleKey === moduleKey);
  const order = ['可操作', '只读', '隐藏'];
  const idx = cur ? order.indexOf(cur.visibility) : 2;
  const next = order[(idx + 1) % 3];
  try {
    await YFSC.put('/api/permissions', { roleCode, moduleKey, visibility: next });
    PERMISSIONS = (await YFSC.get('/api/permissions')).list;
    showToast(roleCode + ' / ' + moduleKey + ' → ' + next);
    rerender('role-perm');
  } catch (e) { showToast('调整失败：' + e.message); }
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
