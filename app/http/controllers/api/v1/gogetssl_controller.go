package v1

import (
	"github.com/gin-gonic/gin"

	"api/pkg/gogetssl"
	"api/pkg/response"
)

type GoGetSllController struct {
	BaseAPIController
}

// Index 数据
func (ctrl *GoGetSllController) Index(c *gin.Context) {
	tag := c.DefaultQuery("type", "")
	var err error
	var data interface{}

	switch tag {
	case "products/ssl":
		data, err = gogetssl.GetProductsSll()
	case "products/all_prices":
		data, err = gogetssl.GetProductsAllPrices()
	default:
		data, err = gogetssl.GetGoGetSll(tag)
	}

	if err != nil {
		response.NormalVerificationError(c, err.Error())
		return
	}

	response.Data(c, data)
}

// Operation 操作
func (ctrl *GoGetSllController) Operation(c *gin.Context) {
	tag := c.DefaultQuery("type", "")
	var err error
	var data interface{}

	switch tag {
	case "tools/csr/decode":
		data, err = gogetssl.PostCsrDecode(c)
	default:
		data, err = gogetssl.PostGoGetSll(tag, c)
	}

	if err != nil {
		response.NormalVerificationError(c, err.Error())
		return
	}

	response.Data(c, data)
}
