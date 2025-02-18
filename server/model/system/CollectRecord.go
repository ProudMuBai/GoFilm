package system

import (
	"gorm.io/gorm"
	"log"
	"server/config"
	"server/plugin/db"
)

// FailureRecord 失败采集记录信息机构体
type FailureRecord struct {
	gorm.Model
	OriginId    string       `json:"originId"`    // 采集站唯一ID
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
	if err := db.Mdb.Create(&fl).Error; err != nil {
		log.Println("Add failure record failed:", err)
	}
}

// FailureRecordList 获取所有的采集失效记录
func FailureRecordList(page *Page) []FailureRecord {
	var count int64
	db.Mdb.Model(&FailureRecord{}).Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
	// 获取分页查询的数据
	var list []FailureRecord
	if err := db.Mdb.Limit(page.PageSize).Offset((page.Current - 1) * page.PageSize).Find(&list).Error; err != nil {
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
	db.Mdb.First(fr)
	return &fr
}

// RetryRecord 修改重试采集成功的记录
func RetryRecord(id uint, status int64) error {
	// 查询id对应的失败记录
	fr := FindRecordById(id)
	// 将本次更新成功的记录数据状态修改为成功 0
	return db.Mdb.Model(&FailureRecord{}).Where("update_at > ?", fr.UpdatedAt).Update("status", 0).Error

}
