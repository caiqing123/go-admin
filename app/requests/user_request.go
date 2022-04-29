package requests

import (
	"mime/multipart"

	"api/app/requests/validators"
	"api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type UserUpdateProfileRequest struct {
	NickName     string `valid:"nick_name" json:"nick_name"`
	City         string `valid:"city" json:"city"`
	Introduction string `valid:"introduction" json:"introduction"`
}

func UserUpdateProfile(data interface{}, c *gin.Context) map[string][]string {

	// 查询用户名重复时，过滤掉当前用户 ID
	uid := auth.CurrentUID(c)
	rules := govalidator.MapData{
		"nick_name":    []string{"required", "between:3,20", "not_exists:users,nick_name," + uid},
		"introduction": []string{"min_cn:4", "max_cn:240"},
		"city":         []string{"min_cn:2", "max_cn:20"},
	}

	messages := govalidator.MapData{
		"nick_name": []string{
			"required:用户名为必填项",
			"between:用户名长度需在 3~20 之间",
			"not_exists:昵称已被占用",
		},
		"introduction": []string{
			"min_cn:描述长度需至少 4 个字",
			"max_cn:描述长度不能超过 240 个字",
		},
		"city": []string{
			"min_cn:城市需至少 2 个字",
			"max_cn:城市不能超过 20 个字",
		},
	}
	return validate(data, rules, messages)
}

type UserUpdateEmailRequest struct {
	Email      string `json:"email,omitempty" valid:"email"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

func UserUpdateEmail(data interface{}, c *gin.Context) map[string][]string {

	currentUser := auth.CurrentUser(c)
	rules := govalidator.MapData{
		"email": []string{
			"required", "min:4",
			"max:30",
			"email",
			"not_exists:users,email," + currentUser.GetStringID(),
			"not_in:" + currentUser.Email,
		},
		"verify_code": []string{"required", "digits:6"},
	}
	messages := govalidator.MapData{
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
			"not_exists:Email 已被占用",
			"not_in:新的 Email 与老 Email 一致",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*UserUpdateEmailRequest)
	errs = validators.ValidateVerifyCode(_data.Email, _data.VerifyCode, errs)

	return errs
}

type UserUpdatePhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

func UserUpdatePhone(data interface{}, c *gin.Context) map[string][]string {

	currentUser := auth.CurrentUser(c)

	rules := govalidator.MapData{
		"phone": []string{
			"required",
			"digits:11",
			"not_exists:users,phone," + currentUser.GetStringID(),
			"not_in:" + currentUser.Phone,
		},
		"verify_code": []string{"required", "digits:6"},
	}
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
			"not_exists:手机号已被占用",
			"not_in:新的手机与老手机号一致",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*UserUpdatePhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

type UserUpdatePasswordRequest struct {
	Password           string `valid:"password" json:"password,omitempty"`
	NewPassword        string `valid:"new_password" json:"new_password,omitempty"`
	NewPasswordConfirm string `valid:"new_password_confirm" json:"new_password_confirm,omitempty"`
}

func UserUpdatePassword(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"password":     []string{"required", "min:6"},
		"new_password": []string{"required", "min:6", "max:20"},
	}
	messages := govalidator.MapData{
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"new_password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
			"max:密码长度需小于 20",
		},
	}

	// 确保 comfirm 密码正确
	errs := validate(data, rules, messages)
	_data := data.(*UserUpdatePasswordRequest)
	errs = validators.ValidatePasswordConfirm(_data.NewPassword, _data.NewPasswordConfirm, errs)

	return errs
}

type UserUpdateAvatarRequest struct {
	Avatar *multipart.FileHeader `valid:"avatar" form:"avatar"`
}

func UserUpdateAvatar(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		// size 的单位为 bytes
		// - 1024 bytes 为 1kb
		// - 1048576 bytes 为 1mb
		// - 20971520 bytes 为 20mb
		"file:avatar": []string{"ext:png,jpg,jpeg", "size:20971520", "required"},
	}
	messages := govalidator.MapData{
		"file:avatar": []string{
			"ext:ext头像只能上传 png, jpg, jpeg 任意一种的图片",
			"size:头像文件最大不能超过 20MB",
			"required:必须上传图片",
		},
	}

	return validateFile(c, data, rules, messages)
}

type UserDeleteRequest struct {
	Ids []int `valid:"ids" form:"ids"`
}

func UserDelete(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"ids": []string{"required"},
	}
	messages := govalidator.MapData{
		"ids": []string{
			"required:id为空",
		},
	}
	errs := validate(data, rules, messages)

	_data := data.(*UserDeleteRequest)
	errs = validators.ValidateInAdminArray(_data.Ids, 1, errs)

	return errs
}

// UserPaginationRequest 分页查询
type UserPaginationRequest struct {
	Sort    string `valid:"sort" form:"sort" search:"-"`
	Order   string `valid:"order" form:"order" search:"-"`
	PerPage string `valid:"per_page" form:"per_page" search:"-"`

	Role   `search:"type:left;on:id:role_id;table:users;join:roles"`
	Phone  string `form:"phone" search:"type:contains;column:phone;table:users"`
	Name   string `form:"name" search:"type:contains;column:name;table:users"`
	Status string `form:"status" search:"type:exact;column:status;table:users"`
}

type Role struct {
	RoleName string `search:"type:contains;column:role_name;table:roles" form:"role_name"`
}

// UserRequest AddRequest
type UserRequest struct {
	Phone        string `json:"phone,omitempty" valid:"phone" form:"phone"`
	Email        string `json:"email,omitempty" valid:"email" form:"email"`
	Name         string `valid:"name" json:"name" form:"name"`
	Password     string `valid:"password" json:"password,omitempty" form:"password"`
	Status       string `valid:"status" json:"status" form:"status"`
	RoleID       int    `valid:"role_id" json:"role_id" form:"role_id"`
	NickName     string `json:"nick_name" form:"nick_name"`
	Introduction string `json:"introduction" form:"introduction"`
	Id           string `valid:"id" json:"id" form:"id"`
}

func UserSave(data interface{}, c *gin.Context) map[string][]string {
	_data := data.(*UserRequest)
	rules := govalidator.MapData{
		"phone":    []string{"required", "digits:11", "not_exists:users,phone," + _data.Id},
		"name":     []string{"required", "alpha_num", "between:3,20", "not_exists:users,name," + _data.Id},
		"password": []string{"min:6"},
		"email":    []string{"required", "min:4", "max:30", "email", "not_exists:users,email," + _data.Id},
		"role_id":  []string{"required", "exists:roles,id"},
		"status":   []string{"required", "in:1,2"},
		"id":       []string{"exists:users,id"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
			"not_exists:phone 已被占用",
		},
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
			"not_exists:name 已被占用",
		},
		"password": []string{
			"min:密码长度需大于 6",
		},
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
			"not_exists:email 已被占用",
		},
		"status": []string{
			"required:状态为必填项",
			"in:状态格式错误",
		},
		"role_id": []string{
			"required:角色为必填项",
			"exists:角色不存在",
		},
		"id": []string{
			"exists:用户不存在",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
