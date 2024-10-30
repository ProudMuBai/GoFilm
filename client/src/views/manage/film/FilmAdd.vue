<template>
  <div class="container">
    <h2 style="text-align: start">添加影片</h2>
    <el-form :model="data.form"   class="film_add_form">
      <el-form-item>
        <div class="el-input-group__prepend" style="border: 1px solid #dcdfe6;border-right: none;border-radius: 3px;height: 32px">影片分类: </div>
          <el-select v-model="data.currentClass" style="width: calc(100% - 103px)" @change="changeClass"	  placeholder="影片分类选择">
            <el-option v-for="item in data.options.category" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
      </el-form-item>
      <el-form-item>
        <el-input  v-model="data.form.name" placeholder="请输入影片名称" clearable >
          <template #prepend>影片名称: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.subTitle"  placeholder="影片别名, 可留空" clearable >
          <template #prepend>影片别名: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.initial"  placeholder="影片检索首字母, 大写" clearable >
          <template #prepend>首字母: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.classTag"  placeholder="影片剧情标签(多标签以逗号分隔): 奇幻,校园,爱情" clearable >
          <template #prepend>剧情Tag: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.director"  placeholder="导演名, 多个名称以逗号进行分隔" clearable >
          <template #prepend>导演: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input  v-model="data.form.actor" placeholder="主演名, 多个名称以逗号进行分隔" clearable >
          <template #prepend>主演: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.writer"  placeholder="作者名, 多个名称以逗号进行分隔" clearable >
          <template #prepend>作者: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.remarks"  placeholder="影片更新进度信息, 完结, HD, 更新至xx集" clearable >
          <template #prepend>更新状态: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.releaseDate"  placeholder="影片上映时间: YYYY-MM-DD" clearable >
          <template #prepend>上映时间: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.area"  placeholder="影片来源地区信息" clearable >
          <template #prepend>地区: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.lang"  placeholder="影片语言信息" clearable >
          <template #prepend>语言: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.year"  placeholder="影片上映年份信息: YYYY" clearable >
          <template #prepend>年份: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.state"  placeholder=" 影片状态: 正片 | 预告片" clearable >
          <template #prepend>影片状态: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.dbId"  placeholder="豆瓣ID" clearable >
          <template #prepend>豆瓣Id: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.dbScore"  placeholder="豆瓣评分" clearable >
          <template #prepend>豆瓣评分: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.hits"  placeholder="影片热度(播放数)" clearable >
          <template #prepend>影片热度: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.picture"  placeholder="输入图片URL链接或点击上传到服务器并自动生成URL连接信息)" clearable >
          <template #prepend>影片海报: </template>
          <template #append>
            <!--<el-button> 上传图片</el-button>-->
            <el-upload class="upload-demo" :show-file-list="false" action="#" :http-request="customUpload">
              <el-button type="primary">上传图片</el-button>
            </el-upload>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-input v-model="data.form.playForm"  placeholder="影片播放资源来源: xxXm3u8" clearable >
          <template #prepend>播放来源: </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <template #label>
          <span class="el-input-group__prepend cus_label" >剧情简介: </span>
        </template>
          <el-input v-model="data.form.content" :autosize="{ minRows: 2, maxRows: 5 }" type="textarea" placeholder="影片剧情描述信息" />
      </el-form-item>
      <el-form-item label="播放地址:">
        <template #label>
          <span class="el-input-group__prepend cus_label" >播放地址: </span>
        </template>
        <el-input v-model="data.form.playLink" :autosize="{ minRows: 2, maxRows: 5 }" type="textarea"
                  placeholder="影片播放地址信息: &#10;格式: 第01集$https://xxx/xxx/index.m3u8#第02集$https://xxx/xxx/index.m3u8" />
      </el-form-item>
      <el-form-item class="form_btn">
        <el-button type="primary" @click="addFilm" >添加影片</el-button>
        <el-button >清空信息</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {PictureFilled} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage} from "element-plus";

// 表单数据初始化
const formInit = {
  id: 0,
  cid: 0,
  pid: 0,
  name: '',
  picture: '',
  subTitle: '',
  cName: '',
  enName: '',
  initial: '',
  classTag: '',
  actor: '',
  director: '',
  writer: '',
  blurb: '',
  content: '',
  remarks: '',
  releaseDate: '',
  area: '',
  lang: '',
  year: '',
  state: '',
  updateTime: '',
  addTime: '',
  dbId: 0,
  dbScore: '',
  hits: 0,
  playForm: '',
  playLink: '',
}
const data = reactive({
  form: formInit,
  options: {
    category: [{id:0,name:'分类名称',pid: 0}],
  },
})

const customUpload = (options:any)=>{
  let file = options.file
  let formData = new FormData();
  formData.append("file", file)
  ApiPost(`/manage/file/upload`, formData).then((resp:any)=>{
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      data.form.picture = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 选择影片分类
const changeClass = (value:any)=>{
  data.options.category.forEach(item=>{
    if (item.id == value) {
      data.form.cid = item.id
      data.form.pid = item.pid
      data.form.cName =item.name
    }
  })
}

// 添加影片
const addFilm = ()=>{
  // 对数字类型参数进行转换
  let params = data.form
  params.dbId = params.dbId - 0
  params.hits = params.hits - 0
  ApiPost(`/manage/film/add`,{...data.form}).then((resp: any) => {
    if (resp.code === 0) {
      // 更新成功后关闭弹窗, 重新获取最新的分类表格信息
      ElMessage.success({message: resp.msg})
      // 重置表单数据
      data.form = formInit
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  ApiGet(`/manage/film/class/tree` ).then((resp: any) => {
    if (resp.code === 0) {
      let l = [{id:0,name:'分类名称',pid: 0}]
      l.pop()
      resp.data.children.forEach((item:any)=>{
        if (item.children && item.children.length > 0) {
          l = [...l,...item.children]
        }
      })
      data.options.category = l
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
})
</script>

<style scoped>
.container {
  background: var(--bg-light);
}
.film_add_form{
  width: 100%;
  flex-flow: wrap;
  display: flex;
  justify-content: start;
}
:deep(.el-form-item) {
  --el-fill-color-light: var(--bg-fill-light);
  width: calc(50% - 120px);
  margin: 15px 60px;
}
.form_btn{
  width: 100%!important;
  margin: 40px auto;
}

:deep(.form_btn .el-form-item__content){
  justify-content: center;
}
:deep(.el-form-item__label){
  padding-right: 0!important;
}
.cus_label{
  border: 1px solid #dcdfe6;
  border-right: none;
  border-radius: 3px;
}

</style>