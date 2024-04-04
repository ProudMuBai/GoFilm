import {
    createRouter,
    createWebHistory,
} from "vue-router";

// 1.定义路由组件
import IndexHome from "../views/IndexHome.vue";
import Home from "../views/index/Home.vue";
import FilmDetails from "../views/index/FilmDetails.vue";
import Play from "../views/index/Play.vue";
import SearchFilm from "../views/index/SearchFilm.vue";
import NotFound from '../views/error/Error404.vue'
import FilmClassifySearch from "../views/index/FilmClassifySearch.vue";
import FilmClassify from "../views/index/FilmClassify.vue";
import ManageIndex from "../views/manage/Index.vue"
import Login from "../views/Login.vue"
import ManageHome from "../views/manage/ManageHome.vue";
import {getToken} from "../utils/token";
import CollectManage from "../views/manage/collect/CollectManage.vue";
import SiteConfig from "../views/manage/system/SiteConfig.vue";
import CronManage from "../views/manage/cron/CronManage.vue";
import Temp from "../views/manage/file/Temp.vue";
import FilmClass from "../views/manage/film/FilmClass.vue";
import Film from "../views/manage/film/Film.vue";
import FileUpload from "../views/manage/file/FileUpload.vue";
import FilmAdd from "../views/manage/film/FilmAdd.vue";
import CustomPlay from "../views/index/CustomPlay.vue";


// 2. 定义一个路由
const routes = [
    {
        path: '/',
        component: IndexHome,
        redirect: '/index',
        children: [
            {path: 'index', component: Home},
            {path: 'filmDetail', component: FilmDetails},
            {path: 'play', component: Play},
            {path: 'search', component: SearchFilm},
            {path: 'filmClassify', component: FilmClassify},
            {path: 'filmClassifySearch', component: FilmClassifySearch},
            {path: '/custom/player', component: CustomPlay},
        ]
    },
    {path: '/login', component: Login},
    {
        path: '/manage',
        component: ManageHome,
        redirect: '/manage/index',
        children: [
            {path: 'index', component: ManageIndex},
            {path: 'collect/index', component: CollectManage},
            {path: 'system/webSite', component: SiteConfig},
            {path: 'cron/index', component: CronManage},
            {path: 'file/upload', component: FileUpload},
            {path: 'file/gallery', component: Temp},
            {path: 'film', component: Film},
            {path: 'film/class', component: FilmClass},
            {path: 'film/add', component: FilmAdd},
            {path: 'film/detail', component: Temp},

        ]
    },
    {path: `/:pathMatch(.*)*`, component: NotFound},
]

// 创建路由实例并传递 routes配置
const router = createRouter({
    history: createWebHistory(),
    routes
})

// 添加全局前置守卫拦截未登录的跳转
router.beforeEach((to, from, next) =>{
    // 如果访问的是 /manage 下的路由, 且 token信息为空 则跳转到登录界面
    let matchPath = new RegExp(/^\/manage\//).test(to.path)
    let token = getToken()
    if ( matchPath && !token ) {
        next('/login')
    } else {
        next()
    }
})




export {router}

