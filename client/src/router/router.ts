
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
import CategoryFilm from "../views/index/CategoryFilm.vue";
import NotFound from '../views/error/Error404.vue'
import FilmClassifySearch from "../views/index/FilmClassifySearch.vue";
import FilmClassify from "../views/index/FilmClassify.vue";


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
            {path: 'CategoryFilm', component: CategoryFilm},
            {path: 'filmClassify', component: FilmClassify},
            {path: 'filmClassifySearch', component: FilmClassifySearch},
        ]
    },
    {path: `/:pathMatch(.*)*`, component: NotFound},
]

// 创建路由实例并传递 routes配置
const router = createRouter({
    history: createWebHistory(),
    routes
})


export {router}

