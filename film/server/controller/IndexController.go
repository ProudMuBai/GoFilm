package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model/system"
	"strconv"
	"strings"
)

const (
	StatusOk     = "ok"
	StatusFailed = "failed"
)

// Index 首页数据
func Index(c *gin.Context) {
	// 获取首页所需数据
	data := logic.IL.IndexPage()
	system.Success(data, "首页数据获取成功", c)
}

// CategoriesInfo 分类信息获取
func CategoriesInfo(c *gin.Context) {
	//data := logic.IL.GetCategoryInfo()
	data := logic.IL.GetNavCategory()

	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			`status`:  StatusFailed,
			`message`: `暂无分类信息!!!`,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		`status`: StatusOk,
		`data`:   data,
	})
}

// FilmDetail 影片详情信息查询
func FilmDetail(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		system.Failed("请求异常,影片请求参数异常!!!", c)
		return
	}
	// 获取影片详情信息
	detail := logic.IL.GetFilmDetail(id)
	// 获取相关推荐影片数据
	page := system.Page{Current: 0, PageSize: 14}
	relateMovie := logic.IL.RelateMovie(detail.MovieDetail, &page)
	system.Success(gin.H{
		"detail": detail,
		"relate": relateMovie,
	}, "影片详情信息获取成功", c)
}

// FilmPlayInfo 影视播放页数据
func FilmPlayInfo(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	playFrom := c.DefaultQuery("playFrom", "")
	episode, err := strconv.Atoi(c.DefaultQuery("episode", "0"))
	if err != nil {
		system.Failed("请求异常,暂无影片信息!!!", c)
		return
	}
	// 获取影片详情信息
	detail := logic.IL.GetFilmDetail(id)
	// 如果 playFrom 为空, 则设置默认播放源和默认影片数据
	if len(playFrom) <= 1 && len(detail.List) > 0 {
		playFrom = detail.List[0].Id

	}
	// 获取当前影片播放信息
	var currentPlay system.MovieUrlInfo
	for _, v := range detail.List {
		if v.Id == playFrom {
			currentPlay = v.LinkList[episode]
		}
	}

	// 推荐影片信息
	page := system.Page{Current: 0, PageSize: 14}
	relateMovie := logic.IL.RelateMovie(detail.MovieDetail, &page)
	system.Success(gin.H{
		"detail":          detail,
		"current":         currentPlay,
		"currentPlayFrom": playFrom,
		"currentEpisode":  episode,
		"relate":          relateMovie,
	}, "影片播放信息获取成功", c)
}

// SearchFilm 通过片名模糊匹配库存中的信息
func SearchFilm(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	currStr := c.DefaultQuery("current", "1")
	current, _ := strconv.Atoi(currStr)
	page := system.Page{PageSize: 10, Current: current}
	bl := logic.IL.SearchFilmInfo(strings.TrimSpace(keyword), &page)

	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"list": bl,
			"page": page,
		},
	})
}

// FilmTagSearch 通过tag获取满足条件的对应影片
func FilmTagSearch(c *gin.Context) {
	params := system.SearchTagsVO{}
	pidStr := c.DefaultQuery("Pid", "")
	cidStr := c.DefaultQuery("Category", "")
	yStr := c.DefaultQuery("Year", "")
	if pidStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "缺少分类信息",
		})
		return
	}
	params.Pid, _ = strconv.ParseInt(pidStr, 10, 64)
	params.Cid, _ = strconv.ParseInt(cidStr, 10, 64)
	params.Plot = c.DefaultQuery("Plot", "")
	params.Area = c.DefaultQuery("Area", "")
	params.Language = c.DefaultQuery("Language", "")
	params.Year, _ = strconv.ParseInt(yStr, 10, 64)
	params.Sort = c.DefaultQuery("Sort", "update_stamp")

	// 设置分页信息
	currentStr := c.DefaultQuery("current", "1")
	current, _ := strconv.Atoi(currentStr)
	page := system.Page{PageSize: 49, Current: current}
	logic.IL.GetFilmsByTags(params, &page)
	// 获取当前分类Title
	// 返回对应信息
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"title":  logic.IL.GetPidCategory(params.Pid).Category,
			"list":   logic.IL.GetFilmsByTags(params, &page),
			"search": logic.IL.SearchTags(params.Pid),
			"params": map[string]string{
				"Pid":      pidStr,
				"Category": cidStr,
				"Plot":     params.Plot,
				"Area":     params.Area,
				"Language": params.Language,
				"Year":     yStr,
				"Sort":     params.Sort,
			},
		},
		"page": page,
	})
}

// FilmClassify  影片分类首页数据展示
func FilmClassify(c *gin.Context) {
	pidStr := c.DefaultQuery("Pid", "")
	if pidStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "缺少分类信息",
		})
		return
	}
	// 1. 顶部Title数据
	pid, _ := strconv.ParseInt(pidStr, 10, 64)
	title := logic.IL.GetPidCategory(pid)
	// 2. 设置分页信息
	page := system.Page{PageSize: 21, Current: 1}
	// 3. 获取当前分类下的 最新上映, 排行榜, 最近更新 影片信息
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"title":   title,
			"content": logic.IL.GetFilmClassify(pid, &page),
		},
		"page": page,
	})
}
