<template>
  <div class="header_container">
    <div class="left">
      <a href="javascript:;" @click="collapse.changeCollapse"
         :class="`iconfont ${ collapse.collapse.value ? 'icon-unfold': 'icon-fold'}`"></a>
      <h3>后台管理中心</h3>
    </div>
    <div class="right">
      <el-dropdown placement="bottom">
        <div class="dropdown_user">
          <el-avatar class="avatar" :size="35" :src="data.userInfo.avatar.toString()" alt="admin"/>
          <span>{{ data.userInfo.nickName }}</span>
          <el-icon class="el-icon--right">
            <arrow-down/>
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="a"><em class="iconfont icon-user-info"/>个人信息</el-dropdown-item>
            <el-dropdown-item command="a" @click="dialogV.changePwd = true"><em class="iconfont icon-change-pwd2"/>修改密码
            </el-dropdown-item>
            <el-dropdown-item command="e" divided @click="logout"><em class="iconfont icon-logout"/>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!--密码修改弹窗-->
    <el-dialog v-model="dialogV.changePwd" width="480px" title="用户密码修改">
      <el-form :model="form.changePwd" :rules="rules" label-width="80px">
        <el-form-item label="原始密码" prop="password">
          <el-input v-model="form.changePwd.password" :type="form.type.password?'text':'password'"/>
          <i :class="`cus-pwd iconfont ${form.type.password?'icon-eye2':'icon-eye'}`"
             @click="form.type.password = !form.type.password"/>
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="form.changePwd.newPassword" :type="form.type.newPassword?'text':'password'"/>
          <i :class="`cus-pwd iconfont ${form.type.newPassword?'icon-eye2':'icon-eye'}`"
             @click="form.type.newPassword = !form.type.newPassword"/>

        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.changePwd.confirmPassword" :type="form.type.confirmPassword?'text':'password'"/>
          <i :class="`cus-pwd iconfont ${form.type.confirmPassword?'icon-eye2':'icon-eye'}`"
             @click="form.type.confirmPassword = !form.type.confirmPassword"/>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="changePassword">确认</el-button>
        <el-button @click="cancelDialog">取消</el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>


<script setup lang="ts"  >
import {ArrowDown} from "@element-plus/icons-vue";
import {inject, onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../utils/request";
import {ElMessage} from "element-plus";
import {useRouter} from "vue-router";
import {clearAuthToken} from "../../utils/token";

const router = useRouter()

const data = reactive({
  userInfo: {id: Number, userName: String, email: String, gender: Number, nickName: String, avatar: String, status: Number}
})

// 侧边菜单栏展开状态
const collapse = inject('collapse')



// 用户密码修改对话框
const dialogV = reactive({
  changePwd: false,
})

// 用户密码修改表单
const form = reactive({
  changePwd: {password: '', newPassword: '', confirmPassword: ''},
  type: {password: false, newPassword: false, confirmPassword: false},
})

// 校验密码一致性
const regex = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[$@$!%*?&])[A-Za-z\\d$@$!%*?&]{8,12}$`
const validateNewPwd = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('新密码不能为空'))
  } else if (!value.match(regex)) {
    // ruleFormRef.value.validateField('checkPass', () => null)
    callback(new Error('密码必须为8-12位且包含大小写字母数字和特殊字符'))
  }
  callback()
}
const validateConfirmPwd = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('确认密码不能为空'))
  } else if (form.changePwd.newPassword !== '' && form.changePwd.newPassword != form.changePwd.confirmPassword) {
    callback(new Error('新密码与确认密码不一致'))
  }
  callback()
}
// 表单校验
const rules = reactive({
  password: [{required: true, message: '原始密码信息不能为空', trigger: 'blur'}],
  newPassword: [{required: true, validator: validateNewPwd, trigger: 'blur'}],
  confirmPassword: [{required: true, validator: validateConfirmPwd, trigger: 'blur'}],
})


// 修改密码
const changePassword = ()=>{
  ApiPost(`/changePassword`, {password: form.changePwd.password, newPassword: form.changePwd.newPassword}).then((resp: any) => {
    if (resp.code === 0) {
      // 退出登录成功则删除本地的token信息并返回到 登录 /login 界面
      // clearAuthToken()
      // router.push(`/login`)
      form.changePwd = {password: '', newPassword: '', confirmPassword: ''}
      dialogV.changePwd = false
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
// 对话框关闭
const cancelDialog = () => {
  dialogV.changePwd = false
  form.changePwd = {password: '', newPassword: '', confirmPassword: ''}
}

// 退出登录
const logout = () => {
  // 发送请求使当前的token信息失效
  ApiGet(`/logout`).then((resp: any) => {
    if (resp.code === 0 ) {
      // 退出登录成功则删除本地的token信息并返回到 登录 /login 界面
      clearAuthToken()
      router.push(`/login`)
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 获取用户信息
const getUserInfo = ()=>{
  ApiGet(`/manage/user/info` ).then((resp: any) => {
    if (resp.code === 0) {
      resp.data.avatar = resp.data.avatar == 'empty'?'https://s2.loli.net/2023/12/05/O2SEiUcMx5aWlv4.jpg': resp.data.avatar
      data.userInfo = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(()=>{
  // 获取用户信息, 初始化组件数据
  getUserInfo()
})

</script>

<style scoped>

/*密码输入框后缀*/
.cus-pwd {
  color: #b07ada;
  position: absolute;
  right: 8px;
  cursor: pointer;
}

.header_container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: space-between;
}

.left {
  display: flex;
  justify-content: center;
  align-items: center;
}

.left a {
  font-size: 30px;
  color: #9b49e7;
}

.left h3 {
  color: #d9ecff;
}


.left a:hover {
  color: #9b49e7b8;
}


.right {
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.el-dropdown) {
  outline: none;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 20px !important;
  --el-dropdown-menuItem-hover-color: #8b40ff;
}

.dropdown_user {
  outline: none;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}


.iconfont {
  margin-right: 10px;
}

.avatar {
  margin-right: 13px;
}
</style>