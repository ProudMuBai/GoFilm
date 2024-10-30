<template>
  <div class="container">
    <el-table
        :data="data.classTree" style="width: 100%" border size="default"
        table-layout="auto" max-height="calc(90vh - 20px)"
        row-key="id"
        :row-class-name="'cus-tr'" >
      <el-table-column prop="name" label="分类名称">
        <template #default="scope">
          <el-tag :type="scope.row.pid==0?'success':'warning'" disable-transitions>{{ scope.row.name}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="show" align="center" label="是否展示">
        <template #default="scope">
          <el-switch v-if="scope.row.pid == 0"  v-model="scope.row.show" inline-prompt active-text="展示" inactive-text="隐藏" @change="changeClassState(scope.row.id, scope.row.show)" />
          <el-switch v-else v-model="scope.row.show" inline-prompt active-text="屏蔽" inactive-text="恢复" @change="changeClassState(scope.row.id, scope.row.show)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scope">
          <el-button type="primary" :icon="Edit" @click="openEditDialog(scope.row.id)" plain circle />
          <el-button type="danger" :icon="Delete" @click="delClass(scope.row.id)" plain circle/>
        </template>
      </el-table-column>
    </el-table>
    <!--功能按钮-->
    <div class="cus_util">
      <el-button color="#9b49e7" :icon="RefreshLeft" @click="resetFilmClass">重置分类信息</el-button>
    </div>
    <!--影片分类信息修改弹窗-->
    <el-dialog v-model="dialog.editV" @close="cancelDialog" width="480px" title="更新分类信息">
      <el-form :model="dialog.editForm">
        <el-form-item label="分类名称">
          <el-input v-model="dialog.editForm.name" placeholder="分类名称,用于首页导航展示"/>
        </el-form-item>
        <el-form-item label="分类层级">
          <el-tag :type="dialog.editForm.pid==0?'success':'warning'"  disable-transitions>{{ dialog.editForm.pid == 0 ? '一级分类':'二级分类' }}</el-tag>
        </el-form-item>
        <el-form-item label="是否展示">
          <el-switch v-model="dialog.editForm.show" inline-prompt active-text="展示" inactive-text="隐藏"/>
        </el-form-item>
        <el-form-item class="class_sub" label="拓展分类" v-if="dialog.editForm.children">
          <el-tag class="class_sub_tag" v-for="c in dialog.editForm.children" type="warning" disable-transitions>{{ c.name }}</el-tag>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="updateClass" >更新</el-button>
        <el-button @click="dialog.editV = false">取消</el-button>
      </span>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import {Delete, Edit, RefreshLeft, Warning,} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage,ElMessageBox} from "element-plus";

// table数据
const data = reactive({
  classTree: []
})


// dialog 弹窗数据
const dialog = reactive({
  editV: false,
  editForm: {id: -99, pid: -99, name:'', show: true, children:[]},
})

// 删除分类信息
const delClass = (id:number)=>{
  ApiGet(`/manage/film/class/del`, {id: id}).then((resp: any) => {
    if (resp.code === 0) {
      getFilmClassTree()
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 打开修改弹窗
const openEditDialog = (id:number)=>{
  dialog.editV = true
  ApiGet(`/manage/film/class/find`,{id:id}).then((resp: any) => {
    if (resp.code === 0) {
      dialog.editForm = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
// 更新影片分类信息
const updateClass = ()=>{
  let { id, name, show } = dialog.editForm
  ApiPost(`/manage/film/class/update`,{ id: id, name:name, show:show }).then((resp: any) => {
    if (resp.code === 0) {
      // 更新成功后关闭弹窗, 重新获取最新的分类表格信息
      dialog.editV = false
      getFilmClassTree()
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
// 影片是否展示Switch按钮
const changeClassState = (id:string, show: number)=>{
  ApiPost(`/manage/film/class/update`,{ id: id, show:show }).then((resp: any) => {
    if (resp.code === 0) {
      // 更新成功后关闭弹窗, 重新获取最新的分类表格信息
      dialog.editV = false
      getFilmClassTree()
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 对话框关闭时重置表单数据为初始值
const cancelDialog = ()=>{
  dialog.editForm = {id: -99, pid: -99, name:'', show: true, children:[]}
}

// 重置分类信息
const resetFilmClass = ()=>{
  ApiGet(`/manage/spider/class/cover`).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 获取影片分类树信息
const getFilmClassTree = ()=>{
  ApiGet(`/manage/film/class/tree`).then((resp: any) => {
    if (resp.code === 0) {
      data.classTree = resp.data.children
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  getFilmClassTree()
})


</script>



<style scoped>
:deep(.el-table){
  --el-table-row-hover-bg-color: #9b49e71a;
}
.class_sub {
  display: flex;
  justify-content: start;
}
.class_sub_tag{
  width: calc(20% - 8px);
  margin: 5px 4px;
}

.cus_util {
  display: flex;
  padding: 10px 8px;
  border-left: 2px solid #9b49e733;
  border-right: 2px solid #9b49e733;
  border-bottom: 2px solid #9b49e733;
  background: #ffffff;
  justify-content: end;
}
</style>