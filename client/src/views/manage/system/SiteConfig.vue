<template>
  <div class="container">
    <h2 class="title">网站基础参数配置</h2>
    <div class="content">
      <el-form size="large" :model="c.site" label-width="120px">
        <el-form-item label="网站名称">
          <el-input v-model="c.site.siteName"/>
        </el-form-item>
        <el-form-item label="网站域名">
          <el-input v-model="c.site.domain"/>
        </el-form-item>
        <el-form-item label="网站Logo">
          <el-input v-model="c.site.logo"/>
        </el-form-item>
        <el-form-item label="搜索关键字">
          <el-input v-model="c.site.keyword"/>
        </el-form-item>
        <el-form-item label="网站描述">
          <el-input v-model="c.site.describe"/>
        </el-form-item>
        <el-form-item label="网站状态">
          <el-switch v-model="c.site.state" inline-prompt active-text="开启" inactive-text="关闭"/>
        </el-form-item>
        <el-form-item label="维护提示">
          <el-input v-model="c.site.hint"/>
        </el-form-item>
        <el-form-item>
          <el-button color="#9b49e7" @click="updateBasicConfig" >更新</el-button>
          <el-button @click="getBasicInfo" >重置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">


import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage} from "element-plus";

const c = reactive({
  site: {
    siteName: '',
    domain: '',
    logo: '',
    keyword: '',
    describe: '',
    state: true,
    hint: '',
  }
})

const updateBasicConfig = ()=>{
  ApiPost(`/manage/config/basic/update`, c.site).then((resp: any) => {
    if (resp.code === 0) {
      // console.log(window.location.hostname)
      ElMessage.success({message: resp.msg})
      getBasicInfo()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
  // 更新后重新从后端获取最新数据
}
const getBasicInfo = ()=>{
  ApiGet(`/manage/config/basic`).then((resp: any) => {
    if (resp.code === 0) {
      // console.log(window.location.hostname)
      c.site = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}


onMounted(() => {
  getBasicInfo()
})
</script>


<style scoped>
.container {
  background: #ffffff;
  padding: 20px 0;
  height: 100%;
}
.title {
  color: #2b333fb3;
  padding-bottom: 20px;
  border-bottom: 2px solid #00000005;
}

.content {
  width: 60%;
  margin: 36px auto;
}

:deep(.el-form-item__label) {
  color: #888888;
  font-size: 18px;

}

</style>