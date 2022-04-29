package dict_type

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idstr string) (dictType DictType) {
	database.DB.Where("dict_id", idstr).First(&dictType)
	return
}

func GetBy(field, value string) (dictType DictType) {
	database.DB.Where("? = ?", field, value).First(&dictType)
	return
}

func All() (dictTypes []DictType) {
	database.DB.Find(&dictTypes)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(DictType{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

func Paginate(c *gin.Context, perPage int, q interface{}) (dictTypes []DictType, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(DictType{}))
	//query := database.DB.Model(User{}).Where("id = ?", c.Param("io"))
	paging = paginator.Paginate(
		c,
		db,
		&dictTypes,
		app.V1URL(database.TableName(&DictType{})),
		perPage,
		"*",
	)
	return
}

func DeleteIds(ids []int, dictType DictType) (rowsAffected int64) {
	result := database.DB.Delete(dictType, ids)
	return result.RowsAffected
}
