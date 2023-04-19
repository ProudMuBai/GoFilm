package config

import "time"

/*
定义一些数据库存放的key值
*/
const (

	// MAXGoroutine max goroutine, 执行spider中对协程的数量限制
	MAXGoroutine = 10

	// CategoryTreeKey 分类树 key
	CategoryTreeKey     = "CategoryTree"
	CategoryTreeExpired = time.Hour * 24 * 90
	// MovieListInfoKey movies分类列表 key
	MovieListInfoKey = "MovieList:Cid%d"

	// MovieDetailKey movie detail影视详情信息 可以
	MovieDetailKey = "MovieDetail:Cid%d:Id%d"
	// MovieBasicInfoKey 影片基本信息, 简略版本
	MovieBasicInfoKey = "MovieBasicInfoKey:Cid%d:Id%d"

	// MultipleSiteDetail 多站点影片信息存储key
	MultipleSiteDetail = "MultipleSource:%s"

	// SearchCount Search scan 识别范围
	SearchCount = 3000
	// SearchKeys Search Key Hash
	SearchKeys = "SearchKeys"
	// SearchScoreListKey 根据评分检索的key
	SearchScoreListKey = "Search:SearchScoreList"
	SearchTimeListKey  = "Search:SearchTimeList"
	SearchHeatListKey  = "Search:SearchHeatList"
	// SearchInfoTemp redis暂存检索数据信息
	SearchInfoTemp = "Search:SearchInfoTemp"

	// CornMovieUpdate 影片更新定时任务间隔
	CornMovieUpdate = "0 0/20 * * * ?"
	// UpdateInterval 获取最近几小时更新的影片 (h 小时) 默认3小时
	UpdateInterval = "3"
	// CornUpdateAll 每月28执行一次清库更新
	CornUpdateAll = "0 0 2 28 * ?"

	// SpiderCipher 设置Spider触发指令
	SpiderCipher = "Life in a different world from zero"
)

/*API相关redis key*/
const (
	IndexCacheKey = "IndexCache"
)

const (

	// SearchTableName 存放检索信息的数据表名
	SearchTableName = "search"

	//mysql服务配置信息 root:root 设置mysql账户的用户名和密码

	MysqlDsn = "root:root@(192.168.20.10:3307)/FilmSite?charset=utf8mb4&parseTime=True&loc=Local"

	// MysqlDsn docker compose 环境下的链接信息 mysql:3306 为 docker compose 中 mysql服务对应的网络名称和端口
	//MysqlDsn = "root:root@(mysql:3306)/FilmSite?charset=utf8mb4&parseTime=True&loc=Local"

	/*
		redis 配置信息
		RedisAddr host:port
		RedisPassword redis访问密码
		RedisDBNo 使用第几号库
	*/
	RedisAddr     = `192.168.20.10:6379`
	RedisPassword = `root`
	RedisDBNo     = 0

	// RedisAddr docker compose 环境下运行使用如下配置信息
	//RedisAddr     = `redis:6379`
	//RedisPassword = `root`
	//RedisDBNo     = 0
)
