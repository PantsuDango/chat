package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
)

const (
	Success              = 0    // 请求正常
	IllegalRequestParams = -100 // 请求参数错误
	OperationDBError     = -200 // 操作数据库失败
	IpAnalysisError      = -300 // IP地址解析错误
)

const (
	SuccessMessage          = "成功"
	OperationDBErrMessage   = "数据库错误"
	IpAnalysisErrMessage    = "IP地址解析错误"
	RequestParamsErrMessage = "请求参数错误"
)

const (
	MessageTypeFirst    = "First"
	MessageTypeOption   = "Option"
	MessageTypeManual   = "Manual"
	MessageTypeKeyword  = "Keyword"
	MessageTypeCustomer = "Customer"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

// 请求成功的返回
func JSONSuccess(ctx *gin.Context, result interface{}) {

	ctx.JSON(http.StatusOK, gin.H{
		"Code":      Success,
		"Status":    SuccessMessage,
		"RequestID": uuid.NewV4(),
		"Result":    result,
	})
}

// 请求失败的返回
func JSONFail(ctx *gin.Context, retCode int, errorMsg string) {

	ctx.JSON(http.StatusOK, gin.H{
		"Code":      retCode,
		"Status":    "Fail",
		"RequestID": uuid.NewV4(),
		"ErrorMsg":  errorMsg,
	})
}

// 校验IP地址
func CheckIP(ip string) (err error) {

	// 校验IP地址
	regex := "((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}"
	ok, err := regexp.MatchString(regex, ip)
	if err != nil || !ok {
		err = errors.New(IpAnalysisErrMessage)
		return
	}

	return
}
