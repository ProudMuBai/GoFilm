package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/model/system"
	"strconv"
)

// CollectFilm 开启ID对应的资源站的采集任务
func CollectFilm(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	hourStr := c.DefaultQuery("h", "0")
	if id == "" || hourStr == "0" {
		system.Failed("采集任务开启失败, 缺乏必要参数", c)
		return
	}
	h, err := strconv.Atoi(hourStr)
	if err != nil {
		system.Failed("采集任务开启失败, 采集(时长)不符合规范", c)
		return
	}
	// 执行采集逻处理逻辑
	if err = logic.SL.StartCollect(id, h); err != nil {
		system.Failed(fmt.Sprint("采集任务开启失败: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("采集任务已成功开启!!!", c)
}

// StarSpider 开启并执行采集任务
func StarSpider(c *gin.Context) {
	var cp system.CollectParams
	// 获取请求参数
	if err := c.ShouldBindJSON(&cp); err != nil {
		system.Failed("请求参数异常!!!", c)
		return
	}
	if cp.Time == 0 {
		system.Failed("采集开启失败,采集时长不能为0", c)
		return
	}
	// 根据 Batch 执行对应的逻辑
	if cp.Batch {
		// 执行批量采集
		if cp.Ids == nil || len(cp.Ids) <= 0 {
			system.Failed("批量采集开启失败, 关联的资源站信息为空", c)
			return
		}
		// 执行批量采集
		logic.SL.BatchCollect(cp.Time, cp.Ids)
	} else {
		if len(cp.Id) <= 0 {
			system.Failed("批量采集开启失败, 资源站Id获取失败", c)
			return
		}
		if err := logic.SL.StartCollect(cp.Id, cp.Time); err != nil {
			system.Failed(fmt.Sprint("采集任务开启失败: ", err.Error()), c)
			return
		}
	}
	// 返回成功执行的信息
	system.SuccessOnlyMsg("采集任务已成功开启!!!", c)
}

// SpiderReset 重置影视数据, 清空库存, 从零开始
func SpiderReset(c *gin.Context) {
	// 清空采集数据进行重新采集前校验输入的密码是否正确
	key := c.DefaultQuery("accessKey", "")
	// 如密码错误则不执行后续操作
	if len(key) <= 0 || key != "Re0" {
		system.Failed("重置失败, 密钥校验失败!!!", c)
		return
	}
	// 前置校验通过则清空采集数据并对已启用站点进行 全量采集
	logic.SL.ZeroCollect(-1)
	system.SuccessOnlyMsg("影视数据已重置, 请耐心等待采集完成!!!", c)
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
