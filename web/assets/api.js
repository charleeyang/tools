// ============================================================
// 实时 API 客户端 — 对接 Go 后端 (/api/*)
// ============================================================
const YFSC = {
  base: '',
  token: localStorage.getItem('yfsc_token') || '',
  user: JSON.parse(localStorage.getItem('yfsc_user') || 'null'),
  permissions: JSON.parse(localStorage.getItem('yfsc_perms') || '{}'),

  async request(method, path, body) {
    const headers = { 'Content-Type': 'application/json' };
    if (this.token) headers['Authorization'] = 'Bearer ' + this.token;
    const res = await fetch(this.base + path, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
    let data = {};
    try { data = await res.json(); } catch (e) {}
    if (res.status === 401) {
      // token 失效，回到登录
      YFSC.logout();
      throw new Error(data.message || '登录已失效');
    }
    if (data.code !== 0) {
      const err = new Error(data.message || ('请求失败 ' + res.status));
      err.code = data.code;
      throw err;
    }
    return data.data;
  },
  get(p) { return this.request('GET', p); },
  post(p, b) { return this.request('POST', p, b); },
  put(p, b) { return this.request('PUT', p, b); },
  del(p) { return this.request('DELETE', p); },

  async login(username, password) {
    const data = await this.post('/api/auth/login', { username, password });
    this.token = data.token;
    this.user = data.user;
    this.permissions = data.permissions || {};
    localStorage.setItem('yfsc_token', this.token);
    localStorage.setItem('yfsc_user', JSON.stringify(this.user));
    localStorage.setItem('yfsc_perms', JSON.stringify(this.permissions));
    return data;
  },

  logout() {
    this.token = ''; this.user = null; this.permissions = {};
    localStorage.removeItem('yfsc_token');
    localStorage.removeItem('yfsc_user');
    localStorage.removeItem('yfsc_perms');
  },
};
