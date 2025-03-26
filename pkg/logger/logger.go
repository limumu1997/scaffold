package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MyHandler 自定义 slog 处理器，实现分级日志输出
type MyHandler struct {
	appHandler   slog.Handler
	errorHandler slog.Handler
}

// RotateFileWriter 是一个支持日志轮转的 io.Writer 实现
type RotateFileWriter struct {
	filename   string
	maxSize    int64 // 单位：MB
	maxBackups int
	maxAge     int // 单位：天
	size       int64
	file       *os.File
}

// NewRotateFileWriter 创建一个新的日志轮转写入器
func NewRotateFileWriter(filename string, maxSize int, maxBackups, maxAge int) (*RotateFileWriter, error) {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 打开日志文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// 获取当前文件大小
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	w := &RotateFileWriter{
		filename:   filename,
		maxSize:    int64(maxSize) * 1024 * 1024, // 转换为字节
		maxBackups: maxBackups,
		maxAge:     maxAge,
		size:       info.Size(),
		file:       file,
	}

	// 清理过期日志
	go w.cleanOldLogs()

	return w, nil
}

// Write 实现 io.Writer 接口
func (w *RotateFileWriter) Write(p []byte) (n int, err error) {
	// 检查是否需要轮转
	if w.size+int64(len(p)) >= w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

// Close 关闭文件
func (w *RotateFileWriter) Close() error {
	return w.file.Close()
}

// rotate 执行日志轮转
func (w *RotateFileWriter) rotate() error {
	// 关闭当前文件
	if err := w.file.Close(); err != nil {
		return err
	}

	// 生成时间戳后缀
	timestamp := time.Now().Format("2006-01-02-150405")

	// 备份当前日志文件
	backupName := w.filename + "." + timestamp
	if err := os.Rename(w.filename, backupName); err != nil {
		return err
	}

	// 打开新的日志文件
	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.size = 0

	// 异步清理过期日志
	go w.cleanOldLogs()

	return nil
}

// cleanOldLogs 清理过期的日志文件
func (w *RotateFileWriter) cleanOldLogs() {
	dir := filepath.Dir(w.filename)
	base := filepath.Base(w.filename)

	// 查找所有备份文件
	matches, err := filepath.Glob(filepath.Join(dir, base+".*"))
	if err != nil {
		return
	}

	// 检查数量限制
	if len(matches) <= w.maxBackups {
		return
	}

	// 按修改时间排序
	type backupFile struct {
		name string
		time time.Time
	}

	backups := make([]backupFile, 0, len(matches))
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			continue
		}

		// 检查文件年龄
		if w.maxAge > 0 {
			cutoff := time.Now().AddDate(0, 0, -w.maxAge)
			if info.ModTime().Before(cutoff) {
				os.Remove(m)
				continue
			}
		}

		backups = append(backups, backupFile{name: m, time: info.ModTime()})
	}

	// 根据时间排序
	for i := range backups {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].time.After(backups[j].time) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}

	// 删除多余的备份
	for i := range len(backups)-w.maxBackups {
		os.Remove(backups[i].name)
	}
}

// 自定义格式化处理器
type customHandler struct {
	w          io.Writer
	level      slog.Level
	withSource bool
}

func (h *customHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("2006-01-02T15:04:05.000")
	levelStr := r.Level.String()

	var prefix string
	// 尝试获取前缀
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "prefix" {
			prefix = a.Value.String()
		}
		return true
	})

	var builder strings.Builder
	builder.WriteString("[")
	builder.WriteString(timestamp)
	builder.WriteString("] [")
	builder.WriteString(strings.ToUpper(levelStr))
	builder.WriteString("]")

	if prefix != "" {
		builder.WriteString(" [")
		builder.WriteString(prefix)
		builder.WriteString("]")
	}

	builder.WriteString(" ")
	builder.WriteString(r.Message)

	// 添加其他属性（排除已处理的前缀）
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "prefix" {
			builder.WriteString(" ")
			builder.WriteString(a.Key)
			builder.WriteString("=")

			// 如果值包含空格，加上引号
			val := a.Value.String()
			if strings.ContainsAny(val, " \t\n") {
				builder.WriteString("\"")
				builder.WriteString(val)
				builder.WriteString("\"")
			} else {
				builder.WriteString(val)
			}
		}
		return true
	})

	builder.WriteString("\r\n")

	_, err := h.w.Write([]byte(builder.String()))
	return err
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 创建一个新的处理器，属性会在日志记录时处理
	return h
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	// 简单实现，忽略分组
	return h
}

// 创建自定义处理器
func newTextHandler(w io.Writer) slog.Handler {
	return &customHandler{
		w:          w,
		level:      slog.LevelDebug,
		withSource: false,
	}
}

// Handle 实现 slog.Handler 接口
func (h *MyHandler) Handle(ctx context.Context, r slog.Record) error {
	// 所有日志都写入应用日志
	err := h.appHandler.Handle(ctx, r)

	// 错误级别以上的日志也写入错误日志
	if r.Level >= slog.LevelError {
		errErr := h.errorHandler.Handle(ctx, r)
		if err == nil {
			err = errErr
		}
	}

	return err
}

// WithAttrs 实现 slog.Handler 接口
func (h *MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MyHandler{
		appHandler:   h.appHandler.WithAttrs(attrs),
		errorHandler: h.errorHandler.WithAttrs(attrs),
	}
}

// WithGroup 实现 slog.Handler 接口
func (h *MyHandler) WithGroup(name string) slog.Handler {
	return &MyHandler{
		appHandler:   h.appHandler.WithGroup(name),
		errorHandler: h.errorHandler.WithGroup(name),
	}
}

// Enabled 实现 slog.Handler 接口
func (h *MyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.appHandler.Enabled(ctx, level)
}

// InitMyLog 初始化日志系统
func InitMyLog() error {
	// 获取可执行文件所在目录
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	res, err := filepath.EvalSymlinks(filepath.Dir(executable))
	if err != nil {
		return err
	}

	// 创建应用日志写入器
	appLogWriter, err := NewRotateFileWriter(
		filepath.Join(res, "logs/app.log"),
		20,   // 20MB
		1024, // 最多1024个备份
		512,  // 保留512天
	)
	if err != nil {
		return err
	}

	// 创建错误日志写入器
	errorLogWriter, err := NewRotateFileWriter(
		filepath.Join(res, "logs/error.log"),
		20,   // 20MB
		1024, // 最多1024个备份
		512,  // 保留512天
	)
	if err != nil {
		return err
	}

	// 创建多输出写入器（应用日志同时输出到标准输出）
	multiWriter := io.MultiWriter(appLogWriter, os.Stdout)

	// 创建自定义处理器
	handler := &MyHandler{
		appHandler:   newTextHandler(multiWriter),
		errorHandler: newTextHandler(errorLogWriter),
	}

	// 设置全局日志处理器
	slog.SetDefault(slog.New(handler))

	return nil
}

// 添加前缀功能的便捷方法
func WithPrefix(prefix string) *slog.Logger {
	return slog.With("prefix", prefix)
}
