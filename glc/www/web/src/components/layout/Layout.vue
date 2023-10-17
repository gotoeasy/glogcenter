<template>
  <AppContainer>
    <el-container class="layout">
      <el-header
        v-show="!themeState.tabsMaximize && !themeState.pageMaximize && (themeStore.layout == 'HeaderAsideMainFooter' || themeStore.layout == 'HeaderMainFooter')"
        class="header">
        <PageHeader />
      </el-header>
      <el-container>
        <el-aside class="aside">
          <AsideMenu />
        </el-aside>
        <el-container>
          <el-header
            v-show="!themeState.tabsMaximize && !themeState.pageMaximize && (themeStore.layout == 'AsideHeaderMainFooter')"
            class="header">
            <PageHeader />
          </el-header>
          <el-main class="main">
            <!-- <PageTabs v-if="!themeState.pageMaximize && themeStore.showPageTabs" /> -->
            <div class="scroll">
              <div class="panel">
                <PageMain />
              </div>
              <div class="footer">
                <PageFooter v-show="!themeState.tabsMaximize && !themeState.pageMaximize && themeStore.showFooter" />
              </div>
            </div>
          </el-main>
        </el-container>
      </el-container>
    </el-container>
  </AppContainer>
</template>

<script setup>
import { useThemeStore, useThemeState } from '~/pkgs';
import AppContainer from '~/components/layout/app-container/AppContainer.vue';
import PageHeader from '~/components/layout/page-header/PageHeader.vue';
import AsideMenu from '~/components/layout/page-menu/AsideMenu.vue';
// import PageTabs from '~/components/layout/page-tabs/PageTabs.vue';
import PageFooter from '~/components/layout/page-footer/PageFooter.vue';
import PageMain from '~/components/layout/page-main/PageMain.vue';

const themeState = useThemeState();
const themeStore = useThemeStore();

const headerHeight = computed(() => `${((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight)}px`);
const mainHeight = computed(() => `calc(100vh - ${((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight)}px)`);
// const scrollHeight = computed(() => (!themeState.pageMaximize && themeStore.showPageTabs) ? `calc(100vh - ${((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight) + themeStore.pageTabsHeight}px)` : `calc(100vh - ${((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight)}px)`);
const scrollHeight = computed(() => `calc(100vh - ${((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight)}px)`);
const panelHeight = computed(() => {
  const header = ((themeState.pageMaximize || themeState.tabsMaximize) ? 0 : themeStore.headerHeight);
  // const tab = (!themeState.pageMaximize && themeStore.showPageTabs) ? themeStore.pageTabsHeight : 0;
  const tab = 0;
  const mainPadTop = themeStore.mainPanelMarginTop;
  const mainPadBottom = themeStore.mainPanelMarginBottom;
  // const footer = themeStore.showFooter ? themeStore.footerHeight : mainPadBottom;
  const footer = (!themeState.tabsMaximize && !themeState.pageMaximize && themeStore.showFooter) ? themeStore.footerHeight : mainPadBottom;
  const panelPadding = themeStore.mainPanelPaddingTop + themeStore.mainPanelPaddingBottom;
  return `calc(100vh - ${header + tab + mainPadTop + panelPadding + footer}px)`;
});
const footerHeight = computed(() => (`${(!themeState.tabsMaximize && !themeState.pageMaximize && themeStore.showFooter) ? themeStore.footerHeight : themeStore.mainPanelMarginBottom}px`));
const footerMarginTop = computed(() => (!themeState.tabsMaximize && !themeState.pageMaximize && themeStore.showFooter) ? `-${themeStore.mainPanelMarginBottom}px` : '0px');
const mainPanelPadding = computed(() =>
  `${themeStore.mainPanelPaddingTop}px ${themeStore.mainPanelPaddingRight}px ${themeStore.mainPanelPaddingBottom}px ${themeStore.mainPanelPaddingLeft}px`
);
const menuWidth = computed(() => `${((themeState.pageMaximize || themeState.tabsMaximize || themeStore.layout === 'HeaderMainFooter') ? 0 : themeStore.getMenuWidth)}px`);
const scrollPadding = computed(() => `${themeStore.mainPanelMarginTop}px ${themeStore.mainPanelMarginRight}px ${themeStore.mainPanelMarginBottom}px ${themeStore.mainPanelMarginLeft}px`);
const asideBorderTop = computed(() => {
  if (themeStore.layout === 'AsideHeaderMainFooter' || (themeStore.headerBgColor.toLowerCase() !== '#ffffff' || themeStore.menuBgColor.toLowerCase() !== '#ffffff')) {
    return 0;
  }
  let border = 0;
  let color = 'lightgray';
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    border = 1;
  } else {
    color = themeStore.menuBgColor;
  }
  return `${border}px solid ${color}`;
});
const headerBorderLeft = computed(() => {
  if (themeStore.layout === 'AsideHeaderMainFooter' && themeStore.headerBgColor.toLowerCase() === '#ffffff' && themeStore.menuBgColor.toLowerCase() === '#ffffff') {
    return `1px solid var(--el-border-color)`;
  }
  return 0;
});
const mainBorderLeft = computed(() => {
  if (themeStore.layout !== 'HeaderMainFooter' && themeStore.menuBgColor.toLowerCase() === '#ffffff') {
    return `1px solid var(--el-border-color)`;
  }
  return 0;
});
</script>

<style scoped lang="scss">
.layout {
  --gotoeasy-admin-footer-height: v-bind('footerHeight');

  .header {
    position: relative;
    min-width: 480px;
    height: v-bind('headerHeight');
    padding: 0;
    color: v-bind('themeStore.headerColor');
    background-color: v-bind('themeStore.headerBgColor');
    border-bottom: 1px solid lightgray;
    border-left: v-bind('headerBorderLeft');
  }

  .aside {
    width: v-bind('menuWidth');
    margin-top: -1px;
    margin-left: -1px;
    color: v-bind('themeStore.menuColor');
    background-color: v-bind('themeStore.menuBgColor');
    border-top: v-bind('asideBorderTop');

    // border-right: 1px solid lightgray;
  }

  .main {
    min-width: 500px;
    height: v-bind('mainHeight');
    padding: 0;
    margin: 0;
    overflow: hidden;
    border-left: v-bind('mainBorderLeft');

    .scroll {
      min-height: v-bind('scrollHeight');
      max-height: v-bind('scrollHeight');
      padding: 0;
      margin: 0;

      // padding: v-bind('scrollPadding');
      overflow: auto;
      background-color: #eee;

      .panel {
        min-height: v-bind('panelHeight');
        padding: v-bind('mainPanelPadding');
        margin: v-bind('scrollPadding');
        background-color: #fff;
      }

      .footer {
        margin-top: v-bind('footerMarginTop');
      }
    }
  }
}
</style>
