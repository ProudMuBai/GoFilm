<template>
  <div class="container">
    <h3>管理后台首页, 不知道放点啥, 先空着</h3>
  </div>
</template>

<script setup lang="ts">

import {onMounted} from "vue";
import {ApiGet} from "../../utils/request";
import {ElMessage} from "element-plus";

onMounted(()=>{
  ApiGet('/manage/index').then((resp:any)=>{
      if (resp.code == 0) {
        ElMessage.success({message: resp.msg})
      } else {
        ElMessage.error({message: resp.msg})
      }
  })
})

</script>


<style scoped>
.container {

}


</style>