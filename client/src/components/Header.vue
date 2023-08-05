<template>
    <div class="header">
        <!-- 左侧logo以及搜索 -->
        <div class="nav_left">
            <!--        <img class="logo" src="/src/assets/logo.png">-->
            <a href="/" class="site" >GoFilm</a>
            <div class="search_group">
                <input v-model="keyword" @keydown="(e)=>{e.keyCode == 13 && searchFilm()}" placeholder="搜索 动漫,剧集,电影 " class="search"/>
                <el-button @click="searchFilm" :icon="Search"/>
            </div>
        </div>
        <!--右侧顶级分类导航 -->
        <div class="nav_right">
            <el-link :underline="false" href="/">首页</el-link>
            <el-link :underline="false" :href="`/filmClassify?Pid=${nav.film.id}`">电影</el-link>
            <el-link :underline="false" :href="`/filmClassify?Pid=${nav.tv.id}`">剧集</el-link>
            <el-link :underline="false" :href="`/filmClassify?Pid=${nav.cartoon.id}`">动漫</el-link>
            <el-link :underline="false" :href="`/filmClassify?Pid=${nav.variety.id}`">综艺</el-link>
            <!--        <span style="color:#777; font-weight: bold">|</span>-->
            <el-link href="/search" class="hidden-md-and-up" :underline="false">
                <el-icon style="font-size: 18px">
                    <Search/>
                </el-icon>
            </el-link>
        </div>
    </div>

</template>

<script lang="ts" setup>
import {onMounted, reactive, ref, watch} from "vue";
import {useRouter} from "vue-router";
import {Search, CircleClose} from '@element-plus/icons-vue'
import {ElMessage} from "element-plus";
import {ApiGet} from "../utils/request";

//
const keyword = ref<string>('')


// 从父组件获取当前路由对象
const router = useRouter()
// 影片搜索
const searchFilm = () => {
    if (keyword.value.length <= 0) {
        ElMessage.error({message: "请先输入影片名称关键字再进行搜索", duration: 1500})
        return
    }
    location.href = `/search?search=${keyword.value}`
    // router.push({path: '/search', query:{search: keyword.value}, replace: true})
}

// 导航栏挂载完毕时发送一次请求拿到对应的分类id数据
const nav = reactive({
    cartoon: {},
    film: {},
    tv: {},
    variety: {},
})
onMounted(() => {
    ApiGet('/navCategory').then((resp: any) => {
        if (resp.status === 'ok') {
            nav.tv = resp.data.tv
            nav.film = resp.data.film
            nav.cartoon = resp.data.cartoon
            nav.variety = resp.data.variety
        } else {
            ElMessage.error({message: "请先输入影片名称关键字再进行搜索", duration: 1000})
        }
    })
})


</script>


<!--移动端适配-->
<style>
/*小尺寸时隐藏状态栏*/
@media (max-width: 768px) {
    .nav_right {
        display: flex;
        justify-content: space-between;
        /*display: none!important;*/
    }

    .nav_right a {
        color: #ffffff;
        flex-basis: calc(19% - 5px);
        padding: 0 10px;
        line-height: 40px;
        /*border-radius: 5px;*/
        /*border: 1px solid rebeccapurple;*/

    }

    .nav_right a:hover {
        color: #ffffff;
        /*background-color: transparent;*/
    }

    .header {
        width: 100% !important;
        height: 40px;
        background: radial-gradient(circle, #d275cd, rgba(155, 73, 231, 0.72), #4ad1e5);
    }

    .nav_left {
        display: none !important;
        width: 90% !important;
        margin: 0 auto;
    }
}
</style>

<style scoped>


@media (min-width: 768px) {
    .header {
        width: 78%;
        z-index: 0;
        max-height: 40px;
        line-height: 60px;
        margin: 0 auto;
        display: flex;
        justify-content: space-between;
    }

    .nav_left {
        display: flex;
    }
    /*site标志样式*/
    .site{
        font-weight: 600;
        font-style: italic;
        font-size: 24px;
        margin-right: 5px;
        background: linear-gradient(118deg, #e91a90, #c965b3, #988cd7, #00acfd);
        -webkit-background-clip: text;
        background-clip: text;
        color: transparent;
    }

    /*搜索栏*/
    .search_group {
        width: 80%;
        margin: 10px auto;
        display: flex;
    }

    .search {
        flex: 10;
        background-color: #2e2e2e !important;
        border: none !important;
        height: 40px;
        border-radius: 6px 0 0 6px;
        padding-left: 20px;
        color: #c9c4c4;
        font-size: 15px;
        font-weight: bold;
        line-height: 60px;
    }

    .search::placeholder {
        font-size: 15px;
        color: #999999;
    }

    .search:focus {
        outline: none;
    }

    .search_group button {
        flex: 1;
        margin: 0;
        background-color: #2e2e2e;
        color: rgb(171, 44, 68);
        border: none !important;
        height: 40px;
        border-radius: 0 6px 6px 0;
        font-size: 20px;
        /*margin-bottom: 2px*/
    }

    .nav_right {
        display: flex;
        height: 60px;
        flex-direction: row;
    }

    .nav_right a {
        min-width: 60px;
        height: 40px;
        line-height: 60px;
        margin: 10px 10px;
        font-size: 15px;
        text-align: center;
        font-weight: bold;
    }

    .nav_right a:hover {
        color: orange;
    }

    .logo {
        height: 40px;
        margin-top: 10px;
    }

}


</style>