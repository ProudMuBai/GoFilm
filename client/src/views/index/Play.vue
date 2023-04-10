<template>
    <div class="player_area" v-show="data.loading">
        <div class="player_p">
            <video-player  @mounted="handleBtn" :src="data.options.src" poster="https://s2.loli.net/2023/04/10/H23F6qDXlnwmPG9.png" controls :loop="false"
                          @keydown="handlePlay"
                          :volume="data.options.volume"
                          crossorigin="anonymous" playsinline class="video-player"
                          :playback-rates="[0.5, 1.0, 1.5, 2.0]" />
        </div>
        <div class="current_play_info">
            <div class="play_info_left">
                <h3 class="current_play_title">{{ `${data.detail.name}&emsp;${data.current.episode}` }}</h3>
                <div class="tags">
                    <b>
                        <el-icon>
                            <Promotion/>
                        </el-icon>
                        {{ data.detail.descriptor.cName }}</b>
                    <span>{{ data.detail.descriptor.classTag }}</span>
                    <span>{{ data.detail.descriptor.year }}</span>
                    <span>{{ data.detail.descriptor.area }}</span>
                </div>
            </div>

        </div>
        <!-- 播放选集   -->
        <div class="play_list">
            <h2 class="hidden-md-and-down">播放列表:(右侧切换播放源)</h2>
            <el-tabs type="card" v-model="data.currentTabName" class="plya_tabs">
                <el-tab-pane v-for="(p,i) in data.detail.playList"
                             :name="`tab-${i}`"
                             :label="`播放列表${i+1}`">
                    <div class="play_content">
                        <a v-for="(item,index) in p" href="javascript:void(false)"
                           @click="playChange({sourceIndex: i, episodeIndex: index, target: this})"
                           :class="item.link == data.current.link?'play_active':''">{{ item.episode }}</a>

                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
        <div class="correlation">
            <RelateList :relate-list="data.relate"/>
        </div>
    </div>

</template>

<script lang="ts" setup>

import {onBeforeMount, onMounted, reactive, ref, withDirectives} from "vue";
import {useRouter} from "vue-router";
import {ApiGet} from "../../utils/request";
import {ElMessage} from "element-plus";
import RelateList from "../../components/RelateList.vue";
import {Promotion} from "@element-plus/icons-vue";
// 引入视频播放器组件
import {VideoPlayer} from '@videojs-player/vue'
import 'video.js/dist/video-js.css'
const handlePlay = (e:any)=>{
    e.preventDefault()
    switch (e.keyCode) {
        case 32:
            console.log(e.target.paused)
            if(e.target.paused) {
                e.target.play()
            } else {
                e.target.pause()
            }
            break
        case 37:
            e.target.currentTime = e.target.currentTime-5 < 0 ? 0 : e.target.currentTime-5
            break
        case 39:
            e.target.currentTime = e.target.currentTime+5 > e.target.duration ? e.target.duration : e.target.currentTime+5
            break
        case 38:
            data.options.volume = data.options.volume + 0.05 > 1 ? 1 : data.options.volume + 0.05
            break
        case 40:
            data.options.volume = data.options.volume - 0.05 < 0 ? 0 : data.options.volume - 0.05
            break
    }

}
// 主动触发快捷键
const tiggerKeyMap = (keyCode:number)=>{
    let player = document.getElementsByTagName("video")[0]
    player.focus()
    const event = document.createEvent('HTMLEvents');
    event.initEvent('keydown', true, false);
    event.keyCode = keyCode; // 设置键码
    player.dispatchEvent(event)
}
const handleBtn = (e:any)=> {
    let btns = document.getElementsByClassName('vjs-button')
    for (let el of btns) {
        el.addEventListener('keydown', function (t:any){
            t.preventDefault()
            tiggerKeyMap(t.keyCode)
        })

    }
}


// 播放页所需数据
const data = reactive({
    loading: false,
    detail: {descriptor: {}, playList: [[{episode: '', link: ''}]]},
    current: {index: 0, episode: '', link: ''},
    currentTabName: '',
    currentPlayFrom: 0,
    currentEpisode: 0,
    relate: [{}],
// @videojs-player 播放属性设置
    options: {
        title: "", //视频名称
        src: "", //视频源
        volume: 0.6, // 音量
        currentTime: 0,
    }

})

// 获取路由信息
const router = useRouter()
onBeforeMount(() => {
    let query = router.currentRoute.value.query
    ApiGet(`/filmPlayInfo`, {id: query.id, playFrom: query.source, episode: query.episode}).then((resp: any) => {
        if (resp.status === 'ok') {
            data.detail = resp.data.detail
            data.current = {index: resp.data.currentEpisode, ...resp.data.current}
            data.currentPlayFrom = resp.data.currentPlayFrom
            data.currentEpisode = resp.data.currentEpisode
            data.relate = resp.data.relate
            // 设置当前选中的播放源
            data.currentTabName = `tab-${query.source}`
            // 设置当前的视频播放url
            data.options.src = data.current.link
            data.loading = true
        } else {
            ElMessage.error("影片信息加载失败,请尝试刷新页面!!!")
        }
    })
})

// 点击播集数播放对应影片
const playChange = (play: { sourceIndex: number, episodeIndex: number, target: any }) => {
    let currPlay = data.detail.playList[play.sourceIndex][play.episodeIndex]
    data.current = {index: play.episodeIndex, episode: currPlay.episode, link: currPlay.link}
    data.options.src = currPlay.link
    data.options.title = data.detail.name + "  " + currPlay.episode
}


</script>

<style scoped>
/*vue3-video-play 相关设置*/
/*//播放器控件区域大小*/
.video-player {
    width: 100% !important;
    height: 100% !important;
    position: absolute;
    border-radius: 6px;

}

:deep(.vjs-big-play-button){
    line-height: 2em;
    height: 2em;
    width: 2em;
    border-radius: 50%;
    border: none;
    background: rgba(0,0,0,0.65);
}
:deep(.vjs-control-bar){
    background: rgba(0,0,0,0.32);
}
/*取消video被选中的白边*/
:deep(video:focus) {
    border: none!important;
    outline: none;
}
:deep(.data-vjs-player:focus){
    border: none!important;
    outline: none;
}
:deep(.vjs-tech){
    border-radius: 6px;
}

:deep(img) {
    border-radius: 6px;
}

/*当前播放的影片信息展示*/
.current_play_info {
    width: 100%;
    padding: 15px 5px;
    text-align: left;
}

.current_play_title {
    font-weight: 500;
    color: rgb(201, 196, 196);
    margin: 0 0 5px 0;
}


/* 播放区域*/
.player_area {
    width: 100%;
    min-height: 100%;
}


@media (min-width: 650px) {
    .player_area {
        padding: 10px 6%;
    }

    .tags b {
        padding: 5px 10px;
        background-color: rgba(155, 73, 231, 0.72);
        font-size: 13px;
        border-radius: 6px;
        margin-right: 15px;
    }

    .tags span {
        padding: 6px 12px;
        background-color: #404042;
        color: #b5b2b2;
        border-radius: 5px;
        margin: 0 8px;
        font-size: 12px;
    }

    .play_content a {
        font-size: 12px;
        flex-basis: calc(10% - 24px);
        padding: 6px 10px;
        color: #ffffff;
        border-radius: 6px;
        margin: 8px 12px;
        background: #888888;
    }

}


.player_p {
    width: 100%;
    /*height: 700px;*/
    margin: 0;
    padding-bottom: 56.25% !important;
    position: relative;
    border-radius: 6px;
    display: flex;
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


/*集数选中效果*/
.play_active {
    color: orange !important;
    background: #424242 !important;
}

.play_content .play_tabs {
    background: #2e2e2e;
}

:deep(.el-tabs__nav-scroll) {
    display: flex;
    justify-content: end;
}

:deep(.el-tabs__header) {
    /*border-bottom: 1px solid #888888!important;*/
    margin-bottom: 0;
    border-bottom: none !important;
    background: rgb(34, 34, 34);
    height: 50px !important;
}

:deep(.el-tabs__nav) {
    border: none !important;
}

:deep(.el-tabs__item.is-active) {
    color: #ee9600;
}

:deep(.el-tabs__item:hover) {
    color: orange;
    background: #484646;
}

:deep(.el-tabs__item) {
    height: 50px;
    line-height: 50px;
    margin-left: 2px;
    border-radius: 8px 8px 0 0;
    border: none !important;
    color: #ffffff;
    background: #2e2e2e;
}

/*推荐列表区域*/
.correlation {
    width: 100%;
}

</style>

<!--移动端-->
<style scoped>

/*适应小尺寸*/
@media (max-width: 650px) {
    .player_area {
        padding: 5px 10px;
    }

    .tags b {
        padding: 5px 10px;
        background-color: rgba(155, 73, 231, 0.72);
        font-size: 13px;
        border-radius: 6px;
        margin-right: 3px;
    }


    .play_content a {
        color: #ffffff;
        border-radius: 6px;
        margin: 6px 8px;
        background: #888888;
        flex-basis: calc(25% - 16px);
        font-size: 12px;
        padding: 6px 12px !important;
    }

    .tags span {
        padding: 6px 10px;
        background-color: #404042;
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

}
</style>