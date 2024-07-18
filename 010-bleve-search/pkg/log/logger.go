package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newExampleLogger(options ...zap.Option) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.DebugLevel)
	return zap.New(core).WithOptions(options...)
}

// Production

func newProductionEncoderConfig0() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newProductionConfig() zap.Config {

	encoderCfg := newProductionEncoderConfig()

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return cfg
}

func newDevelopementJSONConfig() zap.Config {

	encoderCfg := newProductionEncoderConfig()

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg
}

/*
func newProductionLogger0(options ...zap.Option) (*zap.Logger, error) {

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build(options...)
}
*/

func newProductionLogger(options ...zap.Option) (*zap.Logger, error) {
	return newProductionConfig().Build(options...)
}

// Standard

func newStdLogger(options ...zap.Option) (*zap.Logger, error) {

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "ts",
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "logger",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		//EncodeTime:  zapcore.EpochTimeEncoder,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build(options...)
}

func newCliLogger(options ...zap.Option) (*zap.Logger, error) {

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "ts",
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "logger",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		//EncodeTime:  zapcore.EpochTimeEncoder,
		EncodeTime: zapcore.RFC3339TimeEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build(options...)
}

func newClientLogger(options ...zap.Option) (*zap.Logger, error) {

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "ts",
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "logger",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		//EncodeTime:  zapcore.EpochTimeEncoder,
		EncodeTime: zapcore.RFC3339TimeEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build(options...)
}

func NewLogger(name ...string) ILogger {
	logger, _ := newProductionLogger() //zap.NewProduction()
	slogger := logger.Sugar()

	if len(name) > 0 {
		return slogger.Named(name[0])
	} else {
		return slogger
	}
}

func NewFullLogger(name ...string) ILogger {
	logger, _ := newStdLogger()
	slogger := logger.Sugar()

	if len(name) > 0 {
		return slogger.Named(name[0])
	} else {
		return slogger
	}
}

func NewCliLogger(name ...string) ILogger {
	logger, _ := newCliLogger()
	slogger := logger.Sugar()

	if len(name) > 0 {
		return slogger.Named(name[0])
	} else {
		return slogger
	}
}

func NewClientLogger(name ...string) ILogger {
	logger, _ := newClientLogger()
	slogger := logger.Sugar()

	if len(name) > 0 {
		return slogger.Named(name[0])
	} else {
		return slogger
	}
}

// func NewFileLogger(filename string, blah ...bool) ILogger {
// 	logger, _ := newCliLogger()
// 	slogger := logger.Sugar()

// 	logger.WithOptions()

// 	return slogger
// }

func NewDevelopment(name ...string) ILogger {
	logger, _ := zap.NewDevelopment()
	slogger := logger.Sugar()

	if len(name) > 0 {
		return slogger.Named(name[0])
	} else {
		return slogger
	}
}

func NewDevelopmentJSON(options ...zap.Option) ILogger {

	logger, _ := newDevelopementJSONConfig().Build(options...)
	slogger := logger.Sugar()

	return slogger
}

func With(logger ILogger, args ...interface{}) ILogger {

	if _log, ok := logger.(*zap.SugaredLogger); ok {
		return _log.With(args...)
	}

	return logger
}

func SetLogLevel(logger ILogger, args ...interface{}) ILogger {

	if _log, ok := logger.(*zap.SugaredLogger); ok {
		return _log.With(args...)

	}

	return logger
}

func ZapSugaredLogger(logger ILogger) *zap.SugaredLogger {

	if _log, ok := logger.(*zap.SugaredLogger); ok {
		return _log
	}

	return nil
}
