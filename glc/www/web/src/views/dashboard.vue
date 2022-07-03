<template>
  <el-container>
    
    <el-main style="padding-top:0">
      <div id="dashboard">
        <el-card>
          <template #header>

            <div class="header">
              <div style="display:flex;justify-content:space-between;">
                <div>
                  <el-input @keyup.enter="search()" v-model="params.searchKey" placeholder="请输入关键词检索" style="width:600px">
                    <template #append>
                      <el-button type="primary" @click="search()" class="x-search">全文检索</el-button>
                    </template>
                  </el-input>
                </div>
              </div>
            </div>

          </template>
 

          <el-table :stripe="true" v-loading="loading" :data="data" :height="tableHeight" style="width: 100%">

            <el-table-column fixed type="expand" width="60">
              <template #default="scope">
                <div class="x-detail">
                  <el-scrollbar :class="{'x-scrollbar':(scope.row.detail && scope.row.detail.split('\n').length>20)}">
                    <div v-html="(scope.row.detail || scope.row.text).replace(/\n/g, '<br>')" style="word-break: break-all;"></div>
                  </el-scrollbar>
                </div>
              </template>
            </el-table-column>

            <el-table-column prop="system" label="分类" width="120"/>
            <el-table-column prop="date" label="日期时间" width="208"/>
            <el-table-column prop="text" label="内容" :show-overflow-tooltip="true">
              <template #default="scope">
                <span v-html="scope.row.text"></span>
              </template>
            </el-table-column>

          </el-table>

   
            <div class="header">
              <div style="display:flex;justify-content:space-between;">
                <div v-html="info" class="x-info"></div>
              </div>
            </div>

 
        </el-card>
      </div>
    </el-main>
  </el-container>
</template>

<script>
import api from '../api'
//import jsonViewer from 'vue-json-viewer'

const FixHeight = 215  // 177

export default {
  name: 'dashboard',
  components: {  },
  data() {
    return {
      loading: false,
      tableHeight: (window.innerHeight - FixHeight) + 'px',
      params: {
        storeName: '',
        searchKey: '',
        pageSize: 100,
        currentId: '',
        forward: true,
      },
      data: [],
      info: '',
    }
  },
  mounted(){

    // 自适应表格高度
    let that = this;
    window.onresize = () => {
      that.tableHeight = (window.innerHeight - FixHeight) + 'px';
    };

    // 
    let scrollWrap = document.querySelector('.el-scrollbar__wrap')
    scrollWrap.addEventListener('scroll', function(e) {
      // e.target.scrollTop 为0时是滚动到顶部，不触发自动加载前一页数据，重新检索吧
      let scrollDistance = e.target.scrollHeight - e.target.scrollTop - e.target.clientHeight
      if (scrollDistance <= 0) {
        that.searchMore(); // 滚动到了底部，自动加载后一页
      }
    })

    this.search() 

  },
  methods: {
    searchMore() {
      if (this.data.length >= 5000) {
        if (this.info.indexOf('请考虑') < 0){
          this.info += ` （数据太多不再自动加载，请考虑添加条件）`
        }
        return
      }
      let params = Object.assign({}, this.params);
      params.forward = true
      params.currentId = this.data[this.data.length-1].id; // 相对最后条id，继续找后面的日志

      api.search(params).then(rs => {
        let res = rs.data
        if (res.success) {
          if (res.result.data) {
            this.data.push(...res.result.data)
            if (res.result.data.length < params.pageSize) {
              this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${this.data.length} 条，正展示前 ${this.data.length} 条`
            }else{
              this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${res.result.count} 条，正展示前 ${this.data.length} 条`
            }
          }else{
            this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${this.data.length} 条，正展示前 ${this.data.length} 条`
          }
        }
      })
    },
    search() {
      this.loading = true

      //console.info("----------this.params",this.params)
      this.params.currentId = ''
      api.search(this.params).then(rs => {
        let res = rs.data
        if (res.success) {
         // console.info(res,"xxxxxxxxxxxxxxxxxxxxxxx")
          this.data = res.result.data || [];
          document.querySelector('.el-scrollbar__wrap').scrollTop = 0; // 滚动到顶部
          this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${res.result.count} 条，正展示前 ${this.data.length} 条`
        }

      }).finally(() => {
        this.loading = false
      })
    },
  },
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
}

.database .item {
  line-height: 35px;
  border-bottom: 1px solid var(--el-card-border-color);
  cursor: pointer;
  transition: background-color .3s;
  padding: 0 5px;
}

.database .item:hover {
  background-color: #f3f6f9;
}

.database .item .name {
  margin-left: 10px;
}


.database .active {
  background-color: #fef5ea;
  color: var(--el-color-primary);
}
button.el-button.x-search{
    background-color: #0081dd;
    color: white;
}

.x-detail{
    padding-top:5px;
    padding-left:30px;
    padding-right:5px;
    padding-bottom:5px;
    background-color: floralwhite;
}
.x-scrollbar{
    height:420px;
}

.x-info{
  min-height:42px;
  line-height:42px;
  padding-top:5px;
  font-size: 14px;
  font-weight: 500;
  color: #909399;
}
</style>

<style>
.el-popper.is-dark{
  display: none;
}
</style>
