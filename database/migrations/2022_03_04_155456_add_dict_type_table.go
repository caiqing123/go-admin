package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type DictType struct {
		DictId   uint64 `gorm:"column:dict_id;primaryKey;autoIncrement;"`
		DictName string `gorm:"type:varchar(100);"`
		DictType string `gorm:"type:varchar(60);"`
		Status   int    `gorm:"type:smallint(2);default:2"`
		Remark   string `gorm:"type:varchar(200);"`

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&DictType{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&DictType{})
	}

	migrate.Add("2022_03_04_155456_add_dict_type_table", up, down)
}
