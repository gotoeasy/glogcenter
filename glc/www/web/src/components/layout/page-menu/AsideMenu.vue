<template>
  <el-menu unique-opened class="aside-menu" :collapse="true" :background-color="themeStore.menuBgColor"
    :text-color="themeStore.menuColor" :active-text-color="themeStore.menuActiveColor" :default-active="$route.path"
    router>
    <MenuTree />
  </el-menu>
</template>

<script setup>
import { useThemeStore } from '~/pkgs';
import MenuTree from './MenuTree.vue';

const themeStore = useThemeStore();

</script>

<style >
.el-menu a {
  text-decoration: none;
}

.el-menu-item {
  width: 48px;
  height: 48px;
  color: #0081dd;
  background-color: white;
  border-bottom: solid 1px lightgray;
}

.side-menu-icon {
  margin-left: -9px;
}
</style>
