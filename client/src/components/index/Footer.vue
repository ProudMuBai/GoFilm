<template>
  <div class="custom-footer">
    <!--Mobile  bottom Btn Group-->
    <div v-if="global.isMobile" class="nav-bottom">
      <a href="/"><b class="iconfont icon-ziyuan"/>首页</a>
      <a href="javascript:void(0)"><b class="iconfont icon-gengxin"/>更新</a>
      <a href="/history"><b class="iconfont icon-lishi"/>历史</a>
      <a href="javascript:void(0)"><b class="iconfont icon-shoucang2"/>收藏</a>
      <a href="javascript:void(0)"><b class="iconfont icon-ziyuan-copy"/>我的</a>
    </div>
    <p>本站所有内容均来自互联网分享站点所提供的公开引用资源，未提供资源上传、存储服务。</p>
  </div>
</template>

<script lang="ts" setup>

// 获取全局状态对象
import {inject} from "vue";

const global = inject('global')

</script>

<style scoped>


.custom-footer {
  text-align: center;
  width: 100%;
  height: 30px;
  margin-top: 25px;
}

.custom-footer p {
  line-height: 15px;
  font-size: 15px;
  color: #888888;
}

/*底部导航*/
.nav-bottom {
  display: flex;
  justify-content: space-between;
  position: fixed;
  width: 100%;
  bottom: 0;
  background: rgba(0, 0, 0, 0.65);
}

.nav-bottom > a {
  display: flex;
  flex-direction: column;
  color: var(--content-text-color);
  flex-basis: 20%;
  font-size: 10px;
}

.iconfont {
  font-size: 24px;
}

.icon-ziyuan {
  color: #B1AFFF;
}

.icon-gengxin {
  color: #5AB2FF;
}

.icon-lishi {
  color: #FFBE98;
}

.icon-shoucang2 {
  color: #6AD4DD;
}

.icon-ziyuan-copy {
  color: #E89ABEFF;
}

</style>