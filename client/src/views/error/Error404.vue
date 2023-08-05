<template>
  <div class="container">

    <el-empty v-if="false" style="height: 100%; line-height: 100%; font-size: 20px" :image-size="400"
              :image="notFoundImg">
      <template #description>
        <p style="font-size: 32px;width: 100%;color: #a574b7;margin-bottom: 20px">你好像走错地方了哦!!!</p>
        <button @click="handleError('home')">主页</button>
        <button @click="handleError('back')">返回</button>
      </template>
    </el-empty>


    <img :src="notFoundImg" alt="404" class="notfound-img"/>
    <p class="warning-txt">你好像走错地方了哦!!!</p>
    <div class="btn-group">
      <button @click="handleError('home')">主页</button>
      <button @click="handleError('back')">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import notFoundImg from '../../assets/image/404.png'

const router = useRouter()
const handleError = (t: string) => {
  switch (t) {
    case 'home':
      router.push('/')
      break
    case 'back':
      router.go(-1)
      break
  }
}
</script>

<style scoped>
@media (min-width: 768px) {
  .container {
    padding: 50px 0!important;
  }
  .notfound-img {
    margin: 0 auto;
    max-height: 76vh;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 80px 0;
  }

  .notfound-img {
    max-width: 96%;
    margin: 0 auto;
  }
}

.notfound-img {
  max-width: 90%;
  margin: 0 auto;
}

.container {
  width: 100vw;
  height: 100vh;
}

.warning-txt {
  font-size: 32px;
  color: #79bbff;
  margin: 5px;
}

.btn-group {
  margin: 0 auto;
}

button {
  margin: 5px 15px;
  background: rgb(103, 217, 232);
  color: #ffffff;
  width: 120px;
  border-radius: 36px;
}

button:hover {
  border: none;
}
</style>