import { $emitter } from '~/pkgs';

// 首页路由
export const HomeRouter = {
  path: '/',
  redirect: '/glc/search',
  // meta: {
  //   title: '首页',
  //   componentName: 'MemoMain', // 根据实际更换（组件名称要符合规则才能被缓存）
  // },
  // component: () => import('~/components/layout/Layout.vue'),
  // children: [
  //   {
  //     path: '/',
  //     component: () => import('~/views/dashboard/Dashboard.vue'), // 根据实际更换
  //   },
  // ],
};

$emitter.on('$HomeRouter', () => HomeRouter); // 首页路由提供解耦的获取方式

// 静态路由
export const staticRouters = [
  HomeRouter,
  {
    path: '/glc',
    redirect: '/glc/search',
  },
  // {
  //   meta: {
  //     title: '日志检索',
  //     componentName: 'GlcMain',
  //   },
  //   component: () => import('~/components/layout/Layout.vue'),
  //   children: [
  //     {
  //       path: '/glc/search',
  //       component: () => import('~/views/system/glc/GlcMain.vue'),
  //     },
  //   ],
  // },
  {
    path: '/login',
    component: () => import('~/views/login/Login.vue'), // 登录路由，必有
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('~/components/error-pages/404.vue'),
    // meta: {
    //   title: '找不到页面',
    //   componentName: '404',
    // },
    // children: [
    //   {
    //     path: '/:pathMatch(.*)*',
    //     component: () => import('~/components/error-pages/404.vue'), // 根据实际更换
    //   },
    // ],
  },
];
