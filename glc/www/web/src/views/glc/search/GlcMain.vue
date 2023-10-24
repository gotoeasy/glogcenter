<template>
  <div v-show="visible">

    <GxToolbar style="margin-bottom: 13px" class="c-btn">
      <template #left>
        <SearchForm :data="formData" class="c-search-form" @search="search">
          <el-row>
            <el-form-item label="选择日志仓">
              <el-select v-model="formData.storage" filterable placeholder="请选择" style="width:420px;">
                <el-option v-for="item in storageOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="系统名">
              <el-select v-model="formData.system" :multiple="false" filterable allow-create default-first-option
                style="width:420px;" clearable :reserve-keyword="true" placeholder="请输入系统名">
                <el-option v-for="item in systemOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志级别">
              <el-select v-model="formData.loglevel" multiple clearable :reserve-keyword="true" style="width:420px;"
                placeholder="请选择...">
                <el-option label="ERROR" value="error" />
                <el-option label="WARN" value="warn" />
                <el-option label="INFO" value="info" />
                <el-option label="DEBUG" value="debug" />
              </el-select>
            </el-form-item>
            <el-form-item label="时间范围">
              <el-date-picker v-model="formData.datetime" type="datetimerange" :shortcuts="shortcuts" range-separator="～"
                value-format="YYYY-MM-DD HH:mm:ss" start-placeholder="开始时间" end-placeholder="结束时间"
                popper-class="c-datapicker" />
            </el-form-item>
          </el-row>
        </SearchForm>
      </template>

      <template #right>
        <el-tooltip v-if="showTestBtn" content="生成测试数据" placement="top">
          <el-button circle @click="genTestData">
            <SvgIcon name="quick" />
          </el-button>
        </el-tooltip>
        <el-tooltip :content="autoSearchMode ? '停止自动查询' : '开始自动查询'" placement="top">
          <el-button circle @click="switchAutoSearchMode">
            <SvgIcon v-if="!autoSearchMode" name="play" />
            <SvgIcon v-if="autoSearchMode" name="stop" />
          </el-button>
        </el-tooltip>
        <el-tooltip content="缩放" placement="top">
          <el-button circle @click="emitter.emit('main:switchMaximizePage')">
            <SvgIcon name="zoom" />
          </el-button>
        </el-tooltip>
        <el-tooltip content="下载当前检索结果" placement="top">
          <el-button circle @click="fnDownload">
            <SvgIcon name="download" />
          </el-button>
        </el-tooltip>
        <GxPageTableConfig :tid="tid" :page-config="pageSettingStore" />
      </template>
    </GxToolbar>

    <GxTable ref="table" v-loading="showTableLoadding" scrollbar-always-on :enable-header-contextmenu="false"
      :enable-first-expand="true" stripe :tid="tid" :data="tableData" :height="tableHeight" class="c-gx-table c-glc-table"
      row-key="id">
    </GxTable>

    <div>
      <div style="display:flex;justify-content:space-between;">
        <div style="min-height:30px;padding-top:5px; line-height:30px; color: #909399;" v-text="info"></div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { useEmitter, usePageMainHooks, useTabsState } from "~/pkgs";
import { userLogout } from "~/api";

const tabsState = useTabsState();
const emitter = useEmitter(tabsState.activePath);
const router = useRouter();

const opt = {
  emitter,
  withoutSearchForm: true,
};
const { formData, visible, tableData, tableHeight, pageSettingStore, showTableLoadding } = usePageMainHooks(opt);

const showTestBtn = ref(false); // 是否显示生成测试数据按钮
const autoSearchMode = ref(false); // 自动查询
const table = ref(); // 表格实例
const tid = ref('glcSearchMain'); // 表格ID
const info = ref(''); // 底部提示信息
const storageOptions = ref([]) // 日志仓
const systemSet = new Set();
const systemOptions = ref([]) // 系统名
const shortcuts = ref([
  {
    text: '近5分钟',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 5 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近10分钟',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 10 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近15分钟',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 15 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近20分钟',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 20 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近30分钟',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 30 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近1小时',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 60 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近2小时',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 2 * 60 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近3小时',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 3 * 60 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
  {
    text: '近4小时',
    value: () => {
      const start = new Date()
      start.setTime(start.getTime() - 4 * 60 * 60 * 1000)
      const end = new Date()
      return [start, end]
    },
  },
]);

// 初期默认检索
onMounted(() => {
  const configStore = $emitter.emit('$table:config', { id: tid.value });
  !configStore.columns.length && $emitter.emit('$table:config', { id: tid.value, update: true }); // 首次使用开启默认布局
  // 日志仓列表查询
  const url = `/v1/store/names`;
  $post(url, {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const names = rs.result || [];
      for (let i = 0; i < names.length; i++) {
        storageOptions.value.push({ value: names[i], label: `日志仓：${names[i]}` })
      }
    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
      router.push('/login');
    }
  });

  // 查询是否测试模式
  $post('/v1/store/mode', {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    showTestBtn.value = rs.success && rs.result
  });

  // 滚动底部触发检索更多
  const scrollWrap = document.querySelector('.c-glc-table .el-scrollbar__wrap')
  scrollWrap.addEventListener('scroll', (e) => {
    // e.target.scrollTop 为0时是滚动到顶部，不触发自动加载前一页数据，重新检索吧
    const scrollDistance = e.target.scrollHeight - e.target.scrollTop - e.target.clientHeight
    if (scrollDistance <= 0) {
      searchMore(); // 滚动到了底部，自动加载后一页
    }
  });

  // 默认检索
  search();
});

// 生成测试数据
function genTestData() {
  $post('/v1/log/addTestData', {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    $msg.info(rs.message);
  });
}

function isAutoSearchMode() {
  return autoSearchMode.value
}

function switchAutoSearchMode(changMode = true) {
  changMode && (autoSearchMode.value = !autoSearchMode.value);
  if (autoSearchMode.value) {
    search();
    setTimeout(() => {
      isAutoSearchMode() && switchAutoSearchMode(false);
    }, 5000);
  }
}

function search() {
  autoSearchMode.value ? (showTableLoadding.value = false) : (showTableLoadding.value = true);
  const url = `/v1/log/search`;
  const data = {};
  data.searchKey = formData.value.searchKeys;
  data.storeName = formData.value.storage;
  data.system = formData.value.system;
  data.loglevel = (formData.value.loglevel || []).join(',');
  data.datetimeFrom = (formData.value.datetime || ['', ''])[0];
  data.datetimeTo = (formData.value.datetime || ['', ''])[1];

  $post(url, data, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const resultData = rs.result.data || [];
      const pagesize = rs.result.pagesize - 0;
      tableData.value.splice(0, tableData.value.length);  // 删除原全部元素，nextTick时再插入新查询结果
      document.querySelector('.c-glc-table .el-scrollbar__wrap').scrollTop = 0; // 滚动到顶部

      nextTick(() => {
        resultData.forEach(item => {
          tableData.value.push(item);
          item.system && !systemSet.has(item.system) && systemSet.add(item.system) && systemOptions.value.push({ value: item.system, label: item.system });
        });

        if (resultData.length < pagesize) {
          info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${tableData.value.length} 条，正展示前 ${tableData.value.length} 条`
        } else {
          info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${rs.result.count} 条，正展示前 ${tableData.value.length} 条`
        }
      });

    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
      router.push('/login');
    }
  }).finally(() => {
    showTableLoadding.value = false;
  });
}

function searchMore() {
  if (tableData.value.length >= 5000) {
    if (info.value.indexOf('请考虑') < 0) {
      info.value += ` （数据太多不再自动加载，请考虑添加条件）`
    }
    return
  }

  const url = `/v1/log/search`;
  const data = {};
  data.searchKey = formData.value.searchKeys;
  data.storeName = formData.value.storage;
  data.system = formData.value.system;
  data.loglevel = (formData.value.loglevel || []).join(',');
  data.datetimeFrom = (formData.value.datetime || ['', ''])[0];
  data.datetimeTo = (formData.value.datetime || ['', ''])[1];
  data.forward = true
  data.currentId = tableData.value[tableData.value.length - 1].id; // 相对最后条id，继续找后面的日志

  $post(url, data, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const resultData = rs.result.data || [];
      const pagesize = rs.result.pagesize - 0;
      tableData.value.push(...resultData)

      if (resultData.length < pagesize) {
        info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${tableData.value.length} 条，正展示前 ${tableData.value.length} 条`
      } else {
        info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${rs.result.count} 条，正展示前 ${tableData.value.length} 条`
      }

      nextTick(() => {
        resultData.forEach(item => {
          item.system && !systemSet.has(item.system) && systemSet.add(item.system) && systemOptions.value.push({ value: item.system, label: item.system });
        });
      });
    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
    }
  })
}

// 下载当前检索结果
function fnDownload() {
  let fileContent = '';
  const tableConfigStore = $emitter.emit('$table:config', { id: tid.value })
  tableData.value.forEach(item => {
    let flg = false;
    tableConfigStore.columns.forEach(oCol => {
      if (!oCol.hidden && !oCol.editType.startsWith('$') && oCol.editType != 'text') {
        flg && (fileContent += ',');
        fileContent += item[oCol.field];
        flg = true;
      }
    })
    fileContent += ',';
    fileContent += (item.detail || item.text);
    fileContent += '\r\n';
  })

  const blob = new Blob([fileContent], { type: 'text/plain' });  // 创建Blob对象
  const downloadLink = document.createElement('a');
  downloadLink.href = URL.createObjectURL(blob);
  downloadLink.download = 'logfile.txt'; // 文件名
  downloadLink.click(); // 模拟点击下载链接
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
