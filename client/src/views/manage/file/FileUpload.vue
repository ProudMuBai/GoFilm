<template>
  <div class="container">
<!--    <div class="title_container">
      <h3>海报墙预览</h3>
    </div>-->
    <div class="content">
      <el-upload v-model:file-list="data.photoWall" action="#" list-type="picture-card"
          :http-request="customUpload">
        <template #file="{ file }">
            <el-image class="el-upload-list__item-thumbnail" style="width: 100%;height: 100%" :src="file.link"  fit="cover" />
            <span class="el-upload-list__item-actions">
          <span class="el-upload-list__item-preview" @click="handlePictureCardPreview(file)">
            <el-icon><zoom-in /></el-icon>
          </span>
          <span class="el-upload-list__item-delete" v-if="false">
            <el-icon><Download /></el-icon>
          </span>
          <span class="el-upload-list__item-delete" @click="delImage(file)" >
            <el-icon><Delete /></el-icon>
          </span>
        </span>
        </template>
        <el-icon><Plus /></el-icon>
      </el-upload>

      <div class="pagination">
        <el-pagination  background layout="prev, pager, next"
                       :total="data.page.total" v-model:page-size="data.page.pageSize"
                       v-model:current-page="data.page.current"
                       @change="getPhotoPage" hide-on-single-page/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import {Delete, Download, Plus, UploadFilled, ZoomIn} from "@element-plus/icons-vue";


import {ElMessage, UploadProps, UploadUserFile} from 'element-plus'
import {onMounted, reactive, ref} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {Preview} from "../../../components/Global/preview";

const data = reactive({
  photoWall: [],
  page: {current:1, pageSize: 39, pageNumber: 0, total: 0},
  imgList:[""]
})
const customUpload = (options:any)=>{
  let file = options.file
  let formData = new FormData();
  formData.append("file", file)
  ApiPost(`/manage/file/upload`, formData).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getPhotoPage()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 分页数据获取
const getPhotoPage = ()=>{
  ApiGet(`/manage/file/list`, {current: data.page.current} ).then((resp: any) => {
    if (resp.code === 0) {
      data.photoWall = resp.data.list
      data.page = resp.data.page
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  getPhotoPage()
})
const delImage = (file:any)=>{
  ApiGet(`/manage/file/del`, {id: file.ID} ).then((resp: any) => {
    if (resp.code === 0) {
      getPhotoPage()
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 图片放大预览
const handlePictureCardPreview  = (currentFile:any) => {
  let list = data.photoWall.map((item:any)=>{
    return item.link
  })
  Preview({list:list,currentLink: currentFile.link})
}
</script>

<style scoped>
.container {
  background: var(--bg-light);
}
.content {
  width: 100%;
  padding: 10px 0;
}
.title_container {
  margin: 10px 0 10px 0;
}
:deep(.el-upload-list--picture-card ) {
  padding: 10px 10px;
}
:deep(.el-upload-list__item ) {
  margin: 10px 10px!important;
}
:deep(.el-upload--picture-card){
  margin: 10px auto;
}

.pagination {
  padding: 20px 0;
  text-align: center;
}

:deep(.el-pagination) {
  width: 100% !important;
  justify-content: center;
  --el-color-primary: var(--paging-parmary-color);
}

</style>