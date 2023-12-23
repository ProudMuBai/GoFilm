package system

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
)

type User struct {
	gorm.Model
	UserName string `json:"userName"` // 用户名
	Password string `json:"password"` // 密码
	Salt     string `json:"salt"`     // 盐值
	Email    string `json:"email"`    // 邮箱
	Gender   int    `json:"gender"`   // 性别
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Status   int    `json:"status"`   // 状态
	Reserve1 string `json:"reserve1"` // 预留字段 3
	Reserve2 string `json:"reserve2"` // 预留字段 2
	Reserve3 string `json:"reserve3"` // 预留字段 1
	//LastLongTime time.Time `json:"LastLongTime"` // 最后登录时间
}

// TableName 设置user表的表名
func (u *User) TableName() string {
	return config.UserTableName
}

// CreateUserTable 创建存储检索信息的数据表
func CreateUserTable() {
	var u = &User{}
	// 如果不存在则创建表 并设置自增ID初始值为10000
	if !ExistUserTable() {
		err := db.Mdb.AutoMigrate(u)
		db.Mdb.Exec(fmt.Sprintf("alter table %s auto_Increment=%d", u.TableName(), config.UserIdInitialVal))
		if err != nil {
			log.Println("Create Table SearchInfo Failed: ", err)
		}
	}
}

// ExistUserTable 判断表中是否存在User表
func ExistUserTable() bool {
	return db.Mdb.Migrator().HasTable(&User{})
}

// InitAdminAccount 初始化admin用户密码
func InitAdminAccount() {
	// 先查询是否已经存在admin用户信息, 存在则直接退出
	user := GetUserByNameOrEmail("admin")
	if user != nil {
		return
	}
	// 不存在管理员用户则进行初始化创建
	u := &User{
		UserName: "admin",
		Password: "admin",
		Salt:     util.GenerateSalt(),
		Email:    "administrator@gmail.com",
		Gender:   2,
		NickName: "Zero",
		Avatar:   "empty",
		Status:   0,
	}

	u.Password = util.PasswordEncrypt(u.Password, u.Salt)
	db.Mdb.Create(u)
}

// GetUserByNameOrEmail 查询 username || email 对应的账户信息
func GetUserByNameOrEmail(userName string) *User {
	var u *User
	if err := db.Mdb.Where("user_name = ? OR email = ?", userName, userName).First(&u).Error; err != nil {
		log.Println(err)
		return nil
	}
	return u
}

func GetUserById(id uint) User {
	var user = User{Model: gorm.Model{ID: id}}
	db.Mdb.First(&user)
	return user
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(u User) {
	// 值更新允许修改的部分字段, 零值会在更新时被自动忽略
	db.Mdb.Model(&u).Updates(User{Password: u.Password, Email: u.Email, NickName: u.NickName, Status: u.Status})
}
