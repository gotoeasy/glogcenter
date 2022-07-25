
<template>
  <el-container>
    <el-container>

      <el-header style="display:flex;justify-content: space-between;border-bottom: 1px solid lightgray;height:52px;">
      
      <div style="display:flex;align-items: center" class="header-logo-title">
        <div style="text-align:center;width: 34px;color:white;cursor:pointer;" @click="clickLogo">
          <img src="/glc.png" style="width:34px;margin-top:9px;"/>
        </div> 
        <div style="text-align:center;">
          <img src="/title.png" style="margin-top:11px;margin-left:26px;height:32px"/>
        </div> 

      </div>

      <el-link v-if="inited && (enableLogin && isLogin)" @click="logout">
        <el-icon :size="26"><expand/></el-icon>
      </el-link>

      </el-header>

      <el-aside v-if="inited && (!enableLogin || isLogin)" class="menubar x-menubar" style="width:48px">
        <Menu :isCollapsed="true"></Menu>
      </el-aside>

      <el-main v-if="inited && (!enableLogin || isLogin)" class="main x-main" >
        <router-view></router-view> 
      </el-main>

      <el-dialog v-if="inited && enableLogin && !isLogin" title="登录" class="x-login" v-model="dialogVisible" width="420px"
        :close-on-click-modal="false" :close-on-press-escape="false" :show-close="false">

        <el-input placeholder="请输入用户名" v-model="username" maxlength="100"></el-input><p/>
        <el-input placeholder="请输入密码" type="password" v-model="password" maxlength="100"></el-input>
        
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="login" style="width:100%;">确 定</el-button>
          </span>
        </template>
      </el-dialog>

    </el-container>
  </el-container>
</template>

<script>
import Menu from './components/Menu.vue'
import { Expand, Fold } from '@element-plus/icons-vue'
import api from './api'

export default {
  components: {
    Fold,
    Expand,
    Menu
  },
  data() {
    return {
      inited: false,
      enableLogin: true,
      isLogin: false,
      dialogVisible: true,
      username: "",
      password: "",
    }
  },
  created() {
    api.initApp(this);
    api.enableLogin().then(rs => {
      let res = rs.data
      if (res.success) {
        this.inited = true;
        this.enableLogin = res.result;
        if (res.result) {
          this.isLogin = !!sessionStorage.getItem('glctoken');
        }
      }
    });
  },
  methods: {
    clickLogo() {
      location.href = "https://github.com/gotoeasy/glogcenter"
    },
    login(){
      this.loading = true
      api.login(this.username,this.password).then(rs => {
        let res = rs.data
        if (res.success) {
          sessionStorage.setItem("glctoken", res.result);
          this.dialogVisible = false;
          this.isLogin = true;
          this.username = '';
          this.password = '';
          this.$router.push({name:'dashboard'})
        }else{
          this.$message({type: 'error', message: res.message});
        }
      }).finally(() => {
        this.loading = false
      })
    },
    logout() {
      sessionStorage.clear();
      this.isLogin = false;
      this.dialogVisible = true;
      this.username = '';
      this.password = '';
    },
  }
}
</script>

<style scoped>
.x-menubar{
  position: fixed;
  top: 52px;
  overflow: hidden;
}
.main.x-main{
  margin-left: 48px;
  padding-left: 0px;
  padding-top: 20px;
  padding-right: 0px;
}

.header-logo-title{
  margin-left: -14px;
}
</style>