package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/config"
	"server/logic"
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
	go logic.SL.ReZero()
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "ReZero 任务执行已成功开启",
	})
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
	go logic.SL.FixDetail()
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "FixDetail 任务执行已成功开启",
	})
}

// RefreshSitePlay 清空附属站点影片数据并重新获取
func RefreshSitePlay(c *gin.Context) {
	// 获取指令参数
	cip := c.Query("cipher")
	if cip != config.SpiderCipher {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "指令错误无法进行此操作",
		})
		return
	}

	// 执行多站点播放数据重置
	go logic.SL.SpiderMtPlayRe()
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "SpiderMtPlayRe 任务执行已成功开启",
	})
}
