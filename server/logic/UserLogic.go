package logic

import (
	"errors"
	"server/model/system"
	"server/plugin/common/util"
)

type UserLogic struct {
}

var UL *UserLogic

// UserLogin 用户登录
func (ul *UserLogic) UserLogin(account, password string) (token string, err error) {
	// 根据 username 或 email 查询用户信息
	var u *system.User = system.GetUserByNameOrEmail(account)
	// 用户信息不存在则返回提示信息
	if u == nil {
		return "", errors.New(" 用户信息不存在!!!")
	}
	// 校验用户信息, 执行账号密码校验逻辑
	if util.PasswordEncrypt(password, u.Salt) != u.Password {
		return "", errors.New("用户名或密码错误")
	}
	// 密码校验成功后下发token
	token, err = system.GenToken(u.ID, u.UserName)
	err = system.SaveUserToken(token, u.ID)
	return
}

// UserLogout 用户退出登录 注销
func (ul *UserLogic) UserLogout() {
	// 通过用户ID清除Redis中的token信息

}

// ChangePassword 修改密码
func (ul *UserLogic) ChangePassword(account, password, newPassword string) error {
	// 根据 username 或 email 查询用户信息
	var u *system.User = system.GetUserByNameOrEmail(account)
	// 用户信息不存在则返回提示信息
	if u == nil {
		return errors.New(" 用户信息不存在!!!")
	}
	// 首先校验用户的旧密码是否正确
	if util.PasswordEncrypt(password, u.Salt) != u.Password {
		return errors.New("原密码校验失败")
	}
	// 密码校验正确则生成新的用户信息
	newUser := system.User{}
	newUser.ID = u.ID
	// 将新密码进行加密
	newUser.Password = util.PasswordEncrypt(newPassword, u.Salt)
	// 更新用户信息
	system.UpdateUserInfo(newUser)
	return nil
}

func (ul *UserLogic) GetUserInfo(id uint) system.UserInfoVo {
	// 通过用户ID查询对应的用户信息
	u := system.GetUserById(id)
	// 去除user信息中的不必要信息
	var vo = system.UserInfoVo{Id: u.ID, UserName: u.UserName, Email: u.Email, Gender: u.Gender, NickName: u.NickName, Avatar: u.Avatar, Status: u.Status}
	return vo
}
