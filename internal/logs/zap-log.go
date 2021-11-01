package logs

import (
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
		Log: z.getLogger(),
	}
}

func (z *zapConfig) getLogger() *zapplugin.WrapLogger {
	l, _ := zap.NewProduction()
	logger := zapplugin.WrapWithContext(l)
	return logger
}