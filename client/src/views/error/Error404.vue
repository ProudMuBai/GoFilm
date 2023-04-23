<template>
  <div class="container">
      <el-empty style="height: 100%; line-height: 100%; font-size: 20px" :image-size="400"
                :image="notFoundImg">
          <template #description>
              <p style="font-size: 32px;width: 100%;color: #a574b7;margin-bottom: 20px">你好像走错地方了哦!!!</p>
              <button @click="handleError('home')" >主页</button>
              <button @click="handleError('back')">返回</button>
          </template>
      </el-empty>
  </div>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import notFoundImg from '../../assets/image/404.png'

const router = useRouter()
const handleError = (t:string) =>{
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
.container {
    height: 100vh;
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