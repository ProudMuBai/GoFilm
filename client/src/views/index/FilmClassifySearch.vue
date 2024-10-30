<template>
  <div class="container" v-if="d.title.name">
    <div class="title">
      <a :href="`/filmClassify?Pid=${d.title.id}`" >{{ d.title.name }}</a>
      <span class="line"/>
      <a :href="`/filmClassifySearch?Pid=${d.title.id}`" class="h_active">{{ `${d.title.name}库` }}</a>
    </div>
    <!--影片分类检索-->
    <div class="t_container">
      <div class="t_item" v-for="k in d.search.sortList ">
        <div class="t_title">{{d.search.titles[k]}} <b class="iconfont icon-triangle"/> </div>
        <div class="tag_group">
          <a href="javascript:void(false)" :class="`tag ${t['Value'] === d.searchParams[k]?'t_active':''}`" v-for="t in d.search.tags[k]" @click="handleTag(k,t['Value'])" >
            {{t['Name']}}
          </a>
        </div>
      </div>
    </div>

    <!--影片列表展示-->
    <FilmList :col="7" :list="d.list"/>
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
import {onMounted, reactive, toRefs} from "vue";
import {useRouter} from "vue-router";
import {ApiGet} from "../../utils/request";
import {ElMessage} from "element-plus";
import {ArrowRightBold, ArrowLeftBold} from '@element-plus/icons-vue'
import FilmList from "../../components/index/FilmList.vue";

// 页面所需数据
const d = reactive({
  title: {},
  list: [],
  search: {
    sortList:[],
    titles: [],
    tags: [],
  },
  page: {
    current: 0,
  },
  searchParams:{
    Pid: '',
    Category: '',
    Plot: '',
    Area: '',
    Language: '',
    Year: '',
    Sort: '',
  },

})

// 获取路由参数查询对应数据
const router = useRouter()

// 点击分页按钮事件 current-change
const changeCurrent = (currentVal: number) => {
  handleParams()
}

// 分类tag点击事件
const handleTag = (k:string,v:string)=>{
  // 设置被点击的tag属性值
  d.searchParams[k as keyof typeof d.searchParams] = v
  // searchTag改变, 重置 current当前页码
  d.page.current = 1
  handleParams()
}
const handleParams = ()=> {
  let q = ''
  for (let k in d.searchParams) {
    let val = d.searchParams[k as keyof typeof d.searchParams]
    if (val != '') {
      q += `&${k}=${val}`
    }
  }
  location.href = '/filmClassifySearch?'+q.slice(1)+`&current=${d.page.current}`
}

// 请求数据
const getFilmData = () => {
  let query = router.currentRoute.value.query
  ApiGet(`/filmClassifySearch`, {...query}).then((resp: any) => {
    if (resp.code === 0) {
      d.title = resp.data.title
      d.list = resp.data.list
      d.page = resp.data.page
      d.search = resp.data.search
      d.searchParams = resp.data.params
    } else {
      ElMessage.error({message: "影片搜索结果异常,请稍后刷新重试", duration: 1000})
    }
  })
}

onMounted(() => {
  getFilmData()
})



</script>


<style scoped>
@import "/src/assets/css/classify.css";
@import "/src/assets/css/pagination.css";

@media (min-width: 768px) {
  .tag {
    margin: 0 8px;
    padding: 6px 12px;
  }
  .t_title {
    padding: 3px 0;
  }
}
@media (max-width: 768px) {
  .tag {
    margin: 0 5px;
    padding: 4px 10px;
    font-size: 12px;
  }
}

.t_container {
  display: block;
  font-size: 14px;
  padding-bottom: 10px;
  margin-bottom: 30px;
  border-bottom: 1px solid rgba(255,255,255, 0.12);
}

.t_item {
  display: flex;
  justify-content: start;
  margin: 14px 0;
  white-space: nowrap;
}

.t_title {
  display: inline-block;
  font-size: 17px;
  font-weight: 700;
  text-align: left;
  color: rgba(255,255,255,0.35);
  border-radius: 6px;
  margin-right: 12px;
}
.t_title b{
  color: rgba(255,255,255,0.15);
}

.tag_group {
  display: flex;
  justify-content: start;
  flex-flow: nowrap;
  overflow: auto;
}

.tag {
  display: inline-block;
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 5px;
  text-align: center;
}

.t_active {
  background: rgba(255,255,255,0.12);
  color: #ffa500cc!important;
  border: none!important;
}


</style>
<!--移动端修改-->
<style scoped>
@media (max-width: 768px) {

  /*顶部内容区域*/
  .header {
    width: 100%;
    margin-bottom: 100px;
    background: none !important;
  }


}
</style>

<style scoped>
.container {
  max-width: 100vw;
}

@media (min-width: 768px) {

  /*顶部内容区域*/
  .header {
    width: 100%;
  }




}
</style>

