<template>
    <div class="film" v-show="data.loading">
        <!-- hidden-sm-and-up 移动端title   -->
        <div class="hidden-sm-and-up">
            <div class="title_mt  ">
                <a class="picture_mt" href="" :style="{backgroundImage: `url('${data.detail.picture}')`}"></a>
                <div class="title_mt_right">
                    <h3>{{ data.detail.name }}</h3>
                    <ul class="tags">
                        <li style="margin: 2px 0">{{ data.detail.descriptor.classTag?`${data.detail.descriptor.classTag}`.replaceAll(",", " | "):'未知' }}</li>
                    </ul>
                    <p><span>导演:</span> {{data.detail.descriptor.director }}</p>
                    <p><span>主演:</span> {{handleLongText(data.detail.descriptor.actor)}}</p>
                    <p><span>上映:</span> {{ data.detail.descriptor.releaseDate }}</p>
                    <p><span>地区:</span> {{ data.detail.descriptor.area }}</p>
                    <p v-if="data.detail.descriptor.remarks"><span>连载:</span>{{ data.detail.descriptor.remarks }}</p>
                    <!--<p><span>评分:</span><b id="score">{{ data.detail.descriptor.dbScore }}</b></p>-->
                </div>
            </div>
            <div class="mt_content">
                <p v-html="`${data.detail.descriptor.content}`.replaceAll('　　', '')"></p>
            </div>
        </div>
        <!-- pc端title-->
        <div class="title hidden-sm-and-down ">
            <a class="picture" href="" :style="{backgroundImage: `url('${data.detail.picture}')`}"></a>
            <h2>{{ data.detail.name }}</h2>
            <ul class="tags">
                <li class="t_c">
                    <el-icon>
                        <Promotion/>
                    </el-icon>
                    {{ data.detail.descriptor.cName }}
                </li>
                <li v-if="data.detail.descriptor.classTag">{{ `${data.detail.descriptor.classTag}`.replaceAll(",", "&emsp;") }}</li>
                <li>{{ data.detail.descriptor.year }}</li>
                <li>{{ data.detail.descriptor.area }}</li>
            </ul>
            <p><span>导演:</span> {{ data.detail.descriptor.director }}</p>
            <p><span>主演:</span> {{ data.detail.descriptor.actor }}</p>
            <p><span>上映:</span> {{ data.detail.descriptor.releaseDate }}</p>
            <p v-if="data.detail.descriptor.remarks"><span>连载:</span>{{ data.detail.descriptor.remarks }}</p>
            <p><span>评分:</span><b id="score">{{ data.detail.descriptor.dbScore }}</b></p>
            <div class="cus_wap">
                <p style="min-width: 40px"><span>剧情:</span></p>
                <p ref="textContent" class="text_content">
                    <el-button v-if="`${data.detail.descriptor.content}`.length > 140" class="multi_text"
                               style="color:#a574b7;"
                               @click="showContent(multiBtn.state)" link>{{ multiBtn.text }}
                    </el-button>
                    <span class="cus_info" v-html="data.detail.descriptor.content"></span>
                </p>
            </div>
            <p>
                <el-button type="warning" class="player" size="large" @click="play({episode:0,source:0})" round>
                    <el-icon>
                        <CaretRight/>
                    </el-icon>
                    立即播放
                </el-button>
            </p>
        </div>
        <!--播放列表-->
        <div class="play_list">
            <h2 class="hidden-md-and-down">播放列表:(右侧切换播放源)</h2>
            <el-tabs type="card" class="play_tabs">
                <el-tab-pane v-for="(p,i) in data.detail.playList" :label="`播放地址${i+1}`">
                    <div class="play_content">
                        <a v-for="(item,index) in p" href="javascript:;"
                           @click="play({source: i, episode: index})">{{ item.episode }}</a>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
        <!--相关系列影片-->
        <div class="correlation">
            <RelateList :relate-list="data.relate"/>
        </div>
    </div>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import {onBeforeMount, reactive, ref,} from "vue";
import {ApiGet} from "../../utils/request";
import {ElMessage, ElLoading} from 'element-plus'
import {Promotion, CaretRight} from "@element-plus/icons-vue";
import RelateList from "../../components/RelateList.vue";

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
        descriptor: {
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
        }
    },
    relate: [],
    loading: false,
})

// 对部分信息过长进行处理
const handleLongText = (t:string):string=>{
    let res = ''
    t.split(',').forEach((s,i)=>{
        console.log(s)
        if (i <3) {
            res += `${s} `
        }
    })
    return res.trimEnd()
}

onBeforeMount(() => {
    let link = router.currentRoute.value.query.link
    ApiGet('/filmDetail', {id: link}).then((resp: any) => {
        if (resp.status === "ok") {
            data.detail = resp.data.detail
            // 去除影视简介中的无用内容和特殊标签格式等
            data.detail.name = data.detail.name.replace(/(～.*～)/g, '')
            data.detail.descriptor.content = data.detail.descriptor.content.replace(/(&.*;)|( )|(　　)|(\n)|(<[^>]+>)/g, '')
            data.relate = resp.data.relate
            // 处理过长数据
            data.detail.descriptor.actor = handleLongText(data.detail.descriptor.actor)
            console.log(handleLongText(data.detail.descriptor.actor))
            data.detail.descriptor.director = handleLongText(data.detail.descriptor.director)
            data.loading = true
        } else {
            ElMessage({
                type: "error",
                dangerouslyUseHTMLString: true,
                message: resp.message,
            })
        }
    })

})

// 选集播放点击事件
const play = (change: { source: number, episode: number }) => {
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

</script>

<!--移动端适配-->
<style scoped>
@media (max-width: 650px) {
    .title_mt {
        padding: 0 3px;
        display: flex;
        flex-direction: row;
        flex-flow: nowrap;
        overflow: hidden;
    }

    .picture_mt {
        width: 100%;
        height: 180px;
        margin-right: 12px;
        border-radius: 8px;
        background-size: cover;
        flex: 1;
    }

    .title_mt_right {
        flex: 1.5;
        text-align: left;
    }

    .title_mt_right h3 {
        font-size: 14px;
        margin: 0 0 5px 0;
    }

    .title_mt_right p {
        font-size: 12px;
        margin: 3px 2px;
        white-space: nowrap;
    }

    .mt_content {
        margin-top: 5px;
        border-top: 1px solid #777777;
        border-bottom: 1px solid #777777;
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

    :deep(.el-tabs__item) {
        width: 70px;
        height: 35px !important;
        margin: 17px 5px 0 0 !important;
        font-size: 13px;
    }
}
</style>

<style scoped>
.correlation {
    width: 100%;
}

.film {
    width: 100%;
    padding: 0 1%;
}

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
    padding: 20px 10px 20px 18px;

}

.play_list > h2 {
    position: absolute;
    left: 10px;
    top: -10px;
    z-index: 50;
}

@media (min-width: 650px) {
    .play_content a {
        white-space: nowrap;
        font-size: 12px;
        min-width: calc(10% - 24px);
        padding: 6px 15px;
        color: #ffffff;
        border-radius: 6px;
        margin: 8px 12px;
        background: #888888;
    }

}


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


/*顶部影片信息显示区域*/

.title {
    width: 100%;
    border-radius: 10px;
    background: #2e2e2e;
    padding: 5px 30px 30px 30px;
    position: relative;
}

.title > h2 {
    text-align: left;
    color: #888888;
}

.picture {
    position: absolute;
    width: 190px;
    height: 270px;
    right: 30px;
    top: 30px;
    border-radius: 8px;
    background-size: cover;

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
    background: rgba(66, 66, 66);
    margin: 0 8px;
    font-size: 12px;
    color: #888888;
}

.tags > .t_c {
    background: rgba(155, 73, 231, 0.72);
    color: #c4c2c2;
    margin-left: 0;
}

.title p {
    text-align: left;
    font-size: 14px;
    margin: 20px 0;
    max-width: 60%;
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.title p span {
    font-size: 15px;
    font-weight: bold;
    color: #888888;
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

</style>