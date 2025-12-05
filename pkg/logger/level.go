package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	red     = 1
	green   = 2
	yellow  = 3
	blue    = 4
	magenta = 5
	cyan    = 6
	gray    = 7

	fmtColorEscape   = "\033[3%dm"
	colorEscapeReset = "\033[0m"
	boldEscape       = "\033[1m"
)

// ColorizeLevel applies color and bold formatting to log level strings.
func ColorizeLevel(level zapcore.Level, levelStr string) string {
	var color int
	switch level {
	case zapcore.DebugLevel:
		color = gray
	case zapcore.InfoLevel:
		color = green
	case zapcore.WarnLevel:
		color = yellow
	case zapcore.ErrorLevel:
		color = magenta
	case zapcore.FatalLevel:
		color = red
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		color = cyan
	default:
		color = gray
	}
	return fmt.Sprintf(fmtColorEscape+boldEscape+levelStr+colorEscapeReset, color)
}

// String returns a lower-case ASCII representation of the log level.
func LevelString(l zapcore.Level) string {
	switch l {
	case zapcore.DebugLevel:
		return "[debug]"
	case zapcore.InfoLevel:
		return "[info]"
	case zapcore.WarnLevel:
		return "[warn]"
	case zapcore.ErrorLevel:
		return "[error]"
	case zapcore.DPanicLevel:
		return "[dpanic]"
	case zapcore.PanicLevel:
		return "[panic]"
	case zapcore.FatalLevel:
		return "[fatal]"
	default:
		return fmt.Sprintf("[Level(%d)]", l)
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func LevelCapitalString(l zapcore.Level) string {
	switch l {
	case zapcore.DebugLevel:
		return "[DEBUG]"
	case zapcore.InfoLevel:
		return "[INFO]"
	case zapcore.WarnLevel:
		return "[WARN]"
	case zapcore.ErrorLevel:
		return "[ERROR]"
	case zapcore.DPanicLevel:
		return "[DPANIC]"
	case zapcore.PanicLevel:
		return "[PANIC]"
	case zapcore.FatalLevel:
		return "[FATAL]"
	default:
		return fmt.Sprintf("[LEVEL(%d)]", l)
	}
}
