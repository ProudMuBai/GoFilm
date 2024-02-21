package config

import "time"

/*
 定义一些数据库存放的key值, 以及程序运行时的相关参数配置
*/

// -------------------------System Config-----------------------------------
const (

	// ListenerPort web服务监听的端口
	ListenerPort = "3601"

	// MAXGoroutine max goroutine, 执行spider中对协程的数量限制
	MAXGoroutine = 10

	FilmPictureUploadDir = "./static/upload/gallery"
	FilmPictureUrlPath   = "/upload/pic/poster/"
	FilmPictureAccess    = "/api/upload/pic/poster/"
)

// -------------------------redis key-----------------------------------
const (
	// CategoryTreeKey 分类树 key
	CategoryTreeKey     = "CategoryTree"
	CategoryTreeExpired = time.Hour * 24 * 90
	// MovieListInfoKey movies分类列表 key
	MovieListInfoKey = "MovieList:Cid%d"

	// MovieDetailKey movie detail影视详情信息 可以
	MovieDetailKey = "MovieDetail:Cid%d:Id%d"
	// MovieBasicInfoKey 影片基本信息, 简略版本
	MovieBasicInfoKey = "MovieBasicInfo:Cid%d:Id%d"

	// MultipleSiteDetail 多站点影片信息存储key
	MultipleSiteDetail = "MultipleSource:%s"

	// SearchInfoTemp redis暂存检索数据信息
	SearchInfoTemp = "Search:SearchInfoTemp"

	// SearchTitle 影片分类标题key
	SearchTitle = "Search:Pid%d:Title"
	// SearchTag 影片剧情标签key
	SearchTag = "Search:Pid%d:%s"

	// VirtualPictureKey 待同步图片临时存储 key
	VirtualPictureKey = "VirtualPicture"
	// MaxScanCount redis Scan 操作每次扫描的数据量, 每次最多扫描300条数据
	MaxScanCount = 300
)

const (
	AuthUserClaims = "UserClaims"
)

// -------------------------manage 管理后台相关key----------------------------------
const (
	// FilmSourceListKey 采集 API 信息列表key
	FilmSourceListKey = "Config:Collect:FilmSource"
	// ManageConfigExpired 管理配置key 长期有效, 暂定10年
	ManageConfigExpired = time.Hour * 24 * 365 * 10
	// SiteConfigBasic 网站参数配置
	SiteConfigBasic = "SystemConfig:SiteConfig:Basic"

	// FilmCrontabKey 定时任务列表信息
	FilmCrontabKey = "Cron:Task:Film"
	// DefaultUpdateSpec 每20分钟执行一次
	DefaultUpdateSpec = "0 */20 * * * ?"
	// DefaultUpdateTime 每次采集最近 3 小时内更新的影片
	DefaultUpdateTime = 3
)

// -------------------------Web API相关redis key-----------------------------------
const (
	// IndexCacheKey , 首页数据缓存
	IndexCacheKey = "IndexCache"
)

// -------------------------Database Connection Params-----------------------------------
const (
	// SearchTableName 存放检索信息的数据表名
	SearchTableName  = "search"
	UserTableName    = "users"
	UserIdInitialVal = 10000
	FileTableName    = "files"

	//mysql服务配置信息 root:root 设置mysql账户的用户名和密码

	MysqlDsn = "root:root@(mysql:3306)/FilmSite?charset=utf8mb4&parseTime=True&loc=Local"

	// Redis连接信息
	RedisAddr     = `redis:6379`
	RedisPassword = `root`
	RedisDBNo     = 0
)
