package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (c Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c Config) ZapLoggerLevel() zap.AtomicLevel {
	var level zapcore.Level

	switch c.Logging.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	return zap.NewAtomicLevelAt(level)
}
