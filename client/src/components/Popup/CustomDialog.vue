<template>
  <div>
    <el-dialog v-model="remind.visible"
               :width="remind.width"
               align-center
               @close="remindClose"
               :close-on-click-modal="false"
               draggable>
      <template #header>
        <span class="remind-title">
          温馨提示
        </span>
      </template>
      <span class="remind-content">{{remind.msg}}</span>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive} from "vue";

const remind = reactive({
  visible: false,
  width: "50%",
  msg: "您已经持续观看很长时间了,请休息一下哦!!!",
  interval: 1000 * 60 * 45,
})
const remindClose = () => {
  setTimeout(() => {
    remind.visible = true
  }, remind.interval)
}

onMounted(() => {
  const userAgent = navigator.userAgent.toLowerCase();
  let isMobile = /mobile|android|iphone|ipad|phone/i.test(userAgent)
  remind.width = isMobile ? "90%" : "30%"
  /*开启定时器, 每隔指定时间出现一次弹窗*/
  setTimeout(() => {
    remind.visible = true
  }, remind.interval)
})
</script>

<style scoped>
:deep(.el-overlay) {
  /*
  --el-overlay-color-lighter: rgba(255,255,255,0.35);
  */
}

/*:deep(.el-dialog){
  border-radius: 8px;
  background-image: linear-gradient( 135deg, #81FFEF 10%, #F067B4 100%);
}
:deep(.el-dialog__header){
  border-bottom: 1px solid rgba(0,0,0,0.1);
  margin-right: 0;
}
:deep(.el-dialog__body){
  padding: 45px!important;
}*/
.remind-title {
  font-size: 26px;
  color: rgb(172 56 191);

}

.remind-content {
  font-size: 18px;
  color: rgb(134 68 135);
}
</style>