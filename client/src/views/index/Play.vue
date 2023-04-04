<template>

    <div class="player_area">
        <!--    视频播放区域-->
        <div class="player_p">
            <iframe ref="iframe" class="player" :src="data.current.link"
                    :name="data.detail.name"
                    marginheight="0"
                    marginwidth="0"
                    framespacing="0"
                    vspale="0"
                    frameborder="0" allowfullscreen="true" scolling="no"
                    sandbox="allow-scripts allow-same-origin allow-downloads">
            </iframe>
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
            <el-tabs type="card" v-model="data.currentTabName" class="plya_tabs" @tab-change="changeSource">
                <el-tab-pane v-for="(p,i) in data.detail.playList"
                             :name="`tab-${i}`"
                             :label="`播放地址${i+1}`">
                    <div class="play_content">
                        <a v-for="(item,index) in p" href="javascript:void(false)" @click="playChange(item)"
                           :class="data.current.link.search(item.link) !== -1?'play_active':''">{{ item.episode }}</a>
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

import {onMounted, reactive, ref, withDirectives} from "vue";
import {useRouter} from "vue-router";
import {ApiGet} from "../../utils/request";
import RelateList from "../../components/RelateList.vue";
import {Promotion} from "@element-plus/icons-vue";

// 播放页所需数据
const data = reactive({
    detail: {descriptor: {}, playList: [[{episode: '', link: ''}]]},
    current: {episode: '', link: ''},
    currentTabName: '',
    currentPlayFrom: 0,
    currentEpisode: 0,
    relate: [],

})

// 获取路由信息
const router = useRouter()
onMounted(() => {
    let query = router.currentRoute.value.query
    ApiGet(`/filmPlayInfo`, {id: query.id, playFrom: query.source, episode: query.episode}).then((resp: any) => {
        if (resp.status === 'ok') {
            data.detail = resp.data.detail
            resp.data.current.link = converLink(resp.data.current.link)
            data.current = resp.data.current
            data.currentPlayFrom = resp.data.currentPlayFrom
            data.currentEpisode = resp.data.currentEpisode
            data.relate = resp.data.relate
            // 设置当前选中的播放源
            data.currentTabName = `tab-${query.source}`
        }
    })
})


// ===============================视频播放处理=======================================
// 视频解析接口地址, 默认使用第一个
const resolver = [
    // m3u8使用此解析
    'https://jx.jsonplayer.com/player/?url=',
    //   'https://jx.m3u8.tv/jiexi/?url=',


    //   'https://jx.jsonplayer.com/player/?url=',
    //   'https://vip.bljiex.com/?url=',
    // 'https://jx.bozrc.com:4433/player/?url=',
    // html视频使用此解析
    'http://www.82190555.com/index/qqvod.php?url=',
    // 'https://jx.bozrc.com:4433/player/?url=',
    // 'https://vip.bljiex.com/?url=',

    // Google上随便找的
    'https://vip.bljiex.com/?url=',
    'https://jx.kingtail.xyz/?url=',
    'http://www.82190555.com/index/qqvod.php?url=',
    'https://www.nxflv.com/?url=',
    'http://www.wmxz.wang/video.php?url=',
    'https://www.feisuplayer.com/m3u8/?url=',
    // tampermonkey 脚本使用的解析
    'https://jx.bozrc.com:4433/player/?url=',
    'https://z1.m1907.top/?jx=',
    'https://jx.aidouer.net/?url=',
    'https://www.gai4.com/?url=',
    'https://okjx.cc/?url=',
    'https://jx.rdhk.net/?v=',
    'https://jx.blbo.cc:4433/?url=',
    'https://jsap.attakids.com/?url=',
    'https://jx.dj6u.com/?url=',
]

// 添加视频解析前缀
const converLink = (link: string): string => {
    // 视频统一使用第三方解析
    if (link.search("m3u8") != -1) {
        return `${resolver[0] + link}`
    }
    // return `${resolver[1]+link}`
    return `${link}`
}

// 点击播放对应影片
const playChange = (info: { link: string, episode: string }) => {
    // 判断是否是m3u8播放器, 如果是则添加前缀
    data.current.link = converLink(info.link)
    data.current.episode = info.episode
}
// 点击播放源标签事件
const changeSource = (tabName: any) => {
    data.currentTabName = tabName
    data.detail.playList.find((item, index) => {
        if (tabName.split("-")[1] - index == 0) {
            item.find(i => {
                if (i.episode == data.current.episode) {
                    data.current.link = converLink(i.link)
                }
            })
        }
    })
}

// 测试滑动音量调节
const iframe = ref()
// document.querySelector('#brightnessSlider').addEventListener('change', function() {
//     // 获取滑块的值
//     var brightness = this.value;
//
//     // 设置亮度为滑块的值
//     iframe.style.filter = 'brightness(' + brightness + '%)';
// })

</script>

<style scoped>
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
    /*height: 1400px;*/
    /*background: red;*/

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

}


.player_p {
    width: 100%;
    min-height: 200px;
    margin: 0;
    padding-bottom: 56.25% !important;
    /*padding-bottom: 42.25% !important;*/
    position: relative;
    background-image: url("/src/assets/image/play.png");
    background-size: cover;
    border-radius: 6px;
}

iframe {
    border-radius: 6px;
    left: 0;
    width: 100%; /* 设置iframe元素的宽度为父容器的100% */
    height: 100%; /* 设置iframe元素的高度为0，以便自适应高度 */
    /*padding-bottom: 56.25%; !* 使用padding-bottom属性计算iframe元素的高度，这里假设视频的宽高比为16:9 *!*/
    /*border: none; !* 去除iframe元素的边框 *!*/
    /*transform: scale(1);*/
    position: absolute;

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
    padding: 10px 10px 10px 18px;

}

.play_list > h2 {
    position: absolute;
    left: 10px;
    top: -10px;
    z-index: 50;
}

.play_content a {
    font-size: 12px;
    min-width: 65px;
    padding: 6px 15px;
    color: #ffffff;
    border-radius: 6px;
    margin: 8px 8px;
    background: #888888;
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
    .tags span {
        padding: 6px 10px;
        background-color: #404042;
        color: #b5b2b2;
        border-radius: 5px;
        margin: 0 3px;
        font-size: 12px;
    }

}
</style>