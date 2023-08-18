<template>
  <el-config-provider :locale="zhCn" :size="'small'">
    <router-view />
  </el-config-provider>
</template>

<script setup>
import zhCn from "element-plus/es/locale/lang/zh-cn";

// 停止首次加载的loadding
onMounted(() => {
  document.getElementById('app-loading')?.remove();
  document.getElementById('app-loading-link')?.remove();
});

</script>
