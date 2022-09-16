package user

import (
	"github.com/gin-gonic/gin"

	"api/pkg/app"
	"api/pkg/database"
	"api/pkg/logger"
	"api/pkg/paginator"
)

// IsEmailExist 判断 Email 已被注册
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// IsPhoneExist 判断手机号已被注册
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// GetByPhone 通过手机号来获取用户
func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone = ?", phone).First(&userModel)
	return
}

// GetByMulti 通过 手机号/Email/用户名 来获取用户
func GetByMulti(loginID string) (userModel User) {
	database.DB.
		Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	return
}

// GetByName GetByMulti 通过 用户名 来获取用户
func GetByName(loginID string) (userModel User) {
	database.DB.
		Where("name = ?", loginID).
		First(&userModel)
	return
}

// Get 通过 ID 获取用户
func Get(idstr int) (userModel User) {
	database.DB.Where("id", idstr).First(&userModel)
	return
}

// IsAdmin id是否超级管理员
func IsAdmin(id int) bool {
	if id == 1 {
		return false
	}
	return true
}

// GetByEmail 通过 Email 来获取用户
func GetByEmail(email string) (userModel User) {
	database.DB.Where("email = ?", email).First(&userModel)
	return
}

// All 获取所有用户数据
func All() (users []User) {
	database.DB.Find(&users)
	return
}

type UserRoleName struct {
	RoleName string `json:"role_name"`
	User
}

func (UserRoleName) TableName() string {
	return "users"
}

// Paginate 分页内容
func Paginate(c *gin.Context, perPage int, q interface{}) (users []UserRoleName, paging paginator.Paging) {
	err := c.Bind(&q)
	if err != nil {
		logger.Error(err.Error())
	}
	db := database.MakeCondition(q, database.DB.Model(User{}))
	//query := database.DB.Model(User{}).Where("id = ?", c.Param("io"))
	paging = paginator.Paginate(
		c,
		db,
		&users,
		app.V1URL(database.TableName(&User{})),
		perPage,
		"users.*,roles.role_name",
	)
	return
}

func DeleteIds(ids []int, user User) (rowsAffected int64) {
	result := database.DB.Delete(user, ids)
	return result.RowsAffected
}
