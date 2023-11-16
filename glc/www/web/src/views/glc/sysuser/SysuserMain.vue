<template>
  <div v-show="visible">
    <el-row style="align-items: center;height:26px;">
      <SvgIcon name="detail" height="16" width="16" style="margin: 0 6px 0 0;color:var(--el-color-primary)" />
      <span style="font-size:16px">用户信息列表</span>
    </el-row>

    <el-divider style="margin: 0 0 8px" />

    <GxToolbar style="margin-bottom: 8px" class="c-btn">
      <template #left>
        <GxButton icon="refresh-right" @click="emitter.emit('main:search', 1)">刷 新</GxButton>
        <GxButton icon="new" @click="emitter.emit('main:create')">新 建</GxButton>
      </template>

      <template #right>
        <el-tooltip content="缩放页面" placement="top">
          <el-button circle @click="emitter.emit('main:switchMaximizePage')">
            <SvgIcon name="zoom" />
          </el-button>
        </el-tooltip>
        <GxPageTableConfig :tid="tid" :page-config="pageSettingStore" />
      </template>
    </GxToolbar>

    <GxTable ref="table" v-loading="showTableLoadding" :tid="tid" :data="tableData" :max-height="tableHeight"
      :height="tableHeight">
      <template #status="{ row }">
        <SvgIcon v-if="row.status == '1'" name="success-filled" style="padding-top: 5px;color:limegreen;" />
        <SvgIcon v-else name="close-bold" style="padding-top: 5px;" />
      </template>
      <template #$operation="{ row }">
        <el-button size="small" @click="emitter.emit('main:edit', row)">编辑</el-button>
        <el-button size="small" type="warning" @click="() => fnDelete(row)">删除</el-button>
      </template>
    </GxTable>

  </div>
  <PageDetail />
</template>

<script setup>
import { useEmitter, usePageMainHooks, useTabsState } from "~/pkgs";
import usePageSettingStore from "./SysuserPageSetting.store";
import PageDetail from "./SysuserPageDetail.vue";

const tabsState = useTabsState();
const emitter = useEmitter(tabsState.activePath);

const opt = { emitter, pageSettingStore: usePageSettingStore() };
const { visible, tableData, getTableHeight, pageSettingStore, showTableLoadding } = usePageMainHooks(opt);

const table = ref(); // 表格实例
const tid = ref('userMain231126'); // 表格ID

const tableHeight = computed(() => getTableHeight(false, true, 19)); // 表格高度

// 挂载事件
onMounted(() => emitter.emit('main:search'));

// 删除
const fnDelete = async (row) => {
  const confirmMsg = `确定删除用户 ${row.username} ？`;
  if (await $msg.confirm(confirmMsg)) {
    $emitter.emit('$axios:post', pageSettingStore.pageOptions.urlDelete, { username: row.username }).then(rs => {
      if (!rs.success) {
        $msg.error(rs.message);
      } else {
        emitter.emit('main:search');
        $msg.info('操作成功！');
      }
    });
  }
};
</script>
