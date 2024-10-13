package base

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// LogFormatter is used to format log entry.
type LogFormatter struct{}

// Format formats a given log entry, returns byte slice and error.
func (c *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	if level == "WARNING" {
		level = "WARN"
	}
	if len(level) < 5 {
		level = strings.Repeat(" ", 5-len(level)) + level
	}

	return []byte(fmt.Sprintf(
		"[%s #%d] %s -- : %s\n",
		time.Now().Format("2006-01-02T15:04:05.000Z"),
		os.Getpid(),
		level,
		entry.Message)), nil
}

// SetOutput set the destination for the log output
func SetOutput(out io.Writer) {
	Log.Out = out
}

// CheckLevel checks whether the log level is valid.
func CheckLevel(level string) error {
	if _, err := logrus.ParseLevel(level); err != nil {
		return fmt.Errorf(`log level not valid: "%s"`, level)
	}
	return nil
}

// GetLevel get the log level string.
func GetLevel() string {
	return Log.Level.String()
}

// SetLevel sets the log level. Valid levels are "debug", "info", "warn", "error", and "fatal".
func SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		Fatal(fmt.Sprintf(`log level not valid: "%s"`, level))
	}
	Log.Level = lvl
}

// Debug logs a message with severity DEBUG.
func Debug(format string, v ...interface{}) {
	output(Log.Debug, format, v...)
}

// Info logs a message with severity INFO.
func Info(format string, v ...interface{}) {
	output(Log.Info, format, v...)
}

// Warn logs a message with severity WARN.
func Warn(format string, v ...interface{}) {
	output(Log.Warn, format, v...)
}

// Error logs a message with severity ERROR.
func Error(format string, v ...interface{}) {
	output(Log.Error, format, v...)
}

// Fatal logs a message with severity ERROR followed by a call to os.Exit().
func Fatal(format string, v ...interface{}) {
	output(Log.Fatal, format, v...)
}

func output(origin func(...interface{}), format string, v ...interface{}) {
	if len(v) > 0 {
		origin(fmt.Sprintf(format, v...))
	} else {
		origin(format)
	}
}

func NewLogger(dir, name string, level int) *logrus.Logger {
	NewLog := logrus.New()
	NewLog.Formatter = &LogFormatter{}
	NewLog.Out = os.Stderr
	switch level {
	case 0:
		NewLog.SetLevel(logrus.PanicLevel)
	case 1:
		NewLog.SetLevel(logrus.FatalLevel)
	case 2:
		NewLog.SetLevel(logrus.ErrorLevel)
	case 3:
		NewLog.SetLevel(logrus.WarnLevel)
	case 4:
		NewLog.SetLevel(logrus.InfoLevel)
	case 5:
		NewLog.SetLevel(logrus.DebugLevel)
	case 6:
		NewLog.SetLevel(logrus.TraceLevel)
	default:
		NewLog.SetLevel(logrus.WarnLevel)
	}
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v", dir, name),
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
	}
	NewLog.SetOutput(io.MultiWriter(logger, os.Stdout))
	return NewLog
}

func InitLogger(dir, name string, level int) {
	Log = logrus.New()
	Log.Formatter = &LogFormatter{}
	Log.Out = os.Stderr
	switch level {
	case 0:
		Log.SetLevel(logrus.PanicLevel)
	case 1:
		Log.SetLevel(logrus.FatalLevel)
	case 2:
		Log.SetLevel(logrus.ErrorLevel)
	case 3:
		Log.SetLevel(logrus.WarnLevel)
	case 4:
		Log.SetLevel(logrus.InfoLevel)
	case 5:
		Log.SetLevel(logrus.DebugLevel)
	case 6:
		Log.SetLevel(logrus.TraceLevel)
	default:
		Log.SetLevel(logrus.WarnLevel)
	}
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v", dir, name),
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
	}
	Log.SetOutput(io.MultiWriter(logger, os.Stdout))
}
