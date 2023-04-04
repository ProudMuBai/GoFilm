package spider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"server/config"
	"server/model"
	"server/plugin/common"
	"strings"
	"sync"
	"time"
)

/*
		公共资源采集站点
	 1. 视频列表请求参数
	    ac=list 列表数据, t 影视类型ID, pg 页码, wd 关键字, h 几小时内数据
	 2. 视频详情请求参数
	    ac=detail 详情数据, ids 影片id列表, h, pg, t 影视类型ID
*/
const (
	LZ_MOVIES_URL    = "https://cj.lziapi.com/api.php/provide/vod/"
	LZ_MOVIES_Bk_URL = "https://cj.lzcaiji.com/api.php/provide/vod/"
	TK_MOVIES_URL    = "https://api.tiankongapi.com/api.php/provide/vod"
	KC_MOVIES_URL    = "https://caiji.kczyapi.com/api.php/provide/vod/"
	FS_MOVIES_URL    = "https://www.feisuzyapi.com/api.php/provide/vod/"

	// FILM_COLLECT_SITE 当前使用的采集URL
	FILM_COLLECT_SITE = "https://www.feisuzyapi.com/api.php/provide/vod/"
)

// 定义一个同步等待组
var wg = &sync.WaitGroup{}

func StartSpider() {
	// 1. 先拉取全部分类信息
	CategoryList()

	//2. 拉取所有分类下的影片基本信息
	tree := model.GetCategoryTree()
	AllMovies(&tree)
	wg.Wait()
	log.Println("AllMovies 影片列表获取完毕")

	// 3. 获取入库的所有影片详情信息
	// 3.2 获取入库的所有影片的详情信息
	AllMovieInfo()
	log.Println("AllMovieInfo 所有影片详情获取完毕")

	// 4. mysql批量插入与数据爬取同时进行容易出现主键冲突, 因此滞后
	// 4.1 先一步将输入存入redis中, 待网络io结束后再进行分批扫描入库
	// 3.1 先查找并创建search数据库
	time.Sleep(time.Second * 10)
	model.CreateSearchTable()
	SearchInfoToMdb()
	log.Println("SearchInfoToMdb 影片检索信息保存完毕")
	time.Sleep(time.Second * 10)
}

// CategoryList 获取分类数据
func CategoryList() {
	// 设置请求参数信息
	r := RequestInfo{Uri: FILM_COLLECT_SITE, Params: url.Values{}}
	r.Params.Set(`ac`, "list")
	r.Params.Set(`pg`, "1")
	r.Params.Set(`t`, "1")
	// 执行请求, 获取一次list数据
	ApiGet(&r)
	// 解析resp数据
	movieListInfo := model.MovieListInfo{}
	if len(r.Resp) <= 0 {
		log.Println("MovieListInfo数据获取异常 : Resp Is Empty")
	}
	_ = json.Unmarshal(r.Resp, &movieListInfo)
	// 获取分类列表信息
	classList := movieListInfo.Class
	// 组装分类数据信息树形结构
	categoryTree := common.CategoryTree(classList)
	// 序列化tree
	data, _ := json.Marshal(categoryTree)
	// 保存 tree 到redis
	err := model.SaveCategoryTree(string(data))
	if err != nil {
		log.Println("SaveCategoryTree Error: ", err)
	}
}

// AllMovies 遍历所有分类, 获取所有二级分类数据
func AllMovies(tree *model.CategoryTree) {
	// 遍历一级分类
	for _, c := range tree.Children {
		// 遍历二级分类, 屏蔽主页不需要的影片信息, 只获取 电影1 电视剧2 综艺3 动漫4等分类下的信息
		//len(c.Children) > 0 && c.Id <= 4
		if len(c.Children) > 0 {
			for _, cInfo := range c.Children {
				//go CategoryAllMovie(cInfo.Category)
				CategoryAllMoviePlus(cInfo.Category)
			}
		}
	}
}

// CategoryAllMovie 获取指定分类的所有影片基本信息
func CategoryAllMovie(c *model.Category) {
	// 添加一个等待任务, 执行完减去一个任务
	wg.Add(1)
	defer wg.Done()
	// 设置请求参数
	r := &RequestInfo{Uri: FILM_COLLECT_SITE, Params: url.Values{}}
	r.Params.Set(`ac`, "list")
	r.Params.Set(`t`, fmt.Sprint(c.Id))
	ApiGet(r)
	// 解析请求数据
	listInfo := model.MovieListInfo{}
	_ = json.Unmarshal(r.Resp, &listInfo)
	// 获取pageCount信息, 循环获取所有页数据
	pageCount := listInfo.PageCount
	// 开始获取所有信息, 使用协程并发获取数据
	for i := 1; i <= int(pageCount); i++ {
		// 使用新的 请求参数
		r.Params.Set(`pg`, fmt.Sprint(i))
		// 保存当前分类下的影片信息
		info := model.MovieListInfo{}
		ApiGet(r)
		// 如果返回数据中的list为空,则直接结束本分类的资源获取
		if len(r.Resp) <= 0 {
			log.Println("SaveMoves Error Response Is Empty")
			break
		}
		_ = json.Unmarshal(r.Resp, &info)
		if info.List == nil {
			log.Println("MovieList Is Empty")
			break
		}
		// 处理影片信息
		list := common.ProcessMovieListInfo(info.List)
		// 保存影片信息至redis
		_ = model.SaveMoves(list)
	}
}

// CategoryAllMoviePlus  部分分类页数很多,因此采用单分类多协程拉取
func CategoryAllMoviePlus(c *model.Category) {
	// 设置请求参数
	r := &RequestInfo{Uri: FILM_COLLECT_SITE, Params: url.Values{}}
	r.Params.Set(`ac`, "list")
	r.Params.Set(`t`, fmt.Sprint(c.Id))
	ApiGet(r)
	// 解析请求数据
	listInfo := model.MovieListInfo{}
	_ = json.Unmarshal(r.Resp, &listInfo)
	// 获取pageCount信息, 循环获取所有页数据
	pageCount := listInfo.PageCount
	// 使用chan + goroutine 进行并发获取
	chPg := make(chan int, int(pageCount))
	chClose := make(chan int)
	// 开始获取所有信息, 使用协程并发获取数据
	for i := 1; i <= int(pageCount); i++ {
		// 将当前分类的所有页码存入chPg
		chPg <- i
	}
	close(chPg)
	// 开启MAXGoroutine数量的协程进行请求
	for i := 0; i < config.MAXGoroutine; i++ {
		go func() {
			// 当前协程结束后向 chClose中写入一次数据
			defer func() { chClose <- 0 }()
			for {
				pg, ok := <-chPg
				if !ok {
					return
				}
				// 使用新的 请求参数
				req := RequestInfo{Uri: FILM_COLLECT_SITE, Params: url.Values{}}
				req.Params.Set(`ac`, "list")
				req.Params.Set(`t`, fmt.Sprint(c.Id))
				req.Params.Set(`pg`, fmt.Sprint(pg))
				// 保存当前分类下的影片信息
				info := model.MovieListInfo{}
				ApiGet(&req)
				// 如果返回数据中的list为空,则直接结束本分类的资源获取
				if len(r.Resp) <= 0 {
					log.Println("SaveMoves Error Response Is Empty")
					return
				}
				_ = json.Unmarshal(r.Resp, &info)
				if info.List == nil {
					log.Println("MovieList Is Empty")
					return
				}
				// 处理影片信息
				list := common.ProcessMovieListInfo(info.List)
				// 保存影片信息至redis
				_ = model.SaveMoves(list)
			}
		}()
	}
	// 使用chClose等待当前分类列表数据请求完毕
	for i := 0; i < config.MAXGoroutine; i++ {
		<-chClose
	}
}

// AllMovieInfo 拉取全部影片的基本信息
func AllMovieInfo() {
	keys := model.AllMovieInfoKey()
	for _, key := range keys {
		// 获取当前分类下的sort set数据集合
		movies := model.GetMovieListByKey(key)
		ids := ""
		for i, m := range movies {
			// 反序列化获取影片基本信息
			movie := model.Movie{}
			err := json.Unmarshal([]byte(m), &movie)
			if err == nil && movie.Id != 0 {
				// 拼接ids信息
				ids = fmt.Sprintf("%s,%d", ids, movie.Id)
			}
			// 每20个id执行一次请求, limit 最多20
			if (i+1)%20 == 0 {
				// ids对应影片的详情信息
				go MoviesDetails(strings.Trim(ids, ","))
				ids = ""
			}
		}
		// 如果ids != "" , 将剩余id执行一次请求
		MoviesDetails(strings.Trim(ids, ","))
	}
}

// MoviesDetails 获取影片详情信息, ids 影片id,id,....
func MoviesDetails(ids string) {
	// // 添加一个等待任务, 执行完减去一个任务
	//wg.Add(1)
	//defer wg.Done()
	// 如果ids为空数据则直接返回
	if len(ids) <= 0 {
		return
	}
	// 设置请求参数
	r := RequestInfo{
		Uri:    FILM_COLLECT_SITE,
		Params: url.Values{},
	}
	r.Params.Set("ac", "detail")
	r.Params.Set("ids", ids)
	ApiGet(&r)
	// 映射详情信息
	details := model.DetailListInfo{}
	// 如果返回数据为空则直接结束本次方法
	if len(r.Resp) <= 0 {
		return
	}
	// 序列化详情数据
	err := json.Unmarshal(r.Resp, &details)
	if err != nil {
		log.Println("DetailListInfo Unmarshal Error: ", err)
		return
	}
	// 处理details信息
	list := common.ProcessMovieDetailList(details.List)
	// 保存影片详情信息到redis
	err = model.SaveDetails(list)
	if err != nil {
		log.Println("SaveDetails Error: ", err)
	}
}

// SearchInfoToMdb 扫描redis中的检索信息, 并批量存入mysql
func SearchInfoToMdb() {
	// 1. 从redis的Zset集合中scan扫描数据, 每次100条
	var cursor uint64 = 0
	var count int64 = 100
	for {
		infoList, nextStar := model.ScanSearchInfo(cursor, count)
		// 2. 将扫描到的数据插入mysql中
		model.BatchSave(infoList)
		// 3.设置下次开始的游标
		cursor = nextStar
		// 4. 判断迭代是否已经结束 cursor为0则表示已经迭代完毕
		if cursor == 0 {
			return
		}
	}

}

// GetRecentMovie 获取最近更的影片, 默认最近3小时
func GetRecentMovie() {
	// 请求URL URI?ac=list&h=6
	r := RequestInfo{Uri: FILM_COLLECT_SITE, Params: url.Values{}}
	r.Params.Set("ac", "list")
	r.Params.Set("pg", "1")
	r.Params.Set("h", config.UpdateInterval)
	// 执行请求获取分页信息
	ApiGet(&r)
	if len(r.Resp) < 0 {
		log.Println("更新数据获取失败")
		return
	}
	pageInfo := model.MovieListInfo{}
	_ = json.Unmarshal(r.Resp, &pageInfo)
	// 获取分页数据
	ids := ""
	// 存储检索信息
	var tempSearchList []model.SearchInfo
	// 获取影片详细数据,并保存到redis中
	for i := 1; i <= int(pageInfo.PageCount); i++ {
		// 执行获取影片基本信息
		r.Params.Set("pg", fmt.Sprint(i))
		ApiGet(&r)
		// 解析请求的结果
		if len(r.Resp) < 0 {
			log.Println("更新数据获取失败")
			return
		}
		info := model.MovieListInfo{}
		_ = json.Unmarshal(r.Resp, &info)
		// 将影片信息保存到 movieList
		list := common.ProcessMovieListInfo(info.List)
		_ = model.SaveMoves(list)
		// 拼接ids 用于请求detail信息
		for _, m := range list {
			ids = fmt.Sprintf("%s,%d", ids, m.Id)
			// 保存一份id切片用于添加mysql检索信息
			tempSearchList = append(tempSearchList, model.SearchInfo{Mid: m.Id, Cid: m.Cid})
		}
		// 执行获取详情请求
		MoviesDetails(strings.Trim(ids, ","))
		ids = ""
	}
	// 根据idList 补全对应影片的searInfo信息
	var sl []model.SearchInfo
	for _, s := range tempSearchList {
		// 通过id 获取对应的详情信息
		sl = append(sl, model.ConvertSearchInfo(model.GetDetailByKey(fmt.Sprintf(config.MovieDetailKey, s.Cid, s.Mid))))
	}
	// 调用批量保存或更新方法, 如果对应mid数据存在则更新, 否则执行插入
	model.BatchSaveOrUpdate(sl)
}

// StartSpiderRe 清空存储数据,从零开始获取
func StartSpiderRe() {
	// 删除已有的存储数据, redis 和 mysql中的存储数据全部清空
	model.RemoveAll()
	// 执行完整数据获取
	StartSpider()
}
