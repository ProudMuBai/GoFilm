<template>
    <div class="container">
        <div class="header">
            <p>{{ d.category.name }}</p>
            <div class="c_header">
                <!--默认全部-->
                <a :class="`nav ${d.cid == -1?'active':''}`" :href="`/categoryFilm?pid=${d.category.id}`">全部{{d.category.name}}</a>
                <!--子分类-->
                <a :class="`nav ${d.cid == c.id?'active':''}`" :href="`/categoryFilm?pid=${c.pid}&cid=${c.id}&current=1`"
                   v-for="c in d.category.children">{{ c.name }}</a>
            </div>
        </div>

        <!--影片列表展示-->
        <FilmList :list="d.list"/>
        <!--分页展示区域-->
        <div class="pagination_container ">
            <el-pagination background layout="prev, pager, next"
                           v-model:current-page="d.page.current"
                           @current-change="changeCurrent"
                           :pager-count="5"
                           :background="true"
                           :page-size="d.page.pageSize"
                           :total="d.page.total"
                           :prev-icon="ArrowLeftBold"
                           :next-icon="ArrowRightBold"
                           hide-on-single-page
                           class="pagination"/>
        </div>
    </div>
</template>

<script setup lang="ts">


import {onMounted, reactive} from "vue";
import {useRouter} from "vue-router";
import {ApiGet} from "../../utils/request";
import {ElMessage} from "element-plus";
import {ArrowRightBold, ArrowLeftBold} from '@element-plus/icons-vue'
import FilmList from "../../components/FilmList.vue";

// 页面所需数据
const d = reactive({
    category: {},
    list: [],
    page: {
        current: 0,
    },
    cid: -1

})
// 获取路由参数查询对应数据
const router = useRouter()

// 点击分页按钮事件 current-change
const changeCurrent = (currentVal: number) => {
    let query = router.currentRoute.value.query
    // router.push({path: '/categoryFilm', query:{pid: query.pid, cid: query.cid, current: currentVal}})
    if (query.cid && query.cid != "") {
        location.href = `/categoryFilm?pid=${query.pid}&cid=${query.cid}&current=${currentVal}`
    } else {
        location.href = `/categoryFilm?pid=${query.pid}&&current=${currentVal}`
    }
}

// 点击分类事件
// const changeCategory = (cid?: any) => {
//     let params = new URLSearchParams(location.search)
//     location.href = cid ? `/categoryFilm?pid=${params.get('pid')}&cid=${cid}&current=1` : `/categoryFilm?pid=${params.get('pid')}`
// }


const getFilmData = (param: any) => {
    ApiGet('/filmCategory', {pid: param.pid, cid: param.cid, current: param.current}).then((resp: any) => {
        if (resp.status === 'ok') {
            d.category = resp.data.category
            d.list = resp.data.list
            d.page = resp.page
            // 设置当前请求的分类id
            d.cid = param.cid? param.cid:-1
        } else {
            ElMessage.error({message: "请先输入影片名称关键字再进行搜索", duration: 1000})
        }
    })
}

onMounted(() => {
    let query = router.currentRoute.value.query
    getFilmData({pid: query.pid, cid: query.cid, current: query.current})
})

</script>

<!--移动端修改-->
<style scoped>
@media (max-width: 650px) {
    .container {
        padding: 0 10px;

    }

    /*顶部内容区域*/
    .header {
        width: 100%;
        margin-bottom: 100px;
        background: none !important;
    }

    .header p {
        text-align: left;
        font-weight: 600;
        font-size: 24px;
        color: #c9c4c4;
        margin-top: 0;
        padding-left: 10px;
        height: 100%;
        /*background: rgb(34, 34, 34);*/
    }

    .c_header {
        max-width: 100%;
        display: flex;
        /*justify-content: start;*/
        flex-wrap: nowrap;
        overflow-x: scroll;
        margin-bottom: 20px;
    }

    .c_header a {
        white-space: nowrap;
        margin-right: 10px;
        color: #000;
        font-weight: 400;
        background: rgba(255, 255, 255, 0.94);
        padding: 6px 10px;
        border-radius: 10px;
    }

    .c_header a:hover {
        color: orange;
    }

    .nav:before {
        width: 36px;
        height: 4px;
        background: orange;
        content: '';
        position: absolute;
        left: 35%;
        bottom: 12px;
        border-radius: 50px;
        transform: scaleX(0);
        transition: transform 0.5s ease-out;
    }

    .nav:hover:before {
        width: 36px;
        height: 4px;
        background: orange;
        content: '';
        position: absolute;
        left: 35%;
        bottom: 12px;
        border-radius: 50px;
        transform: scaleX(1);
    }

    .active {
        background: rgb(249 230 195) !important;
        color: #e52424 !important;
    }

}
</style>

<style scoped>
.container {
    max-width: 100vw;
}

@media (min-width: 650px) {

    /*顶部内容区域*/
    .header {
        width: 100%;
    }

    .header p {
        text-align: left;
        font-weight: 800;
        font-size: 32px;
        color: #c9c4c4;
        margin-top: 0;
        /*border-bottom: 2px solid #ffffff;*/
        /*padding-bottom: 30px;*/
    }

    .c_header {
        width: 100%;
        display: flex;
        justify-content: start;
        margin-bottom: 20px;
    }

    .c_header a {
        flex-basis: calc(14% - 16px);
        white-space: nowrap;
        margin-right: 20px;
        color: #000;
        font-weight: 800;
        background: rgba(255, 255, 255, 0.94);
        padding: 1.35% 0;
        border-radius: 10px;
        position: relative;
    }

    .c_header a:hover {
        color: orange;
    }

    .nav:before {
        width: 36px;
        height: 4px;
        background: orange;
        content: '';
        position: absolute;
        left: 35%;
        bottom: 12px;
        border-radius: 50px;
        transform: scaleX(0);
        transition: transform 0.5s ease-out;
    }

    .nav:hover:before {
        width: 36px;
        height: 4px;
        background: orange;
        content: '';
        position: absolute;
        left: 35%;
        bottom: 12px;
        border-radius: 50px;
        transform: scaleX(1);
    }

    .active {
        background: rgb(249 230 195) !important;
        color: #e52424 !important;
    }
}
</style>

<style scoped>
/*分页插件区域*/
.pagination_container {
    width: 100%;
    margin-top: 30px;
    /*background: deepskyblue;*/
    text-align: center;
    /*display: flex;*/
}

:deep(.el-pagination) {
    width: 100% !important;
    margin: 0 auto !important;
    justify-content: center;
}

/*分页器样式修改*/
:deep(.number) {
    font-weight: bold;
    width: 45px;
    height: 45px;
    background: #2e2e2e !important;
    color: #ffffff;
    border-radius: 50%;
    /*margin: 0 3px!important;*/
}

:deep(.number:hover) {
    color: #67d9e8;
}

:deep(.btn-prev) {
    font-weight: bold;
    width: 45px;
    height: 45px;
    background: #2e2e2e !important;
    color: #ffffff;
    border-radius: 50%;
}

:deep(.btn-next) {
    font-weight: bold;
    width: 45px;
    height: 45px;
    background: #2e2e2e !important;
    color: #ffffff;
    border-radius: 50%;
}

:deep(.more) {
    font-weight: bold;
    width: 45px;
    height: 45px;
    background: #2e2e2e !important;
    color: #ffffff;
    border-radius: 50%;
}

:deep(.is-active) {
    background: #67d9e8 !important;
}

/*移动端缩小*/
@media (max-width: 650px) {
    :deep(.number) {
        width: 35px;
        height: 35px;
    }

    :deep(.btn-prev) {
        width: 35px;
        height: 35px;
    }

    :deep(.btn-next) {
        width: 35px;
        height: 35px;
    }

    :deep(.more) {
        width: 35px;
        height: 35px;
    }
}

</style>