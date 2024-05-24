<template>
  <div class="container">
    <!-- 错误提示图片和文本 -->
    <img :src="notFoundImg" alt="404" class="notfound-img"/>
    <p class="warning-txt">你好像走错地方了哦!!!</p>
    <div class="btn-group">
      <button @click="handleError('home')">主页</button>
      <button @click="handleError('back')">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import notFoundImg from '../../assets/image/404.png';

const router = useRouter();

const handleError = (t: string) => {
  switch (t) {
    case 'home':
      router.push('/');
      break;
    case 'back':
      router.go(-1);
      break;
  }
};
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
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.warning-txt {
  font-size: 32px;
  color: #79bbff;
  margin: 20px 0;
}

.btn-group {
  display: flex;
  gap: 15px;
}

button {
  padding: 10px 20px;
  background: rgb(103, 217, 232);
  color: #ffffff;
  border-radius: 36px;
  border: none;
  cursor: pointer;
  transition: background 0.3s ease;
}

button:hover {
  background: hotpink;
  border: none;
}
</style>
