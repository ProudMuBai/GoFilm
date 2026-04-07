# GoFilm

一个基于 vue 和 gin 实现的在线观影网站

效果展示: <a href="https://m.mubai.link/" target="_blank">点击访问演示站点</a>  

新版本测试访问站点: [新版本站点](https://www.mubai.cn.mt/)  [备用域名](https://m2.mubai.link/)

## 简介

**GoFilm** 

项目采用vite + vue作为前端技术栈, 使用 ElementPlus 作为UI 框架进行开发

后端程序使用 Gin + gorm + go-redis 等相关框架提供接口服务, 使用 gocolly 和 robfig/cron 进行公共影视资源采集和定时更新功能



## 项目部署

**部署方式**

- [Docker部署](https://github.com/ProudMuBai/GoFilm/blob/main/film/README.md) 
- [1Panel部署(可视化面板操作)](https://blog.mubai.link/procedure/application/github/GoFilm/)

**使用指南**

- 程序使用文档:  提供项目安装部署以及相应的初始化使用步骤说明  [GoFilm 部署文档](https://blog.mubai.link/procedure/application/github/GoFilm/)

- 后端项目路径 `GoFilm/server`, 包含 项目结构说明,  后端程序源码, API接口说明, 本地启动注意事项 [查看](https://github.com/ProudMuBai/GoFilm/tree/main/server)
- 前端项目路径 `GoFilm/client`,  包含 项目结构说明,  前端项目源码,  配置文件说明, 本地启动方式  [查看](https://github.com/ProudMuBai/GoFilm/tree/main/client)
- 部署文件 `GoFilm/film`, 包含项目部署所需的所有文件以及相应的说明文件  [查看](https://github.com/ProudMuBai/GoFilm/tree/main/film)



## 新版本说明

**网站前台**

- 对新功能和目前功能有需要改善和补充的可以去issue #71下方留言, 一月中旬休假可以集中处理

- 前台部分对网站名称以及播放源等部分信息与后台数据进行关联, 可通过后台进行修改
- 影片详情部分以及首页导航数据结构发生变化, 样式保持一致
- 默认访问地址: `服务器IP:默认端口 [http://127.0.0.1/index]`

**管理后台**

- 新增管理后台功能组, 主要用于对 `采集站点`, `定时更新`, `网站基本信息`, `图片同步`, `影片分类`, `影片信息` 等进行管理 (部分功能正在完善中, 不影响已有功能使用)
- 管理后台访问需进行登录, 默认账号/密码: `admin admin` (登录成功后自行通过右上下拉弹窗进行密码修改)
- 具体情况请自行搭建访问
- 默认访问地址: `服务器IP:默认端口/manage [http://127.0.0.1:3600/manage]`

**更新说明**

- 更换主体数据存储方式, 精简部分无意义字段, 节省服务器内存
- 更换视频播放器, 优化部分功能键,新增播放列表和下一集控键
- 新版本安装方法以及使用说明请查看本项目 film 文件夹下的说明文件

>更新内容:
>
>- 采集优化, 采集逻辑大幅度改动, 主站点采集前进行分类设置(将不需要的分类设置为不展示), 可精简采集内容
>- 数据存储方式变更, 主附站点数据存储转为MySQL存储, 减少Redis内存消耗
>- 更新采集内容缓存, 通过定时任务(每周日凌晨四点)统一同步到MySQL, 降低MySQL数据更新频率
>- 资源匹配逻辑修改, 变更对不同站点的影片匹配逻辑
>- DockerFile 更新, 由于golang版本以及相关依赖升级导致旧文件无法正常构建镜像. 
>- UI 不分列表展示区域数据重复问题
>- 播放页逻辑变更, 更换播放器为 `xgplayer`
>
>后续计划:
>
>- 有无前端UI大佬救一救啊, 界面改的没思绪了.
>- 新增功能测试 && buf修复
>- 优化UI组件
>- 针对播放器进行BUG测试优化以及功能完善
>- 采集方式细节化, 实现采集数据的精简
>- 聚合采集, 实现通过影片名称搜索动态添加数据到站点中

## 目录结构

- client 客户端项目目录 [Client简介](./client/README.md)
- server 服务端接口项目目录 [Server简介](./client/README.md)
- film 项目部署相关配置目录 [film 项目安装](./film/README.md)
- 详细说明请查看具体目录中的README文件

```text
GoFilm-main                            
├─ client                              
│  ├─ public                           
│  │  └─ favicon.ico                   
│  ├─ src                              
│  │  ├─ assets                        
│  │  │  ├─ css                        
│  │  │  │  ├─ classify.css            
│  │  │  │  ├─ film.css                
│  │  │  │  └─ pagination.css          
│  │  │  └─ image                      
│  │  │     ├─ 404.png                 
│  │  │     └─ play.png                
│  │  ├─ components                    
│  │  │  ├─ Loading                    
│  │  │  │  ├─ index.ts                
│  │  │  │  └─ Loading.vue             
│  │  │  ├─ FilmList.vue               
│  │  │  ├─ Footer.vue                 
│  │  │  ├─ Header.vue                 
│  │  │  ├─ RelateList.vue             
│  │  │  └─ Util.vue                   
│  │  ├─ router                        
│  │  │  └─ router.ts                  
│  │  ├─ utils                         
│  │  │  ├─ cookie.ts                  
│  │  │  └─ request.ts                 
│  │  ├─ views                         
│  │  │  ├─ error                      
│  │  │  │  └─ Error404.vue            
│  │  │  ├─ index                      
│  │  │  │  ├─ FilmClassify.vue        
│  │  │  │  ├─ FilmClassifySearch.vue  
│  │  │  │  ├─ FilmDetails.vue         
│  │  │  │  ├─ Home.vue                
│  │  │  │  ├─ Play.vue                
│  │  │  │  └─ SearchFilm.vue          
│  │  │  └─ IndexHome.vue              
│  │  ├─ App.vue                       
│  │  ├─ main.ts                       
│  │  ├─ style.css                     
│  │  └─ vite-env.d.ts                 
│  ├─ auto-imports.d.ts                
│  ├─ components.d.ts                  
│  ├─ index.html                       
│  ├─ package.json                     
│  ├─ README.md                        
│  ├─ tsconfig.json                    
│  ├─ tsconfig.node.json               
│  └─ vite.config.ts                   
├─ film                                
│  ├─ data                             
│  │  ├─ nginx                         
│  │  │  ├─ html                       
│  │  │  │  ├─ assets                  
│  │  │  │  │  ├─ 404-b813c94a.png     
│  │  │  │  │  ├─ index-984712d6.js    
│  │  │  │  │  ├─ index-de4c7ff5.css   
│  │  │  │  │  └─ play-bb9c8990.png    
│  │  │  │  ├─ favicon.ico             
│  │  │  │  └─ index.html              
│  │  │  └─ nginx.conf                 
│  │  └─ redis                         
│  │     └─ redis.conf                 
│  ├─ server                           
│  │  ├─ config                        
│  │  │  └─ DataConfig.go              
│  │  ├─ controller                    
│  │  │  ├─ IndexController.go         
│  │  │  └─ SpiderController.go        
│  │  ├─ logic                         
│  │  │  ├─ IndexLogic.go              
│  │  │  └─ SpiderLogic.go             
│  │  ├─ model                         
│  │  │  ├─ Categories.go              
│  │  │  ├─ Movies.go                  
│  │  │  ├─ RequestParams.go           
│  │  │  ├─ ResponseJson.go            
│  │  │  └─ Search.go                  
│  │  ├─ plugin                        
│  │  │  ├─ common                     
│  │  │  │  ├─ dp                      
│  │  │  │  │  ├─ ProcessCategory.go   
│  │  │  │  │  └─ ProcessMovies.go     
│  │  │  │  └─ param                   
│  │  │  │     └─ SimpleParam.go       
│  │  │  ├─ db                         
│  │  │  │  ├─ mysql.go                
│  │  │  │  └─ redis.go                
│  │  │  └─ spider                     
│  │  │     ├─ Spider.go               
│  │  │     ├─ SpiderCron.go           
│  │  │     └─ SpiderRequest.go        
│  │  ├─ router                        
│  │  │  └─ router.go                  
│  │  ├─ go.mod                        
│  │  ├─ go.sum                        
│  │  ├─ main.go                       
│  │  └─ README.md                     
│  ├─ docker-compose.yml               
│  ├─ Dockerfile                       
│  └─ README.md                        
├─ server                              
│  ├─ config                           
│  │  └─ DataConfig.go                 
│  ├─ controller                       
│  │  ├─ IndexController.go            
│  │  └─ SpiderController.go           
│  ├─ logic                            
│  │  ├─ IndexLogic.go                 
│  │  └─ SpiderLogic.go                
│  ├─ model                            
│  │  ├─ Categories.go                 
│  │  ├─ Movies.go                     
│  │  ├─ RequestParams.go              
│  │  ├─ ResponseJson.go               
│  │  └─ Search.go                     
│  ├─ plugin                           
│  │  ├─ common                        
│  │  │  ├─ dp                         
│  │  │  │  ├─ ProcessCategory.go      
│  │  │  │  └─ ProcessMovies.go        
│  │  │  ├─ param                      
│  │  │  │  └─ SimpleParam.go          
│  │  │  └─ util                       
│  │  │     ├─ FileDownload.go         
│  │  │     └─ Request.go              
│  │  ├─ db                            
│  │  │  ├─ mysql.go                   
│  │  │  └─ redis.go                   
│  │  └─ spider                        
│  │     ├─ Spider.go                  
│  │     └─ SpiderCron.go              
│  ├─ router                           
│  │  └─ router.go                     
│  ├─ go.mod                           
│  ├─ go.sum                           
│  ├─ main.go                          
│  └─ README.md                        
├─ LICENSE                             
└─ README.md                           
```



## 起源

从正式接触编程语言到第一次动手敲代码, , 当时有动手做一些东西的想法,也正是在那时喜欢追番迷二次元, 曾想过做一个自己的动漫站,

但因为知识面匮乏, 总是在进行到某一步时就会遇到一些盲区, 从最开始的静态页面到后面的伪数据, 也实现过一些当时能做到的部分, 

后面慢慢学习的过程中也渐渐遗忘了这个想法, 但因为一些偶然的因素, 想要做一个自己的开源项目, 于是就从零开始慢慢实现并完善了这个

影视站的各个部分, 期间也一点点修改颠覆了一些最开始的思路, 但目前主体功能基本完善, 后续也会定期进行一些bug修复和新功能的更新

如有发现Bug, 或者有好的建议, 可以进行反馈, 欢迎各位大佬来指点一二



## 更新迭代计划

- 目前用户界面的一些功能有待开发和完善, 大家也可以继续提供一些好的建议
- 目前pc端的历史记录写了一个简单的测试版, 后面有时间会同步完善pc和wrap端的历史记录和收藏功能
- 前台功能目前基本满足观看的需求, 后续考虑切入一些登录和账户以及管理后台的功能,慢慢完善这个项目.



