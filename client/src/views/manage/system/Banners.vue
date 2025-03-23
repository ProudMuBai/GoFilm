<template>
  <h2 style="color: #8e48b4">首页横幅管理界面</h2>
  <el-table
      :data="data.banners" style="width: 100%" border size="default"
      :row-class-name="'cus-tr'" table-layout="auto">
    <el-table-column prop="name" label="影片名称"/>
    <el-table-column prop="collectType" align="center" label="影片类型">
      <template #default="scope">
        <el-tag type="warning">{{ scope.row.cName }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="上映年份">
      <template #default="scope">
        <el-tag type="warning">{{ scope.row.year }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="影片海报">
      <template #default="scope">
        <el-image style="width: 180px; height: 80px" :src="scope.row.poster" :preview-src-list="[scope.row.poster]"
                  preview-teleported fit="contain"/>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="影片封面">
      <template #default="scope">
        <el-image style="width: 60px; height: 80px" :src="scope.row.picture" :preview-src-list="[scope.row.picture]"
                  preview-teleported fit="cover"/>
      </template>
    </el-table-column>
    <el-table-column prop="collectType" align="center" label="排序">
      <template #default="scope">
        <el-tag disable-transitions>{{ scope.row.sort }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="resultModel" align="center" label="连载状态">
      <template #default="scope">
        <el-tag v-if="(scope.row.remarks+'').search('更新') == -1" type="success">{{ scope.row.remark }}</el-tag>
        <el-tag v-else type="primary">{{ scope.row.remark }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" align="center">
      <template #default="scope">
        <el-tooltip content="绑定影片信息" placement="top">
          <el-button type="success" :icon="Link" plain circle @click="openBindV(scope.row)"/>
        </el-tooltip>
        <el-tooltip content="修改海报信息" placement="top">
          <el-button type="primary" :icon="Edit" plain circle @click="openEditV(scope.row)"/>
        </el-tooltip>
        <el-tooltip content="删除海报信息" placement="top">
          <el-button type="danger" :icon="Delete" plain circle @click="delBanner(scope.row)"/>
        </el-tooltip>
      </template>
    </el-table-column>
  </el-table>
  <div class="cus_util">
    <el-button color="#9b49e7" :icon="CirclePlus" @click="openAddV">添加海报</el-button>
    <el-button type="danger" :icon="TakeawayBox" @click="clearCache">清除缓存</el-button>
  </div>


  <!--Banner添加弹窗-->
  <el-dialog v-model="data.dialogV.addV" width="680px" title="添加海报">
    <el-form :model="data.banner">
      <el-form-item label="影片ID&emsp;">
        <el-input v-model.number="data.banner.mid" placeholder="影片唯一ID"/>
      </el-form-item>
      <el-form-item label="影片名称">
        <el-input v-model="data.banner.name" placeholder="影片名称"/>
      </el-form-item>
      <el-form-item label="影片分类">
        <el-input v-model="data.banner.cName" placeholder="影片所属分类"/>
      </el-form-item>
      <el-form-item label="影片海报">
        <el-input v-model="data.banner.poster" placeholder="影片海报访问URL" class="upload_input" />
        <el-upload :show-file-list="false" action="#" :http-request="customUpload" :data="{type: 0}" class="upload" >
          <el-button color="#626aef" round plain :icon="UploadFilled" class="upload_btn">Upload</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="影片封面">
        <el-input v-model="data.banner.picture" placeholder="影片封面访问URL" class="upload_input" />
        <el-upload :show-file-list="false" action="#" :http-request="customUpload" :data="{type: 1}" class="upload"  >
          <el-button color="#626aef" round plain :icon="UploadFilled" class="upload_btn" >Upload</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="更新状态">
        <el-input v-model="data.banner.remark" placeholder="影片更新状态"/>
      </el-form-item>
      <el-form-item label="上映年份">
        <el-input-number v-model="data.banner.year" :min="0" :step="1" :max="2100" step-strictly/>
      </el-form-item>
      <el-form-item label="排序分值">
        <el-input-number v-model="data.banner.sort" :min="-100" :step="1" :max="100" step-strictly/>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button color="#cf48be" @click="data.dialogV.addBindV = true">绑定影片</el-button>
        <el-button color="#9b49e7" @click="add">确认添加</el-button>
        <el-button @click="data.dialogV.addV = false">取消</el-button>
      </span>
    </template>
    <!--影片绑定弹窗-->
    <el-dialog v-model="data.dialogV.addBindV" width="620px" title="绑定影片" align-center>
      <el-form :model="data.banner">
        <el-form-item label="搜索影片">
          <el-select-v2 v-model="data.FilmId" filterable :props="{label:'name', value: 'id'}" remote
                        :remote-method="loadingFilm" clearable
                        :options="data.options" :loading="data.loading" placeholder="请输入需要绑定的影片名称"
                        @change="changeFilm">
          </el-select-v2>
        </el-form-item>
        <el-form-item v-if="data.film.id">
          <div class="film_view">
            <a href="javascript:void(0);" :style="{backgroundImage: `url('${data.film.picture}')`}"></a>
            <div class="film_intro">
              <h3>{{ data.film.name }}</h3>
              <p class="tags">
                <span class="tag_c">{{ data.film.cName }}</span>
                <span>{{ data.film.year }}</span>
                <span>{{ data.film.area }}</span>
              </p>
              <p><em>导演:</em>{{ data.film.director }}</p>
              <p><em>主演:</em>{{ data.film.actor }}</p>
              <p class="blurb"><em>剧情:</em>{{ (data.film.blurb + '').replaceAll('　　', '') }}</p>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="bindAddBanner">确认绑定</el-button>
        <el-button @click="data.dialogV.addBindV = false">取消</el-button>
      </span>
      </template>
    </el-dialog>
  </el-dialog>

  <!--影片修改搜索弹窗-->
  <el-dialog v-model="data.dialogV.editV" width="680px" title="修改海报信息">
    <el-form :model="data.banner">
      <el-form-item label="影片ID&emsp;">
        <el-input v-model.number="data.banner.mid"  placeholder="影片唯一ID"/>
      </el-form-item>
      <el-form-item label="影片名称">
        <el-input v-model="data.banner.name" placeholder="影片名称"/>
      </el-form-item>
      <el-form-item label="影片分类">
        <el-input v-model="data.banner.cName" placeholder="影片所属分类"/>
      </el-form-item>
      <el-form-item label="影片海报">
        <el-input v-model="data.banner.poster" placeholder="影片海报访问URL" class="upload_input" />
        <el-upload :show-file-list="false" action="#" :http-request="customUpload" :data="{type: 0}" class="upload"  >
          <el-button color="#626aef" round plain :icon="UploadFilled" class="upload_btn" >Upload</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="影片封面">
        <el-input v-model="data.banner.picture" placeholder="影片封面访问URL" class="upload_input" />
        <el-upload :show-file-list="false" action="#" :http-request="customUpload" :data="{type: 1}" class="upload"  >
          <el-button color="#626aef" round plain :icon="UploadFilled" class="upload_btn" >Upload</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="更新状态">
        <el-input v-model="data.banner.remark" placeholder="影片更新状态"/>
      </el-form-item>
      <el-form-item label="上映年份">
        <el-input-number v-model="data.banner.year" :min="0" :step="1" :max="2100" step-strictly/>
      </el-form-item>
      <el-form-item label="排序分值">
        <el-input-number v-model="data.banner.sort" :min="-100" :step="1" :max="100" step-strictly/>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
         <el-button color="#cf48be" @click="data.dialogV.editBindV = true">绑定影片</el-button>
        <el-button color="#9b49e7" @click="edit">保存</el-button>
        <el-button @click="data.dialogV.editV = false">取消</el-button>
      </span>
    </template>
    <!--重新绑定影片-->
    <el-dialog v-model="data.dialogV.editBindV" width="620px" title="绑定影片" align-center>
      <el-form :model="data.banner">
        <el-form-item label="搜索影片">
          <el-select-v2 v-model="data.FilmId" filterable :props="{label:'name', value: 'id'}" remote
                        :remote-method="loadingFilm" clearable
                        :options="data.options" :loading="data.loading" placeholder="请输入需要绑定的影片名称"
                        @change="changeFilm">
          </el-select-v2>
        </el-form-item>
        <el-form-item v-if="data.film.id">
          <div class="film_view">
            <a href="javascript:void(0);" :style="{backgroundImage: `url('${data.film.picture}')`}"></a>
            <div class="film_intro">
              <h3>{{ data.film.name }}</h3>
              <p class="tags">
                <span class="tag_c">{{ data.film.cName }}</span>
                <span>{{ data.film.year }}</span>
                <span>{{ data.film.area }}</span>
              </p>
              <p><em>导演:</em>{{ data.film.director }}</p>
              <p><em>主演:</em>{{ data.film.actor }}</p>
              <p class="blurb"><em>剧情:</em>{{ (data.film.blurb + '').replaceAll('　　', '') }}</p>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="bindAddBanner">确认绑定</el-button>
        <el-button @click="data.dialogV.addBindV = false">取消</el-button>
      </span>
      </template>
    </el-dialog>
  </el-dialog>

  <!--搜索绑定的影片-->
  <el-dialog v-model="data.dialogV.bindV" width="680px" title="绑定影片">
    <el-form :model="data.banner">
      <el-form-item label="搜索影片">
        <el-select-v2 v-model="data.FilmId" filterable :props="{label:'name', value: 'id'}" remote
                      :remote-method="loadingFilm" clearable
                      :options="data.options" :loading="data.loading" placeholder="请输入需要绑定的影片名称"
                      @change="changeFilm">
        </el-select-v2>
      </el-form-item>
      <el-form-item v-if="data.film.id">
        <div class="film_view">
          <a :href="`/filmDetail?link=${data.film.id}`" :style="{backgroundImage: `url('${data.film.picture}')`}"></a>
          <div class="film_intro">
            <h3>{{ data.film.name }}</h3>
            <p class="tags">
              <span class="tag_c">{{ data.film.cName }}</span>
              <span>{{ data.film.year }}</span>
              <span>{{ data.film.area }}</span>
            </p>
            <p><em>导演:</em>{{ data.film.director }}</p>
            <p><em>主演:</em>{{ data.film.actor }}</p>
            <p class="blurb"><em>剧情:</em>{{ (data.film.blurb + '').replaceAll('　　', '') }}</p>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button color="#9b49e7" @click="bindFilm">确认绑定</el-button>
        <el-button @click="data.dialogV.bindV = false">取消</el-button>
      </span>
    </template>
  </el-dialog>

</template>

<script setup lang="ts">
import {
  CirclePlus,
  Delete,
  Edit,
  Link, TakeawayBox, UploadFilled,
} from "@element-plus/icons-vue";
import {onMounted, reactive} from "vue";
import {ApiGet, ApiPost} from "../../../utils/request";
import {ElMessage} from "element-plus";

// 渲染数据维护
const data = reactive({
  banners: [],
  banner: {id: '', mid: 0, name: '', cName: '', poster: '', picture: '', year: 0, remark: '', sort: 0},
  loading: false,
  FilmId: '',
  film: {},
  options: [{}],
  dialogV: {
    addV: false,
    editV: false,
    bindV: false,
    addBindV: false,
    editBindV: false,
  }
})

// banner添加功能组
const openAddV = () => {
  data.banner = {id: '', mid: 0, name: '', cName: '', poster: '', picture: '', year: 0, remark: '', sort: 0}
  data.dialogV.addV = true
}
const bindAddBanner = () => {
  // 同步绑定的影片信息到当前Banner
  data.banner.mid = data.film.id
  data.banner.name = data.film.name
  data.banner.cName = data.film.cName
  data.banner.picture = data.film.picture
  data.banner.year = parseInt(data.film.year)
  data.banner.remark = data.film.remarks
  data.dialogV.addBindV = false
  data.dialogV.editBindV = false
  ElMessage.success({message: "影片信息绑定成功!!!"})
}
const add = () => {
  ApiPost('/manage/banner/add', data.banner).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      data.banner = {id: '', mid: 0, name: '', cName: '', poster: '', picture: '', year: 0, remark: '', sort: 0}
      data.dialogV.addV = false
      getBanners()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 修改功能组
const openEditV = (b: any) => {
  data.banner = b
  data.dialogV.editV = true

}
const edit = () => {
  ApiPost('/manage/banner/update', data.banner).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      data.banner = {id: '', mid: 0, name: '', cName: '', poster: '', picture: '', year: 0, remark: '', sort: 0}
      data.dialogV.editV = false
      getBanners()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 绑定功能组
const openBindV = (b: any) => {
  data.banner = b
  data.dialogV.bindV = true
}
const loadingFilm = (query: string) => {
  if (query) {
    data.loading = true
    setTimeout(() => {
      data.loading = false
      ApiGet('/searchFilm', {keyword: query, current: 0}).then((resp: any) => {
        if (resp.code == 0) {
          data.options = resp.data.list
        } else {
          ElMessage.warning({message: resp.msg, duration: 1000})
          data.options = []
        }
      })
    }, 1500)
  }

}
const changeFilm = (val: any) => {
  data.options.forEach(item => {
    if (item.id == val) {
      data.film = item
    }
  })
}
const bindFilm = () => {
  // 同步绑定的影片信息到当前Banner
  data.banner.mid = data.film.id
  data.banner.name = data.film.name
  data.banner.cName = data.film.cName
  data.banner.picture = data.film.picture
  data.banner.year = parseInt(data.film.year)
  data.banner.remark = data.film.remarks
  ApiPost('/manage/banner/update', data.banner).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      data.banner = {id: '', mid: 0, name: '', cName: '', poster: '', picture: '', year: 0, remark: '', sort: 0}
      data.dialogV.bindV = false
      getBanners()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 删除海报信息
const delBanner = (b: any) => {
  ApiGet('/manage/banner/del', {id: b.id}).then((resp: any) => {
    if (resp.code === 0) {
      ElMessage.success({message: resp.msg})
      getBanners()
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 清除海报信息
const clearCache = () => {
  ApiGet('/cache/del').then((resp: any) => {
    if (resp.code == 0) {
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

// 上传并回填图片信息
const customUpload = (options: any) => {
  let file = options.file
  let formData = new FormData();
  formData.append("file", file)
  ApiPost(`/manage/file/upload`, formData).then((resp: any) => {
    if (resp.code === 0) {
      switch (options.data.type) {
        case 0:
          data.banner.poster = resp.data
          break
        case 1:
          data.banner.picture = resp.data
          break
      }
      ElMessage.success({message: resp.msg})
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}
const changePicture = (r: any) => {
  console.log(r)
}

const getBanners = () => {
  ApiGet(`/manage/banner/list`).then((resp: any) => {
    if (resp.code === 0) {
      data.banners = resp.data
    } else {
      ElMessage.error({message: resp.msg})
    }
  })
}

onMounted(() => {
  // 初始化banners数据
  getBanners()
})


</script>

<style scoped>

.upload_input {
  width: 76%;
}

.upload {
  height: 32px;
  margin-left: 10px;
}

.upload_btn {
  margin: 0 auto;
}

.film_view {
  max-width: 100%;
  display: flex;
  background: rgba(255, 255, 255, 0.25);
  padding: 16px;
  min-height: 200px;
  max-height: 200px;
  border-radius: 10px;
  margin: 16px 0;
}

.film_view a {
  flex: 1;
  border-radius: 8px;
  background-size: cover;
}

.film_intro {
  max-width: 75%;
  margin-left: 10px;
  flex: 3;
  /*flex-grow: 4;*/
  text-align: left;
  padding: 0 10px;
  font-size: 15px;
  position: relative;
}

.film_view h3, p, button {
  margin: 3px 0 3px 0;
}

.film_view p {
  max-width: 90%;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.film_view p em {
  font-weight: bold;
  margin-right: 8px;
}

.film_view button {
  background-color: orange;
  border-radius: 20px;
  border: none !important;
  color: #ffffff;
  font-weight: bold;
  position: absolute;
  margin-bottom: 2px;
  bottom: 0;
}

.tags {
  display: flex;
  width: 90%;
  justify-content: space-between;
}

.tags .tag_c {
  background: rgba(155, 73, 231, 0.72);
}

.tags span {
  border-radius: 5px;
  padding: 3px 5px;
  background: rgba(66, 66, 66);
  color: #c9c4c4;
  margin-right: 10px;
}

</style>