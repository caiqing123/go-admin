package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type Config struct {
		models.BaseModel

		ConfigName  string `gorm:"type:varchar(128);not null"`
		ConfigKey   string `gorm:"type:varchar(128)"`
		ConfigValue string `gorm:"type:varchar(255);not null"`
		ConfigType  string `gorm:"type:varchar(30)"`
		IsFrontend  int    `gorm:"type:tinyint(1);default:0"`
		Remark      string `gorm:"type:varchar(255)"`

		models.CommonTimestampsField
		models.DeletedTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&Config{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&Config{})
	}

	migrate.Add("2022_03_02_171949_add_config_table", up, down)
}
