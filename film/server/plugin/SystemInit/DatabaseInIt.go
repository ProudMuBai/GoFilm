package SystemInit

import "server/model/system"

// TableInIt 初始化 mysql 数据库相关数据
func TableInIt() {
	// 创建 User Table
	system.CreateUserTable()
	// 初始化管理员账户
	system.InitAdminAccount()
	// 创建 Search Table
	system.CreateSearchTable()
	// 创建图片信息管理表
	system.CreateFileTable()
}
