package job

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (job Job) {
	database.DB.Where("id", idstr).First(&job)
	return
}

func GetBy(field, value string) (job Job) {
	database.DB.Where("? = ?", field, value).First(&job)
	return
}

func GetId(idstr string) (job Job) {
	database.DB.Where("job_id", idstr).First(&job)
	return
}

func GetList() (job []Job) {
	database.DB.Where("status = ?", 2).Find(&job)
	return
}

func RemoveAllEntryID() (rowsAffected int64) {
	result := database.DB.Table("jobs").Where("entry_id > ?", 0).Update("entry_id", 0)
	return result.RowsAffected
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Job{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int, q interface{}) (job []Job, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(Job{}))
	paging = paginator.Paginate(
		c,
		db,
		&job,
		app.V1URL(database.TableName(&Job{})),
		perPage,
		"*",
	)
	return
}

func DeleteIds(ids []int, job Job) (rowsAffected int64) {
	result := database.DB.Delete(job, ids)
	return result.RowsAffected
}
