package migrations

import (
	"database/sql"

	"api/app/models"
	"api/pkg/migrate"

	"gorm.io/gorm"
)

func init() {

	type User struct {
		models.BaseModel

		Name     string `gorm:"type:varchar(255);not null;index"`
		RoleId   int    `gorm:"column:role_id;type:int(11);not null;index"`
		Status   int    `gorm:"type:smallint(2);default:2"`
		Email    string `gorm:"type:varchar(255);index;default:null"`
		Phone    string `gorm:"type:varchar(20);index;default:null"`
		NickName string `gorm:"type:varchar(50);default:null"`
		Password string `gorm:"type:varchar(255)"`

		models.CommonTimestampsField
	}

	up := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.AutoMigrate(&User{})
	}

	down := func(migrator gorm.Migrator, DB *sql.DB) {
		_ = migrator.DropTable(&User{})
	}

	migrate.Add("2022_01_11_212406_add_users_table", up, down)
}
