package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/config"
	"server/model"
	"server/plugin/db"
)

/*
*
 IndexController数据处理
*/

type IndexLogic struct {
}

var IL *IndexLogic

// IndexPage 首页数据处理
func (i *IndexLogic) IndexPage() gin.H {
	Info := gin.H{}
	// 首页分类数据处理

	// 1. 导航分类数据处理, 只提供 电影 电视剧 综艺 动漫 四大顶级分类和其子分类
	tree := model.CategoryTree{Category: &model.Category{Id: 0, Name: "分类信息"}}
	sysTree := model.GetCategoryTree()
	//  由于采集源数据格式不一,因此采用名称匹配
	for _, c := range sysTree.Children {
		switch c.Category.Name {
		case "电影", "电影片", "连续剧", "电视剧", "综艺", "综艺片", "动漫", "动漫片":
			tree.Children = append(tree.Children, c)
		}
	}
	Info["category"] = tree
	// 2. 提供用于首页展示的顶级分类影片信息, 每分类 14条数据
	var list []gin.H
	for _, c := range tree.Children {
		page := model.Page{PageSize: 14, Current: 1}
		movies := model.GetMovieListByPid(c.Id, &page)
		item := gin.H{"nav": c, "movies": movies}
		list = append(list, item)
	}
	Info["content"] = list

	return Info
}

// GetFilmDetail 影片详情信息页面处理
func (i *IndexLogic) GetFilmDetail(id int) model.MovieDetail {
	// 通过Id 获取影片search信息
	search := model.SearchInfo{}
	db.Mdb.Where("mid", id).First(&search)
	// 获取redis中的完整影视信息 MovieDetail:Cid11:Id24676
	movieDetail := model.GetDetailByKey(fmt.Sprintf(config.MovieDetailKey, search.Cid, search.Mid))
	return movieDetail
}

// GetCategoryInfo 分类信息获取, 组装导航栏需要的信息
func (i *IndexLogic) GetCategoryInfo() gin.H {
	// 组装nav导航所需的信息
	nav := gin.H{}
	// 1.获取所有分类信息
	tree := model.GetCategoryTree()
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

// SearchFilmInfo 获取关键字匹配的影片信息
func (i *IndexLogic) SearchFilmInfo(key string, page *model.Page) []model.MovieBasicInfo {
	// 1. 从mysql中获取满足条件的数据, 每页10条
	sl := model.SearchFilmKeyword(key, page)
	// 2. 获取redis中的basicMovieInfo信息
	var bl []model.MovieBasicInfo
	for _, s := range sl {
		bl = append(bl, model.GetBasicInfoByKey(fmt.Sprintf(config.MovieBasicInfoKey, s.Cid, s.Mid)))
	}
	return bl
}

// GetFilmCategory 根据Pid或Cid获取指定的分页数据
func (i *IndexLogic) GetFilmCategory(id int64, idType string, page *model.Page) []model.MovieBasicInfo {
	// 1. 根据不同类型进不同的查找
	var basicList []model.MovieBasicInfo
	switch idType {
	case "pid":
		basicList = model.GetMovieListByPid(id, page)
	case "cid":
		basicList = model.GetMovieListByCid(id, page)
	}
	return basicList
}

// GetPidCategory 获取pid对应的分类信息
func (i *IndexLogic) GetPidCategory(pid int64) *model.CategoryTree {
	tree := model.GetCategoryTree()
	for _, t := range tree.Children {
		if t.Id == pid {
			return t
		}
	}
	return nil
}

// RelateMovie 根据当前影片信息匹配相关的影片
func (i *IndexLogic) RelateMovie(detail model.MovieDetail, page *model.Page) []model.MovieBasicInfo {
	/*
		根据当前影片信息匹配相关的影片
		1. 分类Cid,
		2. 影片名Name
		3. 剧情内容标签class_tag
		4. 地区 area
		5. 语言 Language
	*/
	search := model.SearchInfo{
		Cid:      detail.Cid,
		Name:     detail.Name,
		ClassTag: detail.ClassTag,
		Area:     detail.Area,
		Language: detail.Language,
	}
	return model.GetRelateMovieBasicInfo(search, page)
}
