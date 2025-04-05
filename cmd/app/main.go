package main

import (
	"log/slog"
	"os"
	"template/internal/config"
	"template/internal/logging"
)

func main() {
	cfg := config.MustLoad()
	mustSetupLogging(cfg.Env, cfg.LogPath)
	slog.Info("Logging started")
}

func mustSetupLogging(env, logPath string) {
	var logger *slog.Logger

	switch env {
	case "dev":
		logger = slog.New(logging.NewTerminalHandler(os.Stdout, slog.LevelDebug))
	case "prod":
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			panic(err)
		}
		logger = slog.New(slog.NewJSONHandler(f, nil))
	}

	if logger == nil {
		panic("Logging setup failed")
	}
	slog.SetDefault(logger)
}
