import { useTabsState, useTokenStore, useMenuState } from '~/pkgs';

// 是否需要登录
export const enableLogin = async () => {
  const rs = await $post('/v1/user/enableLogin', {});
  console.log(rs);

  const tokenStore = useTokenStore();
  if (rs.success) {
    if (rs.result) {
      tokenStore.needLogin = 'true';
    } else {
      tokenStore.needLogin = 'false';
    }
  } else {
    tokenStore.needLogin = 'unknow';
  }
  return rs;
};

// 用户登录
export const userLogin = async params => {
  const { username, password } = params;
  const rs = await $post('/v1/user/login', { username, password }, null, { 'Content-Type': 'application/x-www-form-urlencoded' });
  if (!rs.success) {
    return $msg.notify(rs.message, 'error');
  }

  const tokenStore = useTokenStore();
  tokenStore.token = rs.result; // 登录成功后设定令牌
  tokenStore.loginUserName = username;
  tokenStore.time = new Date().getTime();
  return rs.result;
};

// 用户登出
export const userLogout = () => {
  useTokenStore().$reset();
  useMenuState().$reset();
  useTabsState().$reset();
};

// 取用户菜单
export const getUserMenu = async () => {
  const menuState = useMenuState();
  const menus = [];
  menus.push({
    path: '/glc/search',
    title: '检索',
    icon: 'search',
    component: '/glc/search/GlcMain',
    openInner: true,
  });
  menus.push({
    path: '/glc/storages',
    title: '日志仓',
    icon: 'db',
    component: '/glc/storages/StoragesMain',
    openInner: true,
  });
  return (menuState.aryMenu = menus);
};
