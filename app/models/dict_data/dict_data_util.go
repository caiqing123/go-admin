package dict_data

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

// Get find 查询多个值必须DictDat[]
func Get(dictType string) (dictData []DictData) {
	database.DB.Where("dict_type", dictType).Find(&dictData)
	return
}

func GetId(idstr string) (dictType DictData) {
	database.DB.Where("dict_code", idstr).First(&dictType)
	return
}

func GetBy(field, value string) (dictData DictData) {
	database.DB.Where("? = ?", field, value).First(&dictData)
	return
}

func All() (dictData []DictData) {
	database.DB.Find(&dictData)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(DictData{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int, q interface{}) (dictData []DictData, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(DictData{}))
	paging = paginator.Paginate(
		c,
		db,
		&dictData,
		app.V1URL(database.TableName(&DictData{})),
		perPage,
		"*",
	)
	return
}

func DeleteIds(ids []int, dictData DictData) (rowsAffected int64) {
	result := database.DB.Delete(dictData, ids)
	return result.RowsAffected
}
