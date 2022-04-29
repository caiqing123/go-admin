//Package config 模型
package config

import (
	"api/app/models"
	"api/pkg/database"
)

type Config struct {
	models.BaseModel

	ConfigName  string `json:"config_name"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	IsFrontend  int    `json:"is_frontend"`
	Remark      string `json:"remark"`

	models.CommonTimestampsField
}

// GetSetSysConfigReq 增、改使用的结构体
type GetSetSysConfigReq struct {
	ConfigKey   string `json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
}

func (config *Config) Create() {
	database.DB.Create(&config)
}

func (config *Config) Save() (rowsAffected int64) {
	result := database.DB.Save(&config)
	return result.RowsAffected
}

func (config *Config) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&config)
	return result.RowsAffected
}
