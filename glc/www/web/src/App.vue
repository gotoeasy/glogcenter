
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
      </el-header>

    <el-aside class="menubar x-menubar" style="width:48px">
      <Menu :isCollapsed="true"></Menu>
    </el-aside>

      <el-main class="main x-main" >
        <router-view></router-view> 
      </el-main>

    </el-container>
  </el-container>
</template>

<script>
import Menu from './components/Menu.vue'
import { Expand, Fold } from '@element-plus/icons-vue'

export default {
  components: {
    Fold,
    Expand,
    Menu,
  },
  data() {
  },
  methods: {
    clickLogo() {
      location.href = "https://github.com/gotoeasy/glogcenter"
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