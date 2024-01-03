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
	"server/plugin/common/conver"
	"server/plugin/common/util"
	"time"
)

/*
	采集逻辑 v3

*/

var spiderCore = &JsonCollect{}

// ======================================================= 通用采集方法  =======================================================

// HandleCollect 影视采集  id-采集站ID h-时长/h
func HandleCollect(id string, h int) error {
	// 1. 首先通过ID获取对应采集站信息
	s := system.FindCollectSourceById(id)
	if s == nil {
		log.Println("Cannot Find Collect Source Site")
		return errors.New(" Cannot Find Collect Source Site ")
	} else if !s.State {
		log.Println(" The acquisition site was disabled ")
		return errors.New(" The acquisition site was disabled ")
	}
	// 如果是主站点且状态为启用则先获取分类tree信息
	if s.Grade == system.MasterCollect && s.State {
		// 是否存在分类树信息, 不存在则获取
		if !system.ExistsCategoryTree() {
			CollectCategory(s)
		}
	}

	// 生成 RequestInfo
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	// 如果 h == 0 则直接返回错误信息
	if h == 0 {
		log.Println(" Collect time cannot be zero ")
		return errors.New(" Collect time cannot be zer ")
	}
	// 如果 h = -1 则进行全量采集
	if h > 0 {
		r.Params.Set("h", fmt.Sprint(h))
	}
	// 2. 首先获取分页采集的页数
	pageCount, err := spiderCore.GetPageCount(r)
	// 分页页数失败 则再进行一次尝试
	if err != nil {
		// 如果第二次获取分页页数依旧获取失败则关闭当前采集任务
		pageCount, err = spiderCore.GetPageCount(r)
		if err != nil {
			return err
		}
	}
	// 通过采集类型分别执行不同的采集方法
	switch s.CollectType {
	case system.CollectVideo:
		// 采集视频资源
		// 如果采集源参数中采集间隔参数大于500ms,则使用单线程采集
		if s.Interval > 500 {
			// 少量数据不开启协程
			for i := 1; i <= pageCount; i++ {
				collectFilm(s, h, i)
				// 执行一次采集后休眠指定时长
				time.Sleep(time.Duration(s.Interval) * time.Millisecond)
			}
		} else if pageCount <= config.MAXGoroutine*2 {
			// 少量数据不开启协程
			for i := 1; i <= pageCount; i++ {
				collectFilm(s, h, i)
			}
		} else {
			// 如果分页数量较大则开启协程
			ConcurrentPageSpider(pageCount, s, h, collectFilm)
		}
		// 视频数据采集完成后同步相关信息到mysql
		if s.Grade == system.MasterCollect {
			// 执行影片信息更新操作
			if h > 0 {
				// 执行数据更新操作
				system.SyncSearchInfo(1)
			} else {
				// 清空searchInfo中的数据并重新添加, 否则执行
				system.SyncSearchInfo(0)
			}
			// 开启图片同步
			if s.SyncPictures {
				system.SyncFilmPicture()
			}
			// 每次成功执行完都清理redis中的相关API接口数据缓存
			clearCache()
		}

	case system.CollectArticle, system.CollectActor, system.CollectRole, system.CollectWebSite:
		log.Println("暂未开放此采集功能!!!")
		return errors.New("暂未开放此采集功能")
	}
	log.Println("Spider Task Exercise Success")
	return nil
}

// CollectCategory 影视分类采集
func CollectCategory(s *system.FilmSource) {
	// 获取分类树形数据
	categoryTree, err := spiderCore.GetCategoryTree(util.RequestInfo{Uri: s.Uri, Params: url.Values{}})
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

// 影视详情采集
func collectFilm(s *system.FilmSource, h, pg int) {
	// 生成请求参数
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	// 设置分页页数
	r.Params.Set("pg", fmt.Sprint(pg))
	// 如果 h = -1 则进行全量采集
	if h > 0 {
		r.Params.Set("h", fmt.Sprint(h))
	}
	// 执行采集方法 获取影片详情list
	list, err := spiderCore.GetFilmDetail(r)
	if err != nil || len(list) <= 0 {
		log.Println("GetMovieDetail Error: ", err)
		return
	}
	// 通过采集站 Grade 类型, 执行不同的存储逻辑
	switch s.Grade {
	case system.MasterCollect:
		// 主站点 	保存完整影片详情信息到 redis
		if err = system.SaveDetails(list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
		// 如果主站点开启了图片同步, 则将图片url以及对应的mid存入ZSet集合中
		if s.SyncPictures {
			if err = system.SaveVirtualPic(conver.ConvertVirtualPicture(list)); err != nil {
				log.Println("SaveVirtualPic Error: ", err)
			}
		}
	case system.SlaveCollect:
		// 附属站点	仅保存影片播放信息到redis
		if err = system.SaveSitePlayList(s.Id, list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
	}
}

// ConcurrentPageSpider 并发分页采集, 不限类型
func ConcurrentPageSpider(capacity int, s *system.FilmSource, h int, collectFunc func(s *system.FilmSource, hour, pageNumber int)) {
	// 开启协程并发执行
	ch := make(chan int, capacity)
	waitCh := make(chan int)
	for i := 1; i <= capacity; i++ {
		ch <- i
	}
	close(ch)
	// 开启 MAXGoroutine 数量的协程, 如果分页页数小于协程数则将协程数限制为分页页数
	var GoroutineNum = config.MAXGoroutine
	if capacity < GoroutineNum {
		GoroutineNum = capacity
	}
	for i := 0; i < GoroutineNum; i++ {
		go func() {
			defer func() { waitCh <- 0 }()
			for {
				// 从channel中获取 pageNumber
				pg, ok := <-ch
				if !ok {
					break
				}
				// 执行对应的采集方法
				collectFunc(s, h, pg)
			}
		}()
	}
	for i := 0; i < config.MAXGoroutine; i++ {
		<-waitCh
	}
}

// BatchCollect 批量采集, 采集指定的所有站点最近x小时内更新的数据
func BatchCollect(h int, ids ...string) {
	for _, id := range ids {
		// 如果查询到对应Id的资源站信息, 且资源站处于启用状态
		if fs := system.FindCollectSourceById(id); fs != nil && fs.State {
			// 采用协程并发执行, 每个站点单独开启一个协程执行
			go func() {
				err := HandleCollect(fs.Id, h)
				if err != nil {
					log.Println(err)
				}
			}()
			// 执行当前站点的采集任务
			//if err := HandleCollect(fs.Id, h); err != nil {
			//	log.Println(err)
			//}
		}
	}
}

// AutoCollect 自动进行对所有已启用站点的采集任务
func AutoCollect(h int) {
	// 获取采集站中所有站点, 进行遍历
	for _, s := range system.GetCollectSourceList() {
		// 如果当前站点为启用状态 则执行 HandleCollect 进行数据采集
		if s.State {
			if err := HandleCollect(s.Id, h); err != nil {
				log.Println(err)
			}
		}
	}
}

// StarZero 情况站点内所有影片信息
func StarZero(h int) {
	// 首先清除影视信息
	system.FilmZero()
	// 开启自动采集
	AutoCollect(h)
}

// ======================================================= 公共方法  =======================================================

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
