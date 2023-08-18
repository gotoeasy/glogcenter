<template>
  <div style="display: flex;flex-wrap: wrap;">
    <div v-for="name in icons" :key="name" class="c-svgicon-item" style="width: 90px;height: 90px;"
      @click="copyName(name)">
      <div style="display: flex;align-items: center;justify-content: center;width: 90px;height: 90px;cursor: pointer;">
        <div style="text-align: center;">
          <svg-icon :name="name" size="24" /> <br> <span>{{ name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup >
// src/assets/icons目录下的svg文件名  dir /b > e:/aaa.txt
const icons = ref(['account', 'add', 'apply', 'approve', 'arrow-left', 'close-bold', 'close', 'cloud-download', 'contract', 'dept', 'detail', 'dict', 'download', 'dragable',
  'edit-password', 'exit', 'home', 'item', 'log', 'maximize', 'menu', 'message', 'minimize', 'new', 'normalize', 'minus', 'refresh-left',
  'refresh-right', 'reload', 'role', 'save', 'scan', 'search', 'select', 'setting', 'side-menu', 'sort', 'success-filled', 'to-all', 'to-left', 'to-other',
  'send', 'flow-manage', 'flow', 'task-manage', 'task', 'finish',
  'to-right', 'upload', 'user', 'users', 'zoom'])

const props = defineProps({
  copyName: {
    type: Boolean,
    default: true,
  },
});
const emit = defineEmits(['select']);

const copyName = async (name) => {
  emit('select', name);
  props.copyName && navigator.clipboard?.writeText(name).then(() => {
    $msg.info(`已复制图标名称：${name}`);
  });
};
</script>

<style   lang="scss">
.c-svgicon-item:hover {
  background-color: var(--el-color-primary-light-9);
}
</style>
