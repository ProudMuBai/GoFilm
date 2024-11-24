<template>
  <div class="container">
    <div class="card" v-for="item in data.historyList">
      <div class="card-left">
        <a :href="item.link"></a>
      </div>
      <div class="card-right">
        <h5>{{item.name}}</h5>
        <span>{{item.episode}}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive} from "vue";
import {COOKIE_KEY_MAP, cookieUtil} from "../../utils/cookie";

const data = reactive({
  historyList: [{}]
})

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
  border: 1px solid red;
  display: flex;
  flex-direction: row;
}

.card-left {
  flex-basis: 20%;
  border: 1px solid greenyellow;
}

.card-right {
  flex-basis: 80%;
  border: 1px solid deepskyblue;
}

</style>