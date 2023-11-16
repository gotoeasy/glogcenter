<template>
  <div class="gotoeasy-admin">
    <slot></slot>
    <ThemeDrawer />
    <GxTableConfig :width="1150" />
    <GxEditTableConfig :width="1150" />
  </div>
</template>

<script setup>
import { gxUtil, useThemeStore, useThemeState } from '~/pkgs';
import ThemeDrawer from '~/components/theme-drawer/ThemeDrawer.vue';

const themeStore = useThemeStore();
const themeState = useThemeState();

// 主题色
const primaryColor = computed(() => themeStore.primaryColor);
const lightPrimaryColor1 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 1 / 10))
const lightPrimaryColor2 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 2 / 10))
const lightPrimaryColor3 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 3 / 10))
const lightPrimaryColor4 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 4 / 10))
const lightPrimaryColor5 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 5 / 10))
const lightPrimaryColor6 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 6 / 10))
const lightPrimaryColor7 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 7 / 10))
const lightPrimaryColor8 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 8 / 10))
const lightPrimaryColor9 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 9 / 10))
const darkPrimaryColor1 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 1 / 10))
const darkPrimaryColor2 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 2 / 10))
const darkPrimaryColor3 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 3 / 10))
const darkPrimaryColor4 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 4 / 10))
const darkPrimaryColor5 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 5 / 10))
const darkPrimaryColor6 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 6 / 10))
const darkPrimaryColor7 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 7 / 10))
const darkPrimaryColor8 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 8 / 10))
const darkPrimaryColor9 = computed(() => gxUtil.mixColor(themeStore.primaryColor, "#000000", 9 / 10))

// // 菜单
const menuHeight = computed(() => `${themeStore.menuHeight}px`);
const menuBgColor = computed(() => themeStore.menuBgColor);

// 菜单
const menuActiveBgColor = computed(() => {
  if (themeStore.customMenuColor) {
    return themeStore.menuActiveBgColor;
  }
  if (themeStore.menuBgColor.value === themeStore.primaryColor.value) {
    return gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 3 / 10);
  }
  return themeStore.primaryColor;
});
const menuHoverBgColor = computed(() => {
  if (themeStore.customMenuColor) {
    return themeStore.menuHoverBgColor;
  }
  if (themeStore.menuBgColor.value === themeStore.primaryColor.value) {
    return gxUtil.mixColor(themeStore.primaryColor, "#ffffff", 3 / 10);
  }
  return themeStore.primaryColor;
});

const menuColor = computed(() => themeStore.menuColor);
const menuActiveColor = computed(() => themeStore.menuActiveColor);

// Tabs
const pageTabsHeight = computed(() => `${themeStore.pageTabsHeight}px`);

$emitter.on('$layout:SwitchMaximizePage', () => {
  themeState.pageMaximize = !themeState.pageMaximize;
});

// ------------ iframe ------------
const iframeHeight = computed(() => {
  const header = ((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight);
  // const tab = (!themeState.pageMaximize && themeStore.showPageTabs) ? themeStore.pageTabsHeight : 0;
  const tab = 0;
  const mainPadTop = themeStore.mainPanelMarginTop;
  const mainPadBottom = themeStore.mainPanelMarginBottom;
  const footer = (!themeState.tabsMaximize && !themeState.pageMaximize && themeStore.showFooter) ? themeStore.footerHeight : mainPadBottom;
  const panelPadding = themeStore.mainPanelPaddingTop + themeStore.mainPanelPaddingBottom;
  return `calc(100vh - ${header + tab + mainPadTop + panelPadding + footer + 12}px)`;
});

const iframeWidth = computed(() => {
  let menu = 0;
  if (themeStore.layout !== 'HeaderMainFooter' && !themeState.tabsMaximize && !themeState.pageMaximize) {
    menu = themeStore.menuCollapse ? themeStore.menuCollapseWidth : themeStore.menuExpandWidth;
  }
  const margin = themeStore.mainPanelMarginLeft + themeStore.mainPanelMarginRight;
  const padding = themeStore.mainPanelPaddingLeft + themeStore.mainPanelPaddingRight;
  return `calc(100vw - ${menu + margin + padding + 8}px)`
});
</script>

<style lang="scss">
:root {
  --el-menu-hover-bg-color: v-bind('menuActiveBgColor');
}

.gotoeasy-admin {
  // 主题色
  --el-color-primary: v-bind('primaryColor');
  --el-color-primary-light-1: v-bind('lightPrimaryColor1');
  --el-color-primary-light-2: v-bind('lightPrimaryColor2');
  --el-color-primary-light-3: v-bind('lightPrimaryColor3');
  --el-color-primary-light-4: v-bind('lightPrimaryColor4');
  --el-color-primary-light-5: v-bind('lightPrimaryColor5');
  --el-color-primary-light-6: v-bind('lightPrimaryColor6');
  --el-color-primary-light-7: v-bind('lightPrimaryColor7');
  --el-color-primary-light-8: v-bind('lightPrimaryColor8');
  --el-color-primary-light-9: v-bind('lightPrimaryColor9');
  --el-color-primary-dark-1: v-bind('darkPrimaryColor1');
  --el-color-primary-dark-2: v-bind('darkPrimaryColor2');
  --el-color-primary-dark-3: v-bind('darkPrimaryColor3');
  --el-color-primary-dark-4: v-bind('darkPrimaryColor4');
  --el-color-primary-dark-5: v-bind('darkPrimaryColor5');
  --el-color-primary-dark-6: v-bind('darkPrimaryColor6');
  --el-color-primary-dark-7: v-bind('darkPrimaryColor7');
  --el-color-primary-dark-8: v-bind('darkPrimaryColor8');
  --el-color-primary-dark-9: v-bind('darkPrimaryColor9');

  // 菜单
  --el-menu-active-color: v-bind('menuActiveColor');
  --el-menu-text-color: v-bind('menuColor');
  --el-menu-hover-text-color: v-bind('menuActiveColor');
  --el-menu-bg-color: v-bind('menuBgColor');
  --el-menu-hover-bg-color: v-bind('menuHoverBgColor');
  --el-menu-item-height: v-bind('menuHeight');
  --el-menu-sub-item-height: var(--el-menu-item-height); // --el-menu-sub-item-height: calc(var(--el-menu-item-height) - 4px);
  --el-menu-item-hover-fill: v-bind('menuActiveColor');
  --el-menu-base-level-padding: 4px;
  --el-menu-level-padding: 16px;

  // 组件高度（如输入框等）
  --el-component-size: 28px;

  // 组件高度（如输入框左边的标签文本等）
  .el-form-item--default .el-form-item__label {
    height: 28px;
    line-height: 28px;
  }

  // Tabs
  .page-tabs {
    .el-tabs {
      --el-tabs-header-height: v-bind('pageTabsHeight');
    }
  }

  .el-menu {
    border-right: 0;

    .el-menu-item {
      &.is-active {
        color: var(--el-menu-active-color);
        background-color: v-bind('menuActiveBgColor');
      }
    }

    .el-menu-item:hover,
    .el-sub-menu__title:hover {
      color: var(--el-menu-active-color);
      background-color: v-bind('menuHoverBgColor');
    }
  }

}
</style>

<style lang="scss">
// 检索一览页面
.c-pagging {
  margin-top: 10px;

  &.el-pagination {
    --el-pagination-button-height: 22px;
  }

  .el-input__inner {
    height: 22px;
    line-height: 22px;
  }

  .el-input__wrapper {
    padding: 0 11px;
  }

  .el-pagination .el-select .el-input {
    width: 100px;
  }

  .el-pagination__editor.el-input {
    width: 50px;
  }
}

// // 弹出的详细查询条件面板
// .c-down-panel {
//   position: absolute;
//   z-index: 100;
//   width: v-bind('downPanelWidth');
//   padding: v-bind('downPanelPadding');
//   margin-left: v-bind('downPanelMarginLeft');
//   background-color: white;
//   box-shadow: 0 2px 12px 0 rgb(0 0 0 / 30%);
// }

.c-popover-table .el-input-number.is-controls-right .el-input__wrapper {
  padding-right: 32px;
  padding-left: 5px;
}

.c-align .el-button--small {
  padding: 5px;
}

// 表格操作按钮微调
.gotoeasy-admin .el-button.is-link {
  padding: 4px 2px;
}

// iframe高宽
.c-iframe {
  width: v-bind('iframeWidth');
  height: v-bind('iframeHeight');
}
</style>

