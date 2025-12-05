package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// LowercaseLevelEncoder serializes a Level to a lowercase string with color.
func LowercaseLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorizeLevel(l, LevelString(l)))
}

// CapitalLevelEncoder serializes a Level to an all-caps string with color.
func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorizeLevel(l, LevelCapitalString(l)))
}

func ShortCallerEncoderCustom(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(colorizeCaller(caller.TrimmedPath()))
}

func colorizeCaller(callerPath string) string {
	return fmt.Sprintf("at "+fmtColorEscape+boldEscape+callerPath+colorEscapeReset, blue)
}
