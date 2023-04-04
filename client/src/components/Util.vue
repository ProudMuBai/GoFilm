<template>
    <div class="util">
        <el-collapse-transition>
            <div v-show="control.show">
                <a href="javascript:;" @click="changeStyle('top')">
                    <el-icon>
                        <ArrowUp />
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
            </div>
        </el-collapse-transition>
        <a href="javascript:;" @click="changeStyle('more')" class="more">
            <el-icon>
                <MoreFilled/>
            </el-icon>
        </a>
    </div>
</template>

<script setup lang="ts">
import {ArrowUp,Sunny, Moon, MoreFilled} from '@element-plus/icons-vue'
import {onBeforeMount, onMounted, onUnmounted, reactive} from "vue";

const control = reactive({
    show: false,
    darkTheme: true,
})

//
onMounted(()=>{

    changeStyle(localStorage.getItem("theme")+'')
})


const changeStyle = (type:string)=>{
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
          document.getElementsByClassName('main')[0].style.background = `radial-gradient(circle, #C147E9, #810CA8, #2D033B)`
          break
      case  'dark':
          control.darkTheme = true
          localStorage.setItem("theme", 'dark')
          document.getElementsByClassName('main')[0].style.background = `rgb(34,34,34)`
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
}
.util a {
    display: block;
    width: 100%;
    margin-bottom: 3px;
    height: 35px;
    border-radius: 50%;
    background: rgba(0,0,0,0.35);
}
.util a:hover{
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