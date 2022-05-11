package opera_log

import (
	"encoding/json"

	"github.com/adjust/rmq/v4"

	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
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
	//todo 失败重试处理
	log.Create()
	if err := delivery.Ack(); err != nil {
		// handle ack error
	}
}
