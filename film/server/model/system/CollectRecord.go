package system

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"server/config"
	"server/plugin/db"
)

// FailureRecord 失败采集记录信息机构体
type FailureRecord struct {
	gorm.Model
	OriginId    string       `json:"originId"`    // 采集站唯一ID
	OriginName  string       `json:"originName"`  // 采集站唯一ID
	Uri         string       `json:"uri"`         // 采集源链接
	CollectType ResourceType `json:"collectType"` // 采集类型
	PageNumber  int          `json:"pageNumber"`  // 页码
	Hour        int          `json:"hour"`        // 采集参数 h 时长
	Cause       string       `json:"cause"`       // 失败原因
	Status      int          `json:"status"`      // 重试状态
}

// TableName 采集失败记录表表名
func (fr FailureRecord) TableName() string {
	return config.FailureRecordTableName
}

// CreateFailureRecordTable 创建失效记录表
func CreateFailureRecordTable() {
	var fl = &FailureRecord{}
	// 不存在则创建FailureRecord表
	if !db.Mdb.Migrator().HasTable(fl) {
		if err := db.Mdb.AutoMigrate(fl); err != nil {
			log.Println("Create Table failure_record failed:", err)
		}
	}
}

// SaveFailureRecord 添加采集失效记录
func SaveFailureRecord(fl FailureRecord) {
	// 数据量不多但存在并发问题, 开启事务
	err := db.Mdb.Transaction(func(tx *gorm.DB) error {
		// 将采集失败信息存储到record表中
		if err := tx.Create(&fl).Error; err != nil {
			log.Println("Add failure record failed:", err)
			return err
		}
		return nil
	})
	// 如果事务提交失败, 则输出相应信息, (存一份数据到Redis??)
	if err != nil {
		log.Println("Save failure record affairs failed:", err)
	}
}

// FailureRecordList 获取所有的采集失效记录
func FailureRecordList(vo RecordRequestVo) []FailureRecord {
	// 通过RecordRequestVo,生成查询条件
	qw := db.Mdb.Model(&FailureRecord{})
	if vo.OriginId != "" {
		qw.Where("origin_id = ?", vo.OriginId)
	}
	if !vo.BeginTime.IsZero() && !vo.EndTime.IsZero() {
		qw.Where("created_at BETWEEN ? AND ? ", vo.BeginTime, vo.EndTime)
	}
	//if vo.CollectType >= 0 {
	//	qw.Where("collect_type = ?", vo.CollectType)
	//}
	//if vo.Hour != 0 {
	//	qw.Where("hour = ?", vo.Hour)
	//}
	if vo.Status >= 0 {
		qw.Where("status = ?", vo.Status)
	}

	// 获取分页数据
	GetPage(qw, vo.Paging)
	// 获取分页查询的数据
	var list []FailureRecord
	if err := qw.Limit(vo.Paging.PageSize).Offset((vo.Paging.Current - 1) * vo.Paging.PageSize).Order("updated_at DESC").Find(&list).Error; err != nil {
		log.Println(err)
		return nil
	}
	return list
}

// FindRecordById 获取id对应的失效记录
func FindRecordById(id uint) *FailureRecord {
	var fr FailureRecord
	fr.ID = id
	// 通过ID查询对应的数据
	db.Mdb.First(&fr)
	return &fr
}

// PendingRecord 查询所有待处理的记录信息
func PendingRecord() []FailureRecord {
	var list []FailureRecord
	// 1. 获取 hour > 4320 || hour < 0  && status = 1 的影片信息
	db.Mdb.Where("(hour > 4320 OR hour < 0) AND status = 1").Find(&list)
	// 2. 获取 hour > 0 && hour < 4320 && status = 1 的影片信息(只获取最早的一条记录)
	var fr FailureRecord
	db.Mdb.Where("hour > 0 AND hour < 4320 AND status = 1").Order("hour DESC, created_at ASC").First(&fr)
	// 3. 将 fr 添加到 list中
	list = append(list, fr)
	return list
}

// ChangeRecord 修改已完成二次采集的记录状态
func ChangeRecord(fr *FailureRecord, status int) {
	switch {
	case fr.Hour > 168 && fr.Hour < 360:
		db.Mdb.Model(&FailureRecord{}).Where("hour > 168 AND hour < 360 AND created_at >= ?", fr.CreatedAt).Update("status", status)
	case fr.Hour < 0, fr.Hour > 4320:
		db.Mdb.Model(&FailureRecord{}).Where("id = ?", fr.ID).Update("status", status)
	default:
		// 其余范围,暂不处理
		break
	}
}

// RetryRecord 修改重试采集成功的记录
func RetryRecord(id uint, status int64) error {
	// 查询id对应的失败记录
	fr := FindRecordById(id)
	// 将本次更新成功的记录数据状态修改为成功 0
	return db.Mdb.Model(&FailureRecord{}).Where("update_at > ?", fr.UpdatedAt).Update("status", 0).Error

}

// DelDoneRecord 删除已处理的记录信息 -- 逻辑删除
func DelDoneRecord() {
	if err := db.Mdb.Where("status = ?", 0).Delete(&FailureRecord{}).Error; err != nil {
		log.Println("Delete failure record failed:", err)
	}
}

// TruncateRecordTable  截断 record table
func TruncateRecordTable() {
	var s FailureRecord
	err := db.Mdb.Exec(fmt.Sprintf("TRUNCATE Table %s", s.TableName())).Error
	if err != nil {
		log.Println("TRUNCATE TABLE Error: ", err)
	}
}
