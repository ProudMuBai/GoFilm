<template>
  <div class="container">

    <div class="content">
      <el-table
          :data="data.records" style="width: 100%" border size="default"
          table-layout="auto" max-height="calc(68vh - 20px)"
          row-key="id" fit
          :row-class-name="'cus-tr'">
        <el-table-column type="index"  align="left" min-width="35px" label="序列">
          <template #default="scope">
            <span style="color: #8b40ff">{{ scope.row.ID }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="originId" align="center" label="采集站">
          <template #default="scope">
            <el-tag type="primary" disable-transitions>{{ scope.row.originName }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="originId" align="center" min-width="100px" label="采集源ID">
          <template #default="scope">
            <el-tag type="success" disable-transitions>{{ scope.row.originId }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="collectType" align="center" label="采集类型" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" disable-transitions>{{ scope.row.collectType == 0 ? '影片详情' : '未知' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pageNumber" align="center" label="分页页码">
          <template #default="scope">
            <el-tag type="warning" disable-transitions>{{ scope.row.pageNumber }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hour" align="center" label="采集时长">
          <template #default="scope">
            <el-tag type="warning" disable-transitions>{{ scope.row.hour }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cause" align="center" label="失败原因" min-width="150px" >
          <template #default="scope">
            <el-tag type="danger" disable-transitions>{{ scope.row.cause }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" align="center" label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status == 1" type="warning">待重试</el-tag>
            <el-tag v-else type="success">已处理</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="UpdatedAt" align="center" label="执行时间" min-width="100px" >
          <template #default="scope">
            <el-tag :type="`${scope.row.status == 1 ? 'warning':'success' }`" disable-transitions>{{ scope.row.timeFormat }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" min-width="100px">
          <template #default="scope">
            <el-tooltip content="采集重试" placement="top"><el-button type="success" :icon="RefreshRight" @click="collectRecover(scope.row.ID)" plain circle/></el-tooltip>
            <el-tooltip content="删除记录" placement="top"><el-button type="danger" :icon="Delete" @click="delRecord(scope.row.ID)" plain circle/></el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination :page-sizes="[10, 20, 50, 100, 500]" background layout="prev, pager, next, sizes, total, jumper"
                       :total="data.page.total" v-model:page-size="data.page.pageSize"
                       v-model:current-page="data.page.current"
                       @change="getRecords" hide-on-single-page/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import {Aim, Delete, Edit, RefreshRight} from "@element-plus/icons-vue";
import {fmt} from "../../../utils/format";
import {onMounted, reactive} from "vue";
import {ApiGet} from "../../../utils/request";
import {ElMessage} from "element-plus";


const data = reactive({
  records: [],
  page: {current: 1, pageCount: 0, pageSize: 10, total: 0},
  params: {},
})


// 获取影片分页信息
const getRecords = () => {
  let {current, pageSize} = data.page
  let params = data.params
  ApiGet(`/manage/collect/record/list`, {...params, current, pageSize}).then((resp: any) => {
    if (resp.code === 0) {
      resp.data.list.map((item: any) => {
        // 对数据进行格式化处理
        item.timeFormat = fmt.dateFormat(new Date(item.UpdatedAt).getTime())
        return item
      })
      data.records = resp.data.list
      data.page = resp.data.params.paging
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 恢复采集, 对已采集失败的记录进行重试操作
const collectRecover = (id:number)=>{
  ApiGet(`/manage/collect/record/retry`, {id:id}).then((resp: any) => {
    if(resp.code === 0){
      ElMessage.success({message: resp.msg})
    } else{
       ElMessage.error({message: resp.msg})
    }
  })
}


// 删除当前记录, 貌似不合理
const delRecord = (id:number)=>{

}


// 清除所有失败采集记录
const clearAllRecord = ()=>{

}




onMounted(() => {
  // 获取分页数据
  getRecords()
  console.log(data)
})


</script>

<style scoped>







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