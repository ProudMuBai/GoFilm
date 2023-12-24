package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
)

/*
	影视采集站点信息
*/

type SourceGrade int

const (
	MasterCollect SourceGrade = iota
	SlaveCollect
)

type CollectResultModel int

const (
	JsonResult CollectResultModel = iota
	XmlResult
)

type ResourceType int

func (rt ResourceType) GetActionType() string {
	var ac string = ""
	switch rt {
	case CollectVideo:
		ac = "detail"
	case CollectArticle:
		ac = "article"
	case CollectActor:
		ac = "actor"
	case CollectRole:
		ac = "role"
	case CollectWebSite:
		ac = "web"
	default:
		ac = "detail"
	}
	return ac
}

const (
	CollectVideo = iota
	CollectArticle
	CollectActor
	CollectRole
	CollectWebSite
)

// FilmSource 影视站点信息保存结构体
type FilmSource struct {
	Id           string             `json:"id"`           // 唯一ID
	Name         string             `json:"name"`         // 采集站点备注名
	Uri          string             `json:"uri"`          // 采集链接
	ResultModel  CollectResultModel `json:"resultModel"`  // 接口返回类型, json || xml
	Grade        SourceGrade        `json:"grade"`        // 采集站等级 主站点 || 附属站
	SyncPictures bool               `json:"syncPictures"` // 是否同步图片到服务器
	CollectType  ResourceType       `json:"collectType"`  // 采集资源类型
	State        bool               `json:"state"`        // 是否启用
	Interval     int                `json:"interval"`     // 采集时间间隔 单位/ms
}

// SaveCollectSourceList 保存采集站Api列表
func SaveCollectSourceList(list []FilmSource) error {
	var zl []redis.Z
	for _, v := range list {
		m, _ := json.Marshal(v)
		zl = append(zl, redis.Z{Score: float64(v.Grade), Member: m})
	}
	return db.Rdb.ZAdd(db.Cxt, config.FilmSourceListKey, zl...).Err()
}

// GetCollectSourceList 获取采集站API列表
func GetCollectSourceList() []FilmSource {
	l, err := db.Rdb.ZRange(db.Cxt, config.FilmSourceListKey, 0, -1).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return getCollectSource(l)
}

// GetCollectSourceListByGrade 返回指定类型的采集Api信息 Master | Slave
func GetCollectSourceListByGrade(grade SourceGrade) []FilmSource {
	s := fmt.Sprintf("%d", grade)
	zl, err := db.Rdb.ZRangeByScore(db.Cxt, config.FilmSourceListKey, &redis.ZRangeBy{Max: s, Min: s}).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return getCollectSource(zl)
}

// FindCollectSourceById 通过Id标识获取对应的资源站信息
func FindCollectSourceById(id string) *FilmSource {
	for _, v := range GetCollectSourceList() {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

// 将 []string 转化为 []FilmSourceApi
func getCollectSource(sl []string) []FilmSource {
	var l []FilmSource
	for _, s := range sl {
		f := FilmSource{}
		_ = json.Unmarshal([]byte(s), &f)
		l = append(l, f)
	}
	return l
}

// DelCollectResource 通过Id删除对应的采集站点信息
func DelCollectResource(id string) {
	for _, v := range GetCollectSourceList() {
		if v.Id == id {
			data, _ := json.Marshal(v)
			db.Rdb.ZRem(db.Cxt, config.FilmSourceListKey, data)
		}
	}
}

// AddCollectSource 添加采集站信息
func AddCollectSource(s FilmSource) error {
	for _, v := range GetCollectSourceList() {
		if v.Uri == s.Uri {
			return errors.New("当前采集站点信息已存在, 请勿重复添加")
		}
	}
	// 生成一个短uuid
	s.Id = util.GenerateSalt()
	data, _ := json.Marshal(s)
	return db.Rdb.ZAddNX(db.Cxt, config.FilmSourceListKey, redis.Z{Score: float64(s.Grade), Member: data}).Err()
}

// UpdateCollectSource 更新采集站信息
func UpdateCollectSource(s FilmSource) error {
	for _, v := range GetCollectSourceList() {
		if v.Id != s.Id && v.Uri == s.Uri {
			return errors.New("当前采集站链接已存在其他站点中, 请勿重复添加")
		} else if v.Id == s.Id {
			// 删除当前旧的采集信息
			DelCollectResource(s.Id)
			// 将新的采集信息存入list中
			data, _ := json.Marshal(s)
			db.Rdb.ZAdd(db.Cxt, config.FilmSourceListKey, redis.Z{Score: float64(s.Grade), Member: data})
		}
	}
	return nil
}

// ClearAllCollectSource 删除所有采集站信息
func ClearAllCollectSource() {
	db.Rdb.Del(db.Cxt, config.FilmSourceListKey)
}

// ExistCollectSourceList 查询是否已经存在站点list相关数据
func ExistCollectSourceList() bool {
	if db.Rdb.Exists(db.Cxt, config.FilmSourceListKey).Val() == 0 {
		return false
	}
	return true

}
