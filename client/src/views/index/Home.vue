<template>
    <div class="container">
        <div class="content_item" v-for="item in data.info.content">
            <template v-if="item.nav.name !='综艺' & item.nav.name !='综艺片'">
                <el-row class="row-bg  cus_nav" justify="space-between">
                    <el-col :span="12" class="title">
                        <span :class="`iconfont ${item.nav.name.search('电影') != -1?'icon-film':item.nav.name.search('剧') != -1?'icon-tv':'icon-cartoon'}`"
                              style="color: #79bbff;font-size: 32px;margin-right: 10px; line-height: 130%" />
                        <a :href="`/categoryFilm?pid=${item.nav.id}`">{{ item.nav.name }}</a>
                    </el-col>
                    <el-col :span="12">
                        <ul class="nav_ul">
                            <li v-for="c in item.nav.children" class="nav_category hidden-md-and-down"><a
                                    :href="`/categoryFilm?pid=${c.pid}&cid=${c.id}`">{{ c.name }}</a></li>
                            <li class=" hidden-md-and-down">更多 ></li>
                        </ul>
                    </el-col>
                </el-row>
                <el-row class="cus_content">
                    <el-col :md="24" :lg="20" :xl="20" class="cus_content">
                        <el-row style="max-width: 100%">
                            <template v-for=" (m,i) in item.movies">
                                <el-col :md="4" :sm="6" :xs="8" v-if="i <12" class="cus_content_item">
                                    <a :href="`/filmDetail?link=${m.id}`" class="cus_content_link"
                                       @error="handleImgError"
                                       :style="{backgroundImage: `url('${m.picture}')`}">
                                        <span class="cus_tag hidden-md-and-down">{{ m.year }}</span>
                                        <span class="cus_tag hidden-md-and-down">{{ m.cName }}</span>
                                        <span class="cus_tag hidden-md-and-down">{{ m.area }}</span>
                                    </a>
                                    <a :href="`/filmDetail?link=${m.id}`"
                                       class="content_text content_text_tag">{{ m.name }}</a>
                                    <span class="cus_remark">{{ m.remarks }}</span>
                                </el-col>
                            </template>
                        </el-row>
                    </el-col>
                    <el-col :md="0" :lg="4" :xl="4" class="hidden-md-and-down content_right">
                        <template v-for="(m,i) in item.movies">
                            <div class="content_right_item">
                                <a :href="`/filmDetail?link=${m.id}`"><b class="top_item">{{ i + 1 + '.' }}</b>
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
import 'element-plus/theme-chalk/display.css'
import {onBeforeMount, onMounted, reactive} from "vue";
import {ApiGet} from "../../utils/request";

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
    color: #999;
}

.nav_category > a:hover {
    color: #1890ff;
}

.nav_ul > li {
    min-width: 60px;
    line-height: 40px;
    text-align: center;
    color: #999;
    font-size: 14px;
    font-weight: 400;
}

/*svg图标*/
embed {
    width: 2rem;
    height: 2rem;
    margin-right: 8px;
    margin-top: 5px;
}

/*影片简介区域*/
.cus_content {
    display: flex;
    padding-top: 15px;
}

.cus_content_link {
    border-radius: 5px;
    display: flex;
    /*position: relative;*/
    padding-top: 125%;
    background-size: cover;
}

.cus_tag {
    text-align: center;
    color: rgb(255, 255, 255);
    padding: 0 3px;
    margin: 0 0 10px 8px;
    background: rgba(0, 0, 0, 0.55);
    font-size: 12px;
    border-radius: 5px;
}

.content_text_tag {
    font-size: 15px !important;
    color: rgb(221, 221, 221);
    padding: 2px 10px 2px 2px !important;
}

.cus_remark {
    display: block;
    width: 100%;
    padding-left: 3px;
    font-size: 12px;
    color: #999999;
    text-align: left;
}

.content_text {
    display: block;
    width: 100%;
    padding: 2px 10px 10px 2px;
    font-size: 12px;
    overflow: hidden;
    text-overflow: ellipsis;
    -o-text-overflow: ellipsis;
    white-space: nowrap;
    text-align: left;
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