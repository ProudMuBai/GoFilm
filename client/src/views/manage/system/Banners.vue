<template>
  <h2 style="color: #8e48b4">首页横幅管理界面</h2>
  <el-table
      :data="data.banners" style="width: 100%" border size="default"
      :row-class-name="'cus-tr'" table-layout="auto">
    <el-table-column prop="name" label="影片名称"/>
    <el-table-column prop="collectType" align="center" label="影片类型">
      <template #default="scope">
        <el-tag type="warning">{{ scope.row.cName }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="上映年份">
      <template #default="scope">
        <el-tag type="warning">{{ scope.row.year }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="影片海报">
      <template #default="scope">
        <el-image style="width: 120px; height: 80px" :src="scope.row.poster" :preview-src-list="[scope.row.poster]" preview-teleported fit="contain" />
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="影片封面">
      <template #default="scope">
        <el-image style="width: 60px; height: 80px" :src="scope.row.picture" :preview-src-list="[scope.row.picture]" preview-teleported fit="cover" />
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="排序">
      <template #default="scope">
        <el-tag disable-transitions>{{ scope.row.sort }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="resultModel" align="center" label="连载状态">
      <template #default="scope">
        <el-tag v-if="(scope.row.remarks+'').search('更新') == -1" type="success" >{{ scope.row.remark}}</el-tag>
        <el-tag v-else type="primary" >{{ scope.row.remark}}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" align="center">
      <template #default="scope">
        <el-button type="success" :icon="SwitchButton" plain circle @click="" />
        <el-button type="primary" :icon="Edit" plain circle @click="" />
        <el-button type="danger" :icon="Delete" plain circle @click="" />
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import {Delete, Edit, SwitchButton} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet} from "../../../utils/request";
import {ElMessage} from "element-plus";

const data = reactive({
  banners: [],

})

const getBanners = ()=>{
  ApiGet(`/manage/banner/list`).then((resp:any)=>{
    if(resp.code === 0){
      data.banners = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  // 初始化banners数据
  getBanners()
})


</script>

<style scoped>

</style>