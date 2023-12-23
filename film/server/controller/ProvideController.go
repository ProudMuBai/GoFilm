package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model/system"
	"strconv"
)

// 提供用于第三方站点采集的API

// HandleProvide 返回视频列表信息
func HandleProvide(c *gin.Context) {
	// 将请求参数封装为一个map
	var params = map[string]string{
		"t": c.DefaultQuery("t", ""),
		//"pg":  c.DefaultQuery("pg", ""),
		"wd":  c.DefaultQuery("wd", ""),
		"h":   c.DefaultQuery("h", ""),
		"ids": c.DefaultQuery("ids", ""),
	}
	// 设置分页信息
	currentStr := c.DefaultQuery("pg", "1")
	pageSizeStr := c.DefaultQuery("limit", "20")
	current, _ := strconv.Atoi(currentStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	page := system.Page{PageSize: pageSize, Current: current}
	// ac-请求类型 t-类别ID pg-页码 wd-搜索关键字 h=几小时内的数据 ids-数据ID 多个ID逗好分割
	var ac string = c.DefaultQuery("ac", "")
	switch ac {
	case "list":
		c.JSON(http.StatusOK, logic.PL.GetFilmListPage(params, &page))
	case "detail", "videolist":
		c.JSON(http.StatusOK, logic.PL.GetFilmDetailPage(params, &page))
	default:
		c.JSON(http.StatusOK, logic.PL.GetFilmListPage(params, &page))
	}
}

// HandleProvideXml 处理返回xml格式的数据
func HandleProvideXml(c *gin.Context) {
	// 将请求参数封装为一个map
	var params = map[string]string{
		"t": c.DefaultQuery("t", ""),
		//"pg":  c.DefaultQuery("pg", ""),
		"wd":  c.DefaultQuery("wd", ""),
		"h":   c.DefaultQuery("h", ""),
		"ids": c.DefaultQuery("ids", ""),
	}
	// 设置分页信息
	currentStr := c.DefaultQuery("pg", "1")
	pageSizeStr := c.DefaultQuery("limit", "20")
	current, _ := strconv.Atoi(currentStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	page := system.Page{PageSize: pageSize, Current: current}
	// ac-请求类型 t-类别ID pg-页码 wd-搜索关键字 h=几小时内的数据 ids-数据ID 多个ID逗好分割
	var ac string = c.DefaultQuery("ac", "")
	switch ac {
	case "list":
		c.XML(http.StatusOK, logic.PL.GetFilmListXmlPage(params, &page))
	case "detail", "videolist":
		c.XML(http.StatusOK, logic.PL.GetFilmDetailXmlPage(params, &page))
	default:
		c.XML(http.StatusOK, logic.PL.GetFilmListXmlPage(params, &page))
	}
}
