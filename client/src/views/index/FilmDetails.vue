<template>
  <div class="film" v-show="data.loading">
    <!-- 移动端title   -->
    <div v-if="global.isMobile">
      <div class="title_mt  ">
        <a class="picture_mt" href="" >
          <img class="picture-c" :src="data.detail.picture" alt="load error" @error="handleImg" />
        </a>
        <div class="title_mt_right">
          <h3>{{ data.detail.name }}</h3>
          <ul class="tags">
            <li style="margin: 2px 0">{{
                data.detail.classTag ? `${data.detail.classTag}`.replaceAll(",", " | ") : '未知'
              }}
            </li>
          </ul>
          <p><span>导演:</span> {{ data.detail.director }}</p>
          <p><span>主演:</span> {{ handleLongText(data.detail.actor) }}</p>
          <p><span>上映:</span> {{ data.detail.releaseDate }}</p>
          <p><span>地区:</span> {{ data.detail.area }}</p>
          <p v-if="data.detail.remarks"><span>连载:</span>{{ data.detail.remarks }}</p>
          <!--<p><span>评分:</span><b id="score">{{ data.detail.dbScore }}</b></p>-->
        </div>
      </div>
      <div class="mt_content">
        <p v-html="`${data.detail.content}`.replaceAll('　　', '')"></p>
      </div>
    </div>
    <!-- pc端title-->
    <div v-else class="title" >
      <a href="" class="picture">
        <img class="pic" :src="data.detail.picture" alt="load error" @error="handleImg" >
      </a>
      <h2>{{ data.detail.name }}</h2>
      <ul class="tags">
        <li class="t_c">
          <a :href="`/filmClassifySearch?Pid=${data.detail.pid}&Category=${data.detail.cid}`">
            <el-icon>
              <Promotion/>
            </el-icon>
            {{ data.detail.cName.replace(/\s/g, '') ? data.detail.cName : '暂无分类' }}
          </a>
        </li>
        <li v-if="data.detail.classTag">
          {{ `${data.detail.classTag}`.replaceAll(",", "&emsp;") }}
        </li>
        <li>{{ data.detail.year }}</li>
        <li>{{ data.detail.area }}</li>
      </ul>
      <p><span>导演:</span> {{ data.detail.director.replace(/\s/g, '') ? data.detail.director : '未知' }}</p>
      <p><span>主演:</span> {{ data.detail.actor.replace(/\s/g, '') ? data.detail.actor : '未知' }}</p>
      <p><span>上映:</span> {{ data.detail.releaseDate }}</p>
      <p v-if="data.detail.remarks"><span>连载:</span>{{ data.detail.remarks }}</p>
      <p><span>评分:</span><b id="score">{{ data.detail.dbScore }}</b></p>
      <div class="cus_wap">
        <p style="min-width: 40px"><span>剧情:</span></p>
        <p ref="textContent" class="text_content">
          <el-button v-if="`${data.detail.content}`.length > 140" class="multi_text" style="color:#a574b7;"
                     @click="showContent(multiBtn.state)" link>{{ multiBtn.text }}
          </el-button>
          <span class="cus_info" v-html="data.detail.content"></span>
        </p>
      </div>
      <p>
        <el-button type="warning" class="player" size="large" @click="play({episode:0,source:data.detail.list[0].id})"
                   round>
          <el-icon>
            <CaretRight/>
          </el-icon>
          立即播放
        </el-button>
        <el-button color="#9b49e7" class="player" size="large" round plain>
          <el-icon>
            <Star/>
          </el-icon>
          收藏
        </el-button>
      </p>
    </div>
    <!--播放列表-->
    <div class="play-module">
      <div class="play-module-item">
        <div class="module-heading">
          <p class=" play-module-title">播放列表</p>
          <div class="play-tab-group">
            <a href="javascript:;" :class="`play-tab-item ${data.currentTabId == item.id ? 'tab-active':''}`"
               v-for="item in data.detail.list" @click="changeTab(item.id)">{{ item.name }}</a>
          </div>
        </div>
        <div class="play-list">
          <div class="play-list-item" v-show="data.currentTabId == item.id " v-for="item in data.detail.list">
            <a class="play-link" v-for="(v,i) in item.linkList" href="javascript:;"
               @click="play({source: item.id, episode: i})">{{ v.episode }}</a>
          </div>
        </div>
      </div>
    </div>

    <!--相关系列影片-->
    <div class="correlation">
      <RelateList :relate-list="data.relate"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import {inject, onBeforeMount, reactive, ref,} from "vue";
import {ApiGet} from "../../utils/request";
import {ElMessage} from 'element-plus'
import {Promotion, CaretRight, Star} from "@element-plus/icons-vue";
import RelateList from "../../components/index/RelateList.vue";
import notFoundImg from '/src/assets/image/404.png'

// 获取路由对象
const router = useRouter()
const data = reactive({
  detail: {
    id: '',
    cid: '',
    pid: '',
    name: '',
    picture: '',
    playFrom: [],
    DownFrom: '',
    playList: [[]],
    downloadList: '',
    subTitle: '',
    cName: '',
    enName: '',
    initial: '',
    classTag: '',
    actor: '',
    director: '',
    writer: '',
    blurb: '',
    remarks: '',
    releaseDate: '',
    area: '',
    language: '',
    year: '',
    state: '',
    updateTime: '',
    addTime: '',
    dbId: '',
    dbScore: '',
    hits: '',
    content: '',
    list: []
  },
  relate: [],
  loading: false,
  currentTabId: '',
})

// 获取全局属性
const global = inject('global')

// 对部分信息过长进行处理
const handleLongText = (t: string): string => {
  let res = ''
  t.split(',').forEach((s, i) => {
    if (i < 3) {
      res += `${s} `
    }
  })
  return res.trimEnd()
}

// 图片加载失败事件
const handleImg = (e: Event) =>{
  e.target && (e.target.src = notFoundImg)
}

// 播放源切换
const changeTab = (id: string) => {
  data.currentTabId = id
}

// 选集播放点击事件
const play = (change: { source: string, episode: number }) => {
  router.push({path: `/play`, query: {id: `${router.currentRoute.value.query.link}`, ...change}})
}

// 内容展开收起效果
const multiBtn = ref({state: false, text: '展开'})
const textContent = ref()
const showContent = (flag: boolean) => {
  if (flag) {
    multiBtn.value = {state: !flag, text: '展开'}
    textContent.value.style.webkitLineClamp = 2
    return
  }
  multiBtn.value = {state: !flag, text: '收起'}
  textContent.value.style.webkitLineClamp = 8
}

// 页面加载数据初始化
onBeforeMount(() => {
  let link = router.currentRoute.value.query.link
  ApiGet('/filmDetail', {id: link}).then((resp: any) => {
    if (resp.code === 0) {
      data.detail = resp.data.detail
      // 去除影视简介中的无用内容和特殊标签格式等
      data.detail.name = data.detail.name.replace(/(～.*～)/g, '')
      data.detail.content = data.detail.content.replace(/(&.*;)|( )|(　　)|(\n)|(<[^>]+>)/g, '')
      data.relate = resp.data.relate
      // 处理过长数据
      data.detail.actor = handleLongText(data.detail.actor)
      data.detail.director = handleLongText(data.detail.director)
      data.currentTabId = resp.data.detail.list[0].id
      data.loading = true
    } else {
      ElMessage({
        type: "error",
        dangerouslyUseHTMLString: true,
        message: resp.msg,
      })
    }
  })

})
</script>


<!--移动端适配-->
<style scoped>
@media (max-width: 768px) {
  .title_mt {
    width: 100%;
    padding: 0 3px;
    display: flex;
    flex-direction: row;
    flex-flow: nowrap;
    overflow: hidden;
  }

  .picture_mt {
    min-width: 35%;
    max-width: 35%;
    margin-right: 12px;
    border-radius: 5px;
    background-size: cover;
  }

  .picture_mt:active {
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.15);
  }

  .picture-c {
    width: 100%;
    aspect-ratio: 3/4;
    border-radius: 5px;
  }

  .title_mt_right {
    flex: 1;
    text-align: left;
  }

  .title_mt_right h3 {
    max-width: 90%;
    font-size: 14px;
    margin: 0 0 5px 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .title_mt_right p {
    font-size: 12px;
    margin: 3px 2px;
    white-space: nowrap;
  }

  .mt_content {
    margin-top: 5px;
    border-top: 1px solid var(--border-color-highlight);
    border-bottom: 1px solid  var(--border-color-highlight);
    width: 100%;
    padding: 5px;
  }

  .mt_content p {
    max-width: 96%;
    margin: 0 auto;
    font-size: 12px;
    text-align: left;
    word-wrap: break-word;
  }

  .play_content a {
    white-space: nowrap;
    color: #ffffff;
    border-radius: 6px;
    margin: 6px 8px;
    background: #888888;
    min-width: calc(25% - 16px);
    font-size: 12px;
    padding: 6px 12px !important;
  }
}
</style>

<style scoped>
@import "/src/assets/css/film.css";

.correlation {
  width: 100%;
}

.film {
  width: 100%;
  padding: 0 1%;
}

/*影片播放列表信息展示*/

/*顶部影片信息显示区域*/

.title {
  width: 100%;
  background: linear-gradient(#fff2, transparent);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 5px;
  padding: 5px 30px 30px 30px;
  position: relative;
}

.title > h2 {
  text-align: left;
  color: var(--text-content-color-light);
}

.picture {
  position: absolute;
  width: 220px;
  aspect-ratio: 3/4;
  right: 30px;
  top: 30px;
  border-radius: 8px;
  background-size: cover;
  overflow: hidden;
}

.picture:hover {
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.15);
}

.picture .pic {
  width: 100%;
  aspect-ratio: 3/4;

}


.picture::before {
  content: '';
  position: absolute;
  top: 0;
  left: -60%;
  width: 25%;
  height: 100%;
  opacity: 1;
  background: linear-gradient(to right, rgba(255, 255, 255, 0) 0%, rgb(255, 255, 255, 0.5) 50%, rgba(255, 255, 255, 0) 100%);
  transform: skewX(-25deg);
}

.picture:hover::before {
  left: 130%;
  transition: all 0.6s ease-in-out;
  opacity: 1;
}

.tags {
  list-style-type: none;
  display: flex;
  justify-content: left;
  margin: 0;
  padding: 0;
}

.tags > li {
  padding: 6px 10px;
  border-radius: 5px;
  background: linear-gradient(#ffffff14, transparent);
  border: 1px solid rgba(255, 255, 255, 0.1);
  margin: 0 8px;
  font-size: 12px;
  color: var(--text-content-color-light);
}

.tags > .t_c {
  background: linear-gradient(#9b49e7, #9b49e7bf);
  margin-left: 0;
}

.t_c a {
  color: #c4c2c2;
}

.title p {
  text-align: left;
  font-size: 14px;
  margin: 20px 0;
  max-width: 60%;
  color: var(--text-content-color-light);
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.title p span {
  font-size: 15px;
  font-weight: bold;
  color: var(--text-content-color-light);
  margin-right: 5px;
}

#score {
  color: #1cbeb9;
}

.cus_wap {
  display: flex;
}

.title .text_content {
  max-width: 70%;
  margin: 20px 3px;
  line-height: 22.5px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  vertical-align: top;
  /*margin-top: 5px;*/

}

.text_content::before {
  content: '';
  float: right;
  width: 0; /*设置为0，或者不设置宽度*/
  height: calc(100% - 20px); /*先随便设置一个高度*/
}

.text_content .cus_info {
  height: 100%;
  margin: 0;
  font-size: 15px !important;
  font-weight: normal;

}

.multi_text {
  float: right;
  clear: both;
  margin-right: 10px;
}

.el-icon {
  margin-right: 3px;
}

</style>

