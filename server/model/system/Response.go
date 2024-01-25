package system

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

/*
	对 http response 做简单的封装
*/

const (
	SUCCESS = 0
	FAILED  = -1
)

// Response http返回数据结构体
type Response struct {
	Code int    `json:"code"` // 状态 ok | failed
	Data any    `json:"data"` // 数据
	Msg  string `json:"msg"`  // 提示信息
	//Count  int    `json:"count"`  // 内容长度
}

// PagingData 分页基本数据通用格式
type PagingData struct {
	List   []any `json:"list"`
	Paging Page  `json:"paging"`
}

// Page 分页信息结构体
type Page struct {
	PageSize  int `json:"pageSize"`  // 每页大小
	Current   int `json:"current"`   // 当前页
	PageCount int `json:"pageCount"` // 总页数
	Total     int `json:"total"`     // 总记录数
	//List      []interface{} `json:"list"`      // 数据
}

// Result 构建response返回数据结构
func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// Success 成功响应 数据 + 成功提示
func Success(data any, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// SuccessOnlyMsg 成功响应, 只返回成功信息
func SuccessOnlyMsg(message string, c *gin.Context) {
	Result(SUCCESS, nil, message, c)
}

// Failed 响应失败 只返回错误信息
func Failed(message string, c *gin.Context) {
	Result(FAILED, nil, message, c)
}

// FailedWithData 返回错误信息以及必要数据
func FailedWithData(data any, message string, c *gin.Context) {
	Result(FAILED, data, message, c)
}

// CustomResult 自定义返回状态以及相关数据, 用于异常返回情况
func CustomResult(statusCode int, code int, data any, msg string, c *gin.Context) {
	c.JSON(statusCode, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// ExceptionResult 异常状态返回
func ExceptionResult(statusCode int, message string, c *gin.Context) {
	CustomResult(statusCode, SUCCESS, nil, message, c)
}

// GetPage 获取分页相关数据
func GetPage(db *gorm.DB, page *Page) {
	var count int64
	db.Count(&count)
	page.Total = int(count)
	page.PageCount = int((page.Total + page.PageSize - 1) / page.PageSize)
}
