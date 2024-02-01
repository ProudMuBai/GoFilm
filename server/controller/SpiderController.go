package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model/system"
	"strconv"
)

// CollectFilm 开启ID对应的资源站的采集任务
func CollectFilm(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	hourStr := c.DefaultQuery("h", "0")
	if id == "" || hourStr == "0" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "采集任务开启失败, 缺乏必要参数",
		})
		return
	}
	h, err := strconv.Atoi(hourStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "采集任务开启失败, hour(时长)参数不符合规范",
		})
		return
	}
	// 执行采集逻处理逻辑
	if err = logic.SL.StartCollect(id, h); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("采集任务开启失败: ", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "采集任务已成功开启!!!",
	})
}

// StarSpider 开启并执行采集任务
func StarSpider(c *gin.Context) {
	var cp system.CollectParams
	// 获取请求参数
	if err := c.ShouldBindJSON(&cp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	if cp.Time == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "采集开启失败,采集时长不能为0",
		})
		return
	}
	// 根据 Batch 执行对应的逻辑
	if cp.Batch {
		// 执行批量采集
		if cp.Ids == nil || len(cp.Ids) <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": "批量采集开启失败, 关联的资源站信息为空",
			})
			return
		}
		// 执行批量采集
		logic.SL.BatchCollect(cp.Time, cp.Ids)
	} else {
		if len(cp.Id) <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": "批量采集开启失败, 资源站Id获取失败",
			})
			return
		}
		if err := logic.SL.StartCollect(cp.Id, cp.Time); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": fmt.Sprint("采集任务开启失败: ", err.Error()),
			})
			return
		}
	}
	// 返回成功执行的信息
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "采集任务已成功开启!!!",
	})
}

// SpiderReset 重置影视数据, 清空库存, 从零开始
func SpiderReset(c *gin.Context) {
	var cp system.CollectParams
	// 获取请求参数
	if err := c.ShouldBindJSON(&cp); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	if cp.Time == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "采集开启失败,采集时长不能为0",
		})
		return
	}
	// 后期加入一些前置验证
	if len(cp.Id) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "SpiderReset Failed, 资源站Id获取失败",
		})
		return
	}
	logic.SL.ZeroCollect(cp.Time)
}

// CoverFilmClass 重置覆盖影片分类信息
func CoverFilmClass(c *gin.Context) {
	// 执行分类采集, 覆盖当前分类信息
	if err := logic.SL.FilmClassCollect(); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("影视分类信息重置成功, 请稍等片刻后刷新页面", c)
}
