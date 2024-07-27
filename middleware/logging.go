// middleware/logging.go

package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 跳过前端资源请求
		if strings.EqualFold(r.URL.Path, "/") ||
			strings.HasPrefix(r.URL.Path, "/static") ||
			strings.HasPrefix(r.URL.Path, "/favicon.ico") {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		var params string
		var body string

		switch r.Method {
		case "GET":
			// 对于 GET 请求，直接使用 URL 中的查询参数
			params = r.URL.RawQuery
		case "POST", "PUT", "PATCH":
			// 对于 POST, PUT, PATCH 请求，解析表单数据
			err := r.ParseForm()
			if err != nil {
				logrus.Error("Error parsing form: ", err)
			}
			params = r.Form.Encode()

			// 读取请求体
			if r.Body != nil {
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				body = string(bodyBytes)
			}
		}

		// 创建一个包含所有字段的消息字符串
		logMessage := fmt.Sprintf(
			"HTTP %s %s - Status: %d, Duration: %v, Remote Addr: %s, Params: %s",
			r.Method,
			r.URL.Path,
			wrappedWriter.status,
			duration,
			r.RemoteAddr,
			params,
		)

		// 如果有请求体，添加到日志消息中
		if body != "" {
			logMessage += fmt.Sprintf(", Body: %s", body)
		}

		// 使用 WithField 来添加一个前缀，这样可以与您的自定义格式器配合使用
		logrus.WithField("prefix", "HTTP").Info(logMessage)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
