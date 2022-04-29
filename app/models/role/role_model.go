//Package role 模型
package role

import (
	"api/app/models"
	"api/pkg/database"
)

type Role struct {
	models.BaseModel

	RoleName string `json:"role_name,omitempty"`
	Status   int    `json:"status,omitempty"`
	RoleKey  string `json:"role_key,omitempty"`
	Admin    bool   `json:"admin,omitempty"`
	RoleSort int    `json:"role_sort"`
	Remark   string `json:"remark"`

	models.CommonTimestampsField
}

func (role *Role) Create() {
	database.DB.Create(&role)
}

func (role *Role) Save() (rowsAffected int64) {
	result := database.DB.Save(&role)
	return result.RowsAffected
}

func (role *Role) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&role)
	return result.RowsAffected
}
