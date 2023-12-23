package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/config"
	"server/logic"
	"server/model/system"
	"server/plugin/common/util"
)

// Login 管理员登录接口
func Login(c *gin.Context) {
	var u system.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "数据格式异常!!!",
		})
		return
	}
	if len(u.UserName) <= 0 || len(u.Password) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "用户名和密码信息不能为空!!!",
		})
		return
	}
	token, err := logic.UL.UserLogin(u.UserName, u.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": err.Error(),
		})
		return
	}
	c.Header("new-token", token)
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "登录成功!!!",
	})
	return
}

// Logout 退出登录
func Logout(c *gin.Context) {
	// 获取已登录的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "登录信息异常!!!",
		})
		return
	}
	// 清除redis中存储的对应token
	uc, ok := v.(*system.UserClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "登录信息异常!!!",
		})
		return
	}
	err := system.ClearUserToken(uc.UserID)
	if err != nil {
		log.Println("user logOut err: ", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "logout success!!!",
	})
}

// UserPasswordChange 修改用户密码
func UserPasswordChange(c *gin.Context) {
	// 接收密码修改参数
	var params map[string]string
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "数据格式异常!!!",
		})
		return
	}
	// 校验参数是否存在空值
	if params["password"] == "" || params["newPassword"] == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "原密码和新密码不能为空!!!",
		})
		return
	}
	// 校验新密码是否符合规范
	if err := util.ValidPwd(params["newPassword"]); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("密码格式校验失败: ", err.Error()),
		})
		return
	}
	// 获取已登录的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": "登录信息异常!!!",
		})
		return
	}
	// 从context中获取用户的登录信息
	uc := v.(*system.UserClaims)
	if err := logic.UL.ChangePassword(uc.UserName, params["password"], params["newPassword"]); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  StatusFailed,
			"message": fmt.Sprint("密码修改失败: ", err.Error()),
		})
		return
	}
	// 密码修改成功后不主动使token失效, 以免影响体验
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusOk,
		"message": "密码修改成功",
	})
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
