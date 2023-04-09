package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"server/config"
	"server/model"
	"server/plugin/common"
	"time"
)

/*
	舍弃第一版的数据处理思路, v2版本
	直接分页获取采集站点的影片详情信息


*/

/*
 1. 选择一个采集主站点, mysql检索表中只存储主站点检索的信息
 2. 采集多个站点数据
    2.1 主站点的采集数据完整地保存相关信息, basicInfo movieDetail search 等信息
    2.2 其余站点数据只存储 name(影片名称), playUrl(播放url), 存储形式 Key<hash(name)>:value([]MovieUrlInfo)
 3. api数据格式不变, 获取影片详情时通过subTitle 去redis匹配其他站点的对应播放源并整合到主站详情信息的playUrl中
 4. 影片搜索时不再使用name进行匹配, 改为使用 subTitle 进行匹配
*/
const (
	MainSite = "https://www.feisuzyapi.com/api.php/provide/vod/"
)

type Site struct {
	Name string
	Uri  string
}

// SiteList 播放源采集站
var SiteList = []Site{
	//{"tk", "https://api.tiankongapi.com/api.php/provide/vod"},
	//{"yh", "https://m3u8.apiyhzy.com/api.php/provide/vod/"},
	{"su", "https://subocaiji.com/api.php/provide/vod/at/json"},
	{"lz", "https://cj.lziapi.com/api.php/provide/vod/"},
	{"ff", "https://cj.ffzyapi.com/api.php/provide/vod/"},
	//{"fs", "https://www.feisuzyapi.com/api.php/provide/vod/"},
}

// StartSpider 执行多源spider
func StartSpider() {
	// 保存分类树
	CategoryList()
	log.Println("CategoryList 影片分类信息保存完毕")
	// 爬取主站点数据
	MainSiteSpider()
	log.Println("MainSiteSpider 主站点影片信息保存完毕")
	// 查找并创建search数据库
	time.Sleep(time.Second * 10)
	model.CreateSearchTable()
	SearchInfoToMdb()
	log.Println("SearchInfoToMdb 影片检索信息保存完毕")
	// 获取其他站点数据
	go MtSiteSpider()
	log.Println("Spider End , 数据保存执行完成")
	//time.Sleep(time.Second * 10)
}

// CategoryList 获取分类数据
func CategoryList() {
	// 设置请求参数信息
	r := RequestInfo{Uri: MainSite, Params: url.Values{}}
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

// MainSiteSpider 主站点数据处理
func MainSiteSpider() {
	// 获取分页页数
	pageCount, err := GetPageCount(RequestInfo{Uri: MainSite, Params: url.Values{}})
	// 主站点分页出错直接终止程序
	if err != nil {
		panic(err)
	}
	// 开启协程加快分页请求速度
	ch := make(chan int, pageCount)
	waitCh := make(chan int)
	for i := 1; i <= pageCount; i++ {
		ch <- i
	}
	close(ch)
	for i := 0; i < config.MAXGoroutine; i++ {
		go func() {
			defer func() { waitCh <- 0 }()
			for {
				pg, ok := <-ch
				if !ok {
					break
				}
				list, e := GetMovieDetail(pg, RequestInfo{Uri: MainSite, Params: url.Values{}})
				if e != nil {
					log.Println("GetMovieDetail Error: ", err)
					continue
				}
				// 保存影片详情信息到redis
				if err = model.SaveDetails(list); err != nil {
					log.Println("SaveDetails Error: ", err)
				}
			}
		}()
	}
	for i := 0; i < config.MAXGoroutine; i++ {
		<-waitCh
	}
}

// MtSiteSpider 附属数据源处理
func MtSiteSpider() {
	for _, s := range SiteList {
		// 执行每个站点的播放url缓存
		PlayDetailSpider(s)
		log.Println(s.Name, "playUrl 爬取完毕!!!")
	}
}

// PlayDetailSpider SpiderSimpleInfo 获取单个站点的播放源
func PlayDetailSpider(s Site) {
	// 获取分页页数
	pageCount, err := GetPageCount(RequestInfo{Uri: s.Uri, Params: url.Values{}})
	// 出错直接终止当前站点数据获取
	if err != nil {
		log.Println(err)
		return
	}

	// 开启协程加快分页请求速度
	ch := make(chan int, pageCount)
	waitCh := make(chan int)
	for i := 1; i <= pageCount; i++ {
		ch <- i
	}
	close(ch)
	for i := 0; i < config.MAXGoroutine; i++ {
		go func() {
			defer func() { waitCh <- 0 }()
			for {
				pg, ok := <-ch
				if !ok {
					break
				}
				list, e := GetMovieDetail(pg, RequestInfo{Uri: s.Uri, Params: url.Values{}})
				if e != nil || len(list) <= 0 {
					log.Println("GetMovieDetail Error: ", err)
					continue
				}
				// 保存影片播放信息到redis
				if err = model.SaveSitePlayList(s.Name, list); err != nil {
					log.Println("SaveDetails Error: ", err)
				}
			}
		}()
	}
	for i := 0; i < config.MAXGoroutine; i++ {
		<-waitCh
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

// UpdateMovieDetail 定时更新主站点和其余播放源信息
func UpdateMovieDetail() {
	// 更新主站系列信息
	UpdateMainDetail()
	// 更新播放源数据信息
	UpdatePlayDetail()
}

// UpdateMainDetail 更新主站点的最新影片
func UpdateMainDetail() {
	// 获取分页页数
	r := RequestInfo{Uri: MainSite, Params: url.Values{}}
	r.Params.Set("h", config.UpdateInterval)
	pageCount, err := GetPageCount(r)
	if err != nil {
		log.Printf("Update MianStieDetail failed")
	}
	// 保存本次更新的所有详情信息
	var ds []model.MovieDetail
	// 获取分页数据
	for i := 1; i <= pageCount; i++ {
		list, err := GetMovieDetail(i, r)
		if err != nil {
			continue
		}
		// 保存更新的影片信息, 同类型直接覆盖
		if err = model.SaveDetails(list); err != nil {
			log.Printf("Update MianStieDetail failed, SaveDetails Error ")
		}
		ds = append(ds, list...)
	}

	// 整合详情信息切片
	var sl []model.SearchInfo
	for _, d := range ds {
		// 通过id 获取对应的详情信息
		sl = append(sl, model.ConvertSearchInfo(d))
	}
	// 调用批量保存或更新方法, 如果对应mid数据存在则更新, 否则执行插入
	model.BatchSaveOrUpdate(sl)
}

// UpdatePlayDetail 更新最x小时的影片播放源数据
func UpdatePlayDetail() {
	for _, s := range SiteList {
		// 获取单个站点的分页数
		r := RequestInfo{Uri: MainSite, Params: url.Values{}}
		r.Params.Set("h", config.UpdateInterval)
		pageCount, err := GetPageCount(r)
		if err != nil {
			log.Printf("Update %s playDetail failed", s.Name)
		}
		for i := 1; i <= pageCount; i++ {
			// 获取详情信息, 保存到对应hashKey中
			list, e := GetMovieDetail(i, r)
			if e != nil || len(list) <= 0 {
				log.Println("GetMovieDetail Error: ", err)
				continue
			}
			// 保存影片播放信息到redis
			if err = model.SaveSitePlayList(s.Name, list); err != nil {
				log.Println("SaveDetails Error: ", err)
			}
		}
	}
}

// StartSpiderRe 清空存储数据,从零开始获取
func StartSpiderRe() {
	// 删除已有的存储数据, redis 和 mysql中的存储数据全部清空
	model.RemoveAll()
	// 执行完整数据获取
	StartSpider()
}

// =========================公共方法==============================

// GetPageCount 获取总页数
func GetPageCount(r RequestInfo) (count int, err error) {
	// 发送请求获取pageCount
	r.Params.Set("ac", "detail")
	r.Params.Set("pg", "1")
	ApiGet(&r)
	//  判断请求结果是否为空, 如果为空直接输出错误并终止
	if len(r.Resp) <= 0 {
		err = errors.New("response is empty")
		return
	}
	// 获取pageCount
	res := model.DetailListInfo{}
	err = json.Unmarshal(r.Resp, &res)
	if err != nil {
		return
	}
	count = int(res.PageCount)
	return
}

// GetMovieDetail 处理详情接口请求返回的数据
func GetMovieDetail(pageNumber int, r RequestInfo) (list []model.MovieDetail, err error) {
	// 防止json解析异常引发panic
	defer func() {
		if e := recover(); e != nil {
			log.Println("GetMovieDetail Failed : ", e)
		}
	}()
	// 设置分页请求参数
	r.Params.Set(`ac`, `detail`)
	r.Params.Set(`pg`, fmt.Sprint(pageNumber))
	ApiGet(&r)
	// 影视详情信息
	details := model.DetailListInfo{}
	// 如果返回数据为空则直接结束本次循环
	if len(r.Resp) <= 0 {
		err = errors.New("response is empty")
		return
	}
	// 序列化详情数据
	if err = json.Unmarshal(r.Resp, &details); err != nil {
		return
	}
	// 处理details信息
	list = common.ProcessMovieDetailList(details.List)
	return
}
