// Package dict_data 模型
package dict_data

import (
	"api/app/models"
	"api/pkg/database"
)

type DictData struct {
	DictCode  int    `gorm:"column:dict_code;primaryKey;autoIncrement;" json:"dict_code,omitempty"`
	DictSort  int    `json:"dict_sort" title:"排序"`
	DictLabel string `json:"dict_label" title:"数据标签"`
	DictValue string `json:"dict_value" title:"数据键值"`
	DictType  string `json:"dict_type" title:"字典类型"`
	Status    int    `json:"status" title:"状态"`
	Remark    string `json:"remark" title:"备注"`

	models.CommonTimestampsField
}

type DictDataGetAllResp struct {
	DictLabel string `json:"label"`
	DictValue string `json:"value"`
}

func (dictData *DictData) Create() {
	database.DB.Create(&dictData)
}

func (dictData *DictData) Save() (rowsAffected int64) {
	result := database.DB.Save(&dictData)
	return result.RowsAffected
}

func (dictData *DictData) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&dictData)
	return result.RowsAffected
}
