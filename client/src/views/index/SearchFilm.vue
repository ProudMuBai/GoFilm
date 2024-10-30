<template>
    <div  class="container">
        <div class="search_group">
                <input v-model="data.search" @keydown="e=>{e.keyCode==13 && searchMovie()}" placeholder="输入关键字搜索 动漫,剧集,电影 " class="search"/>
                <el-button @click="searchMovie" :icon="Search"  style="" />
        </div>
        <div v-if="data.list && data.list.length > 0 " class="search_res">
            <div class="title">
                <h2>{{ data.oldSearch }}</h2>
                <p>共找到{{ data.page.total }}部与"{{ data.oldSearch }}"相关的影视作品</p>
            </div>
            <div class="content">
                <div class="film_item" v-for="m in data.list">
                    <a :href="`/filmDetail?link=${m.id}`" :style="{backgroundImage: `url('${m.picture}')`}"></a>
                    <div class="film_intro">
                        <h3>{{ m.name }}</h3>
                        <p class="tags">
                            <span class="tag_c">{{ m.cName }}</span>
                            <span>{{ m.year }}</span>
                            <span>{{ m.area }}</span>
                        </p>
                        <p><em>导演:</em>{{ m.director }}</p>
                        <p><em>主演:</em>{{ m.actor }}</p>
                        <p class="blurb"><em>剧情:</em>{{ m.blurb.replaceAll('　　', '') }}</p>
                        <el-button :icon="CaretRight" @click="play(m.id)">立即播放</el-button>
                    </div>
                </div>
            </div>
            <div class="pagination_container">
                <el-pagination background layout="prev, pager, next"
                               v-model:current-page="data.page.current"
                               @current-change="changeCurrent"
                               :pager-count="5"
                               :background="true"
                               :page-size="data.page.pageSize"
                               :total="data.page.total"
                               :prev-icon="ArrowLeftBold"
                               :next-icon="ArrowRightBold"
                               hide-on-single-page
                               class="pagination"/>
            </div>
        </div>
    </div>
    <el-empty v-if="data.oldSearch != '' && (!data.list || data.list.length == 0) " description="未查询到对应影片"/>
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
const searchMovie = ()=>{
    if (data.search.length <=0) {
        ElMessage.error({message: '搜索信息不能为空', duration:1000})
        return
    }
    location.href = location.href = `/search?search=${data.search}`
}

// 执行搜索请求
const refreshPage = (keyword: any, current: any) => {
    ApiGet('/searchFilm', {keyword: keyword, current: current}).then((resp: any) => {
      if (resp.code == 0) {
        data.list = resp.data.list
        data.page = resp.data.page
        data.oldSearch = keyword
      } else {
        ElMessage.warning({message: resp.msg, duration: 1000})
      }

    })
}

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

<!--移动端-->
<style scoped>
@import "/src/assets/css/pagination.css";
@media (max-width: 768px) {
    .title h2 {
        margin: 8px auto;
    }
    .film_item {
        flex-basis: calc(100% - 20px);
        margin: 0 10px 25px 10px;
        display: flex;
        background: #2e2e2e;
        padding: 10px;
        min-height: 180px;
        max-height: 200px;
        border-radius: 16px;
    }
    .film_item a {
        flex: 2;
        border-radius: 8px;
        background-size: cover;
    }

    .film_intro {
        max-width: 60%;
        margin-left: 10px;
        flex: 3;
        text-align: left;
        padding: 0 10px;
        font-size: 15px;
        position: relative;
    }

    .film_intro h3 {
        font-size: 16px;

        font-weight: bold;
    }

    .film_item h3, p, button {
        margin: 2px 0 2px 0;
    }

    .film_item p {
        max-width: 90%;
        display: -webkit-box;
        -webkit-line-clamp: 1;
        -webkit-box-orient: vertical;
        overflow: hidden;
        font-size: 13px;

    }

    .film_item p em {
        font-weight: bold;
        margin-right: 8px;
    }

    .film_item button {
        background-color: orange;
        border-radius: 20px;
        border: none !important;
        color: #ffffff;
        font-weight: bold;
        position: absolute;
        margin-bottom: 2px;
        bottom: 0;
    }
    .blurb{
        display: none!important;
    }
    .tags {
        display: flex;
        width: 90%;
        justify-content: space-between;
    }
    .tags .tag_c{
        background: rgba(155, 73, 231, 0.72);
    }
    .tags span {
        border-radius: 5px;
        padding: 3px 5px;
        background: rgba(66, 66, 66);
        color: #c9c4c4;
        margin-right: 5px;
    }

    .search_group {
        width: 80%;
        margin: 0 auto;
        display: flex;
    }
    .search {
        flex: 10;
        background-color: #2e2e2e!important;
        border: none!important;
        height: 40px;
        border-radius: 6px 0 0 6px;
        padding-left: 20px;
        color: #c9c4c4;
        font-size: 15px;
        font-weight: bold;
    }
    .search::placeholder{
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
        color: rgba(155,73,231,0.72);
        border: none!important;
        height: 40px;
        border-radius: 0 8px 8px 0;
        font-size: 20px;
        /*margin-bottom: 2px*/
    }

}

</style>
<!--pc端-->
<style scoped>

.title {
    margin-bottom: 20px;
}
.container {
    width: 100%;
}

.content {
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
}
.search_res {
    width: 100%;
}


@media (min-width: 768px) {
    .film_item {
        flex-basis: calc(50% - 18px);
        max-width: 50%;
        display: flex;
        background: #2e2e2e;
        padding: 16px;
        min-height: 250px;
        max-height: 280px;
        border-radius: 16px;
        margin-bottom: 25px;
    }
    .film_item a {
        flex: 1;
        border-radius: 8px;
        background-size: cover;
    }

    .film_intro {
        max-width: 75%;
        margin-left: 10px;
        flex: 3;
        /*flex-grow: 4;*/
        text-align: left;
        padding: 0 10px;
        font-size: 15px;
        position: relative;
    }

    .film_item h3, p, button {
        margin: 3px 0 12px 0;
    }

    .film_item p {
        max-width: 90%;
        display: -webkit-box;
        -webkit-line-clamp: 1;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .film_item p em {
        font-weight: bold;
        margin-right: 8px;
    }

    .film_item button {
        background-color: orange;
        border-radius: 20px;
        border: none !important;
        color: #ffffff;
        font-weight: bold;
        position: absolute;
        margin-bottom: 2px;
        bottom: 0;
    }

    .tags {
        display: flex;
        width: 90%;
        justify-content: space-between;
    }
    .tags .tag_c{
        background: rgba(155, 73, 231, 0.72);
    }
    .tags span {
        border-radius: 5px;
        padding: 3px 5px;
        background: rgba(66, 66, 66);
        color: #c9c4c4;
        margin-right: 10px;
    }

    .search_group {
        width: 45%;
        margin: 20px auto;
        display: flex;
    }
    .search {
        flex: 10;
        background-color: #2e2e2e!important;
        border: none!important;
        height: 40px;
        border-radius: 6px 0 0 6px;
        padding-left: 20px;
        color: #c9c4c4;
        font-size: 15px;
        font-weight: bold;
    }
    .search::placeholder{
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
        color: rgba(155,73,231,0.72);
        border: none!important;
        height: 40px;
        border-radius: 0 6px 6px 0;
        font-size: 20px;
    }

}





</style>

