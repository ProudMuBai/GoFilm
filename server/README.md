# Film Server

## 简介

- server 是本项目的后端项目
- 主要用于提供前端项目需要的 API数据接口, 以及数据搜集和更新
- 实现思路 : 
  - 使用 gocolly 获取公开的影视资源, 
  - 将请求数据通过程序处理整合成统一格式后使用redis进行暂存
  - 使用 mysql 存储收录的影片的检索信息, 用于影片检索, 分类
  - 使用 gin 作为web服务, 提供相应api接口
- 项目依赖

```go
# gin web服务框架, 用于处理与前端工程的交互
github.com/gin-gonic/gin v1.9.0
# gocolly go语言爬虫框架, 用于搜集公共影视资源
github.com/gocolly/colly/v2 v2.1.0
# go-redis redis交互程序
github.com/redis/go-redis/v9 v9.0.2
# gorm 用于处理与mysql数据库的交互
gorm.io/gorm v1.24.6
gorm.io/driver/mysql v1.4.7
```



## 项目结构

> 项目主要目录结构

- config				用于存放项目中使用的配置信息和静态常量
- controller	      请求处理控制器
- logic                请求处理逻辑实现
- model              数据模型结构体以及与数据库交互
- plugin              项目所需的插件工具集合
  - common     公共依赖
  - db              数据库配置信息
  - spider         gocolly配置, 执行逻辑, 数据前置处理等

```text
server                          
├─ config                       
│  └─ DataConfig.go             
├─ controller                   
│  ├─ IndexController.go        
│  └─ SpiderController.go       
├─ logic                        
│  ├─ IndexLogic.go             
│  └─ SpiderLogic.go            
├─ model                        
│  ├─ Categories.go             
│  ├─ Movies.go                 
│  ├─ RequestParams.go          
│  ├─ ResponseJson.go           
│  └─ Search.go                 
├─ plugin                       
│  ├─ common                    
│  │  ├─ dp                     
│  │  │  ├─ ProcessCategory.go  
│  │  │  └─ ProcessMovies.go    
│  │  ├─ param                  
│  │  │  └─ SimpleParam.go      
│  │  └─ util                   
│  │     ├─ FileDownload.go     
│  │     └─ Request.go          
│  ├─ db                        
│  │  ├─ mysql.go               
│  │  └─ redis.go               
│  └─ spider                    
│     ├─ Spider.go              
│     └─ SpiderCron.go          
├─ router                       
│  └─ router.go                 
├─ go.mod                       
├─ go.sum                       
├─ main.go                      
└─ README.md                    
```



## 启动方式

### 本地运行

1.  修改 /server/plugin/db 目录下的 mysql.go 和 redis.go 中的连接地址和用户名密码
2. 在 server 目录下执行 `go run main.go`





## 数据库信息简介

#### 1.Mysql  

> 连接信息(以docker compose部署为例) :

```yaml
 mysql:
 	ip: 部署的服务器IP
    port: 3610
    username: root
    password: root
    database: FilmSite
```

> 数据库结构

- 数据库: FilmSite
  - 数据表 search

> search 表 (用于记录影片的相关检索信息, 主要用于影片的 搜索, 分类, 排序 等)

| 字段名称     | 类型     | 字段释义               |
| ------------ | -------- | ---------------------- |
| id           | bigint   | 自增主键               |
| created_at   | datetime | 记录创建时间           |
| updated_at   | datetime | 记录更新时间           |
| deleted_at   | datetime | 逻辑删除字段           |
| mid          | bigint   | 影片ID                 |
| cid          | bigint   | 二级分类ID             |
| pid          | bigint   | 一级分类ID             |
| name         | varchar  | 影片名称               |
| sub_title    | varchar  | 子标题(影片别名)       |
| c_name       | varchar  | 分类名称               |
| class_tag    | varchar  | 剧情标签               |
| area         | varchar  | 地区                   |
| language     | varchar  | 语言                   |
| year         | bigint   | 上映年份               |
| initial      | varchar  | 首字母                 |
| score        | double   | 豆瓣评分               |
| update_stamp | bigint   | 影片更新时间戳         |
| hits         | bigint   | 热度(播放次数)         |
| state        | varchar  | 状态(正片)             |
| remarks      | varchar  | 更新状态(完结 \| xx集) |
| release_data | bigint   | 上映时间戳             |



#### 2.Redis

> 连接信息(以docker compose部署为例) : 

```yaml
  ## 部署时默认使用如下信息
  redis:
  	ip: 部署的服务器IP
    port: 3620
    password: root
    DB: 0  ##使用的redis数据库为0号库
```



## 服务端API数据示例

### 1. 网站前台API

#### 1. API接口基本信息

- 响应结构

```text
{
    code: 0|1,		// 成功|失败
	data: {},	// 数据内容
    msg: "",		// 提示信息
}
```



| 名称               | URL                 | client component                              | Method | Params                                                       |
| ------------------ | :------------------ | --------------------------------------------- | ------ | ------------------------------------------------------------ |
| 首页数据           | /index              | client/src/views/index/Home.vue               | GET    | 无                                                           |
| 网站基本配置信息   | /config/basic       | client/src/components/index/Header.vue        | GET    | 无                                                           |
| 影片分类导航       | /navCategory        | client/src/components/index/Header.vue        | GET    | 无                                                           |
| 影片详情           | /filmDetail         | client/src/views/index/FilmDetails.vue        | GET    | id   (int, 影片ID)                                           |
| 影片播放页数据     | /filmPlayInfo       | client/src/views/index/Play.vue               | GET    | id   (int, 影片ID) <br>playFrom   (string, 播放源ID)<br>episode   (int, 集数索引) |
| 影片检索(名称搜索) | /searchFilm         | client/src/views/index/SearchFilm.vue         | GET    | keyword   (string, 影片名)                                   |
| 影片分类首页       | /filmClassify       | client/src/views/index/FilmClassify.vue       | GET    | Pid   (int, 一级分类ID)                                      |
| 影片分类详情页     | /filmClassidySearch | client/src/views/index/FilmClassifySearch.vue | GET    | Pid   (int, 一级分类ID)<br>Category   (int, 二级分类ID)<br>Plot   (string, 剧情)<br>Area   (string, 地区)<br>Language   (string, 语言)<br>Year   (string, 年份)<br>Sort   (string, 排序方式) |

#### 2. 接口响应数据示例:

-  `/index` 首页数据

```text
{
    "code": 0,		// 状态码
    "data": {		// 数据内容
        "category": {				// 分类信息
            "id": 0,				// 分类ID
            "name": "xxx",			// 分类名称
            "pid": 0,				// 上级分类ID
            "show": false,			// 是否展示
            "children": [], 			// 子分类信息
        },
        "content": [				// 内容区数据
            {
                "hot": [			// 热播影片
                    {
                        "CreatedAt": "2024-01-13T19:04:01+08:00",		// 创建时间
                        "DeletedAt": null,				// 删除时间
                        "ID": 100,						// ID
                        "UpdatedAt": "2024-01-13T19:04:01+08:00",	// 更新时间
                        "area": "xxx",					// 地区
                        "cName": "xxx",					// 分类名称
                        "cid": 45,					// 分类ID
                        "classTag": "xxx",				// 剧情标签
                        "hits": 0,					// 热度
                        "initial": "X",					// 首字母
                        "language": "xxx",				// 语言 
                        "mid": 10000,					// 影片ID
                        "name": "xxx",					// 影片名称
                        "pid": 1,					// 上级分类ID
                        "releaseStamp": 1704880403,		// 上映时间戳
                        "remarks": "xxx",			 	// 备注信息 [预告|完结|更新至xx集]
                        "score": 0,						// 评分
                        "state": "xx",					// 状态 正片|预告
                        "subTitle": "xxx",				// 子标题, 别名
                        "updateStamp": 1704880403,		// 更新时间戳
                        "year": 2024,					// 年份
                    }
                ],	
        		"movies": [			// 近期更新影片
                    {
                        "id": 10000,						// 影片ID
                        "cid": 6,					// 分类ID
                        "pid": 1,					// 上级分类ID
                        "name": "xxxx",						// 影片名称
                        "subTitle": "xxxx",					// 子标题, 别名
                        "cName": "xxx",						// 分类名称
                        "state": "正片",						// 影片状态
                        "picture": "http://xxxx.jpg",		// 海报图片url
                        "actor": "xxx,xxx", 				// 演员
                        "director": "xxx,xxx",				// 导演
                        "blurb": "",						// 剧情简介
                        "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                        "area": "xxx",						// 地区
                        "year": "2024" 						// 年份
                    }
                ],
        		"nav": [						// 导航信息
                    {	
                    	"id": 0,				// 分类ID
       					"name": "xxxx", 		// 分类名称
       					"pid": 0,				//上级分类ID
       					"show": false,			// 是否展示
       					"children": [], 		//子分类信息
                    }
                ]
            },
        ]
    },
    msg: "", 	// 提示信息
}
```

- `/config/basic` 网站基本配置信息

```text
{
    "code": 0,
    "data": {
        "siteName": "GoFilm",					// 网站名称
        "domain": "http://127.0.0.1:3600",			// 域名
        "logo": "https://xxx.jpg",				// 网站logo
        "keyword": "xxxx, xxxx",				// 网站搜索关键字
        "describe": "xxxxxxx",					// 网站描述信息
        "state": true,						//站点状态
        "hint": "网站升级中, 暂时无法访问 !!!" 		// 网站关闭时提示信息
    },
    "msg": ""
}
```

- `/navCategory` 首页头部分类信息

```text
{
    "code": 0,
    "data": [
             {	
                    "id": 0,				// 分类ID
       				"name": "xxxx", 			// 分类名称
       				"pid": 0,					// 上级分类ID
       				"show": false,				// 是否展示
              },
    ],
    "msg": ""
}
```

- `  /filmDetail` 影片详情信息

```text
 {
    "code": 0,
    "data": {
        "detail": {									// 影片详情信息
            "id": 100000,							// 影片ID
            "cid": 30,								// 影片分类ID
            "pid": 4,								// 上级分类ID
            "name": "xxx",							// 影片名称
            "picture": "https://xxx.jpg",			// 海报封面url
            "playFrom": [ "xxx","xxx" ],			// 播放来源
            "DownFrom": "http",						// 下载方式
            "playList": [ 							// 播放地址列表(主站点)
                {
                    "episode": "第xx集",					// 集数
                    "link": "https://xxx/index.m3u8"		// 播放地址url
                },
            ],	
            "downloadList": [ 						// 下载地址列表 
            	 {
                    "episode": "第xx集",					// 集数
                    "link": "https://xxx/index.m3u8"			// 播放地址url
               	 },
            ],	
            "descriptor": { 						// 影片详情
            	"subTitle": "",						// 副标题, 别名
                "cName": "xxxx",					// 分类名称
                "enName": "xxx",					// 影片名称中文拼音
                "initial": "X",						// 影片名称首字母
                "classTag": "xxxx",					// 内容标签
                "actor": "xxx,xxx",					// 演员
                "director": "xxx",					// 导演
                "writer": "xxx",					// 作者
                "blurb": "xxx",						// 简介(缺省)
                "remarks": "更新至第xx集",			// 更新进度
                "releaseDate": "2024-01-06",		// 上映日期
                "area": "xxx",						// 地区
                "language": "xxx",					// 语言
                "year": "2024",						// 年份
                "state": "正片",					// 状态 正片|预告
                "updateTime": "2024-01-13 00:51:21",		// 更新时间
                "addTime": 1704511497,				// 添加时间戳
                "dbId": 26373174,					// 豆瓣ID
                "dbScore": "0.0",					// 豆瓣评分
                "hits": 0,							// 热度
                "content": "xxx"					//影片内容简介(全)
            },	
            "list": [ 								// 播放地址列表(全站点)
            	{
                    "id": "xxxxxxxxxxxx",			// 播放源ID
                    "name": "HD(xxx)",				// 播放源别名
                    "linkList": [					// 播放地址列表
                         {
                            "episode": "第xx集",			// 集数
                            "link": "https://xxx/index.m3u8"		// 播放地址url
                         },
                    ]
                },
            ]	
        },
        "relate": [ 		// 相关影片推荐
        	{
              	"id": 10000,					// 影片ID
                "cid": 6,						// 分类ID
                "pid": 1,						// 上级分类ID
                "name": "xxxx",					// 影片名称
                "subTitle": "xxxx",				// 子标题, 别名
                "cName": "xxx",					// 分类名称
                "state": "xxx",					// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 			// 演员
                "director": "xxx,xxx",			// 导演
                "blurb": "",					// 剧情简介
                "remarks": "HD", 				// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",					// 地区
                "year": "2024" 					// 年份
            },
        ]
    },
    "msg": "xxx"
}
```

- `  /filmPlayInfo` 影片播放页信息

```text
{
    "code": 0,
    "data": {
        "current": { 	// 当前播放信息
            "episode": "第xx集",					   // 当前播放集数
            "link": "https://xxx/index.m3u8"		// 当前播放地址url
        },
        "currentEpisode": 0,			// 当前播放集数索引
        "currentPlayFrom": "xxx",		// 当前播放源ID
        "detail": { 		// 影片详情
            "id": 100000,							// 影片ID
            "cid": 30,								// 影片分类ID
            "pid": 4,								// 上级分类ID
            "name": "xxx",							// 影片名称
            "picture": "https://xxx.jpg",			// 海报封面url
            "playFrom": [ "xxx","xxx" ],			// 播放来源
            "DownFrom": "http",						// 下载方式
            "playList": [ 	// 播放地址列表(主站点)
                {
                    "episode": "第xx集",							   // 集数
                    "link": "https://xxx/index.m3u8"// 播放地址url
                },
            ],	
            "downloadList": [ 	// 下载地址列表 
                 {
                    "episode": "第xx集",				// 集数
                    "link": "https://xxx/index.m3u8"	// 播放地址url
                 },
            ],	
            "descriptor": { 	// 影片详情
                "subTitle": "",						// 副标题, 别名
                "cName": "xxxx",					// 分类名称
                "enName": "xxx",					// 影片名称中文拼音
                "initial": "X",						// 影片名称首字母
                "classTag": "xxxx",					// 内容标签
                "actor": "xxx,xxx",					// 演员
                "director": "xxx",					// 导演
                "writer": "xxx",					// 作者
                "blurb": "xxx",						// 简介(缺省)
                "remarks": "更新至第xx集",			// 更新进度
                "releaseDate": "2024-01-06",		// 上映日期
                "area": "xxx",						// 地区
                "language": "xxx",					// 语言
                "year": "2024",						// 年份
                "state": "xxx",						// 状态 正片|预告
                "updateTime": "2024-01-13 00:51:21",	// 更新时间
                "addTime": 1704511497,				// 添加时间戳
                "dbId": 26373174,					// 豆瓣ID
                "dbScore": "0.0",					// 豆瓣评分
                "hits": 0,							// 热度
                "content": "xxx"					//影片内容简介(全)
            },	
            "list": [ 		// 播放地址列表(全站点)
                {
                    "id": "xxxxxxxxxxxx",			// 播放源ID
                    "name": "HD(xxx)",				// 播放源别名
                    "linkList": [					// 播放地址列表
                         {
                            "episode": "第xx集",					   // 集数
                            "link": "https://xxx/index.m3u8"		// 播放地址url
                         },
                    ]
                },
            ]	
        },							
        "relate": [ 		// 相关影片推荐
            {
                "id": 10000,						// 影片ID
                "cid": 6,							// 分类ID
                "pid": 1,							// 上级分类ID
                "name": "xxxx",						// 影片名称
                "subTitle": "xxxx",					// 子标题, 别名
                "cName": "xxx",						// 分类名称
                "state": "xxx",						// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 				// 演员
                "director": "xxx,xxx",				// 导演
                "blurb": "",						// 剧情简介
                "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",						// 地区
                "year": "2024" 						// 年份
            },
        ]
    },
    "msg": "影片播放信息获取成功"
}
```

- `/filmClassify` 分类影片首页数据

```text
{
    "code": 0,
    "data": {
        "content": {		// 内容区数据
            "news": [		//最新上映
            	 "id": 10000,						// 影片ID
                "cid": 6,							// 分类ID
                "pid": 1,							// 上级分类ID
                "name": "xxxx",						// 影片名称
                "subTitle": "xxxx",					// 子标题, 别名
                "cName": "xxx",						// 分类名称
                "state": "xxx",						// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 				// 演员
                "director": "xxx,xxx",				// 导演
                "blurb": "",						// 剧情简介
                "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",						// 地区
                "year": "2024" 						// 年份
            ],
            "recent": [ 	// 近期更新
            	 "id": 10000,						// 影片ID
                "cid": 6,							// 分类ID
                "pid": 1,							// 上级分类ID
                "name": "xxxx",						// 影片名称
                "subTitle": "xxxx",					// 子标题, 别名
                "cName": "xxx",						// 分类名称
                "state": "xxx",						// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 				// 演员
                "director": "xxx,xxx",				// 导演
                "blurb": "",						// 剧情简介
                "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",						// 地区
                "year": "2024" 						// 年份
            ],
            "top": [ 		// 热度排行
            	 "id": 10000,						// 影片ID
                "cid": 6,							// 分类ID
                "pid": 1,							// 上级分类ID
                "name": "xxxx",						// 影片名称
                "subTitle": "xxxx",					// 子标题, 别名
                "cName": "xxx",						// 分类名称
                "state": "xxx",						// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 				// 演员
                "director": "xxx,xxx",				// 导演
                "blurb": "",						// 剧情简介
                "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",						// 地区
                "year": "2024" 						// 年份
            ]
        },
        "title": { 			// 头部标题区数据(暂未使用)
        	 "id": 0,						// 分类ID
            "name": "xxx", 					// 分类名称
            "pid": 0,						// 上级分类ID
            "show": false,					// 是否展示
            "children": [], 				// 子分类信息
        }
    },
    "msg": ""
}
```

-  ` /filmClassidySearch` 影片分类检索页数据

```text
{
    "code": 0,
    "data": {
        "list": [ 		// 影片信息集合
        	{
                "id": 10000,						// 影片ID
                "cid": 6,							// 分类ID
                "pid": 1,							// 上级分类ID
                "name": "xxxx",						// 影片名称
                "subTitle": "xxxx",					// 子标题, 别名
                "cName": "xxx",						// 分类名称
                "state": "xxx",						// 影片状态
                "picture": "http://xxxx.jpg",		// 海报图片url
                "actor": "xxx,xxx", 				// 演员
                "director": "xxx,xxx",				// 导演
                "blurb": "",						// 剧情简介
                "remarks": "HD", 					// 备注信息 [预告|完结|更新至xx集]
                "area": "xxx",						// 地区
                "year": "2024" 						// 年份
            }
        ],
        "page": { 		// 分页信息
        	"pageSize": 49,							// 每页页数						
            "current": 1,							// 当前页
            "pageCount": xx,						// 总页数
            "total": xx,							// 总数据量
        },
        "params": { 	// 请求参数
        	"Area": "",								// 地区
            "Category": "",							// 分类ID
            "Language": "",							// 语言
            "Pid": "1",								// 上级分类ID
            "Plot": "",								// 剧情
            "Sort": "xxx",							// 排序方式
            "Year": "",								// 年份
        },
        "search": { 	// 分类标签组信息
        	"sortList": [ "Category","Plot","Area","Language","Year","Sort" ], 		// 标签数据排序, 固定值
            "tags": { 			// 标签组, 用于页面筛选Tag渲染
            	"Area": [ { Name:"", Value:"" } ],								
                "Category": [ { Name:"", Value:"" } ],
                "Initial": [ { Name:"", Value:"" } ],
                "Language": [ { Name:"", Value:"" } ],
                "Plot": [ { Name:"", Value:"" } ],
                "Sort": [ { Name:"", Value:"" } ],
                "Year": [ { Name:"", Value:"" } ]
            },
            "titles": { 		// 标签组标题映射(固定值)
            	"Area": "地区",
                "Category": "类型",
                "Initial": "首字母",
                "Language": "语言",
                "Plot": "剧情",
                "Sort": "排序",
                "Year": "年份"
            }
        },				
        "title": { 		// 当前一级分类信息
        	"id": 1,				// 分类ID
            "pid": 0,				// 上级分类ID
            "name": "xxx",			// 分类名称
            "show": true,			// 是否展示
        }
    },
    "msg": ""
}
```















