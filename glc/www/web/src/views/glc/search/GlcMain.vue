<template>
  <div v-show="visible">

    <GxToolbar style="margin-bottom: 13px" class="c-btn">
      <template #left>
        <SearchForm :data="formData" class="c-search-form" @search="() => search()">
          <el-row>
            <el-form-item label="选择日志仓">
              <el-select v-model="formData.storage" clearable filterable placeholder="请选择" style="width:420px;"
                :disabled="searchNearMode" @clear="reGetStorageOptions">
                <el-option v-for="item in storageOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="系统名">
              <el-select v-model="formData.system" :multiple="false" filterable allow-create default-first-option
                :disabled="searchNearMode" style="width:420px;" clearable :reserve-keyword="false" placeholder="请输入系统名">
                <el-option v-for="item in systemOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志级别">
              <!-- <el-select v-model="formData.loglevel" multiple clearable :reserve-keyword="true" style="width:420px;"
                placeholder="请选择...">
                <el-option label="ERROR" value="error" />
                <el-option label="WARN" value="warn" />
                <el-option label="INFO" value="info" />
                <el-option label="DEBUG" value="debug" />
              </el-select> -->
              <el-checkbox-group v-model="formData.loglevel">
                <el-checkbox label="DEBUG" value="debug" border />
                <el-checkbox label="INFO" value="info" border style="margin-left: 6px;" />
                <el-checkbox label="WARN" value="warn" border style="margin-left: 6px;" />
                <el-checkbox label="ERROR" value="error" border style="margin-left: 6px;" />
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="时间范围">
              <el-date-picker v-model="formData.datetime" type="datetimerange" :shortcuts="shortcuts"
                range-separator="～" value-format="YYYY-MM-DD HH:mm:ss" start-placeholder="开始时间" end-placeholder="结束时间"
                popper-class="c-datapicker" />
            </el-form-item>
            <el-form-item label="用户">
              <el-input v-model="formData.user" placeholder="请输入用户" maxlength="100" style="width:420px;"
                @keyup.enter="() => { $emitter.emit('fnSearch'); }" />
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
          <el-button circle :disabled="searchNearMode" @click="switchAutoSearchMode">
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

    <GxTable ref="table" v-loading="showTableLoadding" :row-class-name="tableRowClassName" scrollbar-always-on
      :enable-header-contextmenu="false" :enable-first-expand="true" stripe :tid="tid" :data="tableData"
      :height="tableHeight" class="c-gx-table c-glc-table" row-key="id">
      <template #$operation="{ row }">
        <div title="定位相邻检索">
          <SvgIcon name="near-search" class="hand" @click="fnSearchNear(row)" />
        </div>
      </template>
    </GxTable>

    <div>
      <div style="display:flex;justify-content:space-between;">
        <div style="min-height:30px;padding-top:5px; line-height:30px; color: #909399;" v-text="info"> </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useEmitter, $emitter, usePageMainHooks, useTabsState } from "~/pkgs";
import { userLogout, enableLogin } from "~/api";

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
const tid = ref('glcMain241130'); // 表格ID
const info = ref(''); // 底部提示信息
const storageOptions = ref([]) // 日志仓
const systemSet = new Set();
const systemOptions = ref([]) // 系统名
const lastStoreName = ref('') // 检索结果中最久远的一条日志所属的日志仓名称
const maxMatchCount = ref('0') // 最大匹配件数(字符串)
const moreConditon = ref(null)
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

// 相邻检索模式
const searchNearMode = ref(false);
const setSearchNearMode = (near = false) => {
  searchNearMode.value = near;
}
$emitter.on("searchNearMode", setSearchNearMode);

const oldNearId = ref(0);
const newNearId = ref(0);
const nearStoreName = ref('');
const tableRowClassName = ({ row }) => {
  if (searchNearMode.value && row.id == newNearId.value) {
    return 'search-near-row';
  }
  return '';
}

// 定位相邻检索
const fnSearchNear = (row) => {
  $emitter.emit("searchFormNormalMode", false); // 收起条件等处理

  if (searchNearMode.value) {
    // 相邻检索模式下，继续进行相邻检索，需使用输入的条件进一步筛选
    if (newNearId.value != row.id) {
      oldNearId.value = newNearId.value;
      newNearId.value = row.id;
    }
  } else {
    // 普通检索模式下，进行相邻检索
    oldNearId.value = 0;
    newNearId.value = row.id;
    nearStoreName.value = row.storename;

    formData.value.storage = row.storename;
    formData.value.system = row.system;
    formData.value.loglevel = [];
    formData.value.datetime = [];
    formData.value.user = '';
    formData.value.searchKeys = '';

    $emitter.emit("searchNearMode", true);
  }

  // 检索
  search();
}

// 初期默认检索
onMounted(async () => {
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
      if (names[0]) {
        $emitter.emit("defaultStorageCondtion", names[0]); // 小蓝点提示判断用
        formData.value.storage = names[0]; // 选中第一个日志仓作为默认条件
      }

      // 默认检索（查取好日志仓后再做检索）
      search();

    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
      router.push('/glc/login');
    }
  });

  // 查询有权限的系统名
  $post('/v1/store/systems', {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    if (rs.success && rs.result?.length) {
      for (let i = 0; i < rs.result.length; i++) {
        systemSet.add(rs.result[i]);
        systemSet.add(rs.result[i]) && systemOptions.value.push({ value: rs.result[i], label: rs.result[i] });
      }
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

  // 检查下，避免不需要登录又还显示着登录状态
  await enableLogin();
});

// 清除日志仓条件时，重新拉取最新日志仓列表
function reGetStorageOptions() {
  const url = `/v1/store/names`;
  $post(url, {}, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const names = rs.result || [];
      if (names.length) {
        storageOptions.value.splice(0, storageOptions.value.length)
        for (let i = 0; i < names.length; i++) {
          storageOptions.value.push({ value: names[i], label: `日志仓：${names[i]}` })
        }
      }
      if (names[0]) {
        $emitter.emit("defaultStorageCondtion", names[0]); // 小蓝点提示判断用
      }
    }
  });
}

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
  if (autoSearchMode.value && !searchNearMode.value) {
    search();
    setTimeout(() => {
      isAutoSearchMode() && switchAutoSearchMode(false);
    }, 5000);
  }
}

function search() {
  autoSearchMode.value ? (showTableLoadding.value = false) : (showTableLoadding.value = true);
  const url = `/v1/log/search`;

  // 检索条件
  const data = {};
  data.searchKey = formData.value.searchKeys;
  data.storeName = formData.value.storage;
  data.system = formData.value.system;
  data.loglevel = (formData.value.loglevel || []).join(',');
  data.datetimeFrom = (formData.value.datetime || ['', ''])[0];
  data.datetimeTo = (formData.value.datetime || ['', ''])[1];
  data.user = formData.value.user;
  if (searchNearMode.value) {
    data.oldNearId = oldNearId.value;
    data.newNearId = newNearId.value;
    data.nearStoreName = nearStoreName.value;
  }

  // 保存好滚动检索的输入条件，保持和检索时一致，避免修改输入再滚动查询而出现矛盾结果
  moreConditon.value = data;

  $post(url, data, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    if (rs.success) {
      const resultData = rs.result.data || [];
      const pagesize = rs.result.pagesize - 0;
      tableData.value.splice(0, tableData.value.length);  // 删除原全部元素，nextTick时再插入新查询结果
      !searchNearMode.value && (document.querySelector('.c-glc-table .el-scrollbar__wrap').scrollTop = 0); // 普通检索时滚动到顶部
      rs.result.laststorename && (lastStoreName.value = rs.result.laststorename); // 查到有结果时，更新
      maxMatchCount.value = rs.result.count; // 最大匹配件数

      nextTick(() => {
        if (searchNearMode.value) {
          // 相邻检索
          tableData.value.push(...resultData);
          info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${maxMatchCount.value} 条，正展示相邻 ${tableData.value.length - 1} 条，查询${rs.result.timemessage}`

          nextTick(() => {
            const step = 30;
            const down = oldNearId.value == 0 || newNearId.value <= oldNearId.value; // 向下定位检索
            let scrollTop = down ? (-7 * step) : (-15 * step); // 滚动到合适位置（向下时上部7行新日志，剩余下部较大空间给旧日志展示;向上时上部固定留出15行显示新日志）
            for (let i = 0; i < tableData.value.length; i++) {
              if (newNearId.value == tableData.value[i].id) {
                break;
              }
              scrollTop += step;
            }
            // scrollTop <= (step * 15) && (scrollTop = 0);
            document.querySelector('.c-glc-table .el-scrollbar__wrap').scrollTop = scrollTop; // 滚动到焦点行
            // if (oldNearId.value > 0 && newNearId.value > oldNearId.value) {
            //   document.querySelector('.c-glc-table .el-scrollbar__wrap').scrollTop = 20000; // 滚动到底部
            // }
          });
        } else {
          // 普通检索
          resultData.forEach(item => {
            tableData.value.push(item);
            item.system && !systemSet.has(item.system) && systemSet.add(item.system) && systemOptions.value.push({ value: item.system, label: item.system });
          });

          if (resultData.length < pagesize) {
            info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${tableData.value.length} 条，正展示前 ${tableData.value.length} 条，查询${rs.result.timemessage}`
          } else {
            info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${maxMatchCount.value} 条，正展示前 ${tableData.value.length} 条，查询${rs.result.timemessage}`
          }
        }
      });

    } else if (rs.code == 403) {
      userLogout(); // 403 时登出
      router.push('/glc/login');
    }
  }).finally(() => {
    showTableLoadding.value = false;
  });
}

function searchMore() {
  if (searchNearMode.value) {
    return; // 相邻检索模式不支持滚动检索
  }

  if (tableData.value.length >= 5000) {
    if (info.value.indexOf('请考虑') < 0) {
      info.value += ` （数据太多不再自动加载，请考虑添加条件）`
    }
    return;
  }

  const url = `/v1/log/search`;
  moreConditon.value.forward = true
  moreConditon.value.currentId = tableData.value[tableData.value.length - 1].id; // 相对最后条id，继续找后面的日志
  moreConditon.value.currentStoreName = lastStoreName.value;  // 相对最后条id的所属日志仓

  $post(url, moreConditon.value, null, { 'Content-Type': 'application/x-www-form-urlencoded' }).then(rs => {
    console.log(rs)
    if (rs.success) {
      const resultData = rs.result.data || [];
      const pagesize = rs.result.pagesize - 0;
      tableData.value.push(...resultData)
      rs.result.laststorename && (lastStoreName.value = rs.result.laststorename); // 查到有结果时，更新

      if (resultData.length < pagesize) {
        info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${tableData.value.length} 条，正展示前 ${tableData.value.length} 条，查询${rs.result.timemessage}`
      } else {
        (rs.result.count - 0 < maxMatchCount.value - 0) && (maxMatchCount.value = rs.result.count) // 控制maxMatchCount只会更小以减少误差
        info.value = `日志总量 ${rs.result.total} 条，当前条件最多匹配 ${maxMatchCount.value} 条，正展示前 ${tableData.value.length} 条，查询${rs.result.timemessage}`
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
      if (!oCol.hidden && !oCol.editType.startsWith('$') && oCol.field != 'text') {
        flg && (fileContent += ',');
        fileContent += item[oCol.field];
        flg = true;
      }
    })
    flg && (fileContent += ',');
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
.el-table--striped .el-table__body tr.el-table__row--striped.search-near-row td.el-table__cell {
  background: lemonchiffon;
}

.el-table__body tr.el-table__row.search-near-row td.el-table__cell {
  background: lemonchiffon;
}

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

.c-search-form .el-form-item--small .el-form-item__label {
  height: 30px;
  line-height: 30px;
}

.c-datapicker.el-popper.is-pure {
  margin-left: -100px;
}
</style>
