<template>
  <el-container>
    
    <el-main style="padding-top:0">
      <div id="dashboard">
        <el-card>
          <template #header>

            <div class="header">
              <div style="display:flex;justify-content:space-between;">
                <div>
                  日志仓信息列表 <el-button type="primary" @click="search">刷新</el-button>
                </div>
              </div>
            </div>

          </template>
 

          <el-table :stripe="true" v-loading="loading" :data="data" :height="tableHeight" style="width: 100%">

            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="name" label="名称"/>
            <el-table-column prop="logCount" label="日志数量" />
            <el-table-column prop="indexCount" label="已建索引数量" />
            <el-table-column prop="fileCount" label="文件数量" />
            <el-table-column prop="totalSize" label="空间占用" />
            <el-table-column fixed="right" label="操作" width="100">
              <template #default="scope">
                <el-button size="small" type="warning" @click="remove(scope.row)">删除</el-button>
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
        searchKey: '',
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

    this.search() 

  },
  methods: {
    remove(row) {

      this.$confirm('确定删除日志仓 ' + row.name + ' 吗？', '确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        this.loading = true
        api.deleteStorage(row.name).then(rs => {
          let res = rs.data
          if (res.success) {
            this.$message({type: 'info', message: "已删除日志仓 " + row.name});
            this.search();
          }else if (res.code == 403){
            api.logout();
          }else{
            this.$message({type: 'error', message: res.message});
          }
        }).finally(() => {
          this.loading = false
        })
      }).catch(() => {
        // ignore
      })

    },
    search() {
      this.loading = true

      api.searchStorages(this.params).then(rs => {
        let res = rs.data
        if (res.success) {
          this.data = res.result.data || [];
          this.info =  res.result.info;
          // document.querySelector('.el-scrollbar__wrap').scrollTop = 0; // 滚动到顶部
        }else if (res.code == 403){
          api.logout();
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

.x-info{
  min-height:42px;
  line-height:42px;
  padding-top:5px;
  font-size: 14px;
  font-weight: 500;
  color: #909399;
}
</style>
