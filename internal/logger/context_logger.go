package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey int

const (
	loggerID loggerKey = iota
)

type contextLogger struct {
	zapLogger *zap.SugaredLogger
}

func NewLogger() *contextLogger {
	zapLogger, _ := zap.NewProduction()

	return &contextLogger{
		zapLogger: zapLogger.Sugar(),
	}
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerID, NewLogger())
}

func LoggerFromContext(ctx context.Context) *contextLogger {
	if ctx == nil || ctx.Value(loggerID) == nil {
		return NewLogger()
	}
	ctxLogger := ctx.Value(loggerID).(*contextLogger)
	return ctxLogger
}

func (c *contextLogger) Info(args ...interface{}) {
	c.zapLogger.Info(args...)
}

func (c *contextLogger) Infof(template string, args ...interface{}) {
	c.zapLogger.Infof(template, args...)
}

func (c *contextLogger) Debug(args ...interface{}) {
	c.zapLogger.Debug(args...)
}

func (c *contextLogger) Debugf(template string, args ...interface{}) {
	c.zapLogger.Debugf(template, args...)
}

func (c *contextLogger) Error(args ...interface{}) {
	c.zapLogger.Error(args...)
}

func (c *contextLogger) Errorf(template string, args ...interface{}) {
	c.zapLogger.Errorf(template, args...)
}

func (c *contextLogger) Fatal(args ...interface{}) {
	c.zapLogger.Fatal(args...)
}

func (c *contextLogger) Fatalf(template string, args ...interface{}) {
	c.zapLogger.Fatalf(template, args...)
}

func (c *contextLogger) Panic(args ...interface{}) {
	c.zapLogger.Panic(args...)
}

func (c *contextLogger) Panicf(template string, args ...interface{}) {
	c.zapLogger.Panicf(template, args...)
}
