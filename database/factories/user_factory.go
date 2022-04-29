// Package factories 存放工厂方法
package factories

import (
	"api/app/models"
	"api/pkg/hash"
	"api/pkg/helpers"

	"github.com/bxcodec/faker/v3"
)

type User struct {
	models.BaseModel

	Name string `json:"name,omitempty"`

	RoleID       int    `json:"role_id,omitempty"`
	City         string `json:"city,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	NickName     string `json:"nick_name"`

	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"-"`

	models.CommonTimestampsField
}

func MakeUsers(times int) []User {

	var objs []User

	// 设置唯一值
	faker.SetGenerateUniqueValues(true)

	for i := 0; i < times; i++ {
		model := User{
			Name:     faker.Username(),
			Email:    faker.Email(),
			RoleID:   2,
			Phone:    helpers.RandomNumber(11),
			Password: "$2a$14$oPzVkIdwJ8KqY0erYAYQxOuAAlbI/sFIsH0C0R4MPc.3JbWWSuaUe",
		}
		objs = append(objs, model)
	}

	return objs
}

func MakeAdminUsers() []User {

	var objs []User

	// 设置唯一值
	faker.SetGenerateUniqueValues(true)

	model := User{
		Name:     "admin",
		NickName: "admin",
		Email:    faker.Email(),
		RoleID:   1,
		Phone:    helpers.RandomNumber(11),
		Password: hash.BcryptHash("123456"),
	}
	objs = append(objs, model)

	return objs
}
