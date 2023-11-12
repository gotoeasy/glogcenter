<template>
  <template v-for="item in menuState.aryMenu">
    <!-- 一级目录菜单 -->
    <el-sub-menu v-if="item.children" :key="item.id" :index="item.path">
      <template v-if="item.children" #title>
        <svg-icon :name="item.icon || 'item'" :class="themeStore.menuCollapse ? 'menu-collapse' : ''" />
        <span>{{ item.title }}</span>
      </template>
      <el-menu-item v-if="item.component && item.hidden != '1'" :index="item.path"><svg-icon name="item" />
        {{ item.title }}
      </el-menu-item>

    </el-sub-menu>

    <!-- 一级页面菜单 -->
    <el-menu-item v-if="item.openInner && (!item.role || tokenStore.role?.indexOf(item.role) >= 0)" :key="item.id"
      :index="item.path" :title="item.title" class="menu-collapse-page">
      <svg-icon :name="item.icon || 'item'" class="menu-page-icon" />
      <span>{{ item.title }}</span></el-menu-item>
    <el-menu-item v-if="item.openWindow && (!item.role || tokenStore.role?.indexOf(item.role) >= 0)" :key="item.id"
      class="menu-collapse-page" @click="fnClick(item)">
      <svg-icon :name="item.icon || 'item'" class="menu-page-icon" />
      <span>{{ item.title }}</span></el-menu-item>
  </template>
</template>

<script setup>
import { useMenuState, useThemeStore, useTokenStore } from '~/pkgs';

const themeStore = useThemeStore();
const menuState = useMenuState();
const tokenStore = useTokenStore();

// 点击页面菜单
const fnClick = (menu) => {
  let path = menu.children ? menu.children[0]?.path : menu.path;
  if (/^http[s]?:\/\//.test(path)) {
    path = path.replace(/\{token\}/ig, tokenStore.token);
    window.open(path, "_blank");
  } else {
    window.open(window.location.origin + path, "_blank");
  }
};

</script>

<style lang="scss">
.el-menu--collapse {
  .el-sub-menu__title {
    padding: 0;

    .menu-collapse {
      margin: 0 -10px 0 14px;
    }
  }
}

.menu-collapse-page {
  width: 44px; // 折叠时的菜单宽度，固定吧

  .menu-page-icon {
    width: 18px;
    height: 18px;
    margin: 0 -10px 0 11px;
  }
}

.el-popper.is-pure {
  // 折叠菜单行高。 （//TODO 实际这里的变量取不到，节点插入在body的缘故？）
  --el-menu-item-height: 36px;
  --el-menu-sub-item-height: 36px;

  // margin-top: -1px;

  .el-menu-item:hover {
    background-color: var(--el-menu-hover-bg-color); // 菜单折叠后的弹出菜单，鼠标悬停背景色。 （//TODO 实际这里的变量取不到，节点插入在body的缘故？）
  }
}

.el-menu .el-sub-menu .el-sub-menu__icon-arrow {
  margin-right: 4px; // 菜单目录右边的箭头
}
</style>
