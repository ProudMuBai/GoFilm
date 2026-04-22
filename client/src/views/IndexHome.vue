<template>
  <el-container>
    <el-header v-show="!data.isNavHidden">
      <Header/>
      <!--<NewHeader />-->
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
    <el-footer v-show="!data.isNavHidden">
      <Footer/>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import Header from "../components/index/Header.vue";
import Footer from "../components/index/Footer.vue";
import {onMounted, onUnmounted, provide, reactive} from "vue";
// import NewHeader from "../components/index/NewHeader.vue";

// 页面数据
const data = reactive({
  lastScrollTop: 0,
  isNavHidden: false,
})

// 在全局注入一个当前是pc还是wrap的状态
const userAgent = navigator.userAgent.toLowerCase()
let isMobile = /Mobile|Tablet|Android|iPhone|iPad|iPod|BlackBerry|webOS|Windows Phone|SymbianOS|IEMobile|Opera Mini/i.test(userAgent)
// 传递一个全局状态对象
provide('global', {isMobile: isMobile})

//
// 导航栏显示状态控制
const handleScroll = ()=>{
  // 获取当前滚动条距离顶部的距离
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop;

  // 只有当滚动超过一定距离（例如100px）后才开始判断，避免在顶部微小滚动时触发
  if (scrollTop > 200) {
    data.isNavHidden = scrollTop > data.lastScrollTop;
  } else {
    // 滚动到顶部附近时，始终显示导航栏
    data.isNavHidden = false;
  }

  // 更新上一次的滚动位置
  data.lastScrollTop = scrollTop <= 0 ? 0 : scrollTop; // 防止负值
}

// 页面挂载完成时添加监听事件
onMounted(()=>{
  // 添加页面滚动监听, 用于导航栏隐藏状态
  window.addEventListener('scroll', handleScroll);
})
// 卸载时取消页面滚动监听
onUnmounted(()=>{
  window.removeEventListener('scroll', handleScroll);
})

</script>


<style scoped>

:deep(.el-main) {
  padding-top: 70px !important;
  padding-bottom: 30px !important;
  min-height: 85vh;
}

:deep(.el-header) {
  padding: 0 !important;
  position: fixed !important;
  width: 100% !important;
  min-height: 60px;
  transform: translateZ(0);
  z-index: 1000;
  background-color: rgba(0, 0, 0, 0.45);
  top: 0;
}

:deep(.el-footer) {
  --el-footer-padding: 0 0;
}


@media (min-width: 768px) {
  .el-main {
    margin: 0 auto;
    padding: 100px 0;
    /*padding-top: 100px!important;*/
  }
}

@media (max-width: 768px) {
  .el-main {
    /*margin: 0 auto;*/
    padding: 55px 0 !important;
    /*padding-top: 100px!important;*/
  }

  :deep(.el-header) {
    height: 40px !important;
    min-height: 40px !important;
  }
}

:deep(.el-menu--horizontal) {
  border-bottom: 1px solid rgb(46, 46, 46);
}

/*@media (min-width: 768px){ //>=768的设备 }*/
/*@media (min-width: 992px){ //>=992的设备 }*/
/*@media (min-width: 1200){ //>=1200的设备 }*/

@media (min-width: 1024px) {
  .el-main {
    width: 1023px
  }
}

@media (min-width: 990px) {
  .el-main {
    width: 970px
  }
}

@media (min-width: 1200px) {
  .el-main {
    width: 1180px
  }
}

@media (min-width: 1400px) {
  .el-main {
    width: 1400px
  }
}

@media (min-width: 1560px) {
  .el-main {
    width: 1500px
  }
}

</style>


