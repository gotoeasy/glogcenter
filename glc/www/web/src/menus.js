import { Coin, DataAnalysis, DataLine, Discount, Document, Search, Setting } from '@element-plus/icons-vue'

const menus = [
  {
    path: '/',
    redirect: '/glc/search',
    hidden: true,
  },{
    path: '/glc',
    redirect: '/glc/search',
    hidden: true,
  },{
    path: '/glc/search',
    name: 'search',
    icon: Search,
    label: '日志检索',
    color: '#0081dd',
    component: () => import('./views/search.vue'),
  },{
    path: '/glc/storages',
    name: 'storages',
    icon: Coin,
    label: '日志仓管理',
    color: '#0081dd',
    component: () => import('./views/storages.vue'),
  },
]
export default menus
