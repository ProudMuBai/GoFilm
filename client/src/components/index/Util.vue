<template>
  <div class="util">
    <el-collapse-transition>
      <div v-show="control.show">
        <a href="javascript:;" @click="changeStyle('top')">
          <el-icon>
            <ArrowUp/>
          </el-icon>
        </a>
        <a v-if="control.darkTheme" href="javascript:;" @click="changeStyle('light')">
          <el-icon>
            <Sunny/>
          </el-icon>
        </a>
        <a v-if="!control.darkTheme" href="javascript:;" @click="changeStyle('dark')">
          <el-icon>
            <Moon/>
          </el-icon>
        </a>
        <a href="/custom/player" @click="changeStyle('dark')">
          <el-icon>
            <VideoCamera />
          </el-icon>
        </a>
      </div>
    </el-collapse-transition>
    <a href="javascript:;" @click="changeStyle('more')" class="more">
      <el-icon>
        <MoreFilled/>
      </el-icon>
    </a>
    <!--<CustomDialog />-->
  </div>
</template>

<script setup lang="ts">
import {ArrowUp, Sunny, Moon, MoreFilled, VideoCamera} from '@element-plus/icons-vue'
import {onMounted, reactive} from "vue";
import CustomDialog from "../Popup/CustomDialog.vue";

const control = reactive({
  show: false,
  darkTheme: true,

})

onMounted(() => {
  changeStyle(localStorage.getItem("theme") + '')
})
const changeStyle = (type: string) => {
  switch (type) {
    case 'top':
      let top = document.documentElement.scrollTop
      if (top > 0) {
        // 创建定时器，平滑滚动
        const interval = setInterval(() => {
          document.documentElement.scrollTop -= 10;
          if (document.documentElement.scrollTop === 0) {
            clearInterval(interval);
          }
        }, 5);
      }
      break
    case 'light':
      control.darkTheme = false
      localStorage.setItem("theme", 'light')
      document.getElementsByClassName('main')[0].style.background = `linear-gradient(45deg, #356697, rgb(105, 68, 140), rgb(151, 109, 133), rgb(92 104 149))`
      break
    case  'dark':
      control.darkTheme = true
      localStorage.setItem("theme", 'dark')
      document.getElementsByClassName('main')[0].style.background = `#16161a`
      break
    case 'more':
      control.show = !control.show
      break
  }
}
</script>

<style scoped>
/*窗口工具栏设置*/
.util {
  position: fixed;
  right: 10px;
  bottom: 15%;
  width: 35px;
  z-index: 20;
}

.util a {
  display: block;
  width: 100%;
  margin-bottom: 3px;
  height: 35px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.35);
}

.util a:hover {
  background: #d329a4;
}

:deep(.el-icon) {
  font-size: 18px;
  height: 100%;
  color: #ffffff;
}

.more {
  background: rgb(238, 150, 0) !important;
}

</style>