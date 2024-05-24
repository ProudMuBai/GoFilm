<template>
  <div class="container">
    <!--头部工具栏-->
    <div class="player_header">
      <div class="player_link">
        <input type="text" v-model="data.link" @keyup.enter="play" placeholder="请输入视频播放地址, mp4 或 m3u8 格式" class="cus-input">
        <button class="iconfont icon-play" @click="play" />
      </div>
    </div>
    <!--播放器区域-->
    <div class="player_area">
      <video-player @mounted="playerMount" :src="data.options.src" :poster="posterImg" controls
                    :loop="false"
                    @keydown="handlePlay"
                    :bufferedPercent="30"
                    :volume="data.options.volume"
                    crossorigin="anonymous" playsinline class="video-player"
                    :playback-rates="[0.5, 1.0, 1.5, 2.0]"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import posterImg from "../../assets/image/play.png";
import { VideoPlayer } from "@videojs-player/vue";
import { reactive } from "vue";
import { ElMessage } from "element-plus";

const data = reactive({
  link: '',
  options: {
    title: "", //视频名称
    src: "", //视频源
    volume: 0.6, // 音量
    currentTime: 50,
  },
});

const handlePlay = (e: any) => {
  e.preventDefault();
  switch (e.keyCode) {
    case 32:
      if (e.target.paused) {
        e.target.play();
      } else {
        e.target.pause();
      }
      break;
    case 37:
      e.target.currentTime = e.target.currentTime - 5 < 0 ? 0 : e.target.currentTime - 5;
      break;
    case 39:
      e.target.currentTime = e.target.currentTime + 5 > e.target.duration ? e.target.duration : e.target.currentTime + 5;
      break;
    case 38:
      data.options.volume = data.options.volume + 0.05 > 1 ? 1 : data.options.volume + 0.05;
      break;
    case 40:
      data.options.volume = data.options.volume - 0.05 < 0 ? 0 : data.options.volume - 0.05;
      break;
  }
};

const triggerKeyMap = (keyCode: number) => {
  let player = document.getElementsByTagName("video")[0];
  player.focus();
  const event = document.createEvent('HTMLEvents');
  event.initEvent('keydown', true, false);
  event.keyCode = keyCode; // 设置键码
  player.dispatchEvent(event);
};

const handleBtn = (e: any) => {
  let btns = document.getElementsByClassName('vjs-button');
  for (let el of btns) {
    el.addEventListener('keydown', function (t: any) {
      t.preventDefault();
      triggerKeyMap(t.keyCode);
    });
  }
};

const playerMount = (e: any) => {
  handleBtn(e);
};

const play = () => {
  const pattern = /(^http[s]?:\/\/[^\s]+\.m3u8$)|(^http[s]?:\/\/[^\s]+\.mp4$)/;
  if (!pattern.test(data.link)) {
    ElMessage.error({ message: '视频链接格式异常, 请输入正确的播放链接!!!' });
    return;
  }
  data.options.src = data.link;
  document.getElementsByTagName("video")[0].play();
};
</script>

<style scoped>
:deep(.el-main) {
  padding-bottom: 0 !important;
}
.container {
  margin: 0 auto;
  height: 80%;
}
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
  padding: 0 40px;
  border: none;
  outline-style: none;
  border-radius: 16px 0 0 16px;
  min-height: 40px;
  background: rgba(255, 255, 255, 0.8);
}
.cus-input:focus {
  border: 0;
}
.icon-play {
  height: 100%;
  font-size: 16px;
  border-radius: 0 16px 16px 0;
  background: deeppink;
  color: rgba(255, 255, 255, 0.8);
  outline-style: none;
}
.icon-play:hover {
  background: hotpink;
}
.player_area {
  width: 100%;
  margin: 0;
  padding-bottom: 56.25% !important;
  position: relative;
  border-radius: 6px;
  display: flex;
}
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
:deep(.video-js .vjs-load-progress div) {
  background: rgba(255, 255, 255, 0.55) !important;
}
:deep(.video-js .vjs-play-progress) {
  background: #44c8cf;
}
:deep(.video-js .vjs-slider) {
  background-color: hsla(0, 0%, 100%, .2);
}
</style>
