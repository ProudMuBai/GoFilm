<template>
  <el-config-provider :locale="zhCn">
    <div class="main" >
      <router-view></router-view>
      <Util/>
    </div>
  </el-config-provider>
</template>

<script lang="ts" setup>
import Util from "./components/index/Util.vue";
import zhCn from "element-plus/dist/locale/zh-cn.min.js";
</script>

<style>
html, body, #app {
    width: 100%;
    margin: 0;
    padding: 0;
}

/*背景色切换*/
#app, .main {
    max-width: 100%;
    min-height: 100vh;
    color: rgb(221, 221, 221);
    background: #16161a;
}

* {
    box-sizing: border-box;
}
/*全局a标签默认样式去除*/
a {
    outline: none;
    text-decoration: none;
    color: rgba(255,255,255,0.38);
}

a:hover {
    color: #888888;
}
</style>