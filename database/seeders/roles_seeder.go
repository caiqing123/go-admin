package seeders

import (
	"fmt"

	"api/database/factories"
	"api/pkg/console"
	"api/pkg/logger"
	"api/pkg/seed"

	"gorm.io/gorm"
)

func init() {

	seed.Add("SeedRolesTable", func(db *gorm.DB) {

		roles := factories.MakeRoles()

		result := db.Table("roles").Create(&roles)

		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
