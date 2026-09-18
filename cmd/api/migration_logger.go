package main

import (
	"fmt"
	"log/slog"
	"strings"
)

type migrationLogger struct {
	logger *slog.Logger
}

func (l migrationLogger) Printf(format string, args ...any) {
	l.logger.Info(
		"database migration",
		"message", strings.TrimSuffix(fmt.Sprintf(format, args...), "\n"),
	)
}

func (migrationLogger) Verbose() bool {
	return false
}
