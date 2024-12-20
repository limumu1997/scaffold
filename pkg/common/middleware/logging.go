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
	maxLogSize = 256 // byte
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

		// 获取客户端 IP
		clientIP := getClientIP(r)

		// 请求开始时打印日志
		logrus.WithField("prefix", "HTTP").Infof(
			"Request Started - %s %s from %s",
			r.Method,
			r.URL.Path,
			clientIP,
		)

		start := time.Now()

		var bodyBuffer bytes.Buffer
		var bodyReader io.Reader = r.Body

		// 如果有请求体，设置 TeeReader
		if r.Body != nil {
			bodyReader = io.TeeReader(r.Body, &bodyBuffer)
		}

		// 创建一个新的 request，复制原始请求的所有字段
		newRequest := *r
		// 设置新的 Body
		newRequest.Body = io.NopCloser(bodyReader)

		wrappedWriter := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrappedWriter, &newRequest)

		duration := time.Since(start)

		// 请求结束时的日志消息
		logMessage := fmt.Sprintf(
			"Request Completed - %s %s - status: %d, duration: %v, remote addr: %s",
			r.Method,
			r.URL.Path,
			wrappedWriter.status,
			duration,
			clientIP,
		)

		// 添加请求体到日志消息中，但限制日志大小
		if bodyBuffer.Len() > 0 {
			logBody := bodyBuffer.Bytes()
			if len(logBody) > maxLogSize {
				logBody = logBody[:maxLogSize]
				logBody = append(logBody, []byte("... (truncated)")...)
			}
			logMessage += fmt.Sprintf(", body: %s", string(logBody))
		}

		// 请求结束时打印日志
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
