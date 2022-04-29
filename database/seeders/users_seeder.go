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

	// 添加 Seeder
	seed.Add("SeedUsersTable", func(db *gorm.DB) {

		// 创建 管理 用户对象
		users := factories.MakeAdminUsers()

		// 批量创建用户（注意批量创建不会调用模型钩子）
		result := db.Table("users").Create(&users)

		// 记录错误
		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		// 打印运行情况
		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.RowsAffected))
	})
}
