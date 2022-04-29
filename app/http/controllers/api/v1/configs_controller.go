package v1

import (
	"fmt"

	"api/app/models/config"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigsController struct {
	BaseAPIController
}

// GetApp 获取配置
func (ctrl *ConfigsController) GetApp(c *gin.Context) {
	isFrontend := c.DefaultQuery("type", "0")
	data := config.GetByAppConfig(isFrontend)
	mp := make(map[string]string)
	for i := 0; i < len(data); i++ {
		key := data[i].ConfigKey
		if key != "" {
			mp[key] = data[i].ConfigValue
		}
	}
	response.Data(c, mp)
}

// ConfigKey 根据key获取配置
func (ctrl *ConfigsController) ConfigKey(c *gin.Context) {
	data := config.GetKey(c.DefaultQuery("configKey", ""))
	fmt.Println(data)
	mp := make(map[string]string)
	mp["ConfigKey"] = data.ConfigKey
	mp["configValue"] = data.ConfigValue
	response.Data(c, mp)
}

// GetAll 获取配置
func (ctrl *ConfigsController) GetAll(c *gin.Context) {
	data := config.All()
	mp := make(map[string]string)
	for i := 0; i < len(data); i++ {
		key := data[i].ConfigKey
		if key != "" {
			mp[key] = data[i].ConfigValue
		}
	}
	response.Data(c, mp)
}

//SetConfig 设置
func (ctrl *ConfigsController) SetConfig(c *gin.Context) {
	request := make([]config.GetSetSysConfigReq, 0)
	// 1. 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(&request); err != nil {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return
	}
	err := config.UpdateForSet(&request)
	if err != nil {
		response.NormalVerificationError(c, "设置失败")
		return
	}
	response.Success(c)
}
