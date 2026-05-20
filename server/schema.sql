-- 袁夫稻田智慧园区综合管理平台 — 数据库结构 (SQLite 版, 对应 PRD V3.0 27 张表)
-- 由 MySQL DDL 适配: BIGINT UNSIGNED -> INTEGER, JSON -> TEXT, ON UPDATE 由应用层维护

PRAGMA foreign_keys = ON;

-- ============ 基础数据 (001) ============
CREATE TABLE IF NOT EXISTS park (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    address TEXT,
    contact_person TEXT,
    contact_phone TEXT,
    business_hours TEXT,
    logo_url TEXT,
    description TEXT,
    status TEXT NOT NULL DEFAULT '营业中',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT
);

CREATE TABLE IF NOT EXISTS shop_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    description TEXT,
    status TEXT NOT NULL DEFAULT '启用',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_by TEXT,
    updated_at TEXT
);

CREATE TABLE IF NOT EXISTS shop (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    park_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    type_id INTEGER,
    status TEXT NOT NULL DEFAULT '营业中',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    FOREIGN KEY (park_id) REFERENCES park(id),
    FOREIGN KEY (type_id) REFERENCES shop_type(id)
);
CREATE INDEX IF NOT EXISTS idx_shop_park ON shop(park_id);
CREATE INDEX IF NOT EXISTS idx_shop_type ON shop(type_id);

CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT,
    phone TEXT,
    openid TEXT UNIQUE,
    avatar_url TEXT,
    member_level TEXT NOT NULL DEFAULT '普通用户',
    balance REAL NOT NULL DEFAULT 0,
    consumption_total REAL NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0,
    face_status TEXT NOT NULL DEFAULT '未录入',
    registered_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_user_phone ON user(phone);
CREATE INDEX IF NOT EXISTS idx_user_level ON user(member_level);

-- ============ 业务数据 (002) ============
CREATE TABLE IF NOT EXISTS product_category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shop_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    icon TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT '启用',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (shop_id) REFERENCES shop(id)
);

CREATE TABLE IF NOT EXISTS product (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER,
    shop_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    original_price REAL,
    image TEXT,
    description TEXT,
    specs TEXT,
    stock INTEGER NOT NULL DEFAULT 0,
    sales_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT '已上架',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    FOREIGN KEY (shop_id) REFERENCES shop(id)
);
CREATE INDEX IF NOT EXISTS idx_p_shop ON product(shop_id);

CREATE TABLE IF NOT EXISTS "order" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_no TEXT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    shop_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    pay_method TEXT,
    status TEXT NOT NULL DEFAULT '待支付',
    source TEXT NOT NULL DEFAULT 'self',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (shop_id) REFERENCES shop(id)
);
CREATE INDEX IF NOT EXISTS idx_o_user ON "order"(user_id);
CREATE INDEX IF NOT EXISTS idx_o_shop ON "order"(shop_id);
CREATE INDEX IF NOT EXISTS idx_o_status ON "order"(status);

CREATE TABLE IF NOT EXISTS recharge (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biz_no TEXT,
    user_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    gift_amount REAL NOT NULL DEFAULT 0,
    pay_method TEXT,
    txn_no TEXT,
    remark TEXT,
    status TEXT NOT NULL DEFAULT '成功',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id)
);
CREATE INDEX IF NOT EXISTS idx_r_user ON recharge(user_id);

CREATE TABLE IF NOT EXISTS refund (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biz_no TEXT,
    user_id INTEGER NOT NULL,
    order_id INTEGER,
    order_no TEXT,
    amount REAL NOT NULL,
    type TEXT NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT '待审核',
    applied_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    reviewed_at TEXT,
    reviewer TEXT,
    review_remark TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id)
);
CREATE INDEX IF NOT EXISTS idx_ref_user ON refund(user_id);

CREATE TABLE IF NOT EXISTS withdraw (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biz_no TEXT,
    shop_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    bank_name TEXT NOT NULL,
    bank_account TEXT NOT NULL,
    bank_holder TEXT NOT NULL,
    reconciliation_id TEXT,
    status TEXT NOT NULL DEFAULT '待审核',
    applied_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    reviewed_at TEXT,
    reviewer TEXT,
    review_remark TEXT,
    paid_at TEXT,
    FOREIGN KEY (shop_id) REFERENCES shop(id)
);
CREATE INDEX IF NOT EXISTS idx_w_shop ON withdraw(shop_id);
CREATE INDEX IF NOT EXISTS idx_w_status ON withdraw(status);

-- ============ 员工与权限 (003) ============
CREATE TABLE IF NOT EXISTS position (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    responsibilities TEXT,
    requirements TEXT,
    status TEXT NOT NULL DEFAULT '启用',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS employee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    phone TEXT,
    gender TEXT,
    park_id INTEGER,
    shop_id INTEGER,
    position_id INTEGER,
    role_code TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT '启用',
    online_status INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    FOREIGN KEY (park_id) REFERENCES park(id),
    FOREIGN KEY (shop_id) REFERENCES shop(id),
    FOREIGN KEY (position_id) REFERENCES position(id)
);

CREATE TABLE IF NOT EXISTS role_permission (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_code TEXT NOT NULL,
    module_key TEXT NOT NULL,
    visibility TEXT NOT NULL DEFAULT '可操作',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT,
    UNIQUE (role_code, module_key)
);

-- ============ 闸机 (004) ============
CREATE TABLE IF NOT EXISTS gate (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    park_id INTEGER,
    name TEXT NOT NULL,
    type TEXT,
    direction TEXT,
    device_sn TEXT NOT NULL UNIQUE,
    enabled INTEGER NOT NULL DEFAULT 1,
    online_status INTEGER NOT NULL DEFAULT 0,
    today_pass INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS face (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    feature_id TEXT,
    authorized_gates TEXT,
    status TEXT NOT NULL DEFAULT '已录入',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE IF NOT EXISTS entry_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    gate_id INTEGER,
    direction TEXT,
    method TEXT,
    passed_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (gate_id) REFERENCES gate(id)
);

-- ============ 活动 (005) ============
CREATE TABLE IF NOT EXISTS activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    type TEXT,
    park_scope TEXT,
    shop_scope TEXT,
    start_at TEXT,
    end_at TEXT,
    status TEXT NOT NULL DEFAULT '待审核',
    exposure INTEGER NOT NULL DEFAULT 0,
    joins INTEGER NOT NULL DEFAULT 0,
    orders INTEGER NOT NULL DEFAULT 0,
    verify_rate TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS activity_park (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER NOT NULL,
    park_id INTEGER NOT NULL,
    FOREIGN KEY (activity_id) REFERENCES activity(id),
    FOREIGN KEY (park_id) REFERENCES park(id)
);

CREATE TABLE IF NOT EXISTS activity_shop (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER NOT NULL,
    shop_id INTEGER NOT NULL,
    FOREIGN KEY (activity_id) REFERENCES activity(id),
    FOREIGN KEY (shop_id) REFERENCES shop(id)
);

CREATE TABLE IF NOT EXISTS coupon (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER,
    name TEXT NOT NULL,
    type TEXT,
    discount_amount REAL,
    threshold REAL,
    total_count INTEGER NOT NULL DEFAULT 0,
    used_count INTEGER NOT NULL DEFAULT 0,
    valid_from TEXT,
    valid_until TEXT,
    status TEXT NOT NULL DEFAULT '启用',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS user_coupon (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    coupon_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT '未使用',
    received_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    used_at TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (coupon_id) REFERENCES coupon(id)
);

CREATE TABLE IF NOT EXISTS points_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    change_points INTEGER NOT NULL,
    type TEXT,
    remark TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

-- ============ 第三方平台 (006) ============
CREATE TABLE IF NOT EXISTS platform_verify_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biz_no TEXT,
    platform TEXT NOT NULL,
    coupon_code TEXT NOT NULL,
    product TEXT,
    origin_price REAL,
    sale_price REAL,
    user_name TEXT,
    phone TEXT,
    shop_id INTEGER,
    shop_name TEXT,
    coupon_type TEXT,
    verify_count INTEGER NOT NULL DEFAULT 1,
    max_count INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT '已核销',
    verify_by TEXT,
    verify_at TEXT,
    refundable INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_pvr_platform ON platform_verify_records(platform);

CREATE TABLE IF NOT EXISTS platform_settlements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    biz_no TEXT,
    platform TEXT NOT NULL,
    shop_name TEXT,
    settle_date TEXT,
    order_count INTEGER NOT NULL DEFAULT 0,
    total_amount REAL NOT NULL DEFAULT 0,
    commission REAL NOT NULL DEFAULT 0,
    net_amount REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT '待结算',
    period TEXT,
    rate TEXT
);
CREATE INDEX IF NOT EXISTS idx_pstl_platform ON platform_settlements(platform);

CREATE TABLE IF NOT EXISTS platform_store_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    shop_id INTEGER,
    shop_name TEXT,
    app_auth_token TEXT,
    status TEXT NOT NULL DEFAULT '待绑定',
    bind_at TEXT,
    last_sync TEXT,
    verify_count INTEGER NOT NULL DEFAULT 0,
    scope TEXT
);
CREATE INDEX IF NOT EXISTS idx_psc_platform ON platform_store_configs(platform);

-- ============ 系统 (007) ============
CREATE TABLE IF NOT EXISTS mini_program (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT,
    type_label TEXT,
    app_secret TEXT,
    callback_url TEXT,
    mch_id TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    version TEXT,
    published_at TEXT
);

CREATE TABLE IF NOT EXISTS operation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operator TEXT,
    role_code TEXT,
    module TEXT,
    action TEXT,
    detail TEXT,
    ip TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);
