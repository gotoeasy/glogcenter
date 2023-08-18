<template>
  <div>
    <router-view v-slot="{ Component, route }">
      <keep-alive :include="tabsState.keepAliveNames">
        <component :is="Component" v-if="showPage && !route.meta.iframeSrc" :key="route.path" />
      </keep-alive>
    </router-view>

    <div v-for="item in tabsState.iframeTabs" :key="item.path">
      <div v-if="item.meta?.iframeSrc">
        <GxIframe v-show="$route.meta?.iframeSrc == item.meta?.iframeSrc" :src="item.meta?.iframeSrc" />
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 【页面缓存设计效果】
 * 仅限于加入页签的页面才做keepalive处理
 * 页签相互切换时会保留状态，新开或关闭后再开的页签都是刷新显示的效果
 *
 * 【注意】
 * 1）是否缓存是通过组件名控制，组件名默认是首字母大写的驼峰式名称（勿轻易自定义修改）
 * 2）页面的默认组件名称作为路由的数据属性存放于route.meta.componentName，在路由初始化时自动设定
 * 3）路由变化时，动态修改keepAliveNames，达到控制缓存的目的
 * 4）页签打开的iframe通过v-show实现缓存效果
 */
import { $emitter, useTabsState } from '~/pkgs';

const router = useRouter();
const tabsState = useTabsState();

// 以下控制刷新页签
const showPage = ref(true);
$emitter.on("$PageMain:reflesh", tabPath => {
  if (!tabPath) {
    tabPath = tabsState.activePath;
  }
  const targetTab = tabsState.aryTabs[tabsState.aryTabs.findIndex(item => item.path === tabPath)];

  if (targetTab.meta?.iframeSrc) {
    const iframeIdx = tabsState.iframeTabs.findIndex(item => item.path === targetTab.path);
    // 取消显示
    tabsState.iframeTabs.splice(iframeIdx, 1);  // 删除指定页签缓存
    // 重新显示
    nextTick(() => {
      console.log("刷新iframe页签....", targetTab);
      tabsState.iframeTabs.push(targetTab); // 重新放回指定页签缓存
      if (tabsState.activePath !== targetTab.path) {
        router.push(targetTab.path); // 右击刷新别的页签
      }
    });
  } else {
    const nameIdx = tabsState.keepAliveNames.findIndex(item => item === targetTab.meta.componentName);
    // 取消显示
    tabsState.keepAliveNames.splice(nameIdx, 1);  // 删除指定页签缓存
    showPage.value = false; // 取消显示触发删除页面节点
    // 重新显示
    nextTick(() => {
      console.log("刷新页签....", tabPath);
      tabsState.keepAliveNames.push(targetTab.meta.componentName); // 重新放回指定页签缓存
      showPage.value = true; // 重新显示
      if (tabsState.activePath !== targetTab.path) {
        router.push(targetTab.path); // 右击刷新别的页签
      }
    });
  }
});
</script>
