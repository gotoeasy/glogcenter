<template>
  <el-menu v-if="themeStore.layout === 'HeaderMainFooter'" class="header-menu" :background-color="themeStore.menuBgColor"
    :text-color="themeStore.menuColor" mode="horizontal" :active-text-color="themeStore.menuActiveColor"
    :default-active="$route.path" router>
    <MenuTree />
  </el-menu>
</template>

<script  setup>
import { useThemeStore } from '~/pkgs';
import MenuTree from './MenuTree.vue';

const themeStore = useThemeStore();

</script>

<style lang="scss">
.header-menu {
  // min-width: 300px;

  &.header-menu.el-menu--horizontal>.el-menu-item {
    border-bottom: 0;
  }

  &.header-menu.el-menu--horizontal>.el-sub-menu .el-sub-menu__title {
    border-bottom: 0;
  }
}
</style>
