<template>
  <div class="container">

    <el-table
        :data="data.siteList" style="width: 100%" border size="default"
        :row-class-name="'cus-tr'" table-layout="auto">
      <el-table-column prop="name" label="资源名称"/>
      <el-table-column prop="resultModel" align="center" label="数据类型">
        <template #default="scope">
          <el-tag disable-transitions>{{ scope.row.resultModel == 0 ? 'JSON' : 'XML' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="collectType" align="center" label="资源类型">
        <template #default="scope">
          <el-tag disable-transitions>{{ scope.row.collectTypeText }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="uri" label="资源站">
        <template #default="scope">
          <el-link :href="scope.row.uri" target="_blank">{{ scope.row.uri }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="syncPictures" align="center" label="同步图片">
        <template #default="scope">
          <el-switch @change="changeSourceState(scope.row)" :disabled="scope.row.grade == 1" v-model="scope.row.syncPictures" inline-prompt active-text="开启" inactive-text="关闭"/>
        </template>
      </el-table-column>
      <el-table-column prop="state" align="center" label="是否启用">
        <template #default="scope">
          <el-switch @change="changeSourceState(scope.row)" v-model="scope.row.state" inline-prompt active-text="启用" inactive-text="禁用"/>
        </template>
      </el-table-column>
      <el-table-column prop="grade" align="center" label="站点权重">
        <template #default="scope">
          <el-tag disable-transitions :type="`${scope.row.grade == 0 ? 'success': 'info'}`">
            {{ scope.row.grade == 0 ? '采集主站' : '附属站点' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="grade" align="center" label="采集间隔">
        <template #default="scope">
          <el-tag disable-transitions type="success">
            {{scope.row.interval >0 ?`${scope.row.interval} ms`:`无限制` }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="采集方式">
        <template #default="scope">
          <el-select v-model="scope.row.cd" class="m-2" placeholder="Select" size="small">
            <el-option
                v-for="item in data.collectDuration"
                :key="item.time"
                :label="item.label"
                :value="item.time"
            />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scope">
          <el-button type="success" :icon="SwitchButton" plain circle @click="startTask(scope.row)" />
          <el-button type="primary" :icon="Edit" plain circle @click="openEditDialog(scope.row.id)" />
          <el-button type="danger" :icon="Delete" plain circle @click="delSourceSite(scope.row.id)" />
        </template>
      </el-table-column>
    </el-table>
    <div class="cus_util">
      <el-button color="#9b49e7" :icon="CirclePlus" @click="dialogV.addV = true">添加采集站</el-button>
      <el-button color="#d942bf" @click="openBatchCollect" :icon="Promotion">一键采集</el-button>
      <el-button type="danger" :icon="DeleteFilled" @click="dialogV.clear = true" >RemoveAll</el-button>
      <el-button type="primary" :icon="BellFilled" @click="dialogV.reCollect = true" >AutoCollect</el-button>
    </div>
    <!--站点添加弹窗-->
    <el-dialog v-model="dialogV.addV" title="添加采集站点">
      <el-form :model="form.add">
        <el-form-item label="资源名称">
          <el-input v-model="form.add.name" placeholder="自定义资源名称(禁用汉字)"/>
        </el-form-item>
        <el-form-item label="接口地址">
          <el-input v-model="form.add.uri" placeholder="资源采集链接,本站只采集综合资源或m3u8资源"/>
        </el-form-item>
        <el-form-item label="间隔时长">
          <el-tooltip class="box-item" effect="dark" content="单次采集请求的时间间隔, 单位/ms" placement="top">
            <el-input-number v-model="form.add.interval" :min="0" :step="100" step-strictly />
          </el-tooltip>
        </el-form-item>
        <el-form-item label="接口类型">
          <el-radio-group v-model="form.add.resultModel">
            <el-radio :label="0">JSON</el-radio>
            <el-radio disabled  :label="1">XML</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-radio-group fill="#9b49e7" v-model="form.add.collectType">
            <el-radio fill="#9b49e7" :label="0">视频</el-radio>
            <el-radio disabled  :label="1">文章</el-radio>
            <el-radio disabled  :label="2">演员</el-radio>
            <el-radio disabled  :label="3">角色</el-radio>
            <el-radio disabled  :label="4">网站</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="站点权重">
          <el-radio-group @change="restrict(0)" fill="#9b49e7" v-model="form.add.grade">
            <el-radio :label="0">主站点</el-radio>
            <el-radio :label="1">附属站点</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="图片同步">
          <el-switch v-model="form.add.syncPictures" @change="restrict(0)" inline-prompt active-text="开启"   inactive-text="关闭"/>
        </el-form-item>
        <el-form-item label="是否启用">
          <el-switch v-model="form.add.state" inline-prompt active-text="启用" inactive-text="禁用"/>
        </el-form-item>

      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#cf48be" @click="apiTest(form.add)" >测试</el-button>
        <el-button color="#9b49e7" @click="addSite" >添加</el-button>
        <el-button @click="dialogV.addV = false">取消</el-button>
      </span>
      </template>
    </el-dialog>
    <!--站点修改弹窗-->
    <el-dialog v-model="dialogV.editV" title="修改资源站信息">
      <el-form :model="form.edit">
        <el-form-item label="资源名称">
          <el-input v-model="form.edit.name" placeholder="自定义资源名称(禁用汉字)"/>
        </el-form-item>
        <el-form-item label="接口地址">
          <el-input v-model="form.edit.uri" placeholder="资源采集链接,本站只采集综合资源或m3u8资源"/>
        </el-form-item>
        <el-form-item label="间隔时长">
          <el-tooltip class="box-item" effect="dark" content="单次采集请求的时间间隔, 单位/ms" placement="top">
            <el-input-number v-model="form.edit.interval" :min="0" :step="100" step-strictly />
          </el-tooltip>
        </el-form-item>
        <el-form-item label="接口类型">
          <el-radio-group v-model="form.edit.resultModel">
            <el-radio :label="0">JSON</el-radio>
            <el-radio disabled  :label="1">XML</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-radio-group fill="#9b49e7" v-model="form.edit.collectType">
            <el-radio fill="#9b49e7" :label="0">视频</el-radio>
            <el-radio disabled  :label="1">文章</el-radio>
            <el-radio disabled  :label="2">演员</el-radio>
            <el-radio disabled  :label="3">角色</el-radio>
            <el-radio disabled  :label="4">网站</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="站点权重">
          <el-radio-group fill="#9b49e7" @change="restrict(1)" v-model="form.edit.grade">
            <el-radio :label="0">主站点</el-radio>
            <el-radio :label="1">附属站点</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="图片同步">
          <el-switch v-model="form.edit.syncPictures" @change="restrict(1)" inline-prompt active-text="开启"  inactive-text="关闭"/>
        </el-form-item>
        <el-form-item label="是否启用">
          <el-switch v-model="form.edit.state" inline-prompt active-text="启用" inactive-text="禁用"/>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#cf48be" @click="apiTest(form.edit)" >测试</el-button>
        <el-button color="#9b49e7" @click="updateSite(form.edit)" >更新</el-button>
        <el-button @click="dialogV.editV = false">取消</el-button>
      </span>
      </template>
    </el-dialog>
    <!--一键执行功能弹窗-->
    <el-dialog v-model="dialogV.batchV" width="450px" title="多资源站一键采集">
      <el-form :model="form.batch">
        <el-form-item label="执行站点">
          <el-select v-model="form.batch.ids" multiple collapse-tags collapse-tags-tooltip placeholder="Select" style="width: 240px">
            <el-option v-for="item in form.options" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="采集时长">
          <el-tooltip class="box-item" effect="dark" content="采集最近x小时更新的影片,负数则默认采集所有资源" placement="top">
            <el-input-number v-model="form.batch.time" :step="1" step-strictly />
          </el-tooltip>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="startBatchCollect" >确认执行</el-button>
        <el-button @click="cancelDialog">取消</el-button>
      </span>
      </template>
    </el-dialog>

    <!--影片删除提示弹窗-->
    <el-dialog v-model="dialogV.clear" title="是否清除所有影视数据 ?" width="500">
      <el-form :model="form">
        <el-form-item label="确认密码" >
          <el-input v-model="data.password" type="password" placeholder="请输入账户密码并开确认执行" autocomplete="off" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogV.clear = false">取消</el-button>
          <el-button type="primary" @click="clearFilms">确认执行</el-button>
        </div>
      </template>
    </el-dialog>

    <!--Re0 从零开始的自动全量采集-->
    <el-dialog v-model="dialogV.reCollect" title="是否清除影片数据并重新采集 ?" width="500">
      <el-form :model="form">
        <el-form-item label="确认密码" >
          <el-input v-model="data.password" type="password" placeholder="请输入账户密码并开确认执行" autocomplete="off" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogV.reCollect = false">取消</el-button>
          <el-button type="primary" @click="reCollectFilm">确认执行</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">


import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage} from "element-plus";
import {Delete, Edit, SwitchButton, CirclePlus, Promotion, BellFilled, DeleteFilled} from "@element-plus/icons-vue";

const data = reactive({
  siteList: [],
  collectDuration: [
    {time: 24, label: '采集今日'},
    {time: 24 * 7, label: '采集本周'},
    {time: -1, label: '采集全部'},
  ],
  password: '',
})

const dialogV = reactive({
  addV: false,
  editV: false,
  batchV:false,
  clear: false,
  reCollect: false,
})

interface FilmSource {
  id: string
  name: string
  uri: string
  resultModel: number
  grade: number
  collectType: number
  syncPictures: boolean
  state: boolean
}

const form = reactive({
  add: {name: '', uri: '', resultModel: 0, grade: 1, collectType: 0, syncPictures: false, state: false, interval: 0},
  edit: {id:'', name: '', uri: '', resultModel: 0, grade: 1, collectType: 0, syncPictures: false, state: false,interval:0},
  batch: {ids:[],time: 0},
  options:[]

})

const openBatchCollect = ()=>{
  dialogV.batchV = true
  ApiGet(`/manage/collect/options`, ).then((resp: any) => {
    if (resp.code === 0) {
      form.options = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

const startBatchCollect = ()=>{
  ApiPost(`/manage/spider/start`, {ids: form.batch.ids, time: form.batch.time, batch: true}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      cancelDialog()
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

//  开启采集
const startTask = (row:any)=>{
  ApiPost(`/manage/spider/start`, {id:row.id, time: row.cd, batch: false}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 弹窗图片同步开关限制
const restrict = (t:number)=>{
  // t 弹窗类型 0 - add | 1 - edit
  switch (t){
    case 0:
      // 只有 主站点才能开启图片同步, 否则自动为false
      form.add.syncPictures = (form.add.syncPictures)&&(form.add.grade == 0)
    break
    case 1:
      // 只有 主站点才能开启图片同步, 否则自动为false
      form.edit.syncPictures = (form.edit.syncPictures)&&(form.edit.grade == 0)
    break
  }

}


// 添加采集资源站
const addSite = ()=>{
  ApiPost(`/manage/collect/add`, form.add).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      cancelDialog()
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
// 测试添加的采集接口是否可用
const apiTest = (params:any)=>{
  ApiPost(`/manage/collect/test`, params).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 编辑按钮事件, 打开修改对话框
const openEditDialog = (id:string)=>{
  // 从后台获取采集站信息
  ApiGet(`/manage/collect/find`, {id:id}).then((resp: any) => {
    if (resp.code === 0 ) {
      form.edit = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
  dialogV.editV = true
}

// switch 开关
const changeSourceState = (s:any)=>{
  ApiPost(`/manage/collect/change`, {id:s.id, state: s.state, syncPictures: s.syncPictures}).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

//更新资源站点信息
const updateSite = (params:FilmSource)=>{
  ApiPost(`/manage/collect/update`, params).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      dialogV.editV = false
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 删除采集资源站
const delSourceSite = (id:string) =>{
  ApiGet(`/manage/collect/del`, {id:id}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getCollectList()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
const cancelDialog = ()=>{
  // 关闭对话框
  dialogV.addV = false
  dialogV.editV = false
  dialogV.batchV = false
  // 还原表单状态
  form.add = {name: '', uri: '', resultModel: 0, grade: 1, collectType: 0, syncPictures: false, state: false, interval: 0}
}

// 获取采集列表信息
const getCollectList = ()=>{
  ApiGet(`/manage/collect/list`).then((resp: any) => {
    if (resp.code === 0) {
      data.siteList = resp.data.map((item: any) => {
        switch (item.collectType) {
          case 0:
            item.collectTypeText = "视频"
            break
          case 1:
            item.collectTypeText = "文章"
            break
          case 2:
            item.collectTypeText = "演员"
            break
          case 3:
            item.collectTypeText = "角色"
            break
          case 4:
            item.collectTypeText = "网站"
            break
        }
        item.cd = 24
        return item
      })
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 清除影片信息
const clearFilms = ()=>{
  if (data.password.length <= 0) {
    ElMessage.error({message: '操作失败, 密钥信息缺失'})
    return
  }
  ApiGet(`/manage/spider/clear`, {password: data.password}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
    dialogV.clear = false
    data.password = ''
  })
}

// 从零开始重新采集
const reCollectFilm = ()=>{
  if (data.password.length <= 0) {
    ElMessage.error({message: '操作失败, 密钥信息缺失'})
    return
  }
  ApiGet(`/manage/spider/zero`, {password: data.password}).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
    dialogV.reCollect = false
    data.password = ''
  })
}

onMounted(() => {
  getCollectList()
})

</script>

<style scoped>




</style>