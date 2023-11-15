import { $emitter } from '~/pkgs';

// 首页路由
export const HomeRouter = {
  path: '/',
  redirect: '/glc/search',
};

$emitter.on('$HomeRouter', () => HomeRouter); // 首页路由提供解耦的获取方式

// 静态路由
export const staticRouters = [
  HomeRouter,
  {
    path: '/glc',
    redirect: '/glc/search',
  },
  {
    path: '/login',
    component: () => import('~/views/login/Login.vue'), // 登录路由，必有
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('~/components/error-pages/404.vue'),
  },
];
