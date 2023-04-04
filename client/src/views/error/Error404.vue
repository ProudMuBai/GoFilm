<template>
  <div class="container">
      <el-empty style="height: 100%; line-height: 100%; font-size: 20px" :image-size="400"
                image="/src/assets/image/404.png">
          <template #description>
              <p style="font-size: 32px;width: 100%;color: #a574b7">你好像走错地方了哦!!!</p>
          </template>
      </el-empty>
  </div>
</template>

<script>

</script>

<style scoped>
.container {
    height: 100vh;

}
</style>