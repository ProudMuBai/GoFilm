package spider

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"server/config"
	"server/model/collect"
	"server/model/system"
	"server/plugin/common/conver"
	"server/plugin/common/util"
	"strconv"
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
	// 通过 采集时长 h 的不同来执行不同出处理方式
	switch {
	case h < 0:
		// 采集时长为负数则先执行对应数据表的重置
		if s.Grade == system.MasterCollect {
			system.FilmZero()
		} else {
			// 如果所处站点是次级站点, 则删除对应站点在表中的数据
			system.ResetSlaveMovieInfoTable()
		}
	case h > 0:
		// 如果采集时长是正常数值, 则设置参数 h
		r.Params.Set("h", fmt.Sprint(h))
	default:
		log.Println("Params Collect time Exception !!!")
		return errors.New(" Params Collect time Exception !!! ")
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
		// 视频数据采集完成后 对暂存数据进行处理和优化
		if s.Grade == system.MasterCollect {
			// 如果采集时长为负, (全量采集), 则在数据采集完成后为search表添加索引
			if h < 0 {
				// 全量采集时进行数据同步(仅保存)
				system.SyncMovieDetail(s.Id, s.Grade)
				system.AddSearchIndex()
			}
			// 采集时长在一定阈值内时执行redis数据同步 (存在则更新, 不存在则新增)

			// 开启图片同步
			if s.SyncPictures {
				system.SyncFilmPicture()
			}
			// 每次成功执行完都清理redis中的相关API接口数据缓存
			ClearCache()
		} else if s.Grade == system.SlaveCollect {
			// 如果采集时长为负, (全量采集), 则在数据采集完成后为search表添加索引
			if h < 0 {
				// 全量采集时进行数据同步
				system.SyncMovieDetail(s.Id, s.Grade)
			}
		}

	case system.CollectArticle, system.CollectActor, system.CollectRole, system.CollectWebSite:
		log.Println("暂未开放此采集功能!!!")
		return errors.New("暂未开放此采集功能")
	}
	log.Println("Spider Task Exercise Success")
	return nil
}

func HandleCollectRefine(id string, h int) error {
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
	// 通过 采集时长 h 的不同来执行不同前置出处理方式
	switch {
	case h < 0:
		// 采集时长为负数则先执行对应数据表的重置
		if s.Grade == system.MasterCollect {
			system.FilmZero()
		} else {
			// 如果所处站点是次级站点, 则删除对应站点在表中的数据
			system.ResetSlaveMovieInfoTable()
		}
	case h > 0:
		// 如果采集时长是正常数值, 则设置参数 h
		r.Params.Set("h", fmt.Sprint(h))
	default:
		log.Println("Params Collect time Exception !!!")
		return errors.New(" Params Collect time Exception !!! ")
	}
	// 通过采集类型分别执行不同的采集方法
	switch s.CollectType {
	case system.CollectVideo:
		// 采集视频资源 根据采集站类型进行不同逻辑
		switch s.Grade {
		case system.MasterCollect:
			// 获取展示的分类切片信息
			cl := system.GetRevealCategoryList()
			for _, c := range cl {
				// 获取分类采集页数
				r.Params.Set("t", fmt.Sprint(c.Id))
				pageCount, err := spiderCore.GetPageCount(r)
				if err != nil {
					// 如果第二次获取分页页数依旧获取失败则关闭当前采集任务
					pageCount, err = spiderCore.GetPageCount(r)
					if err != nil {
						return err
					}
				}
				// 如果采集源参数中采集间隔参数大于500ms,则使用单线程采集
				if s.Interval > 500 {
					// 少量数据不开启协程
					for i := 1; i <= pageCount; i++ {
						// 设置采集参数pg
						r.Params.Set("pg", fmt.Sprint(i))
						collectFilmRefine(s, r)
						// 执行一次采集后休眠指定时长
						time.Sleep(time.Duration(s.Interval) * time.Millisecond)
					}
				} else if pageCount <= config.MAXGoroutine*5 {
					// 少量数据不开启协程
					for i := 1; i <= pageCount; i++ {
						r.Params.Set("pg", fmt.Sprint(i))
						collectFilmRefine(s, r)
					}
				} else {
					// 如果分页数量较大则开启协程
					collectFilmMT(pageCount, s, r, collectFilmRefine)
				}
			}
		case system.SlaveCollect:
			pageCount, err := spiderCore.GetPageCount(r)
			if err != nil {
				// 如果第二次获取分页页数依旧获取失败则关闭当前采集任务
				pageCount, err = spiderCore.GetPageCount(r)
				if err != nil {
					return err
				}
			}
			// 如果采集源参数中采集间隔参数大于500ms,则使用单线程采集
			if s.Interval > 500 {
				// 少量数据不开启协程
				for i := 1; i <= pageCount; i++ {
					// 设置采集参数pg
					r.Params.Set("pg", fmt.Sprint(i))
					collectFilmRefine(s, r)
					// 执行一次采集后休眠指定时长
					time.Sleep(time.Duration(s.Interval) * time.Millisecond)
				}
			} else if pageCount <= config.MAXGoroutine*5 {
				// 少量数据不开启协程
				for i := 1; i <= pageCount; i++ {
					r.Params.Set("pg", fmt.Sprint(i))
					collectFilmRefine(s, r)
				}
			} else {
				// 如果分页数量较大则开启协程
				collectFilmMT(pageCount, s, r, collectFilmRefine)
			}
		}

		// 视频数据采集完成后 对暂存数据进行处理和优化
		if s.Grade == system.MasterCollect {
			// 如果采集时长为负, (全量采集), 则在数据采集完成后为search表添加索引
			if h < 0 {
				// 全量采集时进行数据同步以及添加索引(仅保存)
				system.SyncMovieDetail(s.Id, s.Grade)
				system.AddSearchIndex()
				system.AddMovieDetailIndex()
			}
			// 采集时长在一定阈值内时执行redis数据同步 (存在则更新, 不存在则新增)

			// 开启图片同步
			if s.SyncPictures {
				system.SyncFilmPicture()
			}
			// 每次成功执行完都清理redis中的相关API接口数据缓存
			ClearCache()
		} else if s.Grade == system.SlaveCollect {
			// 如果采集时长为负, (全量采集), 则在数据采集完成后为search表添加索引
			if h < 0 {
				// 全量采集时进行数据同步
				system.SyncMovieDetail(s.Id, s.Grade)
				system.AddSlaveMovieInfoIndex()
			}
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

// collectFilm 影视详情采集 (单一源分页全采集)
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
		// 添加采集失败记录
		fr := system.FailureRecord{OriginId: s.Id, OriginName: s.Name, Uri: s.Uri, CollectType: system.CollectVideo, PageNumber: pg, Hour: h, Cause: fmt.Sprintln(err), Status: 1}
		system.SaveFailureRecord(fr)
		log.Println("GetMovieDetail Error: ", err)
		return
	}
	// 通过采集站 Grade 类型, 执行不同的存储逻辑
	switch s.Grade {
	case system.MasterCollect:
		// 将数据缓存到redis中
		if err = system.MovieDetailCache(list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
		//break
		// 如果 采集时长 h 小于阈值, 则将主体数据缓存到redis
		//if h > 0 && h < config.FilmSaveCacheThreshold {
		//	// 主站点 	执行保存或更新
		//	if err = system.BatchUpdateDetails(list); err != nil {
		//		log.Println("SaveDetails Error: ", err)
		//	}
		//} else {
		//	// 主站点 	从零开始只执行保存逻辑
		//	if err = system.MovieDetailCache(list); err != nil {
		//		log.Println("SaveDetails Error: ", err)
		//	}
		//}
		// 如果主站点开启了图片同步, 则将图片url以及对应的mid存入ZSet集合中
		if s.SyncPictures {
			if err = system.SaveVirtualPic(conver.ConvertVirtualPicture(list)); err != nil {
				log.Println("SaveVirtualPic Error: ", err)
			}
		}
	case system.SlaveCollect:
		// 将采集数据缓存到redis中
		if err = system.SlaveDetailCache(s.Id, list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
		//if h > 0 && h < config.FilmSaveCacheThreshold {
		//	// 附属站点	仅保存影片播放信息到mysql
		//	if err = system.UpdateSitePlayList(s.Id, list); err != nil {
		//		log.Println("SaveDetails Error: ", err)
		//	}
		//} else {
		//	// 附属站点	仅保存影片播放信息到mysql
		//	if err = system.SlaveDetailCache(s.Id, list); err != nil {
		//		log.Println("SaveDetails Error: ", err)
		//	}
		//}

	}
}

// collectFilmById 采集指定ID的影片信息
func collectFilmById(ids string, s *system.FilmSource) {
	// 生成请求参数
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	// 设置分页页数
	r.Params.Set("pg", "1")
	// 设置影片IDS参数信息
	r.Params.Set("ids", ids)
	// 执行采集方法 获取影片详情list
	list, err := spiderCore.GetFilmDetail(r)
	if err != nil || len(list) <= 0 {
		log.Println("GetMovieDetail Error: ", err)
		return
	}
	// 通过采集站 Grade 类型, 执行不同的存储逻辑
	switch s.Grade {
	case system.MasterCollect:
		// 主站点 	保存完整影片详情信息到 redis 和 mysql 中
		if err = system.SaveDetail(list[0]); err != nil {
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
		if err = system.UpdateSitePlayList(s.Id, list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
	}
}

// 影片信息采集, 改进版
func collectFilmRefine(s *system.FilmSource, r util.RequestInfo) {
	// 执行采集方法 获取影片详情list
	//log.Printf("%s?%s", r.Uri, r.Params.Encode())
	list, err := spiderCore.GetFilmDetail(r)
	if err != nil || len(list) <= 0 {
		// 添加采集失败记录
		pg, _ := strconv.Atoi(r.Params.Get("pg"))
		h, _ := strconv.Atoi(r.Params.Get("h"))
		fr := system.FailureRecord{OriginId: s.Id, OriginName: s.Name, Uri: s.Uri, CollectType: system.CollectVideo, PageNumber: pg, Hour: h, Cause: fmt.Sprintln(err), Status: 1}
		system.SaveFailureRecord(fr)
		log.Println("GetMovieDetail Error: ", err)
		return
	}
	// 通过采集站 Grade 类型, 执行不同的存储逻辑
	switch s.Grade {
	case system.MasterCollect:
		// 将数据缓存到redis中
		if err = system.MovieDetailCache(list); err != nil {
			log.Println("SaveDetails Error: ", err)
		}
		// 如果主站点开启了图片同步, 则将图片url以及对应的mid存入ZSet集合中
		if s.SyncPictures {
			if err = system.SaveVirtualPic(conver.ConvertVirtualPicture(list)); err != nil {
				log.Println("SaveVirtualPic Error: ", err)
			}
		}
	case system.SlaveCollect:
		// 将采集数据缓存到redis中
		if err = system.SlaveDetailCache(s.Id, list); err != nil {
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
	for i := 0; i < GoroutineNum; i++ {
		<-waitCh
	}
}

// collectFilmMT 并发采集影片信息
func collectFilmMT(capacity int, s *system.FilmSource, r util.RequestInfo, collectFunc func(s *system.FilmSource, r util.RequestInfo)) {
	// 初始化 channel, 容量为 capacity
	ch := make(chan int, capacity)

	// 收集结束标识
	waitCh := make(chan int)
	// 循环将所有需采集的页码写入 ch
	for i := 1; i <= capacity; i++ {
		ch <- i
	}
	close(ch)
	// 开启 MAXGoroutine 数量的协程, 如果分页页数小于设定的最大线程数, 则将线程数设置为1
	var GoroutineNum = config.MAXGoroutine
	if capacity < GoroutineNum*5 {
		GoroutineNum = 1
	}
	// 如果满足开启并发的条件, 则开启GoroutineNum数量的协程进行并发采集
	for i := 0; i < GoroutineNum; i++ {
		go func() {
			defer func() { waitCh <- 0 }()
			for {
				// 从channel中获取 pageNumber
				pg, ok := <-ch
				if !ok {
					break
				}
				// 执行对应的采集方法, 并发时不同使用同一个requestInfo
				requestInfo := util.CopyRequestInfo(r)
				requestInfo.Params.Set("pg", fmt.Sprint(pg))
				collectFunc(s, requestInfo)
			}
		}()
	}
	// 等待所有协程执行完毕
	for i := 0; i < GoroutineNum; i++ {
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
				err := HandleCollectRefine(fs.Id, h)
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
			if err := HandleCollectRefine(s.Id, h); err != nil {
				log.Println(err)
			}
		}
	}
}

// ClearSpider  删除所有已采集的影片信息
func ClearSpider() {
	system.FilmZero()
}

// StarZero 情况站点内所有影片信息
func StarZero(h int) {
	// 首先清除影视信息
	system.FilmZero()
	// 开启自动采集
	AutoCollect(h)
}

// CollectSingleFilm 通过影片唯一ID获取影片信息
func CollectSingleFilm(ids string) {
	// 获取采集站列表信息
	fl := system.GetCollectSourceList()
	// 循环遍历所有采集站信息
	for _, f := range fl {
		// 目前仅对主站点进行处理
		if f.Grade == system.MasterCollect && f.State {
			collectFilmById(ids, &f)
			return
		}
	}
}

// ======================================================= 采集拓展内容  =======================================================

// SingleRecoverSpider 二次采集
func SingleRecoverSpider(fr *system.FailureRecord) {
	// 通过采集时长范围执行不同的采集方式
	switch {
	case fr.Hour > 168 && fr.Hour < 360:
		// 将此记录之后的所有同类采集记录变更为已重试
		system.ChangeRecord(fr, 0)
		// 如果采集的内容是 7~15 天之内更新的内容,则采集此记录之后的所有更新内容
		// 获取采集参数h, 采集时长变更为 原采集时长 + 采集记录距现在的时长
		h := fr.Hour + int(math.Ceil(time.Since(fr.CreatedAt).Hours()))
		// 对当前所有已启用的站点 更新最新 h 小时的内容
		AutoCollect(h)
	case fr.Hour < 0, fr.Hour > 4320:
		// 将此记录状态修改为已重试
		system.ChangeRecord(fr, 0)
		// 如果采集的是 最近180天内更新的内容 或全部内容, 则只对当前一条记录进行二次采集
		s := system.FindCollectSourceById(fr.OriginId)
		collectFilm(s, fr.Hour, fr.PageNumber)
	default:
		// 其余范围,暂不处理
		break
	}
}

// FullRecoverSpider 扫描记录表中的失败记录, 并进行处理 (用于定时任务定期处理失败采集)
func FullRecoverSpider() {
	/*
		获取待处理的记录数据
		1. 采集时长 > 168h (一周,7天)  状态-1 待处理, | 只获取满足条件的最早的待处理记录
		2. 采集时长 > 4320h (半年,180天)  状态-1 待处理,   | 获取满足条件的所有数据
	*/
	list := system.PendingRecord()

	// 遍历记录信息切片, 针对不同时长进行不同处理
	for _, fr := range list {
		switch {
		case fr.Hour > 0 && fr.Hour < 4320:
			// 将此记录之后的所有同类采集记录变更为已重试
			system.ChangeRecord(&fr, 0)
			// 如果采集的内容是 0~180 天之内更新的内容,则采集此记录之后的所有更新内容
			// 获取采集参数h, 采集时长变更为 原采集时长 + 采集记录距现在的时长
			h := fr.Hour + int(math.Ceil(time.Since(fr.CreatedAt).Hours()))
			// 对当前所有已启用的站点 更新最新 h 小时的内容
			AutoCollect(h)
		case fr.Hour < 0, fr.Hour > 4320:
			// 将此记录状态修改为已重试
			system.ChangeRecord(&fr, 0)
			// 如果采集的是 180天之前更新的内容 或全部内容, 则只对当前一条记录进行二次采集
			s := system.FindCollectSourceById(fr.OriginId)
			collectFilm(s, fr.Hour, fr.PageNumber)
		default:
			// 其余范围,暂不处理
			break
		}
	}

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
