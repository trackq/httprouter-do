package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ZapLogger struct {
	next   http.Handler
	logger *zap.Logger
}

func NewZapLogger(next http.Handler) *ZapLogger {
	logger, err := zap.NewDevelopment() // or NewProduction, or NewDevelopment

	if err != nil {
		fmt.Println(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	return &ZapLogger{next: next, logger: logger}
}

func (l *ZapLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()

	l.next.ServeHTTP(w, r)

	l.logger.Info(fmt.Sprintf("%s request", r.Method),
		zap.String("method", r.Method),
		zap.String("path", r.URL.EscapedPath()),
		zap.Duration("duration", time.Since(startedAt)),
	)
}
