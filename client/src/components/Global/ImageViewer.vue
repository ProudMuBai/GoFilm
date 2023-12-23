<template>
  <el-image-viewer
      v-if="data.show"
      :urlList="data.list"
      :z-index="data.zIndex"
      :initial-index="data.initialIndex"
      :infinite="data.infinite"
      :hideOnClickModal="data.hideOnClickModal"
      @close="data.show = false"
  ></el-image-viewer>
</template>

<script setup lang="ts">

import {onMounted, reactive, watch} from "vue";

const props = defineProps({
  options: {
    type:Object,
    default: {
      list:Array,
      currentLink: String,
      show:Boolean,
    }
  },
  remove: {
    type:Function,
    default: null,
  }
})

const data = reactive({
  show: false,
  list: [{link:''}],
  zIndex: 2000,
  initialIndex: 0,
  infinite: true,
  hideOnClickModal: false,
})

onMounted(()=>{
  data.list = props.options.list
  data.list.forEach((item,index)=>{
    if (item == props.options.currentLink) {
      data.initialIndex = index
    }
  })
  data.show = props.options?.show
})

watch([data],()=>{
 !data.show && props.remove()
})
</script>

<style scoped>

</style>