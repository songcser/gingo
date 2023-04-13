package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/utils"
	"go.uber.org/zap"
	"io"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/admin") {
			c.Next()
			return
		}

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		//开始时间
		startTime := time.Now()

		var request map[string]interface{}

		b, _ := c.Copy().GetRawData()

		_ = json.Unmarshal(b, &request)

		s, _ := json.Marshal(request)

		c.Request.Body = io.NopCloser(bytes.NewReader(b))

		//处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		//结束时间
		endTime := time.Now()

		config.GVA_LOG.Info("请求响应",
			zap.String("request_uri", c.Request.RequestURI),
			zap.String("request_method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_time", utils.TimeFormat(startTime)),
			zap.String("response_time", utils.TimeFormat(endTime)),
			zap.String("request", string(s)),
			zap.String("response", responseBody),
			zap.String("cost_time", endTime.Sub(startTime).String()),
		)
	}
}
