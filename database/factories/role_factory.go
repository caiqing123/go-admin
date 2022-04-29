package factories

import (
	"api/app/models/role"
)

func MakeRoles() []role.Role {

	var objs []role.Role

	roleModel := role.Role{
		RoleName: "系统管理员",
		Status:   2,
		RoleKey:  "admin",
		Admin:    true,
	}
	objs = append(objs, roleModel)

	return objs
}
