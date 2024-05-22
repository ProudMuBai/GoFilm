package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"server/config"
	"server/model/system"
)

/*

 */

// AuthToken 用户登录Token拦截
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authToken := c.Request.Header.Get("auth-token")
		// 如果没有登录信息则直接清退
		if authToken == "" {
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "用户未授权,请先登录", c)
			c.Abort()
			return
		}
		// 解析token中的信息
		uc, err := system.ParseToken(authToken)
		// 从Redis中获取对应的token是否存在, 如果存在则刷新token
		t := system.GetUserTokenById(uc.UserID)
		// 如果 redis中获取的token为空则登录已过期需重新登录
		if len(t) <= 0 {
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "身份验证信息已失效,请重新登录!!!", c)
			c.Abort()
			return
		}
		// 如果redis中存在对应token, 校验authToken是否与redis中的一致
		if t != authToken {
			// 如果不一致则证明authToken已经失效或在其他地方登录, 则需要重新登录
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "账号在其它设备登录,身份验证信息失效,请重新登录!!!", c)
			c.Abort()
			return
		} else if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
			// 如果token已经过期,且redis中的token与authToken 相同则更新 token
			// 生成新token
			newToken, _ := system.GenToken(uc.UserID, uc.UserName)
			// 将新token同步到redis中
			_ = system.SaveUserToken(newToken, uc.UserID)
			// 解析出新的 UserClaims
			uc, _ = system.ParseToken(newToken)
			c.Header("new-token", newToken)
		}
		// 将UserClaims存放到context中
		c.Set(config.AuthUserClaims, uc)
		c.Next()
	}
}
