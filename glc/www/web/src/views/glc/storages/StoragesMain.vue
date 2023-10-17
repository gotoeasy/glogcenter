<template>
  <div v-show="visible">

    <GxToolbar style="margin-bottom: 13px" class="c-btn">
      <template #left>
        <div class="header">
          <div style="display:flex;justify-content:space-between;">
            <div>
              日志仓信息列表 <el-button type="primary" style="height:30px;color: white; background-color:#0081dd"
                @click="search">
                <SvgIcon name="refresh-right" style="margin:0 5px 0 0" /><span>刷新</span>
              </el-button>
            </div>
          </div>
        </div>
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
      :tid="tid" :data="tableData" :height="tableHeight" class="c-gx-table c-glc-table" row-key="id">
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
const { visible, tableData, tableHeight, pageSettingStore, showTableLoadding } = usePageMainHooks(opt);

const table = ref(); // 表格实例
const tid = ref('storagesMain'); // 表格ID
const info = ref(''); // 底部提示信息

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
