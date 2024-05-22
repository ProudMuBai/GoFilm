package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"server/config"
	"server/logic"
	"server/model/system"
	"server/plugin/common/util"
)

// Login 管理员登录接口
func Login(c *gin.Context) {
	var u system.User
	if err := c.ShouldBindJSON(&u); err != nil {
		system.Failed("登录信息异常!!!", c)
		return
	}
	if len(u.UserName) <= 0 || len(u.Password) <= 0 {
		system.Failed("用户名和密码信息不能为空", c)
		return
	}
	token, err := logic.UL.UserLogin(u.UserName, u.Password)
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}
	c.Header("new-token", token)
	system.SuccessOnlyMsg("登录成功!!!", c)
}

// Logout 退出登录
func Logout(c *gin.Context) {
	// 获取已登录的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("请求失败,登录信息获取异常!!!", c)
		return
	}
	// 清除redis中存储的对应token
	uc, ok := v.(*system.UserClaims)
	if !ok {
		system.Failed("注销失败, 身份信息格式化异常!!!", c)
		return
	}
	err := system.ClearUserToken(uc.UserID)
	if err != nil {
		log.Println("user logOut err: ", err)
	}
	system.SuccessOnlyMsg("已退出登录!!!", c)
}

// UserPasswordChange 修改用户密码
func UserPasswordChange(c *gin.Context) {
	// 接收密码修改参数
	var params map[string]string
	if err := c.ShouldBindJSON(&params); err != nil {
		system.Failed("参数校验失败!!!", c)
		return
	}
	// 校验参数是否存在空值
	if params["password"] == "" || params["newPassword"] == "" {
		system.Failed("密码不能为空!!!", c)
		return
	}
	// 校验新密码是否符合规范
	if err := util.ValidPwd(params["newPassword"]); err != nil {
		system.Failed(fmt.Sprint("密码格式校验失败: ", err.Error()), c)
		return
	}
	// 获取已登录的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("操作失败,登录信息异常!!!", c)
		return
	}
	// 从context中获取用户的登录信息
	uc := v.(*system.UserClaims)
	if err := logic.UL.ChangePassword(uc.UserName, params["password"], params["newPassword"]); err != nil {
		system.Failed(fmt.Sprint("密码修改失败: ", err.Error()), c)
		return
	}
	// 密码修改成功后不主动使token失效, 以免影响体验
	system.SuccessOnlyMsg("密码修改成功", c)
}

func UserInfo(c *gin.Context) {
	// 从context中获取用户的相关信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("用户信息获取失败, 未获取到用户授权信息", c)
		return
	}
	uc, ok := v.(*system.UserClaims)
	if !ok {
		system.Failed("用户信息获取失败, 户授权信息异常", c)
		return
	}
	// 通过用户ID获取用户基本信息
	info := logic.UL.GetUserInfo(uc.UserID)
	system.Success(info, "成功获取用户信息", c)
}
