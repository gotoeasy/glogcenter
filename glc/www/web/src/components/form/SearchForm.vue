<template>
  <el-form ref="form" :inline="true" :model="formData" :rules="formRules" label-width="100">
    <el-row>
      <el-form-item class="c-search-form-item">
        <el-input v-model="formData.searchKeys" placeholder="请输入关键词检索，支持多关键词" input-style="width:500px;height: 30px"
          maxlength="1000" @keyup.enter="fnSearch">
          <template #append>
            <el-button type="primary" class="c-btn-search" style="height:30px;color: white; background-color:#0081dd"
              @click="fnSearch">
              <el-icon>
                <Search />
              </el-icon>
              <span>检 索</span>
            </el-button>
          </template>
        </el-input>
        <GxButton icon="refresh-left" style="margin-left:108px;" @click="fnReset">重 置</GxButton>
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
        <el-button type="primary" class="c-btn" @click="fnSearch">
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

const emit = defineEmits(['search']);

const fnSearch = () => {
  moreVisible.value = false;
  emit('search', 1);
}
$emitter.on("fnSearch", fnSearch);

const fnReset = () => {

  const keys = Object.keys(formData.value);
  for (const key of keys) {
    if (!formData.value[key] && !defaultData.value[key]) {
      continue;
    } else if (!formData.value[key] && defaultData.value[key]) {
      formData.value[key] = defaultData.value[key];
    } else if (formData.value[key] && !defaultData.value[key]) {
      formData.value[key] = null;
    } else if (Array.isArray(defaultData.value[key])) {
      formData.value[key] = defaultData.value[key].slice(0);
    } else {
      formData.value[key] = defaultData.value[key];
    }
  }

  moreVisible.value = false;
  emit('search', 1);
}

const noMoreSearchCondition = computed(() => {
  for (const [key, value] of Object.entries(formData.value)) {
    if (key == 'searchKeys') {
      continue; // 忽略不提示小蓝点
    }
    if (key == 'storage') {
      if (!defaultData.value.storage) {
        continue; // 初始化时日志仓还没拿到，忽略不提示小蓝点
      } else if (defaultData.value.storage == value) {
        continue; // 和默认的日志仓条件一样，忽略不提示小蓝点
      }
      if (!value) return false; // 日志仓条件清空时，显示小蓝点
    }
    if (value) {
      if (Array.isArray(value)) {
        if (!defaultData.value[key]) {
          if (!value.length) {
            continue;
          } else {
            return false;
          }
        } else {

          if (value.length != defaultData.value[key].length) {
            return false;
          }
          for (let i = 0; i < value.length; i++) {
            if (value[i] != defaultData.value[key][i]) {
              return false;
            }
          }
        }

      } else if (value != defaultData.value[key]) {
        return false;
      }
    } else if (defaultData.value[key]) {
      return false;
    }
  }
  return true;
})

</script>

