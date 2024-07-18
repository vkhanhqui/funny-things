package log

type noLogger struct{}

func NoILogger(name ...string) ILogger {
	return noLogger{}
}

func (noLogger) DPanic(args ...interface{})                   {}
func (noLogger) DPanicf(template string, args ...interface{}) {}
func (noLogger) DPanicw(template string, args ...interface{}) {}
func (noLogger) Debug(args ...interface{})                    {}
func (noLogger) Debugf(template string, args ...interface{})  {}
func (noLogger) Debugw(template string, args ...interface{})  {}
func (noLogger) Error(args ...interface{})                    {}
func (noLogger) Errorf(template string, args ...interface{})  {}
func (noLogger) Errorw(template string, args ...interface{})  {}
func (noLogger) Fatal(args ...interface{})                    {}
func (noLogger) Fatalf(template string, args ...interface{})  {}
func (noLogger) Fatalw(template string, args ...interface{})  {}
func (noLogger) Info(args ...interface{})                     {}
func (noLogger) Infof(template string, args ...interface{})   {}
func (noLogger) Infow(template string, args ...interface{})   {}
func (noLogger) Panic(args ...interface{})                    {}
func (noLogger) Panicf(template string, args ...interface{})  {}
func (noLogger) Panicw(template string, args ...interface{})  {}
func (noLogger) Warn(args ...interface{})                     {}
func (noLogger) Warnf(template string, args ...interface{})   {}
func (noLogger) Warnw(template string, args ...interface{})   {}
func (noLogger) Sync() error                                  { return nil }
