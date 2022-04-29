// Package middlewares Gin 中间件
package middlewares

import (
	"fmt"

	"github.com/chenhg5/collection"

	"api/app/models/menu"
	"api/app/models/role"
	"api/app/models/user"
	"api/pkg/jwt"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetBool("ssl_white_list") {
			c.Next()
			return
		}
		// 从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
		claims, err := jwt.NewJWT().ParserToken(c)
		// JWT 解析失败，有错误发生
		if err != nil {
			response.Unauthorized(c, "认证失败")
			return
		}

		// JWT 解析成功，设置用户信息
		userModel := user.Get(claims.UserID)
		if userModel.ID == 0 {
			response.Unauthorized(c, "找不到对应用户，用户可能已删除")
			return
		}

		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		c.Set("current_user_id", userModel.GetStringID())
		c.Set("current_user_name", userModel.Name)
		c.Set("current_user", userModel)

		//权限验证
		apiUrl := fmt.Sprintf("%s-%s", c.Request.Method, c.Request.URL)
		mid := menu.GetBy(apiUrl).ID
		if mid > 0 && userModel.ID != 1 {
			mids := role.GetRoleMenuId(userModel.RoleID)
			if !collection.Collect(mids).Contains(mid) {
				response.Abort403(c, "无权限访问")
				return
			}
		}

		c.Next()
	}
}
