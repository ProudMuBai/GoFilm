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
          <p class="card-time "><b :class="`iconfont ${h.devices?'icon-mobile':'icon-pc1'}`" />{{ h.time }}</p>
          <p class="card-episode">{{ h.episode }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {inject, onMounted, reactive} from "vue";
import {COOKIE_KEY_MAP, cookieUtil} from "../../utils/cookie";

const data = reactive({
  historyList: [{}]
})

const global = inject<any>('global')

onMounted(() => {
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
  console.log(arr)
})

</script>

<style scoped>

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
  display: flex;
}

.card-right {
  flex-basis: 73%;
  max-width: 73%;
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


</style>