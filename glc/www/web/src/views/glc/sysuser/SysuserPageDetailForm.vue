<template>
  <el-form ref="form" :model="formData" :rules="rules" label-width="auto" label-position="right" status-icon>
    <el-card shadow="never">

      <el-form-item label="账号" prop="username">
        <el-input v-model="formData.username" :disabled="!isCreateMode" placeholder="请输入账号" maxlength="100" />
      </el-form-item>

      <el-form-item v-if="!readonly" label="密码" prop="password1">
        <el-input v-model="formData.password1" type="password" :placeholder="isCreateMode ? '请输入密码' : '需要修改密码时请输入'"
          autocomplete="new-password" maxlength="100" />
      </el-form-item>

      <el-form-item v-if="!readonly" label="确认密码" prop="password2">
        <el-input v-model="formData.password2" type="password" :placeholder="isCreateMode ? '请输入密码' : '需要修改密码时请输入'"
          autocomplete="new-password" maxlength="100" />
      </el-form-item>

      <el-form-item label="系统权限" prop="roles">
        <el-input v-model="formData.systems" :disabled="readonly"
          :placeholder="readonly ? '' : '请输入可访问系统，多系统时逗号分隔，不填或*代表全部'" maxlength="1000" />
      </el-form-item>

      <el-form-item label="备注">
        <el-input v-model="formData.note" type="textarea" :disabled="readonly" resize="none" maxlength="2000"
          :autosize="{ minRows: 4, maxRows: 4 }" :placeholder="readonly ? '' : '请输入备注'" />
      </el-form-item>

    </el-card>

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
  readonly: {
    type: Boolean,
    default: true,
  },
  isCreateMode: {
    type: Boolean,
    default: false,
  },
});

const formData = ref(props.data);
const form = ref(); // 表单实例

const validatePass2 = (rule, value, callback) => {
  if (props.isCreateMode && !value) {
    callback(new Error('请确认密码'))
  } else if (value !== formData.value.password1) {
    callback(new Error("两次输入密码不一致"))
  } else {
    callback()
  }
}
// 校验规则
const rules = reactive({
  username: [{ required: true, message: '请输入账号', trigger: 'blur' },],
  password1: [{ required: props.isCreateMode, message: '请输入密码', trigger: 'blur' },],
  password2: [{ required: props.isCreateMode, validator: validatePass2, trigger: 'blur' },],
});

const validate = (callback) => form.value.validate(callback)

// 导出方法
defineExpose({ validate });
</script>

