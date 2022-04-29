package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type DictData struct {
		DictCode  uint64 `gorm:"column:dict_code;primaryKey;autoIncrement;"`
		DictSort  int    `gorm:"type:int(11)"`
		DictLabel string `gorm:"type:varchar(100);"`
		DictValue string `gorm:"type:varchar(100);"`
		DictType  string `gorm:"type:varchar(60);"`
		Status    int    `gorm:"type:smallint(2);default:2"`
		Remark    string `gorm:"type:varchar(200);"`

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&DictData{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&DictData{})
	}
	migrate.Add("2022_03_04_155535_add_dict_data_table", up, down)
}
