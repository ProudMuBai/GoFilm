<template>
    <div class="container">

        <div class="content_item" v-for="item in data.info.content">
            <template v-if="item.nav.name !='综艺' & item.nav.name !='综艺片'">
                <el-row class="row-bg  cus_nav" justify="space-between">
                    <el-col :span="12" class="title">
                        <span :class="`iconfont ${item.nav.name.search('电影') != -1?'icon-film':item.nav.name.search('剧') != -1?'icon-tv':'icon-cartoon'}`"
                              style="color: #79bbff;font-size: 32px;margin-right: 10px; line-height: 130%"/>
                        <a :href="`/categoryFilm?pid=${item.nav.id}`">{{ item.nav.name }}</a>
                    </el-col>
                    <el-col :span="12">
                        <ul class="nav_ul">
                            <li v-for="c in item.nav.children" class="nav_category hidden-md-and-down"><a
                                    :href="`/categoryFilm?pid=${c.pid}&cid=${c.id}`">{{ c.name }}</a></li>
                            <li class="nav_category hidden-md-and-down"><a :href="`/categoryFilm?pid=${item.nav.id}`">更多 ></a></li>
                        </ul>
                    </el-col>
                </el-row>
                <el-row class="cus_content">
                    <el-col :md="24" :lg="20" :xl="20" class="cus_content">
                        <!--影片列表-->
                        <FilmList :list="item.movies.slice(0,12)"/>
                    </el-col>
                    <el-col :md="0" :lg="4" :xl="4" class="hidden-md-and-down content_right">
                        <h3 class="hot_title">🔥热播{{item.nav.name}}</h3>
                        <template v-for="(m,i) in item.hot.slice(0,12)">
                            <div class="content_right_item">
                                <a :href="`/filmDetail?link=${m.mid}`"><b class="top_item">{{ i + 1 + '.' }}</b>
                                    <span>{{ m.name }}</span></a>
                            </div>
                        </template>
                    </el-col>
                </el-row>
            </template>
        </div>
    </div>
</template>

<script lang="ts" setup>
// 顶部轮播图
import 'element-plus/theme-chalk/display.css'
import {onBeforeMount, reactive} from "vue";
import {ApiGet} from "../../utils/request";
import FilmList from "../../components/FilmList.vue";

const data = reactive({
    info: {}
})
onBeforeMount(() => {
    ApiGet('/index').then((resp: any) => {
        data.info = resp.data
    })
})

</script>

<style scoped>

.container {
    margin: 0 auto;
}

.content_item {
    padding: 10px;
    margin-bottom: 25px;
}

.title {
    display: flex;
    text-align: left;
    height: 40px;
}

.title > a {
    min-width: 40px;
    color: rgb(221, 221, 221)
}

a {
    color: #333;
    padding-top: 10px;
    text-decoration: none;
    outline: none;
    -webkit-tap-highlight-color: transparent;
}

.cus_nav {
    border-bottom: 1px solid rgb(46, 46, 46);

    height: 40px;
}

.nav_ul {
    list-style-type: none;
    display: flex;
    flex-direction: row;
    justify-content: end;
    margin: 0;
}

.nav_category > a {
    color: #c9c4c4;

}

.nav_category > a:hover {
    color: #1890ff;
}

.nav_ul > li {
    /*min-width: 60px;*/
    white-space: nowrap;
    line-height: 40px;
    margin: 0 8px;
    text-align: center;
    color: #999;
    font-size: 14px;
    font-weight: 400;
}



/*影片简介区域*/
.cus_content {
    display: flex;
    padding-top: 15px;
}

.content_right {
    width: 100%;
}

.content_right_item {
    display: flex;
    padding-left: 10px;
    border-bottom: 1px solid rgb(46, 46, 46);
}

.content_right_item > a {
    padding: 10px 15px 10px 0;
    color: hsla(0, 0%, 100%, .87);
    display: block;
    flex-grow: 1;
    text-align: left;
    overflow: hidden;
    text-overflow: ellipsis;
    -o-text-overflow: ellipsis;
    white-space: nowrap;
}
.hot_title {
    text-align: left;
    margin: 8px 0;
}

:deep(.top_item) {
    color: red;
    /*font-style: oblique 10deg;*/
    font-style: italic;
    /*font-family: Inter;*/
    margin-right: 6px;
}

.content_right_item a span:hover {
    color: orange;
}


</style>

<!--移动端修改-->
<style scoped>
@media (min-width: 650px) {
    .cus_content_item {
        padding: 10px;
        overflow: hidden;
        /*margin-bottom: 10px;*/
    }
}

@media (max-width: 650px) {
    .cus_content_item {
        padding: 0 6px 0 0;
        margin-bottom: 10px;
        overflow: hidden;
    }

    .nav_ul {
        justify-content: end;
    }
}
</style>