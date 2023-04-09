package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/config"
	"server/model"
	"server/plugin/spider"
)

// SpiderRe 数据清零重开
func SpiderRe(c *gin.Context) {
	// 获取指令参数
	cip := c.Query("cipher")
	if cip != config.SpiderCipher {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "指令错误无法进行此操作",
		})
		return
	}
	// 如果指令正确,则执行重制
	spider.StartSpiderRe()
}

// FixFilmDetail 修复因网络异常造成的影片详情数据丢失
func FixFilmDetail(c *gin.Context) {
	// 获取指令参数
	cip := c.Query("cipher")
	if cip != config.SpiderCipher {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "指令错误无法进行此操作",
		})
		return
	}
	// 如果指令正确,则执行详情数据获取
	spider.MainSiteSpider()
	log.Println("FilmDetail 重制完成!!!")
	// 先截断表中的数据
	model.TunCateSearchTable()
	// 重新扫描完整的信息到mysql中
	spider.SearchInfoToMdb()
	log.Println("SearchInfo 重制完成!!!")
}
