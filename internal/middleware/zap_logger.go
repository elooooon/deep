package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// ZapLogger 创建基于 zap 的访问日志中间件（与Gin版保持相同格式）
func ZapLogger() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 记录请求开始时间
			start := time.Now()
			path := r.URL.Path
			clientIP := getRealClientIP(r)

			// 包装ResponseWriter以获取状态码
			lrw := &loggedResponseWriter{ResponseWriter: w}

			// 处理请求
			next(lrw, r)
			// 计算耗时
			cost := time.Since(start)

			// 构建与Gin完全相同的日志格式
			logMsg := fmt.Sprintf("| %s | %s | %-15s | %-20s | %-12s | %s",
				colorStatus(lrw.status),
				colorMethod(r.Method),
				clientIP,
				path,
				cost.String(),
				simplifyUA(r.UserAgent()),
			)

			// 使用logx输出（最终会被zap处理）
			logx.Info(logMsg)
		}
	}
}

// NotFoundHandler 返回一个自定义的404处理函数（新增函数）
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录404日志
		logMsg := fmt.Sprintf("| %s | %s | %-15s | %-20s | %-12s | %s",
			colorStatus(http.StatusNotFound),
			colorMethod(r.Method),
			getRealClientIP(r),
			r.URL.Path,
			"0s", // 404立即返回，没有耗时
			simplifyUA(r.UserAgent()),
		)
		logx.Info(logMsg)

		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// loggedResponseWriter 用于捕获状态码
type loggedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// getRealClientIP 获取真实客户端IP（处理代理情况）
func getRealClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func simplifyUA(ua string) string {
	if strings.Contains(ua, "Chrome") {
		return "Chrome"
	}
	if strings.Contains(ua, "Firefox") {
		return "Firefox"
	}
	if strings.Contains(ua, "Safari") {
		return "Safari"
	}
	return "Other"
}

// ====================== 颜色格式化函数（与Gin版完全一致）======================
func colorStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return fmt.Sprintf("\033[32m%3d\033[0m", status) // 绿色
	case status >= 300 && status < 400:
		return fmt.Sprintf("\033[36m%3d\033[0m", status) // 青色
	case status >= 400 && status < 500:
		return fmt.Sprintf("\033[33m%3d\033[0m", status) // 黄色
	case status >= 500:
		return fmt.Sprintf("\033[31m%3d\033[0m", status) // 红色
	default:
		return fmt.Sprintf("\033[37m%3d\033[0m", status) // 灰色
	}
}

func colorMethod(method string) string {
	switch method {
	case "GET":
		return "\033[36mGET\033[0m" // 青色
	case "POST":
		return "\033[34mPOST\033[0m" // 蓝色
	case "PUT":
		return "\033[35mPUT\033[0m" // 紫色
	case "DELETE":
		return "\033[31mDELETE\033[0m" // 红色
	default:
		return fmt.Sprintf("\033[37m%s\033[0m", method) // 灰色
	}
}
