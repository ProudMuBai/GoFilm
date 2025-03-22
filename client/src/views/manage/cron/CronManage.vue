<template>
  <div class="container">
    <el-table
        :data="data.taskList" style="width: 100%" border size="default"
        :row-class-name="'cus-tr'" table-layout="auto">
      <el-table-column prop="id" label="任务ID">
        <template #default="scope">
          <el-tag disable-transitions>{{ scope.row.id }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="任务描述" />
      <el-table-column prop="model" align="center" label="任务类型">
        <template #default="scope">
          <el-tag disable-transitions>{{ scope.row.model == 0 ? '自动更新':scope.row.model == 0 ?'自定义任务':'采集重试'}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="state" align="center" label="是否启用">
        <template #default="scope">
          <el-switch  v-model="scope.row.state" @change="changeTaskState(scope.row.id, scope.row.state)" inline-prompt active-text="启用" inactive-text="禁用"/>
        </template>
      </el-table-column>
      <el-table-column prop="preV" align="center" label="上次执行时间">
        <template #default="scope">
          <el-tag type="success" disable-transitions>{{ scope.row.preV }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="next" align="center" label="下次执行时间">
        <template #default="scope">
          <el-tag type="warning" disable-transitions>{{ scope.row.next }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scope">
          <el-button type="primary" :icon="Edit" plain circle @click="openEditDialog(scope.row.id)" />
          <el-button type="danger" :icon="Delete" plain circle @click="delTask(scope.row.id)" />
        </template>
      </el-table-column>
    </el-table>
    <div class="cus_util">
      <el-button color="#9b49e7" :icon="Clock" @click="openAddDialog">创建定时任务</el-button>
      <!--<el-button color="#d942bf" :icon="Promotion">一键采集</el-button>-->
      <!--<el-button type="danger" :icon="BellFilled">ReZero</el-button>-->
    </div>
    <!--定时任务添加弹窗-->
    <el-dialog v-model="dialogV.addV" title="创建定时任务">
      <el-form :model="form.add">
        <el-form-item label="任务周期">
          <el-input v-model="form.add.spec" placeholder="定时任务Cron表达式 (例: [0 */20 * * * ?] 每20分钟执行一次)"/>
        </el-form-item>
        <el-form-item label="任务描述">
          <el-input v-model="form.add.remark" placeholder="定时任务描述信息"/>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-radio-group  fill="#9b49e7" v-model="form.add.model">
            <el-tooltip class="box-item" effect="dark" content="执行所有已启用站点的采集任务" placement="top">
              <el-radio :label="0">自动更新</el-radio>
            </el-tooltip>
            <el-tooltip class="box-item" effect="dark" content="只执行指定站点的采集任务" placement="top">
              <el-radio :label="1">自定义更新</el-radio>
            </el-tooltip>
             <el-tooltip class="box-item" effect="dark" content="失败采集重试处理" placement="top">
              <el-radio :label="2">采集重试</el-radio>
            </el-tooltip>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.add.model == 1"  label="资源绑定">
          <el-select v-model="form.add.ids" multiple collapse-tags collapse-tags-tooltip placeholder="Select" style="width: 240px">
            <el-option v-for="item in form.options" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.add.model != 2" label="采集时长">
          <el-tooltip class="box-item" effect="dark" content="采集最近x小时更新的影片,负数则默认采集所有资源" placement="top">
            <el-input-number v-model="form.add.time" :step="1" step-strictly />
          </el-tooltip>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-switch v-model="form.add.state" inline-prompt active-text="开启" inactive-text="禁用"/>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="addTask" >添加</el-button>
        <el-button @click="cancelDialog">取消</el-button>
      </span>
      </template>
    </el-dialog>
    <!--定时任务更新弹窗-->
    <el-dialog v-model="dialogV.editV" title="创建定时任务">
      <el-form :model="form.edit">
        <el-form-item label="任务标识">
          <el-tag type="success" disable-transitions>{{ form.edit.id }}</el-tag>
        </el-form-item>
        <el-form-item label="任务描述">
          <el-input v-model="form.edit.remark" placeholder="定时任务描述信息"/>
        </el-form-item>
        <el-form-item label="任务周期">
          <el-tag  disable-transitions>{{ form.edit.spec }}</el-tag>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-tag  disable-transitions>{{ form.edit.model == 0?'自动更新':form.edit.model == 1?'自定义更新':'采集重试' }}</el-tag>
        </el-form-item>
        <el-form-item v-if="form.edit.model == 1"  label="资源绑定">
          <el-select v-model="form.edit.ids" multiple collapse-tags collapse-tags-tooltip placeholder="Select" style="width: 240px">
            <el-option v-for="item in form.options" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.edit.model != 2" label="采集时长">
          <el-tooltip class="box-item" effect="dark" content="采集最近x小时更新的影片,负数则默认采集所有资源" placement="top">
            <el-input-number v-model="form.edit.time" :step="1" step-strictly />
          </el-tooltip>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-switch v-model="form.edit.state" inline-prompt active-text="开启" inactive-text="禁用"/>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="updateTask" >更新</el-button>
        <el-button @click="cancelDialog">取消</el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>


<script setup lang="ts">
import {Clock, Delete, Edit} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage} from "element-plus";

const data = reactive({
  taskList: []
})

const dialogV = reactive({
  addV: false,
  editV: false,
})

const form = reactive({
  add: { spec:'', remark: '', model: 1,ids:[], time:0, state: false},
  options: [],
  edit: { id: '', cid: '', spec:'', remark: '', model: 1, ids:[], time:0, state: false},
})

const addTask = ()=>{
  ApiPost(`/manage/cron/add`, form.add).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      cancelDialog()
      getTaskList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

const updateTask = ()=>{
  ApiPost(`/manage/cron/update`, {id: form.edit.id, ids: form.edit.ids, time: form.edit.time, state: form.edit.state, remark: form.edit.remark}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      cancelDialog()
      getTaskList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 关闭弹窗还原表单属性
const cancelDialog = ()=>{
  // 关闭对话框
  dialogV.addV = false
  dialogV.editV = false
  // 还原表单状态
  form.add = { spec:'', remark: '', model: 1,ids:[], time:0, state: false}
  form.edit = { id: '', cid: '', spec:'', remark: '', model: 1, ids:[], time:0, state: false}
}

const openAddDialog = ()=>{
  dialogV.addV = true
  getOptions()
}

// 删除定时任务
const delTask = (id:string)=>{
  ApiGet(`/manage/cron/del`, {id:id}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getTaskList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

const changeTaskState = (id:string,state:boolean)=>{
  ApiPost(`/manage/cron/change`, {id:id,state:state}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getTaskList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

const openEditDialog = (id:string)=>{
  dialogV.editV = true
  getOptions()
  ApiGet(`/manage/cron/find`,{id:id}).then((resp: any) => {
    if (resp.code === 0) {
      form.edit = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

const getOptions = ()=>{
  ApiGet(`/manage/collect/options`).then((resp: any) => {
    if (resp.code === 0) {
      form.options = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
const getTaskList = ()=>{
    ApiGet(`/manage/cron/list`).then((resp: any) => {
      if (resp.code === 0) {
        data.taskList = resp.data
      } else {
        data.taskList = []
        ElMessage.warning({message: resp.msg})
      }
    })
}
onMounted(()=>{
  getTaskList()
})

</script>


<style scoped>
.cus_util {
  display: flex;
  padding: 10px 8px;
  border-left: 2px solid #9b49e733;
  border-right: 2px solid #9b49e733;
  border-bottom: 2px solid #9b49e733;
  background: #ffffff;
  justify-content: end;
}

:deep(.el-input-number){
  --el-fill-color-light: #e163ff8f;
  border-radius: var(--el-border-radius-base);
  padding: 0 0 !important;
}
:deep(.el-input-number__decrease){
  left: 0;
  top: 0;
  bottom: 0;
}
:deep(.el-input-number__increase){
  right: 0;
  top: 0;
  bottom: 0;
}
:deep(.el-tag--info){
  --el-fill-color: #67d9e863;
}
</style>