<template>
  <!-- 弹窗详情展示模式 -->
  <GxDialog ref="dialog" :title="title" width="700" height="620" show-footer
    :before-close="() => emitter.emit('detail:close')">
    <template #default>
      <div style="padding:0 60px;">
        <DetailForm ref="formDialog" :data="formData" :is-create-mode="isCreateMode" :readonly="readonly" />
      </div>
    </template>
    <template #footer>
      <el-row style="justify-content: center;">
        <GxButton style="margin-right:20px;margin-left:20px;" @click="emitter.emit('detail:close')">关 闭</GxButton>
        <GxButton v-if="!readonly" icon="select" type="primary" style="margin-right:20px;margin-left:20px;"
          @click="save(formDialog)">保 存</GxButton>
      </el-row>
    </template>
  </GxDialog>
</template>

<script setup>
import { useEmitter, usePageDetailHooks, useTabsState } from "~/pkgs";
import usePageSettingStore from "./SysuserPageSetting.store";
import DetailForm from "./SysuserPageDetailForm.vue";

const tabsState = useTabsState();
const emitter = useEmitter(tabsState.activePath);
const pageSettingStore = usePageSettingStore();
const formDialog = ref();  // 表单实例(弹框)
const dialog = ref();  //  弹框

const opt = { emitter, ...pageSettingStore.pageOptions };
const { formData, title, readonly, isCreateMode, visible } = usePageDetailHooks(opt);

// 监听visible变量，判断控制对话框详情页面的显示和关闭
watch(visible, () => pageSettingStore.detailMode == 'dialog' && dialog.value.show(visible.value));

// 确定保存
const save = (form) => {
  const data = { ...formData.value };
  data.password = data.password1;
  delete data.password1;
  delete data.password2;
  emitter.emit('detail:save', form, data);
};

</script>

