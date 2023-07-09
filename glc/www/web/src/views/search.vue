<template>
  <el-container>
    
    <el-main style="padding-top:0">
      <div>
        <el-card>
          <template #header>

            <div class="header">
              <div style="display:flex;justify-content:space-between;width:100%">
                <div>

                  <el-input @keyup.enter="search()" v-model="params.searchKey" placeholder="请输入关键词检索，支持多关键词" style="width:600px;">
                    <template #append>
                      <el-button type="primary" @click="search" class="x-search">
                        <el-icon>
                          <Search />
                        </el-icon>
                        <span>检 索</span>                      
                      </el-button>
                    </template>
                  </el-input>

                  <el-button class="c-btn" @click="fnResetSearchForm" style="margin-left:10px;">
                    <el-icon>
                      <RefreshLeft />
                    </el-icon>
                    <span>重 置</span>
                  </el-button>

                  <el-badge is-dot :hidden="hasMoreCondition" type="primary" style="margin-left:10px;">
                    <el-button circle @click="() => (showSearchPanel = !showSearchPanel)">
                      <el-icon>
                        <ArrowUp v-if="showSearchPanel" />
                        <ArrowDown v-else />
                      </el-icon>
                    </el-button>
                  </el-badge>

                  <div v-show="showSearchPanel" class="c-down-panel">
                    <el-form ref="form" :inline="true" label-width="100">
                      <el-row>
                      <el-form-item label="选择日志仓">
                        <el-select v-if="storageOptions.length > 0" v-model="storage" filterable placeholder="请选择" style="width:420px;">
                          <el-option
                            v-for="item in storageOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                          />
                        </el-select>
                      </el-form-item>
                    </el-row>
                    <el-row>
                      <el-form-item label="系统名">
                        <el-select v-model="params.system" style="width:420px;" :multiple="false" filterable allow-create default-first-option clearable
                           :reserve-keyword="true" placeholder="请输入系统名">
                          <el-option v-for="item in systemOptions" :key="item.value" :label="item.label" :value="item.value"/>
                        </el-select>
                      </el-form-item>
                    </el-row>
                    <el-row>
                      <el-form-item label="日志级别">
                        <el-select v-model="params.loglevel" style="width:420px;" :multiple="false" clearable
                           :reserve-keyword="true" placeholder="请选择...">
                           <el-option label="ERROR" value="error"/>
                           <el-option label="WARN" value="warning"/>
                           <el-option label="INFO" value="info"/>
                           <el-option label="DEBUG" value="debug"/>
                        </el-select>
                      </el-form-item>
                    </el-row>
                    <el-row>
                      <el-form-item label="时间范围">
                        <el-date-picker v-model="params.datetime" type="datetimerange" :shortcuts="shortcuts"
                          range-separator="～" value-format="YYYY-MM-DD HH:mm:ss" start-placeholder="开始时间" end-placeholder="结束时间"/>
                      </el-form-item>
                    </el-row>
                  </el-form>
                    <el-divider style="margin: 0 0 10px;" />

                    <el-row justify="center">
                      <el-button type="primary" class="x-search" @click="search">
                        <el-icon size="14">
                          <Search />
                        </el-icon>
                        <span>检 索</span>
                      </el-button>
                      <el-button class="c-btn" @click="() => (showSearchPanel = false)">
                        <el-icon size="14">
                          <ArrowUp />
                        </el-icon>
                        <span>收 起</span>
                      </el-button>
                    </el-row>
                  </div>


                </div>

                <div style="display: flex;align-items: end;">

                  <el-dropdown trigger="click">
                    <el-button size="small">
                      <el-icon ><setting /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="switchVisible('system')">{{systemVisible?'隐藏':'显示'}}系统名</el-dropdown-item>
                        <el-dropdown-item @click="switchVisible('servername')">{{serverNameVisible?'隐藏':'显示'}}主机名</el-dropdown-item>
                        <el-dropdown-item @click="switchVisible('serverip')">{{serverIpVisible?'隐藏':'显示'}}主机IP</el-dropdown-item>
                        <el-dropdown-item @click="switchVisible('date')">{{dateVisible?'隐藏':'显示'}}日期时间</el-dropdown-item>
                        <el-dropdown-item @click="switchVisible('loglevel')">{{logLevelVisible?'隐藏':'显示'}}Level</el-dropdown-item>
                        <el-dropdown-item @click="switchVisible('traceid')">{{traceIdVisible?'隐藏':'显示'}}TraceID</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </div>
            </div>

          </template>
 

          <el-table :stripe="true" v-loading="loading" :data="data" :height="tableHeight" style="width: 100%">

            <el-table-column fixed type="expand" width="40">
              <template #default="scope">
                <div class="x-detail">
                  <el-scrollbar :class="{'x-scrollbar':(scope.row.detail && scope.row.detail.split('\n').length>20)}">
                    <div v-html="(scope.row.detail || scope.row.text).replace(/</g, '&amp;lt;').replace(/>/g, '&amp;gt;').replace(/\n/g, '<br>')" style="word-break: break-all;"></div>
                  </el-scrollbar>
                </div>
              </template>
            </el-table-column>

            <el-table-column v-if="systemVisible" prop="system" label="系统名" width="120" :show-overflow-tooltip="true"/>
            <el-table-column v-if="serverNameVisible" prop="servername" label="主机名" width="180" :show-overflow-tooltip="true"/>
            <el-table-column v-if="serverIpVisible" prop="serverip" label="主机IP" width="130" :show-overflow-tooltip="true"/>
            <el-table-column v-if="dateVisible" prop="date" label="日期时间" width="208" :show-overflow-tooltip="true"/>
            <el-table-column v-if="logLevelVisible" prop="loglevel" label="Level" width="75" :show-overflow-tooltip="true"/>
            <el-table-column v-if="traceIdVisible" prop="traceid" label="TraceId" width="300" :show-overflow-tooltip="true"/>
            <el-table-column prop="text" label="日志" :show-overflow-tooltip="true">
              <template #default="scope">
                <span v-text="scope.row.text"></span>
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
import { ref } from 'vue'
import { Search, RefreshLeft, ArrowUp, ArrowDown, Check } from '@element-plus/icons-vue'

const FixHeight = 215  // 177

export default {
  components: {  },
  data() {
    return {
      loading: false,
      tableHeight: (window.innerHeight - FixHeight) + 'px',
      params: {
        storeName: '',
        searchKey: '',
        loglevel: '',
        system: '',
        systems: [],
        datetime: null,
        pageSize: 100,
        currentId: '',
        forward: true,
      },
      data: [],
      systemVisible: true,
      serverNameVisible: false,
      serverIpVisible: false,
      logLevelVisible: true,
      traceIdVisible: false,
      dateVisible: true,
      info: '',
      storage: ref(''),
      storageOptions: ref([]),
      showSearchPanel: ref(false),
      systemOptions: ref([]),
      systemSet: new Set(),
      shortcuts: [
                    {
                      text: '近5分钟',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 5 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近10分钟',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 10 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近15分钟',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 15 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近20分钟',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 20 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近30分钟',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 30 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近1小时',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 60 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近2小时',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 2 * 60 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近3小时',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 3 * 60 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                    {
                      text: '近4小时',
                      value: () => {
                        const start = new Date()
                        start.setTime(start.getTime() - 4 * 60 * 60 * 1000)
                        const end = new Date()
                        return [start, end]
                      },
                    },
                  ],
    }
  },
  created(){
      api.searchStorageNames().then(rs => {
        let res = rs.data
       // console.info(res)
        if (res.success) {
          let names = res.result || [];
          for (let i = 0; i < names.length; i++) {
            this.storageOptions.push({value: names[i], label: '日志仓：' + names[i]})
          }
        }else if (res.code == 403){
          api.logout();
        }
      });
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
  computed:{
    hasMoreCondition(){
      return !this.params.datetime && !this.storage && !this.params.system && !this.params.loglevel;
    }
  },
  methods: {
    switchVisible(name){
      name == 'system' && (this.systemVisible = !this.systemVisible);
      name == 'servername' && (this.serverNameVisible = !this.serverNameVisible);
      name == 'serverip' && (this.serverIpVisible = !this.serverIpVisible);
      name == 'loglevel' && (this.logLevelVisible = !this.logLevelVisible);
      name == 'traceid' && (this.traceIdVisible = !this.traceIdVisible);
      name == 'date' && (this.dateVisible = !this.dateVisible);
    },
    fnResetSearchForm(){
      this.params.searchKey = '';
      this.params.system = '';
      this.params.loglevel = '';
      this.params.datetime = null;
      this.storage = '';
      this.search();
    },
    searchMore() {
      if (this.data.length >= 5000) {
        if (this.info.indexOf('请考虑') < 0){
          this.info += ` （数据太多不再自动加载，请考虑添加条件）`
        }
        return
      }
      let params = Object.assign({}, this.params);
      params.storeName = this.storage
      params.forward = true
      params.currentId = this.data[this.data.length-1].id; // 相对最后条id，继续找后面的日志

      api.search(params).then(rs => {
        let res = rs.data
        if (res.success) {
          if (res.result.data) {
            this.data.push(...res.result.data) 
            if (res.result.data.length < this.params.pageSize) {
              this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${this.data.length} 条，正展示前 ${this.data.length} 条`
            }else{
              this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${res.result.count} 条，正展示前 ${this.data.length} 条`
            }

            this.$nextTick(()=>{
              res.result.data.forEach(item =>{
                item.system && !this.systemSet.has(item.system) && this.systemSet.add(item.system) && this.systemOptions.push({value:item.system,label:item.system});
              });
            });

          }else{
            this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${this.data.length} 条，正展示前 ${this.data.length} 条`
          }
        }else if (res.code == 403){
          api.logout();
        }
      })
    },
    search() {
      this.loading = true

      this.params.storeName = this.storage
      this.params.currentId = ''
      this.params.datetimeFrom = (this.params.datetime || ['', ''])[0]
      this.params.datetimeTo = (this.params.datetime || ['', ''])[1]
      // console.info("----------this.params",this.params)
      api.search(this.params).then(rs => {
        let res = rs.data
        if (res.success) {
         // console.info(res,"xxxxxxxxxxxxxxxxxxxxxxx")
          this.data = res.result.data || [];
          document.querySelector('.el-scrollbar__wrap').scrollTop = 0; // 滚动到顶部
          if (this.data.length < this.params.pageSize) {
            this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${this.data.length} 条，正展示前 ${this.data.length} 条`
          }else{
            this.info = `日志总量 ${res.result.total} 条，当前条件最多匹配 ${res.result.count} 条，正展示前 ${this.data.length} 条`
          }

          this.$nextTick(()=>{
              res.result.data.forEach(item =>{
                item.system && !this.systemSet.has(item.system) && this.systemSet.add(item.system) && this.systemOptions.push({value:item.system,label:item.system});
              });
          });

        }else if (res.code == 403){
          api.logout();
        }

      }).finally(() => {
        this.loading = false
      })

      this.showSearchPanel = false;
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

.c-down-panel {
  position: absolute;
  z-index: 100;
  width: 550px;
  padding: 20px;
  margin-top: 10px;
  margin-left: 425px;
  background-color: white;
  box-shadow: 0 2px 12px 0 rgb(0 0 0 / 30%);
}

button.el-button:focus {
  color: var(--el-button-text-color);
  background-color: var(--el-button-bg-color);
  border-color: var(--el-button-border-color);
}

button.el-button:hover {
  color: var(--el-button-hover-text-color);
  background-color: var(--el-button-hover-bg-color);
  border-color: var(--el-button-hover-border-color);
}

button.el-button:active {
  color: var(--el-button-active-text-color);
  background-color: var(--el-button-active-bg-color);
  border-color: var(--el-button-active-border-color);
}

button.el-button.is-link:focus {
  color: var(--el-button-focus-link-text-color);
}

button.el-button.is-link:hover {
  color: var(--el-button-hover-link-text-color);
}

button.el-button.is-link:active {
  color: var(--el-button-active-link-text-color);
}
</style>
