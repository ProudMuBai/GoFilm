package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/config"
)

var Mdb *gorm.DB

func InitMysql() (err error) {
	// client 相关属性设置
	Mdb, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.MysqlDsn,
		DefaultStringSize:         255,   //string类型字段默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn:   true,  // 用change 重命名列
		SkipInitializeWithVersion: false, // 根据当前Mysql版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_",                          //设置创建表时的前缀
			SingularTable: true, //是否使用 结构体名称作为表名 (关闭自动变复数)
			//NameReplacer:  strings.NewReplacer("spider_", ""), // 替表名和字段中的 Me 为 空
		},
		//Logger: logger.Default.LogMode(logger.Info), //设置日志级别为Info
	})
	return
}
