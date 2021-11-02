package logs

import (
	"fmt"
	zapplugin "github.com/SkyAPM/go2sky-plugins/zap"
	"go.uber.org/zap"
)

type zapConfig struct{}

type Logger struct {
	Log *zapplugin.WrapLogger
}

var Instance *Logger

func NewZapLogger() ILog {
	return &zapConfig{}
}

func (z *zapConfig) Init() {
	Instance = &Logger{
		Log: z.GetLogger(),
	}
}

func (z *zapConfig) GetLogger() *zapplugin.WrapLogger {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Errorf("error creating zap new production")
	}
	logger := zapplugin.WrapWithContext(l)
	return logger
}