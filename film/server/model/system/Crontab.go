package system

import (
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"server/config"
	"server/plugin/db"
)

/*
	定时任务持久化
*/

// FilmCollectTask 影视采集任务
type FilmCollectTask struct {
	Id     string       `json:"id"`     // 唯一标识uid
	Ids    []string     `json:"ids"`    // 采集站id列表
	Cid    cron.EntryID `json:"cid"`    // 定时任务Id
	Time   int          `json:"time"`   // 采集时长, 最新x小时更新的内容
	Spec   string       `json:"spec"`   // 执行周期 cron表达式
	Model  int          `json:"model"`  // 任务类型, 0 - 自动更新已启用站点 || 1 - 更新Ids中的资源站数据
	State  bool         `json:"state"`  // 状态 开启 | 禁用
	Remark string       `json:"remark"` // 任务备注信息
}

// SaveFilmTask 保存影视采集任务信息  {EntryId:FilmCollectTask}
func SaveFilmTask(t FilmCollectTask) {
	data, _ := json.Marshal(t)
	db.Rdb.HSet(db.Cxt, config.FilmCrontabKey, t.Id, data)
}

// GetAllFilmTask 获取所有的任务信息
func GetAllFilmTask() []FilmCollectTask {
	var tl []FilmCollectTask
	tMap := db.Rdb.HGetAll(db.Cxt, config.FilmCrontabKey).Val()
	for _, v := range tMap {
		var t = FilmCollectTask{}
		_ = json.Unmarshal([]byte(v), &t)
		tl = append(tl, t)
	}
	return tl
}

// GetFilmTaskById 通过Id获取当前任务信息
func GetFilmTaskById(id string) (FilmCollectTask, error) {
	var ft = FilmCollectTask{}
	// 如果Id对应的task不存在则返回错误信息
	if !db.Rdb.HExists(db.Cxt, config.FilmCrontabKey, id).Val() {
		return ft, errors.New(" The task does not exist ")
	}
	data := db.Rdb.HGet(db.Cxt, config.FilmCrontabKey, id).Val()
	err := json.Unmarshal([]byte(data), &ft)
	return ft, err
}

// UpdateFilmTask 更新定时任务信息(直接覆盖Id对应的定时任务信息) -- 后续待调整
func UpdateFilmTask(t FilmCollectTask) {
	SaveFilmTask(t)
}

// DelFilmTask 通过Id删除对应的定时任务信息
func DelFilmTask(id string) {
	db.Rdb.HDel(db.Cxt, config.FilmCrontabKey, id)
}

// ExistTask 是否存在定时任务相关信息
func ExistTask() bool {
	return db.Rdb.Exists(db.Cxt, config.FilmCrontabKey).Val() == 1
}
