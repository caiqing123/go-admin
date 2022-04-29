//Package dict_type 模型
package dict_type

import (
	"api/app/models"
	"api/pkg/database"
)

type DictType struct {
	DictId   int    `gorm:"column:dict_id;primaryKey;autoIncrement;" json:"id,omitempty"`
	DictName string `json:"dict_name"`
	DictType string `json:"dict_type"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`

	models.CommonTimestampsField
}

func (dictType *DictType) Create() {
	database.DB.Create(&dictType)
}

func (dictType *DictType) Save() (rowsAffected int64) {
	result := database.DB.Save(&dictType)
	return result.RowsAffected
}

func (dictType *DictType) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&dictType)
	return result.RowsAffected
}
