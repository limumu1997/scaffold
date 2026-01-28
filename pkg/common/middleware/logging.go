package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"scaffold/pkg/logger"
	"strings"
	"time"
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
		logger.WithPrefix("HTTP").Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_ip", clientIP)

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

		// 默认日志参数
		logArgs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.status,
			"duration", duration,
			"remote_ip", clientIP,
		}

		// 如果请求体不为空，记录请求体
		if bodyBuffer.Len() > 0 {
			reqBody := bodyBuffer.String()
			if len(reqBody) > maxLogSize {
				reqBody = reqBody[:maxLogSize] + "... (truncated)"
			}
			logArgs = append(logArgs, "req_body", reqBody)
		}

		// 如果状态码大于等于 400，记录响应体（错误信息）
		if wrappedWriter.status >= http.StatusBadRequest && wrappedWriter.body.Len() > 0 {
			respBody := wrappedWriter.body.String()
			if len(respBody) > maxLogSize {
				respBody = respBody[:maxLogSize] + "... (truncated)"
			}
			logArgs = append(logArgs, "resp_body", respBody)
		}

		// 请求结束时打印日志
		logger.WithPrefix("HTTP").Info("request completed", logArgs...)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// 如果状态码是错误码（>= 400），我们将内容写入 buffer 以便后续打印
	if rw.status >= http.StatusBadRequest {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
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
