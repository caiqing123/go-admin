package config

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (config Config) {
	database.DB.Where("id", idstr).First(&config)
	return
}

// GetByAppConfig 根据条件来获取
func GetByAppConfig(isFrontend string) (config []Config) {
	database.DB.Where("is_frontend = ?", isFrontend).Find(&config)
	return
}

func GetKey(value string) (config Config) {
	database.DB.Where("config_key = ?", value).First(&config)
	return
}

func All() (configs []Config) {
	database.DB.Find(&configs)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Config{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int) (configs []Config, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		database.DB.Model(Config{}),
		&configs,
		app.V1URL(database.TableName(&Config{})),
		perPage,
		"*",
	)
	return
}

func UpdateForSet(c *[]GetSetSysConfigReq) error {
	m := *c
	for _, req := range m {
		var data Config
		if err := database.DB.Where("config_key = ?", req.ConfigKey).
			First(&data).Error; err != nil {
			return err
		}
		if data.ConfigValue != req.ConfigValue {
			data.ConfigValue = req.ConfigValue
			if err := database.DB.Save(&data).Error; err != nil {
				return err
			}
		}
	}
	return nil
}