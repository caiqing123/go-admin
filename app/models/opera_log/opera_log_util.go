package opera_log

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"

	"github.com/gin-gonic/gin"

	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"
)

func Get(idstr string) (operaLog OperaLog) {
	database.DB.Where("id", idstr).First(&operaLog)
	return
}

func GetBy(field, value string) (operaLog OperaLog) {
	database.DB.Where("? = ?", field, value).First(&operaLog)
	return
}

func All() (operaLogs []OperaLog) {
	database.DB.Find(&operaLogs)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(OperaLog{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Clean(operaLogs OperaLog, value string) (rowsAffected int64) {
	result := database.DB.Where("created_at <= ?", value).Delete(&operaLogs)
	return result.RowsAffected
}

func Paginate(c *gin.Context, perPage int, q interface{}) (operaLogs []OperaLog, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(OperaLog{}))
	paging = paginator.Paginate(
		c,
		db,
		&operaLogs,
		app.V1URL(database.TableName(&OperaLog{})),
		perPage,
		"*",
	)
	return
}

type OpLogConsumer struct{}

func (consumer *OpLogConsumer) Consume(delivery rmq.Delivery) {
	var log OperaLog
	if err := json.Unmarshal([]byte(delivery.Payload()), &log); err != nil {
		// handle json error
		if err := delivery.Reject(); err != nil {
			// handle reject error
		}
		return
	}

	//Zinc Search 添加操作日志
	//var zincs = zinc.ZincClient{
	//	ZincClientConfig: &zinc.ZincClientConfig{
	//		ZincHost:     "http://localhost:4080",
	//		ZincUser:     "admin",
	//		ZincPassword: "admin",
	//	},
	//}
	//var data []map[string]interface{}
	//data = append(data, map[string]interface{}{
	//	"index": map[string]interface{}{
	//		"_index": "admin-log",
	//	},
	//}, map[string]interface{}{
	//	"title":          log.Title,
	//	"business_type":  log.BusinessType,
	//	"method":         log.Method,
	//	"RequestMethod":  log.RequestMethod,
	//	"operator_type":  log.OperatorType,
	//	"opera_name":     log.OperaName,
	//	"opera_url":      log.OperaUrl,
	//	"opera_ip":       log.OperaIp,
	//	"opera_time":     log.OperaTime,
	//	"opera_location": log.OperaLocation,
	//	"opera_param":    log.OperaParam,
	//	"status":         log.Status,
	//	"json_result":    log.JsonResult,
	//	"remark":         log.Remark,
	//	"created_at":     log.CreatedAt,
	//	"user_agent":     log.UserAgent,
	//	"latency_time":   log.LatencyTime,
	//})
	//exist, _ := zincs.BulkPutLogDoc(data)
	//if exist {
	//}

	//todo 失败重试处理
	log.Create()
	if err := delivery.Ack(); err != nil {
		// handle ack error
	}
}
