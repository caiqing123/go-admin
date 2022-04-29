package role

import (
	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(id string) (role Role) {
	database.DB.Where("id", id).First(&role)
	return
}

func GetBy(field, value string) (role Role) {
	database.DB.Where("? = ?", field, value).First(&role)
	return
}

func All() (roles []Role) {
	database.DB.Find(&roles)
	return
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Role{}).Where(" = ?", field, value).Count(&count)
	return count > 0
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int, q interface{}) (role []Role, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(Role{}))
	//query := database.DB.Model(User{}).Where("id = ?", c.Param("io"))
	paging = paginator.Paginate(
		c,
		db,
		&role,
		app.V1URL(database.TableName(&Role{})),
		perPage,
		"*",
	)
	return
}

//GetRoleMenuId 角色获取菜单id
func GetRoleMenuId(roleId string) []int {
	menuIds := make([]int, 0)
	database.DB.Table("role_menus").Where("role_id= ?", roleId).Find(&menuIds)
	return menuIds
}

func DeleteIds(ids []int, role Role) (rowsAffected int64) {
	result := database.DB.Delete(role, ids)
	return result.RowsAffected
}
