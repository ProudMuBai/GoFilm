package logic

import (
	"errors"
	"fmt"
	"server/model/system"
	"server/plugin/common/util"
	"server/plugin/spider"
	"time"
)

type CronLogic struct {
}

var CL *CronLogic

// AddFilmCrontab 添加影片更新任务
func (cl *CronLogic) AddFilmCrontab(cv system.FilmCronVo) error {
	// 如果 spec 表达式校验失败则直接返回错误信息并终止
	if err := spider.ValidSpec(cv.Spec); err != nil {
		return err
	}
	// 生成任务信息 生成一个唯一ID 作为Task唯一标识
	task := system.FilmCollectTask{Id: util.GenerateSalt(), Ids: cv.Ids, Time: cv.Time, Spec: cv.Spec, Model: cv.Model, State: cv.State, Remark: cv.Remark}
	// 添加一条定时任务
	switch task.Model {
	case 0:
		cid, err := spider.AddAutoUpdateCron(task.Id, task.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			return errors.New(fmt.Sprint("影视自动更新任务添加失败: ", err.Error()))
		}
		// 将定时任务Id记录到Task中
		task.Cid = cid
	case 1:
		cid, err := spider.AddFilmUpdateCron(task.Id, task.Spec)
		// 如果任务添加失败则直接返回错误信息
		if err != nil {
			return errors.New(fmt.Sprint("影视更新定时任务添加失败: ", err.Error()))
		}
		// 将定时任务Id记录到Task中
		task.Cid = cid
	}
	// 如果没有异常则将当前定时任务信息记录到redis中
	system.SaveFilmTask(task)
	return nil
}

// GetFilmCrontab 获取所有定时任务信息
func (cl *CronLogic) GetFilmCrontab() []system.CronTaskVo {
	var l []system.CronTaskVo
	tl := system.GetAllFilmTask()
	for _, t := range tl {
		e := spider.GetEntryById(t.Cid)
		taskVo := system.CronTaskVo{FilmCollectTask: t, PreV: e.Prev.Format(time.DateTime), Next: e.Next.Format(time.DateTime)}
		l = append(l, taskVo)
	}
	return l
}

// GetFilmCrontabById 通过ID获取对应的定时任务信息
func (cl *CronLogic) GetFilmCrontabById(id string) (system.FilmCollectTask, error) {
	t, err := system.GetFilmTaskById(id)
	//e := spider.GetEntryById(t.Cid)
	//taskVo := system.CronTaskVo{FilmCollectTask: t, PreV: e.Prev.Format(time.DateTime), Next: e.Next.Format(time.DateTime)}
	return t, err
}

// ChangeFilmCrontab 改变定时任务的状态 开启 | 停止
func (cl *CronLogic) ChangeFilmCrontab(id string, state bool) error {
	// 通过定时任务信息的唯一标识获取对应的定时任务信息
	ft, err := system.GetFilmTaskById(id)
	if err != nil {
		return errors.New(fmt.Sprintf("定时任务停止失败: %s", err.Error()))
	}
	// 修改当前定时任务的状态为 false, 则在定时执行方法时不会执行具体逻辑
	ft.State = state
	system.UpdateFilmTask(ft)
	return err
}

// UpdateFilmCron 更新定时任务的状态信息
func (cl *CronLogic) UpdateFilmCron(t system.FilmCollectTask) {
	system.UpdateFilmTask(t)
}

// DelFilmCrontab 删除定时任务
func (cl *CronLogic) DelFilmCrontab(id string) error {
	// 通过定时任务信息的唯一Id标识获取对应的定时任务信息
	ft, err := system.GetFilmTaskById(id)
	if err != nil {
		return errors.New(fmt.Sprintf("定时任务删除失败: %s", err.Error()))
	}
	// 通过定时任务EntryID移出对应的定时任务
	spider.RemoveCron(ft.Cid)
	// 将定时任务相关信息删除
	system.DelFilmTask(id)
	return nil
}
