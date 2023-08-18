<template>
  <el-select v-bind="$attrs" multiple clearable style="width:100%" @clear="clear">
    <el-option v-for="item in data" :key="item.code" :value="item.code" :label="item.name" />
  </el-select>
</template>

<script setup >
const data = ref([]);

const getData = () => {
  const param = { code: 'id', name: 'role_name', dict: 'sys_role' };
  $post(`/common/api/dict`, param).then(rs => {
    if (rs.success) {
      data.value = rs.data;
    }
  });
}

// 初始化时，取数据
onBeforeMount(() => {
  getData();
})

// 清空输入时，重新查取字典列表
const clear = () => {
  getData();
};

</script>
