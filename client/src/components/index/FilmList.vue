<template>
  <div class="c_content" v-if="d.list">
    <template v-if="d.list.length > 0">
      <div class="item film-card" v-for="item in d.list" :style="{width: `calc(${d.width-1}%)`}">
        <!--v2测试-->
        <!--Wrap-->
        <template v-if="global.isMobile && item.id != -99">
          <div class="card">
            <img class="card-img" :src="item.picture" :alt="item.name?.split('[')[0]" @error="handleImg">
            <div class="tag_group">
              <span class="cus_tag " v-if="item.cName.replace(/\s/g, '')" >{{ item.cName }}</span>
              <span class="cus_tag ">{{ item.year ? item.year.slice(0, 4) : '未知' }}</span>
            </div>
            <a class="card-info" :href="`/filmDetail?link=${item.id}`">
              <p class="text-title">{{ item.name }}</p>
              <p v-if="item.blurb == '' || item.blurb == '...' " class="text-body"> 暂无简介 </p>
              <p v-else class="text-body">{{ item.blurb }}</p>
              <el-button class="card-button " :icon="Discount" color="#626aef" plain round @click="toDetail(item.id)">
                详情
              </el-button>
            </a>
          </div>
          <span class="cus_remark">{{ item.remarks }}</span>
          <a :href="`/filmDetail?link=${item.id}`" class="card-external-title"
             v-if="global.isMobile && item.id != -99">{{ item.name.split("[")[0] }}</a>
        </template>
        <!--PC-->
        <template v-if="!global.isMobile && item.id != -99">
          <div class="card">
            <img class="card-img" :src="item.picture" :alt="item.name?.split('[')[0]" @error="handleImg">
            <div class="tag_group">
              <span class="cus_tag ">{{ item.year.replace(/\s/g, '') ? item.year.slice(0, 4) : '未知' }}</span>
              <span v-if="item.cName.replace(/\s/g, '')" class="cus_tag ">{{ item.cName }}</span>
              <span class="cus_tag ">{{ item.area.replace(/\s/g, '') ? item.area.split(',')[0] : '未知' }}</span>
            </div>
            <a class="card-info" :href="`/filmDetail?link=${item.id}`">
              <p class="text-title">{{ item.name }}</p>
              <p v-if="item.blurb == '' || item.blurb == '...' " class="text-body"> 暂无简介 </p>
              <p v-else class="text-body">{{ item.blurb }}</p>
              <el-button class="card-button " :icon="Discount" color="#626aef" plain round @click="toDetail(item.id)">
                详情
              </el-button>
            </a>
          </div>
          <span class="cus_remark">{{ item.remarks }}</span>
          <a :href="`/filmDetail?link=${item.id}`"
             class="content_text_tag">{{ item.name.split("[")[0] }}</a>
        </template>

      </div>
    </template>
    <el-empty v-if="d.list.length <= 0" style="padding: 10px 0;margin: 0 auto" description="暂无相关数据"/>
  </div>
</template>

<script setup lang="ts">

import {inject, reactive, watchEffect} from 'vue'
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
  // e.target.style.display = "none"
  e.target.src = '/src/assets/image/404.png'
}

const toDetail = (id: any) => {
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
  // d.list = l
  d.list = l.map((item: any) => {
    item.blurb = item.blurb.replace(/<[^>]*>/g, '')

    if (item.remarks.match(/\d+/)) {
      if (item.remarks.includes("期")) {
        item.remarks = `${item.remarks.match(/\d+/)}期`
      } else if (item.remarks.includes("集")) {
        item.remarks = `${item.remarks.match(/\d+/)}集`
      }
    }

    // console.log(item.blur)
    return item
  })
})


</script>

<style scoped>
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

.c_content {
  max-width: 99%;
}

.card {
  width: 100%;
  aspect-ratio: 3/4.1;
  /*  padding: 1.9rem;*/
  padding: 0;
  background: #f5f5f5;
  position: relative;
  display: flex;
  align-items: flex-end;
  box-shadow: 0 7px 20px rgba(43, 8, 37, 0.2);
  /*  transition: all 0.3s ease-out;*/
  transition: transform 0.5s cubic-bezier(0.165, 0.84, 0.44, 1) 0.1s;
  overflow: hidden;
}

.card:before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2;
  transition: 0.5s;
}

.card-info {
  position: relative;
  z-index: 3;
  color: #f5f5f5cf;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  opacity: 0;
  width: 100%;
  height: 115%;
  padding-bottom: 16%;
  background: rgba(0, 0, 0, 0.5);
  transform: translateY(6%);
  /*  transition: 0.5s;*/
  transition: transform 0.8s cubic-bezier(0.165, 0.84, 0.44, 1) 0.1s;
}

/*Text*/
.text-title {
  max-width: 80%;
  margin: 0 auto;
  /*  font-size: 1rem;*/
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #FFC570;
}

.text-body {
  letter-spacing: 1px;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

/*Button*/
.card-button {
  margin: 0 auto;
  outline: none;
  border: none;
  font-weight: bold;
  transition: 0.4s ease;
}

/*Image*/
.card-img {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  object-fit: cover;
  object-position: center;
  transition: transform 0.8s ease-out;
}

/*Hover*/
.card:active, .card:hover {
  transform: translateY(-2%);
}

.card:hover .card-img {
  transform: scale(1.1);
}

.card:hover:before {
  opacity: 1;
}

.card:hover .card-info {
  z-index: 10;
  opacity: 1;
  transform: translateY(0);
}

.cus_remark {
  position: absolute;
  right: -5px;
  color: white;
  border-radius: 6px 3px 0 6px;
  font-weight: bold;
  z-index: 2;
}

.cus_remark::after {
  content: "";
  position: absolute;
  width: 0;
  height: 0;
  border-width: 6px;
  border-style: solid;
}
</style>

<style scoped>
/*wrap*/
@media (max-width: 768px) {
  /*展示区域*/
  .c_content {
    width: 98%;
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
  }

  .c_content .item {
    /*  flex-basis: calc(33% - 7px);
      max-width: 33%;*/
    position: relative;
    margin: 0 4px 20px 4px;
    box-sizing: border-box;
  }

  .item .link_content {
    padding-top: 125%;
    position: relative;
    border-radius: 5px;
    display: flex;
    width: 100%;
    background-size: cover;
  }

  .card {
    border-radius: 0.3em;
  }

  .card-img {
    border-radius: 0.3em;
  }

  .card-button {
    width: 56%;
    height: 20px;
    padding: 0.37em !important;
    border-radius: 4px;
    font-size: 12px;
    background: #00EAEAB3;
    color: white;
  }

  .card-button:active {
    background: rgb(89 205 205 / 0.66);
    color: #f5f5f5;
  }

  :deep(.el-button [class*=el-icon]+span) {
    margin-left: 0 !important;
  }

  .tag_group {
    display: block;
    width: 100%;
    position: absolute;
    font-size: 10px;
    overflow: hidden;
    height: 18px;
    z-index: 2;
    line-height: 18px;
    padding-left: 2px;
    margin-bottom: 2px;
    text-align: start;
  }

  .cus_tag {
    border-radius: 3px;
    background: rgba(0,0,0, 0.3);
    margin-right: 5px;
    padding: 2px 3px;
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

  .text-title {
    font-size: 0.8rem;
  }

  .text-body {
    font-size: 0.65rem;
    margin: 5px 0 8px 0;
  }

  .item:hover .cus_remark {
    visibility: hidden;
  }

  .cus_remark {
    top: 6px;
    font-size: 10px;
    padding: 2px 3px;
    background-color: #67d9e8cc;
  }

  .cus_remark::after {
    bottom: -6px;
    right: -6px;
    border-color: transparent transparent transparent #67d9e8cc;
  }

  .card-external-title {
    display: block;
    max-width: 70%;
    margin-top: 5px;
    font-size: 12px;
    font-weight: bold;
    text-align: start;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
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
    position: relative;
  }

  .card {
    border-radius: 0.5em;
  }

  .card-img {
    border-radius: 0.5em;
  }

  .card-button {
    width: 80%;
    padding: 0.7rem;
    border-radius: 4px;
    background: #ee9ca7;
    color: white;
  }

  .card-button:hover {
    background: #f954a6c7;
    color: #f5f5f5;
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

  .text-title {
    font-size: 1rem;
  }

  .text-body {
    font-size: 0.75rem;
    margin: 5px 0 15px 0;
  }

  .item:hover .cus_remark {
    visibility: hidden;
  }

  .cus_remark {
    top: 8px;
    font-size: 12px;
    padding: 3px 6px;
    background-color: #ff69b4fa; /* 粉红色 */
  }

  /* 核心：画小三角 */
  .cus_remark::after {
    bottom: -6px;
    right: -7px;
    border-color: transparent transparent transparent #ff69b4fa;
  }
}
</style>
