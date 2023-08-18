import { getUserMenu } from '~/api/api';
import { router } from '~/router';
import { GxIframe } from '~/pkgs';

const modules = import.meta.glob('~/views/**/**.vue');

// 页面菜单数据转页面路由
function toRouter(aryRouter, menuData) {
  const o = {};

  menuData.forEach(item => {
    if (item.children) {
      toRouter(aryRouter, item.children);
    } else {
      let name = item.name || item.component.substring(item.component.lastIndexOf('/') + 1);
      let componentName = item.component.substring(item.component.lastIndexOf('/') + 1);
      let iframeSrc = '';
      let component;

      if (/^http[s]?:\/\//i.test(item.component)) {
        iframeSrc = item.component;
        name = item.path;
        component = GxIframe;
        componentName = 'GxIframe';
      } else {
        if (o[name]) {
          name += Math.floor(Math.random() * 100000000); // name通常等于componentName才能keepalive相应刷新。若重复会导致冲突不能展示，这里尝试用随机数避免冲突
        }
        o[name] = 1;
        component = modules[`/src/views${item.component}.vue`];
      }

      const r = {
        path: item.path,
        name,
        meta: {
          title: item.title,
          componentName,
        },
        component,
      };
      iframeSrc && (r.meta.iframeSrc = iframeSrc);

      aryRouter.push(r);
    }
  });
}

// 页面菜单初始化成动态路由
export const initRouters = async () => {
  const menus = await getUserMenu(); // 用户菜单数据（菜单数据是响应式，会触发生成页面菜单）

  // 添加用户菜单对应的动态路由
  const routers = [];
  toRouter(routers, menus);
  router.addRoute({
    path: '/',
    name: 'Layout',
    component: () => import('~/components/layout/Layout.vue'),
    children: routers,
  });
};
