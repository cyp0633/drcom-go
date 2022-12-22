package util

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// 初始化 logger
func SetupLog() {
	encConf := zap.NewProductionEncoderConfig()
	encConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var enc zapcore.Encoder
	var writeSyncer zapcore.WriteSyncer
	if CLI.Log == "" {
		enc = zapcore.NewConsoleEncoder(encConf)
		writeSyncer = zapcore.Lock(os.Stdout)
	} else {
		file, err := os.OpenFile(CLI.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("SetLogPath error:", err.Error())
			panic(err)
		}
		writeSyncer = zapcore.AddSync(file)
		enc = zapcore.NewJSONEncoder(encConf)
	}
	var core zapcore.Core
	if CLI.Debug {
		core = zapcore.NewCore(enc, writeSyncer, zapcore.DebugLevel)
	} else {
		core = zapcore.NewCore(enc, writeSyncer, zapcore.InfoLevel)
	}
	Logger = zap.New(core, zap.AddCaller())
}
