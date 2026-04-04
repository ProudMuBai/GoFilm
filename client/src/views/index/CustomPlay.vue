<template>
  <!--自定义播放测试页面-->
  <div class="container">
    <!--头部工具栏-->
    <div class="player_header">
      <!--<el-input class="player_link" placeholder="请输入视频播放地址, mp4 或 m3u8 格式">
        <template #append>
          <el-button :icon="Search"/>
        </template>
      </el-input>-->
      <div class="player_link">
        <input type="text" v-model="data.link" @keyup.enter="play" placeholder="请输入视频播放地址, mp4 或 m3u8 格式" class="cus-input">
        <button class="iconfont icon-play" @click="play" />
      </div>
    </div>
    <!--播放器区域-->
    <div ref="playerContainer" class="player_area" />
  </div>
</template>

<script setup lang="ts">
import {Search} from "@element-plus/icons-vue";
import posterImg from "../../assets/image/play.png";
// import {VideoPlayer} from "@videojs-player/vue";
import {onMounted, reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import Player from "xgplayer"
import 'xgplayer/dist/index.min.css';
import HlsPlugin from 'xgplayer-hls'


const data = reactive({
  link: '',
  options: {
    // width: '100%',
    // height: '100%',
    title: "", //视频名称
    url: "https://vip.dytt-tvs.com/20260401/15005_f4adf97f/index.m3u8", //视频源
    volume: 0.6, // 音量
    currentTime: 50,
    autoplay: false,
  },
})

// 获取播放器渲染节点
const playerContainer = ref<HTMLDivElement | undefined>(undefined)
// 视频组件实例化
const playerInstance = ref()


onMounted(() => {
  playerInstance.value = new Player({
    el: playerContainer.value,
    url: data.options.url,
    poster: posterImg,
    width: "",
    height: "",
    autoplay: data.options.autoplay,
    lang: 'zh-cn', // 设置语言为中文
    volume: 0.7,   // 初始音量
    playbackRate:[3, 2, 1.5, 1, 0.75, 0.5],
    playsinline: true,
    plugins:[HlsPlugin],
    hls: {
      retryCount: 3, // 重试 3 次，默认值
      retryDelay: 1000, // 每次重试间隔 1 秒，默认值
      loadTimeout: 10000, // 请求超时时间为 10 秒，默认值
      fetchOptions: {mode: 'cors'},
      targetLatency: 10, // 直播目标延迟，默认 10 秒
      maxLatency:  20, // 直播允许的最大延迟，默认 20 秒
      disconnectTime: 0, // 直播断流时间，默认 0 秒，（独立使用时等于 maxLatency）
      // preloadTime: 30 // 默认值
    },
    controls: {
      autoHide: true
    },
    keyboard: { playbackRate: 3 },
    mobile: {
      controls: true,
      rotateFullScreen: true,
      playsinline: true,
      hideDefaultControls: true,
      pressRate: 3,//长按倍速
      disablePress: false,
    }
  })
})


// 播放执行
const play = () => {
  const pattern = /(^http[s]?:\/\/[^\s]+\.m3u8$)|(^http[s]?:\/\/[^\s]+\.mp4$)/
  if (!pattern.test(data.link)) {
    ElMessage.error({message: '视频链接格式异常, 请输入正确的播放链接!!!'})
    return
  }
  // 同步 link 为 player src
  data.options.url = data.link
  playerInstance.value.src = data.link
  playerInstance.value.load()
  playerInstance.value.play()
  // 执行播放器播放
}

</script>

<style scoped>
:deep(.el-main) {
  padding-bottom: 0!important;
}
.container {
  margin: 0 auto;
  height: 80%
}



/*========================================================================================================================*/

.player_header {
  margin: 40px auto;
}

.player_link {
  width: 80%;
  height: 45px;
  margin: 0 auto;
  display: flex;
}
.cus-input {
  font-size: 16px;
  width: 100%;
  padding: 0  40px;
  border: none;
  outline-style: none;
  border-radius: 16px 0 0 16px;
  min-height: 40px;
  background: rgba(255,255,255,0.8);
}
.cus-input:focus{
  border: 0;
}

.icon-play {
  height: 100%;
  font-size: 16px;
  border-radius: 0 16px 16px 0;
  background: deeppink;
  color: rgba(255,255,255,0.8);
  outline-style: none;
}
.icon-play:hover {
  background: hotpink;
}

/*播放器样式.*/
.player_area {
  width: 92%;
  margin: 0 auto;
  aspect-ratio: 16 / 9;
/*  padding-bottom: 56.25% !important;*/
  position: relative;
  border-radius: 6px;
  display: flex;
}

/*
.video-player {
  width: 80% !important;
  height: 80% !important;
  left: 10%;
  position: absolute;
  border-radius: 6px;
}

:deep(.vjs-big-play-button) {
  line-height: 2em;
  height: 2em;
  width: 2em;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.65);
}

:deep(.vjs-control-bar) {
  background: rgba(0, 0, 0, 0.32);
}

!*取消video被选中的白边*!
:deep(video:focus) {
  border: none !important;
  outline: none;
}

:deep(.data-vjs-player:focus) {
  border: none !important;
  outline: none;
}

:deep(.vjs-tech) {
  border-radius: 6px;
}

:deep(img) {
  border-radius: 6px;
}

!*进度条配色*!
:deep(.video-js .vjs-load-progress div) {
  background: rgba(255, 255, 255, 0.55) !important;
}

:deep(.video-js .vjs-play-progress) {
  background: #44c8cf;
}

:deep(.video-js .vjs-slider) {
  background-color: hsla(0, 0%, 100%, .2);
}
*/

</style>