<template>
  <div class="container">
    <div class="search_group">
      <div class="InputContainer">
        <input  class="input" placeholder="输入关键字搜索 动漫,剧集,电影" v-model="data.search" @keydown="e=>{e.keyCode==13 && searchMovie()}" />
        <div class="border" />
        <button class="micButton" @click="searchMovie">
          <svg viewBox="0 0 512 512" class="searchIcon"><path d="M416 208c0 45.9-14.9 88.3-40 122.7L502.6 457.4c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L330.7 376c-34.4 25.2-76.8 40-122.7 40C93.1 416 0 322.9 0 208S93.1 0 208 0S416 93.1 416 208zM208 352a144 144 0 1 0 0-288 144 144 0 1 0 0 288z"></path>
          </svg>
        </button>
      </div>
    </div>
    <div class="header-container">
        <h3>{{ data.oldSearch }}</h3>
        <span>共找到<b>{{ data.page.total }}</b>部与"<b>{{ data.oldSearch }}</b>"相关的影视作品</span>
    </div>
    <div class="card-group" v-if="data.list && data.list.length > 0">
      <div class="card"  v-for="(m,i) in data.list" >
        <div class="card-left">
          <div class="loader" v-if="m.loading" />
          <img v-show="!m.loading && !m.error" :src="m.picture" @load="()=>{data.list[i].loading=false}" @error="()=>{data.list[i].error=true}"/>
          <!-- 失败占位 -->
          <img v-if="m.error" src="/src/assets/image/404.png" @load="()=>{data.list[i].loading=false}" />
        </div>
        <div class="card-right">
          <h3>{{ m.name }}</h3>
          <p class="tags">
            <span class="tag_c">{{ `${m.cName?m.cName:'暂未分类'}` }}</span>
            <span>{{  `${m.year?m.year:'未知'}` }}</span>
            <span class="tag-area">{{  `${m.area?m.area:'未知'}`  }}</span>
          </p>
          <p><em>导演:</em>{{ `${m.director?m.director:'未知'}` }}</p>
          <p><em>主演:</em>{{ `${m.actor?m.actor:'未知'}` }}</p>
          <p class="blurb"><em>剧情:</em>{{ `${m.blurb.trim()?m.blurb.replace(/\s/g, ''):'暂无简介'}` }}</p>
          <el-button :icon="CaretRight" @click="play(m.id)">立即播放</el-button>
        </div>
      </div>
    </div>
    <el-empty v-if="data.oldSearch != '' && (!data.list || data.list.length == 0) " description="未查询到对应影片"/>

    <div class="pagination_container" v-if="data.list && data.list.length > 0" >
      <el-pagination background layout="prev, pager, next"
                     v-model:current-page="data.page.current"
                     @current-change="changeCurrent"
                     :pager-count="5"
                     :page-size="data.page.pageSize"
                     :total="data.page.total"
                     :prev-icon="ArrowLeftBold"
                     :next-icon="ArrowRightBold"
                     hide-on-single-page
                     class="pagination"/>
    </div>
  </div>
</template>

<script lang="ts" setup>

import {onMounted, reactive, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {ApiGet} from "../../utils/request";
import {ArrowLeftBold, ArrowRightBold, CaretRight, Search} from '@element-plus/icons-vue'
import {ElMessage} from "element-plus";

const router = useRouter()
const route = useRoute()
const data = reactive({
  list: [],
  page: {
    current: 0,
    pageSize: 0,
    total: 0,
  },
  oldSearch: '',
  search: '',
})
// 监听路由参数的变化
watch(
    [route],
    (oldRoute, newRoute) => {
      refreshPage(router.currentRoute.value.query.search, router.currentRoute.value.query.current)
    },
)

// 点击播放
const play = (id: string | number) => {
  location.href = `/play?id=${id}&episode=0&source=0`
}

// 搜索按钮事件
const searchMovie = () => {
  if (data.search.length <= 0) {
    ElMessage.error({message: '搜索信息不能为空', duration: 1000})
    return
  }
  location.href = location.href = `/search?search=${data.search}`
}

// 执行搜索请求
const refreshPage = (keyword: any, current: any) => {
  ApiGet('/searchFilm', {keyword: keyword, current: current}).then((resp: any) => {
    if (resp.code == 0) {
      data.list = resp.data.list.map((item: any) => {
        item.blurb = item.blurb.replace(/<[^>]*>/g, '')
        // 给列表内的每个元素添加加载中和失败状态
        item.loading = true
        item.error = false
        return item
      })
      data.page = resp.data.page
      data.oldSearch = keyword
      data.search = keyword
    } else {
      ElMessage.warning({message: resp.msg, duration: 1000})
    }

  })
}
// 组件挂载完成后触发数据初始化
onMounted(() => {
  if (router.currentRoute.value.query.search == null) {
    return
  }
  refreshPage(router.currentRoute.value.query.search + '', router.currentRoute.value.query.current)
})
// 分页器
const changeCurrent = (currentVal: number) => {
  let query = router.currentRoute.value.query
  location.href = `/search?search=${query.search}&current=${currentVal}`
}

</script>

<style scoped>
@import "/src/assets/css/pagination.css";


/*-------------------------pc------------------------------*/
@media (min-width: 768px) {
  .InputContainer {
    width: 40%;
  }
  .input {
    width: 90%;
  }
  .micButton {
    width: 10%;
  }
  .card-group {
    justify-content: space-between;
  }
  .card {
    width: calc(50% - 4%);
    min-height: 230px;
  }
  .card-left{
    width: 23%;
  }
  .card-right{
    width: 70%;
  }
  .card-right h3{
    font-size: 20px;
  }
  .card-right p {
    font-size: 15px;
    margin-bottom: 5px;
  }
  .card-right span {
    font-size: 12px;
  }
}
/*-------------------------wrap------------------------------*/
@media (max-width: 768px) {
  .InputContainer {
    width: 80%;
  }
  .input {
    width: 80%;
  }
  .micButton {
    width: 20%;
  }
  .card-group {
    justify-content: center;
  }
  .card {
    width: calc(100% - 10%);
    min-height: 185px;
  }
  .card-left{
    width: 38%;
  }
  .card-right{
    width: 60%;
  }
  .card-right h3{
    font-size: 16px;
  }
  .card-right p {
    font-size: 10px;
  }
  .card-right span {
    font-size: 10px;
  }
}
/*--------------------Search---------------------------*/
.InputContainer {
  height: 40px;
  margin: 20px auto;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  padding-left: 15px;
  box-shadow: 2px 2px 10px rgba(0, 0, 0, 0.075);
}

.input {
  height: 100%;
  border: none;
  outline: none;
  font-size: 0.9em;
  background-color: transparent !important;
  color: #ffffffbd;
  caret-color: #ffffff;
}
.input::placeholder {
  color: rgb(255 255 255 / 0.69);
}
input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0px 1000px transparent inset !important;
  box-shadow: 0 0 0px 1000px transparent inset !important;
}

.searchIcon {
  width: 13px;
}

.border {
  height: 40%;
  width: 1.3px;
  background-color: #ffffff82;
}


.micButton {
  padding: 0px 15px 0px 12px;
  border: none;
  background-color: transparent;
  height: 40px;
  cursor: pointer;
  transition-duration: .3s;
}

.searchIcon path {
  fill: #ffffffbd;
}


.micButton:hover {
  border-radius: 0 10px 10px 0;
  background-color: #fde2e245;
  transition-duration: .3s;
}
/*-----------------------header-container-------------------------------------------*/
.header-container{
  margin: 20px auto;
  padding: 0 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  border-bottom: 1px solid rgba(255,255,255, 0.15);
}
.header-container b {
  color: #FFC570;
}
.header-container h3 {
  margin: 0 auto;
}
/*-----------------------card style---------------------------------*/
.card-group {
  width: 100%;
  display: flex;
  flex-wrap: wrap
  /*  aspect-ratio: 9/16;*/
}
.card {
  display: flex;
  justify-content: start;
  aspect-ratio: 17/5;
  padding: 6px;
  border-radius: 12px;
  margin-bottom: 20px;
  background: rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  transform: translateY(0%);
}
.card:active {
  background-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-3%);
  transition: 0.5s;
}
.card-left {
  height: 100%;
}
.card img {
  width: 100%;
  height: 100%;
  border-radius: 6px;
}
.card-right {
  text-align: start;
  padding-left: 2%;
}
.card-right h3 {
  color: #FFC570;
  max-width: 100%;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.card-right .tags {
  display: flex;
  width: 90%;
  justify-content: start;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-right .tag-area {
  max-width: 39%;
  -webkit-line-clamp: 1;
  white-space: nowrap;
  overflow: hidden;
}
.card-right span {
  margin-right: 8px;
  border-radius: 3px;
  padding: 2px 4px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
.card-right .tag_c{
  background: rgba(155, 73, 231, 0.72);
}
.card-right p:not(.blurb) {
  max-width: 90%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-right .blurb {
  max-width: 90%;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}
.card-right p em {
  font-weight: bold;
  color: #fff;
  margin-right: 5px;
}
.card-right button {
  background-color: orange;
  border-radius: 20px;
  border: none !important;
  color: #ffffff;
  font-weight: bold;
  position: absolute;
  margin-bottom: 10px;
  bottom: 0;
}
.card-right h3, span, p {
  margin: 3px 0;
}
/*--------------------------image loading----------------------------------------------*/
/* 容器样式 */
.loader {
  height: 100%;
  width: 100%;
  background-color: transparent;
  background-image: linear-gradient(
      45deg,
      rgba(255,255,255,0.1) 25%,
      rgba(255,255,255,0.3) 25%,
      rgba(255,255,255,0.3) 50%,
      rgba(255,255,255,0.1) 50%,
      rgba(255,255,255,0.1) 75%,
      rgba(255,255,255,0.3) 75%,
      rgba(255,255,255,0.3) 100%
  );
  /* 3. 设置条纹的大小 (宽和高相等即为正方形条纹) */
  background-size: 20px 20px;
  /* 4. 应用动画 */
  animation: stripes-move 1.5s linear infinite;
}
/* --- 动画定义 --- */
@keyframes stripes-move {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 40% 40%;
  }
}
</style>

