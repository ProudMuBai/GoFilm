<template>
  <div class="relate_container">
    <p class="title">相关推荐</p>
    <FilmList :col="7" :list="relateList"/>
  </div>
</template>

<script lang="ts" setup>
import FilmList from "./FilmList.vue";

defineProps({
  relateList: Array
})

</script>

<style scoped>

@media (max-width: 768px){
  .title {
    font-size: 20px;
  }
}
@media (min-width: 768px){
  .title {
    font-size: 24px;
  }
}
.relate_container {
  width: 100%;
  margin-top: 36px;

}

.title {
  font-weight: 700;
  color: rgba(255, 255, 255, 0.55);
  padding-bottom: 5px;
  text-align: left;
  border-bottom: 2px solid rgba(255, 255, 255, 0.12);
}

</style>