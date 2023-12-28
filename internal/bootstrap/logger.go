package bootstrap

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type MyFormatter struct{}

type MyErrorHook struct{ errorLogger *lumberjack.Logger }

func InitLog() {
	executable, _ := os.Executable()
	res, _ := filepath.EvalSymlinks(filepath.Dir(executable))
	appLogger := &lumberjack.Logger{
		Filename:   filepath.Join(res, "logs/app.log"),
		MaxSize:    20, // megabytes
		MaxBackups: 1024,
		MaxAge:     512, //days
	}
	errorLogger := &lumberjack.Logger{
		Filename:   filepath.Join(res, "logs/error.log"),
		MaxSize:    20, // megabytes
		MaxBackups: 1024,
		MaxAge:     512, //days
	}

	h := &MyErrorHook{errorLogger}
	writers := []io.Writer{
		appLogger,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logrus.SetFormatter(&MyFormatter{})
	logrus.AddHook(h)
	logrus.SetOutput(fileAndStdoutWriter)
	// Will logrus anything that is info or above (warn, error, fatal, panic). Default.
	logrus.SetLevel(logrus.TraceLevel)
}

// Implement Formatter interface
// Format renders a single logrus entry
func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatEntry := formatEntry(entry)
	return formatEntry, nil
}

func formatEntry(entry *logrus.Entry) []byte {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02T15:04:05.000")
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(timestamp)
	sb.WriteString("]")
	sb.WriteString(" ")
	sb.WriteString("[")
	sb.WriteString(entry.Level.String())
	sb.WriteString("]")
	sb.WriteString(" ")
	sb.WriteString(entry.Message)
	sb.WriteString("\r\n")
	fileVal := sb.String()
	b.WriteString(fileVal)
	return b.Bytes()
}

// Implement Hook interface
// Only the logrus level of interest is required
func (h *MyErrorHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (h *MyErrorHook) Fire(entry *logrus.Entry) error {
	formatEntry := formatEntry(entry)
	if _, err := h.errorLogger.Write(formatEntry); err != nil {
		return err
	}
	return nil
}
