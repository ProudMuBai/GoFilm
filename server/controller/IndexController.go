package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model"
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
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   data,
	})
}

// CategoriesInfo 分类信息获取
func CategoriesInfo(c *gin.Context) {
	data := logic.IL.GetCategoryInfo()

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
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求异常,暂无影片信息!!!",
		})
		return
	}
	// 获取影片详情信息
	detail := logic.IL.GetFilmDetail(id)
	// 获取相关推荐影片数据
	page := model.Page{Current: 0, PageSize: 14}
	relateMovie := logic.IL.RelateMovie(detail, &page)
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"detail": detail,
			"relate": relateMovie,
		},
	})
}

// FilmPlayInfo 影视播放页数据
func FilmPlayInfo(c *gin.Context) {
	// 获取请求参数
	id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	playFrom, err := strconv.Atoi(c.DefaultQuery("playFrom", "0"))
	episode, err := strconv.Atoi(c.DefaultQuery("episode", "0"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求异常,暂无影片信息!!!",
		})
		return
	}
	// 获取影片详情信息
	detail := logic.IL.GetFilmDetail(id)
	// 推荐影片信息
	page := model.Page{Current: 0, PageSize: 14}
	relateMovie := logic.IL.RelateMovie(detail, &page)
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"detail":          detail,
			"current":         detail.PlayList[playFrom][episode],
			"currentPlayFrom": playFrom,
			"currentEpisode":  episode,
			"relate":          relateMovie,
		},
	})
}

// SearchFilm 通过片名模糊匹配库存中的信息
func SearchFilm(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	currStr := c.DefaultQuery("current", "1")
	current, _ := strconv.Atoi(currStr)
	page := model.Page{PageSize: 10, Current: current}
	bl := logic.IL.SearchFilmInfo(strings.TrimSpace(keyword), &page)

	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"list": bl,
			"page": page,
		},
	})
}

// FilmCategory 获取指定分类的影片分页数据,
func FilmCategory(c *gin.Context) {
	// 1.1 首先获取Cid 二级分类id是否存在
	cidStr := c.DefaultQuery("cid", "")
	// 1.2 如果pid也不存在直接返回错误信息
	pidStr := c.DefaultQuery("pid", "")
	if pidStr == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "缺少分类信息",
		})
		return
	}
	// 1.3 获取pid对应的分类信息
	pid, _ := strconv.ParseInt(pidStr, 10, 64)
	category := logic.IL.GetPidCategory(pid)

	// 2 设置分页信息
	currentStr := c.DefaultQuery("current", "1")
	current, _ := strconv.Atoi(currentStr)
	page := model.Page{PageSize: 49, Current: current}
	// 2.1 如果不存在cid则根据Pid进行查询
	if cidStr == "" {
		// 2.2 如果存在pid则根据pid进行查找
		c.JSON(http.StatusOK, gin.H{
			"status": StatusOk,
			"data": gin.H{
				"list":     logic.IL.GetFilmCategory(pid, "pid", &page),
				"category": category,
			},
			"page": page,
		})
		return
	}
	// 2.2 如果存在cid 则根据具体的cid去查询数据
	cid, _ := strconv.ParseInt(cidStr, 10, 64)
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data": gin.H{
			"list":     logic.IL.GetFilmCategory(cid, "cid", &page),
			"category": category,
		},
		"page": page,
	})
}
