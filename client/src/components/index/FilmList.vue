<template>
  <div class="c_content" v-if="d.list">
    <template v-if="d.list.length > 0">
      <div class="item film-card" v-for="item in d.list" :style="{width: `calc(${d.width-1}%)`}">
        <div v-if="item.id != -99 && global.isMobile" class="hidden-md-and-up">
          <a :href="`/filmDetail?link=${item.id}`" class="default_image link_content">
            <div class="tag_group">
              <span class="cus_tag ">{{ item.year ? item.year.slice(0, 4) : '未知' }}</span>
              <span class="cus_tag ">{{ item.cName }}</span>
              <span class="cus_tag ">{{ item.area.split(',')[0] }}</span>
            </div>
            <span class="cus_remark hidden-md-and-up">{{ item.remarks }}</span>
            <img :src="item.picture" :alt="item.name?.split('[')[0]" @error="handleImg">
          </a>
          <a :href="`/filmDetail?link=${item.id}`" class="content_text_tag">{{ item.name.split("[")[0] }}</a>
          <span class="cus_remark hidden-md-and-down">{{ item.remarks }}</span>
        </div>

        <div v-if="!global.isMobile"  class="film-card-inner">
          <div class="film-card-front">
            <a :href="`/filmDetail?link=${item.id}`" class="link_content">
              <div class="tag_group">
                <span class="cus_tag ">{{ item.year ? item.year.slice(0, 4) : '未知' }}</span>
                <span class="cus_tag ">{{ item.cName }}</span>
                <span class="cus_tag ">{{ item.area.split(',')[0] }}</span>
              </div>
              <span class="cus_remark hidden-md-and-up">{{ item.remarks }}</span>
              <img :src="item.picture" :alt="item.name?.split('[')[0]" @error="handleImg">
            </a>
          </div>
          <div class="film-card-back" @click="toDetail(item.id)">
            <p class="card-title" >{{item.name}}</p>
            <p v-show="item.blurb != ''" class="card-blurb">{{ item.blurb }}</p>
            <p v-show="item.blurb == ''" class="card-blurb"> 暂无简介 </p>
            <el-button class="card-detail" :icon="Discount" color="#626aef" plain round @click="toDetail(item.id)" >详情</el-button>
          </div>
        </div>
        <a v-if="!global.isMobile" :href="`/filmDetail?link=${item.id}`" class="content_text_tag hidden-sm-and-down">{{ item.name.split("[")[0] }}</a>

      </div>
    </template>
    <el-empty v-if="d.list.length <= 0" style="padding: 10px 0;margin: 0 auto" description="暂无相关数据"/>
  </div>
</template>

<script setup lang="ts">

import {defineProps, inject, reactive, watchEffect} from 'vue'
import {Discount} from "@element-plus/icons-vue";

const props = defineProps({
  list: Array,
  col: Number,
})
const d = reactive({
  col: 0,
  list: Array,
  width: 0,
})

const global = inject('global')
// 图片加载失败事件
const handleImg = (e: Event) => {
  e.target.style.display = "none"
}

const toDetail = (id:any) =>{
  location.href = `/filmDetail?link=${id}`
}

// 监听父组件传递的参数的变化
watchEffect(() => {
  // 首先获取当前设备类型
  const userAgent = navigator.userAgent.toLowerCase();
  let isMobile = /mobile|android|iphone|ipad|phone/i.test(userAgent)
  // 如果是PC, 为防止flex布局最后一行元素不足出现错位, 使用空元素补齐list
  let c = isMobile ? 3 : props.col ? props.col : 0
  let l: any = props.list
  let len = l.length
  d.width = isMobile ? 31 : Math.floor(100 / c)
  if (len % c != 0) {
    for (let i = 0; i < c - len % c; i++) {
      let temp: any = {...l[0] as any}
      temp.id = -99
      l.push(temp)
    }
  }
  d.list = l
})


</script>

<style scoped>
.default_image {
  background: url("/src/assets/image/404.png") no-repeat;
  background-size: cover;
}

:deep(.el-empty) {
  --el-empty-fill-color-1: rgba(155, 73, 231, 0.72);
  --el-empty-fill-color-2: #67d9e891;
  --el-empty-fill-color-3: rgb(106 19 187 / 72%);
  --el-empty-fill-color-4: #67d9e8;
  --el-empty-fill-color-5: #5abcc9;
  --el-empty-fill-color-6: #9fb2d9;
  --el-empty-fill-color-7: #61989f;
  --el-empty-fill-color-8: #697dc5;
  --el-empty-fill-color-9: rgb(43 51 63 / 44%);
}

/*wrap*/
@media (max-width: 768px) {
  /*展示区域*/
  .c_content {
    width: 100%;
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
  }

  .c_content .item {
    /*  flex-basis: calc(33% - 7px);
      max-width: 33%;*/
    margin: 0 4px 20px 4px;
    box-sizing: border-box;
    overflow: hidden;
  }

  .item .link_content {
    padding-top: 125%;
    position: relative;
    border-radius: 5px;
    display: flex;
    width: 100%;
    background-size: cover;
  }

  img {
    position: absolute;
    top: 0;
    left: 0;
    border-radius: 5px;
    object-fit: cover;
    width: 100%;
    height: 100%;
  }

  .tag_group {
    display: none;
  }

  .content_text_tag {
    font-size: 11px !important;
    color: rgb(221, 221, 221);
    width: 96% !important;
    max-height: 40px;
    line-height: 20px;
    padding: 2px 0 2px 0 !important;
    text-align: left;
    display: -webkit-box !important;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
  }

  .cus_remark {
    z-index: 10;
    position: absolute;
    bottom: 0;
    display: block;
    width: 100%;
    font-size: 12px;
    color: #c2c2c2;
    text-align: center;
    background: rgba(0, 0, 0, 0.55);
    border-radius: 0 0 5px 5px;
  }
}

/*pc*/
@media (min-width: 768px) {
  .c_content {
    width: 100%;
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
  }

  .item {
    margin-bottom: 20px;
    box-sizing: border-box;
  }

  .link_content {
    /*    padding-top: 125%;*/
    background-size: cover;
    width: 100%;
    display: flex;
    /*    position: relative;*/
    margin-bottom: 5px;
  }

  img {
    position: absolute;
    top: 0;
    left: 0;
    border-radius: 5px;
    object-fit: cover;
    width: 100%;
    height: 100%;
  }

  .tag_group {
    position: absolute;
    bottom: 3px;
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    overflow: hidden;
    justify-content: start;
    height: 18px;
    z-index: 10;
    line-height: 18px;
    padding-left: 10px;
  }

  .cus_tag {
    flex-shrink: 0; /* 不缩小元素 */
    white-space: nowrap;
    color: rgb(255, 255, 255);
    padding: 0 3px;
    margin-right: 8px;
    background: rgba(0, 0, 0, 0.55);
    font-size: 12px;
    border-radius: 5px;
  }

  .content_text_tag {
    display: block;
    font-size: 14px !important;
    color: rgb(221, 221, 221);
    width: 96% !important;
    padding: 2px 10px 2px 2px !important;
    text-align: left;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .cus_remark {
    display: block;
    width: 100%;
    padding-left: 3px;
    font-size: 12px;
    color: #999999;
    text-align: left;
  }
}
</style>


<style scoped>
.film-card {

  background-color: transparent;
  width: 100%;
  perspective: 1000px;
  font-family: sans-serif;
}

.film-card-inner {

  padding-top: 125%;
  position: relative;
  width: 100%;
  text-align: center;
  transition: transform 0.8s;
  transform-style: preserve-3d;
}

.film-card:hover .film-card-inner {
  transform: rotateY(180deg);
}

.film-card-front, .film-card-back {
  border-radius: 5px;
  position: absolute;
  top: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 100%;
  height: 100%;
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
/*  border: 1px solid rgba(255, 255, 255, 0.1);*/
  box-shadow: 0 25px 25px rgba(0, 0, 0, 0.25);
}

.film-card-front {
  border: none;
  background: url("/src/assets/image/404.png") no-repeat;
  background-size: cover;
}

.film-card-back {
  cursor: pointer;
  transform: rotateY(180deg);
  padding: 0 5px;
  background: linear-gradient(#fff2, transparent);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.card-title {
  max-width: 70%;
  margin: 0 auto;
  font-size: 14px ;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-blurb {
  margin-bottom: 30px;
  display: -webkit-box;
  -webkit-line-clamp: 5; /* 限制显示的行数 */
  -webkit-box-orient: vertical;
  overflow: hidden;
  font-size: 12px;
}

.card-detail {
  position: absolute;
  width: 60%;
  left: 20%;
  bottom: 5px;
}
</style>