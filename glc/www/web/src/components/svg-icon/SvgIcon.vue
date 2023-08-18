<template>
  <svg class="svg-icon" :height="size" :width="size">
    <use :xlink:href="symbolId" />
  </svg>
</template>

<script setup >
const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  size: {
    type: String,
    default: "14",
    required: false,
  },
});

const symbolId = computed(() => `#gx_icon-${props.name}`);
</script>

<style lang="scss">
.svg-icon {
  margin: 0 10px;
}
</style>>
