<template>
  <div v-if="global.isMobile" class="container">
    <div class="card" v-for="h in data.historyList">
      <div class="card-left">
        <a class="card-link" :href="h.link" :style="{backgroundImage: `url(${h.picture})`}"></a>
      </div>
      <div class="card-right">
        <h5 class="card-title"> {{ h.name }}</h5>
        <div class="card-content">
          <p class="card-episode">{{ `已观看: ${h.progress}` }}</p>
          <p class="card-time "><b :class="`iconfont ${h.devices?'icon-mobile':'icon-pc1'}`"/>{{ h.time }}</p>
          <p class="card-episode">{{ h.episode }}</p>
        </div>
      </div>
      <a @click="cleanHistory(h.id)" class="iconfont icon-cancel1"/>
    </div>
    <el-empty v-if="data.historyList && data.historyList.length <= 0" style="padding: 10px 0;"  description="暂无观看记录"/>
  </div>
</template>

<script setup lang="ts">
import {inject, onMounted, reactive} from "vue";
import {COOKIE_KEY_MAP, cookieUtil} from "../../utils/cookie";

const data = reactive({
  historyList: [{}]
})

const global = inject<any>('global')

// 清除对应的历史记录
const cleanHistory = (id: number) => {
  console.log(id)
  // 获取cookie中的filmHistory
  let history = cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY) ? JSON.parse(cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY)) : null
  // 删除 cookie 中的对应记录
  delete history[id]
  // 将需改后的历史记录存储到cookie中
  cookieUtil.setCookie(COOKIE_KEY_MAP.FILM_HISTORY, JSON.stringify(history))
  // 更新记录数据
  getHistory()
}

const getHistory = () => {
  // 获取cookie中的filmHistory
  let historyMap = cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY) ? JSON.parse(cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY)) : null
  let arr = []
  if (historyMap) {
    for (let k in historyMap) {
      arr.push(historyMap[k])
    }
    arr.sort((item1, item2) => item2.timeStamp - item1.timeStamp)
  }
  data.historyList = arr
}

onMounted(() => {
  getHistory()
})

</script>

<style scoped>
.container {
  padding: 0 5px;
}

.card {
  width: 100%;
  max-height: 250px;
  display: flex;
  padding: 5px 5px;
  flex-direction: row;
  background: linear-gradient(#fff2, transparent);
  border: 1px solid rgba(255, 255, 255, 0.1);
  /*  border-bottom: 1px solid rgba(255, 255, 255, 0.1);*/
  margin: 5px auto;
  border-radius: 5px;
}

.card-left {
  flex-basis: 27%;
  min-width: 27%;
  display: flex;
}

.card-right {
  flex-basis: 68%;
  max-width: 68%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  text-align: left;
  padding-left: 5%;
  font-size: var(--text-font-content);
}

.card-link {
  width: 100%;
  padding-top: 125%;
  flex-grow: 1;
  border-radius: 3px;
  background-repeat: no-repeat;
  background-size: cover;

}

.card-title {
  max-width: 80%;
  margin-top: 10px;
  color: var(--text-title-color);
  font-size: var(--text-font-title-md);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-content p {
  margin-top: 3px;
  margin-bottom: 0;
  color: var(--text-content-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}


.iconfont {
  vertical-align: bottom;
  margin-right: 10px;
}


.icon-cancel1 {
  flex-basis: 5%;
  margin-right: 0;
  color: var(--text-content-color);
}



:deep(.el-empty) {
  --el-empty-fill-color-1: rgba(155, 73, 231, 0.72);
  --el-empty-fill-color-2: #67d9e891;
  --el-empty-fill-color-3: rgb(106 19 187 / 72%);
  --el-empty-fill-color-4: #67d9e8;
  --el-empty-fill-color-5: #5abcc9;
  --el-empty-fill-color-6: #9fb2d9;
  --el-empty-fill-color-7: #61989f;
  --el-empty-fill-color-8: #697dc5;
  --el-empty-fill-color-9: rgb(43 51 63 / 44%);
  margin-top: 20vh;
}


</style>