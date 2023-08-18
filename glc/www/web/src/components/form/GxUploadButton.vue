<template>
  <el-upload v-bind="$attrs" :show-file-list="false" :http-request="fnUpload" class="c-upload"
    style="margin-right: 8px;margin-left: 8px;">
    <el-button :type="props.type" style="min-width:80px;">
      <SvgIcon v-if="icon" :name="icon" style="margin:0 5px 0 0" />
      <span>
        <slot>上 传</slot>
      </span>
    </el-button>
  </el-upload>
</template>

<script setup>
import { $upload } from '~/api/request.axios';

const props = defineProps({
  action: {
    type: String,
    default: '',
  },
  type: {
    type: String,
    default: 'primary',
  },
  icon: {
    type: String,
    default: '',
  },
});

const fnUpload = async (options) => {
  const formData = new FormData();
  formData.append("file", options.file);
  $upload(props.action, formData).then(rs => {
    if (rs.success) {
      $msg.info(rs.message);
    } else {
      $msg.error(rs.message);
    }
  });

}
</script>

