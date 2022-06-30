import { Coin, DataLine, Document, Search } from '@element-plus/icons-vue'

const menus = [
  {
    path: '/',
    name: 'dashboard',
    icon: Search,
    label: '日志检索',
    color: '#0081dd',
    component: () => import('./views/dashboard.vue'),
  },
]
export default menus
