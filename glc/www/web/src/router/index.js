import { createRouter, createWebHistory } from 'vue-router';
import nProgress from '~/components/nprogress/nprogress';
import { initRouters } from './dynamic.routers';
import { useTokenStore, useMenuState, useTabsState } from '~/pkgs';
import { enableLogin } from '~/api';
import { staticRouters } from './static.routers';

// 白名单
const WHITE_LIST = ['/glc/login'];

// 此处不能初始化，否则无法持久化
let tabsState = null;
let tokenStore = null;
let menuState = null;

export const router = createRouter({
  history: createWebHistory(),
  routes: [], // 静态路由
  scrollBehavior: () => ({ left: 0, top: 0 }),
});

// 添加基本的静态路由
staticRouters.forEach(item => router.addRoute(item));

// ----------------------- 路由加载前 -----------------------
router.beforeEach(async (to, from, next) => {
  nProgress.start(); // 进度条开启

  !tokenStore && (tokenStore = useTokenStore());
  !tabsState && (tabsState = useTabsState());
  !menuState && (menuState = useMenuState());

  // 1、白名单页面可以随时跳转，通常是没有参数的静态路由，如登入页面这种
  if (WHITE_LIST.indexOf(to.path) !== -1) {
    next();
    return;
  }

  // 2、检查令牌（未登录时跳转登录页面）
  if (tokenStore.needLogin == 'unknow') {
    await enableLogin();
  }
  if (tokenStore.needLogin == 'true') {
    if (!tokenStore.token) {
      const redirect = encodeURIComponent(to.fullPath);
      return next({ path: '/glc/login', query: { redirect } }); // 跳转登录页，并附带目标页参数
    }
  }

  // 3、检查及初始化菜单数据和路由（有令牌也不见得一定有菜单数据和路由，需检查并初始化，取菜单数据则会检查令牌有效性）
  if (!menuState.aryMenu.length) {
    try {
      await initRouters();
    } catch (e) {
      tokenStore.$reset();
      return next({ path: '/glc/login' }); // 系统打开连菜单数据都取不到，多数是服务器停止之类，无法继续，清空登录令牌后跳转停留到登录页面
    }
    return next({ ...to, replace: true });
  }

  // 同步显示页签
  tabsState.showTab(to);
  next(); // 正常跳转
});

router.afterEach(() => {
  nProgress.done(); // 进度条关闭
});

router.onError(error => {
  nProgress.done(); // 进度条关闭
  console.error('路由错误', error.message);
});
