package logger

import (
	"net/http"

	"go.uber.org/zap"
)

var middlewareLoggerFn func(http.Handler) http.Handler

type field map[string]string

func MiddlewareLogger(h http.Handler) http.Handler {
	return middlewareLoggerFn(h)
}

type middlewareLogger struct {
	logger *zap.SugaredLogger
}

func init() {
	zapLogger, _ := zap.NewProduction()

	mwl := middlewareLogger{
		logger: zapLogger.Sugar(),
	}

	middlewareLoggerFn = middlewareHandler(mwl)
}

func middlewareHandler(mwl middlewareLogger) func(next http.Handler) http.Handler {
	return func(hh http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			mwl.logger.Info("http handler",
				field{
					"method": r.Method,
					"URI":    r.RequestURI,
				})
			hh.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
