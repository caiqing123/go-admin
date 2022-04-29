//Package opera_log 模型
package opera_log

import (
	"time"

	"api/app/models"
	"api/pkg/database"
)

type OperaLog struct {
	models.BaseModel
	Title         string    `json:"title"`
	BusinessType  string    `json:"business_type"`
	Method        string    `json:"method"`
	RequestMethod string    `json:"request_method"`
	OperatorType  string    `json:"operator_type"`
	OperaName     string    `json:"opera_name"`
	OperaUrl      string    `json:"opera_url"`
	OperaIp       string    `json:"opera_ip"`
	OperaLocation string    `json:"opera_location"`
	OperaParam    string    `json:"opera_param"`
	Status        string    `json:"status"`
	OperaTime     time.Time `json:"opera_time"`
	JsonResult    string    `json:"json_result"`
	Remark        string    `json:"remark"`
	LatencyTime   string    `json:"latency_time"`
	UserAgent     string    `json:"user_agent"`

	models.CommonTimestampsField
}

func (operaLog *OperaLog) Create() {
	database.DB.Create(&operaLog)
}

func (operaLog *OperaLog) Save() (rowsAffected int64) {
	result := database.DB.Save(&operaLog)
	return result.RowsAffected
}

func (operaLog *OperaLog) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&operaLog)
	return result.RowsAffected
}
