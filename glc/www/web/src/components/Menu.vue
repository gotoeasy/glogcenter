<template>
  <el-menu
      :default-active="active+''"
      :collapse="isCollapsed"
      class="el-menu-vertical-demo"
      background-color="rgb(249,251,253)"
      text-color="#0081dd"
      active-text-color="#0081dd"
  >

    <router-link v-for="(item,index) in menus" :to="{name:item.name}" :key="item.name">
      <el-menu-item v-if="!item.hidden" :index="index">
        <Icon :name="item.icon" :color="item.color" class="side-menu-icon"/>
      </el-menu-item>
    </router-link>

  </el-menu>
</template>

<script>
import { Document, Service } from '@element-plus/icons-vue'
import menus from '../menus'
import Icon from './Icon.vue'

export default {
  name: 'Menu',
  components: { Icon, Service, Document },
  props: {
    isCollapsed: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      menus: menus,
    }
  },
  computed: {
    active() {
      let active = 0
      this.menus.forEach((item, index) => {
        if (item.name === this.$route.name) {
          active = index
        }
      })
      return active
    },
  },
}
</script>

<style scoped>
.el-menu a{
  text-decoration: none;
}
.el-menu-item{
  width: 48px;
  height: 48px;
  border-bottom:solid 1px lightgray;
  background-color: white;
  color: #0081dd;
}
.side-menu-icon{
  margin-left: -9px;
}
</style>
