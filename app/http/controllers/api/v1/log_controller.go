package v1

import (
	"time"

	"api/app/models/login_log"
	"api/app/models/opera_log"
	"api/app/requests"
	"api/pkg/auth"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	BaseAPIController
}

// LoginLog 获取登录log
func (ctrl *LogController) LoginLog(c *gin.Context) {
	request := requests.LoginLogPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.LoginLogPagination); !ok {
		return
	}
	data, pager := login_log.Paginate(c, 10, request)
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

// OperaLog 获取操作log
func (ctrl *LogController) OperaLog(c *gin.Context) {
	request := requests.OperaLogPaginationRequest{}
	if ok := requests.Validate(c, &request, requests.OperaLogPagination); !ok {
		return
	}
	data, pager := opera_log.Paginate(c, 10, request)
	response.JSON(c, gin.H{
		"data":  data,
		"pager": pager,
	})
}

func (ctrl *LogController) Clean(c *gin.Context) {
	if auth.CurrentUID(c) != "1" {
		response.NormalVerificationError(c, "超级管理员才能操作")
		return
	}

	oneTime := time.Now().AddDate(0, -1, 0).Format("2006-01-02 15:04:05")

	rowsAffected := opera_log.Clean(opera_log.OperaLog{}, oneTime)
	if rowsAffected > 0 {
		response.Success(c)
		return
	}
	response.Abort500(c, "删除失败，请稍后尝试~")
}
