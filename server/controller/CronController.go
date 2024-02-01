package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/logic"
	"server/model/system"
	"server/plugin/spider"
	"strings"
)

// ------------------------------------------------------ 定时任务管理 ------------------------------------------------------

// FilmCronTaskList 获取所有的定时任务信息
func FilmCronTaskList(c *gin.Context) {
	tl := logic.CL.GetFilmCrontab()
	if len(tl) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "暂无任务定时任务信息",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   tl,
	})
}

// GetFilmCronTask 通过Id获取对应的定时任务信息
func GetFilmCronTask(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "定时任务信息获取失败,任务Id不能为空",
		})
		return
	}
	task, err := logic.CL.GetFilmCrontabById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("定时任务信息获取失败", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": StatusOk,
		"data":   task,
	})
}

// FilmCronAdd 添加定时任务
func FilmCronAdd(c *gin.Context) {
	var vo = system.FilmCronVo{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 校验请求参数
	if err := validTaskAddVo(vo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	// 去除cron表达式左右空格
	vo.Spec = strings.TrimSpace(vo.Spec)
	// 执行 定时任务信息保存逻辑
	if err := logic.CL.AddFilmCrontab(vo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("定时任务添加失败: ", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "定时任务添加成功",
	})

}

// FilmCronUpdate 更新定时任务信息
func FilmCronUpdate(c *gin.Context) {
	var t = system.FilmCollectTask{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 校验必要参数
	if err := validTaskInfo(t); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	// 获取未更新的task信息
	task, err := logic.CL.GetFilmCrontabById(t.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("更新失败: ", err.Error()),
		})
		return
	}
	// 将task的可变更属性进行变更
	task.Ids = t.Ids
	task.Time = t.Time
	task.State = t.State
	task.Remark = t.Remark
	// 将变更后的task更新到系统中
	logic.CL.UpdateFilmCron(task)
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": fmt.Sprintf("定时任务[%s]更新成功", task.Id),
	})
}

// ChangeTaskState 开启 | 关闭Id 对应的定时任务
func ChangeTaskState(c *gin.Context) {
	var t = system.FilmCollectTask{}
	// 获取请求参数
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "请求参数异常",
		})
		return
	}
	// 获取未更新的task信息
	task, err := logic.CL.GetFilmCrontabById(t.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("更新失败: ", err.Error()),
		})
		return
	}
	// 修改task的状态
	task.State = t.State
	// 将变更后的task更新到系统中
	logic.CL.UpdateFilmCron(task)
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": fmt.Sprintf("定时任务[%s]更新成功", task.Id),
	})
}

// DelFilmCron 删除定时任务
func DelFilmCron(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "删除失败,任务Id不能为空",
		})
		return
	}
	// 如果Id不为空则执行删除逻辑
	if err := logic.CL.DelFilmCrontab(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": fmt.Sprintf("定时任务[%s]已删除", id),
	})

}

// -------------------------------------------------- 参数校验 --------------------------------------------------

// 定时任务必要属性校验
func validTaskInfo(t system.FilmCollectTask) error {
	if len(t.Id) <= 0 {
		return errors.New("参数校验失败, 任务Id信息不能为空")
	}
	if t.Time == 0 {
		return errors.New("参数校验失败, 采集时长不能为零值")
	}
	return nil
}

// 任务添加参数校验
func validTaskAddVo(vo system.FilmCronVo) error {
	if vo.Model != 0 && vo.Model != 1 {
		return errors.New("参数校验失败, 未定义的任务类型")
	}
	if vo.Time == 0 {
		return errors.New("参数校验失败, 采集时长不能为零值")
	}
	if err := spider.ValidSpec(vo.Spec); err != nil {
		return errors.New(fmt.Sprint("参数校验失败 cron表达式校验失败: ", err.Error()))
	}
	if vo.Model == 1 && (vo.Ids == nil || len(vo.Ids) <= 0) {
		return errors.New("参数校验失败, 自定义更新未绑定任何资源站点")
	}
	return nil
}
