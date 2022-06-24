package main

import (
	"flag"

	"github.com/oursky/github-actions-manager/pkg/cmd"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	loglevel := zap.LevelFlag("loglevel", zap.InfoLevel, "log level")
	flag.Parse()

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(*loglevel)
	logger, _ := cfg.Build()
	defer logger.Sync()

	if *configPath == "" {
		logger.Fatal("config is required")
	}

	config, err := NewConfig(*configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	logger.Debug("loaded config", zap.Any("config", config))

	modules, err := initModules(logger, config)
	if err != nil {
		logger.Fatal("failed to init", zap.Error(err))
	}

	err = cmd.Run(logger, modules)
	if err != nil {
		logger.Fatal("fatal error occured", zap.Error(err))
	}
}
