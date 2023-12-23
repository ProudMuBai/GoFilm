package spider

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/url"
	"server/config"
	"server/model/collect"
	"server/model/system"
	"server/plugin/common/util"
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

// StartSpider 执行多源spider
func StartSpider() {
	// 保存分类树
	CategoryList()
	log.Println("CategoryList 影片分类信息保存完毕")
	// 爬取主站点数据
	MainSiteSpider()
	log.Println("MainSiteSpider 主站点影片信息保存完毕")
	// 查找并创建search数据库, 保存search信息, 添加索引
	time.Sleep(time.Second * 10)
	system.SyncSearchInfo(0)
	system.AddSearchIndex()
	log.Println("SearchInfoToMdb 影片检索信息保存完毕")
	//获取其他站点数据
	scl := system.GetCollectSourceListByGrade(system.SlaveCollect)
	go MtSiteSpider(scl...)
	log.Println("Spider End , 数据保存执行完成")
	time.Sleep(time.Second * 10)
}

// CategoryList 获取分类数据
func CategoryList() {
	// 获取主站点uri
	mc := system.GetCollectSourceListByGrade(system.MasterCollect)[0]
	// 获取分类树形数据
	categoryTree, err := spiderCore.GetCategoryTree(util.RequestInfo{Uri: mc.Uri, Params: url.Values{}})
	if err != nil {
		log.Println("GetCategoryTree Error: ", err)
		return
	}
	// 保存 tree 到redis
	err = system.SaveCategoryTree(categoryTree)
	if err != nil {
		log.Println("SaveCategoryTree Error: ", err)
	}
}

// MainSiteSpider 主站点数据处理
func MainSiteSpider() {
	// 获取主站点uri
	mc := system.GetCollectSourceListByGrade(system.MasterCollect)[0]
	// 获取分页页数
	pageCount, err := spiderCore.GetPageCount(util.RequestInfo{Uri: mc.Uri, Params: url.Values{}})
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
				list, e := spiderCore.GetDetail(pg, util.RequestInfo{Uri: mc.Uri, Params: url.Values{}})
				if e != nil {
					log.Println("GetMovieDetail Error: ", err)
					continue
				}
				// 保存影片详情信息到redis
				if err = system.SaveDetails(list); err != nil {
					log.Println("SaveDetails Error: ", err)
				}
			}
		}()
	}
	for i := 0; i < config.MAXGoroutine; i++ {
		<-waitCh
	}
}

// MtSiteSpider 附属站点数据源处理
func MtSiteSpider(scl ...system.FilmSource) {
	for _, s := range scl {
		// 执行每个站点的播放url缓存
		PlayDetailSpider(s)
		log.Println(s.Name, "playUrl 爬取完毕!!!")
	}
}

// PlayDetailSpider SpiderSimpleInfo 获取单个站点的播放源
func PlayDetailSpider(s system.FilmSource) {
	// 获取分页页数
	pageCount, err := spiderCore.GetPageCount(util.RequestInfo{Uri: s.Uri, Params: url.Values{}})
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
				list, e := spiderCore.GetDetail(pg, util.RequestInfo{Uri: s.Uri, Params: url.Values{}})
				if e != nil || len(list) <= 0 {
					log.Println("GetMovieDetail Error: ", err)
					continue
				}
				// 保存影片播放信息到redis
				if err = system.SaveSitePlayList(s.Name, list); err != nil {
					log.Println("SaveDetails Error: ", err)
				}
			}
		}()
	}
	for i := 0; i < config.MAXGoroutine; i++ {
		<-waitCh
	}
}

// UpdateMovieDetail 定时更新主站点和其余播放源信息
func UpdateMovieDetail() {
	// 更新主站系列信息
	UpdateMainDetail()
	// 更新附属播放源数据信息
	scl := system.GetCollectSourceListByGrade(system.SlaveCollect)
	UpdatePlayDetail(scl...)
}

// UpdateMainDetail 更新主站点的最新影片
func UpdateMainDetail() {
	// 获取主站点uri
	l := system.GetCollectSourceListByGrade(system.MasterCollect)
	mc := system.FilmSource{}
	for _, v := range l {
		if len(v.Uri) > 0 {
			mc = v
			break
		}
	}
	// 获取分页页数
	r := util.RequestInfo{Uri: mc.Uri, Params: url.Values{}}
	r.Params.Set("h", config.UpdateInterval)
	pageCount, err := spiderCore.GetPageCount(r)
	if err != nil {
		log.Printf("Update MianStieDetail failed\n")
	}
	// 保存本次更新的所有详情信息
	var ds []system.MovieDetail
	// 获取分页数据
	for i := 1; i <= pageCount; i++ {
		list, err := spiderCore.GetDetail(i, r)
		if err != nil {
			continue
		}
		// 保存更新的影片信息, 同类型直接覆盖
		if err = system.SaveDetails(list); err != nil {
			log.Println("Update MainSiteDetail failed, SaveDetails Error ")
		}
		ds = append(ds, list...)
	}

	// 整合详情信息切片
	var sl []system.SearchInfo
	for _, d := range ds {
		// 通过id 获取对应的详情信息
		sl = append(sl, system.ConvertSearchInfo(d))
	}
	// 调用批量保存或更新方法, 如果对应mid数据存在则更新, 否则执行插入
	system.BatchSaveOrUpdate(sl)
}

// UpdatePlayDetail 更新最x小时的影片播放源数据
func UpdatePlayDetail(scl ...system.FilmSource) {
	for _, s := range scl {
		// 获取单个站点的分页数
		r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
		r.Params.Set("h", config.UpdateInterval)
		pageCount, err := spiderCore.GetPageCount(r)
		if err != nil {
			log.Printf("Update %s playDetail failed\n", s.Name)
		}
		for i := 1; i <= pageCount; i++ {
			// 获取详情信息, 保存到对应hashKey中
			list, e := spiderCore.GetDetail(i, r)
			if e != nil || len(list) <= 0 {
				log.Println("GetMovieDetail Error: ", err)
				continue
			}
			// 保存影片播放信息到redis
			if err = system.SaveSitePlayList(s.Name, list); err != nil {
				log.Println("SaveDetails Error: ", err)
			}
		}
	}
}

// StartSpiderRe 清空存储数据,从零开始获取
func StartSpiderRe() {
	// 删除已有的存储数据, redis 和 mysql中的存储数据全部清空
	system.FilmZero()
	// 执行完整数据获取
	StartSpider()
}

// =========================公共方法==============================

// CollectApiTest 测试采集接口是否可用
func CollectApiTest(s system.FilmSource) error {
	// 使用当前采集站接口采集一页数据
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	r.Params.Set("ac", s.CollectType.GetActionType())
	r.Params.Set("pg", "3")
	err := util.ApiTest(&r)
	// 首先核对接口返回值类型
	if err == nil {
		// 如果返回值类型为Json则执行Json序列化
		if s.ResultModel == system.JsonResult {
			var dp = collect.FilmDetailLPage{}
			if err = json.Unmarshal(r.Resp, &dp); err != nil {
				return errors.New(fmt.Sprint("测试失败, 返回数据异常, JSON序列化失败: ", err))
			}
			return nil
		} else if s.ResultModel == system.XmlResult {
			// 如果返回值类型为XML则执行XML序列化
			var rd = collect.RssD{}
			if err = xml.Unmarshal(r.Resp, &rd); err != nil {
				return errors.New(fmt.Sprint("测试失败, 返回数据异常, XML序列化失败", err))
			}
			return nil
		}
		return errors.New("测试失败, 接口返回值类型不符合规范")
	}
	return errors.New(fmt.Sprint("测试失败, 请求响应异常 : ", err.Error()))
}
