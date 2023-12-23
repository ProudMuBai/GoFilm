package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model/system"
	"server/plugin/SystemInit"
	"server/plugin/common/util"
	"server/plugin/spider"
)

func ManageIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "hahah",
	})
	return
}

// ------------------------------------------------------ 影视采集 ------------------------------------------------------

// FilmSourceList 采集站点信息
func FilmSourceList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   logic.ML.GetFilmSourceList(),
	})
	return
}

// FindFilmSource 通过ID返回对应的资源站数据
func FindFilmSource(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "参数异常, 资源站标识不能为空",
		})
		return
	}
	fs := logic.ML.GetFilmSource(id)
	if fs == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "数据异常,资源站信息不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   fs,
	})
}

func FilmSourceAdd(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "附属站点无法开启图片同步功能",
		})
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "资源接口测试失败, 请确认接口有效再添加",
		})
		return
	}
	// 测试通过后将资源站信息添加到list
	if err := logic.ML.SaveFilmSource(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("资源站添加失败: ", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "添加成功",
	})
}

func FilmSourceUpdate(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "附属站点无法开启图片同步功能",
		})
		return
	}
	// 校验Id信息是否为空
	if s.Id == "" {
		c.JSON(http.StatusOK, gin.H{"status": StatusFailed, "message": "参数异常, 资源站标识不能为空"})
		return
	}
	fs := logic.ML.GetFilmSource(s.Id)
	if fs == nil {
		c.JSON(http.StatusOK, gin.H{"status": StatusFailed, "message": "数据异常,资源站信息不存在"})
		return
	}
	// 如果 uri发生变更则执行spider测试
	if fs.Uri != s.Uri {
		// 执行 spider
		if err := spider.CollectApiTest(s); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": "资源接口测试失败, 请确认更新的数据接口是否有效",
			})
			return
		}
	}
	// 更新资源站信息
	if err := logic.ML.UpdateFilmSource(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("资源站更新失败: ", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "更新成功",
	})

}

func FilmSourceChange(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	if s.Id == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "参数异常, 资源站标识不能为空",
		})
		return
	}
	// 查找对应的资源站点信息
	fs := logic.ML.GetFilmSource(s.Id)
	if fs == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "数据异常,资源站信息不存在",
		})
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (fs.Grade == system.SlaveCollect) {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "附属站点无法开启图片同步功能",
		})
		return
	}
	if s.State != fs.State || s.SyncPictures != fs.SyncPictures {
		// 执行更新操作
		s := system.FilmSource{Id: fs.Id, Name: fs.Name, Uri: fs.Uri, ResultModel: fs.ResultModel,
			Grade: fs.Grade, SyncPictures: s.SyncPictures, CollectType: fs.CollectType, State: s.State}
		// 更新资源站信息
		if err := logic.ML.UpdateFilmSource(s); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": fmt.Sprint("资源站更新失败: ", err.Error()),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "更新成功",
	})
}

func FilmSourceDel(c *gin.Context) {
	id := c.Query("id")
	if len(id) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "资源站ID信息不能为空",
		})
		return
	}
	if err := logic.ML.DelFilmSource(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("删除资源站失败: ", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "删除成功:",
	})
}

// FilmSourceTest 测试影视站点数据是否可用
func FilmSourceTest(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "测试成功!!!",
	})
}

// GetNormalFilmSource 获取状态为启用的采集站信息
func GetNormalFilmSource(c *gin.Context) {
	// 获取所有的采集站信息
	var l []system.FilmTaskOptions
	for _, v := range logic.ML.GetFilmSourceList() {
		if v.State {
			l = append(l, system.FilmTaskOptions{Id: v.Id, Name: v.Name})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   l,
	})
}

// ------------------------------------------------------ 站点基本配置 ------------------------------------------------------

// SiteBasicConfig  网站基本配置
func SiteBasicConfig(c *gin.Context) {
	system.Success(logic.ML.GetSiteBasicConfig(), "网站基本信息获取成功", c)
}

// UpdateSiteBasic 更新网站配置信息
func UpdateSiteBasic(c *gin.Context) {
	// 获取请求参数 && 校验关键配置项
	bc := system.BasicConfig{}
	if err := c.ShouldBindJSON(&bc); err == nil {
		// 对参数进行校验
		if !util.ValidDomain(bc.Domain) && !util.ValidIPHost(bc.Domain) {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": "域名格式校验失败: ",
			})
			return
		}
		if len(bc.SiteName) <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  StatusFailed,
				"message": "网站名称不能为空: ",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": StatusOk, "message": fmt.Sprint("参数提交失败:  ", err)})
		return
	}

	// 保存更新后的配置信息
	if err := logic.ML.UpdateSiteBasic(bc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("网站配置更新失败: ", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "更新成功: ",
	})
	return
}

// ResetSiteBasic 重置网站配置信息为初始化状态
func ResetSiteBasic(c *gin.Context) {
	// 执行配置初始化方法直接覆盖当前基本配置信息
	SystemInit.BasicConfigInit()
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "重置成功: ",
	})
}

// ------------------------------------------------------ 参数校验 ------------------------------------------------------
func validFilmSource(fs system.FilmSource) error {
	// 资源名称不能为空 且长度不能超过20
	if len(fs.Name) <= 0 || len(fs.Name) > 20 {
		return errors.New("资源名称不能为空且长度不能超过20")
	}
	// Uri 采集链接测试格式
	if !util.ValidURL(fs.Uri) {
		return errors.New("资源链接格式异常, 请输入规范的URL链接")
	}
	// 校验接口类型是否是 JSON || XML
	if fs.ResultModel != system.JsonResult && fs.ResultModel != system.XmlResult {
		return errors.New("接口类型异常, 请提交正确的接口类型")
	}
	// 校验采集类型是否符合规范
	switch fs.CollectType {
	case system.CollectVideo, system.CollectArticle, system.CollectActor, system.CollectRole, system.CollectWebSite:
		return nil
	default:
		return errors.New("资源类型异常, 未知的资源类型")
	}
}

func apiValidityCheck() {

}
