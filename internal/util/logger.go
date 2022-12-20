package util

import (
	"fmt"

	cli "github.com/cyp0633/drcom-go/internal/cli"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// 设置日志文件路径并初始化 logger
func SetLogPath() {
	cfg := zap.NewProductionConfig()
	// 未指定日志文件路径时，日志输出到控制台
	if cli.CLI.Log != "" {
		cfg.OutputPaths = []string{cli.CLI.Log}
	}
	logger, err := cfg.Build()
	Logger = logger
	if err != nil {
		fmt.Println("SetLogPath error:", err.Error())
		panic(err)
	}
	logger.Info("Logger initialized")
}
