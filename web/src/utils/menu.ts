export interface MenuItem {
  id: string
  label: string
  icon: string
  page?: string
  sub?: { id: string; page: string; label: string }[]
}

export interface MenuSection {
  section: string
  items: MenuItem[]
}

export const FULL_MENU: MenuSection[] = [
  {
    section: '工作台',
    items: [
      { id: 'overview', label: '首页概览', icon: '⌂', page: 'overview' },
    ],
  },
  {
    section: '园区管理',
    items: [
      { id: 'parks', label: '园区列表', icon: '▤', page: 'parks' },
      { id: 'shops', label: '店铺管理', icon: '◫', page: 'shops' },
    ],
  },
  {
    section: '业务管理',
    items: [
      { id: 'orders', label: '消费订单', icon: '◈', page: 'orders' },
      { id: 'products', label: '商品管理', icon: '☷', page: 'products' },
      { id: 'users', label: '客户管理', icon: '👤', page: 'users' },
      { id: 'finance', label: '财务管理', icon: '¥', page: 'finance' },
    ],
  },
  {
    section: '员工与设备',
    items: [
      { id: 'employees', label: '员工管理', icon: '👥', page: 'employees' },
      { id: 'gates', label: '闸机管理', icon: '⊞', page: 'gates' },
    ],
  },
  {
    section: '运营管理',
    items: [
      { id: 'activities', label: '活动管理', icon: '★', page: 'activities' },
      { id: 'platform', label: '第三方平台', icon: '☰', page: 'platform' },
    ],
  },
  {
    section: '系统设置',
    items: [
      { id: 'settings', label: '权限设置', icon: '⚙', page: 'settings' },
    ],
  },
]

export const VIEW_TO_ROLE: Record<string, string> = {
  platform: 'platform',
  'park-hm': 'park-hm',
  'park-wh': 'park-wh',
  'shop-hcct': 'shop-hcct',
}

export const ROLES: Record<string, { avatar: string; user: string; name: string }> = {
  platform: { avatar: '张', user: '张三', name: '平台超级管理员' },
  'park-hm': { avatar: '李', user: '李四', name: '黄梅园区管理员' },
  'park-wh': { avatar: '王', user: '王五', name: '武汉园区管理员' },
  'shop-hcct': { avatar: '赵', user: '赵六', name: '火车餐厅商户' },
}
