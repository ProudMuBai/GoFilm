<template>
  <div class="container">
    <h3>管理后台首页, 不知道放点啥, 先空着</h3>
  </div>
</template>

<script setup lang="ts">

import {onMounted} from "vue";
import {ApiGet} from "../../utils/request";

onMounted(()=>{
  ApiGet('/manage/index')
})

</script>


<style scoped>
.container {

}


</style>