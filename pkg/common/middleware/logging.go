package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	maxBodySize = 1024 * 10 // 10 KB
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

		// 读取请求体，但限制大小
		var bodyBytes []byte
		if r.Body != nil {
			limitedReader := io.LimitReader(r.Body, int64(maxBodySize)+1)
			bodyBytes, _ = io.ReadAll(limitedReader)

			// 检查是否超过大小限制
			if len(bodyBytes) > maxBodySize {
				bodyBytes = bodyBytes[:maxBodySize]
				bodyBytes = append(bodyBytes, []byte("... (truncated)")...)
			}

			// 重新设置请求体，因为读取后会消耗掉
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		wrappedWriter := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		// 获取真实的客户端 IP 地址
		clientIP := getClientIP(r)

		// 创建一个包含所有字段的消息字符串
		logMessage := fmt.Sprintf(
			"%s %s - Status: %d, Duration: %v, Remote Addr: %s",
			r.Method,
			r.URL.Path,
			wrappedWriter.status,
			duration,
			clientIP,
		)

		// 添加请求体到日志消息中
		if len(bodyBytes) > 0 {
			logMessage += fmt.Sprintf(", Body: %s", string(bodyBytes))
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

// getClientIP 尝试获取客户端的真实 IP 地址
func getClientIP(r *http.Request) string {
	// 检查常见的 HTTP 头部
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// 如果没有找到代理头部，就使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
