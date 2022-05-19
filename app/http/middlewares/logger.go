// Package middlewares 存放系统中间件
package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"api/app/models/opera_log"
	"api/pkg/auth"
	"api/pkg/config"
	"api/pkg/helpers"
	"api/pkg/ip"
	"api/pkg/logger"
	"api/pkg/queue"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取 response 内容
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 设置开始时间
		start := time.Now()
		c.Next()

		requestBodyLog := ""
		bodyLog := ""

		// 开始记录日志的逻辑
		cost := time.Since(start)
		responStatus := c.Writer.Status()

		logFields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("request", c.Request.Method+" "+c.Request.URL.String()),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("time", helpers.MicrosecondsStr(cost)),
		}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			// 请求的内容
			logFields = append(logFields, zap.String("Request Body", string(requestBody)))

			// 响应的内容
			logFields = append(logFields, zap.String("Response Body", w.body.String()))

			requestBodyLog = string(requestBody)
			bodyLog = w.body.String()
		}

		if responStatus > 400 && responStatus <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
			logger.Warn("HTTP Warning "+cast.ToString(responStatus), logFields...)
		} else if responStatus >= 500 && responStatus <= 599 {
			// 除了内部错误，记录 error
			logger.Error("HTTP Error "+cast.ToString(responStatus), logFields...)
		} else {
			logger.Debug("HTTP Access Log", logFields...)
		}

		if c.Request.Method != "OPTIONS" && config.GetBool("log.enabled_db") && responStatus != 404 && c.FullPath() != "" {
			SetDBOperaLog(c, responStatus, c.Request.RequestURI, c.Request.Method, helpers.MicrosecondsStr(cost), requestBodyLog, bodyLog, c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}

func statFromHandlerName(handlerName string) (moduleName, methodName string) {
	lastDotIndex := strings.LastIndex(handlerName, ".")
	if lastDotIndex > 0 {
		moduleName = handlerName[0:lastDotIndex]
		methodName = handlerName[lastDotIndex:]

		moduleName = strings.Replace(moduleName, ".", "", -1)
	} else {
		moduleName = handlerName
	}
	return moduleName, methodName
}

// SetDBOperaLog 写入操作日志表
func SetDBOperaLog(c *gin.Context, statusCode int, reqUri string, reqMethod string, latencyTime string, body string, result string, error string) {
	status := "1"
	if statusCode >= 200 && statusCode < 300 {
		status = "2"
	}
	moduleName, methodName := statFromHandlerName(c.HandlerName())
	log := opera_log.OperaLog{
		Title:         c.FullPath(),
		BusinessType:  moduleName,
		Method:        methodName,
		OperaUrl:      reqUri,
		Status:        status,
		OperaIp:       ip.GetClientIP(c),
		OperaLocation: ip.GetLocation(ip.GetClientIP(c)),
		OperaName:     auth.CurrentName(c),
		RequestMethod: reqMethod,
		OperaParam:    body,
		OperaTime:     time.Now(),
		JsonResult:    result,
		Remark:        error,
		LatencyTime:   latencyTime,
		UserAgent:     c.Request.UserAgent(),
		OperatorType:  strconv.Itoa(statusCode) + " " + c.Request.Method + " " + c.Request.URL.String(),
	}
	if config.GetBool("redis.enable") {
		if err := queue.Queue.Producers(log, "oplog"); err != nil {
			logger.Error(fmt.Sprintf("publish data error:%s", err.Error()))
		}
	} else {
		queue.Producers(log, "oplog")
	}
}