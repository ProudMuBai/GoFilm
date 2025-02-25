<template>
  <div class="container">
    <div class="params_form">
      <el-form :model="data.params" class="cus_form">
        <el-form-item>
          <el-input v-model="data.params.name" style="display: inline-block;text-align: left" placeholder="片名搜素"
                    :suffix-icon="Search"/>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.classId" @change="changeClass" placeholder="影片分类">
            <el-option v-for="item in data.options.class" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.plot" placeholder="剧情筛选">
            <el-option v-for="item in data.options.Plot" :key="item.Value" :label="item.Name" :value="item.Value"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.area" placeholder="地区筛选">
            <el-option v-for="item in data.options.Area" :key="item.Value" :label="item.Name" :value="item.Value"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.language" placeholder="语言筛选">
            <el-option v-for="item in data.options.Language" :key="item.Value" :label="item.Name" :value="item.Value"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.year" placeholder="上映年份">
            <el-option v-for="item in data.options.year" :key="item.Value" :label="item.Name" :value="item.Value"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.remarks" placeholder="更新状态">
            <el-option v-for="item in data.options.remarks" :key="item.Value" :label="item.Name" :value="item.Value"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker v-model="data.dateGroup" value-format="YYYY-MM-DD HH:mm:ss" type="datetimerange" start-placeholder="起始时间"
                          end-placeholder="终止时间"  />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchFilm" >查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="content">
      <el-table
          :data="data.list" style="width: 100%" border size="default"
          table-layout="auto" max-height="calc(68vh - 20px)"
          row-key="id"
          :row-class-name="'cus-tr'">
        <el-table-column type="index" min-width="40px" align="left" label="序号">
          <template #default="scope">
            <span style="color: #8b40ff">{{ serialNum(scope.$index) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="mid" align="center" label="影片ID">
          <template #default="scope">
            <el-tag type="success" disable-transitions>{{ scope.row.mid }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" align="left" label="影片名称" show-overflow-tooltip class-name="col_name"/>
        <!--<el-table-column prop="subTitle" align="center" label="影片别名" />-->
        <el-table-column prop="cName" align="center" label="所属分类">
          <template #default="scope">
            <el-tag type="warning" disable-transitions>{{ scope.row.cName }}</el-tag>
          </template>
        </el-table-column>
        <!--<el-table-column prop="classTag" align="left" label="剧情标签" >-->
        <!--  <template #default="scope">-->
        <!--    <el-tag v-for="t in scope.row.classTag" style="margin: 2px 3px 2px 0" type="warning" disable-transitions>{{ t }}</el-tag>-->
        <!--  </template>-->
        <!--</el-table-column>-->
        <el-table-column prop="year" align="center" label="年份">
          <template #default="scope">
            <el-tag type="warning" disable-transitions>{{ scope.row.year }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column sortable  prop="score" align="center" label="评分">
          <template #default="scope">
            <el-tag type="success" disable-transitions>{{ scope.row.score }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column sortable  prop="hits" align="center" label="热度">
          <template #default="scope">
            <el-tag type="danger" disable-transitions>🔥{{ scope.row.hits }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remarks" align="center" label="更新状态">
          <template #default="scope">
            <el-tag v-if="(scope.row.remarks+'').indexOf('更新') != -1"  type="warning" >{{scope.row.remarks }}</el-tag>
            <el-tag v-else type="success" >{{scope.row.remarks }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column sortable prop="updateStamp" align="center" label="更新时间">
          <template #default="scope">
            <el-tag type="success" disable-transitions>{{fmt.dateFormat(scope.row.updateStamp*1000)}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" min-width="100px">
          <template #default="scope">
            <el-button type="success" :icon="Aim" @click="" plain circle/>
            <el-button type="success" :icon="RefreshRight" @click="UpdateFilm(scope.row.mid)" plain circle/>
            <el-button type="primary" :icon="Edit" @click="" plain circle/>
            <el-button type="danger" :icon="Delete" @click="delFilm(scope.row.ID)" plain circle/>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination :page-sizes="[10, 20, 50, 100, 500]" background layout="prev, pager, next, sizes, total, jumper"
                       :total="data.page.total" v-model:page-size="data.page.pageSize"
                       v-model:current-page="data.page.current"
                       @change="getFilmPage" hide-on-single-page/>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import {Delete, Edit, Aim, Search, Calendar, RefreshRight} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet} from "../../../utils/request";
import {ElMessage} from "element-plus";
import {fmt} from '../../../utils/format'

const data = reactive({
  list: [],
  page: {current: 1, pageCount: 0, pageSize: 10, total: 0},
  params: {name:'', pid: 0, cid: 0, plot: '', area: '', language: '', year: '',remarks: '',beginTime:'', endTime: '',},
  options: {class: [{id: 0, pid: -1, name: '', show: true}], Plot: [], Area: [], Language: [], year: [], remarks: []},
  dateGroup: [],
  classId: 0,
})
let tags = {}

// 选择影片分类时触发连锁事件
const changeClass = (value: any) => {
  for (let i = 0; i < data.options.class.length; i++) {
    if (data.options.class[i].id == value) {
      // 从分类列表中获取匹配的分类信息
      let c = data.options.class[i]
      // 设置请求参数中的分类id和父级分类id信息

      // 如果选择的是一级分类 则直接设置cid作为分类参数
      if (c.pid <= 0) {
        data.params.pid = c.id
        data.params.cid = 0
        return
      } else if (c.pid == data.params.pid){
        // 如果一级分类没有改变则只改变分类ID,并退出
        data.params.pid = c.pid
        data.params.cid = c.id
        return
      }
      // 如果一级分类改变则改变分类ID,并改变关联的tag参数信息
      data.params.pid = c.pid
      data.params.cid = c.id
      // 从tags列表中获取当前分类下的可用tag信息 (一级分类使用id获取, 二级分类使用pid获取)
      let t = c.pid == 0 ? tags[c.id as keyof typeof tags]:tags[c.pid as keyof typeof tags]
      // 匹配成功则设置对应的options参数
      if (t) {
        data.options.Plot = t['Plot']
        data.options.Area = t['Area']
        data.options.Language =  t['Language']
      } else {
        // 匹配失败则清空已有的options参数
        data.options.Plot = []
        data.options.Area = []
        data.options.Language =  []
      }
      // tags改变时清空对应的param参数
      data.params.plot = ''
      data.params.area = ''
      data.params.language = ''
    }
  }
}

// 生成序列号
const serialNum = (index: number) => {
  return (data.page.current - 1) * data.page.pageSize + index + 1
}

// 搜索满足条件的影片
const searchFilm = ()=>{
  let p = data.params
  // 时间选择器参数处理 如果 dateGroup 不为空 则追加时间范围参数
  if (data.dateGroup && data.dateGroup.length == 2) {
    p.beginTime = data.dateGroup[0]
    p.endTime = data.dateGroup[1]
  } else {
    p.beginTime = ''
    p.endTime = ''
  }
  getFilmPage()

}

// 更新影片信息
const UpdateFilm = (id:number)=>{
  let ids = id + ''
  ApiGet(`/manage/spider/update/single`, {ids: id}).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 获取影片分页信息
const getFilmPage = () => {
  let {current, pageSize} = data.page
  let params = data.params
  ApiGet(`/manage/film/search/list`, {...params,current, pageSize}).then((resp: any) => {
    if (resp.code === 0) {
      data.list = resp.data.list ? resp.data.list.map((item: any) => {
        // 对数据进行格式化处理
        item.year = item.year <= 0 ? '未知' : item.year
        item.score = item.score == 0 ? '暂无' : item.score
        // if (item.classTag) {
        //   item.classTag = [...item.classTag.toString().split(',')]
        // } else {
        //   item.classTag = ['未知']
        // }
        return item
      }) : []
      data.page = resp.data.params.paging
      data.options.class = resp.data.options.class
      data.options.remarks = resp.data.options.remarks
      data.options.year = resp.data.options.year
      tags = resp.data.options.tags
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(() => {
  getFilmPage()
})

// 删除影片信息
const delFilm = (id:number) =>{
  ApiGet( `/manage/film/search/del`, {id: id}).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getFilmPage()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

</script>

<style scoped>
.params_form {
  background: var(--bg-light);
  margin-bottom: 20px;
  padding: 10px 20px;
}

.cus_form {
  width: 100%;
  flex-flow: wrap;
  display: flex;
  justify-content: start;
}

:deep(.el-form-item) {
  width: calc(16% - 12px);
  margin: 10px 6px;
}


:deep(.el-table) {
  color: var(--content-text-color);
}

.content {
  border: 1px solid #9b49e733;
  background: var(--bg-light);
  --el-color-primary: var(--paging-parmary-color);
}

.pagination {
  margin: 20px auto;
  max-width: 100%;
  text-align: center;
  padding-right: 50px;
}

:deep(.el-pagination) {
  width: 100% !important;
  justify-content: end;
  --el-color-primary: var(--paging-parmary-color);
}

:deep(.el-pager li) {
  --el-pagination-button-bg-color: var(--btn-bg-linght);
  border: 1px solid var(--border-gray-color);
}

:deep(.el-pagination button) {
  --el-disabled-bg-color: var(--btn-bg-linght);
  --el-pagination-button-bg-color: var(--btn-bg-linght);
  border: 1px solid var(--border-gray-color);
}
</style>
