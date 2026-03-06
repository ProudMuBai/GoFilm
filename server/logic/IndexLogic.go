package logic

import (
	"server/config"
	"server/model/system"
	"server/plugin/common/util"
	"server/plugin/spider"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
*
 IndexController数据处理
*/

type IndexLogic struct {
}

var IL *IndexLogic

// IndexPage 首页数据处理
func (i *IndexLogic) IndexPage() map[string]interface{} {
	// 首页请求时长较高, 采用redis进行缓存, 在定时任务更新影片时清除对应缓存
	// 判断是否存在缓存数据, 存在则直接将数据返回
	Info := system.GetCacheData(config.IndexCacheKey)
	if Info != nil {
		return Info
	}
	Info = make(map[string]interface{})
	// 1. 首页分类数据处理 导航分类数据处理, 只提供 电影 电视剧 综艺 动漫 四大顶级分类和其子分类
	tree := system.CategoryTree{Category: &system.Category{Id: 0, Name: "分类信息"}}
	sysTree := system.GetCategoryTree()
	// 只展示show=true的分页影片信息
	for _, c := range sysTree.Children {
		// 只针对一级分类进行处理
		if c.Show {
			tree.Children = append(tree.Children, c)
		}
	}
	// 返回分类信息
	Info["category"] = tree
	// 2. 提供用于首页展示的顶级分类影片信息, 每分类 14条数据
	var list []map[string]interface{}
	for _, c := range tree.Children {
		// 生成分页参数
		page := system.Page{PageSize: 14, Current: 1}
		// 获取最近上映影片和本月热门影片
		var movies []system.MovieBasicInfo
		var hotMovies []system.SearchInfo
		if c.Children != nil {
			// 如果有子分类, 则通过Pid获取对应影片
			// 获取当前分类的最新上映影片
			movies = system.GetMovieListByPid(c.Id, &page)
			// 获取当前分类的本月热门影片
			hotMovies = system.GetHotMovieByPid(c.Id, &page)
		} else {
			// 如果当前分类为一级分类且没有子分类,则通过Cid获取对应数据
			// 获取当前分类的最新上映影片
			movies = system.GetMovieListByCid(c.Id, &page)
			// 获取当前分类的本月热门影片
			hotMovies = system.GetHotMovieByCid(c.Id, &page)
		}

		item := map[string]interface{}{"nav": c, "movies": movies, "hot": hotMovies}
		list = append(list, item)
	}
	Info["content"] = list
	// 3. 获取首页轮播数据
	Info["banners"] = system.GetBanners()
	// 不存在首页数据缓存时将查询数据缓存到redis中
	system.DataCache(config.IndexCacheKey, Info)
	return Info
}

// ClearIndexCache 删除首页数据缓存
func (i *IndexLogic) ClearIndexCache() {
	// 更新成功后删除首页缓存
	spider.ClearCache()
}

// GetFilmDetail 影片详情信息页面处理
func (i *IndexLogic) GetFilmDetail(id int64) system.MovieDetailVo {
	// 通过mid获取影片的详情信息
	detail := system.GetDetailByMid(id)
	//查找其他站点是否存在影片对应的播放源
	ml := multipleSource(&detail)
	// 转换组合主次站点信息
	return system.ConvertMovieDetailVo(detail, ml)
}

// GetCategoryInfo 分类信息获取, 组装导航栏需要的信息
func (i *IndexLogic) GetCategoryInfo() gin.H {
	// 组装nav导航所需的信息
	nav := gin.H{}
	// 1.获取所有分类信息
	tree := system.GetCategoryTree()
	// 2. 过滤出主页四大分类的tree信息
	for _, t := range tree.Children {
		switch t.Category.Name {
		case "动漫", "动漫片":
			nav["cartoon"] = t
		case "电影", "电影片":
			nav["film"] = t
		case "连续剧", "电视剧":
			nav["tv"] = t
		case "综艺", "综艺片":
			nav["variety"] = t
		}
	}
	// 获取所有的分类
	return nav
}

// GetNavCategory 获取导航分类信息
func (i *IndexLogic) GetNavCategory() []*system.Category {
	// 1.获取所有分类信息
	tree := system.GetCategoryTree()
	// 遍历一级分类返回可展示的分类数据
	var cl []*system.Category
	for _, c := range tree.Children {
		if c.Show {
			cl = append(cl, c.Category)
		}
	}
	// 返回一级分类列表数据
	return cl
}

// SearchFilmInfo 获取关键字匹配的影片信息
func (i *IndexLogic) SearchFilmInfo(key string, page *system.Page) []system.MovieBasicInfo {
	// 1. 从mysql中获取满足条件的数据, 每页10条
	ids := system.SearchFilmKeyword(key, page)
	// 2. 通过ids获取对应的影片信息
	return system.GetBasicInfoByIds(ids)
}

// GetFilmCategory 根据Pid或Cid获取指定的分页数据
func (i *IndexLogic) GetFilmCategory(id int64, idType string, page *system.Page) []system.MovieBasicInfo {
	// 1. 根据不同类型进不同的查找
	var basicList []system.MovieBasicInfo
	switch idType {
	case "pid":
		basicList = system.GetMovieListByPid(id, page)
	case "cid":
		basicList = system.GetMovieListByCid(id, page)
	}
	return basicList
}

// GetPidCategory 获取pid对应的分类信息
func (i *IndexLogic) GetPidCategory(pid int64) *system.CategoryTree {
	tree := system.GetCategoryTree()
	for _, t := range tree.Children {
		if t.Id == pid {
			return t
		}
	}
	return nil
}

// RelateMovie 根据当前影片信息匹配相关的影片
func (i *IndexLogic) RelateMovie(detail system.MovieDetailVo, page *system.Page) []system.MovieBasicInfo {
	/*
		根据当前影片信息匹配相关的影片
		1. 分类Cid,
		2. 影片名Name
		3. 剧情内容标签class_tag
		4. 地区 area
		5. 语言 Language
	*/
	search := system.SearchInfo{
		Cid:      detail.Cid,
		Name:     detail.Name,
		ClassTag: detail.ClassTag,
		Area:     detail.Area,
		Language: detail.Language,
	}
	return system.GetRelateMovieBasicInfo(search, page)
}

// SearchTags 整合对应分类的搜索tag
func (i *IndexLogic) SearchTags(pid int64) map[string]interface{} {
	// 通过pid 获取对应分类的 tags
	return system.GetSearchTag(pid)
}

/*
		将多个站点的对应影视播放源追加到主站点播放列表中
	 1. 将主站点影片的name 和 subtitle 进行处理添加到用于匹配对应播放源的map中
	 2. 仅对主站点影片name进行映射关系处理并将结果添加到map中
	    例如: xxx第一季  xxx
*/
func multipleSource(detail *system.MovieDetail) []system.PlayLinkVo {
	// 生成多站点的播放源信息
	master := system.GetCollectSourceListByGrade(system.MasterCollect)
	var l = []system.PlayLinkVo{{master[0].Id, master[0].Name, detail.PlayList[0]}}
	// 通过 name 以及 subTitle  生成 hash id  和 dbID 匹配次级站点播放信息
	// 使用map 防止清洗后的id重复
	idMap := make(map[string]int)
	idMap[system.GenerateHashKey(detail.Mid)] = 0
	// 将subTitle进行切割
	if len(detail.SubTitle) > 0 {
		for _, s := range strings.Split(util.FormatSpecialChar(detail.SubTitle), ",") {
			idMap[system.GenerateHashKey(s)] = 0
		}
	}
	// 遍历idMqp整合ids
	var ids []string
	for id, _ := range idMap {
		ids = append(ids, id)
	}
	// 获取附属站点的基本信息
	sMap := make(map[string]system.FilmSource)
	for _, c := range system.GetCollectSourceListByGrade(system.SlaveCollect) {
		sMap[c.Id] = c
	}
	// 获取满足条件的次级站点播放数据
	for _, s := range system.GetMultiplePlay(ids, detail.DbId) {
		l = append(l, system.PlayLinkVo{Id: s.Mid, Name: sMap[s.Mid].Name, LinkList: s.PlayList[0]})
	}
	return l
}

// GetFilmsByTags 通过searchTag 返回满足条件的分页影片信息
func (i *IndexLogic) GetFilmsByTags(st system.SearchTagsVO, page *system.Page) []system.MovieBasicInfo {
	// 获取满足条件的影片id 列表
	ids := system.GetSearchInfosByTags(st, page)
	// 通过key 获取对应影片的基本信息
	return system.GetBasicInfoByIds(ids)
}

// GetFilmClassify 通过Pid返回当前所属分类下的首页展示数据
func (i *IndexLogic) GetFilmClassify(pid int64, page *system.Page) map[string]interface{} {
	res := make(map[string]interface{})
	// 最新上映 (上映时间)
	res["news"] = system.GetMovieListBySort(0, pid, page)
	// 排行榜 (暂定为热度排行)
	res["top"] = system.GetMovieListBySort(1, pid, page)
	// 最近更新 (更新时间)
	res["recent"] = system.GetMovieListBySort(2, pid, page)

	return res

}
