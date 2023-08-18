<template>
  <el-drawer v-model="visible" class="theme-drawer" title="布局设置" size="360px" append-to-body>
    <!-- <div class="theme-item" style="margin-top: 0;">
      <div>
        <svg-icon v-if="themeStore.layout !== 'AsideHeaderMainFooter'" name="layout-ahm" size="90" style="cursor:hand;"
          @click="selectLayout('AsideHeaderMainFooter')" />
        <svg-icon v-if="themeStore.layout == 'AsideHeaderMainFooter'" name="layout-ahm-c" size="90" />
      </div>
      <div>
        <svg-icon v-if="themeStore.layout !== 'HeaderAsideMainFooter'" name="layout-ham" size="90" style="cursor:hand;"
          @click="selectLayout('HeaderAsideMainFooter')" />
        <svg-icon v-if="themeStore.layout == 'HeaderAsideMainFooter'" name="layout-ham-c" size="90" />
      </div>
      <div>
        <svg-icon v-if="themeStore.layout !== 'HeaderMainFooter'" name="layout-hm" size="90" style="cursor:hand;"
          @click="selectLayout('HeaderMainFooter')" />
        <svg-icon v-if="themeStore.layout == 'HeaderMainFooter'" name="layout-hm-c" size="90" />
      </div>
    </div> -->
    <div class="theme-item" style="padding: 0 16px;margin-bottom: 10px;">
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[0] }" @click="selectPrimaryColor(0)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[1] }" @click="selectPrimaryColor(1)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[2] }" @click="selectPrimaryColor(2)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[3] }" @click="selectPrimaryColor(3)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[4] }" @click="selectPrimaryColor(4)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[5] }" @click="selectPrimaryColor(5)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[6] }" @click="selectPrimaryColor(6)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[7] }" @click="selectPrimaryColor(7)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[8] }" @click="selectPrimaryColor(8)">
        </div>
      </el-tooltip>
      <el-tooltip content="主题名称" placement="top">
        <div class="theme" :style="{ backgroundColor: config.primaryColors[9] }" @click="selectPrimaryColor(9)">
        </div>
      </el-tooltip>
    </div>
    <div style="padding: 0 16px;margin-bottom: 20px;">
      <div class="theme-item">
        <span class="demonstration">自定义主题色</span>
        <el-color-picker v-model="themeStore.primaryColor" popper-class="theme-color-picker"
          :predefine="config.themeColors" />
      </div>
    </div>

    <el-collapse v-model="activeNames" accordion>
      <!-- <el-collapse-item title="主题设定" name="1">
        <div class="theme-item" style="padding: 0 16px">
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[0] }" @click="selectPrimaryColor(0)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[1] }" @click="selectPrimaryColor(1)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[2] }" @click="selectPrimaryColor(2)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[3] }" @click="selectPrimaryColor(3)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[4] }" @click="selectPrimaryColor(4)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[5] }" @click="selectPrimaryColor(5)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[6] }" @click="selectPrimaryColor(6)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[7] }" @click="selectPrimaryColor(7)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[8] }" @click="selectPrimaryColor(8)">
            </div>
          </el-tooltip>
          <el-tooltip content="主题名称" placement="top">
            <div class="theme" :style="{ backgroundColor: config.primaryColors[9] }" @click="selectPrimaryColor(9)">
            </div>
          </el-tooltip>
        </div>
        <div style="padding: 0 16px">
          <div class="theme-item">
            <span class="demonstration">自定义主题色</span>
            <el-color-picker v-model="themeStore.primaryColor" popper-class="theme-color-picker"
              :predefine="config.themeColors" />
          </div>
        </div>
      </el-collapse-item> -->

      <el-collapse-item title="主面板外边距调整">
        <div style="padding-left: 16px;">
          <div class="theme-item">
            <span>外边距-上(px)</span>
            <el-input-number v-model="themeStore.mainPanelMarginTop" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>外边距-下(px)</span>
            <el-input-number v-model="themeStore.mainPanelMarginBottom" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>外边距-左(px)</span>
            <el-input-number v-model="themeStore.mainPanelMarginLeft" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>外边距-右(px)</span>
            <el-input-number v-model="themeStore.mainPanelMarginRight" :min="0" :max="100" />
          </div>
        </div>
      </el-collapse-item>
      <el-collapse-item title="主面板内边距调整">
        <div style="padding-left: 16px;">
          <div class="theme-item">
            <span>内边距-上(px)</span>
            <el-input-number v-model="themeStore.mainPanelPaddingTop" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>内边距-下(px)</span>
            <el-input-number v-model="themeStore.mainPanelPaddingBottom" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>内边距-左(px)</span>
            <el-input-number v-model="themeStore.mainPanelPaddingLeft" :min="0" :max="100" />
          </div>
          <div class="theme-item">
            <span>内边距-右(px)</span>
            <el-input-number v-model="themeStore.mainPanelPaddingRight" :min="0" :max="100" />
          </div>
        </div>
      </el-collapse-item>

      <el-collapse-item title="页眉调整">
        <div style="padding-left: 16px;">
          <div class="theme-item">
            <span class="demonstration">高度(px)</span>
            <el-input-number v-model="themeStore.headerHeight" :min="40" :max="100" />
          </div>
          <div class="theme-item">
            <span class="demonstration">文字颜色</span>
            <el-color-picker v-model="themeStore.headerColor" popper-class="theme-color-picker"
              :predefine="config.textColors" />
          </div>
          <div class="theme-item">
            <span class="demonstration">背景颜色</span>
            <el-color-picker v-model="themeStore.headerBgColor" popper-class="theme-color-picker"
              :predefine="config.themeColors" />
          </div>
          <div class="theme-item">
            <span>自定义活动颜色</span>
            <el-switch v-model="themeStore.customHeaderColor" />
          </div>
          <div class="theme-item">
            <span class="demonstration">活动文字颜色</span>
            <el-color-picker v-model="themeStore.headerActiveColor" :disabled="!themeStore.customHeaderColor"
              popper-class="theme-color-picker" :predefine="config.themeColors" />
          </div>
          <div class="theme-item">
            <span class="demonstration">活动背景颜色</span>
            <el-color-picker v-model="themeStore.headerActiveBgColor" :disabled="!themeStore.customHeaderColor"
              popper-class="theme-color-picker" :predefine="config.themeColors" />
          </div>
        </div>
      </el-collapse-item>

      <el-collapse-item title="页脚调整">
        <div style="padding-left: 16px;">
          <div class="theme-item">
            <span>显示页脚</span>
            <el-switch v-model="themeStore.showFooter" />
          </div>
          <div class="theme-item">
            <span class="demonstration">高度(px)</span>
            <el-input-number v-model="themeStore.footerHeight" :disabled="!themeStore.showFooter" :min="10" :max="100" />
          </div>
          <div class="theme-item">
            <span class="demonstration">文字颜色</span>
            <el-color-picker v-model="themeStore.footerColor" popper-class="theme-color-picker"
              :predefine="config.textColors" />
          </div>
          <div class="theme-item">
            <span class="demonstration">背景颜色</span>
            <el-color-picker v-model="themeStore.footerBgColor" popper-class="theme-color-picker"
              :predefine="config.themeColors" />
          </div>
        </div>
      </el-collapse-item>
    </el-collapse>

    <template #footer>
      <el-row>
        <GxButton icon="select" type="primary" @click="fnSaveMyLayout">保存为我的布局</GxButton>
        <GxButton icon="refresh-right" @click="fnUseMyLayout">恢复为我的布局</GxButton>
      </el-row>
    </template>
  </el-drawer>
</template>

<script setup>
import { useThemeStore, useTokenStore } from '~/pkgs';
import { config } from '~/config';

const activeNames = '1';

const visible = ref(false);
const themeStore = useThemeStore();
const tokenStore = useTokenStore();

// 布局切换
// const selectLayout = name => (themeStore.layout = name);

// 主题色切换
const selectPrimaryColor = idx => (themeStore.primaryColor = config.primaryColors[idx]);

// 通过事件控制窗口的显示
$emitter.on("openThemeDrawer", () => (visible.value = true));

const fnSaveMyLayout = () => {
  const id = `theme-${tokenStore.loginUserName}`;
  const oState = { ...themeStore.$state }
  const content = JSON.stringify(oState);
  $emitter.emit('$axios:post', '/system/config/save', { id, content }).then(rs => {
    if (rs.success) {
      $msg.info('操作成功');
    }
  })
}

const fnUseMyLayout = () => {
  const id = `theme-${tokenStore.loginUserName}`;
  $emitter.emit('$axios:post', '/system/config/get', { id })?.then(rs => {
    if (rs.success) {
      if (rs.data) {
        const state = JSON.parse(rs.data.content);
        themeStore.$patch(state);
        $msg.info('操作成功');
      } else {
        console.log('没有个人布局，忽略操作', id);
      }
    } else {
      console.log('查取个人布局失败', id);
    }
  });
}

</script>

<style lang="css">
/* 不显示颜色选择器的清空按钮 */
.theme-color-picker .el-button.is-text {
  display: none;
}

/* 风琴行高 */
.theme-drawer .el-collapse {
  --el-collapse-header-height: 40px;
}
</style>

<style scoped lang="css">
/* 布局图标选择 */
.theme {
  width: 20px;
  height: 20px;
  cursor: pointer;
  border: 1px solid var(--el-border-color-lighter);
}
</style>

<style scoped lang="scss">
@import './ThemeDrawer';
</style>
