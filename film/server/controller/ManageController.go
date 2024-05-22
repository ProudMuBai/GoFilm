package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/logic"
	"server/model/system"
	"server/plugin/SystemInit"
	"server/plugin/common/util"
	"server/plugin/spider"
)

func ManageIndex(c *gin.Context) {
	system.SuccessOnlyMsg("后台管理中心", c)
	return
}

// ------------------------------------------------------ 影视采集 ------------------------------------------------------

// FilmSourceList 采集站点信息
func FilmSourceList(c *gin.Context) {
	system.Success(logic.ML.GetFilmSourceList(), "影视源站点信息获取成功", c)
	return
}

// FindFilmSource 通过ID返回对应的资源站数据
func FindFilmSource(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	fs := logic.ML.GetFilmSource(id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	system.Success(fs, "原站点详情信息查找成功", c)
}

// FilmSourceAdd 添加采集源
func FilmSourceAdd(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		system.Failed("资源接口测试失败, 请确认接口有效再添加", c)
		return
	}
	// 测试通过后将资源站信息添加到list
	if err := logic.ML.SaveFilmSource(s); err != nil {
		system.Failed(fmt.Sprint("资源站添加失败: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("添加成功", c)
}

func FilmSourceUpdate(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (s.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	// 校验Id信息是否为空
	if s.Id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	fs := logic.ML.GetFilmSource(s.Id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	// 如果 uri发生变更则执行spider测试
	if fs.Uri != s.Uri {
		// 执行 spider
		if err := spider.CollectApiTest(s); err != nil {
			system.Failed("资源接口测试失败, 请确认更新的数据接口是否有效", c)
			return
		}
	}
	// 更新资源站信息
	if err := logic.ML.UpdateFilmSource(s); err != nil {
		system.Failed(fmt.Sprint("资源站更新失败: ", err.Error()), c)
		return
	}
	system.SuccessOnlyMsg("更新成功", c)
}

func FilmSourceChange(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	if s.Id == "" {
		system.Failed("参数异常, 资源站标识不能为空", c)
		return
	}
	// 查找对应的资源站点信息
	fs := logic.ML.GetFilmSource(s.Id)
	if fs == nil {
		system.Failed("数据异常,资源站信息不存在", c)
		return
	}
	// 如果采集站开启图片同步, 且采集站为附属站点则返回错误提示
	if s.SyncPictures && (fs.Grade == system.SlaveCollect) {
		system.Failed("附属站点无法开启图片同步功能", c)
		return
	}
	if s.State != fs.State || s.SyncPictures != fs.SyncPictures {
		// 执行更新操作
		s := system.FilmSource{Id: fs.Id, Name: fs.Name, Uri: fs.Uri, ResultModel: fs.ResultModel,
			Grade: fs.Grade, SyncPictures: s.SyncPictures, CollectType: fs.CollectType, State: s.State}
		// 更新资源站信息
		if err := logic.ML.UpdateFilmSource(s); err != nil {
			system.Failed(fmt.Sprint("资源站更新失败: ", err.Error()), c)
			return
		}
	}
	system.SuccessOnlyMsg("更新成功", c)
}

func FilmSourceDel(c *gin.Context) {
	id := c.Query("id")
	if len(id) <= 0 {
		system.Failed("资源站ID信息不能为空", c)
		return
	}
	if err := logic.ML.DelFilmSource(id); err != nil {
		system.Failed("删除资源站失败", c)
		return
	}
	system.SuccessOnlyMsg("删除成功", c)
}

// FilmSourceTest 测试影视站点数据是否可用
func FilmSourceTest(c *gin.Context) {
	var s = system.FilmSource{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&s); err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	// 校验必要参数
	if err := validFilmSource(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 执行 spider
	if err := spider.CollectApiTest(s); err != nil {
		system.Failed(err.Error(), c)
		return
	}
	system.SuccessOnlyMsg("测试成功!!!", c)
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
	system.Success(l, "影视源信息获取成功", c)
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
			system.Failed("域名格式校验失败", c)
			return
		}
		if len(bc.SiteName) <= 0 {
			system.Failed("网站名称不能为空", c)
			return
		}
	} else {
		system.Failed(fmt.Sprint("请求参数异常:  ", err), c)
		return
	}

	// 保存更新后的配置信息
	if err := logic.ML.UpdateSiteBasic(bc); err != nil {
		system.Failed(fmt.Sprint("网站配置更新失败:  ", err), c)
		return
	}
	system.SuccessOnlyMsg("更新成功", c)
}

// ResetSiteBasic 重置网站配置信息为初始化状态
func ResetSiteBasic(c *gin.Context) {
	// 执行配置初始化方法直接覆盖当前基本配置信息
	SystemInit.BasicConfigInit()
	system.SuccessOnlyMsg("配置信息重置成功", c)
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
