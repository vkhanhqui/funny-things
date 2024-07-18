package log

import (
	"io"
	"net/url"

	"go.uber.org/zap"
)

type pipeSink struct {
	*io.PipeWriter
}

// Sync implements zap.Sink. The remaining methods are implemented
// by the embedded *io.PipeWriter.
func (pipeSink) Sync() error { return nil }

func NewPipeLogger(logr ILogger, r *io.PipeWriter, name string, debug ...bool) ILogger {

	zap.RegisterSink("pipe", func(u *url.URL) (zap.Sink, error) {
		return pipeSink{r}, nil
	})

	var config zap.Config

	if debug != nil && debug[0] {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = append(config.OutputPaths, "pipe:"+name)

	l, _ := config.Build()

	return l.Sugar()
}
