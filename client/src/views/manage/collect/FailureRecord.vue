<template>
  <div class="container">
    <div class="params_form">
      <el-form :model="data.params" class="cus_form">
        <el-form-item>
          <el-select v-model="data.params.originId" placeholder="采集来源">
            <el-option
              v-for="item in data.options.origin"
              :key="item.value"
              :label="item.name"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="false">
          <el-select v-model="data.params.collectType" placeholder="采集类型">
            <el-option
              v-for="item in data.options.collectType"
              :key="item.value"
              :label="item.name"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="data.params.status" placeholder="记录状态">
            <el-option
              v-for="item in data.options.status"
              :key="item.value"
              :label="item.name"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="data.dateGroup"
            value-format="YYYY-MM-DD HH:mm:ss"
            type="datetimerange"
            start-placeholder="起始时间"
            end-placeholder="终止时间"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="filterRecord">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="content">
      <el-table
        :data="data.records"
        style="width: 100%"
        border
        size="default"
        table-layout="auto"
        max-height="calc(68vh - 20px)"
        row-key="id"
        fit
        :row-class-name="'cus-tr'"
      >
        <el-table-column
          type="index"
          align="left"
          min-width="35px"
          label="序列"
        >
          <template #default="scope">
            <span style="color: #8b40ff">{{ scope.row.ID }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="originId" align="center" label="采集站">
          <template #default="scope">
            <el-tag type="primary" disable-transitions
              >{{ scope.row.originName }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="originId"
          align="center"
          min-width="100px"
          label="采集源ID"
        >
          <template #default="scope">
            <el-tag type="success" disable-transitions
              >{{ scope.row.originId }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="collectType"
          align="center"
          label="采集类型"
          show-overflow-tooltip
        >
          <template #default="scope">
            <el-tag type="success" disable-transitions
              >{{ scope.row.collectType == 0 ? "影片详情" : "未知" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pageNumber" align="center" label="分页页码">
          <template #default="scope">
            <el-tag type="warning" disable-transitions
              >{{ scope.row.pageNumber }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hour" align="center" label="采集时长">
          <template #default="scope">
            <el-tag type="warning" disable-transitions
              >{{ scope.row.hour }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="cause"
          align="center"
          label="失败原因"
          min-width="150px"
        >
          <template #default="scope">
            <el-tag type="danger" disable-transitions
              >{{ scope.row.cause }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" align="center" label="状态">
          <template #default="scope">
            <el-tag v-if="scope.row.status == 1" type="warning">待重试</el-tag>
            <el-tag v-else type="success">已处理</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="UpdatedAt"
          align="center"
          label="执行时间"
          min-width="100px"
        >
          <template #default="scope">
            <el-tag
              :type="`${scope.row.status == 1 ? 'warning' : 'success'}`"
              disable-transitions
              >{{ scope.row.timeFormat }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" min-width="100px">
          <template #default="scope">
            <el-tooltip content="采集重试" placement="top">
              <el-button
                type="success"
                :icon="RefreshRight"
                @click="collectRecover(scope.row.ID)"
                plain
                circle
              />
            </el-tooltip>
            <el-tooltip v-if="false" content="删除记录" placement="top">
              <el-button
                type="danger"
                :icon="Delete"
                @click="delRecord(scope.row.ID)"
                plain
                circle
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <div class="cus_util">
          <el-tooltip content="重试采集所有失败记录" placement="top">
            <el-button
              color="#d942bf"
              :icon="HelpFilled"
              @click="collectRecoverAll"
              >RetryAll</el-button
            >
          </el-tooltip>
          <el-tooltip content="清除已处理记录,保留未处理记录" placement="top">
            <el-button
              type="warning"
              :icon="WarningFilled"
              @click="cleanDoneRecord"
              >CleanDone</el-button
            >
          </el-tooltip>
          <el-tooltip content="清除所有记录" placement="top">
            <el-button type="danger" :icon="BellFilled" @click="cleanAllRecord"
              >CleanAll</el-button
            >
          </el-tooltip>
        </div>
        <el-pagination
          :page-sizes="[10, 20, 50, 100, 500]"
          background
          layout="prev, pager, next, sizes, total, jumper"
          :total="data.page.total"
          v-model:page-size="data.page.pageSize"
          v-model:current-page="data.page.current"
          @change="getRecords"
          hide-on-single-page
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  Delete,
  RefreshRight,
  BellFilled,
  WarningFilled,
  HelpFilled,
} from "@element-plus/icons-vue";
import { fmt } from "../../../utils/format";
import { onMounted, reactive } from "vue";
import { ApiGet } from "../../../utils/request";
import { ElMessage, ElMessageBox } from "element-plus";

const data = reactive({
  records: [],
  page: { current: 1, pageCount: 0, pageSize: 10, total: 0 },
  params: {
    originId: "",
    collectType: -1,
    status: -1,
    betweenTime: "",
    endTime: "",
  },
  dateGroup: [],
  options: { origin: [], collectType: [], status: [] },
});

// 获取影片分页信息
const getRecords = () => {
  let { current, pageSize } = data.page;
  let params = data.params;
  ApiGet(`/manage/collect/record/list`, { ...params, current, pageSize }).then(
    (resp: any) => {
      if (resp.code === 0) {
        resp.data.list.map((item: any) => {
          // 对数据进行格式化处理
          item.timeFormat = fmt.dateFormat(new Date(item.UpdatedAt).getTime());
          return item;
        });
        data.records = resp.data.list;
        data.page = resp.data.params.paging;

        // 初始化select参数
        data.options = resp.data.options;
      } else {
        ElMessage.error({ message: resp.msg });
      }
    }
  );
};

// 处理筛选参数, 获取满足条件的记录
const filterRecord = () => {
  if (data.dateGroup && data.dateGroup.length == 2) {
    data.params.beginTime = data.dateGroup[0];
    data.params.endTime = data.dateGroup[1];
  } else {
    data.params.beginTime = "";
    data.params.endTime = "";
  }
  getRecords();
};

// 恢复采集, 对已采集失败的记录进行重试操作
const collectRecover = (id: number) => {
  ApiGet(`/manage/collect/record/retry`, { id: id }).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({ message: resp.msg });
    } else {
      ElMessage.error({ message: resp.msg });
    }
  });
};

// 删除当前记录, 记录删除貌似不合理, 留以待定
const delRecord = (id: number) => {};

// collectRecoverAll 对目前记录的所有未处理记录进行重新采集
const collectRecoverAll = () => {
  ElMessageBox.confirm("是否对所有失效记录进行重新采集?", "采集失败记录处理", {
    confirmButtonText: "执行",
    cancelButtonText: "取消",
    type: "warning",
    center: true,
  })
    .then(() => {
      ApiGet(`/manage/collect/record/retry/all`).then((resp: any) => {
        if (resp.code === 0) {
          ElMessage.success({ message: resp.msg });
        } else {
          ElMessage.error({ message: resp.msg });
        }
      });
    })
    .catch(() => {
      ElMessage({ type: "warning", message: "采集恢复操作已取消!!!" });
    });
};

// cleanDoneRecord 清除已处理的记录
const cleanDoneRecord = () => {
  ElMessageBox.confirm("是否清除所有已处理的记录?", "记录清除", {
    confirmButtonText: "执行",
    cancelButtonText: "取消",
    type: "warning",
    center: true,
  })
    .then(() => {
      ApiGet(`/manage/collect/record/clear/done`).then((resp: any) => {
        if (resp.code === 0) {
          ElMessage.success({ message: resp.msg });
        } else {
          ElMessage.error({ message: resp.msg });
        }
      });
    })
    .catch(() => {
      ElMessage({ type: "warning", message: "记录清除已取消!!!" });
    });
};

// 清除所有失败采集记录
const cleanAllRecord = () => {
  ElMessageBox.confirm("是否清除所有记录?", "记录清除", {
    confirmButtonText: "执行",
    cancelButtonText: "取消",
    type: "warning",
    center: true,
  }).then(() => {
      ApiGet(`/manage/collect/record/clear/all`).then((resp: any) => {
        if (resp.code === 0) {
          ElMessage.success({ message: resp.msg });
        } else {
          ElMessage.error({ message: resp.msg });
        }
      })
    }).catch(() => {
      ElMessage({ type: "warning", message: "记录清除已取消!!!" })
    })
}

onMounted(() => {
  // 获取分页数据
  getRecords();
  console.log(data);
});
</script>

<style scoped>
/*头部检索表单*/
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
  display: flex;
  justify-content: space-between;
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

.cus_util {
  border: none;
}
</style>