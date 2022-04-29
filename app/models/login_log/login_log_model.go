//Package login_log 模型
package login_log

import (
	"time"

	"api/app/models"
	"api/pkg/database"
)

type LoginLog struct {
	models.BaseModel

	Username      string    `json:"username"`
	Status        string    `json:"status"`
	Ipaddr        string    `json:"ipaddr"`
	LoginLocation string    `json:"login_location"`
	Browser       string    `json:"browser"`
	Os            string    `json:"os"`
	Platform      string    `json:"platform"`
	LoginTime     time.Time `json:"login_time"`
	Remark        string    `json:"remark"`
	Msg           string    `json:"msg"`

	models.CommonTimestampsField
}

func (loginLog *LoginLog) Create() {
	database.DB.Create(&loginLog)
}

func (loginLog *LoginLog) Save() (rowsAffected int64) {
	result := database.DB.Save(&loginLog)
	return result.RowsAffected
}

func (loginLog *LoginLog) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&loginLog)
	return result.RowsAffected
}
