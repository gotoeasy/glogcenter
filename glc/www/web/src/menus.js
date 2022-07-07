import { Coin, DataAnalysis, DataLine, Discount, Document, Search, Setting } from '@element-plus/icons-vue'

const menus = [
  {
    path: '/',
    name: 'dashboard',
    icon: Search,
    label: '日志检索',
    color: '#0081dd',
    component: () => import('./views/dashboard.vue'),
  },
  {
    path: '/',
    name: 'storages',
    icon: Coin,
    label: '日志仓管理',
    color: '#0081dd',
    component: () => import('./views/storages.vue'),
  },
]
export default menus
