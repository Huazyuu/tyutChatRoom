package res

import (
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorCode int

const (
	SettingsError ErrorCode = iota + 1001 // 系统错误
	ArgumentError                         // 参数错误
	SystemError                           // 系统错误
)
const (
	Success = 0
	Error   = 7
)

var (
	ErrorMap = map[ErrorCode]string{
		SettingsError: "系统错误",
		ArgumentError: "参数错误",
	}
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}
type ListResponse[T any] struct {
	Count int64 `json:"count"`
	List  T     `json:"list"`
}

func OkWithData(data any, c *gin.Context) {
	result(Success, data, "成功", c)
}
func OkWithMessage(msg string, c *gin.Context) {
	result(Success, map[string]any{}, msg, c)
}
func OkWithList(list any, count int64, c *gin.Context) {
	OkWithData(ListResponse[any]{
		List:  list,
		Count: count,
	}, c)
}

func FailWithMessage(msg string, c *gin.Context) {
	result(Error, map[string]any{}, msg, c)
}
func FailWithCode(code ErrorCode, c *gin.Context) {
	msg, ok := ErrorMap[code]
	if ok {
		result(int(code), map[string]any{}, msg, c)
		return
	}
	result(Error, map[string]any{}, "未知错误", c)
}

func FailWithError(err error, obj any, c *gin.Context) {
	msg := utils.GetValidMsg(err, obj)
	FailWithMessage(msg, c)
}

func result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}
