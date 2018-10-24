package zlog

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := New()
	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	logger.Fatal("This is an Fatal")
	logger.Panic("This is an Panic")
}

func TestGlobalLogger(t *testing.T) {
	Zlog.Debug("This is a DEBUG message")
	Zlog.Info("This is an INFO message")
	Zlog.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
	Zlog.Warn("This is a WARN message")
	Zlog.Error("This is an ERROR message")
	Zlog.Fatal("This is an Fatal")
	Zlog.Panic("This is an Panic")
}
