import { createRouter, createWebHistory } from 'vue-router'

import menus from './menus'

const router = createRouter({
  history: createWebHistory(),
  routes: menus
})

export default router
