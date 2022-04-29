package gogetssl

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"api/pkg/cache"
	"api/pkg/http"
	"api/pkg/logger"
)

const (
	api  = "https://my.gogetssl.com/api/"
	user = "jiobxn@gmail.com"
	pass = "h2JCL5B8orwu7N0"
)

var token = ""

type CrsResult struct {
	CsrCode string `json:"csr_code"`
	Success bool   `json:"success"`
	CsrKey  string `json:"csr_key"`

	Message     string `json:"message"`
	Error       bool   `json:"error"`
	Description string `json:"description"`
}

type ProductsSLL struct {
	CrsResult
	Products interface{} `json:"products"`
}

func VerificationToken() {
	//token = cache.Cache.Store.Get("ssl_token")
	token = "37eaeeb577f48f3733c3799d64fa939f61e79d19"
	if token == "" {
		token = getAuth()
		cache.Cache.Store.Set("ssl_token", token, 364*(time.Hour*24))
	}
}

func GetProductsSll() (result ProductsSLL, err error) {
	req, err := http.Get(api + fmt.Sprintf("products/ssl?auth_key=%s", token))
	err = json.Unmarshal([]byte(req), &result) //json字符串反解析为结构体.结构体需要提前定义
	if result.Success != true {
		err = fmt.Errorf(result.Message)
	}
	return
}

func GetProductsAllPrices() (result map[string]interface{}, err error) {
	req, err := http.Get(api + fmt.Sprintf("products/all_prices?auth_key=%s", token))
	err = json.Unmarshal([]byte(req), &result) //json字符串反解析为字典
	if result["success"] != true {
		err = fmt.Errorf(result["message"].(string))
	}
	return
}

func GetGoGetSll(name string) (result map[string]interface{}, err error) {
	req, err := http.Get(api + fmt.Sprintf("%s?auth_key=%s", name, token))
	err = json.Unmarshal([]byte(req), &result) //json字符串反解析为字典
	if result["success"] != true {
		err = fmt.Errorf(result["message"].(string))
	}
	return
}

func PostGoGetSll(name string, c *gin.Context) (result map[string]interface{}, err error) {
	err = c.Request.ParseForm()
	req, err := http.Post(api+fmt.Sprintf("%s?auth_key=%s", name, token), strings.NewReader(c.Request.PostForm.Encode()), "application/x-www-form-urlencoded; charset=UTF-8")
	err = json.Unmarshal(req, &result) //json字符串反解析为字典
	if result["success"] != true {
		err = fmt.Errorf(result["message"].(string))
	}
	return
}

func PostCsrDecode(c *gin.Context) (result map[string]interface{}, err error) {
	err = c.Request.ParseForm()
	req, err := http.Post(api+fmt.Sprintf("tools/csr/decode?auth_key=%s", token), strings.NewReader(c.Request.PostForm.Encode()), "application/x-www-form-urlencoded; charset=UTF-8")
	err = json.Unmarshal(req, &result) //json字符串反解析为字典
	if result["success"] != true {
		err = fmt.Errorf(result["message"].(string))
	}
	return
}

func getAuth() string {
	body := strings.NewReader(url.Values{
		"user": []string{user},
		"pass": []string{pass},
	}.Encode())
	req, err := http.Post(api+"api/auth", body, "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		logger.Error(err.Error())
	}
	return string(req)
}
