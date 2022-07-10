<template>
  <el-dialog title="登录" class="x-login" v-model="dialogVisible" width="420px"
    :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false"
    :before-close="handleClose">


    <el-input placeholder="请输入用户名" v-model="username" maxlength="100"></el-input><p/>
    <el-input placeholder="请输入密码" type="password" v-model="password" maxlength="100"></el-input>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="login" style="width:100%;">确 定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import api from '../api'

export default {
  data() {
    return {
      dialogVisible: true,
      username: "",
      password: "",
    };
  },
  methods: {
    login(){
      this.loading = true
      api.login(this.username,this.password).then(rs => {
        let res = rs.data
        // console.info('-------------------------------', res)
        if (res.success) {
          sessionStorage.setItem("glctoken", res.result)
          location.reload()
        }else{
          this.$message({type: 'error', message: res.message});
        }
      }).finally(() => {
        this.loading = false
      })
    }
  }
};
</script>

<style>


</style>