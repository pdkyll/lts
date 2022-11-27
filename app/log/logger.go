package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.opentelemetry.io/otel"

	"gitlab.com/m0ta/lts/app/config"
)

var Logger *zap.Logger
var Tracer = otel.Tracer("lts")

func New(cfg *config.Config) *zap.Logger {
	pe 		:= zap.NewProductionEncoderConfig()
	level 	:= zapcore.Level(cfg.LogLevel)
	//pe := zap.NewDevelopmentEncoderConfig()
    //fileEncoder := zapcore.NewJSONEncoder(pe)
	
    pe.EncodeTime 	= zapcore.ISO8601TimeEncoder
	pe.EncodeLevel 	= zapcore.CapitalColorLevelEncoder
	pe.EncodeCaller = zapcore.FullCallerEncoder
	
    consoleEncoder 	:= zapcore.NewConsoleEncoder(pe)

 	core := zapcore.NewTee(
        //zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zap.DebugLevel),
        zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
    )
	Logger = zap.New(core, zap.AddCaller())
	return Logger
}