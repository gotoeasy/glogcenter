<template>
  <el-form ref="form" :inline="true" :model="formData" :rules="formRules" label-width="100">
    <el-row>
      <el-form-item class="c-search-form-item">
        <el-input v-model="formData.searchKeys" placeholder="请输入关键词检索，支持多关键词" input-style="width:500px;height: 30px"
          maxlength="1000" @keyup.enter="fnSearch">
          <template #append>
            <el-button type="primary" class="c-btn-search" style="height:30px;color:white"
              :style="{ backgroundColor: searchFormNormalMode ? '#0081dd' : '#e6a23c' }" @click="fnSearch">
              <el-icon>
                <Search />
              </el-icon>
              <span>检 索</span>
            </el-button>
          </template>
        </el-input>
        <GxButton icon="refresh-left" style="margin-left:108px;" @click="fnReset">
          <span v-if="searchFormNormalMode">重 置</span>
          <span v-else title="取消定位相邻检索">取 消</span>
        </GxButton>
        <el-badge is-dot :hidden="noMoreSearchCondition" style="margin-left:12px;" type="primary" class="c-btn-badge">
          <el-button circle @click="() => (moreVisible = !moreVisible)">
            <el-icon>
              <ArrowUp v-if="moreVisible" />
              <ArrowDown v-else />
            </el-icon>
          </el-button>
        </el-badge>
      </el-form-item>
    </el-row>

    <div v-show="moreVisible" class="c-down-panel">
      <slot></slot>

      <el-divider style="margin: 0 0 10px" />

      <el-row justify="center">
        <el-button :type="searchFormNormalMode ? 'primary' : 'warning'" class="c-btn" @click="fnSearch">
          <el-icon size="14">
            <Search />
          </el-icon>
          <span>检 索</span>
        </el-button>
        <el-button class="c-btn" @click="() => (moreVisible = false)">
          <el-icon size="14">
            <ArrowUp />
          </el-icon>
          <span>收 起</span>
        </el-button>
      </el-row>
    </div>
  </el-form>
</template>

<script setup>
const props = defineProps({
  data: {
    type: Object,
    default() {
      return {};
    },
  },
  rules: {
    type: Object,
    default() {
      return {};
    },
  },
});

const formData = ref(props.data);
const defaultData = ref({ ...props.data });
const formRules = ref(props.rules);
const form = ref();
const moreVisible = ref(false);

$emitter.on("defaultStorageCondtion", v => defaultData.value.storage = v);

// 检索模式(普通/相邻)
const searchFormNormalMode = ref(true);
const tmpFormDataNormal = {};
const setSearchNormalMode = (normal = false) => {
  if (searchFormNormalMode.value && !normal) {
    // 正切换模式，保存普通检索的条件，用于取消时恢复
    tmpFormDataNormal.storage = formData.value.storage || '';
    tmpFormDataNormal.system = formData.value.system || '';
    tmpFormDataNormal.loglevel = [...(formData.value.loglevel || [])];
    tmpFormDataNormal.datetime = [...(formData.value.datetime || [])];
    tmpFormDataNormal.user = formData.value.user || '';
    tmpFormDataNormal.searchKeys = formData.value.searchKeys || '';
  }
  searchFormNormalMode.value = normal;
  moreVisible.value = false; // 收起，避免点相邻检索时还展开着条件
}
$emitter.on("searchFormNormalMode", setSearchNormalMode);

const emit = defineEmits(['search']);

const fnSearch = () => {
  moreVisible.value = false;
  emit('search', 1);
}
$emitter.on("fnSearch", fnSearch);

const fnReset = () => {

  if (!searchFormNormalMode.value) {
    // 【取消】相邻检索模式 --> 默认检索模式，恢复正常检索的条件
    searchFormNormalMode.value = true;
    $emitter.emit("searchNearMode", false);

    formData.value.storage = tmpFormDataNormal.storage;
    formData.value.system = tmpFormDataNormal.system;
    !formData.value.loglevel && (formData.value.loglevel = []);
    formData.value.loglevel.splice(0, formData.value.loglevel.length);
    formData.value.loglevel.push(...(tmpFormDataNormal.loglevel || []));
    !formData.value.datetime && (formData.value.datetime = []);
    formData.value.datetime.splice(0, formData.value.datetime.length);
    formData.value.datetime.push(...(tmpFormDataNormal.datetime || []));
    formData.value.user = tmpFormDataNormal.user;
    formData.value.searchKeys = tmpFormDataNormal.searchKeys;
  } else {
    // 【重置】正常检索时的重置
    formData.value.storage = defaultData.value.storage;
    formData.value.system = defaultData.value.system;
    !formData.value.loglevel && (formData.value.loglevel = []);
    formData.value.loglevel.splice(0, formData.value.loglevel.length);
    formData.value.loglevel.push(...(defaultData.value.loglevel || []));
    !formData.value.datetime && (formData.value.datetime = []);
    formData.value.datetime.splice(0, formData.value.datetime.length);
    formData.value.datetime.push(...(defaultData.value.datetime || []));
    formData.value.user = defaultData.value.user;
    formData.value.searchKeys = '';
  }

  moreVisible.value = false;
  emit('search', 1);
}

const noMoreSearchCondition = computed(() => {

  // 相邻检索模式下，级别时间用户任一有填写即有修改
  if (!searchFormNormalMode.value) {
    if ((formData.value.loglevel || []).length) return false;
    const ary = formData.value.datetime || ['', ''];
    !ary.length && ary.push(...['', '']);
    if (ary.length != 2 || ary[0] || ary[1]) return false;
    if (formData.value.user) return false;
    return true;
  }

  // 普通检索模式下，逐个比较展开的条件字段
  if ((formData.value.storage || '') != (defaultData.value.storage || '')) {
    return false;
  }
  if ((formData.value.system || '') != (defaultData.value.system || '')) {
    return false;
  }
  if ((formData.value.user || '') != (defaultData.value.user || '')) {
    return false;
  }

  // 日志
  const formDataLoglevel = formData.value.loglevel || [];
  const defaultDataLoglevel = defaultData.value.loglevel || [];
  if (formDataLoglevel.length != defaultDataLoglevel.length) {
    return false;
  }
  for (let i = 0; i < formDataLoglevel.length; i++) {
    if (!defaultDataLoglevel.includes(formDataLoglevel[i])) {
      return false;
    }
  }

  // 时间范围
  const formDataDatetime = formData.value.datetime || [];
  !formDataDatetime.length && formDataDatetime.push(...['', '']);
  const defaultDataDatetime = defaultData.value.datetime || [];
  !defaultDataDatetime.length && defaultDataDatetime.push(...['', '']);
  if (formDataDatetime.length != defaultDataDatetime.length || formDataDatetime[0] != defaultDataDatetime[0] || formDataDatetime[1] != defaultDataDatetime[1]) {
    return false;
  }

  return true;
})

</script>
