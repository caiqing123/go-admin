package login_log

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (loginLog LoginLog) {
	database.DB.Where("id", idstr).First(&loginLog)
	return
}

func GetBy(field, value string) (loginLog LoginLog) {
	database.DB.Where("? = ?", field, value).First(&loginLog)
	return
}

func All() (loginLogs []LoginLog) {
	database.DB.Find(&loginLogs)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(LoginLog{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int, q interface{}) (loginLogs []LoginLog, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(LoginLog{}))
	paging = paginator.Paginate(
		c,
		db,
		&loginLogs,
		app.V1URL(database.TableName(&LoginLog{})),
		perPage,
		"*",
	)
	return
}
