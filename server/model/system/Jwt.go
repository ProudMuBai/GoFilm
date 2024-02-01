package system

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"server/config"
	"server/plugin/common/util"
	"server/plugin/db"
	"time"
)

type UserClaims struct {
	UserID   uint   `json:"userID"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// GenToken 生成token
func GenToken(userId uint, userName string) (string, error) {
	uc := UserClaims{
		UserID:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userName,
			Audience:  jwt.ClaimStrings{"Auth_All"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AuthTokenExpires * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        util.GenerateSalt(),
		},
	}
	priKey, err := util.ParsePriKeyBytes([]byte(config.PrivateKey))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, uc).SignedString(priKey)
	return token, err
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		pub, err := util.ParsePubKeyBytes([]byte(config.PublicKey))
		if err != nil {
			return nil, err
		}
		return pub, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, _ := token.Claims.(*UserClaims)
			return claims, err
		}
		return nil, err
	}
	// 验证token是否有效
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	// 解析token中的claims内容
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid claim type error")
	}
	return claims, err
}

// RefreshToken 刷新 token

// SaveUserToken 将用户登录成功后的token字符串存放到redis中
func SaveUserToken(token string, userId uint) error {
	// 设置redis中token的过期时间为 token过期时间后的7天
	return db.Rdb.Set(db.Cxt, fmt.Sprintf(config.UserTokenKey, userId), token, (config.AuthTokenExpires+7*24)*time.Hour).Err()
}

// GetUserTokenById 从redis中获取指定userId对应的token
func GetUserTokenById(userId uint) string {
	token, err := db.Rdb.Get(db.Cxt, fmt.Sprintf(config.UserTokenKey, userId)).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	return token
}

// ClearUserToken 清楚指定id的用户的登录信息
func ClearUserToken(userId uint) error {
	return db.Rdb.Del(db.Cxt, fmt.Sprintf(config.UserTokenKey, userId)).Err()
}
