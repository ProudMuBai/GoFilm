<template>
  <div class="player_area" v-show="data.loading">
    <div id="player-full" />
    <div ref="playerContainer" class="player_container" />
    <div class="current_play_info">
      <div class="play_info_left">
        <h3 class="current_play_title"><a
            :href="`/filmDetail?link=${data.detail.mid}`">{{ data.detail.name }}</a>{{ data.current.episode }}</h3>
        <div class="tags">
          <a :href="`/filmClassifySearch?Pid=${data.detail.pid}&Category=${data.detail.cid}`">
            <el-icon>
              <Promotion/>
            </el-icon>
            {{ data.detail.cName }}</a>
          <span>{{
              data.detail.classTag ? data.detail.classTag.replaceAll(',', '/') : '未知'
            }}</span>
          <span class="hidden-sm-and-down">{{ data.detail.year }}</span>
          <span class="hidden-sm-and-down">{{ data.detail.area }}</span>
        </div>
      </div>
      <div class="play_info_right">
        <a href="javascript:;" :class="`iconfont icon-play1 ${data.autoplay?'p_r_active':''}`"
           @click="()=>{data.autoplay= !data.autoplay}"></a>
        <a v-show="hasNext" href="javascript:;" class="iconfont icon-iov-next"
           @click="playNext"></a>
      </div>
    </div>
    <!-- 播放选集   -->
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
          <div class="play-list-item" v-show="data.currentTabId == item.id" v-for="item in data.detail.list">
            <a :class="`play-link ${v.link == data.current.link?'play-link-active':''}`" v-for="(v,i) in item.linkList"
               href="javascript:;" @click="playChange({sourceId: item.id, episodeIndex: i, target: this})">{{
                v.episode
              }}
              <div class="loading-wave" v-if="v.link == data.current.link">
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
              </div>
              <div class="loading-wave" v-else></div>
            </a>
          </div>
        </div>
      </div>
    </div>
    <!--相关推荐-->
    <div class="correlation">
      <RelateList :relateList="data.relate"/>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  computed, inject,
  onBeforeMount,
  reactive,
  ref,
  watch,
} from "vue";
import {useRouter} from "vue-router";
import {ApiGet} from "@/utils/request";
import {ElMessage} from "element-plus";
import RelateList from "../../components/index/RelateList.vue";
import {Promotion} from "@element-plus/icons-vue";
import posterImg from '../../assets/image/play.png'
import {cookieUtil, COOKIE_KEY_MAP} from '@/utils/cookie'
// 引入视频播放器组件
import Player, {Events, Plugin} from "xgplayer"
import 'xgplayer/dist/index.min.css';
import HlsPlugin from 'xgplayer-hls'
import {fmt} from "@/utils/format";

// 播放页所需数据
const data = reactive({
  loading: false,
  detail: {
    id: '',
    mid: '',
    cid: '',
    pid: '',
    name: '',
    picture: '',
    playFrom: [],
    DownFrom: '',
    list: [[]],
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
  },
  current: {index: 0, episode: '', link: ''},
  relate: [],
  currentTabId: '', // 当前播放源ID
// @videojs-player 播放属性设置
  autoplay: true,
  options: {
    title: "", //视频名称
    url: "", //视频源
    volume: 0.6, // 音量
    currentTime: 50,
    autoplay: false,
    urls: [],
  },
})

// 获取路由信息
const router = useRouter()
const global = inject<any>('global')
// 是否存在下一集
const hasNext = computed(() => {
  let flag = false
  data.detail.list.forEach((item: any) => {
    if (data.currentTabId == item.id) {
      flag = data.current.index != item.linkList.length - 1
    }
  })
  return flag
})
// 播放源切换事件
const changeTab = (id: string) => {
  data.currentTabId = id
}
// 点击播集数播放对应影片
const playChange = (play: { sourceId: string, episodeIndex: number, target: any }) => {
  data.detail.list.forEach((item: any) => {
    if (item.id == play.sourceId) {
      let currPlay = item.linkList[play.episodeIndex]
      data.current = {index: play.episodeIndex, episode: currPlay.episode, link: currPlay.link}
      data.options.url = currPlay.link
      data.options.title = data.detail.name + "  " + currPlay.episode
      data.currentTabId = play.sourceId
    }
  })
}
// player相关事件
// 点击下一集按钮
const playNext = () => {
  // 如果不存在下一集信息则直接返回
  if (!hasNext.value) {
    return
  }
  if (data.autoplay) {
    setTimeout(() => {
      playChange({sourceId: data.currentTabId, episodeIndex: data.current.index + 1, target: ''})
    }, 100)
  }
}

//影片信息加入本地的观看历史中, 先获取cookie,已存在则追加,否则添加
const saveFilmHistory = () => {
  if (data.options.url.length > 0) {
    // 处理播放历史要记录的影片相关信息
    let history = cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY) ? JSON.parse(cookieUtil.getCookie(COOKIE_KEY_MAP.FILM_HISTORY)) : {}
    let link = `/play?id=${data.detail.mid}&source=${data.currentTabId}&episode=${data.current.index}&currentTime=${mPlayer.currentTime}`
    // 处理播放时长
    let timeStamp = new Date().getTime()
    let time = fmt.dateFormat(timeStamp)
    let progress = `${fmt.secondToTime(mPlayer.currentTime)} / ${fmt.secondToTime(mPlayer.duration)}`
    history[data.detail.mid] = {
      id: data.detail.mid,
      name: data.detail.name,
      picture: data.detail.picture,
      episode: data.current.episode,
      time: time,
      timeStamp: timeStamp,
      source: data.currentTabId,
      link: link,
      currentTime: mPlayer.currentTime,
      duration: mPlayer.duration,
      progress: progress,
      devices: global.isMobile
    }
    // 将历史记录添加到cookie中
    cookieUtil.setCookie(COOKIE_KEY_MAP.FILM_HISTORY, JSON.stringify(history))
  }
}

// 在浏览器关闭前或页面刷新前将当前影片的观看信息存入历史记录中
window.addEventListener('beforeunload', saveFilmHistory)

// 初始化页面数据
onBeforeMount(() => {
  let query = router.currentRoute.value.query
  ApiGet(`/filmPlayInfo`, {id: query.id, playFrom: query.source, episode: query.episode}).then((resp: any) => {
    if (resp.code === 0) {
      data.detail = resp.data.detail
      data.current = {index: resp.data.currentEpisode, ...resp.data.current}
      data.relate = resp.data.relate
      // 设置当前的视频播放url
      data.options.url = data.current.link
      // 设置当前播放源ID信息
      data.currentTabId = resp.data.currentPlayFrom
      data.loading = true
      data.detail.list.forEach((item: any) => {
        if (resp.data.currentPlayFrom == item.id) {
          data.options.urls = item.linkList.map((i: any) => {
            return i.link
          })
        }
      })

    } else {
      ElMessage.error({message: resp.msg})
    }
  }).then(() => {
    // 拿到数据后初始化播放器
    mPlayer = new Player({
      el: playerContainer.value,
      url: data.options.url,
      poster: posterImg,
      width: "100%",
      height: "100%",
      fluid: true,
      videoFillMode: "contain",
      autoplay: data.options.autoplay,
      lang: 'zh-cn', // 设置语言为中文
      volume: 0.7,   // 初始音量
      playbackRate: [3, 2, 1.5, 1, 0.75, 0.5],
      playnext: {
        urlList: data.options.urls,
      },
      playsinline: true,
      miniprogress: true,
      "x5-video-orientation": "landscape",
      "x5-video-player-fullscreen": "true",
      plugins: [HlsPlugin, playListPlugin],
      hls: {
        retryCount: 3, // 重试 3 次，默认值
        retryDelay: 1000, // 每次重试间隔 1 秒，默认值
        loadTimeout: 10000, // 请求超时时间为 10 秒，默认值
        fetchOptions: {mode: 'cors'},
        targetLatency: 10, // 直播目标延迟，默认 10 秒
        maxLatency: 20, // 直播允许的最大延迟，默认 20 秒
        preloadTime: 100 ,// 默认值
        disconnectTime: 0, // 直播断流时间，默认 0 秒，（独立使用时等于 maxLatency）
        // preloadTime: 30 // 默认值
      },
      controls: {
        autoHide: true
      },
      keyboard: {playbackRate: 3},
      mobile: {
        rotateFullscreen: true,
        hideDefaultControls: true,
        gestureX: true,
        gestureY: true,
        scopeR: 0.15,
        pressRate: 3,//长按倍速
        disablePress: false,
      },
      // controls: {
      //   autoHide: false,
      // },
    })
    // 播放器初始化完成时设置播放时长参数
    mPlayer.on(Events.READY, () => {
      // 从router参数中获取时间信息
      let currentTime = router.currentRoute.value.query.currentTime
      if (currentTime) {
        mPlayer.currentTime = currentTime
      }
    })
    // 播放完成事件
    mPlayer.on(Events.ENDED, () => {
      data.autoplay && playNext()
    })
    // 下一集按钮点击事件
    mPlayer.on(Events.PLAYNEXT, () => {
      playNext()
    })
  })
})

// 获取playerContainer挂载节点
const playerContainer = ref<HTMLDivElement | undefined>(undefined)
let mPlayer: any = null

// 监测播放器数据信息变化
watch(data.options, (newVal) => {
  if (mPlayer) {
    mPlayer.pause();
    mPlayer.currentTime = 0
    mPlayer.src = newVal.url
    // mPlayer.load()
    mPlayer.play().then(()=>{
      let playBtn = mPlayer.root.querySelector('.xg-icon-play')
      if (playBtn) {
        playBtn.style.display = 'none'
      }
    })
  }
})
// 自定义播放列表插件
const {POSITIONS} = Plugin
class playListPlugin extends Plugin {
  // 插件的名称，将作为插件实例的唯一key值
  static get pluginName() {
    return 'customPlayList'
  }

  static get defaultConfig() {
    return {
      // 挂载在controls的右侧，如果不指定则默认挂载在播放器根节点上
      position: POSITIONS.CONTROLS_RIGHT
    }
  }

  constructor(args: any) {
    super(args)
  }

  // 定义属性类型
  private listContainer: HTMLElement | null = null;
  private currentIndex: number = 0;

  private renderListItems() {
    if (!this.listContainer) return;
    // 清空现有内容
    this.listContainer.innerHTML = '';
    this.listContainer.className = 'playListContainer'
    // 数据渲染
    let l: any = data.detail.list.find((item: any) => {
      if (item.id == data.currentTabId) {
        return item
      }
    })
    l.linkList.forEach((item: any, index: number) => {
      const el = document.createElement('div')
      el.className = 'playlist-item';
      el.innerText = item.episode;
      // 选中项高亮样式
      const isActive = index === data.current.index
      if (isActive) {
        el.className = 'playlist-item active'
      }
      // 绑定点击切换视频事件
      el.onclick = (e) => {
        e.stopPropagation();
        playChange({sourceId: data.currentTabId, episodeIndex: index, target: el})
        // 重新渲染以更新高亮状态
        this.renderListItems();
        // 播放后自动收起列表（可选体验优化）
        this.toggleList()
      }
      el.addEventListener('wheel', (e:any)=> {
          e.preventDefault()
          this.listContainer && (this.listContainer.scrollTop += e.deltaY)
      });
      this.listContainer && this.listContainer.appendChild(el);
    })
  }

  // 播放器初始化完成后触发
  afterPlayerInit() {
    // TODO 播放器调用start初始化播放源之后的逻辑
    this.bind('click', (e: any) => {
      console.log('---------------------------------click')
      e.stopPropagation(); // 阻止冒泡
      this.toggleList();
    })
    // 点击播放器其他区域时关闭列表
    if (this.player.root) {
      this.player.root.addEventListener('click', () => {
        if (this.listContainer) {
          this.listContainer.style.display = 'none'
        }
      })
    }
    this.listContainer = document.querySelector('.playList-panel')
    this.listContainer && this.listContainer.addEventListener('mouseout', (e) => {
      // 手指离开时的操作
      e.stopPropagation();
      this.toggleList()
    })
    this.renderListItems()
  }

  // 切换列表显示状态
  private toggleList() {
    if (this.listContainer) {
      const isHidden = (this.listContainer.style.display == 'none' || this.listContainer.style.display == '')
      this.listContainer.style.display = isHidden ? 'block' : 'none'
      let mobilePlugin = mPlayer.getPlugin('mobile');
      if (isHidden) {
        // console.log(mPlayer.plugins)
        // console.log(mPlayer.getPlugin('cssfullscreen'))
        // mPlayer.getPlugin('cssfullscreen').show();
        // mobilePlugin.disable();
      } else {
        // mPlayer.getPlugin('cssfullscreen').hide();
        // mobilePlugin.enable();
      }
    }
  }

  afterCreate() {
  }

  destroy() {
    this.listContainer = null;
  }


  render() {
    return `<xg-icon class="iconfont icon-dianying1" ><div class="playList-panel" /></xg-icon>`
  }
}
</script>

<style>
/*xgplayer 样式修改*/
/*播放容器*/
.player_container {
  width: 100%;
  padding-top: 56.29%!important;
  aspect-ratio: 16 / 9;
  margin: 0;
  position: relative;
  border-radius: 6px;
  display: flex;
  box-shadow: 3px 3px 12px rgba(255, 255, 255, 0.2);
}
.player_area .xgplayer-is-fullscreen {
  padding-top: 0 !important;
}
/*进度条颜色*/

.xgplayer .xgplayer-progress-played,.xg-mini-progress xg-mini-progress-played {
  background: linear-gradient(-90deg, #00EAEA80 0%, #E337F780 100%);
}

.xg-right-grid .icon-dianying1 {
  display: block;
  color: #ffffff;
  padding: 0;
  font-size: 16px;
  line-height: 40px;
}
.xgplayer-playlist-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  height: 100%;
}
.xgplayer-playnext .xgplayer-icon svg {
  width: 15px;
}
.xgplayer xg-icon:not(.xgplayer-playnext) svg {
  width: 20px;
}
.xgplayer .xgplayer-progress-btn{
  width: 10px;
  height: 10px;
  border-radius: 10px;
}
.xgplayer .xgplayer-progress-btn:before{
  width: 8px;
  height: 8px;
  border-radius: 8px;
}
.xgplayer .flex-controls .xg-inner-controls{
  bottom: 0;
}

@media only screen and (orientation: landscape) {
  .xgplayer-mobile.xgplayer-is-fullscreen .xg-top-bar, .xgplayer-mobile.xgplayer-is-fullscreen .xg-pos{
    left: 3%;
    right: 3%;
  }
  .xgplayer .xgplayer-playnext svg{
    width: 16px !important;
  }

}

.playListContainer {
  display: none;
  position: absolute;
  bottom: 100%;
  right: 0;
  width: 100px;
  max-height: 180px;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 4px;
  margin-bottom: 6px;
}

.playListContainer .playlist-item {
  font-size: 16px;
  max-height: 30px;
  line-height: 30px;
  color: #fff;
  cursor: pointer;
  border-bottom: 1px solid #333;
  background: transparent;
  font-weight: normal;
}

.playListContainer .playlist-item:hover {
  color: #9a5dd3cc;
}

.playListContainer .active {
  background: #9a5dd3cc;
  font-weight: bold;
}

</style>

<!--公共样式-->
<style scoped>
@import "/src/assets/css/film.css";
/*公共样式区域*/
/*当前播放的影片信息展示*/
.current_play_info {
  width: 100%;
  padding: 15px 5px;
  text-align: left;
  display: flex;
  justify-content: space-between;
}

.current_play_title {
  font-weight: 600;
  color: rgb(201, 196, 196);
  margin: 0 0 12px 0;
}

.current_play_title a {
  color: rgb(201, 196, 196);
  font-weight: 600;
  margin-right: 16px;
}

.current_play_title a:hover {
  color: orange;
}

/*自动播放 & 下一集链接*/

.play_info_right a {
  margin-left: 10px;
  padding: 15px 20px;
  display: inline-block;
  font-size: 20px;
  height: 100%;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
}

.p_r_active {
  color: #FFBB5C;
}


/* 播放区域*/
.player_area {
  width: 100%;
  min-height: 100%;
}




/*右侧播放源选择区域*/
/*影片播放列表信息展示*/
/*影片播放列表信息展示*/
.play_list {
  width: 100%;
  border-radius: 10px;
  background: #2e2e2e;
  margin-top: 50px;
  position: relative;
}

.play_content {
  display: flex;
  flex-flow: row wrap;
  padding: 10px 10px 10px 10px;

}

.play_list > h2 {
  position: absolute;
  left: 10px;
  top: -10px;
  z-index: 50;
}

/*推荐列表区域*/
.correlation {
  width: 100%;
}

</style>

<!--移动端 && pc端-->
<style scoped>

/*适应小尺寸*/
@media (max-width: 768px) {
  :deep(.xgplayer xg-start-inner){
    border-radius: 50%!important;
    background: rgba(0, 0, 0, .38)!important;
  }
  .player_area {
    padding: 5px 10px;
  }

  .tags a {
    padding: 5px 10px;
    color: #c4c2c2;
    background: linear-gradient(#9B49E7B8, #9B49E799);
    font-size: 13px;
    border-radius: 6px;
    margin-right: 3px;
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

  .tags span {
    padding: 6px 10px;
    background: linear-gradient(#fff2, #ffffff1a);
    color: #b5b2b2;
    border-radius: 5px;
    margin: 0 3px;
    font-size: 12px;
  }

  :deep(.el-tabs__item) {
    width: 70px;
    height: 35px;
    margin: 17px 5px 0 0 !important;
    font-size: 13px;
  }

  .play_info_right {
    display: flex;
    flex-direction: row;
  }

  .play_info_right a {
    margin-left: 5px;
    display: inline-block;
    padding: 2px 8px;
    font-size: 20px;
    height: 36px;
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 12px;
  }

  .play_info_right a:active {
    color: #FFBB5C;
    background: rgb(0, 0, 0, 0.2);
  }
}

/*pc端样式*/
@media (min-width: 768px) {
  .player_area {
    padding: 10px 6%;
  }

  .tags a {
    padding: 5px 10px;
    /*    background-color: rgba(155, 73, 231, 0.72);*/
    background: linear-gradient(#9B49E7B8, #9B49E799);
    color: #c4c2c2;
    font-size: 13px;
    border-radius: 6px;
    margin-right: 15px;
  }

  .tags span {
    padding: 6px 12px;
    /*background-color: #404042;*/
    background: linear-gradient(#fff2, #ffffff1a);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #b5b2b2;
    border-radius: 5px;
    margin: 0 8px;
    font-size: 12px;
  }

  .play_content a {
    white-space: nowrap;
    font-size: 12px;
    min-width: calc(10% - 24px);
    padding: 6px 10px;
    color: #ffffff;
    border-radius: 6px;
    margin: 8px 12px;
    background: #888888;
  }

  .play_info_right a:hover {
    color: #FFBB5C;
    background: rgb(0, 0, 0, 0.2);
  }

}
</style>





