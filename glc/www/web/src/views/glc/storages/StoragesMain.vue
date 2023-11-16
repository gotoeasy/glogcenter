<template>
  <div v-show="visible">
    <el-row style="align-items: center;height:26px;">
      <SvgIcon name="detail" height="16" width="16" style="margin: 0 6px 0 0;color:var(--el-color-primary)" />
      <span style="font-size:16px">日志仓信息列表</span>
    </el-row>

    <el-divider style="margin: 0 0 8px" />

    <GxToolbar style="margin-bottom: 8px" class="c-btn">
      <template #left>
        <GxButton icon="refresh-right" @click="search">刷 新</GxButton>
      </template>

      <template #right>
        <el-tooltip content="缩放" placement="top">
          <el-button circle @click="emitter.emit('main:switchMaximizePage')">
            <SvgIcon name="zoom" />
          </el-button>
        </el-tooltip>
        <GxPageTableConfig :tid="tid" :page-config="pageSettingStore" />
      </template>
    </GxToolbar>

    <GxTable ref="table" v-loading="showTableLoadding" scrollbar-always-on stripe :enable-header-contextmenu="false"
      :tid="tid" :data="tableData" :max-height="tableHeight" :height="tableHeight" class="c-gx-table c-glc-table"
      row-key="id">
      <template #$operation="{ row }">
        <el-button size="small" type="warning" @click="remove(row)">删除</el-button>
      </template>
    </GxTable>

    <div>
      <div style="display:flex;justify-content:space-between;">
        <div style="min-height:30px;padding-top:5px; line-height:30px; color: #909399;" v-text="info"></div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { userLogout } from "~/api";
import { useEmitter, usePageMainHooks, useTabsState } from "~/pkgs";

const tabsState = useTabsState();
const emitter = useEmitter(tabsState.activePath);
const router = useRouter();

const opt = {
  emitter,
  withoutSearchForm: true,
};
const { visible, tableData, getTableHeight, pageSettingStore, showTableLoadding } = usePageMainHooks(opt);

const table = ref(); // 表格实例
const tid = ref('storagesMain231126'); // 表格ID
const info = ref(''); // 底部提示信息

const tableHeight = computed(() => getTableHeight(false, true, 19)); // 表格高度

// 初期默认检索
onMounted(() => {
  const configStore = $emitter.emit('$table:config', { id: tid.value });
  !configStore.columns.length && $emitter.emit('$table:config', { id: tid.value, update: true }); // 首次使用开启默认布局
  search()
});

function search() {
  // 日志仓列表查询
  showTableLoadding.value = true;
  const url = `/v1/store/list`;
  $post(url, {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const data = rs.result.data || [];
      tableData.value.splice(0, tableData.value.length, ...data);
      info.value = rs.result.info;
      document.querySelector('.c-glc-table .el-scrollbar__wrap').scrollTop = 0; // 滚动到顶部
    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
      router.push('/login');
    }
  }).finally(() => {
    showTableLoadding.value = false;
  });
}

async function remove(row) {
  // 日志仓删除
  if (await $msg.confirm(`确定要删除日志仓 ${row.name} 吗？`)) {
    const url = `/v1/store/delete`;
    $post(url, { storeName: row.name }, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
      console.log(rs)
      if (rs.success) {
        $msg.info(`已删除日志仓 ${row.name}`);
        search();
      } else if (rs.code == 403) {
        userLogout(); // 403 时登出
        router.push('/login');
      } else {
        $msg.error(rs.message);
      }
    }).finally(() => {
      showTableLoadding.value = false;
    });
  }
}

</script>

<style>
.c-glc-table .el-popper.is-dark {
  display: none;
}

.x-detail {
  padding: 5px 5px 5px 30px;
  background-color: floralwhite;
}

.c-search-form .c-search-form-item {
  margin-bottom: 0;
}

.c-search-form .c-btn-badge.el-badge .el-button--small.is-circle {
  /* width: var(--el-button-size);
  min-width: var(--el-button-size); */
  width: 30px;
  min-width: 30px;
  height: 30px;
  min-height: 30px;
  margin-top: -3px;
}

.c-search-form.el-form--inline .el-input {
  --el-input-width: 100%;
}

.c-datapicker.el-popper.is-pure {
  margin-left: -100px;
}
</style>
