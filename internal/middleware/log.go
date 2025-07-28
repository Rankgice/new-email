package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogMiddleware 日志中间件
func LogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// RequestLogMiddleware 请求日志中间件（详细版）
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 记录日志
		log.Printf("[%s] %s %s %d %v %s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		// 如果是错误请求，记录详细信息
		if c.Writer.Status() >= 400 {
			log.Printf("Request Body: %s", string(requestBody))
			log.Printf("Response Body: %s", writer.body.String())
		}
	}
}

// responseWriter 响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
