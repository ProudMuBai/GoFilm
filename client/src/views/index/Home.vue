<template>
  <div class="container">
    <div v-if="global.isMobile" class=" banner_wrap" @touchstart="touchS" @touchend="touchE" @click="skipLink" >
      <el-carousel  v-model="data.banner.current" ref="wrap" :pause-on-hover="false"   :interval="5000" trigger="hover" height="200px" arrow="never" @change="carousel" >
        <el-carousel-item v-for="item in data.info.banners" :key="item"  >
          <el-image style="width: 100%; height: 100%;" :src="item.poster" fit="fill"/>
          <p class="carousel-title">{{ item.name }}</p>
        </el-carousel-item>
      </el-carousel>
    </div>
    <div v-if="!global.isMobile" class="banner"
         :style="{background:`url(${data.banner.current.poster})`, backgroundRepeat: 'no-repeat', backgroundSize: 'cover'}" @click="skipLink">
      <div class="preview">
        <el-carousel @change="carousel" :interval="5000" height="240px" arrow="always">
          <el-carousel-item v-for="item in data.info.banners" :key="item.id">
            <el-image style="width: 60%; height: 80%;border-radius: 5px;" :src="item.picture" fit="contain"/>
            <div class="carousel-tags">
              <span>{{ item.year }}</span>
              <span>{{ item.cName }}</span>
            </div>
            <p class="carousel-title">{{ item.name }}</p>
          </el-carousel-item>
        </el-carousel>
      </div>
    </div>
    <div class="content_item" v-for="item in data.info.content">
      <template v-if="item.nav.show">
        <el-row class="row-bg  cus_nav" justify="space-between">
          <el-col :span="12" class="title">
            <span :class="`iconfont ${item.nav.name.search('电影') != -1?'icon-film':item.nav.name.search('剧') != -1?'icon-tv':item.nav.name.search('动漫')!= -1?'icon-cartoon':'icon-variety'}`"
                style="color: #79bbff;font-size: 32px;margin-right: 10px; line-height: 130%"/>
            <a :href="`/filmClassify?Pid=${item.nav.id}`">{{ item.nav.name }}</a>
          </el-col>
          <el-col :span="12">
            <ul v-if="!global.isMobile" class="nav_ul">
              <template v-for="(c,i) in item.nav.children">
                <li class="nav_category" v-if="c.show && i < 6"><a
                    :href="`/filmClassifySearch?Pid=${c.pid}&Category=${c.id}`">{{ c.name }}</a></li>
              </template>
              <li class="nav_category"><a :href="`/filmClassify?Pid=${item.nav.id}`">更多 ></a></li>
            </ul>
          </el-col>
        </el-row>
        <el-row class="cus_content">
          <el-col :md="24" :lg="20" :xl="20" class="cus_content">
            <!--影片列表-->
            <FilmList v-if="item.movies" :col="6" :list="item.movies.slice(0,12)"/>
          </el-col>
          <el-col v-if="!global.isMobile" :md="0" :lg="4" :xl="4" class="content_right hidden-sm-and-down">
            <h3 class="hot_title">🔥热播{{ item.nav.name }}</h3>
            <template v-for="(m,i) in item.hot.slice(0,12)">
              <div class="content_right_item">
                <a :href="`/filmDetail?link=${m.mid}`"><b class="top_item">{{ i + 1 + '.' }}</b>
                  <span>{{ m.name }}</span></a>
              </div>
            </template>
          </el-col>
        </el-row>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
// 顶部轮播图
import 'element-plus/theme-chalk/display.css'
import {inject, onBeforeMount, reactive, ref} from "vue";
import {ApiGet} from "../../utils/request";
import FilmList from "../../components/index/FilmList.vue";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";

const data = reactive({
  info: {},
  banner: {
    current: {},
    touch: {index: 0, star: 0, end: 0,}
  }
})
const router = useRouter()
const global = inject<any>('global')

// pc 背景图同步响应
const carousel = (index: number) => {
  data.banner.current = data.info.banners[index]
}
const skipLink = ()=>{
  if (data.banner.current.mid) {
    router.push(`/filmDetail?link=`+data.banner.current.mid)
  }
}

// 滑动开始
const wrap = ref<HTMLFormElement>()
const touchS = (e:any)=>{
  //记录触摸起始位置
  data.banner.touch.star = e.changedTouches[0].pageX
}

//  滑动结束
const touchE = (e:any)=>{
  data.banner.touch.end = e.changedTouches[0].pageX
  let distance = data.banner.touch.end - data.banner.touch.star
  if (distance >= 50) {
    // let index = data.banner.touch.index - 1
    // data.banner.touch.index = index >= 0 ? index : banners.length-1
    wrap.value?.prev()
  } else if (distance <= -50) {
    // let index = data.banner.touch.index + 1
    // data.banner.touch.index = index <= banners.length - 1 ? index : 0
    wrap.value?.next()
  }
  // wrap.value?.setActiveItem(data.banner.touch.index)
}

onBeforeMount(() => {
  ApiGet('/index').then((resp: any) => {
    if (resp.code == 0) {
      data.info = resp.data
      data.banner.current = data.info.banners[0]
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
})
</script>

<style scoped>

.container {
  margin: 0 auto;
}

.content_item {
  padding: 10px;
  margin-bottom: 25px;
}

.title {
  display: flex;
  text-align: left;
  height: 40px;
}

.title > a {
  min-width: 40px;
  color: rgb(221, 221, 221)
}

a {
  color: #333;
  padding-top: 10px;
  text-decoration: none;
  outline: none;
  -webkit-tap-highlight-color: transparent;
}

.cus_nav {
  border-bottom: 1px solid rgb(46, 46, 46);

  height: 40px;
}

.nav_ul {
  list-style-type: none;
  display: flex;
  flex-direction: row;
  justify-content: end;
  margin: 0;
}

.nav_category > a {
  color: #c9c4c4;

}

.nav_category > a:hover {
  color: #1890ff;
}

.nav_ul > li {
  /*min-width: 60px;*/
  white-space: nowrap;
  line-height: 40px;
  margin: 0 8px;
  text-align: center;
  color: #999;
  font-size: 14px;
  font-weight: 400;
}


/*影片简介区域*/
.cus_content {
  display: flex;
  padding-top: 15px;
}

.content_right {
  width: 100%;
  padding-left: 18px;
}

.content_right_item {
  display: flex;
  padding-left: 10px;
  border-bottom: 1px solid rgb(46, 46, 46);
}

.content_right_item > a {
  padding: 10px 15px 10px 0;
  color: hsla(0, 0%, 100%, .87);
  display: block;
  flex-grow: 1;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  -o-text-overflow: ellipsis;
  white-space: nowrap;
}

.hot_title {
  text-align: left;
  margin: 8px 0;
}

:deep(.top_item) {
  color: red;
  /*font-style: oblique 10deg;*/
  font-style: italic;
  /*font-family: Inter;*/
  margin-right: 6px;
}

.content_right_item a span:hover {
  color: orange;
}


</style>

<!--移动端修改-->
<style scoped>
@media (min-width: 768px) {
  .cus_content_item {
    padding: 10px;
    overflow: hidden;
    /*margin-bottom: 10px;*/
  }
}

@media (max-width: 768px) {
  .cus_content_item {
    padding: 0 6px 0 0;
    margin-bottom: 10px;
    overflow: hidden;
  }

  .nav_ul {
    justify-content: end;
  }
}
</style>

<!--轮播图双端样式-->
<style scoped>
@media (max-width: 768px) {
  :deep(.el-carousel) {
    --el-carousel-arrow-size: 30px;
    --el-carousel-arrow-background: rgba(115, 133, 159, 0.5);
  }

  :deep(.el-carousel__arrow) {
    outline: none;
    border: none !important;
  }

  .el-carousel__item h3 {
    color: #475669;
    opacity: 0.75;
    line-height: 200px;
    margin: 0;
    text-align: center;
  }

  .el-carousel__item:nth-child(2n) {
    background-color: transparent;
  }

  .el-carousel__item:nth-child(2n + 1) {
    background-color: transparent;
  }

  :deep(.el-carousel__indicators) {
    width: 100% !important;
    text-align: right;
    height: 20px;
    line-height: 20px;
    padding-right: 10px;
    --el-carousel-indicator-padding-vertical: 0;
  }

  :deep(.el-carousel__button) {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    padding: 0 0!important;
    margin: 0 2px;
  }

  .banner_wrap {
    margin: -15px 0 20px 0;
    position: relative;
    box-shadow: 0 5px 30px 0 rgba(255, 255, 255, 0.15);
  }

  .carousel-tags {
    position: absolute;
    top: 170px;
    left: 25%;
  }

  .carousel-tags span {
    font-size: 12px;
    background: rgba(0, 0, 0, 0.55);
    color: #ffffff;
    padding: 2px 5px;
    margin: 2px 5px;
  }

  .carousel-title {
    font-size: 12px;
    position: absolute;
    bottom: 0;
    height: 20px;
    line-height: 20px;
    background: rgba(0, 0, 0, 0.5);
    text-align: left;
    width: 100%;
    margin: 0 auto;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; /* 显示省略号 */
  }
}


@media (min-width: 768px) {
  :deep(.el-carousel) {
    --el-carousel-arrow-size: 30px;
    --el-carousel-arrow-background: rgba(115, 133, 159, 0.5);
  }

  :deep(.el-carousel__arrow) {
    outline: none;
    border: none !important;
  }

  .container {
/*    padding-top: 660px;*/
  }

  .banner2 {
    height: 600px;
    position: absolute;
    margin-top: 60px;
    left: 0;
    top: 0;
    box-shadow: inset 0 -40px 30px 20px rgba(0, 0, 0, 0.6), 0 5px 30px 0 rgba(255, 255, 255, 0.15);
    padding: 2%;
    margin-bottom: 10px;
    border-radius: 0 0 6px 6px;
    width: 100%;
  }

  .preview2 {
    width: 260px;
    height: 200px;
    position: absolute;
    right: 50px;
    bottom: 60px;
  }

  .banner {
    height: 600px;
    box-shadow: inset 0 -40px 30px 20px rgba(0, 0, 0, 0.6), 0 5px 30px 0 rgba(255, 255, 255, 0.15);
    position: relative;
    padding: 2%;
    margin-bottom: 10px;
    border-radius: 6px;
    width: 100%;
  }

  .preview {
    width: 260px;
    height: 200px;
    position: absolute;
    right: 50px;
    bottom: 60px;
    /*  border: 1px solid skyblue;*/
  }

  .el-carousel__item h3 {
    color: #475669;
    opacity: 0.75;
    line-height: 200px;
    margin: 0;
    text-align: center;
  }

  .el-carousel__item:nth-child(2n) {
    background-color: transparent;
  }

  .el-carousel__item:nth-child(2n + 1) {
    background-color: transparent;
  }

  :deep(.el-carousel__indicators) {
    width: 100% !important;
  }

  :deep(.el-carousel__button) {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin: 0 2px;
  }

  .carousel-tags {
    position: absolute;
    top: 170px;
    left: 25%;
  }

  .carousel-tags span {
    font-size: 12px;
    background: rgba(0, 0, 0, 0.55);
    color: #ffffff;
    padding: 2px 5px;
    margin: 2px 5px;
  }

  .carousel-title {
    font-size: 12px;
    max-width: 50%;

    margin: 0 auto;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; /* 显示省略号 */
  }
}
</style>

