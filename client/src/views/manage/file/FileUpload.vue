<template>
  <div class="container">
    <div class="title_container">
      <h3>文件上传</h3>
    </div>
    <div class="content">
      <el-upload v-model:file-list="data.photoWall" action="#" list-type="picture-card"
          :on-remove="handleRemove" :http-request="customUpload">
        <template #file="{ file }">
          <div>
            <el-image class="el-upload-list__item-thumbnail" :src="file.link"  fit="cover" />
            <span class="el-upload-list__item-actions">
          <span class="el-upload-list__item-preview" @click="handlePictureCardPreview(file)">
            <el-icon><zoom-in /></el-icon>
          </span>
          <span class="el-upload-list__item-delete">
            <el-icon><Download /></el-icon>
          </span>
          <span class="el-upload-list__item-delete">
            <el-icon><Delete /></el-icon>
          </span>
        </span>
          </div>
        </template>
        <el-icon><Plus /></el-icon>
      </el-upload>

      <el-upload v-if="false" class="upload-demo" drag action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15" multiple>
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          删除文件 或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            jpg/png files with a size less than 500kb
          </div>
        </template>
      </el-upload>
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
  imgList:[""]
})
const customUpload = (options:any)=>{
  console.log(options)
  console.log(options.file)
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

const getPhotoPage = ()=>{
  ApiGet(`/manage/file/list`, ).then((resp: any) => {
    if (resp.code === 0) {
      data.photoWall = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  getPhotoPage()
})
const handleRemove: UploadProps['onRemove'] = (uploadFile, uploadFiles) => {
  console.log(uploadFile, uploadFiles)
}

// const handlePictureCardPreview: UploadProps['onPreview'] = (uploadFile) => {
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
  display: flex;
  justify-content: start;
}
</style>