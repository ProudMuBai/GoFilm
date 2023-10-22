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

