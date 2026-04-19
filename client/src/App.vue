<template>
  <el-config-provider :locale="zhCn">
    <div class="main dark" >
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

/*明暗主题背景设置*/
.dark {
  background: #16161a;

  .is-drawer .cus-drawer {
    background-color: rgba(0, 0, 0, 0.55)!important;
  }
  .el-drawer__header{
    border-bottom: 1px solid var(--border-color-highlight);
  }
  .el-drawer__close-btn{
    color: #ddddddbf;
  }
  .drawer-nav a {
    color: #ddd;
  }
  .drawer-nav a:before{
    color: #79bbff;
  }
}
.light{
  background: linear-gradient(45deg, #356697, rgb(105, 68, 140), rgb(151, 109, 133), rgb(92 104 149));

  .is-drawer .cus-drawer {
    background: linear-gradient(315deg, #fcb8e2ad 0.000%, #ffc8edcf 20.000%, #fbe8f9d4 40.000%, #e8f9ffde 60.000%, #d0ecffd6 80.000%, #b8cdfadb 100.000%);
  }
  .el-drawer__header{
    border-bottom: 1px solid var(--border-gray-color);
  }
  .el-drawer__close-btn{
    color: #2222223d;
  }
  .drawer-nav a {
    color: #222222c7;
  }
  .drawer-nav a:before{
    color: #3A9AFF;
  }
}


/*wrap侧边栏菜单样式*/
.drawer-nav{
  display: flex;
  flex-direction: column;
  width: 100%;
  padding: 0 10%;
}
.el-drawer__header{
  margin-left: 8px;
  margin-right: 8px;
  border-bottom: 1px solid var(--border-gray-color);
}
.el-drawer__header .drawer-title {
  margin: 20px 0;
  padding-left: 15px;
  text-align: start !important;
}
.el-drawer__close-btn{
  color: #2222223d;
  padding-right: 0;
}
.drawer-nav a {
  width: 100%;
  color: #22222296;
  text-align: start;
  padding: 6px 10px;
  margin-bottom: 15px;
  transition: transform 0.3s ease;
}
.drawer-nav a:before{
  font-size: 20px;
  color: #3A9AFF;
  margin-right: 6px;
}
.el-drawer__body .drawer-nav a:hover, .el-drawer__body .drawer-nav a:hover a:active{
  color: var(--active--color);
  transform: translateY(-4px);
}






</style>