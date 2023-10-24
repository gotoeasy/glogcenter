<template>
  <GxToolbar :height="themeStore.headerHeight">
    <template #left>
      <div style="display:flex;align-items:center;margin-left:8px;line-height: 30px;">
        <div :title="verInfo" style="width: 34px;color:white;text-align:center;cursor:pointer;" @click="clickLogo">
          <img src="/image/glc.png" style="width:34px;margin-top:9px;" />
        </div>
        <div style="text-align:center;">
          <img src="/image/title.png" style="height:32px;margin-top:11px;margin-left:26px" />
        </div>
      </div>
    </template>
    <template #right>
      <div v-if="needLogin" class="icon-item hover" style="margin-right:5px;" @click="logout">
        <span style="margin-right:-10px;margin-left:16px;">退出</span><svg-icon name="exit" size="16" />
      </div>
    </template>
  </GxToolbar>
</template>

<script setup>
import { gxUtil, useThemeStore, useTokenStore } from '~/pkgs';
import { userLogout } from '~/api';

const router = useRouter();
const verInfo = ref('');
const tokenStore = useTokenStore();
const themeStore = useThemeStore();
const headerHeight = computed(() => `${themeStore.headerHeight}px`);
const needLogin = computed(() => tokenStore.needLogin == 'true');
const headerColor = computed(() => {
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    themeStore.headerColor = '#606266';
  }
  return themeStore.headerColor
});
const headerActiveColor = computed(() => {
  if (!themeStore.customHeaderColor) {
    if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
      themeStore.headerActiveColor = '#606266';
    } else {
      themeStore.headerActiveColor = themeStore.headerColor;
    }
  }
  return themeStore.headerActiveColor
});
const headerActiveBgColor = computed(() => {
  if (!themeStore.customHeaderColor) {
    if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
      themeStore.headerActiveBgColor = gxUtil.darkColor(themeStore.headerBgColor, 0.1);
    } else {
      themeStore.headerActiveBgColor = gxUtil.lightColor(themeStore.headerBgColor);
    }
  }
  return themeStore.headerActiveBgColor;
});
const svgLogoColor = computed(() => {
  if (themeStore.headerBgColor.toLowerCase() === '#ffffff') {
    return '#0081dc';
  }
  return '#eeeeee';
});

onMounted(() => {
  checkVersion();
});

async function logout() {
  if (await $msg.confirm('确定要退出系统吗？')) {
    userLogout();
    router.push('/login');
  }
}

const clickLogo = () => {
  window.open('https://github.com/gotoeasy/glogcenter', '_blank');
};

function checkVersion() {
  if (!window.$checkVersionDone) {
    window.$checkVersionDone = true;
    // 从后台服务读取当前运行版本，避免多处维护版本号
    $post('/v1/version/info', {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
      if (rs.success) {
        verInfo.value = rs.result.version
        // 有新版本时，左上角图标鼠标悬停显示提示（注：最新版本号的查询服务并不保证随时可用）
        fetch(`https://glc.gotoeasy.top/glogcenter/current/version.json?v=${verInfo.value}`)
          .then(response => response.json())
          .then(data => {  // 最新版本（服务不保证可用，可能查不到，仅查到有新版本时更新tip）
            if (data.version && verInfo.value < data.version) {
              verInfo.value = `当前版本 ${verInfo.value} ，有新版本 ${data.version} 可更新`
            }
          })
          .catch(e => console.log(e));
      }
    });
  }
}

</script>

<style  lang="scss">
div.el-popper.el-dropdown__popper.el-popper {
  margin-top: -8px;
  margin-left: 8px;
}

.svg-logo {
  // height: v-bind('svgLogoHeight');
  fill: v-bind('svgLogoColor');
}
</style>

<style scoped lang="scss">
.icon-item {
  color: v-bind('headerColor');

}

.icon-item.hover {
  height: v-bind('headerHeight');
  line-height: v-bind('headerHeight');

  .svg-icon {
    margin: -2px 14px;
  }

  &:hover {
    color: v-bind('headerActiveColor');
    background-color: v-bind('headerActiveBgColor');

  }
}
</style>

