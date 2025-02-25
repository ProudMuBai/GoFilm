<template>
    <div>
      <!--:collapse-transition="false" 关闭过渡动画-->
      <el-menu default-active="2" class="side-nav" router :collapse="collapse.collapse.value" >
        <el-menu-item  index="" @click="toIndex" >
          <el-avatar class="logo" :size="30" :src="data.site.logo.toString()" alt="GoFilm"/>
          <template #title>
            <b class="site_name">{{ data.site.siteName }}</b>
          </template>
        </el-menu-item>
        <el-sub-menu index="/manage/index">
          <template #title>
            <el-icon><HomeFilled /></el-icon>
            <span>网站管理</span>
          </template>
          <el-menu-item index="/manage/system/webSite">站点管理</el-menu-item>
          <el-menu-item index="/manage/system/banners">海报管理</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/manage/collect">
          <template #title>
            <el-icon><MagicStick /></el-icon>
            <span>采集管理</span>
          </template>
          <el-menu-item index="/manage/collect/index">影视采集</el-menu-item>
          <el-menu-item index="/manage/collect/record">失效记录</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/manage/cron">
          <template #title>
            <el-icon><Timer /></el-icon>
            <span>定时任务</span>
          </template>
          <el-menu-item index="/manage/cron/index">任务管理</el-menu-item>
        </el-sub-menu>
        <!--<el-menu-item index="/manage/category/index">-->
        <!--  <el-icon><Menu /></el-icon>-->
        <!--  <template #title>分类管理</template>-->
        <!--</el-menu-item>-->
        <el-sub-menu index="/manage/film">
          <template #title>
            <el-icon><Film /></el-icon>
            <span>影片管理</span>
          </template>
          <el-menu-item index="/manage/film/class">影视分类</el-menu-item>
          <el-menu-item index="/manage/film">影视信息</el-menu-item>
          <el-menu-item index="/manage/film/add">影片添加</el-menu-item>
          <el-menu-item index="/manage/film/detail">视频详情</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/manage/file">
          <template #title>
            <el-icon><FolderOpened /></el-icon>
            <span>文件管理</span>
          </template>
          <el-menu-item index="/manage/file/upload">文件上传</el-menu-item>
          <el-menu-item index="/manage/file/gallery">图库管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </div>
</template>

<script setup lang="ts">
import {Menu , HomeFilled, MagicStick, Timer, Film, FolderOpened} from '@element-plus/icons-vue'
import {inject, onMounted, reactive} from "vue";
import {ApiGet} from "../../utils/request";
import {ElMessage} from "element-plus";

// 菜单栏展开状态
const collapse = inject('collapse')

const data = reactive({
  site: {siteName: String, logo:String},
})

// 网站logo点击事件
const toIndex = ()=>{
  window.open('/index')
}

// 初始化网站图标名称等组件数据
const getSiteInfo = ()=>{
  ApiGet(`/manage/config/basic`).then((resp: any) => {
    if (resp.code == 0) {
      data.site = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  getSiteInfo()
})

</script>



<style scoped>
:deep(.el-menu){
  --el-menu-bg-color: #191a23!important;
  --el-menu-hover-bg-color: rgb(20, 21, 28);
  --el-menu-level: 0;
  --el-menu-text-color: #fff;
  --el-menu-active-color: skyblue;
}
.side-nav {
  padding: 20px 0;
  height: 100vh;
  border-right: none;
}


.side_head{
  display: flex;
  font-size: 16px;
}
.logo{
  margin-right: 10px;
  min-width: 30px;
}
.site_name {
  color: transparent;
  font-size: 20px;
  font-style: italic;
  -webkit-background-clip: text!important;
  background-clip: text;
  background: linear-gradient(118deg, #e91a90, #c965b3, #988cd7, #00acfd);
}


</style>