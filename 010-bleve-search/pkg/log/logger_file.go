package log

import (
	"fmt"
	"net/url"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

// Sync implements zap.Sink. The remaining methods are implemented
// by the embedded *lumberjack.Logger.
func (lumberjackSink) Sync() error { return nil }

func NewFileLogger(filename string, debug ...bool) ILogger {

	zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   u.Opaque,
				MaxSize:    10, // MB
				MaxBackups: 3,  // number of backups
				MaxAge:     28, //days
				LocalTime:  true,
			},
		}, nil
	})

	var config zap.Config

	if debug != nil && debug[0] {
		config = zap.NewDevelopmentConfig()
	} else {
		//config = zap.NewProductionConfig()
		config = newProductionConfig()
	}

	config.OutputPaths = []string{fmt.Sprintf("lumberjack:%s", filename)}

	l, _ := config.Build()

	return l.Sugar()
}

func NewFileTeeLogger(filename string, debug ...bool) ILogger {

	zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   u.Opaque,
				MaxSize:    10, // MB
				MaxBackups: 3,  // number of backups
				MaxAge:     28, //days
				LocalTime:  true,
			},
		}, nil
	})

	var config zap.Config

	if debug != nil && debug[0] {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = append(config.OutputPaths, "lumberjack:"+filename)

	l, _ := config.Build()

	return l.Sugar()
}
