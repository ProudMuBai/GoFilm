package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"server/config"
	"server/model/system"
	"server/plugin/db"
	"strings"
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
		page := system.Page{PageSize: 14, Current: 1}
		movies := system.GetMovieListByPid(c.Id, &page)
		// 获取当前分类的本月热门影片
		HotMovies := system.GetHotMovieByPid(c.Id, &page)
		item := map[string]interface{}{"nav": c, "movies": movies, "hot": HotMovies}
		list = append(list, item)
	}
	Info["content"] = list
	// 不存在首页数据缓存时将查询数据缓存到redis中
	system.DataCache(config.IndexCacheKey, Info)
	return Info
}

// GetFilmDetail 影片详情信息页面处理
func (i *IndexLogic) GetFilmDetail(id int) system.MovieDetailVo {
	// 通过Id 获取影片search信息
	search := system.SearchInfo{}
	db.Mdb.Where("mid", id).First(&search)
	// 获取redis中的完整影视信息 MovieDetail:Cid11:Id24676
	movieDetail := system.GetDetailByKey(fmt.Sprintf(config.MovieDetailKey, search.Cid, search.Mid))
	var res = system.MovieDetailVo{MovieDetail: movieDetail}
	//查找其他站点是否存在影片对应的播放源
	res.List = multipleSource(&movieDetail)
	return res
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
	sl := system.SearchFilmKeyword(key, page)
	// 2. 获取redis中的basicMovieInfo信息
	var bl []system.MovieBasicInfo
	for _, s := range sl {
		bl = append(bl, system.GetBasicInfoByKey(fmt.Sprintf(config.MovieBasicInfoKey, s.Cid, s.Mid)))
	}
	return bl
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
func (i *IndexLogic) RelateMovie(detail system.MovieDetail, page *system.Page) []system.MovieBasicInfo {
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
	var playList = []system.PlayLinkVo{{master[0].Id, master[0].Name, detail.PlayList[0]}}

	// 整合多播放源, 初始化存储key map
	names := make(map[string]int)
	// 1. 判断detail的dbId是否存在, 存在则添加到names中作为匹配条件
	if detail.DbId > 0 {
		names[system.GenerateHashKey(detail.DbId)] = 0
	}
	// 2. 对name进行去除特殊格式处理
	names[system.GenerateHashKey(detail.Name)] = 0
	// 3. 对包含第一季的name进行处理
	names[system.GenerateHashKey(regexp.MustCompile(`第一季$`).ReplaceAllString(detail.Name, ""))] = 0

	// 4. 将subtitle进行切分,放入names中
	if len(detail.SubTitle) > 0 && strings.Contains(detail.SubTitle, ",") {
		for _, v := range strings.Split(detail.SubTitle, ",") {
			names[system.GenerateHashKey(v)] = 0
		}
	}
	if len(detail.SubTitle) > 0 && strings.Contains(detail.SubTitle, "/") {
		for _, v := range strings.Split(detail.SubTitle, "/") {
			names[system.GenerateHashKey(v)] = 0
		}
	}
	// 遍历所有附属站点列表
	sc := system.GetCollectSourceListByGrade(system.SlaveCollect)
	for _, s := range sc {
		for k, _ := range names {
			pl := system.GetMultiplePlay(s.Id, k)
			if len(pl) > 0 {
				// 如果当前站点已经匹配到数据则直接退出当前循环
				//detail.PlayList = append(detail.PlayList, pl)
				playList = append(playList, system.PlayLinkVo{Id: s.Id, Name: s.Name, LinkList: pl})
				break
			}
		}
	}

	return playList
}

// GetFilmsByTags 通过searchTag 返回满足条件的分页影片信息
func (i *IndexLogic) GetFilmsByTags(st system.SearchTagsVO, page *system.Page) []system.MovieBasicInfo {
	// 获取满足条件的影片id 列表
	sl := system.GetSearchInfosByTags(st, page)
	// 通过key 获取对应影片的基本信息
	return system.GetBasicInfoBySearchInfos(sl...)
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
