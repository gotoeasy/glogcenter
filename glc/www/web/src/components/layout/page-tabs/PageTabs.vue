<template>
  <GxToolbar :height="themeStore.pageTabsHeight" class="page-tabs" @selectstart="e => e.preventDefault()">
    <template #left>
      <el-tabs v-model="tabsState.activePath" class="page-tabs_left" :class="themeStore.pageTabsBg ? 'has-bg-color' : ''"
        @tab-click="tabClick" @tab-remove="tabRemove">
        <el-tab-pane v-for="(item, idx) in tabsState.aryTabs" :key="item.path" :name="item.path" :closable="!!idx">
          <!-- 右键菜单开始：自定义标签页显示名称，保证每个标签页都能实现右键菜单 -->
          <template #label>
            <el-dropdown :id="item.path" ref="dropdownRef" size="default" trigger="contextmenu"
              @selectstart="e => e.preventDefault()" @visible-change="handleChange($event, item)">
              <span @selectstart="e => e.preventDefault()">{{ item.title || item.meta?.title }}</span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="refreshTab(item.path)">
                    <svg-icon name="reload" />刷新当前标签页
                  </el-dropdown-item>
                  <el-divider class="divider" />
                  <el-dropdown-item :disabled="tabsState.isDisable(item.path, 'self')"
                    @click="tabRemove(item.path, 'self')">
                    <svg-icon name="close" />关闭当前标签页
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="tabsState.isDisable(item.path, 'left')"
                    @click="tabRemove(item.path, 'left')">
                    <svg-icon name="to-left" />关闭左侧标签页
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="tabsState.isDisable(item.path, 'right')"
                    @click="tabRemove(item.path, 'right')">
                    <svg-icon name="to-right" />关闭右侧标签页
                  </el-dropdown-item>
                  <el-divider class="divider" />
                  <el-dropdown-item :disabled="tabsState.isDisable(item.path, 'other')"
                    @click="tabRemove(item.path, 'other')">
                    <svg-icon name="to-other" />关闭其他标签页
                  </el-dropdown-item>
                  <el-dropdown-item :disabled="tabsState.isDisable(item.path, 'all')"
                    @click="tabRemove(item.path, 'all')">
                    <svg-icon name="to-all" />关闭全部标签页
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <!-- 右键菜单结束 -->
          <component :is="item.componentName" v-bind="item.data" />
        </el-tab-pane>
      </el-tabs>
    </template>
    <template #right>
      <div class="page-tabs_right">
        <el-tooltip content="刷新当前页" placement="bottom">
          <span class="item-icon hover" @click="refreshTab()">
            <svg-icon name="reload" size="16" class="right-icon" />
          </span>
        </el-tooltip>
        <el-tooltip content="缩放页面" :visible="showTooltip" placement="bottom">
          <span v-if="themeState.tabsMaximize" class="item-icon hover" @mouseenter="showTooltip = true"
            @mouseleave="showTooltip = false" @click="switchMaximize">
            <svg-icon name="minimize" size="16" class="right-icon" />
          </span>
          <span v-else class="item-icon hover" @mouseenter="showTooltip = true" @mouseleave="showTooltip = false"
            @click="switchMaximize">
            <svg-icon name="maximize" size="16" class="right-icon" />
          </span>
        </el-tooltip>
      </div>
    </template>
  </GxToolbar>
</template>

<script setup>
import { useThemeStore, useThemeState, useTabsState } from '~/pkgs';

const router = useRouter();
const themeState = useThemeState();
const themeStore = useThemeStore();
const tabsState = useTabsState();

const showTooltip = ref(false); // 默认切换缩放页签时会有不能关闭的现象，故特殊处理单独控制

// 点击页签，显示相应的路由页面
const tabClick = tab => (tabsState.activePath !== tab.paneName) && router.push(tab.paneName);

// 刷新页签
const refreshTab = tabPath => $emitter.emit("$PageMain:reflesh", tabPath);

// 关闭页签，显示相应的路由页面
const tabRemove = (path, type) => {
  tabsState.closeTab(path, type);
  router.push(tabsState.activePath);
}

// 右击弹出菜单时，控制关闭其他已弹出的右键菜单
const dropdownRef = ref()
const handleChange = (event, tabItem) => {
  event && dropdownRef.value.forEach((item) => (item.id !== tabItem.path && item.handleClose())); // 不是正弹出的右键菜单就都关掉
}

// 风格相关响应式控制
const pageTabsHeight = computed(() => `${themeStore.pageTabsHeight}px`);
const pageTabsHeightM4 = computed(() => `${themeStore.pageTabsHeight - 4}px`);
const rightIconMargin = computed(() => {
  const half = parseInt((themeStore.pageTabsHeight - 16) / 2, 10);
  return `${half}px 10px`;
});
const tabsWidth = computed(() => `calc(100vw - ${((themeState.tabsMaximize || themeStore.layout === 'HeaderMainFooter') ? 0 : themeStore.getMenuWidth) + 72}px)`); // 刷新图标+最大化图标=72px
const tabsBorderTop = computed(() => {
  let border = 0;
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    border = 1;
  }
  return `${border}px solid lightgray`;
});
const tabsBorderLeft = computed(() => {
  if (themeStore.layout !== 'HeaderMainFooter' && themeStore.menuBgColor.toLowerCase() === '#ffffff') {
    return `1px solid var(--el-border-color)`;
  }
  return 0;
});

// 放大缩小切换显示主面板（含页签）
function switchMaximize() {
  themeState.tabsMaximize = !themeState.tabsMaximize;
  showTooltip.value = false;
}
</script>

<style lang="scss" scoped>
.right-icon {
  margin: v-bind('rightIconMargin');
}

.divider.el-divider--horizontal {
  margin: 2px 0;
}
</style>

<style lang="scss">
.c-gx-toolbar.page-tabs {
  border-top: v-bind('tabsBorderTop');
  border-left: v-bind('tabsBorderLeft');

  .page-tabs_left {
    width: v-bind('tabsWidth');
    min-width: 300px;
  }

  .page-tabs_right {
    display: flex;
    height: v-bind('pageTabsHeight');

    .item-icon {
      height: v-bind('pageTabsHeight');
      margin-top: -1px;
      line-height: v-bind('pageTabsHeight');
      border-left: 1px solid var(--el-border-color);
      box-shadow: 0 2px 0 var(--el-border-color-light); // 右图标加下边框
    }
  }

  .page-tabs_right_wrap {
    display: flex;
  }

  // 仅显示单行页签头部
  .el-tabs__header {
    padding: 0;
    margin: 0 0 0 10px; // 左边间隔10px
  }

  // 标题垂直居中
  .el-dropdown {
    padding: 0 12px;
    line-height: v-bind('pageTabsHeightM4'); // 调整高度对齐关闭图标，小4px
  }

  // 页签前翻后翻的箭头按钮宽度
  .el-tabs__nav-prev {
    min-width: 26px; // 左边缩小10px
    height: var(--el-tabs-header-height);
    line-height: var(--el-tabs-header-height);
    border-right: 1px solid var(--el-border-color-light);
  }

  .el-tabs__nav-next {
    min-width: 36px;
    height: var(--el-tabs-header-height);
    line-height: var(--el-tabs-header-height);
    border-left: 1px solid var(--el-border-color-light);
  }

  .el-tabs__nav-wrap.is-scrollable {
    padding: 0 36px;
  }

  .el-tabs__item {
    padding: 0;

    .is-icon-close {
      padding-bottom: 2px; // 微调页签的关闭悬停背影
      margin-left: -10px;
      visibility: hidden; // 页签默认不显示关闭图标
    }
  }

  // 页签鼠标悬停时显示关闭图标
  .el-tabs__item:hover {
    .is-icon-close {
      visibility: visible;
    }
  }

  // 页签激活时背景色用最淡主题色
  .has-bg-color {
    .el-tabs__item.is-active {
      background-color: var(--el-color-primary-light-9);
    }
  }

  // // 焦点页签的下边框线调粗
  // .el-tabs__active-bar {
  //   height: 3px;
  // }
}

// 页签右键菜单的图标位置微调
.el-dropdown-menu__item {
  .svg-icon {
    margin: 0 3px 0 0;
  }
}
</style>
