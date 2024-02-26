<template>
  <div class="container">
    <el-container>
      <el-aside :width="collapse?'auto':'190px'" class="side">
        <Sidebar/>
      </el-aside>
      <el-container>
        <el-header class="header">
          <ManageHeader/>

        </el-header>
        <el-main class="view">
          <router-view/>
        </el-main>
        <!--<el-footer>Footer</el-footer>-->
      </el-container>
    </el-container>
  </div>
</template>


<script setup lang="ts">


import ManageHeader from "../../components/Manage/ManageHeader.vue";
import {provide, reactive, ref} from "vue";
import Sidebar from "../../components/Manage/Sidebar.vue";

// 侧边栏菜单状态
const collapse = ref(false)

// 改变侧边栏菜单的 展开|收起 状态
const changeCollapse = () => {
  collapse.value = !collapse.value
}
// 传递给子组件使用
provide('collapse', {
  collapse,
  changeCollapse,
})


</script>


<style scoped>


.header {
  background-image: linear-gradient(135deg, #81bfff 10%, #e267f0 100%);
}

.side {
  color: white;
}

.view {
  background: rgb(245 245 245);
  min-height: calc(100vh - 200px);
}
.el-container{
  max-height: 100vh;
}
:deep(.el-container){

}


/*管理后台公共样式*/
:deep(.el-switch) {
  --el-switch-on-color: #9b49e7b8;
  --el-switch-off-color: #cf2dae80;
}

:deep(.el-table) {
  --el-table-header-bg-color: #9b49e726;
  --el-table-border-color: #9b49e733;
}



:deep(.el-radio-group) {
  --el-color-primary: #9b49e7;
  --el-text-color-placeholder: #73859f;
}

.cus_util {
  display: flex;
  padding: 10px 8px;
  border-left: 2px solid #9b49e733;
  border-right: 2px solid #9b49e733;
  border-bottom: 2px solid #9b49e733;
  background: #ffffff;
  justify-content: end;
}

:deep(.el-input-number) {
  --el-fill-color-light: #e163ff8f;
  border-radius: var(--el-border-radius-base);
  padding: 0 0 !important;
}

:deep(.el-input-number__decrease) {
  left: 0;
  top: 0;
  bottom: 0;
}

:deep(.el-input-number__increase) {
  right: 0;
  top: 0;
  bottom: 0;
}

:deep(.el-tag--info) {
  --el-fill-color: #67d9e863;
}

:deep(.el-form-item__error) {
  --el-color-danger: #b3249a;
}

</style>