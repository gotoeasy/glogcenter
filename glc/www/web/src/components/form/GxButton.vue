<template>
  <el-button v-bind="$attrs" style="min-width:80px;height:30px">
    <SvgIcon v-if="icon" :name="icon" style="margin:0 5px 0 0" />
    <span>
      <slot></slot>
    </span>
  </el-button>
</template>

<script setup>
const props = defineProps({
  icon: {
    type: String,
    default: '',
  },
});

const icon = computed(() => props.icon);
</script>

